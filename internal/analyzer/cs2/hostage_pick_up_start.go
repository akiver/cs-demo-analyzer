package cs2

import (
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/common"
)

func newHostagePickupStart(analyzer *Analyzer, player *common.Player, hostage *common.Hostage) *HostagePickUpStart {
	parser := analyzer.parser

	return &HostagePickUpStart{
		Frame:                  parser.CurrentFrame(),
		Tick:                   analyzer.currentTick(),
		RoundNumber:            analyzer.currentRound.Number,
		PlayerSteamID64:        player.SteamID64,
		IsPlayerControllingBot: player.IsControllingBot(),
		HostageEntityId:        hostage.Entity.ID(),
		X:                      hostage.Position().X,
		Y:                      hostage.Position().Y,
		Z:                      hostage.Position().Z,
	}
}
