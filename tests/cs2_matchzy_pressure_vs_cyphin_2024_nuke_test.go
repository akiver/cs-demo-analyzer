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

// The first round is restored 2 times
func Test_MatchZy_PressurE_vs_cyphin_2024_nuke(t *testing.T) {
	demoName := "matchzy_pressure_vs_cyphin_2024_nuke"
	demoPath := testsutils.GetDemoPath("cs2", demoName)
	match, err := api.AnalyzeDemo(demoPath, api.AnalyzeDemoOptions{
		Source: constants.DemoSourceMatchZy,
	})
	if err != nil {
		t.Error(err)
	}

	expectedRoundCount := 17
	expectedPlayerCount := 10
	expectedScoreTeamA := 13
	expectedScoreTeamB := 4
	expectedScoreFirstHalfTeamA := 10
	expectedScoreFirstHalfTeamB := 2
	expectedScoreSecondHalfTeamA := 3
	expectedScoreSecondHalfTeamB := 2
	expectedTeamNameA := "team_cyphin"
	expectedTeamNameB := "team_PressurE"
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
			StartTick:         17681,
			StartFrame:        18329,
			EndTick:           21025,
			FreezeTimeEndTick: 18832,
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
			StartTick:         21345,
			StartFrame:        22120,
			EndTick:           25468,
			FreezeTimeEndTick: 22497,
			TeamAStartMoney:   18500,
			TeamBStartMoney:   10500,
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
			StartTick:         25788,
			StartFrame:        26848,
			EndTick:           33118,
			FreezeTimeEndTick: 26940,
			TeamAStartMoney:   20650,
			TeamBStartMoney:   21900,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeFull,
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
			StartTick:         33438,
			StartFrame:        34878,
			EndTick:           39707,
			FreezeTimeEndTick: 34590,
			TeamAStartMoney:   24300,
			TeamBStartMoney:   15400,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeEco,
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
			StartTick:         40027,
			StartFrame:        41748,
			EndTick:           45725,
			FreezeTimeEndTick: 41179,
			TeamAStartMoney:   27300,
			TeamBStartMoney:   28600,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonTargetBombed,
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
			StartTick:         46045,
			StartFrame:        48177,
			EndTick:           53082,
			FreezeTimeEndTick: 47197,
			TeamAStartMoney:   20650,
			TeamBStartMoney:   21950,
			TeamAEconomyType:  constants.EconomyTypeFull,
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
			StartTick:         53402,
			StartFrame:        55929,
			EndTick:           60316,
			FreezeTimeEndTick: 54554,
			TeamAStartMoney:   18550,
			TeamBStartMoney:   28950,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonBombDefused,
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
			StartTick:         60636,
			StartFrame:        63661,
			EndTick:           68353,
			FreezeTimeEndTick: 61788,
			TeamAStartMoney:   24550,
			TeamBStartMoney:   27450,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonBombDefused,
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
			StartTick:         68673,
			StartFrame:        72159,
			EndTick:           74328,
			FreezeTimeEndTick: 69825,
			TeamAStartMoney:   29750,
			TeamBStartMoney:   25350,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeFull,
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
			StartTick:         74648,
			StartFrame:        78446,
			EndTick:           79949,
			FreezeTimeEndTick: 75800,
			TeamAStartMoney:   30250,
			TeamBStartMoney:   21100,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeSemi,
			EndReason:         events.RoundEndReasonTerroristsWin,
			WinnerSide:        common.TeamTerrorists,
			TeamAScore:        8,
			TeamBScore:        2,
			WinnerName:        expectedTeamNameB,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            11,
			StartTick:         80269,
			StartFrame:        84409,
			EndTick:           85974,
			FreezeTimeEndTick: 81421,
			TeamAStartMoney:   29750,
			TeamBStartMoney:   20750,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonCTWin,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        9,
			TeamBScore:        2,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            12,
			StartTick:         86294,
			StartFrame:        90772,
			EndTick:           88910,
			FreezeTimeEndTick: 87446,
			TeamAStartMoney:   24500,
			TeamBStartMoney:   21000,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonCTWin,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        10,
			TeamBScore:        2,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            13,
			StartTick:         89870,
			StartFrame:        94644,
			EndTick:           93038,
			FreezeTimeEndTick: 91022,
			TeamAStartMoney:   4000,
			TeamBStartMoney:   4000,
			TeamAEconomyType:  constants.EconomyTypePistol,
			TeamBEconomyType:  constants.EconomyTypePistol,
			EndReason:         events.RoundEndReasonTerroristsWin,
			WinnerSide:        common.TeamTerrorists,
			TeamAScore:        11,
			TeamBScore:        2,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamTerrorists,
			TeamBSide:         common.TeamCounterTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            14,
			StartTick:         93358,
			StartFrame:        98409,
			EndTick:           97123,
			FreezeTimeEndTick: 94510,
			TeamAStartMoney:   18500,
			TeamBStartMoney:   10700,
			TeamAEconomyType:  constants.EconomyTypeSemi,
			TeamBEconomyType:  constants.EconomyTypeSemi,
			EndReason:         events.RoundEndReasonTerroristsWin,
			WinnerSide:        common.TeamTerrorists,
			TeamAScore:        12,
			TeamBScore:        2,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamTerrorists,
			TeamBSide:         common.TeamCounterTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            15,
			StartTick:         97443,
			StartFrame:        102780,
			EndTick:           101859,
			FreezeTimeEndTick: 98595,
			TeamAStartMoney:   22250,
			TeamBStartMoney:   18700,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeForceBuy,
			EndReason:         events.RoundEndReasonCTWin,
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
			StartTick:         102179,
			StartFrame:        107958,
			EndTick:           107587,
			FreezeTimeEndTick: 103331,
			TeamAStartMoney:   15650,
			TeamBStartMoney:   18200,
			TeamAEconomyType:  constants.EconomyTypeSemi,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonBombDefused,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        12,
			TeamBScore:        4,
			WinnerName:        expectedTeamNameB,
			TeamASide:         common.TeamTerrorists,
			TeamBSide:         common.TeamCounterTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            17,
			StartTick:         107907,
			StartFrame:        114080,
			EndTick:           116221,
			FreezeTimeEndTick: 109059,
			TeamAStartMoney:   24850,
			TeamBStartMoney:   22600,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonTerroristsWin,
			WinnerSide:        common.TeamTerrorists,
			TeamAScore:        13,
			TeamBScore:        4,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamTerrorists,
			TeamBSide:         common.TeamCounterTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
	}
	var players = []fake.FakePlayer{
		{
			SteamID64:     76561198005227870,
			Name:          "kdz",
			KillCount:     15,
			AssistCount:   6,
			DeathCount:    6,
			Score:         40,
			Team:          match.TeamA,
			HeadshotCount: 9,
			MvpCount:      4,
			UtilityDamage: 109,
		},
		{
			SteamID64:     76561198065793422,
			Name:          "oer",
			KillCount:     21,
			AssistCount:   4,
			DeathCount:    12,
			Score:         48,
			Team:          match.TeamA,
			HeadshotCount: 12,
			MvpCount:      5,
			UtilityDamage: 106,
		},
		{
			SteamID64:     76561198042767969,
			Name:          "dusty",
			KillCount:     11,
			AssistCount:   4,
			DeathCount:    5,
			Score:         27,
			Team:          match.TeamA,
			HeadshotCount: 2,
			MvpCount:      2,
			UtilityDamage: 43,
		},
		{
			SteamID64:     76561198010988125,
			Name:          "cyphin",
			KillCount:     16,
			AssistCount:   6,
			DeathCount:    10,
			Score:         41,
			Team:          match.TeamA,
			HeadshotCount: 10,
			MvpCount:      1,
			UtilityDamage: 12,
		},
		{
			SteamID64:     76561198082434923,
			Name:          "Yakov",
			KillCount:     10,
			AssistCount:   5,
			DeathCount:    8,
			Score:         29,
			Team:          match.TeamA,
			HeadshotCount: 5,
			MvpCount:      1,
			UtilityDamage: 108,
		},
		{
			SteamID64:     76561198861089344,
			Name:          "bMd",
			KillCount:     14,
			AssistCount:   3,
			DeathCount:    16,
			Score:         35,
			Team:          match.TeamB,
			HeadshotCount: 9,
			MvpCount:      1,
			UtilityDamage: 0,
		},
		{
			SteamID64:     76561197982702431,
			Name:          "-maV.",
			KillCount:     8,
			AssistCount:   1,
			DeathCount:    15,
			Score:         18,
			Team:          match.TeamB,
			HeadshotCount: 3,
			MvpCount:      0,
			UtilityDamage: 4,
		},
		{
			SteamID64:     76561198208988596,
			Name:          "tamarac",
			KillCount:     7,
			AssistCount:   1,
			DeathCount:    14,
			Score:         23,
			Team:          match.TeamB,
			HeadshotCount: 4,
			MvpCount:      1,
			UtilityDamage: 0,
		},
		{
			SteamID64:     76561197974414204,
			Name:          "PressurE",
			KillCount:     7,
			AssistCount:   4,
			DeathCount:    15,
			Score:         22,
			Team:          match.TeamB,
			HeadshotCount: 2,
			MvpCount:      2,
			UtilityDamage: 0,
		},
		{
			SteamID64:     76561197972550155,
			Name:          "davetime",
			KillCount:     5,
			AssistCount:   3,
			DeathCount:    13,
			Score:         17,
			Team:          match.TeamB,
			HeadshotCount: 4,
			MvpCount:      0,
			UtilityDamage: 50,
		},
	}

	assertion.AssertRounds(t, match, rounds)
	assertion.AssertPlayers(t, match, players)
}
