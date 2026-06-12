package tests

import (
	"testing"

	"github.com/akiver/cs-demo-analyzer/pkg/api"
	"github.com/akiver/cs-demo-analyzer/pkg/api/constants"
	"github.com/akiver/cs-demo-analyzer/tests/assertion"
	"github.com/akiver/cs-demo-analyzer/tests/testsutils"
)

// E-frag starts as T, won the knife round and switch to CT.
// There is then a warmup and the first round is cancelled during freezetime (!stop command).
// The first round really starts after a restart by the eBot.
// At the beginning of the second round, a pause is requested. The game is paused and the 2nd round starts at tick 74004.
// https://www.hltv.org/stats/matches/mapstatsid/28406/e-fragnet-vs-faze
func TestEbot_Efrag_Net_VS_Faze_IEM_Oakland_2016_Cache(t *testing.T) {
	demoName := "ebot_efrag_net_vs_faze_iem_oakland_2016_cache"
	demoPath := testsutils.GetDemoPath("csgo", demoName)
	match, err := api.AnalyzeDemo(demoPath, api.AnalyzeDemoOptions{
		Source: constants.DemoSourceEbot,
	})

	if err != nil {
		t.Fatal(err)
	}

	assertion.AssertMatchSnapshot(t, match, demoName)
}
