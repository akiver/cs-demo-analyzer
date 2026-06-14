package tests

import (
	"testing"

	"github.com/akiver/cs-demo-analyzer/pkg/api"
	"github.com/akiver/cs-demo-analyzer/pkg/api/constants"
	"github.com/akiver/cs-demo-analyzer/tests/assertion"
	"github.com/akiver/cs-demo-analyzer/tests/testsutils"
)

// https://www.challengermode.com/s/CsgoAllstars/games/b4d31195-bae0-42c0-bbcc-08dbf6b17863
func TestChallengerMode_6c306e56_8170_4092_b402_08dbf813e452_2023_Anubis(t *testing.T) {
	t.Parallel()

	demoName := "challengermode_6c306e56-8170-4092-b402-08dbf813e452"
	demoPath := testsutils.GetDemoPath("cs2", demoName)
	match, err := api.AnalyzeDemo(demoPath, api.AnalyzeDemoOptions{
		Source: constants.DemoSourceFaceIt,
	})

	if err != nil {
		t.Fatal(err)
	}

	assertion.AssertMatchSnapshot(t, match, demoName)
}
