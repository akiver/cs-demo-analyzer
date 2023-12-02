package tests

import (
	"testing"

	"github.com/akiver/cs-demo-analyzer/pkg/api"
	"github.com/akiver/cs-demo-analyzer/pkg/api/constants"
	"github.com/akiver/cs-demo-analyzer/tests/assertion"
	"github.com/akiver/cs-demo-analyzer/tests/fake"
	"github.com/akiver/cs-demo-analyzer/tests/testsutils"
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/common"
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/events"
)

// Knife round with teams switch.
// The game is stopped a few seconds after the beginning of the 1st round.
// https://www.hltv.org/stats/matches/mapstatsid/28613/nerdrage-vs-galatics
func TestEbot_Galatics_VS_Nerdrage_AlienTech_CSGO_League_Season1_2016_Cache(t *testing.T) {
	demoName := "ebot_galatics_vs_nerdrage_alientech_csgo_league_season1_2016_cache"
	demoPath := testsutils.GetDemoPath("csgo", demoName)
	match, err := api.AnalyzeDemo(demoPath, api.AnalyzeDemoOptions{
		Source: constants.DemoSourceEbot,
	})
	if err != nil {
		t.Error(err)
	}

	expectedRoundCount := 20
	expectedPlayerCount := 10
	expectedScoreTeamA := 16
	expectedScoreTeamB := 4
	expectedScoreFirstHalfTeamA := 13
	expectedScoreFirstHalfTeamB := 2
	expectedScoreSecondHalfTeamA := 3
	expectedScoreSecondHalfTeamB := 2
	expectedTeamNameA := "NerdRage"
	expectedTeamNameB := "GALATICS.CSGO"
	expectedWinnerName := expectedTeamNameA

	if len(match.Rounds) != expectedRoundCount {
		t.Errorf("expected %d rounds but got %d", expectedRoundCount, len(match.Rounds))
	}
	if len(match.Players()) != expectedPlayerCount {
		t.Errorf("expected %d players but got %d", expectedPlayerCount, len(match.Players()))
	}
	if match.TeamA.Score != expectedScoreTeamA {
		t.Errorf("expected score team A to be %d got %d", expectedScoreTeamA, match.TeamA.Score)
	}
	if match.TeamB.Score != expectedScoreTeamB {
		t.Errorf("expected score team B to be %d got %d", expectedScoreTeamB, match.TeamB.Score)
	}
	if match.TeamA.Name != expectedTeamNameA {
		t.Errorf("expected team name A to be %s got %s", expectedTeamNameA, match.TeamA.Name)
	}
	if match.TeamB.Name != expectedTeamNameB {
		t.Errorf("expected team name B to be %s got %s", expectedTeamNameB, match.TeamB.Name)
	}
	if match.TeamA.ScoreFirstHalf != expectedScoreFirstHalfTeamA {
		t.Errorf("expected score first half team A to be %d got %d", expectedScoreFirstHalfTeamA, match.TeamA.ScoreFirstHalf)
	}
	if match.TeamB.ScoreFirstHalf != expectedScoreFirstHalfTeamB {
		t.Errorf("expected score first half team B to be %d got %d", expectedScoreFirstHalfTeamB, match.TeamB.ScoreFirstHalf)
	}
	if match.TeamA.ScoreSecondHalf != expectedScoreSecondHalfTeamA {
		t.Errorf("expected score second half team A to be %d got %d", expectedScoreSecondHalfTeamA, match.TeamA.ScoreSecondHalf)
	}
	if match.TeamB.ScoreSecondHalf != expectedScoreSecondHalfTeamB {
		t.Errorf("expected score second half team B to be %d got %d", expectedScoreSecondHalfTeamB, match.TeamB.ScoreSecondHalf)
	}
	if match.Winner != nil && match.Winner.Name != expectedWinnerName {
		t.Errorf("expected winner to be %s but got %s", expectedWinnerName, match.Winner.Name)
	}

	var rounds = []fake.FakeRound{
		{
			Number:            1,
			StartTick:         29831,
			StartFrame:        3734,
			EndTick:           40145,
			FreezeTimeEndTick: 31367,
			EndOfficiallyTick: 40789,
			TeamAStartMoney:   4000,
			TeamBStartMoney:   4000,
			TeamAEconomyType:  constants.EconomyTypePistol,
			TeamBEconomyType:  constants.EconomyTypePistol,
			EndReason:         events.RoundEndReasonTargetBombed,
			WinnerSide:        common.TeamTerrorists,
			TeamAScore:        0,
			TeamBScore:        1,
			WinnerName:        expectedTeamNameB,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
		},
		{
			Number:            2,
			StartTick:         40789,
			StartFrame:        5103,
			EndTick:           47629,
			FreezeTimeEndTick: 42325,
			EndOfficiallyTick: 48269,
			TeamAStartMoney:   9100,
			TeamBStartMoney:   15750,
			TeamAEconomyType:  constants.EconomyTypeForceBuy,
			TeamBEconomyType:  constants.EconomyTypeSemi,
			EndReason:         events.RoundEndReasonCTWin,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        1,
			TeamBScore:        1,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
		},
	}

	assertion.AssertRounds(t, match, rounds)
}
