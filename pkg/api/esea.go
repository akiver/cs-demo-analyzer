package api

import (
	"fmt"

	"github.com/akiver/cs-demo-analyzer/pkg/api/constants"
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/common"
	events "github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/events"
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/sendtables"
)

const (
	maxPlayers                  = 64
	counterTerroristsTeamNumber = byte(common.TeamCounterTerrorists)
	terroristsTeamNumber        = byte(common.TeamTerrorists)
)

func createEseaAnalyzer(analyzer *Analyzer) {
	parser := analyzer.parser
	match := analyzer.match
	match.gameModeStr = constants.GameModeStrCompetitive
	matchStarted := false
	playerSwapTickDetected := -1      // Keep track of the last tick detected that contains a player Team swap, used to detect teams side switch.
	playerSwappedAtTickCount := 0     // How many players have been swapped at the same tick, used to detect teams side switch.
	consecutiveMatchRestartCount := 0 // ESEA uses the old school "LO3"
	lastMatchStartTick := 0           // Keep track of the last match restart tick to detect LO3
	shouldSwapTeams := false

	analyzer.matchStarted = func() bool {
		return matchStarted
	}

	// Players scores are reset to 0 at the end of the match before the last round really ended.
	// Listen and update PlayersBySteamID scores to have correct values at the end of the match.
	onScoreUpdate := func(value sendtables.PropertyValue) {
		if !matchStarted {
			return
		}

		analyzer.updatePlayersScores()
	}

	parser.RegisterEventHandler(func(event events.MatchStart) {
		analyzer.registerUnknownPlayers()
		currentTick := analyzer.currentTick()

		if analyzer.secondsHasPassedSinceTick(5, lastMatchStartTick) {
			lastMatchStartTick = currentTick
			consecutiveMatchRestartCount = 1
			matchStarted = false
			return
		}

		consecutiveMatchRestartCount += 1
		if consecutiveMatchRestartCount == 3 {
			matchStarted = true
			consecutiveMatchRestartCount = 0
			currentRound := analyzer.currentRound
			currentRound.StartTick = analyzer.currentTick()
			currentRound.StartFrame = parser.CurrentFrame()

			// Some players may have joined the wrong team when entering the server, in that case they are silently
			// switched to the correct team right before the LO3.
			// Update the team of those players only for the first LO3.
			if currentRound.Number == 1 {
				analyzer.updatePlayersCurrentTeam()
				analyzer.updateTeamNames()
			}

			if shouldSwapTeams {
				match.swapTeams()
				analyzer.currentRound.TeamASide = *match.TeamA.CurrentSide
				analyzer.currentRound.TeamBSide = *match.TeamB.CurrentSide
				shouldSwapTeams = false
			}
		}

		lastMatchStartTick = currentTick
	})

	parser.RegisterEventHandler(func(event events.RoundStart) {
		if !analyzer.matchStarted() {
			return
		}

		match.Rounds = append(match.Rounds, analyzer.currentRound)
		analyzer.createRound()
	})

	parser.RegisterEventHandler(func(event events.RoundFreezetimeChanged) {
		freezetimeEnded := !event.NewIsFreezetime
		if freezetimeEnded {
			analyzer.currentRound.FreezeTimeEndTick = analyzer.currentTick()
			analyzer.currentRound.FreezeTimeEndFrame = parser.CurrentFrame()
			analyzer.lastFreezeTimeEndTick = analyzer.currentTick()
		}
	})

	parser.RegisterEventHandler(func(events.DataTablesParsed) {
		parser.ServerClasses().FindByName("CCSGameRulesProxy").OnEntityCreated(func(entity sendtables.Entity) {
			gameRestartProp := entity.Property("cs_gamerules_data.m_bGameRestart")
			if gameRestartProp == nil {
				gameRestartProp = entity.Property("m_pGameRules.m_bGameRestart")
			}
			gameRestartProp.OnUpdate(func(value sendtables.PropertyValue) {
				isRestarting := value.BoolVal()
				if isRestarting && matchStarted {
					matchStarted = false
					if len(match.Rounds) > 1 {
						match.Rounds = append(match.Rounds, analyzer.currentRound)
						analyzer.createRound()
					}
				}

				analyzer.isFirstRoundOfHalf = true
			})
		})

		var playerClass sendtables.ServerClass
		if analyzer.isSource2 {
			playerClass = parser.ServerClasses().FindByName("CCSPlayerController")
		} else {
			playerClass = parser.ServerClasses().FindByName("CCSPlayer")
		}
		playerClass.OnEntityCreated(func(entity sendtables.Entity) {
			var moneyStartProp sendtables.Property
			if analyzer.isSource2 {
				moneyStartProp = entity.Property("m_pInGameMoneyServices.m_iStartAccount")
			} else {
				moneyStartProp = entity.Property("m_iStartAccount")
			}
			moneyStartProp.OnUpdate(func(value sendtables.PropertyValue) {
				analyzer.createPlayersEconomies()
			})

			entity.Property("m_iCoachingTeam").OnUpdate(func(value sendtables.PropertyValue) {
				teamNumber := common.Team(value.Int())
				if teamNumber != common.TeamCounterTerrorists && teamNumber != common.TeamTerrorists {
					return
				}

				// Remove coaches from players
				for _, player := range parser.GameState().Participants().All() {
					if player.EntityID == entity.ID() {
						if analyzer.match.PlayersBySteamID[player.SteamID64] != nil {
							delete(analyzer.match.PlayersBySteamID, player.SteamID64)
						}
						break
					}
				}
			})

			// Teams swap detection, it occurs when all players are switched to the opposite Team at the same tick.
			// The team swapping is delayed at the next LO3 because some events related to team sides may happen between
			// the swap detection and the end of the current round.
			entity.Property("m_iTeamNum").OnUpdate(func(value sendtables.PropertyValue) {
				if len(match.Rounds) < 1 {
					return
				}

				if analyzer.currentTick() == playerSwapTickDetected {
					teamNumber := byte(value.Int())
					if teamNumber == counterTerroristsTeamNumber || teamNumber == terroristsTeamNumber {
						playerSwappedAtTickCount++
						validPlayerCount := len(parser.GameState().Participants().Playing())
						if playerSwappedAtTickCount == validPlayerCount {
							playerSwapTickDetected = 0
							shouldSwapTeams = true
						}
					}
				} else {
					playerSwappedAtTickCount = 1
				}

				playerSwapTickDetected = analyzer.currentTick()
			})

			if analyzer.isSource2 {
				entity.Property("m_iScore").OnUpdate(onScoreUpdate)
			}
		})

		if !analyzer.isSource2 {
			parser.ServerClasses().FindByName("CCSPlayerResource").OnEntityCreated(func(entity sendtables.Entity) {
				for i := 0; i < maxPlayers; i++ {
					iAsString := fmt.Sprintf("%03d", i)
					scoreProp := entity.Property("m_iScore." + iAsString)
					if scoreProp != nil {
						scoreProp.OnUpdate(onScoreUpdate)
					}
				}
			})
		}
	})
}
