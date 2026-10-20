package cs2

import (
	"github.com/akiver/cs-demo-analyzer/internal/converters"
	events "github.com/markus-wa/demoinfocs-golang/v6/pkg/demoinfocs/events"
)

func newBombExploded(analyzer *Analyzer, event events.BombExplode) *BombExploded {
	parser := analyzer.parser
	player := event.Player

	return &BombExploded{
		Frame:                  parser.CurrentFrame(),
		Tick:                   analyzer.currentTick(),
		RoundNumber:            analyzer.currentRound.Number,
		PlanterName:            player.Name,
		PlanterSteamID64:       player.SteamID64,
		IsPlayerControllingBot: player.IsControllingBot(),
		Site:                   converters.BombsiteToString(event.BombEvent.Site),
		X:                      analyzer.bombPlantPosition.X,
		Y:                      analyzer.bombPlantPosition.Y,
		Z:                      analyzer.bombPlantPosition.Z,
	}
}
