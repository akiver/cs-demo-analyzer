package model

import (
	"github.com/akiver/cs-demo-analyzer/pkg/api/constants"
)

type PlayerPosition struct {
	Frame                  int                    `json:"frame"`
	Tick                   int                    `json:"tick"`
	RoundNumber            int                    `json:"roundNumber"`
	IsAlive                bool                   `json:"isAlive"`
	Name                   string                 `json:"name"`
	SteamID64              uint64                 `json:"steamId"`
	X                      float64                `json:"x"`
	Y                      float64                `json:"y"`
	Z                      float64                `json:"z"`
	Yaw                    float32                `json:"yaw"`
	Pitch                  float32                `json:"pitch"`
	FlashDurationRemaining float64                `json:"flashDurationRemaining"`
	Side                   constants.Side         `json:"side"`
	Money                  int                    `json:"money"`
	Health                 int                    `json:"health"`
	Armor                  int                    `json:"armor"`
	HasHelmet              bool                   `json:"hasHelmet"`
	HasBomb                bool                   `json:"hasBomb"`
	HasDefuseKit           bool                   `json:"hasDefuseKit"`
	IsDucking              bool                   `json:"isDucking"`
	IsAirborne             bool                   `json:"isAirborne"`
	IsScoping              bool                   `json:"isScoping"`
	IsDefusing             bool                   `json:"isDefusing"`
	IsPlanting             bool                   `json:"isPlanting"`
	IsGrabbingHostage      bool                   `json:"isGrabbingHostage"`
	ActiveWeaponName       constants.WeaponName   `json:"activeWeaponName"`
	Equipments             []constants.WeaponName `json:"equipments"`
	Grenades               []constants.WeaponName `json:"grenades"`
	Pistols                []constants.WeaponName `json:"pistols"`
	SMGs                   []constants.WeaponName `json:"smgs"`
	Rifles                 []constants.WeaponName `json:"rifles"`
	Heavy                  []constants.WeaponName `json:"heavy"`
	Ping                   int                    `json:"ping"`
}
