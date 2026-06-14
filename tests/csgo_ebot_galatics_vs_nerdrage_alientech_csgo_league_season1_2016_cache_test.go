package tests

import (
	"testing"

	"github.com/akiver/cs-demo-analyzer/pkg/api"
	"github.com/akiver/cs-demo-analyzer/pkg/api/constants"
	"github.com/akiver/cs-demo-analyzer/tests/assertion"
	"github.com/akiver/cs-demo-analyzer/tests/testsutils"
)

// Knife round with teams switch.
// The game is stopped a few seconds after the beginning of the 1st round.
// https://www.hltv.org/stats/matches/mapstatsid/28613/nerdrage-vs-galatics
func TestEbot_Galatics_VS_Nerdrage_AlienTech_CSGO_League_Season1_2016_Cache(t *testing.T) {
	t.Parallel()

	demoName := "ebot_galatics_vs_nerdrage_alientech_csgo_league_season1_2016_cache"
	demoPath := testsutils.GetDemoPath("csgo", demoName)
	match, err := api.AnalyzeDemo(demoPath, api.AnalyzeDemoOptions{
		Source: constants.DemoSourceEbot,
	})

	if err != nil {
		t.Fatal(err)
	}

	assertion.AssertMatchSnapshot(t, match, demoName)
}
