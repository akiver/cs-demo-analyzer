package tests

import (
	"testing"

	"github.com/akiver/cs-demo-analyzer/pkg/api"
	"github.com/akiver/cs-demo-analyzer/pkg/api/constants"
	"github.com/akiver/cs-demo-analyzer/tests/assertion"
	"github.com/akiver/cs-demo-analyzer/tests/testsutils"
)

// Demo with 1 overtime.
// Recording start after the end of the first round's freezetime.
// Pauses at the end of the round 6, 8, 13, 15, 24, 26, 29.
// After a few seconds at the end of the round 20, the round is cancelled (!stop) and really starts at tick 330201.
// https://www.hltv.org/stats/matches/mapstatsid/42335/immortals-vs-north
func TestEbot_Immortals_VS_North_IEM_Katowice_2017_Overpass(t *testing.T) {
	t.Parallel()

	demoName := "ebot_immortals_vs_north_iem_katowice_2017_overpass"
	demoPath := testsutils.GetDemoPath("csgo", demoName)
	match, err := api.AnalyzeDemo(demoPath, api.AnalyzeDemoOptions{
		Source: constants.DemoSourceEbot,
	})

	if err != nil {
		t.Fatal(err)
	}

	assertion.AssertMatchSnapshot(t, match, demoName)
}
