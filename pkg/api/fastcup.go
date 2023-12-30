package api

import (
	events "github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/events"
	st "github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/sendtables"
)

func createFastcupAnalyzer(analyzer *Analyzer) {
	parser := analyzer.parser
	matchStarted := false

	analyzer.matchStarted = func() bool {
		return matchStarted
	}

	parser.RegisterEventHandler(func(events.DataTablesParsed) {
		parser.ServerClasses().FindByName("CCSGameRulesProxy").OnEntityCreated(func(entity st.Entity) {
			// Match start detection - the match really starts when:
			// 1. The warmup period is over
			// 2. Players can buy, otherwise it means it's the knife round
			var cantBuyProp st.Property
			if analyzer.isSource2 {
				cantBuyProp = entity.Property("m_pGameRules.m_bTCantBuy")
			} else {
				cantBuyProp = entity.Property("cs_gamerules_data.m_bTCantBuy")
			}

			cantBuyProp.OnUpdate(func(val st.PropertyValue) {
				if parser.GameState().IsWarmupPeriod() {
					return
				}

				playersCanBuy := !val.BoolVal()
				if playersCanBuy && !analyzer.matchStarted() {
					matchStarted = true
					analyzer.processMatchStart()
				}
			})
		})
	})

	parser.RegisterEventHandler(analyzer.defaultRoundFreezetimeChangedHandler)

	parser.RegisterEventHandler(analyzer.defaultRoundStartHandler)

	parser.RegisterEventHandler(analyzer.defaultRoundEndOfficiallyHandler)

	parser.RegisterEventHandler(func(event events.AnnouncementWinPanelMatch) {
		analyzer.updatePlayersScores()
		matchStarted = false
	})
}
