package cs2

import (
	"github.com/akiver/cs-demo-analyzer/internal/converters"
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/events"
)

func newBombPlanted(analyzer *Analyzer, event events.BombPlanted) *BombPlanted {
	parser := analyzer.parser
	player := event.Player

	return &BombPlanted{
		Frame:                  parser.CurrentFrame(),
		Tick:                   analyzer.currentTick(),
		RoundNumber:            analyzer.currentRound.Number,
		PlanterName:            player.Name,
		PlanterSteamID64:       player.SteamID64,
		IsPlayerControllingBot: player.IsControllingBot(),
		Site:                   converters.BombsiteToString(event.BombEvent.Site),
		X:                      player.LastAlivePosition.X,
		Y:                      player.LastAlivePosition.Y,
		Z:                      player.LastAlivePosition.Z,
	}
}
