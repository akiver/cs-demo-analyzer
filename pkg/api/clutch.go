package api

import (
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/common"
)

type Clutch struct {
	Frame             int         `json:"frame"`
	Tick              int         `json:"tick"`
	RoundNumber       int         `json:"roundNumber"`
	OpponentCount     int         `json:"opponentCount"`
	Side              common.Team `json:"side"`
	HasWon            bool        `json:"hasWon"`
	ClutcherSteamID64 uint64      `json:"clutcherSteamId"`
	ClutcherName      string      `json:"clutcherName"`
	ClutcherSurvived  bool        `json:"clutcherSurvived"`
	ClutcherKillCount int         `json:"clutcherKillCount"`
}

func newClutch(analyzer *Analyzer, clutcher *common.Player, side common.Team, opponentCount int) *Clutch {
	parser := analyzer.parser

	return &Clutch{
		Frame:             parser.CurrentFrame(),
		Tick:              analyzer.currentTick(),
		Side:              side,
		OpponentCount:     opponentCount,
		ClutcherName:      clutcher.Name,
		ClutcherSteamID64: clutcher.SteamID64,
		RoundNumber:       analyzer.currentRound.Number,
		ClutcherSurvived:  true,
	}
}
