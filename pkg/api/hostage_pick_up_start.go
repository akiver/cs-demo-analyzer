package api

import (
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/common"
)

type HostagePickUpStart struct {
	Frame                  int     `json:"frame"`
	Tick                   int     `json:"tick"`
	RoundNumber            int     `json:"roundNumber"`
	PlayerSteamID64        uint64  `json:"playerSteamId"`
	IsPlayerControllingBot bool    `json:"isPlayerControllingBot"`
	HostageEntityId        int     `json:"hostageEntityId"`
	X                      float64 `json:"x"`
	Y                      float64 `json:"y"`
	Z                      float64 `json:"z"`
}

func newHostagePickupStart(analyzer *Analyzer, player *common.Player, hostage *common.Hostage) *HostagePickUpStart {
	parser := analyzer.parser

	return &HostagePickUpStart{
		Frame:                  parser.CurrentFrame(),
		Tick:                   analyzer.currentTick(),
		RoundNumber:            analyzer.currentRound.Number,
		PlayerSteamID64:        player.SteamID64,
		IsPlayerControllingBot: player.IsControllingBot(),
		HostageEntityId:        hostage.Entity.ID(),
		X:                      hostage.Position().X,
		Y:                      hostage.Position().Y,
		Z:                      hostage.Position().Z,
	}
}
