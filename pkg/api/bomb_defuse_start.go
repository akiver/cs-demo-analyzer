package api

import (
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/common"
)

type BombDefuseStart struct {
	Frame                  int     `json:"frame"`
	Tick                   int     `json:"tick"`
	RoundNumber            int     `json:"roundNumber"`
	PlanterSteamID64       uint64  `json:"defuserSteamId"`
	PlanterName            string  `json:"defuserName"`
	IsPlayerControllingBot bool    `json:"isPlayerControllingBot"`
	X                      float64 `json:"x"`
	Y                      float64 `json:"y"`
	Z                      float64 `json:"z"`
}

func newBombDefuseStart(analyzer *Analyzer, player *common.Player) *BombDefuseStart {
	parser := analyzer.parser

	return &BombDefuseStart{
		Frame:                  parser.CurrentFrame(),
		Tick:                   analyzer.currentTick(),
		RoundNumber:            analyzer.currentRound.Number,
		PlanterName:            player.Name,
		PlanterSteamID64:       player.SteamID64,
		IsPlayerControllingBot: player.IsControllingBot(),
		X:                      player.LastAlivePosition.X,
		Y:                      player.LastAlivePosition.Y,
		Z:                      player.LastAlivePosition.Z,
	}
}
