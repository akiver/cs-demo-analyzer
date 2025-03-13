package api

import (
	events "github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/events"
)

func createEsplayAnalyzer(analyzer *Analyzer) {
	parser := analyzer.parser
	matchStarted := false

	analyzer.matchStarted = func() bool {
		return matchStarted
	}

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
