package api

import (
	"github.com/akiver/cs-demo-analyzer/pkg/api/constants"
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/common"
	events "github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/events"
	sendtables "github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/sendtables"
)

func createMatchZyAnalyzer(analyzer *Analyzer) {
	match := analyzer.match
	match.gameModeStr = constants.GameModeStrCompetitive
	parser := analyzer.parser
	matchStarted := true
	// Track the end of the game to not re-compute the last round players economy as the freeze time start event is
	// triggered before the final round end event.
	gameEndTick := -1
	// Used to detect when the game is paused and then resumed after a backup restore
	isPausedDueToBackupRestore := false
	analyzer.matchStarted = func() bool {
		return matchStarted
	}

	parser.RegisterEventHandler(func(events.DataTablesParsed) {
		serverClasses := parser.ServerClasses()
		serverClasses.FindByName("CCSTeam").OnEntityCreated(func(entity sendtables.Entity) {
			entity.Property("m_szClanTeamname").OnUpdate(func(value sendtables.PropertyValue) {
				analyzer.updateTeamNames()
			})
		})

		serverClasses.FindByName("CCSGameRulesProxy").OnEntityCreated(func(entity sendtables.Entity) {
			entity.Property("m_pGameRules.m_eRoundEndReason").OnUpdate(func(value sendtables.PropertyValue) {
				if value.Int() == int(events.RoundEndReasonDraw) {
					roundEndWinnerTeam := entity.PropertyValueMust("m_pGameRules.m_iRoundEndWinnerTeam").Int()
					if roundEndWinnerTeam == int(common.TeamSpectators) {
						// backup restore, the match is now paused, the match will resume when m_bMatchWaitingForResume is set to false
						matchStarted = false
						isPausedDueToBackupRestore = true
						roundPlayedCount := entity.PropertyValueMust("m_pGameRules.m_totalRoundsPlayed").Int()
						if roundPlayedCount == 0 {
							analyzer.reset()
						} else {
							analyzer.resetCurrentRound()
						}
					}
				}
			})

			entity.Property("m_pGameRules.m_bMatchWaitingForResume").OnUpdate(func(value sendtables.PropertyValue) {
				// Resume the match when the backup restore is done
				if isPausedDueToBackupRestore && !value.BoolVal() {
					matchStarted = true
					isPausedDueToBackupRestore = false
					currentRound := analyzer.currentRound
					currentRound.StartFrame = parser.CurrentFrame()
					currentRound.StartTick = analyzer.currentTick()
				}
			})

			entity.Property("m_pGameRules.m_gamePhase").OnUpdate(func(value sendtables.PropertyValue) {
				if value.Int() == int(common.GamePhaseGameEnded) {
					gameEndTick = analyzer.currentTick()
				}
			})
		})

		parser.RegisterEventHandler(func(event events.MatchStartedChanged) {
			analyzer.registerUnknownPlayers()
			matchStarted = event.NewIsStarted
			if matchStarted {
				analyzer.reset()
				currentRound := analyzer.currentRound
				currentRound.StartFrame = parser.CurrentFrame()
				currentRound.StartTick = analyzer.currentTick()
			} else {
				matchStarted = false
			}
		})

		parser.RegisterEventHandler(func(event events.RoundFreezetimeChanged) {
			if !analyzer.matchStarted() || gameEndTick == analyzer.currentTick() {
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
	})

	parser.RegisterEventHandler(func(event events.RoundStart) {
		if !analyzer.matchStarted() || analyzer.currentTick() == 0 || len(match.Rounds) == 0 {
			return
		}

		analyzer.createRound()
	})

	parser.RegisterEventHandler(analyzer.defaultRoundEndOfficiallyHandler)

	parser.RegisterEventHandler(analyzer.defaultAnnouncementWinPanelMatchHandler)
}
