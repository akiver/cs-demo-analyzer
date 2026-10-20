package model

import (
	"github.com/akiver/cs-demo-analyzer/pkg/api/constants"
)

type ChatMessage struct {
	Frame           int            `json:"frame"`
	Tick            int            `json:"tick"`
	RoundNumber     int            `json:"roundNumber"`
	Message         string         `json:"message"`
	SenderSteamID64 uint64         `json:"senderSteamId"`
	SenderName      string         `json:"senderName"`
	SenderSide      constants.Side `json:"senderSide"`
	IsSenderAlive   bool           `json:"isSenderAlive"`
}
