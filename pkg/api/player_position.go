package api

import (
	"slices"

	"github.com/akiver/cs-demo-analyzer/internal/slice"
	"github.com/akiver/cs-demo-analyzer/pkg/api/constants"
	common "github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/common"
)

type PlayerPosition struct {
	Tick                   int                  `json:"tick"`
	RoundNumber            int                  `json:"roundNumber"`
	IsAlive                bool                 `json:"isAlive"`
	Name                   string               `json:"name"`
	SteamID64              uint64               `json:"steamId"`
	X                      float64              `json:"x"`
	Y                      float64              `json:"y"`
	Z                      float64              `json:"z"`
	Yaw                    float32              `json:"yaw"`
	Pitch                  float32              `json:"pitch"`
	FlashDurationRemaining float64              `json:"flashDurationRemaining"`
	Side                   common.Team          `json:"side"`
	Money                  int                  `json:"money"`
	Health                 int                  `json:"health"`
	Armor                  int                  `json:"armor"`
	HasHelmet              bool                 `json:"hasHelmet"`
	HasBomb                bool                 `json:"hasBomb"`
	HasDefuseKit           bool                 `json:"hasDefuseKit"`
	IsDucking              bool                 `json:"isDucking"`
	IsAirborne             bool                 `json:"isAirborne"`
	IsScoping              bool                 `json:"isScoping"`
	IsDefusing             bool                 `json:"isDefusing"`
	IsPlanting             bool                 `json:"isPlanting"`
	IsGrabbingHostage      bool                 `json:"isGrabbingHostage"`
	ActiveWeaponName       constants.WeaponName `json:"activeWeaponName"`
	// Embedded so the JSON fields stay flat. Never nil, its fields wouldn't be in the JSON otherwise.
	*PlayerWeapons
	Ping int `json:"ping"`
}

// PlayerWeapons are the weapon names of a player's inventory by category.
// The same PlayerWeapons is shared by consecutive positions of a player while the inventory doesn't change, it must
// not be mutated.
type PlayerWeapons struct {
	Equipments []constants.WeaponName `json:"equipments"`
	Grenades   []constants.WeaponName `json:"grenades"`
	Pistols    []constants.WeaponName `json:"pistols"`
	SMGs       []constants.WeaponName `json:"smgs"`
	Rifles     []constants.WeaponName `json:"rifles"`
	Heavy      []constants.WeaponName `json:"heavy"`
}

// reset empties the lists, keeping their backing arrays.
func (w *PlayerWeapons) reset() {
	w.Equipments = w.Equipments[:0]
	w.Grenades = w.Grenades[:0]
	w.Pistols = w.Pistols[:0]
	w.SMGs = w.SMGs[:0]
	w.Rifles = w.Rifles[:0]
	w.Heavy = w.Heavy[:0]
}

// internPlayerWeapons returns previous if it holds the same weapons as current (in any order), otherwise a new
// PlayerWeapons with sorted copies of current's lists. Player positions are created for every player on every frame
// while inventories change a few times per round, so sharing the weapons between consecutive positions keeps
// the six lists out of each position.
// The copies are sorted because current is built by ranging the inventory map, whose order differs on every run;
// comparing before sorting keeps the sort out of the common case where the inventory didn't change.
func internPlayerWeapons(previous *PlayerWeapons, current *PlayerWeapons) *PlayerWeapons {
	if previous != nil &&
		sameWeapons(previous.Equipments, current.Equipments) &&
		sameWeapons(previous.Grenades, current.Grenades) &&
		sameWeapons(previous.Pistols, current.Pistols) &&
		sameWeapons(previous.SMGs, current.SMGs) &&
		sameWeapons(previous.Rifles, current.Rifles) &&
		sameWeapons(previous.Heavy, current.Heavy) {
		return previous
	}

	return &PlayerWeapons{
		Equipments: sortedWeapons(current.Equipments),
		Grenades:   sortedWeapons(current.Grenades),
		Pistols:    sortedWeapons(current.Pistols),
		SMGs:       sortedWeapons(current.SMGs),
		Rifles:     sortedWeapons(current.Rifles),
		Heavy:      sortedWeapons(current.Heavy),
	}
}

// sortedWeapons returns a sorted copy of weapons, nil when empty.
func sortedWeapons(weapons []constants.WeaponName) []constants.WeaponName {
	if len(weapons) == 0 {
		return nil
	}

	sorted := slices.Clone(weapons)
	slices.Sort(sorted)

	return sorted
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
			current.Equipments = slice.AppendIfNotInSlice(current.Equipments, weaponName)
		case common.EqClassGrenade:
			if weapon.Type == common.EqFlash {
				flashbangCount := player.FlashbangCount()
				for i := uint64(0); i < flashbangCount; i++ {
					current.Grenades = append(current.Grenades, constants.WeaponFlashbang)
				}
			} else {
				current.Grenades = slice.AppendIfNotInSlice(current.Grenades, weaponName)
			}
		case common.EqClassPistols:
			current.Pistols = slice.AppendIfNotInSlice(current.Pistols, weaponName)
		case common.EqClassSMG:
			current.SMGs = slice.AppendIfNotInSlice(current.SMGs, weaponName)
		case common.EqClassRifle:
			current.Rifles = slice.AppendIfNotInSlice(current.Rifles, weaponName)
		case common.EqClassHeavy:
			current.Heavy = slice.AppendIfNotInSlice(current.Heavy, weaponName)
		}

		if weapon.Type == common.EqBomb {
			hasBomb = true
		}
	}

	weapons := internPlayerWeapons(analyzer.playerWeapons[player], current)
	analyzer.playerWeapons[player] = weapons

	activeWeapon := constants.WeaponUnknown
	if player.ActiveWeapon() != nil {
		activeWeapon = equipmentToWeaponName[player.ActiveWeapon().Type]
	}

	position := player.Position()

	return &PlayerPosition{
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
		PlayerWeapons:          weapons,
		Ping:                   player.Ping(),
	}
}
