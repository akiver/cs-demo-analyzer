package api

import (
	"github.com/akiver/cs-demo-analyzer/pkg/api/constants"
)

type Team struct {
	Name            string               `json:"name"`
	Letter          constants.TeamLetter `json:"letter"`
	Score           int                  `json:"score"`
	ScoreFirstHalf  int                  `json:"scoreFirstHalf"`
	ScoreSecondHalf int                  `json:"scoreSecondHalf"`
	CurrentSide     *constants.Side      `json:"currentSide"`
}

func (team *Team) swap() {
	if *team.CurrentSide == constants.SideCounterTerrorists {
		*team.CurrentSide = constants.SideTerrorists
	} else if *team.CurrentSide == constants.SideTerrorists {
		*team.CurrentSide = constants.SideCounterTerrorists
	}
}
