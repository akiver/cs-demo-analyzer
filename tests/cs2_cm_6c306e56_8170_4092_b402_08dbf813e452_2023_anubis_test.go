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

// https://www.challengermode.com/s/CsgoAllstars/games/b4d31195-bae0-42c0-bbcc-08dbf6b17863
func TestChallengerMode_6c306e56_8170_4092_b402_08dbf813e452_2023_Anubis(t *testing.T) {
	demoName := "challengermode_6c306e56-8170-4092-b402-08dbf813e452"
	demoPath := testsutils.GetDemoPath("cs2", demoName)
	match, err := api.AnalyzeDemo(demoPath, api.AnalyzeDemoOptions{
		Source: constants.DemoSourceFaceIt,
	})
	if err != nil {
		t.Error(err)
	}

	expectedRoundCount := 15
	expectedPlayerCount := 10
	expectedScoreTeamA := 13
	expectedScoreTeamB := 2
	expectedScoreFirstHalfTeamA := 10
	expectedScoreFirstHalfTeamB := 2
	expectedScoreSecondHalfTeamA := 3
	expectedScoreSecondHalfTeamB := 0
	expectedTeamNameA := "Gort Esports"
	expectedTeamNameB := "kriswinenot's party"
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
			StartTick:         24785,
			StartFrame:        25426,
			EndTick:           26812,
			FreezeTimeEndTick: 25553,
			TeamAStartMoney:   4000,
			TeamBStartMoney:   2400, // They are only 3 players
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
			StartTick:         27131,
			StartFrame:        27862,
			EndTick:           28871,
			FreezeTimeEndTick: 27899,
			TeamAStartMoney:   17850,
			TeamBStartMoney:   6650,
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
			StartTick:         29191,
			StartFrame:        29995,
			EndTick:           31598,
			FreezeTimeEndTick: 29959,
			TeamAStartMoney:   22050,
			TeamBStartMoney:   12100,
			TeamAEconomyType:  constants.EconomyTypeSemi,
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
			StartTick:         31918,
			StartFrame:        32871,
			EndTick:           35001,
			FreezeTimeEndTick: 32686,
			TeamAStartMoney:   34400,
			TeamBStartMoney:   11250,
			TeamAEconomyType:  constants.EconomyTypeSemi,
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
			StartTick:         35321,
			StartFrame:        36393,
			EndTick:           38320,
			FreezeTimeEndTick: 36089,
			TeamAStartMoney:   42400,
			TeamBStartMoney:   15150,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeFull,
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
			StartTick:         38640,
			StartFrame:        39801,
			EndTick:           43419,
			FreezeTimeEndTick: 39408,
			TeamAStartMoney:   43250,
			TeamBStartMoney:   14900,
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
			StartTick:         43739,
			StartFrame:        45030,
			EndTick:           50923,
			FreezeTimeEndTick: 44507,
			TeamAStartMoney:   46900,
			TeamBStartMoney:   15450,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonTargetBombed,
			WinnerSide:        common.TeamTerrorists,
			TeamAScore:        6,
			TeamBScore:        1,
			WinnerName:        expectedTeamNameB,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            8,
			StartTick:         51243,
			StartFrame:        52775,
			EndTick:           55035,
			FreezeTimeEndTick: 52011,
			TeamAStartMoney:   36650,
			TeamBStartMoney:   16500,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonTerroristsWin,
			WinnerSide:        common.TeamTerrorists,
			TeamAScore:        6,
			TeamBScore:        2,
			WinnerName:        expectedTeamNameB,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            9,
			StartTick:         55355,
			StartFrame:        57050,
			EndTick:           57716,
			FreezeTimeEndTick: 56123,
			TeamAStartMoney:   24750,
			TeamBStartMoney:   16300,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonCTWin,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        7,
			TeamBScore:        2,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            10,
			StartTick:         58036,
			StartFrame:        59873,
			EndTick:           60535,
			FreezeTimeEndTick: 58804,
			TeamAStartMoney:   18500,
			TeamBStartMoney:   16350,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonCTWin,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        8,
			TeamBScore:        2,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
		{
			Number:            11,
			StartTick:         60855,
			StartFrame:        62735,
			EndTick:           67717,
			FreezeTimeEndTick: 61623,
			TeamAStartMoney:   28200,
			TeamBStartMoney:   15300,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonBombDefused,
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
			StartTick:         68037,
			StartFrame:        70069,
			EndTick:           75645,
			FreezeTimeEndTick: 68805,
			TeamAStartMoney:   41300,
			TeamBStartMoney:   18150,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonBombDefused,
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
			StartTick:         76221,
			StartFrame:        78438,
			EndTick:           79848,
			FreezeTimeEndTick: 77373,
			TeamAStartMoney:   3200, // They are only 4 players
			TeamBStartMoney:   2400,
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
			StartTick:         80168,
			StartFrame:        82527,
			EndTick:           82544,
			FreezeTimeEndTick: 80936,
			TeamAStartMoney:   14750,
			TeamBStartMoney:   6900,
			TeamAEconomyType:  constants.EconomyTypeSemi,
			TeamBEconomyType:  constants.EconomyTypeForceBuy,
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
			StartTick:         82864,
			StartFrame:        85226,
			EndTick:           85627,
			FreezeTimeEndTick: 83632,
			TeamAStartMoney:   24200,
			TeamBStartMoney:   10850,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeForceBuy,
			EndReason:         events.RoundEndReasonTerroristsWin,
			WinnerSide:        common.TeamTerrorists,
			TeamAScore:        13,
			TeamBScore:        2,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamTerrorists,
			TeamBSide:         common.TeamCounterTerrorists,
			TeamAName:         expectedTeamNameA,
			TeamBName:         expectedTeamNameB,
		},
	}

	var players = []fake.FakePlayer{
		{
			SteamID64:     76561198051671412,
			Name:          "Synetic",
			KillCount:     8,
			AssistCount:   2,
			DeathCount:    8,
			Score:         21,
			Team:          match.TeamA,
			HeadshotCount: 4,
			MvpCount:      2,
			UtilityDamage: 180,
		},
		{
			SteamID64:     76561198012776861,
			Name:          "DANDY",
			KillCount:     10,
			AssistCount:   6,
			DeathCount:    6,
			Score:         30,
			Team:          match.TeamA,
			HeadshotCount: 4,
			MvpCount:      4,
			UtilityDamage: 11,
		},
		{
			SteamID64:     76561198123409167,
			Name:          "pApizojT",
			KillCount:     4,
			AssistCount:   4,
			DeathCount:    3,
			Score:         12,
			Team:          match.TeamA,
			HeadshotCount: 2,
			MvpCount:      3,
			UtilityDamage: 0,
		},
		{
			SteamID64:     76561198084521080,
			Name:          "SoSchy",
			KillCount:     10,
			AssistCount:   3,
			DeathCount:    5,
			Score:         23,
			Team:          match.TeamA,
			HeadshotCount: 7,
			MvpCount:      1,
			UtilityDamage: 13,
		},
		{
			SteamID64:     76561198030192915,
			Name:          "Mafiii",
			KillCount:     9,
			AssistCount:   3,
			DeathCount:    3,
			Score:         24,
			Team:          match.TeamA,
			HeadshotCount: 5,
			MvpCount:      3,
			UtilityDamage: 6,
		},
		{
			SteamID64:     76561198971229679,
			Name:          "mampici999",
			KillCount:     0,
			AssistCount:   0,
			DeathCount:    0,
			Score:         0,
			Team:          match.TeamB,
			HeadshotCount: 0,
			MvpCount:      0,
			UtilityDamage: 0,
		},
		{
			SteamID64:     76561199191911268,
			Name:          ".nEZH7",
			KillCount:     0,
			AssistCount:   0,
			DeathCount:    0,
			Score:         0,
			Team:          match.TeamB,
			HeadshotCount: 0,
			MvpCount:      0,
			UtilityDamage: 0,
		},
		{
			SteamID64:     76561198386721272,
			Name:          "Aizputnieks",
			KillCount:     13,
			AssistCount:   1,
			DeathCount:    14,
			Score:         27,
			Team:          match.TeamB,
			HeadshotCount: 4,
			MvpCount:      1,
			UtilityDamage: 90,
		},
		{
			SteamID64:     76561198035979505,
			Name:          "kriswinenot",
			KillCount:     4,
			AssistCount:   4,
			DeathCount:    15,
			Score:         20,
			Team:          match.TeamB,
			HeadshotCount: 1,
			MvpCount:      1,
			UtilityDamage: 325,
		},
		{
			SteamID64:     76561197996341546,
			Name:          "wihu",
			KillCount:     5,
			AssistCount:   4,
			DeathCount:    15,
			Score:         14,
			Team:          match.TeamB,
			HeadshotCount: 3,
			MvpCount:      0,
			UtilityDamage: 175,
		},
	}

	assertion.AssertPlayers(t, match, players)
	assertion.AssertRounds(t, match, rounds)
}
