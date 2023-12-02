package api

import (
	"github.com/akiver/cs-demo-analyzer/pkg/api/constants"
	events "github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/events"
)

func createFaceItAnalyzer(analyzer *Analyzer) {
	parser := analyzer.parser
	match := analyzer.match
	match.gameModeStr = constants.GameModeStrCompetitive
	matchStarted := false

	analyzer.matchStarted = func() bool {
		return matchStarted
	}

	parser.RegisterEventHandler(func(event events.MatchStartedChanged) {
		if !event.OldIsStarted && event.NewIsStarted && !analyzer.isKnifeRound() {
			matchStarted = true

			for _, participant := range parser.GameState().Participants().Playing() {
				for _, player := range match.PlayersBySteamID {
					if player.SteamID64 == participant.SteamID64 {
						var team *Team
						if *match.TeamA.CurrentSide == participant.Team {
							team = match.TeamA
						} else {
							team = match.TeamB
						}
						player.Team = team
					}
				}
			}

			currentRound := analyzer.currentRound
			currentRound.StartFrame = parser.CurrentFrame()
			currentRound.StartTick = analyzer.currentTick()
			currentRound.TeamASide = *match.TeamA.CurrentSide
			currentRound.TeamBSide = *match.TeamB.CurrentSide
			analyzer.updateTeamNames()
			analyzer.createPlayersEconomies()
		}
	})

	parser.RegisterEventHandler(func(event events.RoundFreezetimeChanged) {
		if !analyzer.matchStarted() {
			return
		}

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
		matchStarted = false
	})

	analyzer.postProcess = func() {
		currentRound := analyzer.currentRound
		if len(match.Rounds) < currentRound.Number {
			match.Rounds = append(match.Rounds, currentRound)
		}
	}
}
