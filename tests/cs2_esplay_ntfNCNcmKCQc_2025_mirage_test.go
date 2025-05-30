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

// 5v5 with overtimes.
// https://esplay.com/m/ntfNCNcmKCQc/team-lina-vs-team-qara
func TestEsplay_ntfNCNcmKCQc_2025_Mirage(t *testing.T) {
	demoName := "esplay_ntfNCNcmKCQc_2025_mirage"
	demoPath := testsutils.GetDemoPath("cs2", demoName)
	match, err := api.AnalyzeDemo(demoPath, api.AnalyzeDemoOptions{
		Source: constants.DemoSourceEsplay,
	})

	if err != nil {
		t.Error(err)
	}

	expectedRoundCount := 36
	expectedPlayerCount := 10
	expectedScoreTeamA := 19
	expectedScoreTeamB := 17
	expectedScoreFirstHalfTeamA := 5
	expectedScoreFirstHalfTeamB := 7
	expectedScoreSecondHalfTeamA := 7
	expectedScoreSecondHalfTeamB := 5
	expectedTeamNameA := "Team A"
	expectedTeamNameB := "Team B"
	expectedWinnerName := expectedTeamNameA
	expectedMaxRounds := 24

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
	if match.MaxRounds != expectedMaxRounds {
		t.Errorf("expected max rounds to be %d but got %d", expectedMaxRounds, match.MaxRounds)
	}

	var rounds = []fake.FakeRound{
		{
			Number:            1,
			StartTick:         0,
			StartFrame:        1175,
			EndTick:           6229,
			FreezeTimeEndTick: 895,
			TeamAStartMoney:   4000,
			TeamBStartMoney:   4000,
			TeamAEconomyType:  constants.EconomyTypePistol,
			TeamBEconomyType:  constants.EconomyTypePistol,
			EndReason:         events.RoundEndReasonTargetBombed,
			WinnerSide:        common.TeamTerrorists,
			TeamAScore:        0,
			TeamBScore:        1,
			WinnerName:        expectedTeamNameB,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            2,
			StartTick:         6677,
			StartFrame:        9312,
			EndTick:           12719,
			FreezeTimeEndTick: 7573,
			TeamAStartMoney:   11150,
			TeamBStartMoney:   19700,
			TeamAEconomyType:  constants.EconomyTypeEco,
			TeamBEconomyType:  constants.EconomyTypeSemi,
			EndReason:         events.RoundEndReasonTerroristsWin,
			WinnerSide:        common.TeamTerrorists,
			TeamAScore:        0,
			TeamBScore:        2,
			WinnerName:        expectedTeamNameB,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            3,
			StartTick:         13167,
			StartFrame:        17384,
			EndTick:           20878,
			FreezeTimeEndTick: 14063,
			TeamAStartMoney:   20100,
			TeamBStartMoney:   20100,
			TeamAEconomyType:  constants.EconomyTypeForceBuy,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonTerroristsWin,
			WinnerSide:        common.TeamTerrorists,
			TeamAScore:        0,
			TeamBScore:        3,
			WinnerName:        expectedTeamNameB,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            4,
			StartTick:         21326,
			StartFrame:        27198,
			EndTick:           27036,
			FreezeTimeEndTick: 22222,
			TeamAStartMoney:   15750,
			TeamBStartMoney:   26150,
			TeamAEconomyType:  constants.EconomyTypeSemi,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonTerroristsWin,
			WinnerSide:        common.TeamTerrorists,
			TeamAScore:        0,
			TeamBScore:        4,
			WinnerName:        expectedTeamNameB,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            5,
			StartTick:         27484,
			StartFrame:        35004,
			EndTick:           33956,
			FreezeTimeEndTick: 28380,
			TeamAStartMoney:   28550,
			TeamBStartMoney:   28000,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonCTWin,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        1,
			TeamBScore:        4,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            6,
			StartTick:         34404,
			StartFrame:        43524,
			EndTick:           39596,
			FreezeTimeEndTick: 35300,
			TeamAStartMoney:   18550,
			TeamBStartMoney:   20500,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonTerroristsWin,
			WinnerSide:        common.TeamTerrorists,
			TeamAScore:        1,
			TeamBScore:        5,
			WinnerName:        expectedTeamNameB,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            7,
			StartTick:         40044,
			StartFrame:        50489,
			EndTick:           47736,
			FreezeTimeEndTick: 40940,
			TeamAStartMoney:   15750,
			TeamBStartMoney:   19550,
			TeamAEconomyType:  constants.EconomyTypeSemi,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonTargetBombed,
			WinnerSide:        common.TeamTerrorists,
			TeamAScore:        1,
			TeamBScore:        6,
			WinnerName:        expectedTeamNameB,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            8,
			StartTick:         48184,
			StartFrame:        60712,
			EndTick:           54612,
			FreezeTimeEndTick: 49080,
			TeamAStartMoney:   27650,
			TeamBStartMoney:   25850,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonTerroristsWin,
			WinnerSide:        common.TeamTerrorists,
			TeamAScore:        1,
			TeamBScore:        7,
			WinnerName:        expectedTeamNameB,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            9,
			StartTick:         55060,
			StartFrame:        69193,
			EndTick:           60299,
			FreezeTimeEndTick: 55956,
			TeamAStartMoney:   26050,
			TeamBStartMoney:   30650,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonCTWin,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        2,
			TeamBScore:        7,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            10,
			StartTick:         60747,
			StartFrame:        76285,
			EndTick:           69837,
			FreezeTimeEndTick: 61643,
			TeamAStartMoney:   19250,
			TeamBStartMoney:   28000,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonBombDefused,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        3,
			TeamBScore:        7,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            11,
			StartTick:         74153,
			StartFrame:        93486,
			EndTick:           77290,
			FreezeTimeEndTick: 75049,
			TeamAStartMoney:   21850,
			TeamBStartMoney:   19450,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeForceBuy,
			EndReason:         events.RoundEndReasonCTWin,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        4,
			TeamBScore:        7,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            12,
			StartTick:         77738,
			StartFrame:        97812,
			EndTick:           82342,
			FreezeTimeEndTick: 78634,
			TeamAStartMoney:   22700,
			TeamBStartMoney:   13350,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeForceBuy,
			EndReason:         events.RoundEndReasonBombDefused,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        5,
			TeamBScore:        7,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            13,
			StartTick:         83302,
			StartFrame:        104899,
			EndTick:           86088,
			FreezeTimeEndTick: 84198,
			TeamAStartMoney:   4000,
			TeamBStartMoney:   4000,
			TeamAEconomyType:  constants.EconomyTypePistol,
			TeamBEconomyType:  constants.EconomyTypePistol,
			EndReason:         events.RoundEndReasonCTWin,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        5,
			TeamBScore:        8,
			WinnerName:        expectedTeamNameB,
			TeamASide:         common.TeamTerrorists,
			TeamBSide:         common.TeamCounterTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            14,
			StartTick:         86536,
			StartFrame:        108996,
			EndTick:           90719,
			FreezeTimeEndTick: 87432,
			TeamAStartMoney:   10450,
			TeamBStartMoney:   18450,
			TeamAEconomyType:  constants.EconomyTypeEco,
			TeamBEconomyType:  constants.EconomyTypeSemi,
			EndReason:         events.RoundEndReasonCTWin,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        5,
			TeamBScore:        9,
			WinnerName:        expectedTeamNameB,
			TeamASide:         common.TeamTerrorists,
			TeamBSide:         common.TeamCounterTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            15,
			StartTick:         91167,
			StartFrame:        114627,
			EndTick:           99178,
			FreezeTimeEndTick: 92063,
			TeamAStartMoney:   22650,
			TeamBStartMoney:   19350,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonTerroristsWin,
			WinnerSide:        common.TeamTerrorists,
			TeamAScore:        6,
			TeamBScore:        9,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamTerrorists,
			TeamBSide:         common.TeamCounterTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            16,
			StartTick:         99626,
			StartFrame:        124754,
			EndTick:           105130,
			FreezeTimeEndTick: 100522,
			TeamAStartMoney:   18550,
			TeamBStartMoney:   13950,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeForceBuy,
			EndReason:         events.RoundEndReasonTerroristsWin,
			WinnerSide:        common.TeamTerrorists,
			TeamAScore:        7,
			TeamBScore:        9,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamTerrorists,
			TeamBSide:         common.TeamCounterTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            17,
			StartTick:         105578,
			StartFrame:        132395,
			EndTick:           109921,
			FreezeTimeEndTick: 106474,
			TeamAStartMoney:   19700,
			TeamBStartMoney:   11950,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeEco,
			EndReason:         events.RoundEndReasonTerroristsWin,
			WinnerSide:        common.TeamTerrorists,
			TeamAScore:        8,
			TeamBScore:        9,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamTerrorists,
			TeamBSide:         common.TeamCounterTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            18,
			StartTick:         110369,
			StartFrame:        138176,
			EndTick:           116186,
			FreezeTimeEndTick: 111265,
			TeamAStartMoney:   23700,
			TeamBStartMoney:   20950,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeForceBuy,
			EndReason:         events.RoundEndReasonTargetBombed,
			WinnerSide:        common.TeamTerrorists,
			TeamAScore:        9,
			TeamBScore:        9,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamTerrorists,
			TeamBSide:         common.TeamCounterTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            19,
			StartTick:         116634,
			StartFrame:        145994,
			EndTick:           122361,
			FreezeTimeEndTick: 117530,
			TeamAStartMoney:   28500,
			TeamBStartMoney:   15600,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeSemi,
			EndReason:         events.RoundEndReasonTerroristsWin,
			WinnerSide:        common.TeamTerrorists,
			TeamAScore:        10,
			TeamBScore:        9,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamTerrorists,
			TeamBSide:         common.TeamCounterTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            20,
			StartTick:         122809,
			StartFrame:        153658,
			EndTick:           130227,
			FreezeTimeEndTick: 123705,
			TeamAStartMoney:   33150,
			TeamBStartMoney:   26450,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonTerroristsWin,
			WinnerSide:        common.TeamTerrorists,
			TeamAScore:        11,
			TeamBScore:        9,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamTerrorists,
			TeamBSide:         common.TeamCounterTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            21,
			StartTick:         130675,
			StartFrame:        163703,
			EndTick:           134120,
			FreezeTimeEndTick: 131571,
			TeamAStartMoney:   39050,
			TeamBStartMoney:   19650,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeForceBuy,
			EndReason:         events.RoundEndReasonCTWin,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        11,
			TeamBScore:        10,
			WinnerName:        expectedTeamNameB,
			TeamASide:         common.TeamTerrorists,
			TeamBSide:         common.TeamCounterTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            22,
			StartTick:         134568,
			StartFrame:        168517,
			EndTick:           139987,
			FreezeTimeEndTick: 135464,
			TeamAStartMoney:   27950,
			TeamBStartMoney:   18750,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonCTWin,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        11,
			TeamBScore:        11,
			WinnerName:        expectedTeamNameB,
			TeamASide:         common.TeamTerrorists,
			TeamBSide:         common.TeamCounterTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            23,
			StartTick:         140435,
			StartFrame:        175870,
			EndTick:           145534,
			FreezeTimeEndTick: 141331,
			TeamAStartMoney:   13950,
			TeamBStartMoney:   26150,
			TeamAEconomyType:  constants.EconomyTypeEco,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonCTWin,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        11,
			TeamBScore:        12,
			WinnerName:        expectedTeamNameB,
			TeamASide:         common.TeamTerrorists,
			TeamBSide:         common.TeamCounterTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            24,
			StartTick:         145982,
			StartFrame:        182571,
			EndTick:           152214,
			FreezeTimeEndTick: 146878,
			TeamAStartMoney:   21800,
			TeamBStartMoney:   30800,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonTerroristsWin,
			WinnerSide:        common.TeamTerrorists,
			TeamAScore:        12,
			TeamBScore:        12,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamTerrorists,
			TeamBSide:         common.TeamCounterTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			OvertimeNumber:    1,
			Number:            25,
			StartTick:         153174,
			StartFrame:        191420,
			EndTick:           160132,
			FreezeTimeEndTick: 154070,
			TeamAStartMoney:   50000,
			TeamBStartMoney:   50000,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonCTWin,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        12,
			TeamBScore:        13,
			WinnerName:        expectedTeamNameB,
			TeamASide:         common.TeamTerrorists,
			TeamBSide:         common.TeamCounterTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			OvertimeNumber:    1,
			Number:            26,
			StartTick:         160580,
			StartFrame:        200737,
			EndTick:           166754,
			FreezeTimeEndTick: 161476,
			TeamAStartMoney:   33750,
			TeamBStartMoney:   37550,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonTerroristsWin,
			WinnerSide:        common.TeamTerrorists,
			TeamAScore:        13,
			TeamBScore:        13,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamTerrorists,
			TeamBSide:         common.TeamCounterTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			OvertimeNumber:    1,
			Number:            27,
			StartTick:         167202,
			StartFrame:        208806,
			EndTick:           172497,
			FreezeTimeEndTick: 168098,
			TeamAStartMoney:   27300,
			TeamBStartMoney:   30650,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonCTWin,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        13,
			TeamBScore:        14,
			WinnerName:        expectedTeamNameB,
			TeamASide:         common.TeamTerrorists,
			TeamBSide:         common.TeamCounterTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			OvertimeNumber:    1,
			Number:            28,
			StartTick:         173457,
			StartFrame:        216726,
			EndTick:           179731,
			FreezeTimeEndTick: 174353,
			TeamAStartMoney:   50000,
			TeamBStartMoney:   50000,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonTerroristsWin,
			WinnerSide:        common.TeamTerrorists,
			TeamAScore:        13,
			TeamBScore:        15,
			WinnerName:        expectedTeamNameB,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			OvertimeNumber:    1,
			Number:            29,
			StartTick:         180179,
			StartFrame:        225668,
			EndTick:           183773,
			FreezeTimeEndTick: 181075,
			TeamAStartMoney:   31200,
			TeamBStartMoney:   43450,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonCTWin,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        14,
			TeamBScore:        15,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			OvertimeNumber:    1,
			Number:            30,
			StartTick:         184221,
			StartFrame:        230486,
			EndTick:           189353,
			FreezeTimeEndTick: 185117,
			TeamAStartMoney:   21850,
			TeamBStartMoney:   40700,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonBombDefused,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        15,
			TeamBScore:        15,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			OvertimeNumber:    2,
			Number:            31,
			StartTick:         190313,
			StartFrame:        237907,
			EndTick:           198980,
			FreezeTimeEndTick: 191209,
			TeamAStartMoney:   50000,
			TeamBStartMoney:   50000,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonBombDefused,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        16,
			TeamBScore:        15,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			OvertimeNumber:    2,
			Number:            32,
			StartTick:         199428,
			StartFrame:        249104,
			EndTick:           204533,
			FreezeTimeEndTick: 200324,
			TeamAStartMoney:   39750,
			TeamBStartMoney:   39100,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonCTWin,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        17,
			TeamBScore:        15,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			OvertimeNumber:    2,
			Number:            33,
			StartTick:         204981,
			StartFrame:        256123,
			EndTick:           211540,
			FreezeTimeEndTick: 205877,
			TeamAStartMoney:   34650,
			TeamBStartMoney:   26900,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonCTWin,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        18,
			TeamBScore:        15,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			OvertimeNumber:    2,
			Number:            34,
			StartTick:         212500,
			StartFrame:        265883,
			EndTick:           216431,
			FreezeTimeEndTick: 213396,
			TeamAStartMoney:   50000,
			TeamBStartMoney:   50000,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonCTWin,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        18,
			TeamBScore:        16,
			WinnerName:        expectedTeamNameB,
			TeamASide:         common.TeamTerrorists,
			TeamBSide:         common.TeamCounterTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			OvertimeNumber:    2,
			Number:            35,
			StartTick:         216879,
			StartFrame:        271428,
			EndTick:           221390,
			FreezeTimeEndTick: 217775,
			TeamAStartMoney:   35100,
			TeamBStartMoney:   39050,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonCTWin,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        18,
			TeamBScore:        17,
			WinnerName:        expectedTeamNameB,
			TeamASide:         common.TeamTerrorists,
			TeamBSide:         common.TeamCounterTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			OvertimeNumber:    2,
			Number:            36,
			StartTick:         221838,
			StartFrame:        277561,
			EndTick:           227687,
			FreezeTimeEndTick: 222734,
			TeamAStartMoney:   18550,
			TeamBStartMoney:   37100,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonTerroristsWin,
			WinnerSide:        common.TeamTerrorists,
			TeamAScore:        19,
			TeamBScore:        17,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamTerrorists,
			TeamBSide:         common.TeamCounterTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
	}

	var players = []fake.FakePlayer{
		{
			SteamID64:     76561197965932705,
			Name:          "DanDolme",
			KillCount:     24,
			AssistCount:   10,
			DeathCount:    25,
			Score:         61,
			Team:          match.TeamA,
			HeadshotCount: 12,
			MvpCount:      4,
			UtilityDamage: 141,
		},
		{
			SteamID64:     76561197960361674,
			Name:          "foxm",
			KillCount:     11,
			AssistCount:   11,
			DeathCount:    24,
			Score:         46,
			Team:          match.TeamA,
			HeadshotCount: 6,
			MvpCount:      1,
			UtilityDamage: 297,
		},
		{
			SteamID64:     76561197961094454,
			Name:          "steveness",
			KillCount:     21,
			AssistCount:   5,
			DeathCount:    29,
			Score:         50,
			Team:          match.TeamA,
			HeadshotCount: 9,
			MvpCount:      2,
			UtilityDamage: 157,
		},
		{
			SteamID64:     76561198123034769,
			Name:          "Qara",
			KillCount:     35,
			AssistCount:   9,
			DeathCount:    26,
			Score:         86,
			Team:          match.TeamA,
			HeadshotCount: 20,
			MvpCount:      7,
			UtilityDamage: 19,
		},
		{
			SteamID64:     76561198078476229,
			Name:          "olii1",
			KillCount:     35,
			AssistCount:   6,
			DeathCount:    21,
			Score:         79,
			Team:          match.TeamA,
			HeadshotCount: 10,
			MvpCount:      5,
			UtilityDamage: 67,
		},
		{
			SteamID64:     76561198119455532,
			Name:          "spaceFocus",
			KillCount:     24,
			AssistCount:   10,
			DeathCount:    25,
			Score:         63,
			Team:          match.TeamB,
			HeadshotCount: 12,
			MvpCount:      5,
			UtilityDamage: 210,
		},
		{
			SteamID64:     76561198259994010,
			Name:          "Lina",
			KillCount:     27,
			AssistCount:   8,
			DeathCount:    26,
			Score:         60,
			Team:          match.TeamB,
			HeadshotCount: 13,
			MvpCount:      3,
			UtilityDamage: 139,
		},
		{
			SteamID64:     76561199045841130,
			Name:          "Snurri",
			KillCount:     20,
			AssistCount:   6,
			DeathCount:    25,
			Score:         48,
			Team:          match.TeamB,
			HeadshotCount: 9,
			MvpCount:      2,
			UtilityDamage: 39,
		},
		{
			SteamID64:     76561199060069995,
			Name:          "henke88",
			KillCount:     28,
			AssistCount:   7,
			DeathCount:    25,
			Score:         71,
			Team:          match.TeamB,
			HeadshotCount: 19,
			MvpCount:      3,
			UtilityDamage: 64,
		},
		{
			SteamID64:     76561198040691320,
			Name:          "jakoby",
			KillCount:     26,
			AssistCount:   11,
			DeathCount:    25,
			Score:         76,
			Team:          match.TeamB,
			HeadshotCount: 17,
			MvpCount:      4,
			UtilityDamage: 363,
		},
	}

	assertion.AssertPlayers(t, match, players)
	assertion.AssertRounds(t, match, rounds)
}
