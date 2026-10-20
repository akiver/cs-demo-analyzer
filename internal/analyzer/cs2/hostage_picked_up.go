package cs2

import (
	"github.com/markus-wa/demoinfocs-golang/v6/pkg/demoinfocs/common"
)

func newHostagePickedUp(analyzer *Analyzer, hostage *common.Hostage) *HostagePickedUp {
	parser := analyzer.parser

	return &HostagePickedUp{
		Frame:                  parser.CurrentFrame(),
		Tick:                   analyzer.currentTick(),
		RoundNumber:            analyzer.currentRound.Number,
		PlayerSteamID64:        hostage.Leader().SteamID64,
		IsPlayerControllingBot: hostage.Leader().IsControllingBot(),
		HostageEntityId:        hostage.Entity.ID(),
		X:                      hostage.Position().X,
		Y:                      hostage.Position().Y,
		Z:                      hostage.Position().Z,
	}

}
