package model

import (
	"encoding/json"

	"github.com/akiver/cs-demo-analyzer/pkg/api/constants"
)

type Round struct {
	match               *Match
	Number              int                      `json:"number"`
	StartTick           int                      `json:"startTick"`
	StartFrame          int                      `json:"startFrame"`
	FreezeTimeEndTick   int                      `json:"freezeTimeEndTick"`
	FreezeTimeEndFrame  int                      `json:"freezeTimeEndFrame"`
	EndTick             int                      `json:"endTick"`
	EndFrame            int                      `json:"endFrame"`
	EndOfficiallyTick   int                      `json:"endOfficiallyTick"`
	EndOfficiallyFrame  int                      `json:"endOfficiallyFrame"`
	OvertimeNumber      int                      `json:"overtimeNumber"`
	TeamAName           string                   `json:"teamAName"`
	TeamBName           string                   `json:"teamBName"`
	TeamAScore          int                      `json:"teamAScore"`
	TeamBScore          int                      `json:"teamBScore"`
	TeamASide           constants.Side           `json:"teamASide"`
	TeamBSide           constants.Side           `json:"teamBSide"`
	TeamAEquipmentValue int                      `json:"teamAEquipmentValue"`
	TeamBEquipmentValue int                      `json:"teamBEquipmentValue"`
	TeamAMoneySpent     int                      `json:"teamAMoneySpent"`
	TeamBMoneySpent     int                      `json:"teamBmoneySpent"`
	TeamAEconomyType    constants.EconomyType    `json:"teamAEconomyType"`
	TeamBEconomyType    constants.EconomyType    `json:"teamBEconomyType"`
	Duration            int64                    `json:"duration"`
	EndReason           constants.RoundEndReason `json:"endReason"`
	WinnerName          string                   `json:"winnerName"`
	WinnerSide          constants.Side           `json:"winnerSide"`
}

type RoundAlias Round

type RoundJSON struct {
	*RoundAlias
	StartMoneyTeamA int `json:"teamAStartMoney"`
	StartMoneyTeamB int `json:"teamBStartMoney"`
}

func (round *Round) MarshalJSON() ([]byte, error) {
	return json.Marshal(RoundJSON{
		RoundAlias:      (*RoundAlias)(round),
		StartMoneyTeamA: round.StartMoneyTeamA(),
		StartMoneyTeamB: round.StartMoneyTeamB(),
	})
}

func (round *Round) StartMoneyTeamA() int {
	return round.getTeamStartMoney(round.TeamASide)
}

func (round *Round) StartMoneyTeamB() int {
	return round.getTeamStartMoney(round.TeamBSide)
}

func (round *Round) getTeamStartMoney(side constants.Side) int {
	total := 0
	for _, economy := range round.match.PlayerEconomies {
		if economy.RoundNumber == round.Number && economy.PlayerSide == side {
			total += economy.StartMoney
		}
	}

	return total
}

// SetMatch sets the match the round belongs to, used for its computed values. Used by the analyzers.
func (round *Round) SetMatch(match *Match) {
	round.match = match
}
