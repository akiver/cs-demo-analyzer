package model

import (
	"github.com/akiver/cs-demo-analyzer/pkg/api/constants"
)

type ChickenDeath struct {
	Frame         int                  `json:"frame"`
	Tick          int                  `json:"tick"`
	RoundNumber   int                  `json:"roundNumber"`
	KillerSteamID uint64               `json:"killerSteamId"`
	WeaponName    constants.WeaponName `json:"weaponName"`
}
