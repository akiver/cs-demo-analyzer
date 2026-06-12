package tests

import (
	"testing"

	"github.com/akiver/cs-demo-analyzer/pkg/api"
	"github.com/akiver/cs-demo-analyzer/pkg/api/constants"
	"github.com/akiver/cs-demo-analyzer/tests/assertion"
	"github.com/akiver/cs-demo-analyzer/tests/testsutils"
)

// https://www.hltv.org/stats/matches/mapstatsid/151559/astralis-vs-saw
func TestChallengerModeSawVsAstralisM2Paris2023CqInferno(t *testing.T) {
	demoName := "challengermode_saw_vs_astralis_m3_Paris_2023_CQ_inferno"
	demoPath := testsutils.GetDemoPath("csgo", demoName)
	match, err := api.AnalyzeDemo(demoPath, api.AnalyzeDemoOptions{
		Source: constants.DemoSourceFaceIt,
	})

	if err != nil {
		t.Fatal(err)
	}

	assertion.AssertMatchSnapshot(t, match, demoName)
}
