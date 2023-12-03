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

// https://www.hltv.org/matches/2367504/monte-vs-og-roobet-cup-2023
func TestEbot_Monte_VS_OG_Roobet_Cup_2023_Anubis(t *testing.T) {
	demoName := "ebot_monte_vs_og_roobet_cup_2023_anubis"
	demoPath := testsutils.GetDemoPath("cs2", demoName)
	match, err := api.AnalyzeDemo(demoPath, api.AnalyzeDemoOptions{
		Source: constants.DemoSourceEbot,
	})
	if err != nil {
		t.Error(err)
	}

	expectedRoundCount := 20
	expectedPlayerCount := 10
	expectedScoreTeamA := 7
	expectedScoreTeamB := 13
	expectedScoreFirstHalfTeamA := 5
	expectedScoreFirstHalfTeamB := 7
	expectedScoreSecondHalfTeamA := 2
	expectedScoreSecondHalfTeamB := 6
	expectedTeamNameA := "OG Esports"
	expectedTeamNameB := "Monte"
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
			StartTick:         0,
			StartFrame:        15,
			EndTick:           5003,
			FreezeTimeEndTick: 0,
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
			StartTick:         5451,
			StartFrame:        5611,
			EndTick:           14079,
			FreezeTimeEndTick: 6731,
			TeamAStartMoney:   21250,
			TeamBStartMoney:   15300,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeForceBuy,
			EndReason:         events.RoundEndReasonTargetBombed,
			WinnerSide:        common.TeamTerrorists,
			TeamAScore:        1,
			TeamBScore:        1,
			WinnerName:        expectedTeamNameB,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
		},
		{
			Number:            3,
			StartTick:         14527,
			StartFrame:        15102,
			EndTick:           20532,
			FreezeTimeEndTick: 15807,
			TeamAStartMoney:   8400,
			TeamBStartMoney:   19200,
			TeamAEconomyType:  constants.EconomyTypeSemi,
			TeamBEconomyType:  constants.EconomyTypeSemi,
			EndReason:         events.RoundEndReasonTargetBombed,
			WinnerSide:        common.TeamTerrorists,
			TeamAScore:        1,
			TeamBScore:        2,
			WinnerName:        expectedTeamNameB,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
		},
		{
			Number:            4,
			StartTick:         20980,
			StartFrame:        21906,
			EndTick:           23875,
			FreezeTimeEndTick: 22260,
			TeamAStartMoney:   12950,
			TeamBStartMoney:   23950,
			TeamAEconomyType:  constants.EconomyTypeEco,
			TeamBEconomyType:  constants.EconomyTypeFull,
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
			StartTick:         24323,
			StartFrame:        25548,
			EndTick:           30214,
			FreezeTimeEndTick: 25603,
			TeamAStartMoney:   24650,
			TeamBStartMoney:   32600,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonTerroristsWin,
			WinnerSide:        common.TeamTerrorists,
			TeamAScore:        1,
			TeamBScore:        4,
			WinnerName:        expectedTeamNameB,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
		},
		{
			Number:            6,
			StartTick:         30662,
			StartFrame:        32368,
			EndTick:           36772,
			FreezeTimeEndTick: 31942,
			TeamAStartMoney:   15350,
			TeamBStartMoney:   47250,
			TeamAEconomyType:  constants.EconomyTypeSemi,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonTerroristsWin,
			WinnerSide:        common.TeamTerrorists,
			TeamAScore:        1,
			TeamBScore:        5,
			WinnerName:        expectedTeamNameB,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
		},
		{
			Number:            7,
			StartTick:         37220,
			StartFrame:        39364,
			EndTick:           42791,
			FreezeTimeEndTick: 38500,
			TeamAStartMoney:   28300,
			TeamBStartMoney:   48400,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonTerroristsWin,
			WinnerSide:        common.TeamTerrorists,
			TeamAScore:        1,
			TeamBScore:        6,
			WinnerName:        expectedTeamNameB,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
		},
		{
			Number:            8,
			StartTick:         43239,
			StartFrame:        45896,
			EndTick:           50503,
			FreezeTimeEndTick: 44519,
			TeamAStartMoney:   19300,
			TeamBStartMoney:   46250,
			TeamAEconomyType:  constants.EconomyTypeSemi,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonCTWin,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        2,
			TeamBScore:        6,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
		},
		{
			Number:            9,
			StartTick:         53087,
			StartFrame:        56372,
			EndTick:           60234,
			FreezeTimeEndTick: 54367,
			TeamAStartMoney:   27800,
			TeamBStartMoney:   33250,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonCTWin,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        3,
			TeamBScore:        6,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
		},
		{
			Number:            10,
			StartTick:         63462,
			StartFrame:        67423,
			EndTick:           70460,
			FreezeTimeEndTick: 64742,
			TeamAStartMoney:   22650,
			TeamBStartMoney:   17100,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeForceBuy,
			EndReason:         events.RoundEndReasonTargetBombed,
			WinnerSide:        common.TeamTerrorists,
			TeamAScore:        3,
			TeamBScore:        7,
			WinnerName:        expectedTeamNameB,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
		},
		{
			Number:            11,
			StartTick:         73805,
			StartFrame:        78559,
			EndTick:           81141,
			FreezeTimeEndTick: 75085,
			TeamAStartMoney:   27350,
			TeamBStartMoney:   19950,
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
			StartTick:         81589,
			StartFrame:        86731,
			EndTick:           85786,
			FreezeTimeEndTick: 82869,
			TeamAStartMoney:   19850,
			TeamBStartMoney:   8150,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeForceBuy,
			EndReason:         events.RoundEndReasonCTWin,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        5,
			TeamBScore:        7,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
		},
		{
			Number:            13,
			StartTick:         86524,
			StartFrame:        92005,
			EndTick:           90566,
			FreezeTimeEndTick: 88220,
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
		},
		{
			Number:            14,
			StartTick:         91014,
			StartFrame:        96680,
			EndTick:           93684,
			FreezeTimeEndTick: 92294,
			TeamAStartMoney:   10950,
			TeamBStartMoney:   18350,
			TeamAEconomyType:  constants.EconomyTypeEco,
			TeamBEconomyType:  constants.EconomyTypeSemi,
			EndReason:         events.RoundEndReasonCTWin,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        5,
			TeamBScore:        9,
			WinnerName:        expectedTeamNameB,
			TeamASide:         common.TeamTerrorists,
			TeamBSide:         common.TeamCounterTerrorists,
		},
		{
			Number:            15,
			StartTick:         94132,
			StartFrame:        100108,
			EndTick:           102314,
			FreezeTimeEndTick: 95412,
			TeamAStartMoney:   23250,
			TeamBStartMoney:   18200,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeSemi,
			EndReason:         events.RoundEndReasonTerroristsWin,
			WinnerSide:        common.TeamTerrorists,
			TeamAScore:        6,
			TeamBScore:        9,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamTerrorists,
			TeamBSide:         common.TeamCounterTerrorists,
		},
		{
			Number:            16,
			StartTick:         102762,
			StartFrame:        109161,
			EndTick:           108165,
			FreezeTimeEndTick: 104042,
			TeamAStartMoney:   18900,
			TeamBStartMoney:   21000,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeForceBuy,
			EndReason:         events.RoundEndReasonTerroristsWin,
			WinnerSide:        common.TeamTerrorists,
			TeamAScore:        7,
			TeamBScore:        9,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamTerrorists,
			TeamBSide:         common.TeamCounterTerrorists,
		},
		{
			Number:            17,
			StartTick:         108613,
			StartFrame:        115467,
			EndTick:           117253,
			FreezeTimeEndTick: 109893,
			TeamAStartMoney:   21450,
			TeamBStartMoney:   10350,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeEco,
			EndReason:         events.RoundEndReasonTargetSaved,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        7,
			TeamBScore:        10,
			WinnerName:        expectedTeamNameB,
			TeamASide:         common.TeamTerrorists,
			TeamBSide:         common.TeamCounterTerrorists,
		},
		{
			Number:            18,
			StartTick:         120578,
			StartFrame:        127960,
			EndTick:           129218,
			FreezeTimeEndTick: 121858,
			TeamAStartMoney:   8950,
			TeamBStartMoney:   26500,
			TeamAEconomyType:  constants.EconomyTypeForceBuy,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonTargetSaved,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        7,
			TeamBScore:        11,
			WinnerName:        expectedTeamNameB,
			TeamASide:         common.TeamTerrorists,
			TeamBSide:         common.TeamCounterTerrorists,
		},
		{
			Number:            19,
			StartTick:         129666,
			StartFrame:        137434,
			EndTick:           135076,
			FreezeTimeEndTick: 130946,
			TeamAStartMoney:   7800,
			TeamBStartMoney:   24850,
			TeamAEconomyType:  constants.EconomyTypeEco,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonCTWin,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        7,
			TeamBScore:        12,
			WinnerName:        expectedTeamNameB,
			TeamASide:         common.TeamTerrorists,
			TeamBSide:         common.TeamCounterTerrorists,
		},
		{
			Number:            20,
			StartTick:         135524,
			StartFrame:        143680,
			EndTick:           139122,
			FreezeTimeEndTick: 136804,
			TeamAStartMoney:   15100,
			TeamBStartMoney:   44250,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonCTWin,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        7,
			TeamBScore:        13,
			WinnerName:        expectedTeamNameB,
			TeamASide:         common.TeamTerrorists,
			TeamBSide:         common.TeamCounterTerrorists,
		},
	}

	var players = []fake.FakePlayer{
		{
			SteamID64:     76561198098463229,
			Name:          "regali",
			KillCount:     16,
			AssistCount:   5,
			DeathCount:    14,
			Score:         42,
			Team:          match.TeamA,
			HeadshotCount: 8,
			MvpCount:      2,
			UtilityDamage: 259,
		},
		{
			SteamID64:     76561198182462587,
			Name:          "k1to",
			KillCount:     17,
			AssistCount:   7,
			DeathCount:    17,
			Score:         41,
			Team:          match.TeamA,
			HeadshotCount: 12,
			MvpCount:      1,
			UtilityDamage: 97,
		},
		{
			SteamID64:   76561198252434500,
			Name:        "F1KU",
			KillCount:   14,
			AssistCount: 8,
			// It should be 18 because during the round 10 F1KU has been killed by sdy with a molotov but CS2 detected
			// it as a suicide from F1KU. The kill event contains only the victim.
			DeathCount:    17,
			Score:         36,
			Team:          match.TeamA,
			HeadshotCount: 6,
			MvpCount:      1,
			UtilityDamage: 136,
		},
		{
			SteamID64:     76561197963504595,
			Name:          "FASHR",
			KillCount:     10,
			AssistCount:   4,
			DeathCount:    16,
			Score:         25,
			Team:          match.TeamA,
			HeadshotCount: 7,
			MvpCount:      1,
			UtilityDamage: 16,
		},
		{
			SteamID64:     76561197999825422,
			Name:          "nexa",
			KillCount:     9,
			AssistCount:   5,
			DeathCount:    15,
			Score:         25,
			Team:          match.TeamA,
			HeadshotCount: 5,
			MvpCount:      2,
			UtilityDamage: 35,
		},
		{
			SteamID64:     76561198060483793,
			Name:          "br0",
			KillCount:     14,
			AssistCount:   10,
			DeathCount:    16,
			Score:         42,
			Team:          match.TeamB,
			HeadshotCount: 9,
			MvpCount:      4,
			UtilityDamage: 109,
		},
		{
			SteamID64:     76561198975452660,
			Name:          "Woro2k",
			KillCount:     25,
			AssistCount:   1,
			DeathCount:    9,
			Score:         61,
			Team:          match.TeamB,
			HeadshotCount: 10,
			MvpCount:      4,
			UtilityDamage: 160,
		},
		{
			SteamID64:     76561198296417934,
			Name:          "kRaSnaL",
			KillCount:     13,
			AssistCount:   5,
			DeathCount:    15,
			Score:         34,
			Team:          match.TeamB,
			HeadshotCount: 6,
			MvpCount:      1,
			UtilityDamage: 111,
		},
		{
			SteamID64: 76561198040577200,
			Name:      "sdy",
			// It should be 14 because during the round 10 sdy killed F1KU with a molotov but CS2 detected it as a
			// suicide from F1KU. The kill event contains only the victim.
			KillCount:     13,
			AssistCount:   8,
			DeathCount:    14,
			Score:         49,
			Team:          match.TeamB,
			HeadshotCount: 6,
			MvpCount:      3,
			UtilityDamage: 68,
		},
		{
			SteamID64:     76561198126171426,
			Name:          "DemQQ",
			KillCount:     14,
			AssistCount:   7,
			DeathCount:    12,
			Score:         36,
			Team:          match.TeamB,
			HeadshotCount: 8,
			MvpCount:      1,
			UtilityDamage: 130,
		},
	}

	assertion.AssertPlayers(t, match, players)
	assertion.AssertRounds(t, match, rounds)
}
