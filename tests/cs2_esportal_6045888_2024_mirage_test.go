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

// https://esportal.com/en/match/6045888
func TestEsportal_6045888_2024_Mirage(t *testing.T) {
	demoName := "esportal_6045888_2024_mirage"
	demoPath := testsutils.GetDemoPath("cs2", demoName)
	match, err := api.AnalyzeDemo(demoPath, api.AnalyzeDemoOptions{
		Source: constants.DemoSourceEbot,
	})
	if err != nil {
		t.Error(err)
	}

	expectedRoundCount := 18
	expectedPlayerCount := 10
	expectedScoreTeamA := 13
	expectedScoreTeamB := 5
	expectedScoreFirstHalfTeamA := 7
	expectedScoreFirstHalfTeamB := 5
	expectedScoreSecondHalfTeamA := 6
	expectedScoreSecondHalfTeamB := 0
	expectedTeamNameA := "The Pioneers"
	expectedTeamNameB := "The Explorers"
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
			StartTick:         4692,
			StartFrame:        4917,
			EndTick:           9050,
			FreezeTimeEndTick: 5650,
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
			StartTick:         9498,
			StartFrame:        9829,
			EndTick:           13355,
			FreezeTimeEndTick: 10458,
			TeamAStartMoney:   19150,
			TeamBStartMoney:   10800,
			TeamAEconomyType:  constants.EconomyTypeSemi,
			TeamBEconomyType:  constants.EconomyTypeEco,
			EndReason:         events.RoundEndReasonTerroristsWin,
			WinnerSide:        common.TeamTerrorists,
			TeamAScore:        1,
			TeamBScore:        1,
			WinnerName:        expectedTeamNameB,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            3,
			StartTick:         17643,
			StartFrame:        18256,
			EndTick:           21908,
			FreezeTimeEndTick: 18603,
			TeamAStartMoney:   11750,
			TeamBStartMoney:   24850,
			TeamAEconomyType:  constants.EconomyTypeSemi,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonTerroristsWin,
			WinnerSide:        common.TeamTerrorists,
			TeamAScore:        1,
			TeamBScore:        2,
			WinnerName:        expectedTeamNameB,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            4,
			StartTick:         22356,
			StartFrame:        23224,
			EndTick:           26148,
			FreezeTimeEndTick: 23316,
			TeamAStartMoney:   15300,
			TeamBStartMoney:   22650,
			TeamAEconomyType:  constants.EconomyTypeSemi,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonTerroristsWin,
			WinnerSide:        common.TeamTerrorists,
			TeamAScore:        1,
			TeamBScore:        3,
			WinnerName:        expectedTeamNameB,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            5,
			StartTick:         26596,
			StartFrame:        27698,
			EndTick:           30199,
			FreezeTimeEndTick: 27556,
			TeamAStartMoney:   21250,
			TeamBStartMoney:   25150,
			TeamAEconomyType:  constants.EconomyTypeForceBuy,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonTerroristsWin,
			WinnerSide:        common.TeamTerrorists,
			TeamAScore:        1,
			TeamBScore:        4,
			WinnerName:        expectedTeamNameB,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            6,
			StartTick:         30647,
			StartFrame:        31889,
			EndTick:           35831,
			FreezeTimeEndTick: 31607,
			TeamAStartMoney:   15300,
			TeamBStartMoney:   30250,
			TeamAEconomyType:  constants.EconomyTypeSemi,
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
			StartTick:         36279,
			StartFrame:        37914,
			EndTick:           39879,
			FreezeTimeEndTick: 37239,
			TeamAStartMoney:   27550,
			TeamBStartMoney:   42550,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonCTWin,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        2,
			TeamBScore:        5,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            8,
			StartTick:         40327,
			StartFrame:        42021,
			EndTick:           49278,
			FreezeTimeEndTick: 41287,
			TeamAStartMoney:   18600,
			TeamBStartMoney:   27600,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonBombDefused,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        3,
			TeamBScore:        5,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            9,
			StartTick:         49726,
			StartFrame:        51734,
			EndTick:           55900,
			FreezeTimeEndTick: 50686,
			TeamAStartMoney:   23500,
			TeamBStartMoney:   18900,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeSemi,
			EndReason:         events.RoundEndReasonCTWin,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        4,
			TeamBScore:        5,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            10,
			StartTick:         56348,
			StartFrame:        58761,
			EndTick:           60735,
			FreezeTimeEndTick: 57308,
			TeamAStartMoney:   20350,
			TeamBStartMoney:   13450,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeSemi,
			EndReason:         events.RoundEndReasonCTWin,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        5,
			TeamBScore:        5,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            11,
			StartTick:         61183,
			StartFrame:        63926,
			EndTick:           65063,
			FreezeTimeEndTick: 62143,
			TeamAStartMoney:   21600,
			TeamBStartMoney:   22550,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonCTWin,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        6,
			TeamBScore:        5,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            12,
			StartTick:         65511,
			StartFrame:        68389,
			EndTick:           71057,
			FreezeTimeEndTick: 66471,
			TeamAStartMoney:   22850,
			TeamBStartMoney:   17700,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeForceBuy,
			EndReason:         events.RoundEndReasonCTWin,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        7,
			TeamBScore:        5,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            13,
			StartTick:         72017,
			StartFrame:        75203,
			EndTick:           74724,
			FreezeTimeEndTick: 72977,
			TeamAStartMoney:   4000,
			TeamBStartMoney:   4000,
			TeamAEconomyType:  constants.EconomyTypePistol,
			TeamBEconomyType:  constants.EconomyTypePistol,
			EndReason:         events.RoundEndReasonTerroristsWin,
			WinnerSide:        common.TeamTerrorists,
			TeamAScore:        8,
			TeamBScore:        5,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamTerrorists,
			TeamBSide:         common.TeamCounterTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            14,
			StartTick:         75172,
			StartFrame:        78388,
			EndTick:           81915,
			FreezeTimeEndTick: 76132,
			TeamAStartMoney:   18800,
			TeamBStartMoney:   10350,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeEco,
			EndReason:         events.RoundEndReasonTargetBombed,
			WinnerSide:        common.TeamTerrorists,
			TeamAScore:        9,
			TeamBScore:        5,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamTerrorists,
			TeamBSide:         common.TeamCounterTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            15,
			StartTick:         82363,
			StartFrame:        86059,
			EndTick:           86234,
			FreezeTimeEndTick: 83323,
			TeamAStartMoney:   20250,
			TeamBStartMoney:   19500,
			TeamAEconomyType:  constants.EconomyTypeSemi,
			TeamBEconomyType:  constants.EconomyTypeForceBuy,
			EndReason:         events.RoundEndReasonTerroristsWin,
			WinnerSide:        common.TeamTerrorists,
			TeamAScore:        10,
			TeamBScore:        5,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamTerrorists,
			TeamBSide:         common.TeamCounterTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            16,
			StartTick:         86682,
			StartFrame:        90532,
			EndTick:           92594,
			FreezeTimeEndTick: 87642,
			TeamAStartMoney:   19600,
			TeamBStartMoney:   15450,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeSemi,
			EndReason:         events.RoundEndReasonTerroristsWin,
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
			StartTick:         93042,
			StartFrame:        97330,
			EndTick:           95945,
			FreezeTimeEndTick: 94002,
			TeamAStartMoney:   29450,
			TeamBStartMoney:   26050,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeFull,
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
			StartTick:         96393,
			StartFrame:        100683,
			EndTick:           98772,
			FreezeTimeEndTick: 97353,
			TeamAStartMoney:   28950,
			TeamBStartMoney:   18850,
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
			SteamID64:     76561198877382043,
			Name:          "carokann",
			KillCount:     11,
			AssistCount:   4,
			DeathCount:    15,
			Score:         28,
			Team:          match.TeamB,
			HeadshotCount: 8,
			MvpCount:      1,
			UtilityDamage: 86,
		},
		{
			SteamID64:     76561198370755882,
			Name:          "R1C",
			KillCount:     14,
			AssistCount:   2,
			DeathCount:    16,
			Score:         32,
			Team:          match.TeamB,
			HeadshotCount: 10,
			MvpCount:      2,
			UtilityDamage: 30,
		},
		{
			SteamID64:     76561198347044144,
			Name:          "RobdudR",
			KillCount:     14,
			AssistCount:   4,
			DeathCount:    16,
			Score:         32,
			Team:          match.TeamB,
			HeadshotCount: 6,
			MvpCount:      1,
			UtilityDamage: 60,
		},
		{
			SteamID64:     76561198057057588,
			Name:          "Dahlberg",
			KillCount:     11, // 11 instead of 12 round 17 TK carokann?
			AssistCount:   2,
			DeathCount:    14, // 14 insteaf of 13 because of TK?
			Score:         28,
			Team:          match.TeamB,
			HeadshotCount: 9,
			MvpCount:      1,
			UtilityDamage: 28,
		},
		{
			SteamID64:     76561199220179623,
			Name:          "svag",
			KillCount:     3,
			AssistCount:   2,
			DeathCount:    16,
			Score:         11,
			Team:          match.TeamB,
			HeadshotCount: 1,
			MvpCount:      0,
			UtilityDamage: 44,
		},
		{
			SteamID64:     76561199544615952,
			Name:          "tikitak",
			KillCount:     9,
			AssistCount:   9,
			DeathCount:    11,
			Score:         35,
			Team:          match.TeamA,
			HeadshotCount: 5,
			MvpCount:      4,
			UtilityDamage: 4,
		},
		{
			SteamID64:     76561198861806199,
			Name:          "Livsblatte",
			KillCount:     12,
			AssistCount:   4,
			DeathCount:    8,
			Score:         28,
			Team:          match.TeamA,
			HeadshotCount: 9,
			MvpCount:      2,
			UtilityDamage: 9,
		},
		{
			SteamID64:     76561199443155864,
			Name:          "fukushima",
			KillCount:     20,
			AssistCount:   3,
			DeathCount:    13,
			Score:         43,
			Team:          match.TeamA,
			HeadshotCount: 8,
			MvpCount:      3,
			UtilityDamage: 93,
		},
		{
			SteamID64:     76561198988232934,
			Name:          "zyyx",
			KillCount:     21,
			AssistCount:   4,
			DeathCount:    12,
			Score:         49,
			Team:          match.TeamA,
			HeadshotCount: 5,
			MvpCount:      2,
			UtilityDamage: 15,
		},
		{
			SteamID64:     76561199068347238,
			Name:          "incog",
			KillCount:     13,
			AssistCount:   3,
			DeathCount:    10,
			Score:         34,
			Team:          match.TeamA,
			HeadshotCount: 9,
			MvpCount:      2,
			UtilityDamage: 3,
		},
	}

	assertion.AssertPlayers(t, match, players)
	assertion.AssertRounds(t, match, rounds)
}
