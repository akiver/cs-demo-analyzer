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

// Contains a knife round
func Test5EPlay_5eplay_g161_20231231135244670959707_2023_Mirage(t *testing.T) {
	demoName := "5eplay_g161_20231231135244670959707_2023_mirage"
	demoPath := testsutils.GetDemoPath("cs2", demoName)
	match, err := api.AnalyzeDemo(demoPath, api.AnalyzeDemoOptions{
		Source: constants.DemoSourceFiveEPlay,
	})
	if err != nil {
		t.Error(err)
	}

	expectedRoundCount := 16
	expectedPlayerCount := 10
	expectedScoreTeamA := 13
	expectedScoreTeamB := 3
	expectedScoreFirstHalfTeamA := 11
	expectedScoreFirstHalfTeamB := 1
	expectedScoreSecondHalfTeamA := 2
	expectedScoreSecondHalfTeamB := 2
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
			StartTick:         8588,
			StartFrame:        8622,
			EndTick:           11093,
			FreezeTimeEndTick: 9356,
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
			StartTick:         11541,
			StartFrame:        11789,
			EndTick:           14905,
			FreezeTimeEndTick: 12309,
			TeamAStartMoney:   18500,
			TeamBStartMoney:   11500,
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
			StartTick:         15353,
			StartFrame:        15823,
			EndTick:           23633,
			FreezeTimeEndTick: 16121,
			TeamAStartMoney:   23500,
			TeamBStartMoney:   20100,
			TeamAEconomyType:  constants.EconomyTypeSemi,
			TeamBEconomyType:  constants.EconomyTypeSemi,
			EndReason:         events.RoundEndReasonBombDefused,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        3,
			TeamBScore:        0,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            4,
			StartTick:         24081,
			StartFrame:        24842,
			EndTick:           27922,
			FreezeTimeEndTick: 24849,
			TeamAStartMoney:   34800,
			TeamBStartMoney:   23750,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonBombDefused,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        4,
			TeamBScore:        0,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            5,
			StartTick:         28370,
			StartFrame:        29427,
			EndTick:           33385,
			FreezeTimeEndTick: 29138,
			TeamAStartMoney:   35650,
			TeamBStartMoney:   25850,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonTerroristsWin,
			WinnerSide:        common.TeamTerrorists,
			TeamAScore:        4,
			TeamBScore:        1,
			WinnerName:        expectedTeamNameB,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            6,
			StartTick:         33833,
			StartFrame:        35289,
			EndTick:           36625,
			FreezeTimeEndTick: 34601,
			TeamAStartMoney:   30050,
			TeamBStartMoney:   22400,
			TeamAEconomyType:  constants.EconomyTypeSemi,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonCTWin,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        5,
			TeamBScore:        1,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            7,
			StartTick:         37073,
			StartFrame:        38833,
			EndTick:           44081,
			FreezeTimeEndTick: 37841,
			TeamAStartMoney:   29450,
			TeamBStartMoney:   30350,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeSemi,
			EndReason:         events.RoundEndReasonCTWin,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        6,
			TeamBScore:        1,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            8,
			StartTick:         44529,
			StartFrame:        46591,
			EndTick:           51350,
			FreezeTimeEndTick: 45297,
			TeamAStartMoney:   32000,
			TeamBStartMoney:   28750,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonCTWin,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        7,
			TeamBScore:        1,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            9,
			StartTick:         51798,
			StartFrame:        54177,
			EndTick:           54526,
			FreezeTimeEndTick: 52566,
			TeamAStartMoney:   28000,
			TeamBStartMoney:   23800,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeSemi,
			EndReason:         events.RoundEndReasonCTWin,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        8,
			TeamBScore:        1,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            10,
			StartTick:         54974,
			StartFrame:        57640,
			EndTick:           59164,
			FreezeTimeEndTick: 55742,
			TeamAStartMoney:   34600,
			TeamBStartMoney:   23400,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonCTWin,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        9,
			TeamBScore:        1,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            11,
			StartTick:         59612,
			StartFrame:        62550,
			EndTick:           63997,
			FreezeTimeEndTick: 60380,
			TeamAStartMoney:   46850,
			TeamBStartMoney:   21850,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonCTWin,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        10,
			TeamBScore:        1,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            12,
			StartTick:         64445,
			StartFrame:        67667,
			EndTick:           68874,
			FreezeTimeEndTick: 65213,
			TeamAStartMoney:   47650,
			TeamBStartMoney:   20250,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonCTWin,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        11,
			TeamBScore:        1,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            13,
			StartTick:         69418,
			StartFrame:        73009,
			EndTick:           73906,
			FreezeTimeEndTick: 70602,
			TeamAStartMoney:   4000,
			TeamBStartMoney:   4000,
			TeamAEconomyType:  constants.EconomyTypePistol,
			TeamBEconomyType:  constants.EconomyTypePistol,
			EndReason:         events.RoundEndReasonTerroristsWin,
			WinnerSide:        common.TeamTerrorists,
			TeamAScore:        12,
			TeamBScore:        1,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamTerrorists,
			TeamBSide:         common.TeamCounterTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            14,
			StartTick:         74354,
			StartFrame:        78205,
			EndTick:           80177,
			FreezeTimeEndTick: 75122,
			TeamAStartMoney:   18650,
			TeamBStartMoney:   12750,
			TeamAEconomyType:  constants.EconomyTypeSemi,
			TeamBEconomyType:  constants.EconomyTypeForceBuy,
			EndReason:         events.RoundEndReasonCTWin,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        12,
			TeamBScore:        2,
			WinnerName:        expectedTeamNameB,
			TeamASide:         common.TeamTerrorists,
			TeamBSide:         common.TeamCounterTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            15,
			StartTick:         80625,
			StartFrame:        84798,
			EndTick:           85853,
			FreezeTimeEndTick: 81393,
			TeamAStartMoney:   9000,
			TeamBStartMoney:   18800,
			TeamAEconomyType:  constants.EconomyTypeSemi,
			TeamBEconomyType:  constants.EconomyTypeSemi,
			EndReason:         events.RoundEndReasonBombDefused,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        12,
			TeamBScore:        3,
			WinnerName:        expectedTeamNameB,
			TeamASide:         common.TeamTerrorists,
			TeamBSide:         common.TeamCounterTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            16,
			StartTick:         86301,
			StartFrame:        90829,
			EndTick:           89900,
			FreezeTimeEndTick: 87069,
			TeamAStartMoney:   21000,
			TeamBStartMoney:   7300,
			TeamAEconomyType:  constants.EconomyTypeSemi,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonTerroristsWin,
			WinnerSide:        common.TeamTerrorists,
			TeamAScore:        13,
			TeamBScore:        3,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamTerrorists,
			TeamBSide:         common.TeamCounterTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
	}

	var players = []fake.FakePlayer{
		{
			SteamID64:     76561198836396113,
			Name:          "钢镚ㇽ。",
			KillCount:     10,
			AssistCount:   2,
			DeathCount:    15,
			Score:         22,
			Team:          match.TeamB,
			HeadshotCount: 4,
			MvpCount:      0,
			UtilityDamage: 58,
		},
		{
			SteamID64:     76561199104520853,
			Name:          "开什么枪拿刀",
			KillCount:     14,
			AssistCount:   6,
			DeathCount:    10,
			Score:         35,
			Team:          match.TeamA,
			HeadshotCount: 5,
			MvpCount:      2,
			UtilityDamage: 5,
		},
		{
			SteamID64:     76561199064714350,
			Name:          "KG.70",
			KillCount:     13,
			AssistCount:   4,
			DeathCount:    14,
			Score:         32,
			Team:          match.TeamB,
			HeadshotCount: 5,
			MvpCount:      1,
			UtilityDamage: 0,
		},
		{
			SteamID64:     76561198840681151,
			Name:          "Dadddddy",
			KillCount:     23,
			AssistCount:   4,
			DeathCount:    8,
			Score:         59,
			Team:          match.TeamA,
			HeadshotCount: 10,
			MvpCount:      4,
			UtilityDamage: 24,
		},
		{
			SteamID64:     76561198886914468,
			Name:          "和我一起去逃避",
			KillCount:     10,
			AssistCount:   3,
			DeathCount:    9,
			Score:         25,
			Team:          match.TeamA,
			HeadshotCount: 2,
			MvpCount:      2,
			UtilityDamage: 92,
		},
		{
			SteamID64:     76561198872277112,
			Name:          "皮皮猪的小宝贝",
			KillCount:     8,
			AssistCount:   4,
			DeathCount:    14,
			Score:         23,
			Team:          match.TeamB,
			HeadshotCount: 3,
			MvpCount:      1,
			UtilityDamage: 35,
		},
		{
			SteamID64:     76561199238440889,
			Name:          "哥的情绪稳定",
			KillCount:     12,
			AssistCount:   2,
			DeathCount:    14,
			Score:         27,
			Team:          match.TeamB,
			HeadshotCount: 11,
			MvpCount:      0,
			UtilityDamage: 0,
		},
		{
			SteamID64:     76561198994608392,
			Name:          "莫忘风",
			KillCount:     9,
			AssistCount:   3,
			DeathCount:    12,
			Score:         23,
			Team:          match.TeamA,
			HeadshotCount: 6,
			MvpCount:      2,
			UtilityDamage: 103,
		},
		{
			SteamID64:     76561198288657828,
			Name:          "zHunBi_wOw",
			KillCount:     14,
			AssistCount:   5,
			DeathCount:    8,
			Score:         36,
			Team:          match.TeamA,
			HeadshotCount: 4,
			MvpCount:      3,
			UtilityDamage: 0,
		},
		{
			SteamID64:     76561199194343995,
			Name:          "二次元小陈",
			KillCount:     4,
			AssistCount:   3,
			DeathCount:    13,
			Score:         14,
			Team:          match.TeamB,
			HeadshotCount: 2,
			MvpCount:      1,
			UtilityDamage: 55,
		},
	}

	assertion.AssertPlayers(t, match, players)
	assertion.AssertRounds(t, match, rounds)
}
