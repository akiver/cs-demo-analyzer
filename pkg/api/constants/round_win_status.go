package constants

type RoundWinStatus byte

const (
	RoundWinStatusUnassigned RoundWinStatus = 0
	RoundWinStatusDraw       RoundWinStatus = 1
	RoundWinStatusTWon       RoundWinStatus = RoundWinStatus(SideTerrorists)
	RoundWinStatusCTWon      RoundWinStatus = RoundWinStatus(SideCounterTerrorists)
)
