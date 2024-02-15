package api

import (
	"github.com/akiver/cs-demo-analyzer/pkg/api/constants"
	events "github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/events"
)

func createEsportalAnalyzer(analyzer *Analyzer) {
	match := analyzer.match
	match.gameModeStr = constants.GameModeStrCompetitive
	parser := analyzer.parser
	isMatchStarted := false

	analyzer.matchStarted = func() bool {
		return isMatchStarted
	}

	parser.RegisterEventHandler(func(event events.MatchStartedChanged) {
		if !event.OldIsStarted && event.NewIsStarted {
			isMatchStarted = true
			analyzer.processMatchStart()
		}
	})

	parser.RegisterEventHandler(analyzer.defaultRoundFreezetimeChangedHandler)

	parser.RegisterEventHandler(analyzer.defaultRoundStartHandler)

	parser.RegisterEventHandler(analyzer.defaultRoundEndOfficiallyHandler)

	parser.RegisterEventHandler(func(event events.AnnouncementWinPanelMatch) {
		analyzer.updatePlayersScores()
	})
}
