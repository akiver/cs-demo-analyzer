package api

import (
	"errors"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"slices"

	"github.com/akiver/cs-demo-analyzer/internal/visibility"
	"github.com/akiver/cs-demo-analyzer/pkg/api/constants"
	"github.com/golang/geo/r3"
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/common"
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/events"
)

// Time to damage (TTD) is the time between the moment an attacker sees an enemy and the moment they damage them, time
// to kill (TTK) is the time until the kill.
//
// "Seeing" is computed by ray tracing the map collision geometry (see cmd/map-geometry) from the attacker's eyes to
// points of the enemy body, the game's spotted state (m_bSpottedByMask) is not used because it's radar based: delayed
// up to 500ms, only checking a few bones and still true a while after an enemy is no longer visible.
//
// To avoid tracing every pair of players every tick, the players' states of the last seconds are kept in a ring buffer
// and visibility is only computed when an attacker damages or kills an enemy, walking back in time until the enemy was
// not visible.

const (
	// How many seconds of players' states are kept. An enemy visible for longer before being damaged doesn't get a
	// time to damage, it's likely an enemy that didn't move while the attacker was holding an angle on something else.
	visibilityHistorySeconds = 5
	// An enemy that's not visible during less than this duration is still considered visible, it absorbs visibility
	// flickering when an enemy is on the edge of a wall or the screen, ~2 ticks at 64 tick.
	visibilityGapToleranceSeconds = 1.0 / 32
	// A player blinded by a flashbang for longer than this duration doesn't see anything, the end of the effect fades
	// out and lets the player see again.
	blindedMinRemainingSeconds = 1.0
	// Value used when a time to damage/kill is not available.
	noTimeToDamage = -1
)

// visibilityFrame contains the state of the players and smokes at a tick.
type visibilityFrame struct {
	tick      int
	snapshots []playerSnapshot // Indexed by player slot
	smokes    []r3.Vector
}

type playerPair struct {
	attacker *common.Player
	victim   *common.Player
}

// engagement is an attacker damaging an enemy that stayed visible since spottedTick.
type engagement struct {
	firstDamageTick int
	spottedTick     int // -1 when the enemy was visible for longer than the history duration
}

type pendingVisibilityEvent struct {
	tick     int
	attacker *common.Player
	victim   *common.Player
	damage   *Damage
	kill     *Kill
}

type visibilityStatus int

const (
	notVisible visibilityStatus = iota
	spotted
	reachedPreviousEngagement
	visibleForTooLong
)

type visibilityTracker struct {
	tracer     *visibility.Tracer
	tickRate   float64
	slots      map[*common.Player]int
	frames     []visibilityFrame // Ring buffer
	frameCount int
	nextFrame  int
	smokes     map[int]r3.Vector // Active smokes by grenade entity ID
	// Positions of the active smokes, a new slice is created when a smoke starts or expires so frames can share it.
	activeSmokes []r3.Vector
	pending      []pendingVisibilityEvent
	engagements  map[playerPair]engagement
}

// loadMapGeometry returns the BVH of the map from the map geometry folder and the checksum of its file, or nil if it's
// not available.
func loadMapGeometry(mapGeometryDir string, mapName string) (*visibility.BVH, string) {
	if mapGeometryDir == "" || mapName == "" {
		return nil, ""
	}

	path := filepath.Join(mapGeometryDir, mapName+visibility.GeometryFileExtension)
	bvh, checksum, err := visibility.LoadGeometryFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fmt.Fprintf(os.Stderr, "map geometry %s not found, time to damage will not be computed\n", path)
		} else {
			fmt.Fprintf(os.Stderr, "failed to load map geometry %s, time to damage will not be computed: %v\n", path, err)
		}
		return nil, ""
	}

	return bvh, checksum
}

func newVisibilityTracker(bvh *visibility.BVH, tickRate float64) *visibilityTracker {
	tracker := &visibilityTracker{
		tracer:      visibility.NewTracer(bvh),
		slots:       make(map[*common.Player]int),
		smokes:      make(map[int]r3.Vector),
		engagements: make(map[playerPair]engagement),
	}
	tracker.setTickRate(tickRate)

	return tracker
}

func (tracker *visibilityTracker) setTickRate(tickRate float64) {
	if tickRate <= 0 {
		tickRate = 64
	}
	tracker.tickRate = tickRate

	frameCount := int(math.Ceil(tickRate*visibilityHistorySeconds)) + 1
	if frameCount != len(tracker.frames) {
		tracker.frames = make([]visibilityFrame, frameCount)
		tracker.frameCount = 0
		tracker.nextFrame = 0
	}
}

// reset forgets everything about the current round, players are teleported at the start of a round.
func (tracker *visibilityTracker) reset() {
	tracker.frameCount = 0
	tracker.nextFrame = 0
	tracker.pending = tracker.pending[:0]
	clear(tracker.smokes)
	tracker.activeSmokes = nil
	clear(tracker.engagements)
}

func (tracker *visibilityTracker) addSmoke(entityID int, position r3.Vector) {
	tracker.smokes[entityID] = position
	tracker.updateActiveSmokes()
}

func (tracker *visibilityTracker) removeSmoke(entityID int) {
	if _, ok := tracker.smokes[entityID]; ok {
		delete(tracker.smokes, entityID)
		tracker.updateActiveSmokes()
	}
}

func (tracker *visibilityTracker) updateActiveSmokes() {
	smokes := make([]r3.Vector, 0, len(tracker.smokes))
	for _, position := range tracker.smokes {
		smokes = append(smokes, position)
	}
	tracker.activeSmokes = smokes
}

func (tracker *visibilityTracker) slot(player *common.Player) int {
	slot, ok := tracker.slots[player]
	if !ok {
		slot = len(tracker.slots)
		tracker.slots[player] = slot
	}

	return slot
}

// lastFrame returns the most recent frame or nil.
func (tracker *visibilityTracker) lastFrame() *visibilityFrame {
	if tracker.frameCount == 0 {
		return nil
	}

	return &tracker.frames[(tracker.nextFrame-1+len(tracker.frames))%len(tracker.frames)]
}

// frameAt returns the frame recorded age frames ago, 0 being the most recent one.
func (tracker *visibilityTracker) frameAt(age int) *visibilityFrame {
	return &tracker.frames[(tracker.nextFrame-1-age+2*len(tracker.frames))%len(tracker.frames)]
}

// recordFrame stores the state of the given players at the given tick.
// Only the first frame of a tick is recorded since several frames may share the same tick.
func (tracker *visibilityTracker) recordFrame(tick int, players []*common.Player) {
	previous := tracker.lastFrame()
	frame := tracker.startFrame(tick)
	if frame == nil {
		return
	}

	for _, player := range players {
		if player == nil || player.Team == common.TeamSpectators || player.Team == common.TeamUnassigned {
			continue
		}

		if player.IsAlive() {
			tracker.setSnapshot(frame, player, newPlayerSnapshot(player))
			continue
		}

		// Keep the last state of dead players, it's used when they are killed and their state is no longer valid.
		snapshot := playerSnapshot{}
		if slot := tracker.slot(player); previous != nil && slot < len(previous.snapshots) && previous.snapshots[slot].recorded {
			snapshot = previous.snapshots[slot]
			snapshot.isAlive = false
		}
		tracker.setSnapshot(frame, player, snapshot)
	}
}

// startFrame adds a frame for the given tick to the ring buffer and returns it, it returns nil if the tick has already
// been recorded.
func (tracker *visibilityTracker) startFrame(tick int) *visibilityFrame {
	if last := tracker.lastFrame(); last != nil && last.tick >= tick {
		return nil
	}

	frame := &tracker.frames[tracker.nextFrame]
	tracker.nextFrame = (tracker.nextFrame + 1) % len(tracker.frames)
	tracker.frameCount = min(tracker.frameCount+1, len(tracker.frames))

	frame.tick = tick
	frame.smokes = tracker.activeSmokes
	for i := range frame.snapshots {
		frame.snapshots[i].recorded = false
	}

	return frame
}

func (tracker *visibilityTracker) setSnapshot(frame *visibilityFrame, player *common.Player, snapshot playerSnapshot) {
	slot := tracker.slot(player)
	if slot >= len(frame.snapshots) {
		frame.snapshots = slices.Grow(frame.snapshots, slot+1-len(frame.snapshots))[:slot+1]
	}
	frame.snapshots[slot] = snapshot
}

func newPlayerSnapshot(player *common.Player) playerSnapshot {
	pawn := player.PlayerPawnEntity()
	if pawn == nil {
		return playerSnapshot{}
	}

	weapon := player.ActiveWeapon()
	isScoped := player.IsScoped()
	fovTanHorizontal, fovTanVertical := defaultFovTanHorizontal, defaultFovTanVertical
	if isScoped {
		fov := scopedFovDegrees
		if weapon != nil && weapon.ZoomLevel() == common.ZoomFull {
			fov = fullScopedFovDegrees
		}
		fovTanHorizontal, fovTanVertical = fovTangents(fov)
	}

	// Same as Player.PositionEyes() without reading the pawn position twice, the X and Y view offsets are always 0.
	origin := pawn.Position()
	eyes := origin
	if viewOffset, ok := pawn.PropertyValue("m_vecViewOffset.m_vecZ"); ok && viewOffset.Any != nil {
		eyes.Z += float64(viewOffset.Float())
	} else {
		eyes.Z += defaultEyesHeight
	}

	snapshot := playerSnapshot{
		recorded:         true,
		isAlive:          true,
		origin:           origin,
		eyes:             eyes,
		yaw:              float64(player.ViewDirectionX()),
		pitch:            normalizePitch(float64(player.ViewDirectionY())),
		fovTanHorizontal: fovTanHorizontal,
		fovTanVertical:   fovTanVertical,
		isBlinded:        player.FlashDurationTimeRemaining().Seconds() > blindedMinRemainingSeconds,
	}

	stance := visibility.StanceOf(snapshot.eyesHeight(), player.IsAirborne())
	isMoving := visibility.IsMoving(player.Velocity())
	snapshot.hitboxTemplate = visibility.GetHitboxTemplate(stance, isMoving, visibility.WeaponKindOf(weapon, isScoped))

	return snapshot
}

// normalizePitch converts a pitch from the [0, 360[ range to [-180, 180[.
func normalizePitch(pitch float64) float64 {
	if pitch >= 180 {
		return pitch - 360
	}

	return pitch
}

func (analyzer *Analyzer) registerVisibilityHandlers() {
	tracker := analyzer.visibilityTracker
	if tracker == nil {
		return
	}

	parser := analyzer.parser
	parser.RegisterEventHandler(func(event events.TickRateInfoAvailable) {
		tracker.setTickRate(parser.TickRate())
	})

	parser.RegisterEventHandler(func(event events.RoundStart) {
		tracker.reset()
	})

	parser.RegisterEventHandler(func(event events.RoundFreezetimeChanged) {
		if event.NewIsFreezetime {
			tracker.reset()
		}
	})

	parser.RegisterEventHandler(func(event events.SmokeStart) {
		tracker.addSmoke(event.GrenadeEntityID, event.Position)
	})

	parser.RegisterEventHandler(func(event events.SmokeExpired) {
		tracker.removeSmoke(event.GrenadeEntityID)
	})

	parser.RegisterEventHandler(func(event events.FrameDone) {
		if !analyzer.matchStarted() || parser.GameState().IsFreezetimePeriod() {
			return
		}

		tracker.recordFrame(analyzer.currentTick(), parser.GameState().Participants().Playing())
		tracker.resolvePendingEvents()
	})
}

func isTimeToDamageWeapon(weaponType constants.WeaponType) bool {
	switch weaponType {
	case constants.WeaponTypePistol, constants.WeaponTypeSMG, constants.WeaponTypeShotgun, constants.WeaponTypeRifle,
		constants.WeaponTypeSniper, constants.WeaponTypeMachineGun:
		return true
	default:
		return false
	}
}

func isEnemyDamage(attacker *common.Player, victim *common.Player) bool {
	return attacker != nil && victim != nil && attacker != victim && attacker.Team != victim.Team
}

func (tracker *visibilityTracker) onDamage(tick int, attacker *common.Player, victim *common.Player, damage *Damage) {
	if !isEnemyDamage(attacker, victim) || !isTimeToDamageWeapon(damage.WeaponType) {
		return
	}

	tracker.pending = append(tracker.pending, pendingVisibilityEvent{tick: tick, attacker: attacker, victim: victim, damage: damage})
}

func (tracker *visibilityTracker) onKill(tick int, attacker *common.Player, victim *common.Player, kill *Kill) {
	if !isEnemyDamage(attacker, victim) || !isTimeToDamageWeapon(kill.WeaponType) {
		return
	}

	tracker.pending = append(tracker.pending, pendingVisibilityEvent{tick: tick, attacker: attacker, victim: victim, kill: kill})
}

// resolvePendingEvents computes the time to damage/kill of the events that occurred until the last recorded frame,
// damages first since a kill is preceded by its damage at the same tick.
func (tracker *visibilityTracker) resolvePendingEvents() {
	last := tracker.lastFrame()
	if last == nil {
		return
	}

	remaining := tracker.pending[:0]
	for _, event := range tracker.pending {
		if event.tick > last.tick {
			remaining = append(remaining, event)
			continue
		}

		if event.damage != nil {
			tracker.resolveDamage(event)
		} else {
			tracker.resolveKill(event)
		}
	}
	tracker.pending = remaining
}

func (tracker *visibilityTracker) ticksToMilliseconds(ticks int) int {
	return int(math.Round(float64(ticks) / tracker.tickRate * 1000))
}

// timeToDamageStat aggregates times to damage/kill with stat, noTimeToDamage when there are none.
func timeToDamageStat(times []int, stat func([]int) float32) float32 {
	if len(times) == 0 {
		return noTimeToDamage
	}

	return stat(times)
}

func (tracker *visibilityTracker) resolveDamage(event pendingVisibilityEvent) {
	pair := playerPair{attacker: event.attacker, victim: event.victim}
	previous, hasPrevious := tracker.engagements[pair]
	stopTick := -1
	if hasPrevious {
		stopTick = previous.firstDamageTick
	}

	status, spottedTick := tracker.findSpottedTick(event.tick, event.attacker, event.victim, stopTick)
	switch status {
	case spotted:
		tracker.engagements[pair] = engagement{firstDamageTick: event.tick, spottedTick: spottedTick}
		event.damage.TimeToDamage = tracker.ticksToMilliseconds(event.tick - spottedTick)
	case visibleForTooLong:
		tracker.engagements[pair] = engagement{firstDamageTick: event.tick, spottedTick: -1}
	case reachedPreviousEngagement, notVisible:
		// Not the first damage of the engagement or damage done without seeing the enemy (wallbang, through a smoke...)
	}
}

func (tracker *visibilityTracker) resolveKill(event pendingVisibilityEvent) {
	pair := playerPair{attacker: event.attacker, victim: event.victim}
	previous, hasPrevious := tracker.engagements[pair]
	if !hasPrevious || previous.spottedTick == -1 {
		return
	}

	status, _ := tracker.findSpottedTick(event.tick, event.attacker, event.victim, previous.firstDamageTick)
	if status == reachedPreviousEngagement {
		event.kill.TimeToKill = tracker.ticksToMilliseconds(event.tick - previous.spottedTick)
	}
}

// findSpottedTick walks back in time from the event tick while the victim is visible by the attacker.
// It returns reachedPreviousEngagement if the victim has been visible since stopTick (-1 to ignore it).
func (tracker *visibilityTracker) findSpottedTick(eventTick int, attacker *common.Player, victim *common.Player, stopTick int) (visibilityStatus, int) {
	attackerSlot, ok := tracker.slots[attacker]
	if !ok {
		return notVisible, 0
	}
	victimSlot, ok := tracker.slots[victim]
	if !ok {
		return notVisible, 0
	}

	gapToleranceTicks := int(math.Round(visibilityGapToleranceSeconds * tracker.tickRate))
	lastVisibleTick := -1
	for age := range tracker.frameCount {
		frame := tracker.frameAt(age)
		if frame.tick > eventTick {
			continue
		}

		if tracker.isVisible(frame, attackerSlot, victimSlot, frame.tick == eventTick) {
			lastVisibleTick = frame.tick
			if stopTick != -1 && frame.tick <= stopTick {
				return reachedPreviousEngagement, 0
			}
			continue
		}

		referenceTick := eventTick
		if lastVisibleTick != -1 {
			referenceTick = lastVisibleTick
		}
		if referenceTick-frame.tick > gapToleranceTicks {
			break
		}
	}

	if lastVisibleTick == -1 {
		return notVisible, 0
	}

	// The victim was already visible in the oldest frame, we don't know when they have been spotted.
	if oldest := tracker.frameAt(tracker.frameCount - 1); lastVisibleTick-oldest.tick <= gapToleranceTicks {
		return visibleForTooLong, 0
	}

	return spotted, lastVisibleTick
}

// isVisible returns true if any point of the victim's body is visible on the attacker's screen.
// Players may already be dead at the tick of the event (e.g. kill or trade at the same tick), their last state is used.
func (tracker *visibilityTracker) isVisible(frame *visibilityFrame, attackerSlot int, victimSlot int, isEventTick bool) bool {
	if attackerSlot >= len(frame.snapshots) || victimSlot >= len(frame.snapshots) {
		return false
	}

	attacker := &frame.snapshots[attackerSlot]
	victim := &frame.snapshots[victimSlot]
	if !attacker.recorded || !victim.recorded {
		return false
	}
	if !isEventTick && (!attacker.isAlive || !victim.isAlive) {
		return false
	}
	if attacker.isBlinded {
		return false
	}

	forward, right, up := attacker.viewBasis()
	for _, point := range victim.hitboxPoints() {
		position := victim.hitboxPoint(point)
		if !attacker.isInFieldOfView(position, forward, right, up) {
			continue
		}

		isBehindSmoke := false
		for _, smoke := range frame.smokes {
			if segmentIntersectsSmoke(attacker.eyes, position, smoke) {
				isBehindSmoke = true
				break
			}
		}
		if isBehindSmoke {
			continue
		}

		if !tracker.tracer.IsOccluded(attacker.eyes, position) {
			return true
		}
	}

	return false
}
