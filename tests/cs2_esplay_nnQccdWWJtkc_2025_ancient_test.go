package tests

import (
	"testing"

	"github.com/akiver/cs-demo-analyzer/pkg/api"
	"github.com/akiver/cs-demo-analyzer/pkg/api/constants"
	"github.com/akiver/cs-demo-analyzer/tests/assertion"
	"github.com/akiver/cs-demo-analyzer/tests/testsutils"
)

// 2v2
// https://esplay.com/m/nnQccdWWJtkc/midnight-vs-haha123xd
func TestEsplay_nnQccdWWJtkc_2025_Vertigo(t *testing.T) {
	t.Parallel()

	demoName := "esplay_nnQccdWWJtkc_2025_vertigo"
	demoPath := testsutils.GetDemoPath("cs2", demoName)
	match, err := api.AnalyzeDemo(demoPath, api.AnalyzeDemoOptions{
		Source: constants.DemoSourceEsplay,
	})

	if err != nil {
		t.Fatal(err)
	}

	assertion.AssertMatchSnapshot(t, match, demoName)
}
