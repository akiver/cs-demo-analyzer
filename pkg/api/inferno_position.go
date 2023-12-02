package api

import (
	"fmt"

	"github.com/golang/geo/r2"
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/common"
)

type InfernoPosition struct {
	Frame            int        `json:"frame"`
	Tick             int        `json:"tick"`
	RoundNumber      int        `json:"roundNumber"`
	ThrowerSteamID64 uint64     `json:"throwerSteamId"`
	ThrowerName      string     `json:"throwerName"`
	UniqueID         int64      `json:"uniqueId"`
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
		UniqueID:         inferno.UniqueID(),
		ThrowerSteamID64: thrower.SteamID64,
		ThrowerName:      thrower.Name,
		X:                inferno.Entity.Position().X,
		Y:                inferno.Entity.Position().Y,
		Z:                inferno.Entity.Position().Z,
		ConvexHull2D:     inferno.Fires().Active().ConvexHull2D(),
	}
}
