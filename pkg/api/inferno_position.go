package api

import (
	"fmt"
	"slices"

	"github.com/golang/geo/r2"
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/common"
)

type InfernoPosition struct {
	Frame            int        `json:"frame"`
	Tick             int        `json:"tick"`
	RoundNumber      int        `json:"roundNumber"`
	ThrowerSteamID64 uint64     `json:"throwerSteamId"`
	ThrowerName      string     `json:"throwerName"`
	UniqueID         string     `json:"uniqueId"`
	X                float64    `json:"x"`
	Y                float64    `json:"y"`
	Z                float64    `json:"z"`
	ConvexHull2D     []r2.Point `json:"convexHull2D"`
}

func newInfernoPositionFromInferno(analyzer *Analyzer, inferno *common.Inferno) *InfernoPosition {
	thrower := inferno.Thrower()
	if thrower == nil {
		fmt.Println("Thrower nil in inferno")
		return nil
	}

	parser := analyzer.parser

	return &InfernoPosition{
		Frame:            parser.CurrentFrame(),
		Tick:             analyzer.currentTick(),
		RoundNumber:      analyzer.currentRound.Number,
		UniqueID:         inferno.UniqueID().String(),
		ThrowerSteamID64: thrower.SteamID64,
		ThrowerName:      thrower.Name,
		X:                inferno.Entity.Position().X,
		Y:                inferno.Entity.Position().Y,
		Z:                inferno.Entity.Position().Z,
		ConvexHull2D:     analyzer.infernoConvexHull2D(inferno),
	}
}

// infernoHull is the 2D convex hull of an inferno's active fires, see Analyzer.infernoConvexHull2D().
type infernoHull struct {
	fires []common.Fire // the active fires the hull was computed from
	hull  []r2.Point
}

// infernoConvexHull2D returns the 2D convex hull of the inferno's active fires.
// Computing the hull is expensive and the active fires only change when a fire ignites or extinguishes, so the hull
// is cached per inferno and recomputed only when they changed. The returned slice is shared, don't modify it.
func (analyzer *Analyzer) infernoConvexHull2D(inferno *common.Inferno) []r2.Point {
	activeFires := inferno.Fires().Active()
	fires := activeFires.List()

	cached, exists := analyzer.infernoHulls[inferno.UniqueID()]
	if exists && slices.Equal(cached.fires, fires) {
		return cached.hull
	}

	hull := activeFires.ConvexHull2D()
	analyzer.infernoHulls[inferno.UniqueID()] = &infernoHull{
		fires: fires,
		hull:  hull,
	}

	return hull
}
