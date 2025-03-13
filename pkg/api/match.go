package api

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/akiver/cs-demo-analyzer/internal/demo"
	"github.com/akiver/cs-demo-analyzer/internal/math"
	"github.com/akiver/cs-demo-analyzer/internal/slice"
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
	ShareCode                 string                      `json:"shareCode"` // Valve demos only, the .info file must be next to the .dem file to be able to generate it
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
		ShareCode:                 demoInfo.ShareCode,
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
		NetworkProtocol:           demoInfo.NetworkProtocol,
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

func (match *Match) resetRound(roundNumber int) {
	match.BombsPlantStart = slice.Filter(match.BombsPlantStart, func(event *BombPlantStart, index int) bool {
		return event.RoundNumber != roundNumber
	})
	match.BombsPlanted = slice.Filter(match.BombsPlanted, func(event *BombPlanted, index int) bool {
		return event.RoundNumber != roundNumber
	})
	match.BombsDefuseStart = slice.Filter(match.BombsDefuseStart, func(event *BombDefuseStart, index int) bool {
		return event.RoundNumber != roundNumber
	})
	match.BombsDefused = slice.Filter(match.BombsDefused, func(event *BombDefused, index int) bool {
		return event.RoundNumber != roundNumber
	})
	match.BombsExploded = slice.Filter(match.BombsExploded, func(event *BombExploded, index int) bool {
		return event.RoundNumber != roundNumber
	})
	match.ChatMessages = slice.Filter(match.ChatMessages, func(msg *ChatMessage, index int) bool {
		return msg.RoundNumber != roundNumber
	})
	match.ChickenDeaths = slice.Filter(match.ChickenDeaths, func(death *ChickenDeath, index int) bool {
		return death.RoundNumber != roundNumber
	})
	match.ChickenPositions = slice.Filter(match.ChickenPositions, func(position *ChickenPosition, index int) bool {
		return position.RoundNumber != roundNumber
	})
	match.Clutches = slice.Filter(match.Clutches, func(clutch *Clutch, index int) bool {
		return clutch.RoundNumber != roundNumber
	})
	match.Damages = slice.Filter(match.Damages, func(damage *Damage, index int) bool {
		return damage.RoundNumber != roundNumber
	})
	match.DecoysStart = slice.Filter(match.DecoysStart, func(decoy *DecoyStart, index int) bool {
		return decoy.RoundNumber != roundNumber
	})
	match.FlashbangsExplode = slice.Filter(match.FlashbangsExplode, func(flash *FlashbangExplode, index int) bool {
		return flash.RoundNumber != roundNumber
	})
	match.GrenadeBounces = slice.Filter(match.GrenadeBounces, func(grenade *GrenadeBounce, index int) bool {
		return grenade.RoundNumber != roundNumber
	})
	match.GrenadePositions = slice.Filter(match.GrenadePositions, func(position *GrenadePosition, index int) bool {
		return position.RoundNumber != roundNumber
	})
	match.GrenadeProjectilesDestroy = slice.Filter(match.GrenadeProjectilesDestroy, func(grenade *GrenadeProjectileDestroy, index int) bool {
		return grenade.RoundNumber != roundNumber
	})
	match.HeGrenadesExplode = slice.Filter(match.HeGrenadesExplode, func(grenade *HeGrenadeExplode, index int) bool {
		return grenade.RoundNumber != roundNumber
	})
	match.HostagePickUpStart = slice.Filter(match.HostagePickUpStart, func(hostage *HostagePickUpStart, index int) bool {
		return hostage.RoundNumber != roundNumber
	})
	match.HostagePickedUp = slice.Filter(match.HostagePickedUp, func(hostage *HostagePickedUp, index int) bool {
		return hostage.RoundNumber != roundNumber
	})
	match.HostagePositions = slice.Filter(match.HostagePositions, func(position *HostagePosition, index int) bool {
		return position.RoundNumber != roundNumber
	})
	match.HostageRescued = slice.Filter(match.HostageRescued, func(hostage *HostageRescued, index int) bool {
		return hostage.RoundNumber != roundNumber
	})
	match.InfernoPositions = slice.Filter(match.InfernoPositions, func(position *InfernoPosition, index int) bool {
		return position.RoundNumber != roundNumber
	})
	match.Kills = slice.Filter(match.Kills, func(kill *Kill, index int) bool {
		return kill.RoundNumber != roundNumber
	})
	match.PlayerEconomies = slice.Filter(match.PlayerEconomies, func(eco *PlayerEconomy, index int) bool {
		return eco.RoundNumber != roundNumber
	})
	match.PlayerPositions = slice.Filter(match.PlayerPositions, func(position *PlayerPosition, index int) bool {
		return position.RoundNumber != roundNumber
	})
	match.PlayersBuy = slice.Filter(match.PlayersBuy, func(event *PlayerBuy, index int) bool {
		return event.RoundNumber != roundNumber
	})
	match.PlayersFlashed = slice.Filter(match.PlayersFlashed, func(event *PlayerFlashed, index int) bool {
		return event.RoundNumber != roundNumber
	})
	match.Shots = slice.Filter(match.Shots, func(event *Shot, index int) bool {
		return event.RoundNumber != roundNumber
	})
	match.SmokesStart = slice.Filter(match.SmokesStart, func(event *SmokeStart, index int) bool {
		return event.RoundNumber != roundNumber
	})
}

func (match *Match) deleteIncompleteRounds() {
	for i := len(match.Rounds) - 1; i >= 0; i-- {
		round := match.Rounds[i]
		if round.WinnerName != "" {
			continue
		}

		match.Rounds = append(match.Rounds[:i], match.Rounds[i+1:]...)

		match.resetRound(round.Number)
	}
}
