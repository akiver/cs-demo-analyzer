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

// 2v2
// https://esplay.com/m/nnQccdWWJtkc/midnight-vs-haha123xd
func TestEsplay_nnQccdWWJtkc_2025_Vertigo(t *testing.T) {
	demoName := "esplay_nnQccdWWJtkc_2025_vertigo"
	demoPath := testsutils.GetDemoPath("cs2", demoName)
	match, err := api.AnalyzeDemo(demoPath, api.AnalyzeDemoOptions{
		Source: constants.DemoSourceEsplay,
	})

	if err != nil {
		t.Error(err)
	}

	expectedRoundCount := 14
	expectedPlayerCount := 4
	expectedScoreTeamA := 5
	expectedScoreTeamB := 9
	expectedScoreFirstHalfTeamA := 4
	expectedScoreFirstHalfTeamB := 4
	expectedScoreSecondHalfTeamA := 1
	expectedScoreSecondHalfTeamB := 5
	expectedTeamNameA := "haha123xd"
	expectedTeamNameB := "Midnight"
	expectedWinnerName := expectedTeamNameB
	expectedMaxRounds := 16

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
			StartFrame:        422,
			EndTick:           2850,
			FreezeTimeEndTick: 511,
			TeamAStartMoney:   1600,
			TeamBStartMoney:   1600,
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
			StartTick:         3298,
			StartFrame:        3661,
			EndTick:           9316,
			FreezeTimeEndTick: 3810,
			TeamAStartMoney:   7400,
			TeamBStartMoney:   4100,
			TeamAEconomyType:  constants.EconomyTypeSemi,
			TeamBEconomyType:  constants.EconomyTypeEco,
			EndReason:         events.RoundEndReasonBombDefused,
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
			StartTick:         9764,
			StartFrame:        10677,
			EndTick:           11710,
			FreezeTimeEndTick: 10276,
			TeamAStartMoney:   8400,
			TeamBStartMoney:   9300,
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
			StartTick:         12158,
			StartFrame:        13327,
			EndTick:           15540,
			FreezeTimeEndTick: 12670,
			TeamAStartMoney:   5600,
			TeamBStartMoney:   8000,
			TeamAEconomyType:  constants.EconomyTypeForceBuy,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonCTWin,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        3,
			TeamBScore:        1,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            5,
			StartTick:         15988,
			StartFrame:        17460,
			EndTick:           18502,
			FreezeTimeEndTick: 16500,
			TeamAStartMoney:   7150,
			TeamBStartMoney:   6750,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeForceBuy,
			EndReason:         events.RoundEndReasonTerroristsWin,
			WinnerSide:        common.TeamTerrorists,
			TeamAScore:        3,
			TeamBScore:        2,
			WinnerName:        expectedTeamNameB,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            6,
			StartTick:         18950,
			StartFrame:        20673,
			EndTick:           23626,
			FreezeTimeEndTick: 19462,
			TeamAStartMoney:   3350,
			TeamBStartMoney:   7900,
			TeamAEconomyType:  constants.EconomyTypeForceBuy,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonTerroristsWin,
			WinnerSide:        common.TeamTerrorists,
			TeamAScore:        3,
			TeamBScore:        3,
			WinnerName:        expectedTeamNameB,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            7,
			StartTick:         24074,
			StartFrame:        26179,
			EndTick:           25448,
			FreezeTimeEndTick: 24586,
			TeamAStartMoney:   4000,
			TeamBStartMoney:   13650,
			TeamAEconomyType:  constants.EconomyTypeEco,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonTerroristsWin,
			WinnerSide:        common.TeamTerrorists,
			TeamAScore:        3,
			TeamBScore:        4,
			WinnerName:        expectedTeamNameB,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            8,
			StartTick:         25896,
			StartFrame:        28116,
			EndTick:           32444,
			FreezeTimeEndTick: 26408,
			TeamAStartMoney:   8300,
			TeamBStartMoney:   19500,
			TeamAEconomyType:  constants.EconomyTypeForceBuy,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonBombDefused,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        4,
			TeamBScore:        4,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            9,
			StartTick:         33404,
			StartFrame:        36279,
			EndTick:           35256,
			FreezeTimeEndTick: 33916,
			TeamAStartMoney:   1600,
			TeamBStartMoney:   1600,
			TeamAEconomyType:  constants.EconomyTypePistol,
			TeamBEconomyType:  constants.EconomyTypePistol,
			EndReason:         events.RoundEndReasonCTWin,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        4,
			TeamBScore:        5,
			WinnerName:        expectedTeamNameB,
			TeamASide:         common.TeamTerrorists,
			TeamBSide:         common.TeamCounterTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            10,
			StartTick:         35704,
			StartFrame:        38714,
			EndTick:           37867,
			FreezeTimeEndTick: 36216,
			TeamAStartMoney:   4100,
			TeamBStartMoney:   7400,
			TeamAEconomyType:  constants.EconomyTypeEco,
			TeamBEconomyType:  constants.EconomyTypeSemi,
			EndReason:         events.RoundEndReasonCTWin,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        4,
			TeamBScore:        6,
			WinnerName:        expectedTeamNameB,
			TeamASide:         common.TeamTerrorists,
			TeamBSide:         common.TeamCounterTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            11,
			StartTick:         38315,
			StartFrame:        41597,
			EndTick:           41868,
			FreezeTimeEndTick: 38827,
			TeamAStartMoney:   7600,
			TeamBStartMoney:   8350,
			TeamAEconomyType:  constants.EconomyTypeForceBuy,
			TeamBEconomyType:  constants.EconomyTypeSemi,
			EndReason:         events.RoundEndReasonTerroristsWin,
			WinnerSide:        common.TeamTerrorists,
			TeamAScore:        5,
			TeamBScore:        6,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamTerrorists,
			TeamBSide:         common.TeamCounterTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            12,
			StartTick:         42316,
			StartFrame:        45947,
			EndTick:           44937,
			FreezeTimeEndTick: 42828,
			TeamAStartMoney:   7300,
			TeamBStartMoney:   10350,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonCTWin,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        5,
			TeamBScore:        7,
			WinnerName:        expectedTeamNameB,
			TeamASide:         common.TeamTerrorists,
			TeamBSide:         common.TeamCounterTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            13,
			StartTick:         45385,
			StartFrame:        49149,
			EndTick:           53118,
			FreezeTimeEndTick: 45897,
			TeamAStartMoney:   5300,
			TeamBStartMoney:   7800,
			TeamAEconomyType:  constants.EconomyTypeEco,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonBombDefused,
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
			StartTick:         53566,
			StartFrame:        57882,
			EndTick:           58836,
			FreezeTimeEndTick: 54078,
			TeamAStartMoney:   7000,
			TeamBStartMoney:   14200,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeFull,
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
	}

	var players = []fake.FakePlayer{
		{
			SteamID64:     76561198795651874,
			Name:          "Emmy",
			KillCount:     13,
			AssistCount:   2,
			DeathCount:    5,
			Score:         35,
			Team:          match.TeamB,
			HeadshotCount: 10,
			MvpCount:      5,
			UtilityDamage: 9,
		},
		{
			SteamID64:     76561198089429734,
			Name:          "Lowlita",
			KillCount:     8,
			AssistCount:   3,
			DeathCount:    9,
			Score:         22,
			Team:          match.TeamB,
			HeadshotCount: 4,
			MvpCount:      4,
			UtilityDamage: 65,
		},
		{
			SteamID64:     76561199087595827,
			Name:          "stina",
			KillCount:     9,
			AssistCount:   1,
			DeathCount:    10,
			Score:         24,
			Team:          match.TeamA,
			HeadshotCount: 6,
			MvpCount:      3,
			UtilityDamage: 4,
		},
		{
			SteamID64:     76561198117702148,
			Name:          "zietra",
			KillCount:     5,
			AssistCount:   1,
			DeathCount:    11,
			Score:         13,
			Team:          match.TeamA,
			HeadshotCount: 1,
			MvpCount:      2,
			UtilityDamage: 51,
		},
	}

	assertion.AssertPlayers(t, match, players)
	assertion.AssertRounds(t, match, rounds)
}
