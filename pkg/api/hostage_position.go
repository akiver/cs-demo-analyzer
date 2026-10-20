package api

import (
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/common"
)

type HostagePosition struct {
	// Used to identify the hostage, entity IDs are re-used once an entity is destroyed but not with the same serial number.
	// Not exported because they are used only to group positions by hostage.
	entityID    int
	serialNum   int
	Frame       int                 `json:"frame"`
	Tick        int                 `json:"tick"`
	RoundNumber int                 `json:"roundNumber"`
	X           float64             `json:"x"`
	Y           float64             `json:"y"`
	Z           float64             `json:"z"`
	State       common.HostageState `json:"state"`
}

func newHostagePositionFromHostage(analyzer *Analyzer, hostage *common.Hostage) *HostagePosition {
	parser := analyzer.parser

	return &HostagePosition{
		entityID:    hostage.Entity.ID(),
		serialNum:   hostage.Entity.SerialNum(),
		Frame:       parser.CurrentFrame(),
		Tick:        analyzer.currentTick(),
		RoundNumber: analyzer.currentRound.Number,
		X:           hostage.Position().X,
		Y:           hostage.Position().Y,
		Z:           hostage.Position().Z,
		State:       hostage.State(),
	}
}
