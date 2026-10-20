// CSV export dedicated to the application CS Demo Manager.
//
// It's similar to the classic CSV export - the main difference is that there are no headers and some columns differ.
// CS:DM doesn't use the traditional CSV export because column order is crucial as it inserts demos into a database
// from CSV files.
// If CS:DM was using the classic CSV export, any changes to the traditional CSV export would break the application.
package api

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"os"
	"slices"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/akiver/cs-demo-analyzer/internal/converters"
	"github.com/akiver/cs-demo-analyzer/internal/csv"
	"github.com/akiver/cs-demo-analyzer/internal/slice"
	"github.com/akiver/cs-demo-analyzer/pkg/api/constants"
	"github.com/golang/geo/r2"
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/common"
)

// Weapons order may change between frames, sorting them avoids detecting false state changes.
func joinSortedWeaponNames(weaponNames []constants.WeaponName) string {
	names := slice.ToStrings(weaponNames)
	sort.Strings(names)

	return strings.Join(names, ",")
}

// Numbers and booleans of positions are packed as big-endian binary, it's the byte order used by PostgreSQL (int4send,
// float4send) which makes the values easier to query in SQL.
func appendInt32(bytes []byte, value int) []byte {
	return binary.BigEndian.AppendUint32(bytes, uint32(value))
}

func appendFloat32(bytes []byte, value float32) []byte {
	return binary.BigEndian.AppendUint32(bytes, math.Float32bits(value))
}

func appendBool(bytes []byte, value bool) []byte {
	if value {
		return append(bytes, 1)
	}
	return append(bytes, 0)
}

// PostgreSQL bytea hex format
func toBytea(bytes []byte) string {
	return "\\x" + hex.EncodeToString(bytes)
}

var textArrayEscaper = strings.NewReplacer(`\`, `\\`, `"`, `\"`)

// PostgreSQL array literal, elements are quoted because they may contain commas.
func toTextArray(values []string) string {
	elements := make([]string, len(values))
	for index, value := range values {
		elements[index] = `"` + textArrayEscaper.Replace(value) + `"`
	}
	return "{" + strings.Join(elements, ",") + "}"
}

// Positions of an entity packed as big-endian binary.
// A position is written only when it changes, entities such as smokes or chickens don't move most of the time.
// The position of the last tick is always written so that it's possible to know when the entity disappeared.
type packedPositions struct {
	ticks         []byte
	x             []byte
	y             []byte
	z             []byte
	lastTick      int
	lastX         float32
	lastY         float32
	lastZ         float32
	isLastWritten bool
}

func (positions *packedPositions) write(tick int, x float32, y float32, z float32) {
	positions.ticks = appendInt32(positions.ticks, tick)
	positions.x = appendFloat32(positions.x, x)
	positions.y = appendFloat32(positions.y, y)
	positions.z = appendFloat32(positions.z, z)
}

func (positions *packedPositions) add(tick int, x float64, y float64, z float64) {
	x32, y32, z32 := float32(x), float32(y), float32(z)
	if len(positions.ticks) == 0 {
		positions.write(tick, x32, y32, z32)
		positions.isLastWritten = true
	} else {
		// Positions are sorted by frame and there may be multiple frames per tick.
		if tick == positions.lastTick {
			return
		}

		hasMoved := x32 != positions.lastX || y32 != positions.lastY || z32 != positions.lastZ
		if hasMoved {
			// The entity didn't move since the last written position, write the position of the previous tick so that
			// the movement starts at the right tick.
			if !positions.isLastWritten {
				positions.write(positions.lastTick, positions.lastX, positions.lastY, positions.lastZ)
			}
			positions.write(tick, x32, y32, z32)
		}
		positions.isLastWritten = hasMoved
	}

	positions.lastTick = tick
	positions.lastX = x32
	positions.lastY = y32
	positions.lastZ = z32
}

// Writes the position of the last tick if it hasn't been written yet.
func (positions *packedPositions) finish() {
	if !positions.isLastWritten {
		positions.write(positions.lastTick, positions.lastX, positions.lastY, positions.lastZ)
		positions.isLastWritten = true
	}
}

func exportMatchForCSDM(match *Match, outputPath string) error {
	if stat, err := os.Stat(outputPath); err != nil || !stat.IsDir() {
		return errors.New("incorrect output provided, make sure it's a folder that exists and you have write access")
	}

	outputPath = outputPath + string(os.PathSeparator) + match.DemoFileName

	var writeDemo = func() {
		line := []string{
			match.Checksum,
			match.Game.String(),
			match.DemoFileName,
			match.Date.Format(time.RFC3339),
			match.Source.String(),
			match.Type.String(),
			match.ShareCode,
			match.MapName,
			match.ServerName,
			match.ClientName,
			converters.IntToString(match.TickCount),
			converters.Float64ToString(match.TickRate),
			converters.Float64ToString(match.FrameRate),
			converters.Float64ToString(match.Duration.Seconds()),
			converters.IntToString(match.NetworkProtocol),
			converters.IntToString(match.BuildNumber),
		}

		csv.WriteLinesIntoCsvFile(outputPath+"_demo.csv", [][]string{
			line,
		})
	}

	var writeMatch = func() {
		winnerName := ""
		winnerSide := common.TeamUnassigned
		if match.Winner != nil {
			winnerName = match.Winner.Name
			winnerSide = *match.Winner.CurrentSide
		}
		line := []string{
			match.Checksum,
			match.DemoFilePath,
			match.GameType.String(),
			match.GameMode.String(),
			match.GameModeStr().String(),
			converters.BoolToString(match.IsRanked),
			converters.IntToString(match.KillCount()),
			converters.IntToString(match.AssistCount()),
			converters.IntToString(match.DeathCount()),
			converters.IntToString(match.ShotCount()),
			winnerName,
			converters.TeamToString(winnerSide),
			converters.IntToString(match.OvertimeCount),
			converters.IntToString(match.MaxRounds),
			converters.BoolToString(match.HasVacLiveBan),
		}

		csv.WriteLinesIntoCsvFile(outputPath+"_match.csv", [][]string{
			line,
		})
	}

	var buildTeamLine = func(team *Team) []string {
		line := []string{
			team.Name,
			team.Letter.String(),
			converters.IntToString(team.Score),
			converters.IntToString(team.ScoreFirstHalf),
			converters.IntToString(team.ScoreSecondHalf),
			converters.TeamToString(*team.CurrentSide),
			match.Checksum,
		}

		return line
	}

	var writeTeams = func() {
		csv.WriteLinesIntoCsvFile(outputPath+"_teams.csv", [][]string{
			buildTeamLine(match.TeamA),
			buildTeamLine(match.TeamB),
		})
	}

	var writePlayers = func() {
		lines := [][]string{}
		for _, player := range match.Players() {
			line := []string{
				player.Name,
				converters.Uint64ToString(player.SteamID64),
				converters.IntToString(player.UserID + 1),
				converters.IntToString(player.Score),
				player.TeamName(),
				converters.IntToString(player.KillCount()),
				converters.IntToString(player.AssistCount()),
				converters.IntToString(player.DeathCount()),
				converters.IntToString(player.HeadshotCount()),
				converters.Float32ToString(player.KAST()),
				converters.Float32ToString(player.AverageDamagePerRound()),
				converters.Float32ToString(player.AverageKillPerRound()),
				converters.Float32ToString(player.AverageDeathPerRound()),
				converters.Float32ToString(player.UtilityDamagePerRound()),
				converters.IntToString(player.MvpCount),
				converters.IntToString(player.RankType),
				converters.IntToString(player.Rank),
				converters.IntToString(player.OldRank),
				converters.IntToString(player.WinCount),
				converters.IntToString(player.BombPlantedCount()),
				converters.IntToString(player.BombDefusedCount()),
				converters.IntToString(player.HostageRescuedCount()),
				converters.IntToString(player.HealthDamage()),
				converters.IntToString(player.ArmorDamage()),
				converters.IntToString(player.UtilityDamage()),
				converters.IntToString(player.FirstKillCount()),
				converters.IntToString(player.FirstDeathCount()),
				converters.IntToString(player.TradeKillCount()),
				converters.IntToString(player.TradeDeathCount()),
				converters.IntToString(player.FirstTradeKillCount()),
				converters.IntToString(player.FirstTradeDeathCount()),
				converters.IntToString(player.OneKillCount()),
				converters.IntToString(player.TwoKillCount()),
				converters.IntToString(player.ThreeKillCount()),
				converters.IntToString(player.FourKillCount()),
				converters.IntToString(player.FiveKillCount()),
				converters.Float32ToString(player.HltvRating2()),
				converters.Float32ToString(player.HltvRating()),
				player.CrosshairShareCode,
				converters.ColorToString(player.Color),
				converters.IntToString(player.InspectWeaponCount),
				match.Checksum,
			}
			lines = append(lines, line)
		}

		csv.WriteLinesIntoCsvFile(outputPath+"_players.csv", lines)
	}

	// Players positions, 1 row per player and round with:
	// - the movement (position, yaw...) of each tick
	// - the state (health, weapons...) only when it changes, stateTicks contains the tick of each change
	// The first state of each player in a round is always written so that a round can be loaded independently.
	// Numbers and booleans are packed as big-endian binary and weapon names are PostgreSQL text arrays.
	var writePlayerPositions = func() {
		type playerRoundKey struct {
			steamID     uint64
			botName     string
			roundNumber int
		}
		type playerState struct {
			isAlive           bool
			side              common.Team
			money             int
			health            int
			armor             int
			hasHelmet         bool
			hasBomb           bool
			hasDefuseKit      bool
			isScoping         bool
			isDefusing        bool
			isPlanting        bool
			isGrabbingHostage bool
			activeWeaponName  string
			equipments        string
			grenades          string
			pistols           string
			smgs              string
			rifles            string
			heavy             string
		}
		type playerRoundPositions struct {
			steamID                uint64
			name                   string
			roundNumber            int
			lastTick               int
			ticks                  []byte
			x                      []byte
			y                      []byte
			z                      []byte
			yaw                    []byte
			flashDurationRemaining []byte
			isDucking              []byte
			isAirborne             []byte
			lastState              *playerState
			stateTicks             []byte
			isAlive                []byte
			side                   []byte
			money                  []byte
			health                 []byte
			armor                  []byte
			hasHelmet              []byte
			hasBomb                []byte
			hasDefuseKit           []byte
			isScoping              []byte
			isDefusing             []byte
			isPlanting             []byte
			isGrabbingHostage      []byte
			activeWeaponName       []string
			equipments             []string
			grenades               []string
			pistols                []string
			smgs                   []string
			rifles                 []string
			heavy                  []string
		}
		positionsPerPlayer := make(map[playerRoundKey]*playerRoundPositions)
		var orderedPositions []*playerRoundPositions
		for _, position := range match.PlayerPositions {
			// BOTs SteamID is always 0, the name is used to differentiate them.
			// It's not used for real players because they may change their name during a match.
			botName := ""
			if position.SteamID64 == 0 {
				botName = position.Name
			}
			key := playerRoundKey{
				steamID:     position.SteamID64,
				botName:     botName,
				roundNumber: position.RoundNumber,
			}
			positions, exists := positionsPerPlayer[key]
			if !exists {
				positions = &playerRoundPositions{
					steamID:     position.SteamID64,
					name:        position.Name,
					roundNumber: position.RoundNumber,
					lastTick:    -1,
				}
				positionsPerPlayer[key] = positions
				orderedPositions = append(orderedPositions, positions)
			}

			// Positions are sorted by frame and there may be multiple frames per tick.
			if positions.lastTick == position.Tick {
				continue
			}
			positions.lastTick = position.Tick

			positions.ticks = appendInt32(positions.ticks, position.Tick)
			positions.x = appendFloat32(positions.x, float32(position.X))
			positions.y = appendFloat32(positions.y, float32(position.Y))
			positions.z = appendFloat32(positions.z, float32(position.Z))
			positions.yaw = appendFloat32(positions.yaw, position.Yaw)
			positions.flashDurationRemaining = appendFloat32(positions.flashDurationRemaining, float32(position.FlashDurationRemaining))
			positions.isDucking = appendBool(positions.isDucking, position.IsDucking)
			positions.isAirborne = appendBool(positions.isAirborne, position.IsAirborne)

			state := playerState{
				isAlive:           position.IsAlive,
				side:              position.Side,
				money:             position.Money,
				health:            position.Health,
				armor:             position.Armor,
				hasHelmet:         position.HasHelmet,
				hasBomb:           position.HasBomb,
				hasDefuseKit:      position.HasDefuseKit,
				isScoping:         position.IsScoping,
				isDefusing:        position.IsDefusing,
				isPlanting:        position.IsPlanting,
				isGrabbingHostage: position.IsGrabbingHostage,
				activeWeaponName:  position.ActiveWeaponName.String(),
				equipments:        joinSortedWeaponNames(position.Equipments),
				grenades:          joinSortedWeaponNames(position.Grenades),
				pistols:           joinSortedWeaponNames(position.Pistols),
				smgs:              joinSortedWeaponNames(position.SMGs),
				rifles:            joinSortedWeaponNames(position.Rifles),
				heavy:             joinSortedWeaponNames(position.Heavy),
			}
			if positions.lastState != nil && *positions.lastState == state {
				continue
			}
			positions.lastState = &state

			positions.stateTicks = appendInt32(positions.stateTicks, position.Tick)
			positions.isAlive = appendBool(positions.isAlive, state.isAlive)
			positions.side = appendInt32(positions.side, int(state.side))
			positions.money = appendInt32(positions.money, state.money)
			positions.health = appendInt32(positions.health, state.health)
			positions.armor = appendInt32(positions.armor, state.armor)
			positions.hasHelmet = appendBool(positions.hasHelmet, state.hasHelmet)
			positions.hasBomb = appendBool(positions.hasBomb, state.hasBomb)
			positions.hasDefuseKit = appendBool(positions.hasDefuseKit, state.hasDefuseKit)
			positions.isScoping = appendBool(positions.isScoping, state.isScoping)
			positions.isDefusing = appendBool(positions.isDefusing, state.isDefusing)
			positions.isPlanting = appendBool(positions.isPlanting, state.isPlanting)
			positions.isGrabbingHostage = appendBool(positions.isGrabbingHostage, state.isGrabbingHostage)
			positions.activeWeaponName = append(positions.activeWeaponName, state.activeWeaponName)
			positions.equipments = append(positions.equipments, state.equipments)
			positions.grenades = append(positions.grenades, state.grenades)
			positions.pistols = append(positions.pistols, state.pistols)
			positions.smgs = append(positions.smgs, state.smgs)
			positions.rifles = append(positions.rifles, state.rifles)
			positions.heavy = append(positions.heavy, state.heavy)
		}

		lines := [][]string{}
		for _, positions := range orderedPositions {
			line := []string{
				converters.IntToString(positions.roundNumber),
				converters.Uint64ToString(positions.steamID),
				positions.name,
				toBytea(positions.ticks),
				toBytea(positions.x),
				toBytea(positions.y),
				toBytea(positions.z),
				toBytea(positions.yaw),
				toBytea(positions.flashDurationRemaining),
				toBytea(positions.isDucking),
				toBytea(positions.isAirborne),
				toBytea(positions.stateTicks),
				toBytea(positions.isAlive),
				toBytea(positions.side),
				toBytea(positions.money),
				toBytea(positions.health),
				toBytea(positions.armor),
				toBytea(positions.hasHelmet),
				toBytea(positions.hasBomb),
				toBytea(positions.hasDefuseKit),
				toBytea(positions.isScoping),
				toBytea(positions.isDefusing),
				toBytea(positions.isPlanting),
				toBytea(positions.isGrabbingHostage),
				toTextArray(positions.activeWeaponName),
				toTextArray(positions.equipments),
				toTextArray(positions.grenades),
				toTextArray(positions.pistols),
				toTextArray(positions.smgs),
				toTextArray(positions.rifles),
				toTextArray(positions.heavy),
				match.Checksum,
			}
			lines = append(lines, line)
		}

		csv.WriteLinesIntoCsvFile(outputPath+"_positions.csv", lines)
	}

	var writeShots = func() {
		lines := [][]string{}
		for _, shot := range match.Shots {
			line := []string{
				converters.IntToString(shot.Frame),
				converters.IntToString(shot.Tick),
				converters.IntToString(shot.RoundNumber),
				shot.WeaponName.String(),
				shot.WeaponID,
				shot.ProjectileID,
				converters.Float64ToString(shot.X),
				converters.Float64ToString(shot.Y),
				converters.Float64ToString(shot.Z),
				shot.PlayerName,
				converters.Uint64ToString(shot.PlayerSteamID64),
				shot.PlayerTeamName,
				converters.TeamToString(shot.PlayerSide),
				converters.BoolToString(shot.IsPlayerControllingBot),
				converters.Float32ToString(shot.Yaw),
				converters.Float32ToString(shot.Pitch),
				converters.Float64ToString(shot.PlayerVelocityX),
				converters.Float64ToString(shot.PlayerVelocityY),
				converters.Float64ToString(shot.PlayerVelocityZ),
				converters.Float32ToString(shot.RecoilIndex),
				converters.Float64ToString(shot.AimPunchAngleX),
				converters.Float64ToString(shot.AimPunchAngleY),
				converters.Float64ToString(shot.ViewPunchAngleX),
				converters.Float64ToString(shot.ViewPunchAngleY),
				match.Checksum,
			}
			lines = append(lines, line)
		}

		csv.WriteLinesIntoCsvFile(outputPath+"_shots.csv", lines)
	}

	var writeRounds = func() {
		lines := [][]string{}
		for _, round := range match.Rounds {
			line := []string{
				converters.IntToString(round.Number),
				converters.IntToString(round.StartTick),
				converters.IntToString(round.StartFrame),
				converters.IntToString(round.FreezeTimeEndTick),
				converters.IntToString(round.FreezeTimeEndFrame),
				converters.IntToString(round.EndTick),
				converters.IntToString(round.EndFrame),
				converters.IntToString(round.EndOfficiallyTick),
				converters.IntToString(round.EndOfficiallyFrame),
				round.TeamAName,
				round.TeamBName,
				converters.IntToString(round.TeamAScore),
				converters.IntToString(round.TeamBScore),
				converters.TeamToString(round.TeamASide),
				converters.TeamToString(round.TeamBSide),
				converters.IntToString(round.StartMoneyTeamA()),
				converters.IntToString(round.StartMoneyTeamB()),
				converters.IntToString(round.TeamAEquipmentValue),
				converters.IntToString(round.TeamBEquipmentValue),
				converters.IntToString(round.TeamAMoneySpent),
				converters.IntToString(round.TeamBMoneySpent),
				round.TeamAEconomyType.String(),
				round.TeamBEconomyType.String(),
				converters.Int64ToString(round.Duration),
				converters.RoundEndReasonToString(round.EndReason),
				round.WinnerName,
				converters.TeamToString(round.WinnerSide),
				converters.IntToString(round.OvertimeNumber),
				match.Checksum,
			}
			lines = append(lines, line)
		}

		csv.WriteLinesIntoCsvFile(outputPath+"_rounds.csv", lines)
	}

	var writeRoundEconomies = func() {
		lines := [][]string{}
		for _, economy := range match.PlayerEconomies {
			line := []string{
				converters.Uint64ToString(economy.SteamID64),
				economy.Name,
				converters.TeamToString(economy.PlayerSide),
				converters.IntToString(economy.StartMoney),
				converters.IntToString(economy.MoneySpent),
				converters.IntToString(economy.EquipmentValue),
				economy.Type.String(),
				converters.IntToString(economy.RoundNumber),
				match.Checksum,
			}
			lines = append(lines, line)
		}

		csv.WriteLinesIntoCsvFile(outputPath+"_players_economy.csv", lines)
	}

	var writeClutches = func() {
		lines := [][]string{}
		for _, clutch := range match.Clutches {
			line := []string{
				converters.IntToString(clutch.Frame),
				converters.IntToString(clutch.Tick),
				converters.IntToString(clutch.RoundNumber),
				converters.IntToString(clutch.OpponentCount),
				converters.TeamToString(clutch.Side),
				converters.BoolToString(clutch.HasWon),
				converters.Uint64ToString(clutch.ClutcherSteamID64),
				clutch.ClutcherName,
				converters.BoolToString(clutch.ClutcherSurvived),
				converters.IntToString(clutch.ClutcherKillCount),
				match.Checksum,
			}
			lines = append(lines, line)
		}

		csv.WriteLinesIntoCsvFile(outputPath+"_clutches.csv", lines)
	}

	var writeChickenDeaths = func() {
		lines := [][]string{}
		for _, chickenDeath := range match.ChickenDeaths {
			line := []string{
				converters.IntToString(chickenDeath.Frame),
				converters.IntToString(chickenDeath.Tick),
				converters.IntToString(chickenDeath.RoundNumber),
				converters.Uint64ToString(chickenDeath.KillerSteamID),
				chickenDeath.WeaponName.String(),
				match.Checksum,
			}
			lines = append(lines, line)
		}

		csv.WriteLinesIntoCsvFile(outputPath+"_chicken_deaths.csv", lines)
	}

	// Chickens positions, 1 row per chicken and round, see packedPositions.
	var writeChickenPositions = func() {
		type chickenRoundKey struct {
			chickenID   int
			roundNumber int
		}
		type chickenRoundPositions struct {
			roundNumber int
			positions   packedPositions
		}

		positionsPerChicken := make(map[chickenRoundKey]*chickenRoundPositions)
		var orderedPositions []*chickenRoundPositions
		for _, position := range match.ChickenPositions {
			key := chickenRoundKey{
				chickenID:   position.chickenID,
				roundNumber: position.RoundNumber,
			}
			chicken, exists := positionsPerChicken[key]
			if !exists {
				chicken = &chickenRoundPositions{roundNumber: position.RoundNumber}
				positionsPerChicken[key] = chicken
				orderedPositions = append(orderedPositions, chicken)
			}
			chicken.positions.add(position.Tick, position.X, position.Y, position.Z)
		}

		lines := [][]string{}
		for _, chicken := range orderedPositions {
			chicken.positions.finish()
			line := []string{
				converters.IntToString(chicken.roundNumber),
				toBytea(chicken.positions.ticks),
				toBytea(chicken.positions.x),
				toBytea(chicken.positions.y),
				toBytea(chicken.positions.z),
				match.Checksum,
			}
			lines = append(lines, line)
		}

		csv.WriteLinesIntoCsvFile(outputPath+"_chicken_positions.csv", lines)
	}

	var writeDamages = func() {
		lines := [][]string{}
		for _, damage := range match.Damages {
			line := []string{
				converters.IntToString(damage.Frame),
				converters.IntToString(damage.Tick),
				converters.IntToString(damage.RoundNumber),
				converters.IntToString(damage.HealthDamage),
				converters.IntToString(damage.ArmorDamage),
				converters.IntToString(damage.VictimHealth),
				converters.IntToString(damage.VictimNewHealth),
				converters.IntToString(damage.VictimArmor),
				converters.IntToString(damage.VictimNewArmor),
				converters.Uint64ToString(damage.AttackerSteamID64),
				converters.TeamToString(damage.AttackerSide),
				damage.AttackerTeamName,
				converters.BoolToString(damage.IsAttackerControllingBot),
				converters.Uint64ToString(damage.VictimSteamID64),
				converters.TeamToString(damage.VictimSide),
				damage.VictimTeamName,
				converters.BoolToString(damage.IsVictimControllingBot),
				damage.WeaponName.String(),
				string(damage.WeaponType),
				converters.HitgroupToString(damage.HitGroup),
				damage.WeaponUniqueID,
				match.Checksum,
			}
			lines = append(lines, line)
		}

		csv.WriteLinesIntoCsvFile(outputPath+"_damages.csv", lines)
	}

	var writeKills = func() {
		lines := [][]string{}
		for _, kill := range match.Kills {
			line := []string{
				converters.IntToString(kill.Frame),
				converters.IntToString(kill.Tick),
				converters.IntToString(kill.RoundNumber),
				kill.KillerName,
				converters.Uint64ToString(kill.KillerSteamID64),
				converters.TeamToString(kill.KillerSide),
				kill.KillerTeamName,
				kill.VictimName,
				converters.Uint64ToString(kill.VictimSteamID64),
				converters.TeamToString(kill.VictimSide),
				kill.VictimTeamName,
				kill.AssisterName,
				converters.Uint64ToString(kill.AssisterSteamID64),
				converters.TeamToString(kill.AssisterSide),
				kill.AssisterTeamName,
				kill.WeaponName.String(),
				string(kill.WeaponType),
				converters.BoolToString(kill.IsHeadshot),
				converters.IntToString(kill.PenetratedObjects),
				converters.BoolToString(kill.IsAssistedFlash),
				converters.BoolToString(kill.IsKillerControllingBot),
				converters.BoolToString(kill.IsVictimControllingBot),
				converters.BoolToString(kill.IsAssisterControllingBot),
				converters.Float64ToString(kill.KillerX),
				converters.Float64ToString(kill.KillerY),
				converters.Float64ToString(kill.KillerZ),
				converters.BoolToString(kill.IsKillerAirborne),
				converters.BoolToString(kill.IsKillerBlinded),
				converters.Float64ToString(kill.VictimX),
				converters.Float64ToString(kill.VictimY),
				converters.Float64ToString(kill.VictimZ),
				converters.BoolToString(kill.IsVictimAirborne),
				converters.BoolToString(kill.IsVictimBlinded),
				converters.BoolToString(kill.IsVictimInspectingWeapon),
				converters.Float64ToString(kill.AssisterX),
				converters.Float64ToString(kill.AssisterY),
				converters.Float64ToString(kill.AssisterZ),
				converters.BoolToString(kill.IsTradeKill),
				converters.BoolToString(kill.IsTradeDeath),
				converters.BoolToString(kill.IsThroughSmoke),
				converters.BoolToString(kill.IsNoScope),
				converters.Float32ToString(kill.Distance),
				match.Checksum,
			}
			lines = append(lines, line)
		}

		csv.WriteLinesIntoCsvFile(outputPath+"_kills.csv", lines)
	}

	var writeBombsPlanted = func() {
		lines := [][]string{}
		for _, bombPlanted := range match.BombsPlanted {
			line := []string{
				converters.IntToString(bombPlanted.Frame),
				converters.IntToString(bombPlanted.Tick),
				converters.IntToString(bombPlanted.RoundNumber),
				bombPlanted.Site,
				converters.Uint64ToString(bombPlanted.PlanterSteamID64),
				bombPlanted.PlanterName,
				converters.BoolToString(bombPlanted.IsPlayerControllingBot),
				converters.Float64ToString(bombPlanted.X),
				converters.Float64ToString(bombPlanted.Y),
				converters.Float64ToString(bombPlanted.Z),
				match.Checksum,
			}
			lines = append(lines, line)
		}

		csv.WriteLinesIntoCsvFile(outputPath+"_bombs_planted.csv", lines)
	}

	var writeBombsDefuseStart = func() {
		lines := [][]string{}
		for _, bombDefuseStart := range match.BombsDefuseStart {
			line := []string{
				converters.IntToString(bombDefuseStart.Frame),
				converters.IntToString(bombDefuseStart.Tick),
				converters.IntToString(bombDefuseStart.RoundNumber),
				converters.Uint64ToString(bombDefuseStart.PlanterSteamID64),
				bombDefuseStart.PlanterName,
				converters.BoolToString(bombDefuseStart.IsPlayerControllingBot),
				converters.Float64ToString(bombDefuseStart.X),
				converters.Float64ToString(bombDefuseStart.Y),
				converters.Float64ToString(bombDefuseStart.Z),
				match.Checksum,
			}
			lines = append(lines, line)
		}

		csv.WriteLinesIntoCsvFile(outputPath+"_bombs_defuse_start.csv", lines)
	}

	var writeBombsDefused = func() {
		lines := [][]string{}
		for _, bombDefused := range match.BombsDefused {
			line := []string{
				converters.IntToString(bombDefused.Frame),
				converters.IntToString(bombDefused.Tick),
				converters.IntToString(bombDefused.RoundNumber),
				bombDefused.Site,
				converters.Uint64ToString(bombDefused.DefuserSteamID64),
				bombDefused.DefuserName,
				converters.BoolToString(bombDefused.IsPlayerControllingBot),
				converters.Float64ToString(bombDefused.X),
				converters.Float64ToString(bombDefused.Y),
				converters.Float64ToString(bombDefused.Z),
				converters.IntToString(bombDefused.CounterTerroristAliveCount),
				converters.IntToString(bombDefused.TerroristAliveCount),
				match.Checksum,
			}
			lines = append(lines, line)
		}

		csv.WriteLinesIntoCsvFile(outputPath+"_bombs_defused.csv", lines)
	}

	var writeBombsExploded = func() {
		lines := [][]string{}
		for _, bombExploded := range match.BombsExploded {
			line := []string{
				converters.IntToString(bombExploded.Frame),
				converters.IntToString(bombExploded.Tick),
				converters.IntToString(bombExploded.RoundNumber),
				bombExploded.Site,
				converters.Uint64ToString(bombExploded.PlanterSteamID64),
				bombExploded.PlanterName,
				converters.BoolToString(bombExploded.IsPlayerControllingBot),
				converters.Float64ToString(bombExploded.X),
				converters.Float64ToString(bombExploded.Y),
				converters.Float64ToString(bombExploded.Z),
				match.Checksum,
			}
			lines = append(lines, line)
		}

		csv.WriteLinesIntoCsvFile(outputPath+"_bombs_exploded.csv", lines)
	}

	var writeBombsPlantStart = func() {
		lines := [][]string{}
		for _, bombPlantStart := range match.BombsPlantStart {
			line := []string{
				converters.IntToString(bombPlantStart.Frame),
				converters.IntToString(bombPlantStart.Tick),
				converters.IntToString(bombPlantStart.RoundNumber),
				bombPlantStart.Site,
				converters.Uint64ToString(bombPlantStart.PlanterSteamID64),
				bombPlantStart.PlanterName,
				converters.BoolToString(bombPlantStart.IsPlayerControllingBot),
				converters.Float64ToString(bombPlantStart.X),
				converters.Float64ToString(bombPlantStart.Y),
				converters.Float64ToString(bombPlantStart.Z),
				match.Checksum,
			}
			lines = append(lines, line)
		}

		csv.WriteLinesIntoCsvFile(outputPath+"_bombs_plant_start.csv", lines)
	}

	var writePlayersFlashed = func() {
		lines := [][]string{}
		for _, playerFlashed := range match.PlayersFlashed {
			line := []string{
				converters.IntToString(playerFlashed.Frame),
				converters.IntToString(playerFlashed.Tick),
				converters.IntToString(playerFlashed.RoundNumber),
				converters.Float32ToString(playerFlashed.Duration),
				converters.Uint64ToString(playerFlashed.FlashedSteamID64),
				playerFlashed.FlashedName,
				converters.TeamToString(playerFlashed.FlashedSide),
				converters.BoolToString(playerFlashed.IsFlashedControllingBot),
				converters.Uint64ToString(playerFlashed.FlasherSteamID64),
				playerFlashed.FlasherName,
				converters.TeamToString(playerFlashed.FlasherSide),
				converters.BoolToString(playerFlashed.IsFlasherControllingBot),
				match.Checksum,
			}
			lines = append(lines, line)
		}

		csv.WriteLinesIntoCsvFile(outputPath+"_players_flashed.csv", lines)
	}

	var writePlayersBuy = func() {
		lines := [][]string{}
		for _, playerBuy := range match.PlayersBuy {
			line := []string{
				converters.IntToString(playerBuy.Frame),
				converters.IntToString(playerBuy.Tick),
				converters.IntToString(playerBuy.RoundNumber),
				converters.Uint64ToString(playerBuy.PlayerSteamID64),
				converters.TeamToString(playerBuy.PlayerSide),
				playerBuy.PlayerName,
				playerBuy.WeaponName.String(),
				playerBuy.WeaponType.String(),
				playerBuy.WeaponUniqueID,
				converters.BoolToString(playerBuy.HasRefunded),
				match.Checksum,
			}
			lines = append(lines, line)
		}

		csv.WriteLinesIntoCsvFile(outputPath+"_players_buy.csv", lines)
	}

	// Grenades positions, 1 row per projectile, see packedPositions.
	// Thrower values are the ones at the time of the throw.
	var writeGrenadePositions = func() {
		type projectilePositions struct {
			first     *GrenadePosition
			positions packedPositions
		}

		positionsPerProjectile := make(map[string]*projectilePositions)
		var orderedPositions []*projectilePositions
		for _, position := range match.GrenadePositions {
			projectile, exists := positionsPerProjectile[position.ProjectileID]
			if !exists {
				projectile = &projectilePositions{first: position}
				positionsPerProjectile[position.ProjectileID] = projectile
				orderedPositions = append(orderedPositions, projectile)
			}
			projectile.positions.add(position.Tick, position.X, position.Y, position.Z)
		}

		lines := [][]string{}
		for _, projectile := range orderedPositions {
			projectile.positions.finish()
			first := projectile.first
			line := []string{
				converters.IntToString(first.RoundNumber),
				first.GrenadeID,
				first.ProjectileID,
				first.GrenadeName.String(),
				converters.Uint64ToString(first.ThrowerSteamID64),
				first.ThrowerName,
				converters.TeamToString(first.ThrowerSide),
				first.ThrowerTeamName,
				converters.Float64ToString(first.ThrowerVelocityX),
				converters.Float64ToString(first.ThrowerVelocityY),
				converters.Float64ToString(first.ThrowerVelocityZ),
				converters.Float32ToString(first.ThrowerYaw),
				converters.Float32ToString(first.ThrowerPitch),
				toBytea(projectile.positions.ticks),
				toBytea(projectile.positions.x),
				toBytea(projectile.positions.y),
				toBytea(projectile.positions.z),
				match.Checksum,
			}
			lines = append(lines, line)
		}

		csv.WriteLinesIntoCsvFile(outputPath+"_grenade_positions.csv", lines)
	}

	var writeGrenadeBounces = func() {
		lines := [][]string{}
		for _, position := range match.GrenadeBounces {
			line := []string{
				converters.IntToString(position.Frame),
				converters.IntToString(position.Tick),
				converters.IntToString(position.RoundNumber),
				position.GrenadeID,
				position.ProjectileID,
				position.GrenadeName.String(),
				converters.Float64ToString(position.X),
				converters.Float64ToString(position.Y),
				converters.Float64ToString(position.Z),
				converters.Uint64ToString(position.ThrowerSteamID64),
				position.ThrowerName,
				converters.TeamToString(position.ThrowerSide),
				position.ThrowerTeamName,
				converters.Float64ToString(position.ThrowerVelocityX),
				converters.Float64ToString(position.ThrowerVelocityY),
				converters.Float64ToString(position.ThrowerVelocityZ),
				converters.Float32ToString(position.ThrowerYaw),
				converters.Float32ToString(position.ThrowerPitch),
				match.Checksum,
			}
			lines = append(lines, line)
		}

		csv.WriteLinesIntoCsvFile(outputPath+"_grenade_bounces.csv", lines)
	}

	var writeGrenadeProjectilesDestroy = func() {
		lines := [][]string{}
		for _, event := range match.GrenadeProjectilesDestroy {
			line := []string{
				converters.IntToString(event.Frame),
				converters.IntToString(event.Tick),
				converters.IntToString(event.RoundNumber),
				event.GrenadeID,
				event.ProjectileID,
				event.GrenadeName.String(),
				converters.Float64ToString(event.X),
				converters.Float64ToString(event.Y),
				converters.Float64ToString(event.Z),
				converters.Uint64ToString(event.ThrowerSteamID64),
				event.ThrowerName,
				converters.TeamToString(event.ThrowerSide),
				event.ThrowerTeamName,
				converters.Float64ToString(event.ThrowerVelocityX),
				converters.Float64ToString(event.ThrowerVelocityY),
				converters.Float64ToString(event.ThrowerVelocityZ),
				converters.Float32ToString(event.ThrowerYaw),
				converters.Float32ToString(event.ThrowerPitch),
				match.Checksum,
			}
			lines = append(lines, line)
		}

		csv.WriteLinesIntoCsvFile(outputPath+"_grenade_projectiles_destroy.csv", lines)
	}

	var writeSmokesStart = func() {
		lines := [][]string{}
		for _, event := range match.SmokesStart {
			line := []string{
				converters.IntToString(event.Frame),
				converters.IntToString(event.Tick),
				converters.IntToString(event.RoundNumber),
				event.GrenadeID,
				event.ProjectileID,
				converters.Float64ToString(event.X),
				converters.Float64ToString(event.Y),
				converters.Float64ToString(event.Z),
				converters.Uint64ToString(event.ThrowerSteamID64),
				event.ThrowerName,
				converters.TeamToString(event.ThrowerSide),
				event.ThrowerTeamName,
				converters.Float64ToString(event.ThrowerVelocityX),
				converters.Float64ToString(event.ThrowerVelocityY),
				converters.Float64ToString(event.ThrowerVelocityZ),
				converters.Float32ToString(event.ThrowerYaw),
				converters.Float32ToString(event.ThrowerPitch),
				match.Checksum,
			}
			lines = append(lines, line)
		}

		csv.WriteLinesIntoCsvFile(outputPath+"_smokes_start.csv", lines)
	}

	var writeHeGrenadesExplode = func() {
		lines := [][]string{}
		for _, event := range match.HeGrenadesExplode {
			line := []string{
				converters.IntToString(event.Frame),
				converters.IntToString(event.Tick),
				converters.IntToString(event.RoundNumber),
				event.GrenadeID,
				event.ProjectileID,
				converters.Float64ToString(event.X),
				converters.Float64ToString(event.Y),
				converters.Float64ToString(event.Z),
				converters.Uint64ToString(event.ThrowerSteamID64),
				event.ThrowerName,
				converters.TeamToString(event.ThrowerSide),
				event.ThrowerTeamName,
				converters.Float64ToString(event.ThrowerVelocityX),
				converters.Float64ToString(event.ThrowerVelocityY),
				converters.Float64ToString(event.ThrowerVelocityZ),
				converters.Float32ToString(event.ThrowerYaw),
				converters.Float32ToString(event.ThrowerPitch),
				match.Checksum,
			}
			lines = append(lines, line)
		}

		csv.WriteLinesIntoCsvFile(outputPath+"_he_grenades_explode.csv", lines)
	}

	var writeFlashbangsExplode = func() {
		lines := [][]string{}
		for _, event := range match.FlashbangsExplode {
			line := []string{
				converters.IntToString(event.Frame),
				converters.IntToString(event.Tick),
				converters.IntToString(event.RoundNumber),
				event.GrenadeID,
				event.ProjectileID,
				converters.Float64ToString(event.X),
				converters.Float64ToString(event.Y),
				converters.Float64ToString(event.Z),
				converters.Uint64ToString(event.ThrowerSteamID64),
				event.ThrowerName,
				converters.TeamToString(event.ThrowerSide),
				event.ThrowerTeamName,
				converters.Float64ToString(event.ThrowerVelocityX),
				converters.Float64ToString(event.ThrowerVelocityY),
				converters.Float64ToString(event.ThrowerVelocityZ),
				converters.Float32ToString(event.ThrowerYaw),
				converters.Float32ToString(event.ThrowerPitch),
				match.Checksum,
			}
			lines = append(lines, line)
		}

		csv.WriteLinesIntoCsvFile(outputPath+"_flashbangs_explode.csv", lines)
	}

	var writeDecoysStart = func() {
		lines := [][]string{}
		for _, event := range match.DecoysStart {
			line := []string{
				converters.IntToString(event.Frame),
				converters.IntToString(event.Tick),
				converters.IntToString(event.RoundNumber),
				event.GrenadeID,
				event.ProjectileID,
				converters.Float64ToString(event.X),
				converters.Float64ToString(event.Y),
				converters.Float64ToString(event.Z),
				converters.Uint64ToString(event.ThrowerSteamID64),
				event.ThrowerName,
				converters.TeamToString(event.ThrowerSide),
				event.ThrowerTeamName,
				converters.Float64ToString(event.ThrowerVelocityX),
				converters.Float64ToString(event.ThrowerVelocityY),
				converters.Float64ToString(event.ThrowerVelocityZ),
				converters.Float32ToString(event.ThrowerYaw),
				converters.Float32ToString(event.ThrowerPitch),
				match.Checksum,
			}
			lines = append(lines, line)
		}

		csv.WriteLinesIntoCsvFile(outputPath+"_decoys_start.csv", lines)
	}

	// Infernos positions, 1 row per inferno. The position and the thrower of an inferno never change, only its 2D convex
	// hull changes when fires ignite or extinguish.
	// A convex hull is written only when it changes with its tick packed as big-endian binary, the convex hull of the
	// last tick is always written so that it's possible to know when the inferno disappeared.
	// A convex hull is a list of points "x1,y1,x2,y2...", it's empty when the inferno has no active fires.
	var writeInfernoPositions = func() {
		type infernoHulls struct {
			first         *InfernoPosition
			ticks         []byte
			hulls         []string
			lastTick      int
			lastHull      []r2.Point
			isLastWritten bool
		}
		var formatHull = func(hull []r2.Point) string {
			var builder strings.Builder
			for index, point := range hull {
				if index > 0 {
					builder.WriteString(",")
				}
				fmt.Fprintf(&builder, "%f,%f", point.X, point.Y)
			}
			return builder.String()
		}
		var writeHull = func(inferno *infernoHulls, tick int, hull []r2.Point) {
			inferno.ticks = appendInt32(inferno.ticks, tick)
			inferno.hulls = append(inferno.hulls, formatHull(hull))
		}

		hullsPerInferno := make(map[string]*infernoHulls)
		var orderedHulls []*infernoHulls
		for _, position := range match.InfernoPositions {
			inferno, exists := hullsPerInferno[position.UniqueID]
			if !exists {
				inferno = &infernoHulls{first: position}
				hullsPerInferno[position.UniqueID] = inferno
				orderedHulls = append(orderedHulls, inferno)
				writeHull(inferno, position.Tick, position.ConvexHull2D)
				inferno.isLastWritten = true
			} else {
				// Positions are sorted by frame and there may be multiple frames per tick.
				if position.Tick == inferno.lastTick {
					continue
				}
				hasChanged := !slices.Equal(position.ConvexHull2D, inferno.lastHull)
				if hasChanged {
					writeHull(inferno, position.Tick, position.ConvexHull2D)
				}
				inferno.isLastWritten = hasChanged
			}
			inferno.lastTick = position.Tick
			inferno.lastHull = position.ConvexHull2D
		}

		lines := [][]string{}
		for _, inferno := range orderedHulls {
			if !inferno.isLastWritten {
				writeHull(inferno, inferno.lastTick, inferno.lastHull)
			}
			first := inferno.first
			line := []string{
				converters.IntToString(first.RoundNumber),
				first.UniqueID,
				converters.Uint64ToString(first.ThrowerSteamID64),
				first.ThrowerName,
				converters.Float64ToString(first.X),
				converters.Float64ToString(first.Y),
				converters.Float64ToString(first.Z),
				toBytea(inferno.ticks),
				toTextArray(inferno.hulls),
				match.Checksum,
			}
			lines = append(lines, line)
		}

		csv.WriteLinesIntoCsvFile(outputPath+"_inferno_positions.csv", lines)
	}

	var writeChatMessages = func() {
		lines := [][]string{}
		for _, chatMessage := range match.ChatMessages {
			line := []string{
				converters.IntToString(chatMessage.Frame),
				converters.IntToString(chatMessage.Tick),
				converters.IntToString(chatMessage.RoundNumber),
				converters.Uint64ToString(chatMessage.SenderSteamID64),
				chatMessage.SenderName,
				chatMessage.Message,
				converters.BoolToString(chatMessage.IsSenderAlive),
				converters.TeamToString(chatMessage.SenderSide),
				match.Checksum,
			}
			lines = append(lines, line)
		}

		csv.WriteLinesIntoCsvFile(outputPath+"_chat_messages.csv", lines)
	}

	// Hostages positions, 1 row per hostage and round, see packedPositions.
	// The state of a hostage is written only when it changes, stateTicks contains the tick of each change.
	var writeHostagePositions = func() {
		type hostageRoundKey struct {
			entityID    int
			serialNum   int
			roundNumber int
		}
		type hostageRoundPositions struct {
			roundNumber int
			positions   packedPositions
			stateTicks  []byte
			states      []byte
			lastTick    int
			lastState   common.HostageState
		}

		positionsPerHostage := make(map[hostageRoundKey]*hostageRoundPositions)
		var orderedPositions []*hostageRoundPositions
		for _, position := range match.HostagePositions {
			key := hostageRoundKey{
				entityID:    position.entityID,
				serialNum:   position.serialNum,
				roundNumber: position.RoundNumber,
			}
			hostage, exists := positionsPerHostage[key]
			if !exists {
				hostage = &hostageRoundPositions{roundNumber: position.RoundNumber}
				positionsPerHostage[key] = hostage
				orderedPositions = append(orderedPositions, hostage)
			}
			hostage.positions.add(position.Tick, position.X, position.Y, position.Z)
			// Like positions, only the first frame of a tick is used.
			if exists && position.Tick == hostage.lastTick {
				continue
			}
			hostage.lastTick = position.Tick
			if !exists || position.State != hostage.lastState {
				hostage.stateTicks = appendInt32(hostage.stateTicks, position.Tick)
				hostage.states = appendInt32(hostage.states, int(position.State))
				hostage.lastState = position.State
			}
		}

		lines := [][]string{}
		for _, hostage := range orderedPositions {
			hostage.positions.finish()
			line := []string{
				converters.IntToString(hostage.roundNumber),
				toBytea(hostage.positions.ticks),
				toBytea(hostage.positions.x),
				toBytea(hostage.positions.y),
				toBytea(hostage.positions.z),
				toBytea(hostage.stateTicks),
				toBytea(hostage.states),
				match.Checksum,
			}
			lines = append(lines, line)
		}

		csv.WriteLinesIntoCsvFile(outputPath+"_hostage_positions.csv", lines)
	}

	var writeHostagePickUpStart = func() {
		lines := [][]string{}
		for _, event := range match.HostagePickUpStart {
			line := []string{
				converters.IntToString(event.Frame),
				converters.IntToString(event.Tick),
				converters.IntToString(event.RoundNumber),
				converters.Uint64ToString(event.PlayerSteamID64),
				converters.BoolToString(event.IsPlayerControllingBot),
				converters.IntToString(event.HostageEntityId),
				converters.Float64ToString(event.X),
				converters.Float64ToString(event.Y),
				converters.Float64ToString(event.Z),
				match.Checksum,
			}
			lines = append(lines, line)
		}

		csv.WriteLinesIntoCsvFile(outputPath+"_hostage_pick_up_start.csv", lines)
	}

	var writeHostagePickedUp = func() {
		lines := [][]string{}
		for _, event := range match.HostagePickedUp {
			line := []string{
				converters.IntToString(event.Frame),
				converters.IntToString(event.Tick),
				converters.IntToString(event.RoundNumber),
				converters.Uint64ToString(event.PlayerSteamID64),
				converters.BoolToString(event.IsPlayerControllingBot),
				converters.IntToString(event.HostageEntityId),
				converters.Float64ToString(event.X),
				converters.Float64ToString(event.Y),
				converters.Float64ToString(event.Z),
				match.Checksum,
			}
			lines = append(lines, line)
		}

		csv.WriteLinesIntoCsvFile(outputPath+"_hostage_picked_up.csv", lines)
	}

	var writeHostageRescued = func() {
		lines := [][]string{}
		for _, event := range match.HostageRescued {
			line := []string{
				converters.IntToString(event.Frame),
				converters.IntToString(event.Tick),
				converters.IntToString(event.RoundNumber),
				converters.Uint64ToString(event.PlayerSteamID64),
				converters.BoolToString(event.IsPlayerControllingBot),
				converters.IntToString(event.HostageEntityId),
				converters.Float64ToString(event.X),
				converters.Float64ToString(event.Y),
				converters.Float64ToString(event.Z),
				match.Checksum,
			}
			lines = append(lines, line)
		}

		csv.WriteLinesIntoCsvFile(outputPath+"_hostage_rescued.csv", lines)
	}

	var functions = []func(){
		writeMatch,
		writeDemo,
		writeTeams,
		writePlayers,
		writePlayerPositions,
		writeShots,
		writeRounds,
		writeRoundEconomies,
		writeClutches,
		writeChickenDeaths,
		writeChickenPositions,
		writeDamages,
		writeKills,
		writeBombsPlanted,
		writeBombsDefuseStart,
		writeBombsDefused,
		writeBombsExploded,
		writeBombsPlantStart,
		writePlayersFlashed,
		writePlayersBuy,
		writeGrenadePositions,
		writeGrenadeBounces,
		writeGrenadeProjectilesDestroy,
		writeSmokesStart,
		writeHeGrenadesExplode,
		writeFlashbangsExplode,
		writeDecoysStart,
		writeInfernoPositions,
		writeChatMessages,
		writeHostagePositions,
		writeHostagePickUpStart,
		writeHostagePickedUp,
		writeHostageRescued,
	}

	var wg sync.WaitGroup

	for _, function := range functions {
		wg.Add(1)
		go func(function func()) {
			defer wg.Done()
			function()
		}(function)
	}

	wg.Wait()

	return nil
}
