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

// https://www.hltv.org/stats/matches/mapstatsid/40296/astralis-vs-envy
func TestEbot_Astralis_VS_Envyus_Game_Show_Global_eSports_Cup_2016_Cache(t *testing.T) {
	demoName := "ebot_astralis_vs_envyus_game_show_global_esports_cup_2016_cache"
	demoPath := testsutils.GetDemoPath("csgo", demoName)
	match, err := api.AnalyzeDemo(demoPath, api.AnalyzeDemoOptions{
		Source: constants.DemoSourceEbot,
	})
	if err != nil {
		t.Error(err)
	}

	expectedRoundCount := 29
	expectedPlayerCount := 10
	expectedScoreTeamA := 13
	expectedScoreTeamB := 16
	expectedScoreFirstHalfTeamA := 9
	expectedScoreFirstHalfTeamB := 6
	expectedScoreSecondHalfTeamA := 4
	expectedScoreSecondHalfTeamB := 10
	expectedTeamNameA := "Astralis"
	expectedTeamNameB := "EnVyUs"
	expectedWinnerName := expectedTeamNameB
	expectedMaxRounds := 30

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
			StartFrame:        6,
			EndTick:           11929,
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
	}

	var players = []fake.FakePlayer{
		{
			SteamID64:     76561197978241352,
			Name:          "ENVYUS HappyV",
			KillCount:     23,
			AssistCount:   2,
			DeathCount:    16,
			Score:         51,
			Team:          match.TeamB,
			HeadshotCount: 10,
			MvpCount:      1,
			UtilityDamage: 32,
		},
		{
			SteamID64:     76561197990682262,
			Name:          "Xyp9x",
			KillCount:     10,
			AssistCount:   8,
			DeathCount:    25,
			Score:         32,
			Team:          match.TeamA,
			HeadshotCount: 3,
			MvpCount:      3,
			UtilityDamage: 146,
		},
		{
			SteamID64:     76561198024905796,
			Name:          "ENVYUS KENNYS -M-",
			KillCount:     30,
			AssistCount:   4,
			DeathCount:    17,
			Score:         75,
			Team:          match.TeamB,
			HeadshotCount: 6,
			MvpCount:      4,
			UtilityDamage: 42,
		},
		{
			SteamID64:     76561197960710573,
			Name:          "ENVYUS NBK-",
			KillCount:     19,
			AssistCount:   4,
			DeathCount:    19,
			Score:         52,
			Team:          match.TeamB,
			HeadshotCount: 10,
			MvpCount:      3,
			UtilityDamage: 20,
		},
		{
			SteamID64:     76561197989744167,
			Name:          "ENVYUS apEXmousse[D]",
			KillCount:     20,
			AssistCount:   2,
			DeathCount:    20,
			Score:         48,
			Team:          match.TeamB,
			HeadshotCount: 9,
			MvpCount:      3,
			UtilityDamage: 105,
		},
		{
			SteamID64:     76561197987713664,
			Name:          "dEV1CEE -M-",
			KillCount:     12,
			AssistCount:   5,
			DeathCount:    21,
			Score:         32,
			Team:          match.TeamA,
			HeadshotCount: 5,
			MvpCount:      1,
			UtilityDamage: 52,
		},
		{
			SteamID64:     76561197989430253,
			Name:          "kARR1GANN",
			KillCount:     23,
			AssistCount:   3,
			DeathCount:    24,
			Score:         56,
			Team:          match.TeamA,
			HeadshotCount: 11,
			MvpCount:      3,
			UtilityDamage: 36,
		},
		{
			SteamID64:     76561198000782895,
			Name:          "ENVYUS k10$|-|1m@[C]",
			KillCount:     17,
			AssistCount:   3,
			DeathCount:    21,
			Score:         49,
			Team:          match.TeamB,
			HeadshotCount: 6,
			MvpCount:      5,
			UtilityDamage: 349,
		},
		{
			SteamID64:     76561197996352604,
			Name:          "cajunb",
			KillCount:     27,
			AssistCount:   3,
			DeathCount:    21,
			Score:         63,
			Team:          match.TeamA,
			HeadshotCount: 9,
			MvpCount:      2,
			UtilityDamage: 79,
		},
		{
			SteamID64:     76561198004854956,
			Name:          "dupreeh <3 OhLongJohnson",
			KillCount:     21,
			AssistCount:   5,
			DeathCount:    18,
			Score:         56,
			Team:          match.TeamA,
			HeadshotCount: 10,
			MvpCount:      4,
			UtilityDamage: 60,
		},
	}

	assertion.AssertPlayers(t, match, players)
	assertion.AssertRounds(t, match, rounds)
}
