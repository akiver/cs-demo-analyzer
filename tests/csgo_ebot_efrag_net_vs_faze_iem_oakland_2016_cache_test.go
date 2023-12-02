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

// E-frag starts as T, won the knife round and switch to CT.
// There is then a warmup and the first round is cancelled during freezetime (!stop command).
// The first round really starts after a restart by the eBot.
// At the beginning of the second round, a pause is requested. The game is paused and the 2nd round starts at tick 74004.
// https://www.hltv.org/stats/matches/mapstatsid/28406/e-fragnet-vs-faze
func TestEbot_Efrag_Net_VS_Faze_IEM_Oakland_2016_Cache(t *testing.T) {
	demoName := "ebot_efrag_net_vs_faze_iem_oakland_2016_cache"
	demoPath := testsutils.GetDemoPath("csgo", demoName)
	match, err := api.AnalyzeDemo(demoPath, api.AnalyzeDemoOptions{
		Source: constants.DemoSourceEbot,
	})
	if err != nil {
		t.Error(err)
	}

	expectedRoundCount := 21
	expectedPlayerCount := 10
	expectedScoreTeamA := 5
	expectedScoreTeamB := 16
	expectedScoreFirstHalfTeamA := 4
	expectedScoreFirstHalfTeamB := 11
	expectedScoreSecondHalfTeamA := 1
	expectedScoreSecondHalfTeamB := 5
	expectedTeamNameA := "E-Frag.net"
	expectedTeamNameB := "FaZe"
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
			StartTick:         17434,
			StartFrame:        17407,
			EndTick:           23412,
			FreezeTimeEndTick: 19349,
			EndOfficiallyTick: 24052,
			TeamAStartMoney:   4000,
			TeamBStartMoney:   4000,
			TeamAEconomyType:  constants.EconomyTypePistol,
			TeamBEconomyType:  constants.EconomyTypePistol,
			EndReason:         events.RoundEndReasonCTWin,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        1,
			TeamBScore:        0,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            2,
			StartTick:         74004,
			StartFrame:        73957,
			EndTick:           89035,
			FreezeTimeEndTick: 77478,
			EndOfficiallyTick: 89677,
			TeamAStartMoney:   18450,
			TeamBStartMoney:   8150,
			TeamAEconomyType:  constants.EconomyTypeSemi,
			TeamBEconomyType:  constants.EconomyTypeForceBuy,
			EndReason:         events.RoundEndReasonTerroristsWin,
			WinnerSide:        common.TeamTerrorists,
			TeamAScore:        1,
			TeamBScore:        1,
			WinnerName:        expectedTeamNameB,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
		},
		{
			Number:            3,
			StartTick:         89677,
			StartFrame:        89616,
			EndTick:           97974,
			FreezeTimeEndTick: 91597,
			EndOfficiallyTick: 98616,
			TeamAStartMoney:   9650,
			TeamBStartMoney:   18350,
			TeamAEconomyType:  constants.EconomyTypeForceBuy,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonTerroristsWin,
			WinnerSide:        common.TeamTerrorists,
			TeamAScore:        1,
			TeamBScore:        2,
			WinnerName:        expectedTeamNameB,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
		},
		{
			Number:            4,
			StartTick:         98616,
			StartFrame:        98549,
			EndTick:           104255,
			FreezeTimeEndTick: 100536,
			EndOfficiallyTick: 104897,
			TeamAStartMoney:   11100,
			TeamBStartMoney:   19000,
			TeamAEconomyType:  constants.EconomyTypeEco,
			TeamBEconomyType:  constants.EconomyTypeSemi,
			EndReason:         events.RoundEndReasonTerroristsWin,
			WinnerSide:        common.TeamTerrorists,
			TeamAScore:        1,
			TeamBScore:        3,
			WinnerName:        expectedTeamNameB,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
		},
		{
			Number:            5,
			StartTick:         104897,
			StartFrame:        104817,
			EndTick:           123173,
			FreezeTimeEndTick: 106817,
			EndOfficiallyTick: 123813,
			TeamAStartMoney:   22600,
			TeamBStartMoney:   24650,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonBombDefused,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        2,
			TeamBScore:        3,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
		},
		{
			Number:            6,
			StartTick:         123813,
			StartFrame:        123723,
			EndTick:           129067,
			FreezeTimeEndTick: 125733,
			EndOfficiallyTick: 129707,
			TeamAStartMoney:   19650,
			TeamBStartMoney:   26750,
			TeamAEconomyType:  constants.EconomyTypeFull,
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
			StartTick:         129707,
			StartFrame:        129608,
			EndTick:           144873,
			FreezeTimeEndTick: 131627,
			EndOfficiallyTick: 145513,
			TeamAStartMoney:   13450,
			TeamBStartMoney:   21200,
			TeamAEconomyType:  constants.EconomyTypeForceBuy,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonTargetBombed,
			WinnerSide:        common.TeamTerrorists,
			TeamAScore:        2,
			TeamBScore:        5,
			WinnerName:        expectedTeamNameB,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
		},
		{
			Number:            8,
			StartTick:         145513,
			StartFrame:        145409,
			EndTick:           157580,
			FreezeTimeEndTick: 147433,
			EndOfficiallyTick: 158221,
			TeamAStartMoney:   10550,
			TeamBStartMoney:   37800,
			TeamAEconomyType:  constants.EconomyTypeSemi,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonTerroristsWin,
			WinnerSide:        common.TeamTerrorists,
			TeamAScore:        2,
			TeamBScore:        6,
			WinnerName:        expectedTeamNameB,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
		},
		{
			Number:            9,
			StartTick:         158221,
			StartFrame:        158103,
			EndTick:           167465,
			FreezeTimeEndTick: 160141,
			EndOfficiallyTick: 168106,
			TeamAStartMoney:   21450,
			TeamBStartMoney:   40950,
			TeamAEconomyType:  constants.EconomyTypeForceBuy,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonTerroristsWin,
			WinnerSide:        common.TeamTerrorists,
			TeamAScore:        2,
			TeamBScore:        7,
			WinnerName:        expectedTeamNameB,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
		},
		{
			Number:            10,
			StartTick:         168106,
			StartFrame:        167980,
			EndTick:           179634,
			FreezeTimeEndTick: 170026,
			EndOfficiallyTick: 180274,
			TeamAStartMoney:   15550,
			TeamBStartMoney:   55700,
			TeamAEconomyType:  constants.EconomyTypeEco,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonTargetBombed,
			WinnerSide:        common.TeamTerrorists,
			TeamAScore:        2,
			TeamBScore:        8,
			WinnerName:        expectedTeamNameB,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
		},
		{
			Number:            11,
			StartTick:         180274,
			StartFrame:        180142,
			EndTick:           192470,
			FreezeTimeEndTick: 182194,
			EndOfficiallyTick: 193110,
			TeamAStartMoney:   30050,
			TeamBStartMoney:   57150,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonTargetBombed,
			WinnerSide:        common.TeamTerrorists,
			TeamAScore:        2,
			TeamBScore:        9,
			WinnerName:        expectedTeamNameB,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
		},
		{
			Number:            12,
			StartTick:         193110,
			StartFrame:        192964,
			EndTick:           223235,
			FreezeTimeEndTick: 208678,
			EndOfficiallyTick: 223875,
			TeamAStartMoney:   19850,
			TeamBStartMoney:   63050,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonBombDefused,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        3,
			TeamBScore:        9,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
		},
		{
			Number:            13,
			StartTick:         223875,
			StartFrame:        223723,
			EndTick:           236227,
			FreezeTimeEndTick: 225795,
			EndOfficiallyTick: 236867,
			TeamAStartMoney:   21050,
			TeamBStartMoney:   61150,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonTerroristsWin,
			WinnerSide:        common.TeamTerrorists,
			TeamAScore:        3,
			TeamBScore:        10,
			WinnerName:        expectedTeamNameB,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
		},
		{
			Number:            14,
			StartTick:         236867,
			StartFrame:        236710,
			EndTick:           249313,
			FreezeTimeEndTick: 238787,
			EndOfficiallyTick: 249953,
			TeamAStartMoney:   10350,
			TeamBStartMoney:   51950,
			TeamAEconomyType:  constants.EconomyTypeForceBuy,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonTargetBombed,
			WinnerSide:        common.TeamTerrorists,
			TeamAScore:        3,
			TeamBScore:        11,
			WinnerName:        expectedTeamNameB,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
		},
		{
			Number:            15,
			StartTick:         249953,
			StartFrame:        249791,
			EndTick:           257372,
			FreezeTimeEndTick: 251873,
			EndOfficiallyTick: 259330,
			TeamAStartMoney:   10050,
			TeamBStartMoney:   52250,
			TeamAEconomyType:  constants.EconomyTypeForceBuy,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonCTWin,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        4,
			TeamBScore:        11,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
		},
		{
			Number:            16,
			StartTick:         259330,
			StartFrame:        259161,
			EndTick:           271045,
			FreezeTimeEndTick: 261250,
			EndOfficiallyTick: 271685,
			TeamAStartMoney:   4000,
			TeamBStartMoney:   4000,
			TeamAEconomyType:  constants.EconomyTypePistol,
			TeamBEconomyType:  constants.EconomyTypePistol,
			EndReason:         events.RoundEndReasonBombDefused,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        4,
			TeamBScore:        12,
			WinnerName:        expectedTeamNameB,
			TeamASide:         common.TeamTerrorists,
			TeamBSide:         common.TeamCounterTerrorists,
		},
		{
			Number:            17,
			StartTick:         271685,
			StartFrame:        271505,
			EndTick:           280346,
			FreezeTimeEndTick: 273605,
			EndOfficiallyTick: 280988,
			TeamAStartMoney:   12200,
			TeamBStartMoney:   19950,
			TeamAEconomyType:  constants.EconomyTypeForceBuy,
			TeamBEconomyType:  constants.EconomyTypeSemi,
			EndReason:         events.RoundEndReasonCTWin,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        4,
			TeamBScore:        13,
			WinnerName:        expectedTeamNameB,
			TeamASide:         common.TeamTerrorists,
			TeamBSide:         common.TeamCounterTerrorists,
		},
		{
			Number:            18,
			StartTick:         280988,
			StartFrame:        280804,
			EndTick:           294711,
			FreezeTimeEndTick: 282908,
			EndOfficiallyTick: 295353,
			TeamAStartMoney:   10650,
			TeamBStartMoney:   19700,
			TeamAEconomyType:  constants.EconomyTypeForceBuy,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonCTWin,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        4,
			TeamBScore:        14,
			WinnerName:        expectedTeamNameB,
			TeamASide:         common.TeamTerrorists,
			TeamBSide:         common.TeamCounterTerrorists,
		},
		{
			Number:            19,
			StartTick:         295353,
			StartFrame:        295160,
			EndTick:           302598,
			FreezeTimeEndTick: 297273,
			EndOfficiallyTick: 303239,
			TeamAStartMoney:   13550,
			TeamBStartMoney:   24050,
			TeamAEconomyType:  constants.EconomyTypeSemi,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonCTWin,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        4,
			TeamBScore:        15,
			WinnerName:        expectedTeamNameB,
			TeamASide:         common.TeamTerrorists,
			TeamBSide:         common.TeamCounterTerrorists,
		},
		{
			Number:            20,
			StartTick:         303239,
			StartFrame:        303039,
			EndTick:           312242,
			FreezeTimeEndTick: 305159,
			EndOfficiallyTick: 312884,
			TeamAStartMoney:   17150,
			TeamBStartMoney:   23700,
			TeamAEconomyType:  constants.EconomyTypeForceBuy,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonTerroristsWin,
			WinnerSide:        common.TeamTerrorists,
			TeamAScore:        5,
			TeamBScore:        15,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamTerrorists,
			TeamBSide:         common.TeamCounterTerrorists,
		},
		{
			Number:            21,
			StartTick:         312884,
			StartFrame:        312679,
			EndTick:           324268,
			FreezeTimeEndTick: 314804,
			EndOfficiallyTick: 324268,
			TeamAStartMoney:   18750,
			TeamBStartMoney:   20550,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeForceBuy,
			EndReason:         events.RoundEndReasonCTWin,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        5,
			TeamBScore:        16,
			WinnerName:        expectedTeamNameB,
			TeamASide:         common.TeamTerrorists,
			TeamBSide:         common.TeamCounterTerrorists,
		},
	}

	assertion.AssertRounds(t, match, rounds)
}
