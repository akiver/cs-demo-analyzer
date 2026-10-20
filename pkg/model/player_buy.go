package model

import (
	"github.com/akiver/cs-demo-analyzer/pkg/api/constants"
)

type PlayerBuy struct {
	Frame           int                  `json:"frame"`
	Tick            int                  `json:"tick"`
	RoundNumber     int                  `json:"roundNumber"`
	PlayerSteamID64 uint64               `json:"playerSteamId"`
	PlayerSide      constants.Side       `json:"playerSide"`
	PlayerName      string               `json:"playerName"`
	WeaponName      constants.WeaponName `json:"weaponName"`
	WeaponType      constants.WeaponType `json:"weaponType"`
	WeaponUniqueID  string               `json:"weaponUniqueId"`
	HasRefunded     bool                 `json:"hasRefunded"`
}
