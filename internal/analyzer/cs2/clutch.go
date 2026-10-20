package cs2

import (
	"github.com/akiver/cs-demo-analyzer/pkg/api/constants"
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/common"
)

func newClutch(analyzer *Analyzer, clutcher *common.Player, side common.Team, opponentCount int) *Clutch {
	parser := analyzer.parser

	return &Clutch{
		Frame:             parser.CurrentFrame(),
		Tick:              analyzer.currentTick(),
		Side:              constants.Side(side),
		OpponentCount:     opponentCount,
		ClutcherName:      clutcher.Name,
		ClutcherSteamID64: clutcher.SteamID64,
		RoundNumber:       analyzer.currentRound.Number,
		ClutcherSurvived:  true,
	}
}
