package api

import (
	"github.com/akiver/cs-demo-analyzer/pkg/api/constants"
	events "github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/events"
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/msgs2"
)

func createValveAnalyzer(analyzer *Analyzer) {
	parser := analyzer.parser
	match := analyzer.match
	analyzer.matchStarted = parser.GameState().IsMatchStarted

	parser.RegisterNetMessageHandler(func(srvInfo *msgs2.CSVCMsg_ServerInfo) {
		match.gameModeStr = constants.GameModeStr(srvInfo.GameSessionConfig.GetGamemode())
	})

	parser.RegisterEventHandler(func(event events.MatchStart) {
		currentRound := analyzer.currentRound
		currentRound.StartFrame = parser.CurrentFrame()
		currentRound.StartTick = analyzer.currentTick()
		analyzer.updateTeamNames()
	})

	parser.RegisterEventHandler(func(event events.MatchStartedChanged) {
		isMatchStarted := !event.OldIsStarted && event.NewIsStarted
		if isMatchStarted {
			currentRound := analyzer.currentRound
			currentRound.TeamASide = *match.TeamA.CurrentSide
			currentRound.TeamBSide = *match.TeamB.CurrentSide
		}
	})

	parser.RegisterEventHandler(func(event events.GameHalfEnded) {
		analyzer.isFirstRoundOfHalf = true
	})

	parser.RegisterEventHandler(analyzer.defaultRoundFreezetimeChangedHandler)

	parser.RegisterEventHandler(analyzer.defaultRoundStartHandler)

	parser.RegisterEventHandler(analyzer.defaultRoundEndOfficiallyHandler)

	parser.RegisterEventHandler(analyzer.defaultAnnouncementWinPanelMatchHandler)
}
