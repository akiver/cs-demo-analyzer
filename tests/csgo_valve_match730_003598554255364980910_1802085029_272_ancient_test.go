package tests

import (
	"testing"

	"github.com/akiver/cs-demo-analyzer/pkg/api"
	"github.com/akiver/cs-demo-analyzer/pkg/api/constants"
	"github.com/akiver/cs-demo-analyzer/tests/testsutils"
)

// Valve short match (MR 8) demo.
func TestValve_Match730_003598554255364980910_1802085029_272_Ancient(t *testing.T) {
	demoName := "valve_match730_003598554255364980910_1802085029_272_ancient"
	demoPath := testsutils.GetDemoPath("csgo", demoName)
	match, err := api.AnalyzeDemo(demoPath, api.AnalyzeDemoOptions{
		Source: constants.DemoSourceValve,
	})
	if err != nil {
		t.Error(err)
	}
	expectedRoundCount := 14
	expectedPlayerCount := 10
	expectedTeamAScore := 5
	expectedTeamBScore := 9
	expectedScoreFirstHalfTeamA := 2
	expectedScoreFirstHalfTeamB := 6
	expectedScoreSecondHalfTeamA := 3
	expectedScoreSecondHalfTeamB := 3
	expectedTeamNameA := "Team A"
	expectedTeamNameB := "Team B"
	expectedWinnerName := expectedTeamNameB
	expectedMaxRounds := 16

	if match.TeamA.Name != expectedTeamNameA {
		t.Errorf("expected team name A to be %s got %s", expectedTeamNameA, match.TeamA.Name)
	}
	if match.TeamB.Name != expectedTeamNameB {
		t.Errorf("expected team name B to be %s got %s", expectedTeamNameB, match.TeamB.Name)
	}
	if len(match.Rounds) != expectedRoundCount {
		t.Errorf("expected %d rounds but got %d", expectedRoundCount, len(match.Rounds))
	}
	if len(match.Players()) != expectedPlayerCount {
		t.Errorf("expected %d players but got %d", expectedPlayerCount, len(match.Players()))
	}
	if match.TeamA.Score != expectedTeamAScore {
		t.Errorf("expected teamScTeamBScore A to be %d got %d", expectedTeamAScore, match.TeamA.Score)
	}
	if match.TeamB.Score != expectedTeamBScore {
		t.Errorf("expected score team B to be %d got %d", expectedTeamBScore, match.TeamB.Score)
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
	if match.Winner.Name != expectedWinnerName {
		t.Errorf("expected winner to be %s but got %s", expectedWinnerName, match.Winner.Name)
	}
	if match.MaxRounds != expectedMaxRounds {
		t.Errorf("expected max rounds to be %d but got %d", expectedMaxRounds, match.MaxRounds)
	}
}
