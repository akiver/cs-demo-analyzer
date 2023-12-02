package api

import (
	"github.com/akiver/cs-demo-analyzer/pkg/api/constants"
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/common"
	events "github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/events"
)

type PlayerBuy struct {
	Frame           int                  `json:"frame"`
	Tick            int                  `json:"tick"`
	RoundNumber     int                  `json:"roundNumber"`
	PlayerSteamID64 uint64               `json:"playerSteamId"`
	PlayerSide      common.Team          `json:"playerSide"`
	PlayerName      string               `json:"playerName"`
	WeaponName      constants.WeaponName `json:"weaponName"`
	WeaponType      constants.WeaponType `json:"weaponType"`
	WeaponUniqueID  string               `json:"weaponUniqueId"`
	HasRefunded     bool                 `json:"hasRefunded"`
}

func newPlayerBuy(analyzer *Analyzer, event events.ItemPickup) *PlayerBuy {
	parser := analyzer.parser

	return &PlayerBuy{
		Frame:           parser.CurrentFrame(),
		Tick:            analyzer.currentTick(),
		RoundNumber:     analyzer.currentRound.Number,
		PlayerSteamID64: event.Player.SteamID64,
		PlayerName:      event.Player.Name,
		PlayerSide:      event.Player.Team,
		WeaponName:      equipmentToWeaponName[event.Weapon.Type],
		WeaponType:      getEquipmentWeaponType(*event.Weapon),
		WeaponUniqueID:  event.Weapon.UniqueID2().String(),
	}
}
