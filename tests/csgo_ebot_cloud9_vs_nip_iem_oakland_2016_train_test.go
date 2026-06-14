package tests

import (
	"testing"

	"github.com/akiver/cs-demo-analyzer/pkg/api"
	"github.com/akiver/cs-demo-analyzer/pkg/api/constants"
	"github.com/akiver/cs-demo-analyzer/tests/assertion"
	"github.com/akiver/cs-demo-analyzer/tests/testsutils"
)

// Pause requested at the end of the round 3 and effective at the beginning of the round 4 (tick 32229).
// The eBot backup is restored at tick 43560, that's when the round 4 really starts.
// The round 6 is cancelled few seconds after it started, a backup is restored and the round 6 really start.
// https://www.hltv.org/stats/matches/mapstatsid/38754/cloud9-vs-ninjas-in-pyjamas
func TestEbot_Cloud9_VS_NIP_IEM_Oakland_2016_Train(t *testing.T) {
	t.Parallel()

	demoName := "ebot_cloud9_vs_nip_iem_oakland_2016_train"
	demoPath := testsutils.GetDemoPath("csgo", demoName)
	match, err := api.AnalyzeDemo(demoPath, api.AnalyzeDemoOptions{
		Source: constants.DemoSourceEbot,
	})

	if err != nil {
		t.Fatal(err)
	}

	assertion.AssertMatchSnapshot(t, match, demoName)
}
