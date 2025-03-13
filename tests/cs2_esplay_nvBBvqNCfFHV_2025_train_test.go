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

// 5v5 with surrender.
// https://esplay.com/m/nvBBvqNCfFHV/team-pytonorm-vs-team-shawty
func TestEsplay_nvBBvqNCfFHV_2025_Train(t *testing.T) {
	demoName := "esplay_nvBBvqNCfFHV_2025_train"
	demoPath := testsutils.GetDemoPath("cs2", demoName)
	match, err := api.AnalyzeDemo(demoPath, api.AnalyzeDemoOptions{
		Source: constants.DemoSourceEsplay,
	})

	if err != nil {
		t.Error(err)
	}

	expectedRoundCount := 8
	expectedPlayerCount := 10
	expectedScoreTeamA := 8
	expectedScoreTeamB := 0
	expectedScoreFirstHalfTeamA := 8
	expectedScoreFirstHalfTeamB := 0
	expectedScoreSecondHalfTeamA := 0
	expectedScoreSecondHalfTeamB := 0
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
			StartFrame:        1149,
			EndTick:           3214,
			FreezeTimeEndTick: 895,
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
			StartTick:         3662,
			StartFrame:        5598,
			EndTick:           7212,
			FreezeTimeEndTick: 4558,
			TeamAStartMoney:   18450,
			TeamBStartMoney:   11000,
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
			StartTick:         7660,
			StartFrame:        10304,
			EndTick:           12033,
			FreezeTimeEndTick: 8556,
			TeamAStartMoney:   20900,
			TeamBStartMoney:   18650,
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
			StartTick:         12481,
			StartFrame:        16289,
			EndTick:           15226,
			FreezeTimeEndTick: 13377,
			TeamAStartMoney:   34150,
			TeamBStartMoney:   15900,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeSemi,
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
			StartTick:         15674,
			StartFrame:        19949,
			EndTick:           21123,
			FreezeTimeEndTick: 16570,
			TeamAStartMoney:   29200,
			TeamBStartMoney:   21950,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonBombDefused,
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
			StartTick:         21571,
			StartFrame:        27171,
			EndTick:           24466,
			FreezeTimeEndTick: 22467,
			TeamAStartMoney:   40400,
			TeamBStartMoney:   22700,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeFull,
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
			StartTick:         24914,
			StartFrame:        31045,
			EndTick:           30341,
			FreezeTimeEndTick: 25810,
			TeamAStartMoney:   37050,
			TeamBStartMoney:   19200,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeForceBuy,
			EndReason:         events.RoundEndReasonCTWin,
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
			StartTick:         30789,
			StartFrame:        38064,
			EndTick:           38271,
			FreezeTimeEndTick: 31685,
			TeamAStartMoney:   37300,
			TeamBStartMoney:   15150,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeForceBuy,
			EndReason:         events.RoundEndReasonBombDefused,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        8,
			TeamBScore:        0,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
	}

	var players = []fake.FakePlayer{
		{
			SteamID64:     76561198129332650,
			Name:          "Marmaduke",
			KillCount:     10,
			AssistCount:   0,
			DeathCount:    5,
			Score:         20,
			Team:          match.TeamA,
			HeadshotCount: 6,
			MvpCount:      2,
			UtilityDamage: 118,
		},
		{
			SteamID64:     76561198248008728,
			Name:          "Pytonorm",
			KillCount:     9,
			AssistCount:   0,
			DeathCount:    3,
			Score:         20,
			Team:          match.TeamA,
			HeadshotCount: 3,
			MvpCount:      2,
			UtilityDamage: 7,
		},
		{
			SteamID64:     76561197964870354,
			Name:          "ddane",
			KillCount:     7,
			AssistCount:   5,
			DeathCount:    6,
			Score:         16,
			Team:          match.TeamA,
			HeadshotCount: 3,
			MvpCount:      2,
			UtilityDamage: 149,
		},
		{
			SteamID64:     76561199261966104,
			Name:          "klas",
			KillCount:     7,
			AssistCount:   5,
			DeathCount:    4,
			Score:         21,
			Team:          match.TeamA,
			HeadshotCount: 2,
			MvpCount:      2,
			UtilityDamage: 99,
		},
		{
			SteamID64:     76561198114419765,
			Name:          "denfruktade",
			KillCount:     6,
			AssistCount:   3,
			DeathCount:    3,
			Score:         15,
			Team:          match.TeamA,
			HeadshotCount: 3,
			MvpCount:      0,
			UtilityDamage: 3,
		},
		{
			SteamID64:     76561198107761058,
			Name:          "Ruggan",
			KillCount:     3,
			AssistCount:   0,
			DeathCount:    8,
			Score:         9,
			Team:          match.TeamB,
			HeadshotCount: 3,
			MvpCount:      0,
			UtilityDamage: 0,
		},
		{
			SteamID64:     76561198994886574,
			Name:          "teC1",
			KillCount:     5,
			AssistCount:   1,
			DeathCount:    8,
			Score:         13,
			Team:          match.TeamB,
			HeadshotCount: 3,
			MvpCount:      0,
			UtilityDamage: 0,
		},
		{
			SteamID64:     76561198059437550,
			Name:          "Ludicaken",
			KillCount:     7,
			AssistCount:   2,
			DeathCount:    8,
			Score:         16,
			Team:          match.TeamB,
			HeadshotCount: 5,
			MvpCount:      0,
			UtilityDamage: 0,
		},
		{
			SteamID64:     76561198174150442,
			Name:          "Timpelina",
			KillCount:     3,
			AssistCount:   1,
			DeathCount:    8,
			Score:         6,
			Team:          match.TeamB,
			HeadshotCount: 3,
			MvpCount:      0,
			UtilityDamage: 0,
		},
		{
			SteamID64:     76561198018488298,
			Name:          "Shawty",
			KillCount:     3,
			AssistCount:   2,
			DeathCount:    7,
			Score:         8,
			Team:          match.TeamB,
			HeadshotCount: 2,
			MvpCount:      0,
			UtilityDamage: 0,
		},
	}

	assertion.AssertPlayers(t, match, players)
	assertion.AssertRounds(t, match, rounds)
}
