package csgo

import (
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/common"
)

func newHostageRescued(analyzer *Analyzer, hostage *common.Hostage) *HostageRescued {
	parser := analyzer.parser

	return &HostageRescued{
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
