package tests

import (
	"testing"

	"github.com/akiver/cs-demo-analyzer/pkg/api"
	"github.com/akiver/cs-demo-analyzer/pkg/api/constants"
	"github.com/akiver/cs-demo-analyzer/tests/assertion"
	"github.com/akiver/cs-demo-analyzer/tests/fake"
	"github.com/akiver/cs-demo-analyzer/tests/testsutils"
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/common"
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/events"
)

// Pause requested at the end of the round 3 and effective at the beginning of the round 4 (tick 32229).
// The eBot backup is restored at tick 43560, that's when the round 4 really starts.
// The round 6 is cancelled few seconds after it started, a backup is restored and the round 6 really start.
// https://www.hltv.org/stats/matches/mapstatsid/38754/cloud9-vs-ninjas-in-pyjamas
func TestEbot_Cloud9_VS_NIP_IEM_Oakland_2016_Train(t *testing.T) {
	demoName := "ebot_cloud9_vs_nip_iem_oakland_2016_train"
	demoPath := testsutils.GetDemoPath("csgo", demoName)
	match, err := api.AnalyzeDemo(demoPath, api.AnalyzeDemoOptions{
		Source: constants.DemoSourceEbot,
	})
	if err != nil {
		t.Error(err)
	}

	expectedRoundCount := 30
	expectedPlayerCount := 10
	expectedScoreTeamA := 14
	expectedScoreTeamB := 16
	expectedScoreFirstHalfTeamA := 5
	expectedScoreFirstHalfTeamB := 10
	expectedScoreSecondHalfTeamA := 9
	expectedScoreSecondHalfTeamB := 6
	expectedTeamNameA := "Cloud9"
	expectedTeamNameB := "Ninjas in Pyjamas"
	expectedWinnerName := expectedTeamNameB

	if len(match.Rounds) != expectedRoundCount {
		t.Errorf("expected %d rounds but got %d", expectedRoundCount, len(match.Rounds))
	}
	if len(match.Players()) != expectedPlayerCount {
		t.Errorf("expected %d players but got %d", expectedPlayerCount, len(match.Players()))
	}
	if match.TeamA.Score != expectedScoreTeamA {
		t.Errorf("expected score team A to be %d got %d", expectedScoreTeamA, match.TeamA.Score)
	}
	if match.TeamB.Score != expectedScoreTeamB {
		t.Errorf("expected score team B to be %d got %d", expectedScoreTeamB, match.TeamB.Score)
	}
	if match.TeamA.Name != expectedTeamNameA {
		t.Errorf("expected team name A to be %s got %s", expectedTeamNameA, match.TeamA.Name)
	}
	if match.TeamB.Name != expectedTeamNameB {
		t.Errorf("expected team name B to be %s got %s", expectedTeamNameB, match.TeamB.Name)
	}
	if match.TeamA.ScoreFirstHalf != expectedScoreFirstHalfTeamA {
		t.Errorf("expected score first half team A to be %d got %d", expectedScoreFirstHalfTeamA, match.TeamA.ScoreFirstHalf)
	}
	if match.TeamB.ScoreFirstHalf != expectedScoreFirstHalfTeamB {
		t.Errorf("expected score first half team B to be %d got %d", expectedScoreFirstHalfTeamB, match.TeamB.ScoreFirstHalf)
	}
	if match.TeamA.ScoreSecondHalf != expectedScoreSecondHalfTeamA {
		t.Errorf("expected score second half team A to be %d got %d", expectedScoreSecondHalfTeamA, match.TeamA.ScoreSecondHalf)
	}
	if match.TeamB.ScoreSecondHalf != expectedScoreSecondHalfTeamB {
		t.Errorf("expected score second half team B to be %d got %d", expectedScoreSecondHalfTeamB, match.TeamB.ScoreSecondHalf)
	}
	if match.Winner.Name != expectedWinnerName {
		t.Errorf("expected winner to be %s but got %s", expectedWinnerName, match.Winner.Name)
	}

	var rounds = []fake.FakeRound{
		{
			Number:            1,
			StartTick:         0,
			StartFrame:        22,
			EndTick:           8242,
			FreezeTimeEndTick: 0,
			EndOfficiallyTick: 8882,
			TeamAStartMoney:   4000,
			TeamBStartMoney:   4000,
			TeamAEconomyType:  constants.EconomyTypePistol,
			TeamBEconomyType:  constants.EconomyTypePistol,
			EndReason:         events.RoundEndReasonBombDefused,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        1,
			TeamBScore:        0,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
		},
		{
			Number:            2,
			StartTick:         8882,
			StartFrame:        8901,
			EndTick:           20239,
			EndOfficiallyTick: 20878,
			FreezeTimeEndTick: 10802,
			TeamAStartMoney:   19950,
			TeamBStartMoney:   13050,
			TeamAEconomyType:  constants.EconomyTypeSemi,
			TeamBEconomyType:  constants.EconomyTypeSemi,
			EndReason:         events.RoundEndReasonCTWin,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        2,
			TeamBScore:        0,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
		},
		{
			Number:            3,
			StartTick:         20878,
			StartFrame:        20878,
			EndTick:           31587,
			EndOfficiallyTick: 32229,
			FreezeTimeEndTick: 22798,
			TeamAStartMoney:   18900,
			TeamBStartMoney:   18850,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeForceBuy,
			EndReason:         events.RoundEndReasonTerroristsWin,
			WinnerSide:        common.TeamTerrorists,
			TeamAScore:        2,
			TeamBScore:        1,
			WinnerName:        expectedTeamNameB,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
		},
		{
			Number:            4,
			StartTick:         43560,
			StartFrame:        43525,
			EndTick:           58032,
			EndOfficiallyTick: 58672,
			FreezeTimeEndTick: 47172,
			TeamAStartMoney:   12350,
			TeamBStartMoney:   18300,
			TeamAEconomyType:  constants.EconomyTypeEco,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonTerroristsWin,
			WinnerSide:        common.TeamTerrorists,
			TeamAScore:        2,
			TeamBScore:        2,
			WinnerName:        expectedTeamNameB,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
		},
		{
			Number:            5,
			StartTick:         58672,
			StartFrame:        58625,
			EndTick:           66497,
			EndOfficiallyTick: 67137,
			FreezeTimeEndTick: 60592,
			TeamAStartMoney:   21550,
			TeamBStartMoney:   23950,
			TeamAEconomyType:  constants.EconomyTypeForceBuy,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonTerroristsWin,
			WinnerSide:        common.TeamTerrorists,
			TeamAScore:        2,
			TeamBScore:        3,
			WinnerName:        expectedTeamNameB,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
		},
		{
			Number:            6,
			StartTick:         72472,
			StartFrame:        72377,
			EndTick:           91927,
			EndOfficiallyTick: 92567,
			FreezeTimeEndTick: 76523,
			TeamAStartMoney:   12750,
			TeamBStartMoney:   36900,
			TeamAEconomyType:  constants.EconomyTypeEco,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonTerroristsWin,
			WinnerSide:        common.TeamTerrorists,
			TeamAScore:        2,
			TeamBScore:        4,
			WinnerName:        expectedTeamNameB,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
		},
		{
			Number:            7,
			StartTick:         92567,
			StartFrame:        92460,
			EndTick:           99676,
			EndOfficiallyTick: 100317,
			FreezeTimeEndTick: 94487,
			TeamAStartMoney:   26550,
			TeamBStartMoney:   42150,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonCTWin,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        3,
			TeamBScore:        4,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
		},
		{
			Number:            8,
			StartTick:         100317,
			StartFrame:        100195,
			EndTick:           110634,
			EndOfficiallyTick: 111275,
			FreezeTimeEndTick: 102237,
			TeamAStartMoney:   17250,
			TeamBStartMoney:   29050,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonTerroristsWin,
			WinnerSide:        common.TeamTerrorists,
			TeamAScore:        3,
			TeamBScore:        5,
			WinnerName:        expectedTeamNameB,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
		},
		{
			Number:            9,
			StartTick:         111275,
			StartFrame:        111130,
			EndTick:           124473,
			EndOfficiallyTick: 125114,
			FreezeTimeEndTick: 113195,
			TeamAStartMoney:   12900,
			TeamBStartMoney:   23000,
			TeamAEconomyType:  constants.EconomyTypeEco,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonTerroristsWin,
			WinnerSide:        common.TeamTerrorists,
			TeamAScore:        3,
			TeamBScore:        6,
			WinnerName:        expectedTeamNameB,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
		},
		{
			Number:            10,
			StartTick:         125114,
			StartFrame:        124955,
			EndTick:           129950,
			EndOfficiallyTick: 130590,
			FreezeTimeEndTick: 127034,
			TeamAStartMoney:   20300,
			TeamBStartMoney:   31750,
			TeamAEconomyType:  constants.EconomyTypeEco,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonTerroristsWin,
			WinnerSide:        common.TeamTerrorists,
			TeamAScore:        3,
			TeamBScore:        7,
			WinnerName:        expectedTeamNameB,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
		},
		{
			Number:            11,
			StartTick:         130590,
			StartFrame:        130416,
			EndTick:           136801,
			EndOfficiallyTick: 137442,
			FreezeTimeEndTick: 132510,
			TeamAStartMoney:   28700,
			TeamBStartMoney:   41900,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonCTWin,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        4,
			TeamBScore:        7,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
		},
		{
			Number:            12,
			StartTick:         137442,
			StartFrame:        137249,
			EndTick:           150526,
			EndOfficiallyTick: 151166,
			FreezeTimeEndTick: 139362,
			TeamAStartMoney:   18800,
			TeamBStartMoney:   41600,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonTargetBombed,
			WinnerSide:        common.TeamTerrorists,
			TeamAScore:        4,
			TeamBScore:        8,
			WinnerName:        expectedTeamNameB,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
		},
		{
			Number:            13,
			StartTick:         151166,
			StartFrame:        150952,
			EndTick:           160945,
			EndOfficiallyTick: 161584,
			FreezeTimeEndTick: 153086,
			TeamAStartMoney:   16550,
			TeamBStartMoney:   36200,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonTerroristsWin,
			WinnerSide:        common.TeamTerrorists,
			TeamAScore:        4,
			TeamBScore:        9,
			WinnerName:        expectedTeamNameB,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
		},
		{
			Number:            14,
			StartTick:         161584,
			StartFrame:        161355,
			EndTick:           174100,
			EndOfficiallyTick: 174740,
			FreezeTimeEndTick: 163504,
			TeamAStartMoney:   13550,
			TeamBStartMoney:   34800,
			TeamAEconomyType:  constants.EconomyTypeEco,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonTerroristsWin,
			WinnerSide:        common.TeamTerrorists,
			TeamAScore:        4,
			TeamBScore:        10,
			WinnerName:        expectedTeamNameB,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
		},
		{
			Number:            15,
			StartTick:         174740,
			StartFrame:        174492,
			EndTick:           188090,
			EndOfficiallyTick: 192906,
			FreezeTimeEndTick: 176660,
			TeamAStartMoney:   22800,
			TeamBStartMoney:   31650,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonCTWin,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        5,
			TeamBScore:        10,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
		},
		{
			Number:            16,
			StartTick:         192906,
			StartFrame:        192641,
			EndTick:           217033,
			EndOfficiallyTick: 217674,
			FreezeTimeEndTick: 209846,
			TeamAStartMoney:   4000,
			TeamBStartMoney:   4000,
			TeamAEconomyType:  constants.EconomyTypePistol,
			TeamBEconomyType:  constants.EconomyTypePistol,
			EndReason:         events.RoundEndReasonCTWin,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        5,
			TeamBScore:        11,
			WinnerName:        expectedTeamNameB,
			TeamASide:         common.TeamTerrorists,
			TeamBSide:         common.TeamCounterTerrorists,
		},
		{
			Number:            17,
			StartTick:         217674,
			StartFrame:        217391,
			EndTick:           230872,
			EndOfficiallyTick: 231512,
			FreezeTimeEndTick: 219594,
			TeamAStartMoney:   8300,
			TeamBStartMoney:   18400,
			TeamAEconomyType:  constants.EconomyTypeForceBuy,
			TeamBEconomyType:  constants.EconomyTypeSemi,
			EndReason:         events.RoundEndReasonCTWin,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        5,
			TeamBScore:        12,
			WinnerName:        expectedTeamNameB,
			TeamASide:         common.TeamTerrorists,
			TeamBSide:         common.TeamCounterTerrorists,
		},
		{
			Number:            18,
			StartTick:         231512,
			StartFrame:        231212,
			EndTick:           242398,
			EndOfficiallyTick: 243038,
			FreezeTimeEndTick: 233432,
			TeamAStartMoney:   10150,
			TeamBStartMoney:   19100,
			TeamAEconomyType:  constants.EconomyTypeEco,
			TeamBEconomyType:  constants.EconomyTypeSemi,
			EndReason:         events.RoundEndReasonCTWin,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        5,
			TeamBScore:        13,
			WinnerName:        expectedTeamNameB,
			TeamASide:         common.TeamTerrorists,
			TeamBSide:         common.TeamCounterTerrorists,
		},
		{
			Number:            19,
			StartTick:         243038,
			StartFrame:        242722,
			EndTick:           253655,
			EndOfficiallyTick: 254296,
			FreezeTimeEndTick: 244958,
			TeamAStartMoney:   22550,
			TeamBStartMoney:   31350,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonTerroristsWin,
			WinnerSide:        common.TeamTerrorists,
			TeamAScore:        6,
			TeamBScore:        13,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamTerrorists,
			TeamBSide:         common.TeamCounterTerrorists,
		},
		{
			Number:            20,
			StartTick:         254296,
			StartFrame:        253929,
			EndTick:           274777,
			EndOfficiallyTick: 275417,
			FreezeTimeEndTick: 256216,
			TeamAStartMoney:   18400,
			TeamBStartMoney:   23700,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonTargetBombed,
			WinnerSide:        common.TeamTerrorists,
			TeamAScore:        7,
			TeamBScore:        13,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamTerrorists,
			TeamBSide:         common.TeamCounterTerrorists,
		},
		{
			Number:            21,
			StartTick:         275417,
			StartFrame:        275027,
			EndTick:           291042,
			EndOfficiallyTick: 291682,
			FreezeTimeEndTick: 277337,
			TeamAStartMoney:   29400,
			TeamBStartMoney:   10350,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeSemi,
			EndReason:         events.RoundEndReasonTargetBombed,
			WinnerSide:        common.TeamTerrorists,
			TeamAScore:        8,
			TeamBScore:        13,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamTerrorists,
			TeamBSide:         common.TeamCounterTerrorists,
		},
		{
			Number:            22,
			StartTick:         291682,
			StartFrame:        291219,
			EndTick:           306888,
			EndOfficiallyTick: 307527,
			FreezeTimeEndTick: 293602,
			TeamAStartMoney:   33250,
			TeamBStartMoney:   20050,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeSemi,
			EndReason:         events.RoundEndReasonTerroristsWin,
			WinnerSide:        common.TeamTerrorists,
			TeamAScore:        9,
			TeamBScore:        13,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamTerrorists,
			TeamBSide:         common.TeamCounterTerrorists,
		},
		{
			Number:            23,
			StartTick:         307527,
			StartFrame:        307040,
			EndTick:           315495,
			EndOfficiallyTick: 316136,
			FreezeTimeEndTick: 309447,
			TeamAStartMoney:   39100,
			TeamBStartMoney:   30800,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonTerroristsWin,
			WinnerSide:        common.TeamTerrorists,
			TeamAScore:        10,
			TeamBScore:        13,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamTerrorists,
			TeamBSide:         common.TeamCounterTerrorists,
		},
		{
			Number:            24,
			StartTick:         316136,
			StartFrame:        315631,
			EndTick:           327856,
			EndOfficiallyTick: 328496,
			FreezeTimeEndTick: 318056,
			TeamAStartMoney:   40850,
			TeamBStartMoney:   19250,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeEco,
			EndReason:         events.RoundEndReasonTerroristsWin,
			WinnerSide:        common.TeamTerrorists,
			TeamAScore:        11,
			TeamBScore:        13,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamTerrorists,
			TeamBSide:         common.TeamCounterTerrorists,
		},
		{
			Number:            25,
			StartTick:         328496,
			StartFrame:        327967,
			EndTick:           344219,
			EndOfficiallyTick: 344859,
			FreezeTimeEndTick: 330416,
			TeamAStartMoney:   41850,
			TeamBStartMoney:   33250,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonBombDefused,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        11,
			TeamBScore:        14,
			WinnerName:        expectedTeamNameB,
			TeamASide:         common.TeamTerrorists,
			TeamBSide:         common.TeamCounterTerrorists,
		},
		{
			Number:            26,
			StartTick:         344859,
			StartFrame:        344309,
			EndTick:           359052,
			EndOfficiallyTick: 359692,
			FreezeTimeEndTick: 346779,
			TeamAStartMoney:   41950,
			TeamBStartMoney:   21600,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonCTWin,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        11,
			TeamBScore:        15,
			WinnerName:        expectedTeamNameB,
			TeamASide:         common.TeamTerrorists,
			TeamBSide:         common.TeamCounterTerrorists,
		},
		{
			Number:            27,
			StartTick:         359692,
			StartFrame:        359026,
			EndTick:           369869,
			EndOfficiallyTick: 370509,
			FreezeTimeEndTick: 361612,
			TeamAStartMoney:   25900,
			TeamBStartMoney:   24000,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonTerroristsWin,
			WinnerSide:        common.TeamTerrorists,
			TeamAScore:        12,
			TeamBScore:        15,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamTerrorists,
			TeamBSide:         common.TeamCounterTerrorists,
		},
		{
			Number:            28,
			StartTick:         370509,
			StartFrame:        369824,
			EndTick:           384654,
			EndOfficiallyTick: 385295,
			FreezeTimeEndTick: 372429,
			TeamAStartMoney:   18500,
			TeamBStartMoney:   10400,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeEco,
			EndReason:         events.RoundEndReasonTerroristsWin,
			WinnerSide:        common.TeamTerrorists,
			TeamAScore:        13,
			TeamBScore:        15,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamTerrorists,
			TeamBSide:         common.TeamCounterTerrorists,
		},
		{
			Number:            29,
			StartTick:         385295,
			StartFrame:        384593,
			EndTick:           389464,
			EndOfficiallyTick: 390106,
			FreezeTimeEndTick: 387215,
			TeamAStartMoney:   20750,
			TeamBStartMoney:   16400,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeEco,
			EndReason:         events.RoundEndReasonTerroristsWin,
			WinnerSide:        common.TeamTerrorists,
			TeamAScore:        14,
			TeamBScore:        15,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamTerrorists,
			TeamBSide:         common.TeamCounterTerrorists,
		},
		{
			Number:            30,
			StartTick:         390106,
			StartFrame:        389389,
			EndTick:           408156,
			EndOfficiallyTick: 408156,
			FreezeTimeEndTick: 392026,
			TeamAStartMoney:   25400,
			TeamBStartMoney:   25700,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonBombDefused,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        14,
			TeamBScore:        16,
			WinnerName:        expectedTeamNameB,
			TeamASide:         common.TeamTerrorists,
			TeamBSide:         common.TeamCounterTerrorists,
		},
	}

	assertion.AssertRounds(t, match, rounds)
}
