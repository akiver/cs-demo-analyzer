package csgo

import (
	"math"

	"github.com/akiver/cs-demo-analyzer/internal/strings"
	"github.com/akiver/cs-demo-analyzer/pkg/api/constants"
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/common"
)

func (analyzer *Analyzer) registerCoach(player *common.Player, side common.Team) {
	if player.SteamID64 == 0 {
		return
	}

	match := analyzer.match
	for _, coach := range match.Coaches {
		if coach.SteamID64 == player.SteamID64 {
			return
		}
	}

	var teamLetter constants.TeamLetter
	if side == common.TeamCounterTerrorists || side == common.TeamTerrorists {
		teamLetter = match.Team(constants.Side(side)).Letter
	}

	userID := 0
	if player.UserID <= math.MaxUint16 {
		userID = player.UserID & 0xff
	}

	match.Coaches = append(match.Coaches, &Coach{
		SteamID64:     player.SteamID64,
		UserID:        userID,
		Name:          strings.ReplaceUTF8ByteSequences(player.Name),
		CoachingSide:  constants.Side(side),
		TeamLetter:    teamLetter,
		FirstSeenTick: analyzer.currentTick(),
	})
}

// removeCoach removes the coach with the given SteamID, if any.
// It's called when a person is registered as a player so that players and coaches are mutually exclusive, a person
// may briefly be flagged as a coach during warmup before playing the match.
func (analyzer *Analyzer) removeCoach(steamID uint64) {
	coaches := analyzer.match.Coaches
	for i, coach := range coaches {
		if coach.SteamID64 == steamID {
			analyzer.match.Coaches = append(coaches[:i], coaches[i+1:]...)
			return
		}
	}
}
