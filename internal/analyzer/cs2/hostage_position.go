package cs2

import (
	"github.com/akiver/cs-demo-analyzer/pkg/api/constants"
	"github.com/markus-wa/demoinfocs-golang/v6/pkg/demoinfocs/common"
)

func newHostagePositionFromHostage(analyzer *Analyzer, hostage *common.Hostage) *HostagePosition {
	parser := analyzer.parser

	return &HostagePosition{
		EntityID:    hostage.Entity.ID(),
		SerialNum:   hostage.Entity.SerialNum(),
		Frame:       parser.CurrentFrame(),
		Tick:        analyzer.currentTick(),
		RoundNumber: analyzer.currentRound.Number,
		X:           hostage.Position().X,
		Y:           hostage.Position().Y,
		Z:           hostage.Position().Z,
		State:       constants.HostageState(hostage.State()),
	}
}
