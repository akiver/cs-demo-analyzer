package tests

import (
	"testing"

	"github.com/akiver/cs-demo-analyzer/pkg/api"
	"github.com/akiver/cs-demo-analyzer/pkg/api/constants"
	"github.com/akiver/cs-demo-analyzer/tests/assertion"
	"github.com/akiver/cs-demo-analyzer/tests/testsutils"
)

// https://www.hltv.org/stats/matches/mapstatsid/179910/aurora-vs-3dmax
// - Contains a round backup restore at round 15
// - Teams stay after knife round
func Test_MatchZy_Aurora_vs_3dmax_m3_2024_Anubis(t *testing.T) {
	demoName := "matchzy_aurora_vs_3dmax_m3_anubis"
	demoPath := testsutils.GetDemoPath("cs2", demoName)
	match, err := api.AnalyzeDemo(demoPath, api.AnalyzeDemoOptions{
		Source: constants.DemoSourceMatchZy,
	})

	if err != nil {
		t.Fatal(err)
	}

	assertion.AssertMatchSnapshot(t, match, demoName)
}
