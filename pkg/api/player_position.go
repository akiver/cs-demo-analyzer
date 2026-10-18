package api

import (
	"slices"

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
	Pitch                  float32                `json:"pitch"`
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
	Ping                   int                    `json:"ping"`
}

// playerWeapons are the weapon names of a player's inventory by category, as exported in PlayerPosition.
type playerWeapons struct {
	equipments []constants.WeaponName
	grenades   []constants.WeaponName
	pistols    []constants.WeaponName
	smgs       []constants.WeaponName
	rifles     []constants.WeaponName
	heavy      []constants.WeaponName
}

// reset empties the lists, keeping their backing arrays.
func (w *playerWeapons) reset() {
	w.equipments = w.equipments[:0]
	w.grenades = w.grenades[:0]
	w.pistols = w.pistols[:0]
	w.smgs = w.smgs[:0]
	w.rifles = w.rifles[:0]
	w.heavy = w.heavy[:0]
}

// internWeapons returns previous if it holds the same weapons as current (in any order), otherwise a sorted copy of
// current (nil when empty). Player positions are created for every player on every frame while inventories change
// a few times per round, so sharing the lists between consecutive positions avoids allocating six slices per
// position. Positions only read them.
// The copies are sorted because current is built by ranging the inventory map, whose order differs on every run;
// comparing before sorting keeps the sort out of the common case where the inventory didn't change.
func internWeapons(previous []constants.WeaponName, current []constants.WeaponName) []constants.WeaponName {
	if len(current) == 0 {
		return nil
	}

	if sameWeapons(previous, current) {
		return previous
	}

	weapons := slices.Clone(current)
	slices.Sort(weapons)

	return weapons
}

// sameWeapons returns whether a and b hold the same weapon names with the same multiplicities, in any order.
func sameWeapons(a []constants.WeaponName, b []constants.WeaponName) bool {
	if len(a) != len(b) {
		return false
	}

	for _, weapon := range a {
		if count(a, weapon) != count(b, weapon) {
			return false
		}
	}

	return true
}

func count(weapons []constants.WeaponName, weapon constants.WeaponName) int {
	n := 0
	for _, w := range weapons {
		if w == weapon {
			n++
		}
	}

	return n
}

func newPlayerPosition(analyzer *Analyzer, player *common.Player) *PlayerPosition {
	parser := analyzer.parser
	hasBomb := false
	current := &analyzer.playerWeaponsScratch
	current.reset()
	// Player.Inventory is Player.Weapons() without copying the map to a slice first, this runs for every player on
	// every frame.
	for _, weapon := range player.Inventory {
		weaponName := equipmentToWeaponName[weapon.Type]
		// Weird bug encountered with a demo from 2017.
		// Player's weapons may contains duplicates, add it only if it has not been added yet.
		switch weapon.Class() {
		case common.EqClassEquipment:
			current.equipments = slice.AppendIfNotInSlice(current.equipments, weaponName)
		case common.EqClassGrenade:
			if weapon.Type == common.EqFlash {
				flashbangCount := player.FlashbangCount()
				for i := uint64(0); i < flashbangCount; i++ {
					current.grenades = append(current.grenades, constants.WeaponFlashbang)
				}
			} else {
				current.grenades = slice.AppendIfNotInSlice(current.grenades, weaponName)
			}
		case common.EqClassPistols:
			current.pistols = slice.AppendIfNotInSlice(current.pistols, weaponName)
		case common.EqClassSMG:
			current.smgs = slice.AppendIfNotInSlice(current.smgs, weaponName)
		case common.EqClassRifle:
			current.rifles = slice.AppendIfNotInSlice(current.rifles, weaponName)
		case common.EqClassHeavy:
			current.heavy = slice.AppendIfNotInSlice(current.heavy, weaponName)
		}

		if weapon.Type == common.EqBomb {
			hasBomb = true
		}
	}

	weapons, exists := analyzer.playerWeapons[player]
	if !exists {
		weapons = &playerWeapons{}
		analyzer.playerWeapons[player] = weapons
	}

	weapons.equipments = internWeapons(weapons.equipments, current.equipments)
	weapons.grenades = internWeapons(weapons.grenades, current.grenades)
	weapons.pistols = internWeapons(weapons.pistols, current.pistols)
	weapons.smgs = internWeapons(weapons.smgs, current.smgs)
	weapons.rifles = internWeapons(weapons.rifles, current.rifles)
	weapons.heavy = internWeapons(weapons.heavy, current.heavy)

	activeWeapon := constants.WeaponUnknown
	if player.ActiveWeapon() != nil {
		activeWeapon = equipmentToWeaponName[player.ActiveWeapon().Type]
	}

	position := player.Position()

	return &PlayerPosition{
		Frame:                  parser.CurrentFrame(),
		Tick:                   analyzer.currentTick(),
		RoundNumber:            analyzer.currentRound.Number,
		Name:                   player.Name,
		SteamID64:              player.SteamID64,
		IsAlive:                player.IsAlive(),
		X:                      position.X,
		Y:                      position.Y,
		Z:                      position.Z,
		Yaw:                    player.ViewDirectionX(),
		Pitch:                  player.ViewDirectionY(),
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
		Equipments:             weapons.equipments,
		Grenades:               weapons.grenades,
		Pistols:                weapons.pistols,
		SMGs:                   weapons.smgs,
		Rifles:                 weapons.rifles,
		Heavy:                  weapons.heavy,
		Ping:                   player.Ping(),
	}
}
