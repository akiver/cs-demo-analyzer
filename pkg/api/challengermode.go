package api

import (
	"github.com/akiver/cs-demo-analyzer/pkg/api/constants"
	events "github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/events"
)

func createChallengermodeAnalyzer(analyzer *Analyzer) {
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
			analyzer.processMatchStart()
		}
	})

	parser.RegisterEventHandler(analyzer.defaultRoundFreezetimeChangedHandler)

	parser.RegisterEventHandler(analyzer.defaultRoundStartHandler)

	parser.RegisterEventHandler(analyzer.defaultRoundEndOfficiallyHandler)

	parser.RegisterEventHandler(func(event events.AnnouncementWinPanelMatch) {
		analyzer.updatePlayersScores()
		matchStarted = false
	})
}
