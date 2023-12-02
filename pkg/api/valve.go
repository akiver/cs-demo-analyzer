package api

import (
	"github.com/akiver/cs-demo-analyzer/pkg/api/constants"
	events "github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/events"
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/msgs2"
)

func createValveAnalyzer(analyzer *Analyzer, demoFilePath string) {
	parser := analyzer.parser
	match := analyzer.match
	analyzer.matchStarted = parser.GameState().IsMatchStarted

	parser.RegisterNetMessageHandler(func(srvInfo *msgs2.CSVCMsg_ServerInfo) {
		match.gameModeStr = constants.GameModeStr(srvInfo.GameSessionConfig.GetGamemode())
	})

	parser.RegisterEventHandler(func(event events.MatchStart) {
		currentRound := analyzer.currentRound
		currentRound.StartFrame = parser.CurrentFrame()
		currentRound.StartTick = analyzer.currentTick()
		analyzer.updateTeamNames()
	})

	parser.RegisterEventHandler(func(event events.MatchStartedChanged) {
		isMatchStarted := !event.OldIsStarted && event.NewIsStarted
		if isMatchStarted {
			currentRound := analyzer.currentRound
			currentRound.TeamASide = *match.TeamA.CurrentSide
			currentRound.TeamBSide = *match.TeamB.CurrentSide
		}
	})

	parser.RegisterEventHandler(func(event events.GameHalfEnded) {
		analyzer.isFirstRoundOfHalf = true
	})

	parser.RegisterEventHandler(func(event events.RoundFreezetimeChanged) {
		// It may not be accurate to create players economy on round start because it's possible to buy
		// a few ticks before the round start event and so may results in incorrect values.
		// Do it when the freeze time starts, it's updated just before round start events.
		if event.NewIsFreezetime {
			analyzer.createPlayersEconomies()
		} else {
			analyzer.currentRound.FreezeTimeEndTick = analyzer.currentTick()
			analyzer.currentRound.FreezeTimeEndFrame = parser.CurrentFrame()
			analyzer.lastFreezeTimeEndTick = analyzer.currentTick()
		}
	})

	parser.RegisterEventHandler(func(event events.RoundStart) {
		if !analyzer.matchStarted() {
			return
		}

		// No Rounds have been added yet, don't create a new one in this case, it's still the first round.
		if len(match.Rounds) == 0 {
			return
		}

		analyzer.createRound()
	})

	parser.RegisterEventHandler(func(event events.RoundEndOfficial) {
		if !analyzer.matchStarted() {
			return
		}

		match.Rounds = append(match.Rounds, analyzer.currentRound)
	})

	parser.RegisterEventHandler(func(event events.AnnouncementWinPanelMatch) {
		analyzer.updatePlayersScores()
	})

	analyzer.postProcess = func() {
		currentRound := analyzer.currentRound
		if len(match.Rounds) < currentRound.Number {
			match.Rounds = append(match.Rounds, currentRound)
		}
	}
}
