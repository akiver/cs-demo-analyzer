package demo

import (
	"errors"
	"fmt"
	"hash/crc64"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/akiver/cs-demo-analyzer/internal/bitread"
	"github.com/akiver/cs-demo-analyzer/internal/filepath"
	"github.com/akiver/cs-demo-analyzer/pkg/api/constants"
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/msgs2"
	"google.golang.org/protobuf/proto"
)

type Demo struct {
	FilePath string
	FileName string
	Type     constants.DemoType
	// ! Checksums are not generated reading the whole .dem file but only information for the "header" because it would be slow.
	Checksum        string
	Filestamp       string
	Date            time.Time
	ServerName      string
	ClientName      string
	MapName         string
	NetworkProtocol int
	BuildNumber     int           // Source 2 demos only
	TickCount       int           // Not available for Source 2 demos, it's updated during parsing
	TickRate        float64       // Not available for Source 2 demos, it's updated during parsing
	FrameRate       float64       // Not available for Source 2 demos, it's updated during parsing
	Duration        time.Duration // Not available for Source 2 demos, it's updated during parsing
}

var faceItDemoNameRegex = regexp.MustCompile(`/[0-9]+_team[a-z0-9-]+-Team[a-z0-9-]+_de_[a-z0-9]+\.dem/`)
var ebotDemoNameRegex = regexp.MustCompile(`/([0-9]*)_(.*?)-(.*?)_(.*?)(.dem)/`)

func GetDemoFromPath(demoPath string) (*Demo, error) {
	file, err := os.Open(demoPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf(fmt.Sprintf("demo file \"%s\" not found", demoPath))
		}
		return nil, err
	}
	defer file.Close()

	stats, err := file.Stat()
	if err != nil {
		return nil, err
	}

	br := bitread.NewLargeBitReader(file)
	filestamp := br.ReadCString(8)
	isSource2 := filestamp == "PBDEMS2"

	var checksum string
	var mapName string
	var serverName string
	var clientName string
	var tickCount int
	var duration time.Duration
	var frameCount int
	var frameRate float64
	var tickRate float64
	var networkProtocol int
	var buildNumber int
	demoType := constants.DemoTypeGOTV
	if isSource2 {
		br.ReadBytes(8)
		msgType := br.ReadVarInt32()
		// The first proto message should always be EDemoCommands.DEM_FileHeader
		if msgType != 1 {
			return nil, fmt.Errorf("unexpected first proto message type")
		}

		br.ReadVarInt32() // tick
		size := br.ReadVarInt32()
		bytes := br.ReadBytes(int(size))
		var msg msgs2.CDemoFileHeader
		err := proto.Unmarshal(bytes, &msg)
		if err != nil {
			return nil, fmt.Errorf("failed to parse CMsgPlayerInfo msg")
		}

		mapName = msg.GetMapName()
		serverName = msg.GetServerName()
		clientName = msg.GetClientName()
		data := fmt.Sprintf("%s%s%s%d%d%s%s%d", mapName, serverName, clientName, msg.GetNetworkProtocol(), msg.GetBuildNum(), msg.GetDemoVersionGuid(), msg.GetDemoVersionName(), stats.Size())
		checksum = strconv.FormatUint(crc64.Checksum([]byte(data), crc64.MakeTable(crc64.ECMA)), 16)
		networkProtocol = int(msg.GetNetworkProtocol())
		buildNumber = int(msg.GetBuildNum())
		game := msg.GetGame()
		if game != "" {
			demoType = constants.DemoTypePOV
		}
	} else {
		br.ReadSignedInt(32) // demo protocol
		networkProtocol = br.ReadSignedInt(32)
		serverName = br.ReadCString(260)
		clientName = br.ReadCString(260)
		mapName = br.ReadCString(260)
		br.ReadCString(260) // game directory
		duration = time.Duration(br.ReadFloat() * float32(time.Second))
		tickCount = br.ReadSignedInt(32)
		frameCount = br.ReadSignedInt(32)
		signonLength := br.ReadSignedInt(32)

		if duration.Seconds() > 0 {
			frameRate = float64(frameCount) / duration.Seconds()
			tickRate = float64(tickCount) / duration.Seconds()
		}

		data := fmt.Sprintf("%s%s%s%d%d%d%d%d", mapName, serverName, clientName, frameCount, tickCount, networkProtocol, signonLength, stats.Size())
		checksum = strconv.FormatUint(crc64.Checksum([]byte(data), crc64.MakeTable(crc64.ECMA)), 16)
	}

	return &Demo{
		FilePath:        filepath.GetAbsoluteFilePath(demoPath),
		FileName:        filepath.GetFileNameWithoutExtension(demoPath),
		Type:            demoType,
		Filestamp:       filestamp,
		Checksum:        checksum,
		Date:            stats.ModTime(),
		ServerName:      serverName,
		ClientName:      clientName,
		Duration:        duration,
		MapName:         getMapNameFromHeaderMapName(mapName),
		TickCount:       tickCount,
		TickRate:        tickRate,
		FrameRate:       frameRate,
		NetworkProtocol: networkProtocol,
		BuildNumber:     buildNumber,
	}, nil
}

func (demo *Demo) IsSource2() bool {
	return demo.Filestamp == "PBDEMS2"
}

func GetDemoSource(demo *Demo) constants.DemoSource {
	demoName := strings.ToLower(demo.FileName)
	serverName := strings.ToLower(demo.ServerName)

	if strings.Contains(serverName, "faceit") || faceItDemoNameRegex.MatchString(demoName) {
		return constants.DemoSourceFaceIt
	}

	if strings.Contains(serverName, "cevo") {
		return constants.DemoSourceCEVO
	}

	if strings.Contains(serverName, "challengermode") {
		return constants.DemoSourceChallengermode
	}

	if strings.Contains(serverName, "esl") {
		return constants.DemoSourceESL
	}

	if strings.Contains(serverName, "ebot") || ebotDemoNameRegex.MatchString(demoName) {
		return constants.DemoSourceEbot
	}

	if strings.Contains(serverName, "esea") || strings.Contains(demoName, "esea") {
		return constants.DemoSourceESEA
	}

	if strings.Contains(serverName, "popflash") || strings.Contains(demoName, "popflash") {
		return constants.DemoSourcePopFlash
	}

	if strings.Contains(serverName, "esportal") {
		return constants.DemoSourceEsportal
	}

	if strings.Contains(serverName, "fastcup") {
		return constants.DemoSourceFastcup
	}

	if strings.Contains(serverName, "gamersclub") {
		return constants.DemoSourceGamersclub
	}

	if strings.Contains(serverName, "valve") {
		return constants.DemoSourceValve
	}

	if strings.Contains(serverName, "完美世界") {
		return constants.DemoSourcePerfectWorld
	}

	return constants.DemoSourceUnknown
}

func getMapNameFromHeaderMapName(headerMapName string) string {
	// Remove potential "_scrimmagemap" suffix.
	// Noticed with a de_mirage demo, it could be related to the fact that de_mirage has been moved from competitive maps
	// to unranked maps (called "scrimmage" maps back in 2019). It may be the case with others maps too.
	mapName := strings.ReplaceAll(headerMapName, "_scrimmagemap", "")

	// Remove potential workshop identifier prefix (i.e workshop/id/map_name).
	// Noticed with some FACEIT demos:
	// https://www.faceit.com/en/csgo/room/1-e1d60431-30f8-43d3-9985-0135500a7a97
	// https://www.faceit.com/en/csgo/room/1-f6e8327d-ce85-44f4-84ac-74b19b8ec738
	// Every time it's because the match's map is an old version, that's probably the reason.
	workshopRegex := regexp.MustCompile(`workshop\/(\d+\/)(?P<mapName>.*)`)
	match := workshopRegex.FindStringSubmatch(mapName)
	if len(match) == 3 {
		mapName = match[2]
	}

	return mapName
}
