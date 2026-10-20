package api

import (
	"github.com/akiver/cs-demo-analyzer/pkg/api/constants"
)

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
