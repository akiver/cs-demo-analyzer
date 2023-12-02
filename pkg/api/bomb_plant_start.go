package api

import (
	"github.com/akiver/cs-demo-analyzer/internal/converters"
	events "github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/events"
)

type BombPlantStart struct {
	Frame                  int     `json:"frame"`
	Tick                   int     `json:"tick"`
	RoundNumber            int     `json:"roundNumber"`
	Site                   string  `json:"site"`
	PlanterSteamID64       uint64  `json:"defuserSteamId"`
	PlanterName            string  `json:"defuserName"`
	IsPlayerControllingBot bool    `json:"isPlayerControllingBot"`
	X                      float64 `json:"x"`
	Y                      float64 `json:"y"`
	Z                      float64 `json:"z"`
}

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
		X:                      player.LastAlivePosition.X,
		Y:                      player.LastAlivePosition.Y,
		Z:                      player.LastAlivePosition.Z,
	}
}
