package cs2

import (
	"github.com/akiver/cs-demo-analyzer/pkg/api/constants"

	events "github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/events"
)

func newPlayerBuy(analyzer *Analyzer, event events.ItemPickup) *PlayerBuy {
	parser := analyzer.parser

	return &PlayerBuy{
		Frame:           parser.CurrentFrame(),
		Tick:            analyzer.currentTick(),
		RoundNumber:     analyzer.currentRound.Number,
		PlayerSteamID64: event.Player.SteamID64,
		PlayerName:      event.Player.Name,
		PlayerSide:      constants.Side(event.Player.Team),
		WeaponName:      equipmentToWeaponName[event.Weapon.Type],
		WeaponType:      getEquipmentWeaponType(*event.Weapon),
		WeaponUniqueID:  event.Weapon.UniqueID().String(),
	}
}
