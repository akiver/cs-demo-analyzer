package cs2

import (
	"github.com/akiver/cs-demo-analyzer/internal/converters"
	events "github.com/markus-wa/demoinfocs-golang/v6/pkg/demoinfocs/events"
)

func newBombPlantStart(analyzer *Analyzer, event events.BombPlantBegin) *BombPlantStart {
	parser := analyzer.parser
	player := event.Player

	return &BombPlantStart{
		Frame:                  parser.CurrentFrame(),
		Tick:                   analyzer.currentTick(),
		RoundNumber:            analyzer.currentRound.Number,
		PlanterName:            player.Name,
		PlanterSteamID64:       player.SteamID64,
		IsPlayerControllingBot: player.IsControllingBot(),
		Site:                   converters.BombsiteToString(event.BombEvent.Site),
		X:                      player.Position().X,
		Y:                      player.Position().Y,
		Z:                      player.Position().Z,
	}
}
