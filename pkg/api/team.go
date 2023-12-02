package api

import (
	"github.com/akiver/cs-demo-analyzer/pkg/api/constants"
	common "github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/common"
)

type Team struct {
	Name            string               `json:"name"`
	Letter          constants.TeamLetter `json:"letter"`
	Score           int                  `json:"score"`
	ScoreFirstHalf  int                  `json:"scoreFirstHalf"`
	ScoreSecondHalf int                  `json:"scoreSecondHalf"`
	CurrentSide     *common.Team         `json:"currentSide"`
}

func (team *Team) swap() {
	if *team.CurrentSide == common.TeamCounterTerrorists {
		*team.CurrentSide = common.TeamTerrorists
	} else if *team.CurrentSide == common.TeamTerrorists {
		*team.CurrentSide = common.TeamCounterTerrorists
	}
}
