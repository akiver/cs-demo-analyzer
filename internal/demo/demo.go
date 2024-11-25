package demo

import (
	"errors"
	"fmt"
	"hash/crc64"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/akiver/cs-demo-analyzer/internal/bitread"
	"github.com/akiver/cs-demo-analyzer/internal/filepath"
	str "github.com/akiver/cs-demo-analyzer/internal/strings"
	"github.com/akiver/cs-demo-analyzer/pkg/api/constants"
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/msg"
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/msgs2"
	"google.golang.org/protobuf/proto"
)

type Demo struct {
	FilePath string
	FileName string
	Type     constants.DemoType
	// ! Checksums are not generated reading the whole .dem file but only information for the "header" because it would be slow.
	Checksum                      string
	Filestamp                     string
	Date                          time.Time
	ServerName                    string
	ClientName                    string
	MapName                       string
	NetMessageDecryptionPublicKey []byte
	NetworkProtocol               int
	BuildNumber                   int           // Source 2 demos only
	TickCount                     int           // Not available for Source 2 demos, it's updated during parsing
	TickRate                      float64       // Not available for Source 2 demos, it's updated during parsing
	FrameRate                     float64       // Not available for Source 2 demos, it's updated during parsing
	Duration                      time.Duration // Not available for Source 2 demos, it's updated during parsing
	ShareCode                     string        // Valve demos only, the .info file must be next to the .dem file to be able to generate it
}

var faceItDemoNameRegex = regexp.MustCompile(`/[0-9]+_team[a-z0-9-]+-Team[a-z0-9-]+_de_[a-z0-9]+\.dem/`)
var ebotDemoNameRegex = regexp.MustCompile(`/([0-9]*)_(.*?)-(.*?)_(.*?)(.dem)/`)
var fiveEPlayDemoNameRegex = regexp.MustCompile(`^g\d+-(.*)[a-zA-Z0-9_]*$`)

// Default format: {TIME}_{MATCH_ID}_{MAP}_{TEAM1}_vs_{TEAM2}
// https://shobhit-pathak.github.io/MatchZy/gotv/#recording-demos
var matchZyDemoNameRegex = regexp.MustCompile(`^(\d{4}-\d{2}-\d{2}_\d{2}-\d{2}-\d{2})_(\d+)_([a-zA-Z0-9_]+)_(.+?)_vs_(.+)$`)

// Reads the .info file associated with a demo if it exists and returns its content as bytes.
func getMatchInfoProtoBytes(demoFilePath string) []byte {
	infoFilePath := demoFilePath + ".info"
	if _, err := os.Stat(infoFilePath); err != nil {
		return nil
	}

	file, err := os.Open(infoFilePath)
	if err != nil {
		fmt.Printf("Unable to open .info file: %v", err)
		return nil
	}

	bytes, err := io.ReadAll(file)
	if err != nil {
		fmt.Printf("Unable to read .info file: %v", err)
		return nil
	}

	return bytes
}

func getNetMessageDecryptionKeyFromPubKey(clDecryptDataKeyPub uint64) []byte {
	return []byte(strings.ToUpper(fmt.Sprintf("%016x", clDecryptDataKeyPub)))
}

func getDateFromMatchTime(matchTime uint32) time.Time {
	return time.Unix(int64(matchTime), 0)
}

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
	var date = stats.ModTime()
	var shareCode string
	var netMessageDecryptionPublicKey []byte
	demoType := constants.DemoTypeGOTV
	matchInfoBytes := getMatchInfoProtoBytes(demoPath)

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
		data := fmt.Sprintf(
			"%s%s%s%d%d%s%s%d",
			mapName,
			str.RemoveInvalidUTF8Sequences(serverName),
			str.RemoveInvalidUTF8Sequences(clientName),
			msg.GetNetworkProtocol(),
			msg.GetBuildNum(),
			msg.GetDemoVersionGuid(),
			msg.GetDemoVersionName(),
			stats.Size(),
		)
		checksum = strconv.FormatUint(crc64.Checksum([]byte(data), crc64.MakeTable(crc64.ECMA)), 16)

		serverName = str.ReplaceUTF8ByteSequences(serverName)
		clientName = str.ReplaceUTF8ByteSequences(clientName)
		networkProtocol = int(msg.GetNetworkProtocol())
		buildNumber = int(msg.GetBuildNum())
		game := msg.GetGame()
		if game != "" {
			demoType = constants.DemoTypePOV
		}

		if len(matchInfoBytes) > 0 {
			m := new(msgs2.CDataGCCStrike15V2_MatchInfo)
			err = proto.Unmarshal(matchInfoBytes, m)
			if err != nil {
				fmt.Printf("failed to unmarshal MatchInfo message: %v", err)
			} else {
				netMessageDecryptionPublicKey = getNetMessageDecryptionKeyFromPubKey(m.Watchablematchinfo.GetClDecryptdataKeyPub())
				date = getDateFromMatchTime(m.GetMatchtime())
				rounds := m.GetRoundstatsall()
				if len(rounds) > 0 {
					lastRound := rounds[len(rounds)-1]
					shareCode = encodeMatchShareCode(MatchInformation{
						MatchId:       m.GetMatchid(),
						ReservationId: lastRound.GetReservationid(),
						TvPort:        m.GetWatchablematchinfo().GetTvPort(),
					})
				}
			}
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

		data := fmt.Sprintf(
			"%s%s%s%d%d%d%d%d",
			mapName,
			serverName,
			clientName,
			frameCount,
			tickCount,
			networkProtocol,
			signonLength,
			stats.Size(),
		)
		checksum = strconv.FormatUint(crc64.Checksum([]byte(data), crc64.MakeTable(crc64.ECMA)), 16)

		if len(matchInfoBytes) > 0 {
			m := new(msg.CDataGCCStrike15V2_MatchInfo)
			err = proto.Unmarshal(matchInfoBytes, m)
			if err != nil {
				fmt.Printf("failed to unmarshal MatchInfo message: %v", err)
			} else {
				netMessageDecryptionPublicKey = getNetMessageDecryptionKeyFromPubKey(m.Watchablematchinfo.GetClDecryptdataKeyPub())
				date = getDateFromMatchTime(m.GetMatchtime())
				lastRound := m.GetRoundstatsLegacy()
				rounds := m.GetRoundstatsall()
				if lastRound == nil && len(rounds) > 0 {
					lastRound = rounds[len(rounds)-1]
				}
				if lastRound != nil {
					shareCode = encodeMatchShareCode(MatchInformation{
						MatchId:       m.GetMatchid(),
						ReservationId: lastRound.GetReservationid(),
						TvPort:        m.GetWatchablematchinfo().GetTvPort(),
					})
				}
			}
		}
	}

	return &Demo{
		FilePath:                      filepath.GetAbsoluteFilePath(demoPath),
		FileName:                      filepath.GetFileNameWithoutExtension(demoPath),
		Type:                          demoType,
		Filestamp:                     filestamp,
		Checksum:                      checksum,
		Date:                          date,
		ServerName:                    serverName,
		ClientName:                    clientName,
		Duration:                      duration,
		MapName:                       getMapNameFromHeaderMapName(mapName),
		TickCount:                     tickCount,
		TickRate:                      tickRate,
		FrameRate:                     frameRate,
		NetworkProtocol:               networkProtocol,
		BuildNumber:                   buildNumber,
		NetMessageDecryptionPublicKey: netMessageDecryptionPublicKey,
		ShareCode:                     shareCode,
	}, nil
}

func (demo *Demo) IsSource2() bool {
	return demo.Filestamp == "PBDEMS2"
}

func GetDemoSource(demo *Demo) constants.DemoSource {
	demoName := strings.ToLower(demo.FileName)
	serverName := strings.ToLower(demo.ServerName)

	if strings.Contains(serverName, "faceit") || strings.Contains(serverName, "blast") || faceItDemoNameRegex.MatchString(demoName) {
		return constants.DemoSourceFaceIt
	}

	if strings.Contains(serverName, "cevo") {
		return constants.DemoSourceCEVO
	}

	if strings.Contains(serverName, "challengermode") || strings.Contains(serverName, "pgl major cs2") {
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

	if strings.Contains(serverName, "matchzy") || matchZyDemoNameRegex.MatchString(demoName) {
		return constants.DemoSourceMatchZy
	}

	if strings.Contains(serverName, "valve") {
		return constants.DemoSourceValve
	}

	if strings.Contains(serverName, "完美世界") {
		return constants.DemoSourcePerfectWorld
	}

	if fiveEPlayDemoNameRegex.MatchString(demoName) {
		return constants.DemoSourceFiveEPlay
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
