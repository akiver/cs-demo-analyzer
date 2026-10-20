package csgo

import (
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/sendtables"
)

type chickenEntity struct {
	id     int
	entity sendtables.Entity
}

func newChickenPositionFromEntity(analyzer *Analyzer, chicken *chickenEntity) *ChickenPosition {
	parser := analyzer.parser
	entity := chicken.entity

	return &ChickenPosition{
		ChickenID:   chicken.id,
		Frame:       parser.CurrentFrame(),
		Tick:        analyzer.currentTick(),
		RoundNumber: analyzer.currentRound.Number,
		X:           entity.Position().X,
		Y:           entity.Position().Y,
		Z:           entity.Position().Z,
	}
}
