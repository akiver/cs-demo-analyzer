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

// Demo with a surrender, bots and contains 1 tactical timeout.
func TestValve_Match730_003408404295698088038_1541485657_202_Mirage(t *testing.T) {
	demoName := "valve_match730_003408404295698088038_1541485657_202_mirage"
	demoPath := testsutils.GetDemoPath("csgo", demoName)
	match, err := api.AnalyzeDemo(demoPath, api.AnalyzeDemoOptions{
		Source: constants.DemoSourceValve,
	})
	if err != nil {
		t.Error(err)
	}

	expectedRoundCount := 16
	expectedPlayerCount := 10
	expectedTeamAScore := 13
	expectedTeamBScore := 3
	expectedScoreFirstHalfTeamA := 12
	expectedScoreFirstHalfTeamB := 3
	expectedScoreSecondHalfTeamA := 1
	expectedScoreSecondHalfTeamB := 0
	expectedTeamNameA := "Team A"
	expectedTeamNameB := "Team B"
	expectedWinnerName := expectedTeamNameA
	expectedMaxRounds := 30

	if match.TeamA.Name != expectedTeamNameA {
		t.Errorf("expected team name A to be %s got %s", expectedTeamNameA, match.TeamA.Name)
	}
	if match.TeamB.Name != expectedTeamNameB {
		t.Errorf("expected team name B to be %s got %s", expectedTeamNameB, match.TeamB.Name)
	}
	if len(match.Rounds) != expectedRoundCount {
		t.Errorf("expected %d rounds but got %d", expectedRoundCount, len(match.Rounds))
	}
	if len(match.Players()) != expectedPlayerCount {
		t.Errorf("expected %d players but got %d", expectedPlayerCount, len(match.Players()))
	}
	if match.TeamA.Score != expectedTeamAScore {
		t.Errorf("expected teamScTeamBScore A to be %d got %d", expectedTeamAScore, match.TeamA.Score)
	}
	if match.TeamB.Score != expectedTeamBScore {
		t.Errorf("expected score team B to be %d got %d", expectedTeamBScore, match.TeamB.Score)
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
			StartTick:         6606,
			StartFrame:        3307,
			EndTick:           11775,
			FreezeTimeEndTick: 7561,
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
		},
		{
			Number:            2,
			StartTick:         12225,
			StartFrame:        6113,
			EndTick:           20545,
			FreezeTimeEndTick: 13185,
			TeamAStartMoney:   18450,
			TeamBStartMoney:   11400,
			TeamAEconomyType:  constants.EconomyTypeSemi,
			TeamBEconomyType:  constants.EconomyTypeEco,
			EndReason:         events.RoundEndReasonTargetSaved,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        2,
			TeamBScore:        0,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
		},
		{
			Number:            3,
			StartTick:         20993,
			StartFrame:        10492,
			FreezeTimeEndTick: 21953,
			EndTick:           26111,
			TeamAStartMoney:   21450,
			TeamBStartMoney:   18600,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeSemi,
			EndReason:         events.RoundEndReasonTerroristsWin,
			WinnerSide:        common.TeamTerrorists,
			TeamAScore:        2,
			TeamBScore:        1,
			WinnerName:        expectedTeamNameB,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
		},
		{
			Number:            4,
			StartTick:         26565,
			StartFrame:        13274,
			FreezeTimeEndTick: 27521,
			EndTick:           31541,
			TeamAStartMoney:   15350,
			TeamBStartMoney:   20750,
			TeamAEconomyType:  constants.EconomyTypeSemi,
			TeamBEconomyType:  constants.EconomyTypeSemi,
			EndReason:         events.RoundEndReasonCTWin,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        3,
			TeamBScore:        1,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
		},
		{
			Number:            5,
			StartTick:         31994,
			StartFrame:        15987,
			FreezeTimeEndTick: 32950,
			EndTick:           36401,
			TeamAStartMoney:   20250,
			TeamBStartMoney:   17100,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeSemi,
			EndReason:         events.RoundEndReasonTerroristsWin,
			WinnerSide:        common.TeamTerrorists,
			TeamAScore:        3,
			TeamBScore:        2,
			WinnerName:        expectedTeamNameB,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
		},
		{
			Number:            6,
			StartTick:         41402,
			StartFrame:        20686,
			FreezeTimeEndTick: 42362,
			EndTick:           49038,
			TeamAStartMoney:   12700,
			TeamBStartMoney:   27450,
			TeamAEconomyType:  constants.EconomyTypeSemi,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonBombDefused,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        4,
			TeamBScore:        2,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
		},
		{
			Number:            7,
			StartTick:         49486,
			StartFrame:        24726,
			FreezeTimeEndTick: 50447,
			EndTick:           54473,
			TeamAStartMoney:   22850,
			TeamBStartMoney:   29900,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonCTWin,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        5,
			TeamBScore:        2,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
		},
		{
			Number:            8,
			StartTick:         54923,
			StartFrame:        27442,
			FreezeTimeEndTick: 55884,
			EndTick:           59658,
			TeamAStartMoney:   18700,
			TeamBStartMoney:   17100,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeSemi,
			EndReason:         events.RoundEndReasonCTWin,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        6,
			TeamBScore:        2,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
		},
		{
			Number:            9,
			StartTick:         60112,
			StartFrame:        30032,
			FreezeTimeEndTick: 61068,
			EndTick:           64566,
			TeamAStartMoney:   28850,
			TeamBStartMoney:   16400,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeSemi,
			EndReason:         events.RoundEndReasonTerroristsWin,
			WinnerSide:        common.TeamTerrorists,
			TeamAScore:        6,
			TeamBScore:        3,
			WinnerName:        expectedTeamNameB,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
		},
		{
			Number:            10,
			StartTick:         65016,
			StartFrame:        32483,
			FreezeTimeEndTick: 65976,
			EndTick:           69814,
			TeamAStartMoney:   31950,
			TeamBStartMoney:   24000,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonCTWin,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        7,
			TeamBScore:        3,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
		},
		{
			Number:            11,
			StartTick:         70264,
			StartFrame:        35106,
			FreezeTimeEndTick: 71224,
			EndTick:           77126,
			TeamAStartMoney:   21000,
			TeamBStartMoney:   21150,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeSemi,
			EndReason:         events.RoundEndReasonBombDefused,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        8,
			TeamBScore:        3,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
		},
		{
			Number:            12,
			StartTick:         77574,
			StartFrame:        38758,
			FreezeTimeEndTick: 78534,
			EndTick:           81024,
			TeamAStartMoney:   27700,
			TeamBStartMoney:   19450,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeSemi,
			EndReason:         events.RoundEndReasonCTWin,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        9,
			TeamBScore:        3,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
		},
		{
			Number:            13,
			StartTick:         81477,
			StartFrame:        40706,
			FreezeTimeEndTick: 82433,
			EndTick:           85239,
			TeamAStartMoney:   27550,
			TeamBStartMoney:   16500,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeSemi,
			EndReason:         events.RoundEndReasonCTWin,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        10,
			TeamBScore:        3,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
		},
		{
			Number:            14,
			StartTick:         85687,
			StartFrame:        42807,
			FreezeTimeEndTick: 86647,
			EndTick:           89555,
			TeamAStartMoney:   39400,
			TeamBStartMoney:   21500,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeFull,
			EndReason:         events.RoundEndReasonCTWin,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        11,
			TeamBScore:        3,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
		},
		{
			Number:            15,
			StartTick:         90011,
			StartFrame:        44963,
			FreezeTimeEndTick: 90967,
			EndTick:           93492,
			TeamAStartMoney:   54000,
			TeamBStartMoney:   20400,
			TeamAEconomyType:  constants.EconomyTypeFull,
			TeamBEconomyType:  constants.EconomyTypeSemi,
			EndReason:         events.RoundEndReasonCTWin,
			WinnerSide:        common.TeamCounterTerrorists,
			TeamAScore:        12,
			TeamBScore:        3,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamCounterTerrorists,
			TeamBSide:         common.TeamTerrorists,
		},
		{
			Number:            16,
			StartTick:         94452,
			StartFrame:        47182,
			FreezeTimeEndTick: 94949,
			EndTick:           94949,
			TeamAStartMoney:   4000,
			TeamBStartMoney:   4000,
			TeamAEconomyType:  constants.EconomyTypePistol,
			TeamBEconomyType:  constants.EconomyTypePistol,
			EndReason:         events.RoundEndReasonCTSurrender,
			WinnerSide:        common.TeamTerrorists,
			TeamAScore:        13,
			TeamBScore:        3,
			WinnerName:        expectedTeamNameA,
			TeamASide:         common.TeamTerrorists,
			TeamBSide:         common.TeamCounterTerrorists,
		},
	}

	var players = []fake.FakePlayer{
		{
			SteamID64:     76561199044985101,
			Name:          "CORONA VIRUS",
			KillCount:     15,
			AssistCount:   2,
			DeathCount:    8,
			Score:         32,
			Team:          match.TeamA,
			MvpCount:      3,
			HeadshotCount: 8,
			UtilityDamage: 69,
		},
		{
			SteamID64:     76561199034633104,
			Name:          "h0nda",
			KillCount:     13,
			AssistCount:   2,
			DeathCount:    5,
			Score:         27,
			Team:          match.TeamA,
			MvpCount:      2,
			HeadshotCount: 10,
			UtilityDamage: 33,
		},
		{
			SteamID64:     76561198300094916,
			Name:          "🚩SIX1600🚩",
			KillCount:     8,
			AssistCount:   1,
			DeathCount:    14,
			Score:         19,
			Team:          match.TeamB,
			MvpCount:      0,
			HeadshotCount: 3,
			UtilityDamage: 0,
		},
		{
			SteamID64:     76561198874869112,
			Name:          "Landers",
			KillCount:     6,
			AssistCount:   0,
			DeathCount:    14,
			Score:         12,
			Team:          match.TeamB,
			MvpCount:      1,
			HeadshotCount: 3,
			UtilityDamage: 0,
		},
		{
			SteamID64:     76561198355926583,
			Name:          "𝙺𝚒𝚠𝚒𝚜𝚔𝚢",
			KillCount:     5,
			AssistCount:   1,
			DeathCount:    10,
			Score:         14,
			Team:          match.TeamB,
			MvpCount:      0,
			HeadshotCount: 2,
			UtilityDamage: 32,
		},
		{
			SteamID64:     76561198374398001,
			Name:          "lucao the majest1c",
			KillCount:     6,
			AssistCount:   0,
			DeathCount:    10,
			Score:         16,
			Team:          match.TeamB,
			MvpCount:      2,
			HeadshotCount: 4,
			UtilityDamage: 28,
		},
		{
			SteamID64:     76561198245758058,
			Name:          "ZyWoO",
			KillCount:     3,
			AssistCount:   3,
			DeathCount:    12,
			Score:         13,
			Team:          match.TeamB,
			MvpCount:      0,
			HeadshotCount: 0,
			UtilityDamage: 1,
		},
	}

	var playerEconomies = []api.PlayerEconomy{
		{
			RoundNumber:    1,
			SteamID64:      76561199044985101,
			Name:           "CORONA VIRUS",
			StartMoney:     800,
			MoneySpent:     650,
			EquipmentValue: 850,
			Type:           constants.EconomyTypePistol,
		},
		{
			RoundNumber:    2,
			SteamID64:      76561199044985101,
			Name:           "CORONA VIRUS",
			StartMoney:     4300,
			MoneySpent:     3450,
			EquipmentValue: 3750,
			Type:           constants.EconomyTypeSemi,
		},
		{
			RoundNumber:    3,
			SteamID64:      76561199044985101,
			Name:           "CORONA VIRUS",
			StartMoney:     4700,
			MoneySpent:     1800,
			EquipmentValue: 4350,
			Type:           constants.EconomyTypeSemi,
		},
		{
			RoundNumber:    4,
			SteamID64:      76561199044985101,
			Name:           "CORONA VIRUS",
			StartMoney:     4900,
			MoneySpent:     1000,
			EquipmentValue: 1000,
			Type:           constants.EconomyTypeEco,
		},
	}

	assertion.AssertPlayers(t, match, players)
	assertion.AssertRounds(t, match, rounds)
	assertion.AssertPlayerEconomies(t, match, playerEconomies)
}
