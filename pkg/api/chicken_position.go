package api

import (
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/sendtables"
)

type ChickenPosition struct {
	Frame       int     `json:"frame"`
	Tick        int     `json:"tick"`
	RoundNumber int     `json:"roundNumber"`
	X           float64 `json:"x"`
	Y           float64 `json:"y"`
	Z           float64 `json:"z"`
}

func newChickenPositionFromEntity(analyzer *Analyzer, entity sendtables.Entity) *ChickenPosition {
	parser := analyzer.parser

	return &ChickenPosition{
		Frame:       parser.CurrentFrame(),
		Tick:        analyzer.currentTick(),
		RoundNumber: analyzer.currentRound.Number,
		X:           entity.Position().X,
		Y:           entity.Position().Y,
		Z:           entity.Position().Z,
	}
}
