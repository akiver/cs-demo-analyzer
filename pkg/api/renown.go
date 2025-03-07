package api

import (
	"github.com/akiver/cs-demo-analyzer/pkg/api/constants"
	events "github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/events"
	st "github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/sendtables"
)

func createRenownAnalyzer(analyzer *Analyzer) {
	parser := analyzer.parser
	match := analyzer.match
	match.gameModeStr = constants.GameModeStrCompetitive
	matchStarted := false
	isGamePaused := false

	analyzer.matchStarted = func() bool {
		return matchStarted
	}

	parser.RegisterEventHandler(func(events.DataTablesParsed) {
		parser.ServerClasses().FindByName("CCSGameRulesProxy").OnEntityCreated(func(entity st.Entity) {
			entity.Property("m_pGameRules.m_bTechnicalTimeOut").OnUpdate(func(val st.PropertyValue) {
				isGamePaused = val.BoolVal()
			})
		})
	})

	parser.RegisterEventHandler(func(event events.PlayerConnect) {
		if isGamePaused {
			analyzer.createOrUpdatePlayerEconomy(event.Player)
		}
	})

	parser.RegisterEventHandler(func(event events.MatchStartedChanged) {
		matchStarted = event.NewIsStarted
		if matchStarted {
			analyzer.processMatchStart()
		}
	})

	parser.RegisterEventHandler(analyzer.defaultRoundFreezetimeChangedHandler)

	parser.RegisterEventHandler(analyzer.defaultRoundStartHandler)

	parser.RegisterEventHandler(analyzer.defaultRoundEndOfficiallyHandler)

	parser.RegisterEventHandler(analyzer.defaultAnnouncementWinPanelMatchHandler)
}
