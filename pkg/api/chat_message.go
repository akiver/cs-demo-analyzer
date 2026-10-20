package api

import (
	"github.com/akiver/cs-demo-analyzer/pkg/api/constants"

	events "github.com/markus-wa/demoinfocs-golang/v6/pkg/demoinfocs/events"
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

func newChatMessageFromGameEvent(analyzer *Analyzer, event events.ChatMessage) *ChatMessage {
	parser := analyzer.parser
	return &ChatMessage{
		Frame:           parser.CurrentFrame(),
		Tick:            analyzer.currentTick(),
		RoundNumber:     analyzer.currentRound.Number,
		IsSenderAlive:   event.Sender.IsAlive(),
		Message:         event.Text,
		SenderName:      event.Sender.Name,
		SenderSteamID64: event.Sender.SteamID64,
		SenderSide:      constants.Side(event.Sender.Team),
	}
}
