package model

import (
	"github.com/golang/geo/r2"
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
