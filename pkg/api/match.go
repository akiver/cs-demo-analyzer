package api

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/akiver/cs-demo-analyzer/internal/demo"
	"github.com/akiver/cs-demo-analyzer/internal/math"
	"github.com/akiver/cs-demo-analyzer/pkg/api/constants"
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/common"
)

// Match is the root struct that contains relevant data from a demo.
// It excludes data from warmup / halftime / after match.
type Match struct {
	Checksum                  string                      `json:"checksum"`
	Game                      constants.Game              `json:"game"`
	DemoFilePath              string                      `json:"demoFilePath"`
	DemoFileName              string                      `json:"demoFileName"`
	Source                    constants.DemoSource        `json:"source"`
	Type                      constants.DemoType          `json:"type"`
	MapName                   string                      `json:"mapName"`
	TickCount                 int                         `json:"tickCount"`
	TickRate                  float64                     `json:"tickrate"`
	FrameRate                 float64                     `json:"framerate"`
	Date                      time.Time                   `json:"date"`
	Duration                  time.Duration               `json:"duration"`
	ServerName                string                      `json:"serverName"`
	ClientName                string                      `json:"clientName"`
	NetworkProtocol           int                         `json:"networkProtocol"`
	BuildNumber               int                         `json:"buildNumber"` // CS2 only
	GameType                  constants.GameType          `json:"gameType"`
	GameMode                  constants.GameMode          `json:"gameMode"`
	gameModeStr               constants.GameModeStr       // CS2 only, the game mode as a string coming from the CSVCMsg_ServerInfo msg
	IsRanked                  bool                        `json:"isRanked"`
	MaxRounds                 int                         `json:"maxRounds"` // mp_maxrounds if detected or based on final scores
	OvertimeCount             int                         `json:"overtimeCount"`
	HasVacLiveBan             bool                        `json:"hasVacLiveBan"`
	TeamA                     *Team                       `json:"teamA"` // Team A is the Team that started as CT
	TeamB                     *Team                       `json:"teamB"` // Team B is the Team that started as T
	Winner                    *Team                       `json:"winner"`
	PlayersBySteamID          map[uint64]*Player          `json:"players"`
	Kills                     []*Kill                     `json:"kills"` // Includes suicides and bomb explosions too
	Shots                     []*Shot                     `json:"shots"`
	Rounds                    []*Round                    `json:"rounds"`
	Clutches                  []*Clutch                   `json:"clutches"`
	BombsPlanted              []*BombPlanted              `json:"bombsPlanted"`
	BombsDefused              []*BombDefused              `json:"bombsDefused"`
	BombsExploded             []*BombExploded             `json:"bombsExploded"`
	BombsPlantStart           []*BombPlantStart           `json:"bombsPlantStart"`
	BombsDefuseStart          []*BombDefuseStart          `json:"bombsDefuseStart"`
	PlayersFlashed            []*PlayerFlashed            `json:"playersFlashed"`
	GrenadePositions          []*GrenadePosition          `json:"grenadePositions"`
	InfernoPositions          []*InfernoPosition          `json:"infernoPositions"`
	HostagePickUpStart        []*HostagePickUpStart       `json:"hostagePickUpStart"`
	HostagePickedUp           []*HostagePickedUp          `json:"hostagePickedUp"`
	HostageRescued            []*HostageRescued           `json:"hostageRescued"`
	HostagePositions          []*HostagePosition          `json:"hostagePositions"`
	SmokesStart               []*SmokeStart               `json:"smokesStart"`
	DecoysStart               []*DecoyStart               `json:"decoysStart"`
	HeGrenadesExplode         []*HeGrenadeExplode         `json:"heGrenadesExplode"`
	FlashbangsExplode         []*FlashbangExplode         `json:"flashbangsExplode"`
	GrenadeBounces            []*GrenadeBounce            `json:"grenadeBounces"`
	GrenadeProjectilesDestroy []*GrenadeProjectileDestroy `json:"grenadeProjectilesDestroy"`
	ChickenPositions          []*ChickenPosition          `json:"chickenPositions"`
	ChickenDeaths             []*ChickenDeath             `json:"chickenDeaths"`
	Damages                   []*Damage                   `json:"damages"`
	PlayerPositions           []*PlayerPosition           `json:"playerPositions"`
	PlayersBuy                []*PlayerBuy                `json:"playersBuy"`
	PlayerEconomies           []*PlayerEconomy            `json:"playerEconomies"`
	ChatMessages              []*ChatMessage              `json:"chatMessages"`
	scoreTeamA                *int
	scoreTeamB                *int
}

type MatchAlias Match

type MatchJSON struct {
	*MatchAlias
	GameModeStr string `json:"gameModeStr"`
}

func (match *Match) MarshalJSON() ([]byte, error) {

	return json.Marshal(MatchJSON{
		MatchAlias:  (*MatchAlias)(match),
		GameModeStr: match.GameModeStr().String(),
	})
}

func (match Match) GameModeStr() constants.GameModeStr {
	if match.gameModeStr != "" {
		return match.gameModeStr
	}

	gameModeStr := constants.GameModeMapping[match.GameType][match.GameMode]
	if gameModeStr != "" {
		return gameModeStr
	}

	fmt.Printf("Unknown game mode string for game type: %d and game mode: %d", match.GameType, match.GameMode)

	return constants.GameModeStrCompetitive
}

func (match Match) Players() []*Player {
	players := make([]*Player, 0, len(match.PlayersBySteamID))
	for _, player := range match.PlayersBySteamID {
		players = append(players, player)
	}

	return players
}

// This returns the team currently on the given side.
func (match Match) Team(team common.Team) *Team {
	if *match.TeamA.CurrentSide == team {
		return match.TeamA
	}

	return match.TeamB
}

// This returns the number of kills in the match.
// It doesn't include suicides and deaths by bomb explosion.
func (match Match) KillCount() int {
	killCount := 0
	for _, kill := range match.Kills {
		if kill.KillerSteamID64 != 0 && kill.KillerName != kill.VictimName && kill.KillerName != common.EqWorld.String() {
			killCount++
		}
	}
	return killCount
}

// This returns the number of assists excluding BOT assists.
func (match Match) AssistCount() int {
	assistCount := 0
	for _, kill := range match.Kills {
		if kill.AssisterSteamID64 != 0 {
			assistCount++
		}
	}
	return assistCount
}

// This returns the number of player deaths in the match.
// Unlike KillCount, it takes into account ALL deaths whatever the reason is (weapon, suicide, bomb explosion...).
func (match Match) DeathCount() int {
	return len(match.Kills)
}

func (match Match) ShotCount() int {
	return len(match.Shots)
}

func (match *Match) KillsByRound() map[int][]*Kill {
	killsByRound := make(map[int][]*Kill)
	for _, kill := range match.Kills {
		killsByRound[kill.RoundNumber] = append(killsByRound[kill.RoundNumber], kill)
	}

	return killsByRound
}

func (match *Match) GetPlayerEconomyAtRound(playerName string, steamID64 uint64, roundNumber int) *PlayerEconomy {
	for _, economy := range match.PlayerEconomies {
		if economy.RoundNumber == roundNumber && economy.SteamID64 == steamID64 && economy.Name == playerName {
			return economy
		}
	}

	return nil
}

func (match Match) swapTeams() {
	match.TeamA.swap()
	match.TeamB.swap()
}

func (match *Match) computeResultStats() {
	// Compute the teams score based on rounds result.
	// Listening on CCSTeam.m_scoreTotal or the event ScoreUpdated is not reliable with non Valve demos.
	isFirstHalf := true
	for roundIndex, round := range match.Rounds {
		if match.FrameRate > 0 {
			round.Duration = int64(((float64(round.EndFrame) - float64(round.StartFrame)) / match.FrameRate) * 1000)
		}

		if round.WinnerSide == round.TeamASide {
			match.TeamA.Score++
			round.TeamAScore = match.TeamA.Score
			round.TeamBScore = match.TeamB.Score
		} else {
			match.TeamB.Score++
			round.TeamAScore = match.TeamA.Score
			round.TeamBScore = match.TeamB.Score
		}

		if round.OvertimeNumber == 0 {
			if isFirstHalf {
				if round.WinnerSide == round.TeamASide {
					match.TeamA.ScoreFirstHalf++
				} else {
					match.TeamB.ScoreFirstHalf++
				}
			} else {
				if round.WinnerSide == round.TeamASide {
					match.TeamA.ScoreSecondHalf++
				} else {
					match.TeamB.ScoreSecondHalf++
				}
			}
		}

		if roundIndex < len(match.Rounds)-1 {
			nextRound := match.Rounds[roundIndex+1]
			if round.TeamASide != nextRound.TeamASide {
				isFirstHalf = !isFirstHalf
			}
		}
	}

	if match.TeamA.Score > match.TeamB.Score {
		match.Winner = match.TeamA
	} else if match.TeamB.Score > match.TeamA.Score {
		match.Winner = match.TeamB
	}

	if match.MaxRounds == 0 {
		maxScore := math.Max(match.TeamA.Score, match.TeamB.Score)
		match.MaxRounds = maxScore * 2
		if match.TeamA.Score != match.TeamB.Score {
			match.MaxRounds -= 2
		}
	}
}

func (match *Match) initTeams() {
	initialSideTeamA := common.TeamCounterTerrorists
	initialSideTeamB := common.TeamTerrorists
	teamA := &Team{
		Name:        "Team A",
		Letter:      "A",
		CurrentSide: &initialSideTeamA,
	}
	match.TeamA = teamA
	match.scoreTeamA = &teamA.Score

	teamB := &Team{
		Name:        "Team B",
		Letter:      "B",
		CurrentSide: &initialSideTeamB,
	}
	match.TeamB = teamB
	match.scoreTeamB = &teamB.Score
}

func newMatch(source constants.DemoSource, demoInfo *demo.Demo) Match {
	game := constants.CSGO
	if demoInfo.IsSource2() {
		// The build number of CS2 when it was publicly available is 9832, everything below is coming from the limited test.
		if demoInfo.BuildNumber < 9832 {
			game = constants.CS2LT
		} else {
			game = constants.CS2
		}
	}

	match := Match{
		Checksum:                  demoInfo.Checksum,
		Source:                    source,
		Game:                      game,
		Type:                      "GOTV", // By default assume it's a GOTV demo, it will be updated during parsing.
		TickCount:                 demoInfo.TickCount,
		Date:                      demoInfo.Date,
		DemoFilePath:              demoInfo.FilePath,
		DemoFileName:              demoInfo.FileName,
		TickRate:                  demoInfo.TickRate,
		MapName:                   demoInfo.MapName,
		FrameRate:                 demoInfo.FrameRate,
		Duration:                  demoInfo.Duration,
		ServerName:                demoInfo.ServerName,
		ClientName:                demoInfo.ClientName,
		BuildNumber:               demoInfo.BuildNumber,
		PlayersBySteamID:          make(map[uint64]*Player),
		Rounds:                    []*Round{},
		Kills:                     []*Kill{},
		Clutches:                  []*Clutch{},
		BombsPlanted:              []*BombPlanted{},
		BombsDefused:              []*BombDefused{},
		BombsExploded:             []*BombExploded{},
		BombsPlantStart:           []*BombPlantStart{},
		BombsDefuseStart:          []*BombDefuseStart{},
		Damages:                   []*Damage{},
		Shots:                     []*Shot{},
		PlayersFlashed:            []*PlayerFlashed{},
		PlayersBuy:                []*PlayerBuy{},
		PlayerEconomies:           []*PlayerEconomy{},
		PlayerPositions:           []*PlayerPosition{},
		GrenadePositions:          []*GrenadePosition{},
		GrenadeBounces:            []*GrenadeBounce{},
		InfernoPositions:          []*InfernoPosition{},
		SmokesStart:               []*SmokeStart{},
		DecoysStart:               []*DecoyStart{},
		HeGrenadesExplode:         []*HeGrenadeExplode{},
		FlashbangsExplode:         []*FlashbangExplode{},
		HostagePickUpStart:        []*HostagePickUpStart{},
		HostagePickedUp:           []*HostagePickedUp{},
		HostageRescued:            []*HostageRescued{},
		HostagePositions:          []*HostagePosition{},
		ChickenPositions:          []*ChickenPosition{},
		ChatMessages:              []*ChatMessage{},
		ChickenDeaths:             []*ChickenDeath{},
		GrenadeProjectilesDestroy: []*GrenadeProjectileDestroy{},
	}

	match.initTeams()

	return match
}

func (match *Match) reset() {
	match.PlayersBySteamID = make(map[uint64]*Player)
	match.Rounds = []*Round{}
	match.Kills = []*Kill{}
	match.Clutches = []*Clutch{}
	match.BombsPlanted = []*BombPlanted{}
	match.BombsDefused = []*BombDefused{}
	match.BombsExploded = []*BombExploded{}
	match.BombsPlantStart = []*BombPlantStart{}
	match.BombsDefuseStart = []*BombDefuseStart{}
	match.Damages = []*Damage{}
	match.Shots = []*Shot{}
	match.PlayersFlashed = []*PlayerFlashed{}
	match.PlayersBuy = []*PlayerBuy{}
	match.PlayerPositions = []*PlayerPosition{}
	match.GrenadePositions = []*GrenadePosition{}
	match.GrenadeBounces = []*GrenadeBounce{}
	match.InfernoPositions = []*InfernoPosition{}
	match.SmokesStart = []*SmokeStart{}
	match.DecoysStart = []*DecoyStart{}
	match.HeGrenadesExplode = []*HeGrenadeExplode{}
	match.FlashbangsExplode = []*FlashbangExplode{}
	match.HostagePickUpStart = []*HostagePickUpStart{}
	match.HostagePickedUp = []*HostagePickedUp{}
	match.HostageRescued = []*HostageRescued{}
	match.HostagePositions = []*HostagePosition{}
	match.ChickenPositions = []*ChickenPosition{}
	match.ChatMessages = []*ChatMessage{}
	match.ChickenDeaths = []*ChickenDeath{}
	match.GrenadeProjectilesDestroy = []*GrenadeProjectileDestroy{}
	match.PlayerEconomies = []*PlayerEconomy{}
	match.initTeams()
}
