package api

import (
	"github.com/akiver/cs-demo-analyzer/internal/converters"
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/events"
)

type BombPlanted struct {
	Frame                  int     `json:"frame"`
	Tick                   int     `json:"tick"`
	RoundNumber            int     `json:"roundNumber"`
	Site                   string  `json:"site"`
	PlanterSteamID64       uint64  `json:"planterSteamId"`
	PlanterName            string  `json:"planterName"`
	IsPlayerControllingBot bool    `json:"isPlayerControllingBot"`
	X                      float64 `json:"x"`
	Y                      float64 `json:"y"`
	Z                      float64 `json:"z"`
}

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
