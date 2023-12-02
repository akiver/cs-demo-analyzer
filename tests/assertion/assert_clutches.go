package assertion

import (
	"testing"

	"github.com/akiver/cs-demo-analyzer/pkg/api"
)

func AssertClutches(t *testing.T, match *api.Match, clutches []api.Clutch) {
	if len(match.Clutches) != len(clutches) {
		t.Errorf("Expected %d clutches, but got %d", len(clutches), len(match.Clutches))
	}

	for index, expectedClutch := range clutches {
		clutch := match.Clutches[index]
		if clutch.RoundNumber != expectedClutch.RoundNumber {
			t.Errorf("Expected clutch round number to be %d, got %d", expectedClutch.RoundNumber, clutch.RoundNumber)
		}
		if clutch.ClutcherSteamID64 != expectedClutch.ClutcherSteamID64 {
			t.Errorf("Expected clutcher steam id to be %d, got %d for round %d", expectedClutch.ClutcherSteamID64, clutch.ClutcherSteamID64, expectedClutch.RoundNumber)
		}
		if clutch.ClutcherName != expectedClutch.ClutcherName {
			t.Errorf("Expected clutcher name to be %s, got %s for round %d", expectedClutch.ClutcherName, clutch.ClutcherName, expectedClutch.RoundNumber)
		}
		if clutch.HasWon != expectedClutch.HasWon {
			t.Errorf("Expected clutch has won to be %t, got %t for round %d", expectedClutch.HasWon, clutch.HasWon, expectedClutch.RoundNumber)
		}
		if clutch.ClutcherSurvived != expectedClutch.ClutcherSurvived {
			t.Errorf("Expected clutcher survived to be %t, got %t for round %d", expectedClutch.ClutcherSurvived, clutch.ClutcherSurvived, expectedClutch.RoundNumber)
		}
		if clutch.ClutcherKillCount != expectedClutch.ClutcherKillCount {
			t.Errorf("Expected clutcher kill count to be %d, got %d for round %d", expectedClutch.ClutcherKillCount, clutch.ClutcherKillCount, expectedClutch.RoundNumber)
		}
	}
}
