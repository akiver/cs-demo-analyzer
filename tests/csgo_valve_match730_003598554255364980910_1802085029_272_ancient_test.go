package tests

import (
	"testing"

	"github.com/akiver/cs-demo-analyzer/pkg/api"
	"github.com/akiver/cs-demo-analyzer/pkg/api/constants"
	"github.com/akiver/cs-demo-analyzer/tests/assertion"
	"github.com/akiver/cs-demo-analyzer/tests/testsutils"
)

// Valve short match (MR 8) demo.
func TestValve_Match730_003598554255364980910_1802085029_272_Ancient(t *testing.T) {
	t.Parallel()

	demoName := "valve_match730_003598554255364980910_1802085029_272_ancient"
	demoPath := testsutils.GetDemoPath("csgo", demoName)
	match, err := api.AnalyzeDemo(demoPath, api.AnalyzeDemoOptions{
		Source: constants.DemoSourceValve,
	})

	if err != nil {
		t.Fatal(err)
	}

	assertion.AssertMatchSnapshot(t, match, demoName)
}
