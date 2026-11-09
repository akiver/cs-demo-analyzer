package api

import (
	"reflect"
	"testing"

	"github.com/akiver/cs-demo-analyzer/pkg/api/constants"
	"github.com/golang/geo/r2"
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/common"
)

func newSpooledMatch(t *testing.T) *Match {
	t.Helper()

	spool, err := newPositionsSpool()
	if err != nil {
		t.Fatalf("failed to create the spool: %v", err)
	}
	match := &Match{positionsSpool: spool}
	t.Cleanup(func() {
		if err := match.closePositionsSpool(); err != nil {
			t.Errorf("failed to close the spool: %v", err)
		}
	})

	return match
}

// spillRound writes the positions in memory as the round and empties them, like Analyzer.spillPositions().
func spillRound(match *Match, roundNumber int) {
	match.positionsSpool.writeRound(roundNumber, match)
	match.PlayerPositions = nil
	match.GrenadePositions = nil
	match.InfernoPositions = nil
	match.ChickenPositions = nil
	match.HostagePositions = nil
}

func newTestPlayerPosition(tick int, roundNumber int, weapons *PlayerWeapons) *PlayerPosition {
	return &PlayerPosition{
		Tick:                   tick,
		RoundNumber:            roundNumber,
		IsAlive:                true,
		Name:                   "Player \"quoted\" ünïcode",
		SteamID64:              76561198000000000 + uint64(tick),
		X:                      -1234.5678,
		Y:                      0.1,
		Z:                      64,
		Yaw:                    -179.5,
		Pitch:                  12.25,
		FlashDurationRemaining: 1.5,
		Side:                   common.TeamCounterTerrorists,
		Money:                  16000,
		Health:                 -1,
		Armor:                  100,
		HasHelmet:              true,
		HasDefuseKit:           true,
		IsDucking:              true,
		IsScoping:              tick%2 == 0,
		ActiveWeaponName:       constants.WeaponFlashbang,
		PlayerWeapons:          weapons,
		Ping:                   42,
	}
}

func readAllPositions(t *testing.T, match *Match) (players [][]*PlayerPosition, grenades [][]*GrenadePosition, infernos [][]*InfernoPosition, chickens [][]*ChickenPosition, hostages [][]*HostagePosition) {
	t.Helper()

	check := func(err error) {
		if err != nil {
			t.Fatalf("failed to read positions: %v", err)
		}
	}
	check(match.eachPlayerPositions(func(positions []*PlayerPosition) error {
		players = append(players, positions)
		return nil
	}))
	check(match.eachGrenadePositions(func(positions []*GrenadePosition) error {
		grenades = append(grenades, positions)
		return nil
	}))
	check(match.eachInfernoPositions(func(positions []*InfernoPosition) error {
		infernos = append(infernos, positions)
		return nil
	}))
	check(match.eachChickenPositions(func(positions []*ChickenPosition) error {
		chickens = append(chickens, positions)
		return nil
	}))
	check(match.eachHostagePositions(func(positions []*HostagePosition) error {
		hostages = append(hostages, positions)
		return nil
	}))

	return
}

func TestPositionsSpoolRoundTrip(t *testing.T) {
	match := newSpooledMatch(t)

	sharedWeapons := &PlayerWeapons{
		Equipments: []constants.WeaponName{constants.WeaponName("kevlar"), constants.WeaponName("defuse_kit")},
		Grenades:   []constants.WeaponName{constants.WeaponFlashbang, constants.WeaponFlashbang},
		Pistols:    []constants.WeaponName{constants.WeaponName("usp")},
		Rifles:     []constants.WeaponName{constants.WeaponName("ak47")},
		// SMGs and Heavy are nil, which marshals as null.
	}
	emptyWeapons := &PlayerWeapons{
		// Empty but not nil, which marshals as [].
		Equipments: []constants.WeaponName{},
	}
	round1Players := []*PlayerPosition{
		newTestPlayerPosition(10, 1, sharedWeapons),
		newTestPlayerPosition(11, 1, sharedWeapons),
		newTestPlayerPosition(12, 1, emptyWeapons),
	}
	round1Grenades := []*GrenadePosition{{
		Tick:             10,
		RoundNumber:      1,
		GrenadeID:        "123456789",
		ProjectileID:     "987654321",
		X:                1,
		Y:                2.5,
		Z:                -3,
		ThrowerSteamID64: 76561198000000001,
		ThrowerName:      "Thrower",
		ThrowerSide:      common.TeamTerrorists,
		ThrowerTeamName:  "Team A",
		ThrowerVelocityX: 250.5,
		ThrowerVelocityY: -0.25,
		ThrowerVelocityZ: 0,
		ThrowerPitch:     -5,
		ThrowerYaw:       90,
		GrenadeName:      constants.WeaponFlashbang,
	}}
	round1Infernos := []*InfernoPosition{
		{Tick: 10, RoundNumber: 1, ThrowerSteamID64: 1, ThrowerName: "T", UniqueID: "42", X: 1, Y: 2, Z: 3, ConvexHull2D: []r2.Point{{X: 1.5, Y: -2.5}, {X: 3, Y: 4}}},
		{Tick: 11, RoundNumber: 1, ThrowerSteamID64: 1, ThrowerName: "T", UniqueID: "42", X: 1, Y: 2, Z: 3, ConvexHull2D: []r2.Point{}},
		{Tick: 12, RoundNumber: 1, ThrowerSteamID64: 1, ThrowerName: "T", UniqueID: "42", X: 1, Y: 2, Z: 3, ConvexHull2D: nil},
	}
	round1Chickens := []*ChickenPosition{{chickenID: 3, Tick: 10, RoundNumber: 1, X: 1, Y: 2, Z: 3}}
	round1Hostages := []*HostagePosition{{entityID: 7, serialNum: 1234, Tick: 10, RoundNumber: 1, X: 1, Y: 2, Z: 3, State: common.HostageStateGettingPickedUp}}

	match.PlayerPositions = round1Players
	match.GrenadePositions = round1Grenades
	match.InfernoPositions = round1Infernos
	match.ChickenPositions = round1Chickens
	match.HostagePositions = round1Hostages
	spillRound(match, 1)

	// Round 2 is spooled then discarded, like a round restart.
	match.PlayerPositions = []*PlayerPosition{newTestPlayerPosition(20, 2, sharedWeapons)}
	spillRound(match, 2)
	match.resetRound(2)

	round2Players := []*PlayerPosition{newTestPlayerPosition(30, 2, emptyWeapons)}
	match.PlayerPositions = round2Players
	spillRound(match, 2)

	// Round 3 is still in memory, like the current round at the end of an analysis without a final spill.
	round3Players := []*PlayerPosition{newTestPlayerPosition(40, 3, sharedWeapons)}
	match.PlayerPositions = round3Players

	players, grenades, infernos, chickens, hostages := readAllPositions(t, match)

	if len(players) != 3 {
		t.Fatalf("expected 3 chunks of player positions (2 spooled rounds and the one in memory), got %d", len(players))
	}
	if !reflect.DeepEqual(players[0], round1Players) {
		t.Errorf("round 1 player positions differ after the round trip:\ngot  %+v\nwant %+v", players[0], round1Players)
	}
	if !reflect.DeepEqual(players[1], round2Players) {
		t.Errorf("round 2 player positions differ after the round trip:\ngot  %+v\nwant %+v", players[1], round2Players)
	}
	if !reflect.DeepEqual(players[2], round3Players) {
		t.Errorf("round 3 player positions in memory differ:\ngot  %+v\nwant %+v", players[2], round3Players)
	}
	if players[0][0].PlayerWeapons != players[0][1].PlayerWeapons {
		t.Errorf("positions sharing the same PlayerWeapons must share it after the round trip")
	}
	if players[0][0].PlayerWeapons == players[0][2].PlayerWeapons {
		t.Errorf("positions with different PlayerWeapons must not share it after the round trip")
	}

	if len(grenades) != 1 || !reflect.DeepEqual(grenades[0], round1Grenades) {
		t.Errorf("grenade positions differ after the round trip:\ngot  %+v\nwant %+v", grenades, round1Grenades)
	}
	if len(infernos) != 1 || !reflect.DeepEqual(infernos[0], round1Infernos) {
		t.Errorf("inferno positions differ after the round trip:\ngot  %+v\nwant %+v", infernos, round1Infernos)
	}
	if len(chickens) != 1 || !reflect.DeepEqual(chickens[0], round1Chickens) {
		t.Errorf("chicken positions differ after the round trip:\ngot  %+v\nwant %+v", chickens, round1Chickens)
	}
	if len(hostages) != 1 || !reflect.DeepEqual(hostages[0], round1Hostages) {
		t.Errorf("hostage positions differ after the round trip:\ngot  %+v\nwant %+v", hostages, round1Hostages)
	}
}

func TestPositionsSpoolReset(t *testing.T) {
	match := newSpooledMatch(t)
	weapons := &PlayerWeapons{}

	match.PlayerPositions = []*PlayerPosition{newTestPlayerPosition(10, 1, weapons)}
	spillRound(match, 1)
	match.PlayerPositions = []*PlayerPosition{newTestPlayerPosition(20, 2, weapons)}
	spillRound(match, 2)

	// A match restart discards everything spooled so far.
	match.reset()

	afterReset := []*PlayerPosition{newTestPlayerPosition(30, 1, weapons), newTestPlayerPosition(31, 1, weapons)}
	match.PlayerPositions = afterReset
	spillRound(match, 1)

	players, _, _, _, _ := readAllPositions(t, match)
	if len(players) != 1 {
		t.Fatalf("expected only the round spooled after the reset, got %d chunks", len(players))
	}
	if !reflect.DeepEqual(players[0], afterReset) {
		t.Errorf("player positions differ after the reset:\ngot  %+v\nwant %+v", players[0], afterReset)
	}
	if match.positionsSpool.size == 0 {
		t.Errorf("the spool must contain the round written after the reset")
	}
}

func TestPositionsSpoolWithoutSpool(t *testing.T) {
	match := &Match{}
	positions := []*PlayerPosition{newTestPlayerPosition(10, 1, &PlayerWeapons{})}
	match.PlayerPositions = positions

	players, _, _, _, _ := readAllPositions(t, match)
	if len(players) != 1 || !reflect.DeepEqual(players[0], positions) {
		t.Errorf("expected the positions in memory as a single chunk, got %+v", players)
	}
	if err := match.closePositionsSpool(); err != nil {
		t.Errorf("closing a match without a spool must not fail: %v", err)
	}
}

func TestPositionsSpoolRoundInSeveralSegments(t *testing.T) {
	match := newSpooledMatch(t)
	weapons := &PlayerWeapons{Rifles: []constants.WeaponName{constants.WeaponName("ak47")}}

	// A long round is spooled in parts before it's over, see maxPositionsInMemory.
	part1 := []*PlayerPosition{newTestPlayerPosition(10, 1, weapons), newTestPlayerPosition(11, 1, weapons)}
	match.PlayerPositions = part1
	spillRound(match, 1)
	part2 := []*PlayerPosition{newTestPlayerPosition(12, 1, weapons)}
	match.PlayerPositions = part2
	spillRound(match, 1)
	round2 := []*PlayerPosition{newTestPlayerPosition(20, 2, weapons)}
	match.PlayerPositions = round2
	spillRound(match, 2)

	players, _, _, _, _ := readAllPositions(t, match)
	if len(players) != 3 {
		t.Fatalf("expected 3 chunks (2 parts of round 1 and round 2), got %d", len(players))
	}
	if !reflect.DeepEqual(players[0], part1) || !reflect.DeepEqual(players[1], part2) || !reflect.DeepEqual(players[2], round2) {
		t.Errorf("positions differ after the round trip:\ngot  %+v\nwant %+v", players, [][]*PlayerPosition{part1, part2, round2})
	}

	// Discarding the round discards all its parts.
	match.resetRound(1)
	players, _, _, _, _ = readAllPositions(t, match)
	if len(players) != 1 || !reflect.DeepEqual(players[0], round2) {
		t.Errorf("expected only round 2 after discarding round 1, got %+v", players)
	}
}

func TestPositionsSpoolLargeSegment(t *testing.T) {
	match := newSpooledMatch(t)
	weapons := &PlayerWeapons{Rifles: []constants.WeaponName{constants.WeaponName("ak47")}}

	// Enough positions for the records to be written in several blocks, see spoolRecordsBlockSize.
	positions := make([]*PlayerPosition, 0, 20000)
	for tick := 0; tick < 20000; tick++ {
		positions = append(positions, newTestPlayerPosition(tick, 1, weapons))
	}
	match.PlayerPositions = positions
	spillRound(match, 1)

	players, _, _, _, _ := readAllPositions(t, match)
	if len(players) != 1 || !reflect.DeepEqual(players[0], positions) {
		t.Errorf("positions differ after the round trip of a large segment")
	}
	if match.positionsSpool.segments[0].recordsLength <= spoolRecordsBlockSize {
		t.Errorf("expected the records to span several blocks, got %d bytes", match.positionsSpool.segments[0].recordsLength)
	}
}
