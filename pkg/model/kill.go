package model

import (
	"github.com/akiver/cs-demo-analyzer/pkg/api/constants"
)

type Kill struct {
	Frame                    int                  `json:"frame"`
	Tick                     int                  `json:"tick"`
	RoundNumber              int                  `json:"roundNumber"`
	WeaponType               constants.WeaponType `json:"weaponType"`
	WeaponName               constants.WeaponName `json:"weaponName"`
	KillerName               string               `json:"killerName"`
	KillerSteamID64          uint64               `json:"killerSteamId"`
	KillerSide               constants.Side       `json:"killerSide"`
	KillerTeamName           string               `json:"killerTeamName"`
	KillerX                  float64              `json:"killerX"`
	KillerY                  float64              `json:"killerY"`
	KillerZ                  float64              `json:"killerZ"`
	IsKillerAirborne         bool                 `json:"is_killer_airborne"`
	IsKillerBlinded          bool                 `json:"is_killer_blinded"`
	IsKillerControllingBot   bool                 `json:"isKillerControllingBot"`
	VictimName               string               `json:"victimName"`
	VictimSteamID64          uint64               `json:"victimSteamId"`
	VictimSide               constants.Side       `json:"victimSide"`
	VictimTeamName           string               `json:"victimTeamName"`
	VictimX                  float64              `json:"victimX"`
	VictimY                  float64              `json:"victimY"`
	VictimZ                  float64              `json:"victimZ"`
	IsVictimAirborne         bool                 `json:"is_victim_airborne"`
	IsVictimBlinded          bool                 `json:"is_victim_blinded"`
	IsVictimControllingBot   bool                 `json:"isVictimControllingBot"`
	IsVictimInspectingWeapon bool                 `json:"isVictimInspectingWeapon"`
	AssisterName             string               `json:"assisterName"`
	AssisterSteamID64        uint64               `json:"assisterSteamId"`
	AssisterSide             constants.Side       `json:"assisterSide"`
	AssisterTeamName         string               `json:"assisterTeamName"`
	AssisterX                float64              `json:"assisterX"`
	AssisterY                float64              `json:"assisterY"`
	AssisterZ                float64              `json:"assisterZ"`
	IsAssisterControllingBot bool                 `json:"isAssisterControllingBot"`
	IsHeadshot               bool                 `json:"isHeadshot"`
	PenetratedObjects        int                  `json:"penetratedObjects"`
	IsAssistedFlash          bool                 `json:"isAssistedFlash"`
	IsThroughSmoke           bool                 `json:"isThroughSmoke"`
	IsNoScope                bool                 `json:"isNoScope"`
	IsTradeKill              bool                 `json:"isTradeKill"`  // The attacker did a trade kill
	IsTradeDeath             bool                 `json:"isTradeDeath"` // The victim did a trade death
	Distance                 float32              `json:"distance"`
}

func (kill *Kill) IsSuicide() bool {
	return kill.KillerSteamID64 == kill.VictimSteamID64
}

func (kill *Kill) IsTeamKill() bool {
	return kill.KillerSide == kill.VictimSide
}
