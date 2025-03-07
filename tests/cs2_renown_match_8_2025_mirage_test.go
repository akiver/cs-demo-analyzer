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

// https://renown.gg/match/8
func TestRenown_Match_8_2025_Mirage(t *testing.T) {
	demoName := "renown_match_8_2025_mirage"
	demoPath := testsutils.GetDemoPath("cs2", demoName)
	match, err := api.AnalyzeDemo(demoPath, api.AnalyzeDemoOptions{
		Source: constants.DemoSourceRenown,
	})

	if err != nil {
		t.Error(err)
	}

	expectedRoundCount := 18
	expectedPlayerCount := 10
	expectedScoreTeamA := 13
	expectedScoreTeamB := 5
	expectedScoreFirstHalfTeamA := 8
	expectedScoreFirstHalfTeamB := 4
	expectedScoreSecondHalfTeamA := 5
	expectedScoreSecondHalfTeamB := 1
	expectedTeamNameA := "Team FjobX"
	expectedTeamNameB := "Team svennerN"
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
			StartTick:         92,
			StartFrame:        1582,
			EndTick:           4081,
			FreezeTimeEndTick: 1466,
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
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            2,
			StartTick:         4529,
			StartFrame:        7101,
			EndTick:           8736,
			FreezeTimeEndTick: 5489,
			TeamAStartMoney:   20050,
			TeamBStartMoney:   13850,
			TeamAEconomyType:  constants.EconomyTypeSemi,
			TeamBEconomyType:  constants.EconomyTypeEco,
			EndReason:         events.RoundEndReasonCTWin,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        2,
			TeamBScore:        0,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            3,
			StartTick:         9184,
			StartFrame:        12429,
			EndTick:           18382,
			FreezeTimeEndTick: 10144,
			TeamAStartMoney:   18250,
			TeamBStartMoney:   23000,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonTerroristsWin,
			WinnerSide:        common.TeamTerrorists,
			TeamAScore:        2,
			TeamBScore:        1,
			WinnerName:        expectedTeamNameB,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            4,
			StartTick:         18830,
			StartFrame:        24298,
			EndTick:           22593,
			FreezeTimeEndTick: 19790,
			TeamAStartMoney:   11500,
			TeamBStartMoney:   18650,
			TeamAEconomyType:  constants.EconomyTypeSemi,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonTerroristsWin,
			WinnerSide:        common.TeamTerrorists,
			TeamAScore:        2,
			TeamBScore:        2,
			WinnerName:        expectedTeamNameB,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            5,
			StartTick:         23041,
			StartFrame:        29531,
			EndTick:           25684,
			FreezeTimeEndTick: 24001,
			TeamAStartMoney:   15250,
			TeamBStartMoney:   21450,
			TeamAEconomyType:  constants.EconomyTypeEco,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonTerroristsWin,
			WinnerSide:        common.TeamTerrorists,
			TeamAScore:        2,
			TeamBScore:        3,
			WinnerName:        expectedTeamNameB,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            6,
			StartTick:         26132,
			StartFrame:        33439,
			EndTick:           30125,
			FreezeTimeEndTick: 27092,
			TeamAStartMoney:   24400,
			TeamBStartMoney:   28100,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonCTWin,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        3,
			TeamBScore:        3,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            7,
			StartTick:         30573,
			StartFrame:        38920,
			EndTick:           35251,
			FreezeTimeEndTick: 31533,
			TeamAStartMoney:   18450,
			TeamBStartMoney:   29050,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonCTWin,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        4,
			TeamBScore:        3,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            8,
			StartTick:         35699,
			StartFrame:        45303,
			EndTick:           39121,
			FreezeTimeEndTick: 36659,
			TeamAStartMoney:   22200,
			TeamBStartMoney:   14500,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeSemi,
			EndReason:         events.RoundEndReasonCTWin,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        5,
			TeamBScore:        3,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            9,
			StartTick:         39569,
			StartFrame:        49767,
			EndTick:           46358,
			FreezeTimeEndTick: 40529,
			TeamAStartMoney:   19600,
			TeamBStartMoney:   21350,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonCTWin,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        6,
			TeamBScore:        3,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            10,
			StartTick:         46806,
			StartFrame:        58389,
			EndTick:           52987,
			FreezeTimeEndTick: 47766,
			TeamAStartMoney:   33700,
			TeamBStartMoney:   15750,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeSemi,
			EndReason:         events.RoundEndReasonBombDefused,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        7,
			TeamBScore:        3,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            11,
			StartTick:         53435,
			StartFrame:        66258,
			EndTick:           60764,
			FreezeTimeEndTick: 54395,
			TeamAStartMoney:   35100,
			TeamBStartMoney:   31000,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonTerroristsWin,
			WinnerSide:        common.TeamTerrorists,
			TeamAScore:        7,
			TeamBScore:        4,
			WinnerName:        expectedTeamNameB,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            12,
			StartTick:         61212,
			StartFrame:        75524,
			EndTick:           67476,
			FreezeTimeEndTick: 62172,
			TeamAStartMoney:   27600,
			TeamBStartMoney:   22700,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonBombDefused,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        8,
			TeamBScore:        4,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            13,
			StartTick:         70812,
			StartFrame:        87357,
			EndTick:           75247,
			FreezeTimeEndTick: 71772,
			TeamAStartMoney:   4000,
			TeamBStartMoney:   4000,
			TeamAEconomyType:  constants.EconomyTypePistol,
			TeamBEconomyType:  constants.EconomyTypePistol,
			EndReason:         events.RoundEndReasonTerroristsWin,
			WinnerSide:        common.TeamTerrorists,
			TeamAScore:        9,
			TeamBScore:        4,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamTerrorists,
			TeamBSide:         common.TeamCounterTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            14,
			StartTick:         75695,
			StartFrame:        93181,
			EndTick:           80667,
			FreezeTimeEndTick: 76655,
			TeamAStartMoney:   18750,
			TeamBStartMoney:   11400,
			TeamAEconomyType:  constants.EconomyTypeSemi,
			TeamBEconomyType:  constants.EconomyTypeEco,
			EndReason:         events.RoundEndReasonTerroristsWin,
			WinnerSide:        common.TeamTerrorists,
			TeamAScore:        10,
			TeamBScore:        4,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamTerrorists,
			TeamBSide:         common.TeamCounterTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            15,
			StartTick:         81115,
			StartFrame:        99518,
			EndTick:           87698,
			FreezeTimeEndTick: 82075,
			TeamAStartMoney:   18400,
			TeamBStartMoney:   22300,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonCTWin,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        10,
			TeamBScore:        5,
			WinnerName:        expectedTeamNameB,
			TeamASide:         common.TeamTerrorists,
			TeamBSide:         common.TeamCounterTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            16,
			StartTick:         88146,
			StartFrame:        108048,
			EndTick:           95327,
			FreezeTimeEndTick: 89106,
			TeamAStartMoney:   14200,
			TeamBStartMoney:   17500,
			TeamAEconomyType:  constants.EconomyTypeForceBuy,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonTargetBombed,
			WinnerSide:        common.TeamTerrorists,
			TeamAScore:        11,
			TeamBScore:        5,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamTerrorists,
			TeamBSide:         common.TeamCounterTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            17,
			StartTick:         95775,
			StartFrame:        117512,
			EndTick:           100615,
			FreezeTimeEndTick: 96735,
			TeamAStartMoney:   20450,
			TeamBStartMoney:   11250,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeSemi,
			EndReason:         events.RoundEndReasonTerroristsWin,
			WinnerSide:        common.TeamTerrorists,
			TeamAScore:        12,
			TeamBScore:        5,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamTerrorists,
			TeamBSide:         common.TeamCounterTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            18,
			StartTick:         101063,
			StartFrame:        123488,
			EndTick:           104872,
			FreezeTimeEndTick: 102023,
			TeamAStartMoney:   23550,
			TeamBStartMoney:   15300,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeForceBuy,
			EndReason:         events.RoundEndReasonTerroristsWin,
			WinnerSide:        common.TeamTerrorists,
			TeamAScore:        13,
			TeamBScore:        5,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamTerrorists,
			TeamBSide:         common.TeamCounterTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
	}

	var players = []fake.FakePlayer{
		{
			SteamID64:     76561198029485456,
			Name:          "whatsnxt",
			KillCount:     16,
			AssistCount:   6,
			DeathCount:    12,
			Score:         42,
			Team:          match.TeamA,
			HeadshotCount: 10,
			MvpCount:      2,
			UtilityDamage: 101,
		},
		{
			SteamID64:     76561198037147623,
			Name:          "FjobX",
			KillCount:     15,
			AssistCount:   6,
			DeathCount:    10,
			Score:         41,
			Team:          match.TeamA,
			HeadshotCount: 8,
			MvpCount:      5,
			UtilityDamage: 169,
		},
		{
			SteamID64:     76561198047724925,
			Name:          "ERXN",
			KillCount:     14,
			AssistCount:   4,
			DeathCount:    14,
			Score:         32,
			Team:          match.TeamA,
			HeadshotCount: 8,
			MvpCount:      0,
			UtilityDamage: 51,
		},
		{
			SteamID64:     76561198107208516,
			Name:          "Vodka",
			KillCount:     14,
			AssistCount:   6,
			DeathCount:    12, // It's 13 on CS2 and renown because this player had a disconnection
			Score:         34,
			Team:          match.TeamA,
			HeadshotCount: 4,
			MvpCount:      0,
			UtilityDamage: 121,
		},
		{
			SteamID64:     76561198024687209,
			Name:          "sWan_",
			KillCount:     19,
			AssistCount:   5,
			DeathCount:    10,
			Score:         49,
			Team:          match.TeamA,
			HeadshotCount: 6,
			MvpCount:      6,
			UtilityDamage: 210,
		},
		{
			SteamID64:     76561198050649118,
			Name:          "Zach",
			KillCount:     8,
			AssistCount:   4,
			DeathCount:    16,
			Score:         22,
			Team:          match.TeamB,
			HeadshotCount: 4,
			MvpCount:      1,
			UtilityDamage: 33,
		},
		{
			SteamID64:     76561198225661896,
			Name:          "Vit0-1337",
			KillCount:     24,
			AssistCount:   1,
			DeathCount:    13,
			Score:         51,
			Team:          match.TeamB,
			HeadshotCount: 12,
			MvpCount:      2,
			UtilityDamage: 110,
		},
		{
			SteamID64:     76561198028780484,
			Name:          "Hauman",
			KillCount:     7,
			AssistCount:   4,
			DeathCount:    16,
			Score:         20,
			Team:          match.TeamB,
			HeadshotCount: 5,
			MvpCount:      0,
			UtilityDamage: 52,
		},
		{
			SteamID64:     76561198087669880,
			Name:          "andy",
			KillCount:     5,
			AssistCount:   5,
			DeathCount:    16,
			Score:         15,
			Team:          match.TeamB,
			HeadshotCount: 2,
			MvpCount:      1,
			UtilityDamage: 0,
		},
		{
			SteamID64:     76561198063438674,
			Name:          "svennerN",
			KillCount:     14,
			AssistCount:   2,
			DeathCount:    17,
			Score:         31,
			Team:          match.TeamB,
			HeadshotCount: 8,
			MvpCount:      1,
			UtilityDamage: 48,
		},
	}

	assertion.AssertPlayers(t, match, players)
	assertion.AssertRounds(t, match, rounds)
}
