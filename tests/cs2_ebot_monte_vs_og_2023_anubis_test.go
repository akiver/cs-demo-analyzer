package tests

import (
	"testing"

	"github.com/akiver/cs-demo-analyzer/pkg/api"
	"github.com/akiver/cs-demo-analyzer/pkg/api/constants"
	"github.com/akiver/cs-demo-analyzer/tests/assertion"
	"github.com/akiver/cs-demo-analyzer/tests/testsutils"
)

// https://www.hltv.org/matches/2367504/monte-vs-og-roobet-cup-2023
func TestEbot_Monte_VS_OG_Roobet_Cup_2023_Anubis(t *testing.T) {
	t.Parallel()

	demoName := "ebot_monte_vs_og_roobet_cup_2023_anubis"
	demoPath := testsutils.GetDemoPath("cs2", demoName)
	match, err := api.AnalyzeDemo(demoPath, api.AnalyzeDemoOptions{
		Source: constants.DemoSourceEbot,
	})

	if err != nil {
		t.Fatal(err)
	}

	assertion.AssertMatchSnapshot(t, match, demoName)
}
