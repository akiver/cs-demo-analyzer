package cs2

import (
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/common"
)

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
