package api

import (
	"slices"

	common "github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/common"
)

// tickPositions holds the entities whose position has been captured during the current tick, so that the position of
// an entity is captured once per tick.
//
// The parser dispatches a frame for every command of a demo and with CS2 demos several commands may have the same
// tick (4 to 21% of the frames of the demos it has been measured on, none with CS:GO demos). The entities didn't
// change in those extra frames, capturing them again only created duplicated positions.
// It's tracked by entity rather than by frame so that an entity which appears in an extra frame still has a position
// for that tick.
type tickPositions struct {
	tick int
	// The round may change between 2 frames of a tick, a position is identified by both.
	roundNumber int
	// Whether positions have already been captured for this tick, the lists are searched only in that case.
	isExtraFrame bool
	players      []*common.Player
	projectiles  []*common.GrenadeProjectile
	infernos     []*common.Inferno
	chickens     []*chickenEntity
	hostages     []*common.Hostage
}

// startFrame must be called before capturing the positions of a frame.
func (positions *tickPositions) startFrame(tick int, roundNumber int) {
	if tick == positions.tick && roundNumber == positions.roundNumber {
		positions.isExtraFrame = true
		return
	}

	positions.reset()
	positions.tick = tick
	positions.roundNumber = roundNumber
}

// reset forgets the positions captured during the current tick, it must be called when positions are deleted from
// the match so that the next frame captures them again even if it has the same tick.
func (positions *tickPositions) reset() {
	positions.isExtraFrame = false
	positions.players = positions.players[:0]
	positions.projectiles = positions.projectiles[:0]
	positions.infernos = positions.infernos[:0]
	positions.chickens = positions.chickens[:0]
	positions.hostages = positions.hostages[:0]
}

// hasPosition returns whether the position of entity has already been captured during the current tick, captured
// being the list of tickPositions for this kind of entity.
func hasPosition[T comparable](positions *tickPositions, captured []T, entity T) bool {
	return positions.isExtraFrame && slices.Contains(captured, entity)
}
