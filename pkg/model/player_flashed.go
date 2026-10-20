package model

import (
	"github.com/akiver/cs-demo-analyzer/pkg/api/constants"
)

type PlayerFlashed struct {
	Frame                   int            `json:"frame"`
	Tick                    int            `json:"tick"`
	RoundNumber             int            `json:"roundNumber"`
	Duration                float32        `json:"duration"`
	FlashedSteamID64        uint64         `json:"flashedSteamId"`
	FlashedName             string         `json:"flashedName"`
	FlashedSide             constants.Side `json:"flashedSide"`
	IsFlashedControllingBot bool           `json:"isFlashedControllingBot"`
	FlasherSteamID64        uint64         `json:"flasherSteamId"`
	FlasherName             string         `json:"flasherName"`
	FlasherSide             constants.Side `json:"flasherSide"`
	IsFlasherControllingBot bool           `json:"isFlasherControllingBot"`
}
