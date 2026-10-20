package model

import (
	"github.com/akiver/cs-demo-analyzer/pkg/api/constants"
)

type Clutch struct {
	Frame             int            `json:"frame"`
	Tick              int            `json:"tick"`
	RoundNumber       int            `json:"roundNumber"`
	OpponentCount     int            `json:"opponentCount"`
	Side              constants.Side `json:"side"`
	HasWon            bool           `json:"hasWon"`
	ClutcherSteamID64 uint64         `json:"clutcherSteamId"`
	ClutcherName      string         `json:"clutcherName"`
	ClutcherSurvived  bool           `json:"clutcherSurvived"`
	ClutcherKillCount int            `json:"clutcherKillCount"`
}
