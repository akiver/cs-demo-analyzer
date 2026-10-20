package csgo

import (
	"github.com/akiver/cs-demo-analyzer/pkg/api/constants"

	events "github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/events"
)

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
