package api

import (
	"encoding/json"
	"strings"

	"github.com/akiver/cs-demo-analyzer/pkg/api/constants"
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/common"
	events "github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/events"
)

type Round struct {
	analyzer            *Analyzer
	Number              int                   `json:"number"`
	StartTick           int                   `json:"startTick"`
	StartFrame          int                   `json:"startFrame"`
	FreezeTimeEndTick   int                   `json:"freezeTimeEndTick"`
	FreezeTimeEndFrame  int                   `json:"freezeTimeEndFrame"`
	EndTick             int                   `json:"endTick"`
	EndFrame            int                   `json:"endFrame"`
	EndOfficiallyTick   int                   `json:"endOfficiallyTick"`
	EndOfficiallyFrame  int                   `json:"endOfficiallyFrame"`
	OvertimeNumber      int                   `json:"overtimeNumber"`
	TeamAName           string                `json:"teamAName"`
	TeamBName           string                `json:"teamBName"`
	TeamAScore          int                   `json:"teamAScore"`
	TeamBScore          int                   `json:"teamBScore"`
	TeamASide           common.Team           `json:"teamASide"`
	TeamBSide           common.Team           `json:"teamBSide"`
	TeamAEquipmentValue int                   `json:"teamAEquipmentValue"`
	TeamBEquipmentValue int                   `json:"teamBEquipmentValue"`
	TeamAMoneySpent     int                   `json:"teamAMoneySpent"`
	TeamBMoneySpent     int                   `json:"teamBmoneySpent"`
	TeamAEconomyType    constants.EconomyType `json:"teamAEconomyType"`
	TeamBEconomyType    constants.EconomyType `json:"teamBEconomyType"`
	Duration            int64                 `json:"duration"`
	EndReason           events.RoundEndReason `json:"endReason"`
	WinnerName          string                `json:"winnerName"`
	WinnerSide          common.Team           `json:"winnerSide"`
	// Used to detect weapons bought by players during buy time.
	// There is no "player buy" event available, instead we use the "item_pickup" event which occurs when a player pickup a weapon at any time in the game.
	// Since it's possible to buy and drop a weapon to a teammate and so trigger a new "item_pickup" event, it would result in a wrong weapon buy detection.
	// To prevent it, this slice contains all unique weapons ids that players had at the beginning of the round and it's updated when a player buy a weapon.
	weaponsBoughtUniqueIds []string `json:"-"`
}

type RoundAlias Round

type RoundJSON struct {
	*RoundAlias
	StartMoneyTeamA int `json:"teamAStartMoney"`
	StartMoneyTeamB int `json:"teamBStartMoney"`
}

func (round *Round) MarshalJSON() ([]byte, error) {
	return json.Marshal(RoundJSON{
		RoundAlias:      (*RoundAlias)(round),
		StartMoneyTeamA: round.StartMoneyTeamA(),
		StartMoneyTeamB: round.StartMoneyTeamB(),
	})
}

func (round *Round) StartMoneyTeamA() int {
	return round.getTeamStartMoney(round.TeamASide)
}

func (round *Round) StartMoneyTeamB() int {
	return round.getTeamStartMoney(round.TeamBSide)
}

// This indicates if x seconds have passed since the beginning of the round.
func (round *Round) secondsPassedSinceRoundStart(seconds int) bool {
	parser := round.analyzer.parser

	return float64(parser.GameState().IngameTick()-round.StartTick)*parser.TickTime().Seconds() >= float64(seconds)
}

func (round *Round) getTeamStartMoney(side common.Team) int {
	total := 0
	for _, economy := range round.analyzer.match.PlayerEconomies {
		if economy.RoundNumber == round.Number && economy.PlayerSide == side {
			total += economy.StartMoney
		}
	}

	return total
}

func (round *Round) computeTeamsEconomy() {
	analyzer := round.analyzer
	match := analyzer.match
	gameState := analyzer.parser.GameState()

	stateTeamA := gameState.Team(*match.TeamA.CurrentSide)
	round.TeamAMoneySpent = stateTeamA.MoneySpentThisRound()
	round.TeamAEquipmentValue = stateTeamA.CurrentEquipmentValue()
	round.TeamAEconomyType = computeTeamEconomyType(analyzer, stateTeamA)

	stateTeamB := gameState.Team(*match.TeamB.CurrentSide)
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

func (round *Round) updateEndReasonFromRoundEndEvent(event events.RoundEnd, mapName string) {
	round.EndReason = event.Reason

	if round.EndReason == roundEndReasonUnassigned {
		round.EndReason = getEndReasonFromRoundEndMessage(event.Message)
	}

	// Encountered with a demo from 2014.
	// When CTs won because the bomb exploded and the Ts saved their weapons, the reason may be "HostagesRescued" instead of "TargetSaved" on defuse maps.
	if round.EndReason == events.RoundEndReasonHostagesRescued && strings.HasPrefix(mapName, "de_") {
		round.EndReason = events.RoundEndReasonTargetSaved
	}
}

func newRound(number int, analyzer *Analyzer) *Round {
	// Consider all current players weapons as "already bought" so if a player drop his weapon to a teammate, it will not be detected as a buy.
	var weaponsIds []string
	for _, playingPlayer := range analyzer.parser.GameState().Participants().Playing() {
		for _, weapon := range playingPlayer.Weapons() {
			weaponsIds = append(weaponsIds, weapon.UniqueID2().String())
		}
	}

	round := &Round{
		analyzer:               analyzer,
		Number:                 number,
		StartFrame:             analyzer.parser.CurrentFrame(),
		StartTick:              analyzer.currentTick(),
		FreezeTimeEndFrame:     -1,
		FreezeTimeEndTick:      -1,
		OvertimeNumber:         analyzer.match.OvertimeCount,
		TeamAName:              analyzer.match.TeamA.Name,
		TeamAScore:             analyzer.match.TeamA.Score,
		TeamBName:              analyzer.match.TeamB.Name,
		TeamBScore:             analyzer.match.TeamB.Score,
		TeamASide:              *analyzer.match.TeamA.CurrentSide,
		TeamBSide:              *analyzer.match.TeamB.CurrentSide,
		weaponsBoughtUniqueIds: weaponsIds,
	}
	round.computeTeamsEconomy()

	return round
}
