package csgo

import (
	"github.com/akiver/cs-demo-analyzer/internal/modelaccess"
	"github.com/akiver/cs-demo-analyzer/pkg/api/constants"
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/common"
)

func newHostagePositionFromHostage(analyzer *Analyzer, hostage *common.Hostage) *HostagePosition {
	parser := analyzer.parser

	position := &HostagePosition{
		Frame:       parser.CurrentFrame(),
		Tick:        analyzer.currentTick(),
		RoundNumber: analyzer.currentRound.Number,
		X:           hostage.Position().X,
		Y:           hostage.Position().Y,
		Z:           hostage.Position().Z,
		State:       constants.HostageState(hostage.State()),
	}
	modelaccess.SetHostageEntity(position, hostage.Entity.ID(), hostage.Entity.SerialNum())

	return position
}
