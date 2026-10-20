package model

import (
	"github.com/akiver/cs-demo-analyzer/pkg/api/constants"
)

type HostagePosition struct {
	// Used to identify the hostage, entity IDs are re-used once an entity is destroyed but not with the same serial number.
	// Not exported because they are used only to group positions by hostage.
	entityID    int
	serialNum   int
	Frame       int                    `json:"frame"`
	Tick        int                    `json:"tick"`
	RoundNumber int                    `json:"roundNumber"`
	X           float64                `json:"x"`
	Y           float64                `json:"y"`
	Z           float64                `json:"z"`
	State       constants.HostageState `json:"state"`
}
