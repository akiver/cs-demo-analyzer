package api

import (
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/common"
)

type HostagePickedUp struct {
	Frame                  int     `json:"frame"`
	Tick                   int     `json:"tick"`
	RoundNumber            int     `json:"roundNumber"`
	HostageEntityId        int     `json:"hostageEntityId"`
	PlayerSteamID64        uint64  `json:"playerSteamId"`
	IsPlayerControllingBot bool    `json:"isPlayerControllingBot"`
	X                      float64 `json:"x"`
	Y                      float64 `json:"y"`
	Z                      float64 `json:"z"`
}

func newHostagePickedUp(analyzer *Analyzer, hostage *common.Hostage) *HostagePickedUp {
	parser := analyzer.parser

	return &HostagePickedUp{
		Frame:                  parser.CurrentFrame(),
		Tick:                   analyzer.currentTick(),
		RoundNumber:            analyzer.currentRound.Number,
		PlayerSteamID64:        hostage.Leader().SteamID64,
		IsPlayerControllingBot: hostage.Leader().IsControllingBot(),
		HostageEntityId:        hostage.Entity.ID(),
		X:                      hostage.Position().X,
		Y:                      hostage.Position().Y,
		Z:                      hostage.Position().Z,
	}

}
