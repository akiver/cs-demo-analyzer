package constants

import "github.com/markus-wa/demoinfocs-golang/v6/pkg/demoinfocs/common"

type RoundWinStatus byte

const (
	RoundWinStatusUnassigned RoundWinStatus = 0
	RoundWinStatusDraw       RoundWinStatus = 1
	RoundWinStatusTWon       RoundWinStatus = RoundWinStatus(common.TeamTerrorists)
	RoundWinStatusCTWon      RoundWinStatus = RoundWinStatus(common.TeamCounterTerrorists)
)
