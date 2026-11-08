package api

import (
	common "github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/common"
)

// PlayerPositionSample is the position of an alive player taken every second, from the end of the freeze time to the
// end of the round. Unlike PlayerPosition which requires -positions to be set, samples are always exported: they are
// small enough to be kept.
type PlayerPositionSample struct {
	Tick        int         `json:"tick"`
	RoundNumber int         `json:"roundNumber"`
	SteamID64   uint64      `json:"steamId"`
	Name        string      `json:"name"`
	Side        common.Team `json:"side"`
	X           float64     `json:"x"`
	Y           float64     `json:"y"`
	Z           float64     `json:"z"`
}

// positionSampleSecond identifies a second of a round, see samplePlayerPositions().
type positionSampleSecond struct {
	roundNumber    int
	elapsedSeconds int
}

var noPositionSampleSecond = positionSampleSecond{roundNumber: -1, elapsedSeconds: -1}

func newPlayerPositionSample(analyzer *Analyzer, player *common.Player) *PlayerPositionSample {
	position := player.Position()

	return &PlayerPositionSample{
		Tick:        analyzer.currentTick(),
		RoundNumber: analyzer.currentRound.Number,
		SteamID64:   player.SteamID64,
		Name:        player.Name,
		Side:        player.Team,
		X:           position.X,
		Y:           position.Y,
		Z:           position.Z,
	}
}

// samplePlayerPositions adds a position sample of each alive player on the first frame of each second of the round.
// This function is called on every frame.
func (analyzer *Analyzer) samplePlayerPositions() {
	if !analyzer.matchStarted() {
		return
	}

	// Seconds are counted from the end of the freeze time so that the first sample is taken as soon as players can
	// move, or from the start of the round when the end of the freeze time is unknown (recording started after it).
	round := analyzer.currentRound
	originTick := round.FreezeTimeEndTick
	if originTick == -1 {
		originTick = round.StartTick
	}
	tick := analyzer.currentTick()
	if tick < originTick {
		return
	}

	// The frames of a second already sampled are skipped before touching the players. Seconds are compared instead of
	// ticks because some demos don't have a frame on every tick, the first frame of the second is used.
	second := positionSampleSecond{
		roundNumber:    round.Number,
		elapsedSeconds: (tick - originTick) / analyzer.ticksPerSecond,
	}
	if second == analyzer.lastPositionSampleSecond || analyzer.parser.GameState().IsFreezetimePeriod() {
		return
	}
	analyzer.lastPositionSampleSecond = second

	for _, player := range analyzer.parser.GameState().Participants().Playing() {
		if player.IsAlive() {
			analyzer.match.PlayerPositionSamples = append(analyzer.match.PlayerPositionSamples, newPlayerPositionSample(analyzer, player))
		}
	}
}
