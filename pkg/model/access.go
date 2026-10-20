package model

import (
	"github.com/akiver/cs-demo-analyzer/internal/modelaccess/hooks"
	"github.com/akiver/cs-demo-analyzer/pkg/api/constants"
)

// The analyzers reach the unexported operations below through internal/modelaccess.
func init() {
	hooks.ResetMatch = (*Match).reset
	hooks.ResetRound = (*Match).resetRound
	hooks.InitTeams = (*Match).initTeams
	hooks.SwapTeams = func(match *Match) { match.swapTeams() }
	hooks.ComputeResultStats = (*Match).computeResultStats
	hooks.DeleteIncompleteRounds = (*Match).deleteIncompleteRounds
	hooks.SetGameModeStr = func(match *Match, gameModeStr constants.GameModeStr) { match.gameModeStr = gameModeStr }
	hooks.MatchJSON = (*Match).toJSON
	hooks.ResetPlayer = (*Player).reset
	hooks.SetPlayerMatch = func(player *Player, match *Match) { player.match = match }
	hooks.SetRoundMatch = func(round *Round, match *Match) { round.match = match }
	hooks.ChickenID = func(position *ChickenPosition) int { return position.chickenID }
	hooks.SetChickenID = func(position *ChickenPosition, chickenID int) { position.chickenID = chickenID }
	hooks.HostageEntity = func(position *HostagePosition) (int, int) { return position.entityID, position.serialNum }
	hooks.SetHostageEntity = func(position *HostagePosition, entityID int, serialNum int) {
		position.entityID = entityID
		position.serialNum = serialNum
	}
}
