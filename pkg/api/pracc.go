package api

import (
	"github.com/akiver/cs-demo-analyzer/pkg/api/constants"
)

func createPraccAnalyzer(analyzer *Analyzer) {
	match := analyzer.match
	match.gameModeStr = constants.GameModeStrCompetitive
	parser := analyzer.parser
	matchStarted := true // PRACC servers start demo recording only when the match really starts

	analyzer.matchStarted = func() bool {
		return matchStarted
	}

	parser.RegisterEventHandler(analyzer.defaultRoundStartHandler)
	parser.RegisterEventHandler(analyzer.defaultRoundFreezetimeChangedHandler)
	parser.RegisterEventHandler(analyzer.defaultRoundEndOfficiallyHandler)
	parser.RegisterEventHandler(analyzer.defaultAnnouncementWinPanelMatchHandler)
}
