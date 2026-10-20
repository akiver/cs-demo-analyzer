package cs2

import (
	"github.com/akiver/cs-demo-analyzer/internal/modelaccess"
	"github.com/akiver/cs-demo-analyzer/pkg/api/constants"
	events "github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/events"
)

func createEsportligaenAnalyzer(analyzer *Analyzer) {
	parser := analyzer.parser
	match := analyzer.match
	modelaccess.SetGameModeStr(match, constants.GameModeStrCompetitive)
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
