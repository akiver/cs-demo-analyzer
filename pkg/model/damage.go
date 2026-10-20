package model

import (
	"github.com/akiver/cs-demo-analyzer/pkg/api/constants"
)

type Damage struct {
	Frame                    int                  `json:"frame"`
	Tick                     int                  `json:"tick"`
	RoundNumber              int                  `json:"roundNumber"`
	HealthDamage             int                  `json:"healthDamage"`
	ArmorDamage              int                  `json:"armorDamage"`
	AttackerSteamID64        uint64               `json:"attackerSteamId"`
	AttackerSide             constants.Side       `json:"attackerSide"`
	AttackerTeamName         string               `json:"attackerTeamName"`
	IsAttackerControllingBot bool                 `json:"isAttackerControllingBot"`
	VictimHealth             int                  `json:"victimHealth"`
	VictimNewHealth          int                  `json:"victimNewHealth"`
	VictimArmor              int                  `json:"victimArmor"`
	VictimNewArmor           int                  `json:"victimNewArmor"`
	VictimSteamID64          uint64               `json:"victimSteamId"`
	VictimSide               constants.Side       `json:"victimSide"`
	VictimTeamName           string               `json:"victimTeamName"`
	IsVictimControllingBot   bool                 `json:"isVictimControllingBot"`
	HitGroup                 constants.HitGroup   `json:"hitgroup"`
	WeaponName               constants.WeaponName `json:"weaponName"`
	WeaponType               constants.WeaponType `json:"weaponType"`
	WeaponUniqueID           string               `json:"weaponUniqueId"`
}

func (damage *Damage) IsGrenadeWeapon() bool {
	return damage.WeaponType == constants.WeaponTypeGrenade
}

func (damage *Damage) isValidPlayerDamageEvent(player *Player) bool {
	if damage.AttackerSteamID64 != player.SteamID64 {
		return false
	}
	if damage.AttackerSteamID64 == damage.VictimSteamID64 {
		return false
	}
	if damage.VictimSteamID64 == 0 {
		return false
	}
	if damage.IsAttackerControllingBot {
		return false
	}
	if damage.AttackerSide == damage.VictimSide {
		return false
	}

	return true
}
