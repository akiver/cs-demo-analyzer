package tests

import (
	"testing"

	"github.com/akiver/cs-demo-analyzer/pkg/api"
	"github.com/akiver/cs-demo-analyzer/pkg/api/constants"
	"github.com/akiver/cs-demo-analyzer/tests/assertion"
	"github.com/akiver/cs-demo-analyzer/tests/testsutils"
)

func Test_Pracc_CHHyStHPwEnGVsFmrZLWsDQwBX_2026_nuke(t *testing.T) {
	demoName := "pracc_CHHyStHPwEnGVsFmrZLWsDQwBX_2026_nuke"
	demoPath := testsutils.GetDemoPath("cs2", demoName)
	match, err := api.AnalyzeDemo(demoPath, api.AnalyzeDemoOptions{
		Source: constants.DemoSourcePracc,
	})

	if err != nil {
		t.Fatal(err)
	}

	assertion.AssertMatchSnapshot(t, match, demoName)
}
