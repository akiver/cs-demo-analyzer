package tests

import (
	"testing"

	"github.com/akiver/cs-demo-analyzer/pkg/api"
	"github.com/akiver/cs-demo-analyzer/pkg/api/constants"
	"github.com/akiver/cs-demo-analyzer/tests/assertion"
	"github.com/akiver/cs-demo-analyzer/tests/testsutils"
)

// Knife round at the beginning and teams are swapped.
// https://www.hltv.org/matches/2317273/optic-vs-faze-iem-oakland-2017
func TestEbot_Optic_VS_Faze_IEM_Oakland_2017_Overpass(t *testing.T) {
	t.Parallel()

	demoName := "ebot_optic_vs_faze_iem_oakland_2017_overpass"
	demoPath := testsutils.GetDemoPath("csgo", demoName)
	match, err := api.AnalyzeDemo(demoPath, api.AnalyzeDemoOptions{
		Source: constants.DemoSourceEbot,
	})

	if err != nil {
		t.Fatal(err)
	}

	assertion.AssertMatchSnapshot(t, match, demoName)
}
