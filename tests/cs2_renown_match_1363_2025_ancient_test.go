package tests

import (
	"testing"

	"github.com/akiver/cs-demo-analyzer/pkg/api"
	"github.com/akiver/cs-demo-analyzer/pkg/api/constants"
	"github.com/akiver/cs-demo-analyzer/tests/assertion"
	"github.com/akiver/cs-demo-analyzer/tests/testsutils"
)

// Match with 1 overtime.
// The server freezes during round 22, players are disconnected but the round still ends by bomb explosion.
// The match is paused and resumed when all players are reconnected.
// https://renown.gg/match/1363
func TestRenown_Match_1363_2025_Ancient(t *testing.T) {
	demoName := "renown_match_1363_2025_ancient"
	demoPath := testsutils.GetDemoPath("cs2", demoName)
	match, err := api.AnalyzeDemo(demoPath, api.AnalyzeDemoOptions{
		Source: constants.DemoSourceRenown,
	})

	if err != nil {
		t.Fatal(err)
	}

	assertion.AssertMatchSnapshot(t, match, demoName)
}
