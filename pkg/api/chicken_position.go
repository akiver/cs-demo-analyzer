package api

import (
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/sendtables"
)

type chickenEntity struct {
	id     int
	entity sendtables.Entity
}

type ChickenPosition struct {
	// Unique ID of the chicken in the match, not exported because it's used only to group positions by chicken.
	chickenID   int
	Frame       int     `json:"frame"`
	Tick        int     `json:"tick"`
	RoundNumber int     `json:"roundNumber"`
	X           float64 `json:"x"`
	Y           float64 `json:"y"`
	Z           float64 `json:"z"`
}

func newChickenPositionFromEntity(analyzer *Analyzer, chicken *chickenEntity) *ChickenPosition {
	parser := analyzer.parser
	entity := chicken.entity

	return &ChickenPosition{
		chickenID:   chicken.id,
		Frame:       parser.CurrentFrame(),
		Tick:        analyzer.currentTick(),
		RoundNumber: analyzer.currentRound.Number,
		X:           entity.Position().X,
		Y:           entity.Position().Y,
		Z:           entity.Position().Z,
	}
}
