// Package modelaccess gives the analyzers and pkg/api access to the operations of pkg/model that aren't part of its
// public API: resetting and finalising a match while it's analyzed, and the fields used only to group positions.
package modelaccess

import (
	"github.com/akiver/cs-demo-analyzer/internal/modelaccess/hooks"
	"github.com/akiver/cs-demo-analyzer/pkg/api/constants"
	"github.com/akiver/cs-demo-analyzer/pkg/model"
)

func ResetMatch(match *model.Match) {
	hooks.ResetMatch.(func(*model.Match))(match)
}

func ResetRound(match *model.Match, roundNumber int) {
	hooks.ResetRound.(func(*model.Match, int))(match, roundNumber)
}

func InitTeams(match *model.Match) {
	hooks.InitTeams.(func(*model.Match))(match)
}

func SwapTeams(match *model.Match) {
	hooks.SwapTeams.(func(*model.Match))(match)
}

func ComputeResultStats(match *model.Match) {
	hooks.ComputeResultStats.(func(*model.Match))(match)
}

func DeleteIncompleteRounds(match *model.Match) {
	hooks.DeleteIncompleteRounds.(func(*model.Match))(match)
}

// SetGameModeStr overrides the game mode string, for sources that run competitive matches.
func SetGameModeStr(match *model.Match, gameModeStr constants.GameModeStr) {
	hooks.SetGameModeStr.(func(*model.Match, constants.GameModeStr))(match, gameModeStr)
}

// MatchJSON returns the value the match is marshalled as.
func MatchJSON(match *model.Match) model.MatchJSON {
	return hooks.MatchJSON.(func(*model.Match) model.MatchJSON)(match)
}

func ResetPlayer(player *model.Player) {
	hooks.ResetPlayer.(func(*model.Player))(player)
}

// SetPlayerMatch sets the match the player belongs to, used for its computed values.
func SetPlayerMatch(player *model.Player, match *model.Match) {
	hooks.SetPlayerMatch.(func(*model.Player, *model.Match))(player, match)
}

// SetRoundMatch sets the match the round belongs to, used for its computed values.
func SetRoundMatch(round *model.Round, match *model.Match) {
	hooks.SetRoundMatch.(func(*model.Round, *model.Match))(round, match)
}

// ChickenID returns the unique ID of the chicken in the match, used to group positions by chicken.
func ChickenID(position *model.ChickenPosition) int {
	return hooks.ChickenID.(func(*model.ChickenPosition) int)(position)
}

func SetChickenID(position *model.ChickenPosition, chickenID int) {
	hooks.SetChickenID.(func(*model.ChickenPosition, int))(position, chickenID)
}

// HostageEntity returns the entity ID and serial number of the hostage, used to group positions by hostage.
func HostageEntity(position *model.HostagePosition) (entityID int, serialNum int) {
	return hooks.HostageEntity.(func(*model.HostagePosition) (int, int))(position)
}

func SetHostageEntity(position *model.HostagePosition, entityID int, serialNum int) {
	hooks.SetHostageEntity.(func(*model.HostagePosition, int, int))(position, entityID, serialNum)
}
