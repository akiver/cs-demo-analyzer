package api

import (
	"fmt"
	"strings"

	s "github.com/akiver/cs-demo-analyzer/internal/strings"
	"github.com/akiver/cs-demo-analyzer/pkg/api/constants"
	events "github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/events"
	sendtables "github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/sendtables"
)

func createEbotAnalyzer(analyzer *Analyzer) {
	match := analyzer.match
	match.gameModeStr = constants.GameModeStrCompetitive
	parser := analyzer.parser
	matchStarted := true
	matchStartDetected := false
	ctWantStop := false
	tWantStop := false
	lastMatchRoundOfficiallyEndTick := -1 // Used to detect stop command requests
	lastRoundEndTick := -1                // Used to detect backup restore
	playersTeamChangeTick := -1           // Used to detect swap team after a possible knife round

	analyzer.matchStarted = func() bool {
		return matchStarted
	}

	resetCurrentRound := func() {
		currentRound := analyzer.currentRound
		currentRound.StartFrame = analyzer.parser.CurrentFrame()
		currentRound.StartTick = analyzer.currentTick()
		currentRound.FreezeTimeEndFrame = -1
		currentRound.FreezeTimeEndTick = -1
	}

	parser.RegisterEventHandler(func(event events.SayText) {
		if analyzer.currentRound.Number > 1 {
			return
		}

		text := s.RemoveInvisibleChars(event.Text)
		if !tWantStop {
			tWantStop = strings.Contains(text, "(T) want to stop")
		}
		if !ctWantStop {
			ctWantStop = strings.Contains(text, "(CT) want to stop")
		}

		isMatchStopped := ctWantStop && tWantStop
		if isMatchStopped {
			ctWantStop = false
			tWantStop = false
			matchStarted = false
			matchStartDetected = false
			resetCurrentRound()
			analyzer.reset()
		}
	})

	parser.RegisterEventHandler(func(events.DataTablesParsed) {
		parser.ServerClasses().FindByName("CCSTeam").OnEntityCreated(func(entity sendtables.Entity) {
			// Detects possible teams switch after the end of a knife round (!switch command).
			// If the teams are not switched (!stay command) the following players prop of the team entity doesn't change.
			// Teams data such as its name are updated just before this prop change.
			var playersProps []sendtables.Property
			if analyzer.isSource2 {
				// The array is split into 4 props in Source 2
				for i := 0; i < 4; i++ {
					iStr := fmt.Sprintf("%04d", i)
					playersProps = append(playersProps, entity.Property("m_aPlayers."+iStr))
				}
			} else {
				playersProps = append(playersProps, entity.Property("\"player_array\""))
			}

			for _, prop := range playersProps {
				prop.OnUpdate(func(value sendtables.PropertyValue) {
					if value.Any == nil || analyzer.currentTick() <= 1 || analyzer.currentRound.Number > 1 {
						return
					}

					if playersTeamChangeTick == analyzer.currentTick() {
						analyzer.reset()
					}
					playersTeamChangeTick = analyzer.currentTick()
				})
			}
		})

		parser.RegisterEventHandler(func(event events.MatchStartedChanged) {
			analyzer.registerUnknownPlayers()
			matchStarted = event.NewIsStarted
			if matchStarted && !matchStartDetected {
				currentRound := analyzer.currentRound
				currentRound.StartFrame = parser.CurrentFrame()
				currentRound.StartTick = analyzer.currentTick()
				analyzer.updateTeamNames()
				analyzer.createPlayersEconomies()
				matchStartDetected = true
			} else if len(match.Rounds) == 0 {
				matchStartDetected = false
			}

			// When a stop command is requested by a player after the first round there is a round officially end event between game restart and match start:
			// Player says !stop -> m_bGameRestart -> m_bHasMatchStarted -> round_end_officially -> round_start -> m_bHasMatchStarted
			// When it's a restart for a real match live, there is no round official end event:
			// Going live -> m_bGameRestart -> m_bHasMatchStarted -> round_start -> m_bHasMatchStarted
			// Stop commands during the first round don't trigger any game restart. We rely on chat messages for this case.

			if matchStarted && analyzer.currentRound.Number > 1 {
				isMatchStopped := !analyzer.secondsHasPassedSinceTick(5, lastMatchRoundOfficiallyEndTick)
				if isMatchStopped {
					matchStarted = false
				}
			}
		})

		parser.RegisterEventHandler(func(event events.RoundFreezetimeChanged) {
			if !analyzer.matchStarted() {
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
		isBackupRestoration := analyzer.currentTick() == lastRoundEndTick
		if isBackupRestoration {
			matchStarted = true
			resetCurrentRound()
			return
		}

		if !analyzer.matchStarted() || analyzer.currentTick() == 0 || len(match.Rounds) == 0 {
			return
		}

		analyzer.createRound()
	})

	parser.RegisterEventHandler(func(event events.RoundEnd) {
		lastRoundEndTick = analyzer.currentTick()

		knifeKillCount := 0
		killCount := 0
		for _, kill := range analyzer.match.Kills {
			if kill.RoundNumber != analyzer.currentRound.Number {
				continue
			}
			if kill.WeaponName == constants.WeaponKnife {
				knifeKillCount++
			}
			killCount++
		}

		isKnifeRound := killCount > 0 && killCount == knifeKillCount
		if isKnifeRound {
			analyzer.reset()
			matchStarted = false
			matchStartDetected = false
		}
	})

	parser.RegisterEventHandler(func(event events.RoundEndOfficial) {
		lastMatchRoundOfficiallyEndTick = analyzer.currentTick()
		if !analyzer.matchStarted() {
			return
		}

		isBackupRestoration := analyzer.currentTick() == lastRoundEndTick
		if isBackupRestoration {
			return
		}

		match.Rounds = append(match.Rounds, analyzer.currentRound)
	})

	parser.RegisterEventHandler(func(event events.AnnouncementWinPanelMatch) {
		analyzer.updatePlayersScores()
	})
}
