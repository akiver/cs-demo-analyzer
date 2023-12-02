package fake

import (
	"github.com/akiver/cs-demo-analyzer/pkg/api/constants"
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/common"
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/events"
)

type FakeRound struct {
	Number              int
	StartTick           int
	StartFrame          int
	FreezeTimeEndTick   int
	FreezeTimeEndFrame  int
	EndTick             int
	EndFrame            int
	EndOfficiallyTick   int
	EndOfficiallyFrame  int
	OvertimeNumber      int
	TeamAName           string
	TeamBName           string
	TeamASide           common.Team
	TeamBSide           common.Team
	TeamAScore          int
	TeamBScore          int
	TeamAStartMoney     int
	TeamBStartMoney     int
	TeamAEquipmentValue int
	TeamBEquipmentValue int
	TeamAEconomyType    constants.EconomyType
	TeamBEconomyType    constants.EconomyType
	EndReason           events.RoundEndReason
	WinnerName          string
	WinnerSide          common.Team
}
