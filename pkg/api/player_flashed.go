package api

import (
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/common"
	events "github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/events"
)

type PlayerFlashed struct {
	Frame                   int         `json:"frame"`
	Tick                    int         `json:"tick"`
	RoundNumber             int         `json:"roundNumber"`
	Duration                float32     `json:"duration"`
	FlashedSteamID64        uint64      `json:"flashedSteamId"`
	FlashedName             string      `json:"flashedName"`
	FlashedSide             common.Team `json:"flashedSide"`
	IsFlashedControllingBot bool        `json:"isFlashedControllingBot"`
	FlasherSteamID64        uint64      `json:"flasherSteamId"`
	FlasherName             string      `json:"flasherName"`
	FlasherSide             common.Team `json:"flasherSide"`
	IsFlasherControllingBot bool        `json:"isFlasherControllingBot"`
}

func newPlayerFlashed(analyzer *Analyzer, event events.PlayerFlashed) *PlayerFlashed {
	parser := analyzer.parser

	return &PlayerFlashed{
		Frame:                   parser.CurrentFrame(),
		Tick:                    analyzer.currentTick(),
		RoundNumber:             analyzer.currentRound.Number,
		Duration:                event.Player.FlashDuration,
		FlashedName:             event.Player.Name,
		FlashedSteamID64:        event.Player.SteamID64,
		FlashedSide:             event.Player.Team,
		IsFlashedControllingBot: event.Player.IsControllingBot(),
		FlasherName:             event.Attacker.Name,
		FlasherSteamID64:        event.Attacker.SteamID64,
		FlasherSide:             event.Attacker.Team,
		IsFlasherControllingBot: event.Attacker.IsControllingBot(),
	}
}
