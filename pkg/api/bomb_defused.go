package api

import (
	"github.com/akiver/cs-demo-analyzer/internal/converters"
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/common"
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/events"
)

type BombDefused struct {
	Frame                      int     `json:"frame"`
	Tick                       int     `json:"tick"`
	RoundNumber                int     `json:"roundNumber"`
	Site                       string  `json:"site"`
	DefuserSteamID64           uint64  `json:"defuserSteamId"`
	DefuserName                string  `json:"defuserName"`
	IsPlayerControllingBot     bool    `json:"isPlayerControllingBot"`
	X                          float64 `json:"x"`
	Y                          float64 `json:"y"`
	Z                          float64 `json:"z"`
	CounterTerroristAliveCount int     `json:"counterTerroristAliveCount"`
	TerroristAliveCount        int     `json:"terroristAliveCount"`
}

func newBombDefused(analyzer *Analyzer, event events.BombDefused) *BombDefused {
	parser := analyzer.parser
	player := event.Player

	counterTerroristAliveCount := 0
	terroristAliveCount := 0
	for _, player := range parser.GameState().Participants().Playing() {
		if !player.IsAlive() {
			continue
		}
		if player.Team == common.TeamCounterTerrorists {
			counterTerroristAliveCount++
		} else if player.Team == common.TeamTerrorists {
			terroristAliveCount++
		}
	}

	return &BombDefused{
		Frame:                      parser.CurrentFrame(),
		Tick:                       analyzer.currentTick(),
		RoundNumber:                analyzer.currentRound.Number,
		DefuserName:                player.Name,
		DefuserSteamID64:           player.SteamID64,
		IsPlayerControllingBot:     player.IsControllingBot(),
		Site:                       converters.BombsiteToString(event.BombEvent.Site),
		X:                          player.LastAlivePosition.X,
		Y:                          player.LastAlivePosition.Y,
		Z:                          player.LastAlivePosition.Z,
		CounterTerroristAliveCount: counterTerroristAliveCount,
		TerroristAliveCount:        terroristAliveCount,
	}
}
