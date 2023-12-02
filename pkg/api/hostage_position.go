package api

import (
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/common"
)

type HostagePosition struct {
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
		Frame:       parser.CurrentFrame(),
		Tick:        analyzer.currentTick(),
		RoundNumber: analyzer.currentRound.Number,
		X:           hostage.Position().X,
		Y:           hostage.Position().Y,
		Z:           hostage.Position().Z,
		State:       hostage.State(),
	}
}
