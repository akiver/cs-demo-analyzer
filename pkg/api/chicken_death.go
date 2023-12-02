package api

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

func newChickenDeath(frame int, tick int, roundNumber int, killerSteamID uint64, weaponName constants.WeaponName) *ChickenDeath {
	chickenDeath := &ChickenDeath{
		Frame:         frame,
		Tick:          tick,
		RoundNumber:   roundNumber,
		WeaponName:    weaponName,
		KillerSteamID: killerSteamID,
	}

	return chickenDeath
}
