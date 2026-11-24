package api

import (
	"testing"

	"github.com/akiver/cs-demo-analyzer/internal/visibility"
	"github.com/akiver/cs-demo-analyzer/pkg/api/constants"
	"github.com/golang/geo/r3"
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/common"
)

const testTickRate = 64

// addBox appends the 12 triangles of an axis aligned box.
func addBox(vertices []visibility.Vec3, indices []uint32, lower visibility.Vec3, upper visibility.Vec3) ([]visibility.Vec3, []uint32) {
	base := uint32(len(vertices))
	for i := range 8 {
		vertex := lower
		if i&1 != 0 {
			vertex[0] = upper[0]
		}
		if i&2 != 0 {
			vertex[1] = upper[1]
		}
		if i&4 != 0 {
			vertex[2] = upper[2]
		}
		vertices = append(vertices, vertex)
	}
	for _, face := range [][4]uint32{{0, 1, 3, 2}, {4, 5, 7, 6}, {0, 1, 5, 4}, {2, 3, 7, 6}, {0, 2, 6, 4}, {1, 3, 7, 5}} {
		indices = append(indices,
			base+face[0], base+face[1], base+face[2],
			base+face[0], base+face[2], base+face[3],
		)
	}

	return vertices, indices
}

// The attacker stands at (-300, 0) looking at +X, a wall at x = 0 spans y from -200 to 200.
// The victim stands at x = 300, it's hidden behind the wall at y = 0 and visible at y = 600.
var (
	attackerOrigin = r3.Vector{X: -300, Y: 0, Z: 0}
	hiddenOrigin   = r3.Vector{X: 300, Y: 0, Z: 0}
	visibleOrigin  = r3.Vector{X: 300, Y: 600, Z: 0}
)

type ttdTest struct {
	t        *testing.T
	tracker  *visibilityTracker
	attacker *common.Player
	victim   *common.Player
}

func newTTDTest(t *testing.T) *ttdTest {
	var vertices []visibility.Vec3
	var indices []uint32
	vertices, indices = addBox(vertices, indices, visibility.Vec3{-2000, -2000, -10}, visibility.Vec3{2000, 2000, 0})
	vertices, indices = addBox(vertices, indices, visibility.Vec3{-10, -200, 0}, visibility.Vec3{10, 200, 128})

	return &ttdTest{
		t:        t,
		tracker:  newVisibilityTracker(visibility.BuildBVH(vertices, indices), testTickRate),
		attacker: &common.Player{Name: "attacker", Team: common.TeamTerrorists},
		victim:   &common.Player{Name: "victim", Team: common.TeamCounterTerrorists},
	}
}

func snapshotAt(origin r3.Vector, yaw float64) playerSnapshot {
	return playerSnapshot{
		recorded:         true,
		isAlive:          true,
		origin:           origin,
		eyes:             origin.Add(r3.Vector{Z: defaultEyesHeight}),
		yaw:              yaw,
		fovTanHorizontal: defaultFovTanHorizontal,
		fovTanVertical:   defaultFovTanVertical,
	}
}

// recordTicks records the ticks from first to last (included), the victim is visible when isVisible returns true.
func (test *ttdTest) recordTicks(first int, last int, isVisible func(tick int) bool, attackerSnapshot func(tick int) playerSnapshot) {
	for tick := first; tick <= last; tick++ {
		frame := test.tracker.startFrame(tick)
		if frame == nil {
			test.t.Fatalf("tick %d already recorded", tick)
		}

		attacker := snapshotAt(attackerOrigin, 0)
		if attackerSnapshot != nil {
			attacker = attackerSnapshot(tick)
		}
		test.tracker.setSnapshot(frame, test.attacker, attacker)

		victimOrigin := hiddenOrigin
		if isVisible(tick) {
			victimOrigin = visibleOrigin
		}
		test.tracker.setSnapshot(frame, test.victim, snapshotAt(victimOrigin, 180))
		test.tracker.resolvePendingEvents()
	}
}

func (test *ttdTest) damage(tick int) *Damage {
	damage := &Damage{Tick: tick, WeaponType: constants.WeaponTypeRifle, TimeToDamage: noTimeToDamage}
	test.tracker.onDamage(tick, test.attacker, test.victim, damage)

	return damage
}

func (test *ttdTest) kill(tick int) *Kill {
	kill := &Kill{Tick: tick, WeaponType: constants.WeaponTypeRifle, TimeToKill: noTimeToDamage}
	test.tracker.onKill(tick, test.attacker, test.victim, kill)

	return kill
}

func ticksToMs(ticks int) int {
	return int(float64(ticks)/testTickRate*1000 + 0.5)
}

func visibleFrom(firstVisibleTick int) func(int) bool {
	return func(tick int) bool { return tick >= firstVisibleTick }
}

func TestTimeToDamagePeek(t *testing.T) {
	test := newTTDTest(t)
	test.recordTicks(0, 129, visibleFrom(100), nil)
	firstDamage := test.damage(130)
	test.recordTicks(130, 139, visibleFrom(100), nil)
	secondDamage := test.damage(140)
	test.recordTicks(140, 149, visibleFrom(100), nil)
	killDamage := test.damage(150)
	kill := test.kill(150)
	test.recordTicks(150, 150, visibleFrom(100), nil)

	if firstDamage.TimeToDamage != ticksToMs(30) {
		t.Errorf("first damage TTD = %d, want %d", firstDamage.TimeToDamage, ticksToMs(30))
	}
	if secondDamage.TimeToDamage != noTimeToDamage {
		t.Errorf("second damage TTD = %d, want none", secondDamage.TimeToDamage)
	}
	if killDamage.TimeToDamage != noTimeToDamage {
		t.Errorf("kill damage TTD = %d, want none", killDamage.TimeToDamage)
	}
	if kill.TimeToKill != ticksToMs(50) {
		t.Errorf("TTK = %d, want %d", kill.TimeToKill, ticksToMs(50))
	}
}

func TestTimeToDamageOneTap(t *testing.T) {
	test := newTTDTest(t)
	test.recordTicks(0, 119, visibleFrom(100), nil)
	damage := test.damage(120)
	kill := test.kill(120)
	test.recordTicks(120, 120, visibleFrom(100), nil)

	if damage.TimeToDamage != ticksToMs(20) || kill.TimeToKill != ticksToMs(20) {
		t.Errorf("TTD = %d, TTK = %d, want %d for both", damage.TimeToDamage, kill.TimeToKill, ticksToMs(20))
	}
}

func TestTimeToDamageWallbang(t *testing.T) {
	test := newTTDTest(t)
	never := func(int) bool { return false }
	test.recordTicks(0, 129, never, nil)
	damage := test.damage(130)
	kill := test.kill(130)
	test.recordTicks(130, 130, never, nil)

	if damage.TimeToDamage != noTimeToDamage || kill.TimeToKill != noTimeToDamage {
		t.Errorf("TTD = %d, TTK = %d, want none", damage.TimeToDamage, kill.TimeToKill)
	}
}

func TestTimeToDamageGapTolerance(t *testing.T) {
	// Hidden during 1 tick, tolerated.
	test := newTTDTest(t)
	isVisible := func(tick int) bool { return tick >= 100 && tick != 115 }
	test.recordTicks(0, 129, isVisible, nil)
	damage := test.damage(130)
	test.recordTicks(130, 130, isVisible, nil)
	if damage.TimeToDamage != ticksToMs(30) {
		t.Errorf("TTD = %d, want %d", damage.TimeToDamage, ticksToMs(30))
	}

	// Hidden during 10 ticks, the enemy is spotted again after.
	test = newTTDTest(t)
	isVisible = func(tick int) bool { return tick >= 100 && (tick < 110 || tick >= 120) }
	test.recordTicks(0, 129, isVisible, nil)
	damage = test.damage(130)
	test.recordTicks(130, 130, isVisible, nil)
	if damage.TimeToDamage != ticksToMs(10) {
		t.Errorf("TTD = %d, want %d", damage.TimeToDamage, ticksToMs(10))
	}
}

func TestTimeToDamageRepeek(t *testing.T) {
	test := newTTDTest(t)
	isVisible := func(tick int) bool { return (tick >= 100 && tick <= 120) || tick >= 200 }
	test.recordTicks(0, 109, isVisible, nil)
	first := test.damage(110)
	test.recordTicks(110, 229, isVisible, nil)
	second := test.damage(230)
	kill := test.kill(240)
	test.recordTicks(230, 240, isVisible, nil)

	if first.TimeToDamage != ticksToMs(10) {
		t.Errorf("first TTD = %d, want %d", first.TimeToDamage, ticksToMs(10))
	}
	if second.TimeToDamage != ticksToMs(30) {
		t.Errorf("second TTD = %d, want %d", second.TimeToDamage, ticksToMs(30))
	}
	if kill.TimeToKill != ticksToMs(40) {
		t.Errorf("TTK = %d, want %d", kill.TimeToKill, ticksToMs(40))
	}
}

func TestTimeToDamageVisibleForTooLong(t *testing.T) {
	test := newTTDTest(t)
	always := func(int) bool { return true }
	test.recordTicks(0, 499, always, nil)
	damage := test.damage(500)
	kill := test.kill(510)
	test.recordTicks(500, 510, always, nil)

	if damage.TimeToDamage != noTimeToDamage || kill.TimeToKill != noTimeToDamage {
		t.Errorf("TTD = %d, TTK = %d, want none", damage.TimeToDamage, kill.TimeToKill)
	}
}

func TestTimeToDamageFieldOfView(t *testing.T) {
	test := newTTDTest(t)
	// The attacker looks away until tick 110 and turns to the victim.
	attacker := func(tick int) playerSnapshot {
		if tick < 110 {
			return snapshotAt(attackerOrigin, 180)
		}
		return snapshotAt(attackerOrigin, 0)
	}
	test.recordTicks(0, 129, visibleFrom(100), attacker)
	damage := test.damage(130)
	test.recordTicks(130, 130, visibleFrom(100), attacker)

	if damage.TimeToDamage != ticksToMs(20) {
		t.Errorf("TTD = %d, want %d", damage.TimeToDamage, ticksToMs(20))
	}
}

func TestTimeToDamageBlinded(t *testing.T) {
	test := newTTDTest(t)
	attacker := func(tick int) playerSnapshot {
		snapshot := snapshotAt(attackerOrigin, 0)
		snapshot.isBlinded = tick < 120
		return snapshot
	}
	test.recordTicks(0, 129, visibleFrom(100), attacker)
	damage := test.damage(130)
	test.recordTicks(130, 130, visibleFrom(100), attacker)

	if damage.TimeToDamage != ticksToMs(10) {
		t.Errorf("TTD = %d, want %d", damage.TimeToDamage, ticksToMs(10))
	}
}

func TestTimeToDamageSmoke(t *testing.T) {
	test := newTTDTest(t)
	// A smoke between the attacker and the visible position of the victim.
	test.tracker.addSmoke(1, r3.Vector{X: 0, Y: 300, Z: 0})
	test.recordTicks(0, 129, visibleFrom(100), nil)
	damage := test.damage(130)
	test.recordTicks(130, 130, visibleFrom(100), nil)
	if damage.TimeToDamage != noTimeToDamage {
		t.Errorf("TTD through smoke = %d, want none", damage.TimeToDamage)
	}

	// The smoke fades away at tick 140.
	test.recordTicks(131, 139, visibleFrom(100), nil)
	test.tracker.removeSmoke(1)
	test.recordTicks(140, 149, visibleFrom(100), nil)
	damage = test.damage(150)
	test.recordTicks(150, 150, visibleFrom(100), nil)
	if damage.TimeToDamage != ticksToMs(10) {
		t.Errorf("TTD after the smoke = %d, want %d", damage.TimeToDamage, ticksToMs(10))
	}
}

func TestTimeToDamageIgnoredEvents(t *testing.T) {
	test := newTTDTest(t)
	teammate := &common.Player{Name: "teammate", Team: common.TeamTerrorists}
	test.recordTicks(0, 129, visibleFrom(100), nil)

	grenadeDamage := &Damage{Tick: 130, WeaponType: constants.WeaponTypeGrenade, TimeToDamage: noTimeToDamage}
	test.tracker.onDamage(130, test.attacker, test.victim, grenadeDamage)
	teamDamage := &Damage{Tick: 130, WeaponType: constants.WeaponTypeRifle, TimeToDamage: noTimeToDamage}
	test.tracker.onDamage(130, test.attacker, teammate, teamDamage)
	worldDamage := &Damage{Tick: 130, WeaponType: constants.WeaponTypeWorld, TimeToDamage: noTimeToDamage}
	test.tracker.onDamage(130, nil, test.victim, worldDamage)

	if len(test.tracker.pending) != 0 {
		t.Fatalf("expected no pending event, got %d", len(test.tracker.pending))
	}
}

func TestTimeToDamageReset(t *testing.T) {
	test := newTTDTest(t)
	test.recordTicks(0, 129, visibleFrom(100), nil)
	test.damage(130)
	test.tracker.reset()
	if len(test.tracker.pending) != 0 || test.tracker.frameCount != 0 || len(test.tracker.engagements) != 0 {
		t.Fatal("reset must clear the pending events, frames and engagements")
	}

	// The victim is visible since the first frame of the new round, the spotted tick is unknown.
	test.recordTicks(1000, 1009, visibleFrom(1000), nil)
	damage := test.damage(1010)
	test.recordTicks(1010, 1010, visibleFrom(1000), nil)
	if damage.TimeToDamage != noTimeToDamage {
		t.Errorf("TTD = %d, want none", damage.TimeToDamage)
	}
}

func TestFieldOfView(t *testing.T) {
	snapshot := snapshotAt(r3.Vector{}, 90) // Looking at +Y
	forward, right, up := snapshot.viewBasis()
	eyes := snapshot.eyes

	tests := []struct {
		name  string
		point r3.Vector
		want  bool
	}{
		{"in front", eyes.Add(r3.Vector{Y: 100}), true},
		{"behind", eyes.Add(r3.Vector{Y: -100}), false},
		{"45° to the right", eyes.Add(r3.Vector{X: 100, Y: 100}), true},
		{"60° to the right", eyes.Add(r3.Vector{X: 173, Y: 100}), false},
		{"30° up", eyes.Add(r3.Vector{Y: 100, Z: 57}), true},
		{"45° down", eyes.Add(r3.Vector{Y: 100, Z: -100}), false},
	}
	for _, test := range tests {
		if got := snapshot.isInFieldOfView(test.point, forward, right, up); got != test.want {
			t.Errorf("%s: isInFieldOfView = %v, want %v", test.name, got, test.want)
		}
	}

	// Looking down 45° makes the point 45° down visible.
	snapshot.pitch = 45
	forward, right, up = snapshot.viewBasis()
	if !snapshot.isInFieldOfView(eyes.Add(r3.Vector{Y: 100, Z: -100}), forward, right, up) {
		t.Error("expected the point to be visible when looking down")
	}

	// Scoped, the field of view is narrower.
	snapshot.pitch = 0
	snapshot.fovTanHorizontal, snapshot.fovTanVertical = fovTangents(scopedFovDegrees)
	forward, right, up = snapshot.viewBasis()
	if snapshot.isInFieldOfView(eyes.Add(r3.Vector{X: 100, Y: 100}), forward, right, up) {
		t.Error("expected the point to be outside of the scoped field of view")
	}
}

func TestHitboxPointsFollowYawAndCrouch(t *testing.T) {
	standing := snapshotAt(r3.Vector{X: 100, Y: 100, Z: 0}, 90)
	head := standing.hitboxPoint(hitboxPoint{0, 0, 1})
	if head != standing.eyes {
		t.Errorf("head = %v, want the eyes %v", head, standing.eyes)
	}

	// Facing +Y, the right side is +X and the muzzle is in front.
	rightShoulder := standing.hitboxPoint(hitboxPoint{0, 10, 0.82})
	if rightShoulder.X <= standing.origin.X {
		t.Errorf("right shoulder %v should be at +X when facing +Y", rightShoulder)
	}
	muzzle := standing.hitboxPoint(hitboxPoint{26, 4, 0.82})
	if muzzle.Y <= standing.origin.Y+20 {
		t.Errorf("muzzle %v should be in front when facing +Y", muzzle)
	}

	crouching := standing
	crouching.eyes.Z = crouching.origin.Z + 46
	if crouching.hitboxPoint(hitboxPoint{0, 0, 1}).Z != 46 {
		t.Error("the head must follow the eyes height when crouching")
	}
}

func TestSegmentIntersectsSmoke(t *testing.T) {
	smoke := r3.Vector{X: 0, Y: 0, Z: 0}
	tests := []struct {
		name string
		a    r3.Vector
		b    r3.Vector
		want bool
	}{
		{"through the center", r3.Vector{X: -500, Z: 64}, r3.Vector{X: 500, Z: 64}, true},
		{"beside", r3.Vector{X: -500, Y: 200, Z: 64}, r3.Vector{X: 500, Y: 200, Z: 64}, false},
		{"above", r3.Vector{X: -500, Z: 250}, r3.Vector{X: 500, Z: 250}, false},
		{"ends before", r3.Vector{X: -500, Z: 64}, r3.Vector{X: -200, Z: 64}, false},
		{"inside", r3.Vector{X: -10, Z: 64}, r3.Vector{X: 10, Z: 64}, true},
	}
	for _, test := range tests {
		if got := segmentIntersectsSmoke(test.a, test.b, smoke); got != test.want {
			t.Errorf("%s: segmentIntersectsSmoke = %v, want %v", test.name, got, test.want)
		}
	}
}

func TestTimeToDamageStat(t *testing.T) {
	stat := func([]int) float32 { return 42 }
	if got := timeToDamageStat(nil, stat); got != noTimeToDamage {
		t.Errorf("stat of no times = %v, want %d", got, noTimeToDamage)
	}
	if got := timeToDamageStat([]int{100}, stat); got != 42 {
		t.Errorf("stat = %v, want 42", got)
	}
}
