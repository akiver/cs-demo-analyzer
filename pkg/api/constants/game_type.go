package constants

import (
	"github.com/akiver/cs-demo-analyzer/internal/converters"
)

type GameType int

func (gameType GameType) String() string {
	return converters.IntToString(int(gameType))
}

const (
	GameTypeClassic     GameType = 0
	GameTypeGunGame     GameType = 1
	GameTypeTraining    GameType = 2
	GameTypeCustom      GameType = 3
	GameTypeCoOperative GameType = 4
	GameTypeSkirmish    GameType = 5
	GameTypeFFA         GameType = 6
)

type GameMode int

func (gameMode GameMode) String() string {
	return converters.IntToString(int(gameMode))
}

// When GameType is Classic
const (
	GameModeCasual       GameMode = 0
	GameModeCompetitive  GameMode = 1
	GameModeScrimmage2V2 GameMode = 2
	GameModeScrimmage5V5 GameMode = 3
)

// When GameType is GunGame
const (
	GameModeProgressive GameMode = 0
	GameModeBomb        GameMode = 1
	GameModeDeathmatch  GameMode = 2
)

// When GameType is Training
const (
	GameModeTraining GameMode = 0
)

// When GameType is Custom
const (
	GameModeCustom GameMode = 0
)

// When GameType is CoOperative
const (
	GameModeCoOperative        GameMode = 0
	GameModeCoOperativeMission GameMode = 1
)

// When GameType is Skirmish
const (
	GameModeSkirmish GameMode = 0
)

// When GameType is FFA
const (
	GameModeSurvival GameMode = 0
)

// Game mode as a string reported in CSVCMsg_ServerInfo messages.
type GameModeStr string

func (gameMode GameModeStr) String() string {
	return string(gameMode)
}

const (
	GameModeStrCasual             GameModeStr = "casual"
	GameModeStrPremier            GameModeStr = "premier"
	GameModeStrCompetitive        GameModeStr = "competitive"
	GameModeStrScrimmage2V2       GameModeStr = "scrimcomp2v2"
	GameModeStrScrimmage5v5       GameModeStr = "scrimcomp5v5"
	GameModeStrDeathmatch         GameModeStr = "deathmatch"
	GameModeStrGunGameProgressive GameModeStr = "gungameprogressive"
	GameModeStrGunGameBomb        GameModeStr = "gungametrbomb"
	GameModeStrCustom             GameModeStr = "custom"
	GameModeStrCoOperative        GameModeStr = "cooperative"
	GameModeStrCoOperativeMission GameModeStr = "coopmission"
	GameModeStrSkirmish           GameModeStr = "skirmish"
	GameModeStrSurvival           GameModeStr = "survival"
)

var GameModeMapping = map[GameType]map[GameMode]GameModeStr{
	GameTypeClassic: {
		GameModeCasual:       GameModeStrCasual,
		GameModeCompetitive:  GameModeStrCompetitive,
		GameModeScrimmage2V2: GameModeStrScrimmage2V2,
		GameModeScrimmage5V5: GameModeStrScrimmage5v5,
	},
	GameTypeGunGame: {
		GameModeProgressive: GameModeStrGunGameProgressive,
		GameModeBomb:        GameModeStrGunGameBomb,
		GameModeDeathmatch:  GameModeStrDeathmatch,
	},
	GameTypeCustom: {
		GameModeCustom: GameModeStrCustom,
	},
	GameTypeCoOperative: {
		GameModeCoOperative:        GameModeStrCoOperative,
		GameModeCoOperativeMission: GameModeStrCoOperativeMission,
	},
	GameTypeSkirmish: {
		GameModeSkirmish: GameModeStrSkirmish,
	},
	GameTypeFFA: {
		GameModeSurvival: GameModeStrSurvival,
	},
}
