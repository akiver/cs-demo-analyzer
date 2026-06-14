package tests

import (
	"testing"

	"github.com/akiver/cs-demo-analyzer/pkg/api"
	"github.com/akiver/cs-demo-analyzer/pkg/api/constants"
	"github.com/akiver/cs-demo-analyzer/tests/assertion"
	"github.com/akiver/cs-demo-analyzer/tests/testsutils"
)

// 5v5 with surrender.
// https://esplay.com/m/nvBBvqNCfFHV/team-pytonorm-vs-team-shawty
func TestEsplay_nvBBvqNCfFHV_2025_Train(t *testing.T) {
	t.Parallel()

	demoName := "esplay_nvBBvqNCfFHV_2025_train"
	demoPath := testsutils.GetDemoPath("cs2", demoName)
	match, err := api.AnalyzeDemo(demoPath, api.AnalyzeDemoOptions{
		Source: constants.DemoSourceEsplay,
	})

	if err != nil {
		t.Fatal(err)
	}

	assertion.AssertMatchSnapshot(t, match, demoName)
}
