package tests

import (
	"testing"

	"github.com/akiver/cs-demo-analyzer/pkg/api"
	"github.com/akiver/cs-demo-analyzer/pkg/api/constants"
	"github.com/akiver/cs-demo-analyzer/tests/assertion"
	"github.com/akiver/cs-demo-analyzer/tests/testsutils"
)

// ESEA demo with warmup, half time break, coaches and pauses at freeze time.
// https://www.hltv.org/stats/matches/mapstatsid/46797/clg-vs-liquid
func TestEsea_CLG_VS_Liquid_IEM_Cologne_2017_Cbble(t *testing.T) {
	demoName := "esea_clg_vs_liquid_iem_cologne_2017_cbble"
	demoPath := testsutils.GetDemoPath("csgo", demoName)
	match, err := api.AnalyzeDemo(demoPath, api.AnalyzeDemoOptions{
		Source: constants.DemoSourceESEA,
	})

	if err != nil {
		t.Fatal(err)
	}

	assertion.AssertMatchSnapshot(t, match, demoName)
}
