package assertion

import (
	"testing"

	"github.com/akiver/cs-demo-analyzer/pkg/api"
	"github.com/akiver/cs-demo-analyzer/tests/fake"
)

func AssertRounds(t *testing.T, match *api.Match, rounds []fake.FakeRound) {
	for index, expectedRound := range rounds {
		round := match.Rounds[index]
		if round.Number != expectedRound.Number {
			t.Errorf("expected round number %d but got %d", expectedRound.Number, round.Number)
		}
		if round.OvertimeNumber != expectedRound.OvertimeNumber {
			t.Errorf("expected round overtime number to be %d but got %d", expectedRound.OvertimeNumber, round.OvertimeNumber)
		}
		if round.StartTick != expectedRound.StartTick {
			t.Errorf("expected round start tick to be %d but got %d round number %d", expectedRound.StartTick, round.StartTick, round.Number)
		}
		if round.StartFrame != expectedRound.StartFrame {
			t.Errorf("expected round start frame to be %d but got %d round number %d", expectedRound.StartFrame, round.StartFrame, round.Number)
		}
		if round.FreezeTimeEndTick != expectedRound.FreezeTimeEndTick {
			t.Errorf("expected round freeze time end tick to be %d but got %d round number %d", expectedRound.FreezeTimeEndTick, round.FreezeTimeEndTick, round.Number)
		}
		if round.EndTick != expectedRound.EndTick {
			t.Errorf("expected round end tick to be %d but got %d round number %d", expectedRound.EndTick, round.EndTick, round.Number)
		}
		if expectedRound.EndOfficiallyTick != 0 && round.EndOfficiallyTick != expectedRound.EndOfficiallyTick {
			t.Errorf("expected round end officially tick to be %d but got %d round number %d", expectedRound.EndOfficiallyTick, round.EndOfficiallyTick, round.Number)
		}
		if expectedRound.FreezeTimeEndFrame != 0 && round.FreezeTimeEndFrame != expectedRound.FreezeTimeEndFrame {
			t.Errorf("expected round freeze time end frame to be %d but got %d round number %d", expectedRound.FreezeTimeEndFrame, round.FreezeTimeEndFrame, round.Number)
		}
		if expectedRound.EndFrame != 0 && round.EndFrame != expectedRound.EndFrame {
			t.Errorf("expected round end frame to be %d but got %d round number %d", expectedRound.EndFrame, round.EndFrame, round.Number)
		}
		if expectedRound.EndOfficiallyFrame != 0 && round.EndOfficiallyFrame != expectedRound.EndOfficiallyFrame {
			t.Errorf("expected round end officially frame to be %d but got %d round number %d", expectedRound.EndOfficiallyFrame, round.EndOfficiallyFrame, round.Number)
		}
		if expectedRound.TeamAStartMoney != 0 && round.StartMoneyTeamA() != expectedRound.TeamAStartMoney {
			t.Errorf("expected round start money team A to be %d but got %d round number %d", expectedRound.TeamAStartMoney, round.StartMoneyTeamA(), round.Number)
		}
		if expectedRound.TeamBStartMoney != 0 && round.StartMoneyTeamB() != expectedRound.TeamBStartMoney {
			t.Errorf("expected round start money team B to be %d but got %d round number %d", expectedRound.TeamBStartMoney, round.StartMoneyTeamB(), round.Number)
		}
		if expectedRound.TeamAEquipmentValue != 0 && round.TeamAEquipmentValue != expectedRound.TeamAEquipmentValue {
			t.Errorf("expected round equipment value team A to be %d but got %d round number %d", expectedRound.TeamAEquipmentValue, round.TeamAEquipmentValue, round.Number)
		}
		if expectedRound.TeamBEquipmentValue != 0 && round.TeamBEquipmentValue != expectedRound.TeamBEquipmentValue {
			t.Errorf("expected round equipment value team B to be %d but got %d round number %d", expectedRound.TeamBEquipmentValue, round.TeamBEquipmentValue, round.Number)
		}
		if expectedRound.TeamAEconomyType != "" && round.TeamAEconomyType != expectedRound.TeamAEconomyType {
			t.Errorf("expected round economy type team A (%s) to be %s but got %s round number %d", round.TeamAName, expectedRound.TeamAEconomyType, round.TeamAEconomyType, round.Number)
		}
		if expectedRound.TeamBEconomyType != "" && round.TeamBEconomyType != expectedRound.TeamBEconomyType {
			t.Errorf("expected round economy type team B (%s) to be %s but got %s round number %d", round.TeamBName, expectedRound.TeamBEconomyType, round.TeamBEconomyType, round.Number)
		}
		if expectedRound.EndReason != 0 && round.EndReason != expectedRound.EndReason {
			t.Errorf("expected round end reason to be %d but got %d round number %d", expectedRound.EndReason, round.EndReason, round.Number)
		}
		if expectedRound.WinnerSide != 0 && round.WinnerSide != expectedRound.WinnerSide {
			t.Errorf("expected round winner side to be %d but got %d round number %d", expectedRound.WinnerSide, round.WinnerSide, round.Number)
		}
		if expectedRound.WinnerName != "" && round.WinnerName != expectedRound.WinnerName {
			t.Errorf("expected round winner name to be %s but got %s round number %d", expectedRound.WinnerName, round.WinnerName, round.Number)
		}
		if round.TeamAScore != expectedRound.TeamAScore {
			t.Errorf("expected round score team A to be %d but got %d round number %d", expectedRound.TeamAScore, round.TeamAScore, round.Number)
		}
		if round.TeamBScore != expectedRound.TeamBScore {
			t.Errorf("expected round score team B to be %d but got %d round number %d", expectedRound.TeamBScore, round.TeamBScore, round.Number)
		}
		if expectedRound.TeamASide != 0 && round.TeamASide != expectedRound.TeamASide {
			t.Errorf("expected round side team A to be %d but got %d round number %d", expectedRound.TeamASide, round.TeamASide, round.Number)
		}
		if expectedRound.TeamBSide != 0 && round.TeamBSide != expectedRound.TeamBSide {
			t.Errorf("expected round side team B to be %d but got %d round number %d", expectedRound.TeamBSide, round.TeamBSide, round.Number)
		}
		if expectedRound.TeamAName != "" && expectedRound.TeamAName != round.TeamAName {
			t.Errorf("expected team A name to be %s but got %s round number %d", expectedRound.TeamAName, round.TeamAName, round.Number)
		}
		if expectedRound.TeamBName != "" && expectedRound.TeamBName != round.TeamBName {
			t.Errorf("expected team B name to be %s but got %s round number %d", expectedRound.TeamBName, round.TeamBName, round.Number)
		}
	}
}
