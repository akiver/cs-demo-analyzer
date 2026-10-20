package model

import (
	"github.com/akiver/cs-demo-analyzer/pkg/api/constants"
)

type Coach struct {
	SteamID64     uint64               `json:"steamId"`
	UserID        int                  `json:"userId"`
	Name          string               `json:"name"`
	CoachingSide  constants.Side       `json:"coachingSide"` // The side being coached when the coach was first seen
	TeamLetter    constants.TeamLetter `json:"teamLetter"`   // The team being coached when the coach was first seen, empty if unknown
	FirstSeenTick int                  `json:"firstSeenTick"`
}
