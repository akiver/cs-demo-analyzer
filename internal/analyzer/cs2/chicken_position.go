package cs2

import (
	"github.com/akiver/cs-demo-analyzer/internal/modelaccess"
	"github.com/markus-wa/demoinfocs-golang/v6/pkg/demoinfocs/sendtables"
)

type chickenEntity struct {
	id     int
	entity *sendtables.Entity
}

func newChickenPositionFromEntity(analyzer *Analyzer, chicken *chickenEntity) *ChickenPosition {
	parser := analyzer.parser
	entity := chicken.entity

	position := &ChickenPosition{
		Frame:       parser.CurrentFrame(),
		Tick:        analyzer.currentTick(),
		RoundNumber: analyzer.currentRound.Number,
		X:           entity.Position().X,
		Y:           entity.Position().Y,
		Z:           entity.Position().Z,
	}
	modelaccess.SetChickenID(position, chicken.id)

	return position
}
