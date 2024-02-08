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

// Contains players disconnection + BOTs
func Test5EPlay_5eplay_g161_20240107155319726585099_2023_Nuke(t *testing.T) {
	demoName := "5eplay_g161-20240107155319726585099_2023_nuke"
	demoPath := testsutils.GetDemoPath("cs2", demoName)
	match, err := api.AnalyzeDemo(demoPath, api.AnalyzeDemoOptions{
		Source: constants.DemoSourceFiveEPlay,
	})
	if err != nil {
		t.Error(err)
	}

	expectedRoundCount := 22
	expectedPlayerCount := 10
	expectedScoreTeamA := 9
	expectedScoreTeamB := 13
	expectedScoreFirstHalfTeamA := 4
	expectedScoreFirstHalfTeamB := 8
	expectedScoreSecondHalfTeamA := 5
	expectedScoreSecondHalfTeamB := 5
	expectedTeamNameA := "Team A"
	expectedTeamNameB := "Team B"
	expectedWinnerName := expectedTeamNameB
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
			StartTick:         2368,
			StartFrame:        2482,
			EndTick:           6004,
			FreezeTimeEndTick: 3545,
			TeamAStartMoney:   4000,
			TeamBStartMoney:   4000,
			TeamAEconomyType:  constants.EconomyTypePistol,
			TeamBEconomyType:  constants.EconomyTypePistol,
			EndReason:         events.RoundEndReasonTerroristsWin,
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
			StartTick:         6452,
			StartFrame:        6754,
			EndTick:           8948,
			FreezeTimeEndTick: 7220,
			TeamAStartMoney:   11100,
			TeamBStartMoney:   19150,
			TeamAEconomyType:  constants.EconomyTypeSemi,
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
			StartTick:         13752,
			StartFrame:        14540,
			EndTick:           16392,
			FreezeTimeEndTick: 14520,
			TeamAStartMoney:   15050,
			TeamBStartMoney:   20750,
			TeamAEconomyType:  constants.EconomyTypeSemi,
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
			StartTick:         16840,
			StartFrame:        17813,
			EndTick:           24834,
			FreezeTimeEndTick: 17608,
			TeamAStartMoney:   22050,
			TeamBStartMoney:   32900,
			TeamAEconomyType:  constants.EconomyTypeFull,
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
			StartTick:         25282,
			StartFrame:        26552,
			EndTick:           31233,
			FreezeTimeEndTick: 26050,
			TeamAStartMoney:   20600,
			TeamBStartMoney:   40150,
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
			StartTick:         31681,
			StartFrame:        33190,
			EndTick:           36279,
			FreezeTimeEndTick: 32449,
			TeamAStartMoney:   18800,
			TeamBStartMoney:   28050,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonCTWin,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        2,
			TeamBScore:        4,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            7,
			StartTick:         36727,
			StartFrame:        38539,
			EndTick:           45544,
			FreezeTimeEndTick: 37495,
			TeamAStartMoney:   23150,
			TeamBStartMoney:   18400,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeSemi,
			EndReason:         events.RoundEndReasonTerroristsWin,
			WinnerSide:        common.TeamTerrorists,
			TeamAScore:        2,
			TeamBScore:        5,
			WinnerName:        expectedTeamNameB,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            8,
			StartTick:         45992,
			StartFrame:        48070,
			EndTick:           49753,
			FreezeTimeEndTick: 46760,
			TeamAStartMoney:   22100,
			TeamBStartMoney:   23350,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonTerroristsWin,
			WinnerSide:        common.TeamTerrorists,
			TeamAScore:        2,
			TeamBScore:        6,
			WinnerName:        expectedTeamNameB,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            9,
			StartTick:         50201,
			StartFrame:        52534,
			EndTick:           53117,
			FreezeTimeEndTick: 50969,
			TeamAStartMoney:   18200,
			TeamBStartMoney:   21700,
			TeamAEconomyType:  constants.EconomyTypeSemi,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonCTWin,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        3,
			TeamBScore:        6,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            10,
			StartTick:         53565,
			StartFrame:        56168,
			EndTick:           61292,
			FreezeTimeEndTick: 54333,
			TeamAStartMoney:   21250,
			TeamBStartMoney:   17300,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeSemi,
			EndReason:         events.RoundEndReasonTerroristsWin,
			WinnerSide:        common.TeamTerrorists,
			TeamAScore:        3,
			TeamBScore:        7,
			WinnerName:        expectedTeamNameB,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            11,
			StartTick:         61740,
			StartFrame:        64708,
			EndTick:           66041,
			FreezeTimeEndTick: 62508,
			TeamAStartMoney:   26000,
			TeamBStartMoney:   20700,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonTerroristsWin,
			WinnerSide:        common.TeamTerrorists,
			TeamAScore:        3,
			TeamBScore:        8,
			WinnerName:        expectedTeamNameB,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            12,
			StartTick:         66489,
			StartFrame:        69600,
			EndTick:           70745,
			FreezeTimeEndTick: 67257,
			TeamAStartMoney:   21750,
			TeamBStartMoney:   21550,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonCTWin,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        4,
			TeamBScore:        8,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            13,
			StartTick:         71289,
			StartFrame:        74646,
			EndTick:           77613,
			FreezeTimeEndTick: 72473,
			TeamAStartMoney:   4000,
			TeamBStartMoney:   4000,
			TeamAEconomyType:  constants.EconomyTypePistol,
			TeamBEconomyType:  constants.EconomyTypePistol,
			EndReason:         events.RoundEndReasonTargetBombed,
			WinnerSide:        common.TeamTerrorists,
			TeamAScore:        5,
			TeamBScore:        8,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamTerrorists,
			TeamBSide:         common.TeamCounterTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            14,
			StartTick:         78061,
			StartFrame:        81613,
			EndTick:           82331,
			FreezeTimeEndTick: 78829,
			TeamAStartMoney:   20400,
			TeamBStartMoney:   11150,
			TeamAEconomyType:  constants.EconomyTypeSemi,
			TeamBEconomyType:  constants.EconomyTypeSemi,
			EndReason:         events.RoundEndReasonTerroristsWin,
			WinnerSide:        common.TeamTerrorists,
			TeamAScore:        6,
			TeamBScore:        8,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamTerrorists,
			TeamBSide:         common.TeamCounterTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            15,
			StartTick:         82779,
			StartFrame:        86559,
			EndTick:           86989,
			FreezeTimeEndTick: 83547,
			TeamAStartMoney:   23850,
			TeamBStartMoney:   14800,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeSemi,
			EndReason:         events.RoundEndReasonTerroristsWin,
			WinnerSide:        common.TeamTerrorists,
			TeamAScore:        7,
			TeamBScore:        8,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamTerrorists,
			TeamBSide:         common.TeamCounterTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            16,
			StartTick:         87437,
			StartFrame:        91539,
			EndTick:           91380,
			FreezeTimeEndTick: 88205,
			TeamAStartMoney:   26000,
			TeamBStartMoney:   23000,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonCTWin,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        7,
			TeamBScore:        9,
			WinnerName:        expectedTeamNameB,
			TeamASide:         common.TeamTerrorists,
			TeamBSide:         common.TeamCounterTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            17,
			StartTick:         91828,
			StartFrame:        96189,
			EndTick:           98045,
			FreezeTimeEndTick: 92596,
			TeamAStartMoney:   23400,
			TeamBStartMoney:   18650,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeFull,
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
			StartTick:         98493,
			StartFrame:        103124,
			EndTick:           105430,
			FreezeTimeEndTick: 99261,
			TeamAStartMoney:   22900,
			TeamBStartMoney:   25900,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonCTWin,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        8,
			TeamBScore:        10,
			WinnerName:        expectedTeamNameB,
			TeamASide:         common.TeamTerrorists,
			TeamBSide:         common.TeamCounterTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            19,
			StartTick:         105878,
			StartFrame:        110823,
			EndTick:           111125,
			FreezeTimeEndTick: 106646,
			TeamAStartMoney:   12950,
			TeamBStartMoney:   19750,
			TeamAEconomyType:  constants.EconomyTypeSemi,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonCTWin,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        8,
			TeamBScore:        11,
			WinnerName:        expectedTeamNameB,
			TeamASide:         common.TeamTerrorists,
			TeamBSide:         common.TeamCounterTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            20,
			StartTick:         111573,
			StartFrame:        116770,
			EndTick:           116192,
			FreezeTimeEndTick: 112341,
			TeamAStartMoney:   15050,
			TeamBStartMoney:   20000,
			TeamAEconomyType:  constants.EconomyTypeSemi,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonCTWin,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        8,
			TeamBScore:        12,
			WinnerName:        expectedTeamNameB,
			TeamASide:         common.TeamTerrorists,
			TeamBSide:         common.TeamCounterTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            21,
			StartTick:         116640,
			StartFrame:        122210,
			EndTick:           120134,
			FreezeTimeEndTick: 117408,
			TeamAStartMoney:   16250,
			TeamBStartMoney:   25700,
			TeamAEconomyType:  constants.EconomyTypeForceBuy,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonTerroristsWin,
			WinnerSide:        common.TeamTerrorists,
			TeamAScore:        9,
			TeamBScore:        12,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamTerrorists,
			TeamBSide:         common.TeamCounterTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            22,
			StartTick:         120582,
			StartFrame:        126469,
			EndTick:           123848,
			FreezeTimeEndTick: 121350,
			TeamAStartMoney:   12600,
			TeamBStartMoney:   18950,
			TeamAEconomyType:  constants.EconomyTypeSemi,
			TeamBEconomyType:  constants.EconomyTypeForceBuy,
			EndReason:         events.RoundEndReasonCTWin,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        9,
			TeamBScore:        13,
			WinnerName:        expectedTeamNameB,
			TeamASide:         common.TeamTerrorists,
			TeamBSide:         common.TeamCounterTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
	}

	var players = []fake.FakePlayer{
		{
			SteamID64:     76561198285600732,
			Name:          "妮蔻",
			KillCount:     22,
			AssistCount:   6,
			DeathCount:    19,
			Score:         51,
			Team:          match.TeamA,
			HeadshotCount: 10,
			MvpCount:      4,
			UtilityDamage: 82,
		},
		{
			SteamID64:     76561198381364815,
			Name:          "寒王",
			KillCount:     19,
			AssistCount:   3,
			DeathCount:    16,
			Score:         48,
			Team:          match.TeamA,
			HeadshotCount: 7,
			MvpCount:      2,
			UtilityDamage: 124,
		},
		{
			SteamID64: 76561198184530382,
			Name:      "魔男",
			KillCount: 15,
			// It's 4 in-game because this player did an assist while controlling a bot and CS2 increments assists if
			// the player was controlling a bot but doesn't for kills and deaths.
			// For consistency we don't count assists done while controlling a bot.
			AssistCount:   3,
			DeathCount:    18,
			Score:         36,
			Team:          match.TeamA,
			HeadshotCount: 7,
			MvpCount:      2,
			UtilityDamage: 10,
		},
		{
			// This player disconnected at the beginning of the round 2 and replaced by a bot.
			// When he reconnected during round 17 he changed his nickname and the in-game scoreboard shows him as
			// "disconnected" while is actually playing. The scoreboard is not reliable.
			SteamID64:     76561198214055247,
			Name:          "雨神",
			KillCount:     2,
			AssistCount:   1,
			DeathCount:    5,
			Score:         3,
			Team:          match.TeamA,
			HeadshotCount: 2,
			MvpCount:      0,
			UtilityDamage: 0,
		},
		{
			SteamID64:     76561199084592315,
			Name:          "西若",
			KillCount:     15,
			AssistCount:   5,
			DeathCount:    20,
			Score:         36,
			Team:          match.TeamA,
			HeadshotCount: 4,
			MvpCount:      1,
			UtilityDamage: 2,
		},
		{
			SteamID64:     76561199094835698,
			Name:          "森破",
			KillCount:     12,
			AssistCount:   9,
			DeathCount:    18,
			Score:         33,
			Team:          match.TeamB,
			HeadshotCount: 11,
			MvpCount:      1,
			UtilityDamage: 0,
		},
		{
			SteamID64:     76561198448371692,
			Name:          "载物",
			KillCount:     20,
			AssistCount:   3,
			DeathCount:    16,
			Score:         43,
			Team:          match.TeamB,
			HeadshotCount: 7,
			MvpCount:      2,
			UtilityDamage: 0,
		},
		{
			SteamID64:     76561199157434599,
			Name:          "阿电",
			KillCount:     21,
			AssistCount:   4,
			DeathCount:    14,
			Score:         55,
			Team:          match.TeamB,
			HeadshotCount: 9,
			MvpCount:      4,
			UtilityDamage: 64,
		},
		{
			SteamID64:     76561199064836174,
			Name:          "小孩",
			KillCount:     12,
			AssistCount:   8,
			DeathCount:    16,
			Score:         35,
			Team:          match.TeamB,
			HeadshotCount: 5,
			MvpCount:      2,
			UtilityDamage: 201,
		},
		{
			SteamID64:     76561199002382108,
			Name:          "教父",
			KillCount:     22,
			AssistCount:   7,
			DeathCount:    15,
			Score:         53,
			Team:          match.TeamB,
			HeadshotCount: 15,
			MvpCount:      4,
			UtilityDamage: 137,
		},
	}

	assertion.AssertPlayers(t, match, players)
	assertion.AssertRounds(t, match, rounds)
}
