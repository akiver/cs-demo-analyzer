package api

import (
	"github.com/akiver/cs-demo-analyzer/internal/slice"
	"github.com/akiver/cs-demo-analyzer/pkg/api/constants"
	common "github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/common"
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
	FlashDurationRemaining float64                `json:"flashDurationRemaining"`
	Side                   common.Team            `json:"side"`
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
}

func newPlayerPosition(analyzer *Analyzer, player *common.Player) *PlayerPosition {
	parser := analyzer.parser
	hasBomb := false
	var equipments []constants.WeaponName
	var grenades []constants.WeaponName
	var pistols []constants.WeaponName
	var smgs []constants.WeaponName
	var rifles []constants.WeaponName
	var heavy []constants.WeaponName
	for _, weapon := range player.Weapons() {
		weaponName := equipmentToWeaponName[weapon.Type]
		// Weird bug encountered with a demo from 2017.
		// Player's weapons may contains duplicates, add it only if it has not been added yet.
		switch weapon.Class() {
		case common.EqClassEquipment:
			equipments = slice.AppendIfNotInSlice(equipments, weaponName)
		case common.EqClassGrenade:
			if weapon.Type == common.EqFlash {
				flashbangCount := player.FlashbangCount()
				for i := uint64(0); i < flashbangCount; i++ {
					grenades = append(grenades, constants.WeaponFlashbang)
				}
			} else {
				grenades = slice.AppendIfNotInSlice(grenades, weaponName)
			}
		case common.EqClassPistols:
			pistols = slice.AppendIfNotInSlice(pistols, weaponName)
		case common.EqClassSMG:
			smgs = slice.AppendIfNotInSlice(smgs, weaponName)
		case common.EqClassRifle:
			rifles = slice.AppendIfNotInSlice(rifles, weaponName)
		case common.EqClassHeavy:
			heavy = slice.AppendIfNotInSlice(heavy, weaponName)
		}

		if weapon.Type == common.EqBomb {
			hasBomb = true
		}
	}

	activeWeapon := constants.WeaponUnknown
	if player.ActiveWeapon() != nil {
		activeWeapon = equipmentToWeaponName[player.ActiveWeapon().Type]
	}

	return &PlayerPosition{
		Frame:                  parser.CurrentFrame(),
		Tick:                   analyzer.currentTick(),
		RoundNumber:            analyzer.currentRound.Number,
		Name:                   player.Name,
		SteamID64:              player.SteamID64,
		IsAlive:                player.IsAlive(),
		X:                      player.Position().X,
		Y:                      player.Position().Y,
		Z:                      player.Position().Z,
		Yaw:                    player.ViewDirectionX(),
		FlashDurationRemaining: player.FlashDurationTimeRemaining().Seconds(),
		Side:                   player.Team,
		Money:                  player.Money(),
		Health:                 player.Health(),
		Armor:                  player.Armor(),
		HasHelmet:              player.HasHelmet(),
		HasBomb:                hasBomb,
		HasDefuseKit:           player.HasDefuseKit(),
		IsDucking:              player.IsDucking(),
		IsAirborne:             player.IsAirborne(),
		IsScoping:              player.IsScoped(),
		IsDefusing:             player.IsDefusing,
		IsPlanting:             player.IsPlanting,
		IsGrabbingHostage:      player.IsGrabbingHostage(),
		ActiveWeaponName:       activeWeapon,
		Equipments:             equipments,
		Grenades:               grenades,
		Pistols:                pistols,
		SMGs:                   smgs,
		Rifles:                 rifles,
		Heavy:                  heavy,
	}
}
