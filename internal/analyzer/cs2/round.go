package cs2

import (
	"github.com/akiver/cs-demo-analyzer/internal/modelaccess"
	"strings"

	"github.com/akiver/cs-demo-analyzer/pkg/api/constants"
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/common"
	events "github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/events"
)

// This indicates if x seconds have passed since the beginning of the round.
func (analyzer *Analyzer) secondsPassedSinceRoundStart(round *Round, seconds int) bool {
	parser := analyzer.parser

	return float64(parser.GameState().IngameTick()-round.StartTick)*parser.TickTime().Seconds() >= float64(seconds)
}

func (analyzer *Analyzer) computeTeamsEconomy(round *Round) {
	match := analyzer.match
	gameState := analyzer.parser.GameState()

	stateTeamA := gameState.Team(common.Team(*match.TeamA.CurrentSide))
	round.TeamAMoneySpent = stateTeamA.MoneySpentThisRound()
	round.TeamAEquipmentValue = stateTeamA.CurrentEquipmentValue()
	round.TeamAEconomyType = computeTeamEconomyType(analyzer, stateTeamA)

	stateTeamB := gameState.Team(common.Team(*match.TeamB.CurrentSide))
	round.TeamBMoneySpent = stateTeamB.MoneySpentThisRound()
	round.TeamBEquipmentValue = stateTeamB.CurrentEquipmentValue()
	round.TeamBEconomyType = computeTeamEconomyType(analyzer, stateTeamB)
}

// Possible round end reason value reported on round_end event which means that round end reason has not been detected properly.
const roundEndReasonUnassigned = 0

// This returns the RoundEndReason from the round_end string message event.
// Necessary for old demos that contains unassigned end round reason value (i.e. 0).
func getEndReasonFromRoundEndMessage(message string) events.RoundEndReason {
	switch message {
	case "#SFUI_Notice_Target_Saved":
		return events.RoundEndReasonTargetSaved
	case "#SFUI_Notice_Target_Bombed":
		return events.RoundEndReasonTargetBombed
	default:
		return roundEndReasonUnassigned
	}
}

func updateRoundEndReason(round *Round, event events.RoundEnd, mapName string) {
	round.EndReason = constants.RoundEndReason(event.Reason)

	if round.EndReason == roundEndReasonUnassigned {
		round.EndReason = constants.RoundEndReason(getEndReasonFromRoundEndMessage(event.Message))
	}

	// Encountered with a demo from 2014.
	// When CTs won because the bomb exploded and the Ts saved their weapons, the reason may be "HostagesRescued" instead of "TargetSaved" on defuse maps.
	if round.EndReason == constants.RoundEndReasonHostagesRescued && strings.HasPrefix(mapName, "de_") {
		round.EndReason = constants.RoundEndReasonTargetSaved
	}
}

func newRound(number int, analyzer *Analyzer) *Round {
	// Consider all current players weapons as "already bought" so if a player drop his weapon to a teammate, it will not be detected as a buy.
	var weaponsIds []string
	for _, playingPlayer := range analyzer.parser.GameState().Participants().Playing() {
		for _, weapon := range playingPlayer.Weapons() {
			weaponsIds = append(weaponsIds, weapon.UniqueID().String())
		}
	}

	analyzer.weaponsBoughtUniqueIDs = weaponsIds

	round := &Round{
		Number:             number,
		StartFrame:         analyzer.parser.CurrentFrame(),
		StartTick:          analyzer.currentTick(),
		FreezeTimeEndFrame: -1,
		FreezeTimeEndTick:  -1,
		OvertimeNumber:     analyzer.match.OvertimeCount,
		TeamAName:          analyzer.match.TeamA.Name,
		TeamAScore:         analyzer.match.TeamA.Score,
		TeamBName:          analyzer.match.TeamB.Name,
		TeamBScore:         analyzer.match.TeamB.Score,
		TeamASide:          *analyzer.match.TeamA.CurrentSide,
		TeamBSide:          *analyzer.match.TeamB.CurrentSide,
	}
	modelaccess.SetRoundMatch(round, analyzer.match)
	analyzer.computeTeamsEconomy(round)

	return round
}
