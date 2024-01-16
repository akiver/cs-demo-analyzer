package api

import (
	events "github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/events"
	st "github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/sendtables"
)

func createFiveEPlayAnalyzer(analyzer *Analyzer) {
	parser := analyzer.parser
	isRestarting := false
	isMatchStarted := false
	analyzer.matchStarted = func() bool {
		return isMatchStarted
	}

	parser.RegisterEventHandler(func(event events.IsWarmupPeriodChanged) {
		// Match start detection - if there is a knife round the detection will be done with the m_bGameRestart prop.
		if event.OldIsWarmupPeriod && !event.NewIsWarmupPeriod {
			isKnifeRound := false
			for _, player := range analyzer.parser.GameState().Participants().Playing() {
				if player.Money() == 0 {
					isKnifeRound = true
					break
				}
			}

			if isKnifeRound {
				return
			}

			isMatchStarted = true
			analyzer.processMatchStart()
		}
	})

	parser.RegisterEventHandler(func(events.DataTablesParsed) {
		parser.ServerClasses().FindByName("CCSGameRulesProxy").OnEntityCreated(func(entity st.Entity) {
			// Match start detection when there was a knife round - the game really starts after a restart.
			var gameRestartProp st.Property
			if analyzer.isSource2 {
				gameRestartProp = entity.Property("m_pGameRules.m_bGameRestart")
			} else {
				gameRestartProp = entity.Property("cs_gamerules_data.m_bGameRestart")
			}

			gameRestartProp.OnUpdate(func(val st.PropertyValue) {
				newIsRestarting := val.BoolVal()
				if isRestarting && !newIsRestarting {
					isMatchStarted = true
					analyzer.processMatchStart()
				}
				isRestarting = newIsRestarting
			})
		})
	})

	parser.RegisterEventHandler(func(event events.GameHalfEnded) {
		analyzer.isFirstRoundOfHalf = true
	})

	parser.RegisterEventHandler(analyzer.defaultRoundFreezetimeChangedHandler)

	parser.RegisterEventHandler(analyzer.defaultRoundStartHandler)

	parser.RegisterEventHandler(analyzer.defaultRoundEndOfficiallyHandler)

	parser.RegisterEventHandler(func(event events.MatchStartedChanged) {
		// Match end detection
		if isMatchStarted && !event.NewIsStarted {
			isMatchStarted = false
		}
	})

	parser.RegisterEventHandler(func(event events.AnnouncementWinPanelMatch) {
		analyzer.updatePlayersScores()
		isMatchStarted = false
	})
}
