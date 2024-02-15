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

// https://esportal.com/en/match/6008132
func TestEsportal_6008132_2023_Mirage(t *testing.T) {
	demoName := "esportal_6008132_2023_mirage"
	demoPath := testsutils.GetDemoPath("cs2", demoName)
	match, err := api.AnalyzeDemo(demoPath, api.AnalyzeDemoOptions{
		Source: constants.DemoSourceEbot,
	})
	if err != nil {
		t.Error(err)
	}

	expectedRoundCount := 14
	expectedPlayerCount := 10
	expectedScoreTeamA := 13
	expectedScoreTeamB := 1
	expectedScoreFirstHalfTeamA := 11
	expectedScoreFirstHalfTeamB := 1
	expectedScoreSecondHalfTeamA := 2
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
			StartTick:         6811,
			StartFrame:        7216,
			EndTick:           9901,
			FreezeTimeEndTick: 7769,
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
			StartTick:         10349,
			StartFrame:        10943,
			EndTick:           13805,
			FreezeTimeEndTick: 11309,
			TeamAStartMoney:   18400,
			TeamBStartMoney:   10200,
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
			StartTick:         14253,
			StartFrame:        15117,
			EndTick:           18065,
			FreezeTimeEndTick: 15213,
			TeamAStartMoney:   20500,
			TeamBStartMoney:   17700,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeForceBuy,
			EndReason:         events.RoundEndReasonCTWin,
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
			StartTick:         18513,
			StartFrame:        19588,
			EndTick:           23234,
			FreezeTimeEndTick: 19473,
			TeamAStartMoney:   30400,
			TeamBStartMoney:   16050,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeForceBuy,
			EndReason:         events.RoundEndReasonCTWin,
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
			StartTick:         23682,
			StartFrame:        24917,
			EndTick:           29696,
			FreezeTimeEndTick: 24642,
			TeamAStartMoney:   33100,
			TeamBStartMoney:   18800,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeForceBuy,
			EndReason:         events.RoundEndReasonCTWin,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        5,
			TeamBScore:        0,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            6,
			StartTick:         30144,
			StartFrame:        31751,
			EndTick:           34861,
			FreezeTimeEndTick: 31104,
			TeamAStartMoney:   25500,
			TeamBStartMoney:   17400,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeForceBuy,
			EndReason:         events.RoundEndReasonCTWin,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        6,
			TeamBScore:        0,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            7,
			StartTick:         35309,
			StartFrame:        37046,
			EndTick:           43629,
			FreezeTimeEndTick: 36269,
			TeamAStartMoney:   32250,
			TeamBStartMoney:   18850,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeForceBuy,
			EndReason:         events.RoundEndReasonTargetSaved,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        7,
			TeamBScore:        0,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            8,
			StartTick:         44077,
			StartFrame:        46322,
			EndTick:           46998,
			FreezeTimeEndTick: 45037,
			TeamAStartMoney:   34200,
			TeamBStartMoney:   15350,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeSemi,
			EndReason:         events.RoundEndReasonCTWin,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        8,
			TeamBScore:        0,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            9,
			StartTick:         47446,
			StartFrame:        49908,
			EndTick:           49962,
			FreezeTimeEndTick: 48406,
			TeamAStartMoney:   38550,
			TeamBStartMoney:   19600,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeSemi,
			EndReason:         events.RoundEndReasonCTWin,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        9,
			TeamBScore:        0,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            10,
			StartTick:         50410,
			StartFrame:        53097,
			EndTick:           54640,
			FreezeTimeEndTick: 51370,
			TeamAStartMoney:   49050,
			TeamBStartMoney:   16050,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeSemi,
			EndReason:         events.RoundEndReasonTerroristsWin,
			WinnerSide:        common.TeamTerrorists,
			TeamAScore:        9,
			TeamBScore:        1,
			WinnerName:        expectedTeamNameB,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            11,
			StartTick:         55088,
			StartFrame:        58053,
			EndTick:           58260,
			FreezeTimeEndTick: 56048,
			TeamAStartMoney:   45500,
			TeamBStartMoney:   17400,
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
			StartTick:         58708,
			StartFrame:        61812,
			EndTick:           61096,
			FreezeTimeEndTick: 59668,
			TeamAStartMoney:   34750,
			TeamBStartMoney:   16400,
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
			StartTick:         62056,
			StartFrame:        65301,
			EndTick:           66969,
			FreezeTimeEndTick: 63016,
			TeamAStartMoney:   4000,
			TeamBStartMoney:   3200,
			TeamAEconomyType:  constants.EconomyTypePistol,
			TeamBEconomyType:  constants.EconomyTypePistol,
			EndReason:         events.RoundEndReasonTargetBombed,
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
			StartTick:         72051,
			StartFrame:        75518,
			EndTick:           75014,
			FreezeTimeEndTick: 73011,
			TeamAStartMoney:   19050,
			TeamBStartMoney:   8050,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeForceBuy,
			EndReason:         events.RoundEndReasonTerroristsWin,
			WinnerSide:        common.TeamTerrorists,
			TeamAScore:        13,
			TeamBScore:        1,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamTerrorists,
			TeamBSide:         common.TeamCounterTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
	}

	var players = []fake.FakePlayer{
		{
			SteamID64:     76561197997057216,
			Name:          "Knoeln",
			KillCount:     14,
			AssistCount:   2,
			DeathCount:    4,
			Score:         31,
			Team:          match.TeamA,
			HeadshotCount: 6,
			MvpCount:      3,
			UtilityDamage: 78,
		},
		{
			SteamID64:     76561199152902918,
			Name:          "lytikern",
			KillCount:     8,
			AssistCount:   4,
			DeathCount:    5,
			Score:         21,
			Team:          match.TeamA,
			HeadshotCount: 4,
			MvpCount:      1,
			UtilityDamage: 168,
		},
		{
			SteamID64:     76561199468837020,
			Name:          "nukon",
			KillCount:     7,
			AssistCount:   4,
			DeathCount:    3,
			Score:         19,
			Team:          match.TeamA,
			HeadshotCount: 5,
			MvpCount:      1,
			UtilityDamage: 175,
		},
		{
			SteamID64:     76561198376078356,
			Name:          "zoyu",
			KillCount:     21,
			AssistCount:   2,
			DeathCount:    3,
			Score:         48,
			Team:          match.TeamA,
			HeadshotCount: 12,
			MvpCount:      7,
			UtilityDamage: 91,
		},
		{
			SteamID64:     76561197960466083,
			Name:          "dieselkungen",
			KillCount:     9,
			AssistCount:   3,
			DeathCount:    4,
			Score:         22,
			Team:          match.TeamA,
			HeadshotCount: 3,
			MvpCount:      1,
			UtilityDamage: 38,
		},
		{
			SteamID64:     76561199018171213,
			Name:          "Morrizzz",
			KillCount:     3,
			AssistCount:   0,
			DeathCount:    13,
			Score:         6,
			Team:          match.TeamB,
			HeadshotCount: 1,
			MvpCount:      0,
			UtilityDamage: 2,
		},
		{
			SteamID64:     76561198153093657,
			Name:          "R0cKy",
			KillCount:     3,
			AssistCount:   4,
			DeathCount:    14,
			Score:         10,
			Team:          match.TeamB,
			HeadshotCount: 1,
			MvpCount:      0,
			UtilityDamage: 1,
		},
		{
			SteamID64:     76561198990304312,
			Name:          "AsSeLeLe",
			KillCount:     5,
			AssistCount:   1,
			DeathCount:    13,
			Score:         15,
			Team:          match.TeamB,
			HeadshotCount: 3,
			MvpCount:      0,
			UtilityDamage: 13,
		},
		{
			SteamID64:     76561197997455392,
			Name:          "UnknownG0DD",
			KillCount:     0, // 2 on in-game scoreboard but he did 2 teamkills so it's 0
			AssistCount:   1,
			DeathCount:    8,
			Score:         1,
			Team:          match.TeamB,
			HeadshotCount: 1,
			MvpCount:      0,
			UtilityDamage: 0,
		},
		{
			SteamID64:     76561198137170434,
			Name:          "rigby",
			KillCount:     6,
			AssistCount:   1,
			DeathCount:    13,
			Score:         14,
			Team:          match.TeamB,
			HeadshotCount: 3,
			MvpCount:      1,
			UtilityDamage: 56,
		},
	}

	assertion.AssertPlayers(t, match, players)
	assertion.AssertRounds(t, match, rounds)
}
