package assertion

import (
	"testing"

	"github.com/akiver/cs-demo-analyzer/pkg/api"
)

func AssertPlayerEconomies(t *testing.T, match *api.Match, economies []api.PlayerEconomy) {
	for _, expectedEconomy := range economies {
		economy := match.GetPlayerEconomyAtRound(expectedEconomy.Name, expectedEconomy.SteamID64, expectedEconomy.RoundNumber)
		if economy == nil {
			t.Fail()
			t.Logf("could not find economy for player %s at round %d", expectedEconomy.Name, expectedEconomy.RoundNumber)
			break
		}

		if economy.RoundNumber != expectedEconomy.RoundNumber {
			t.Errorf("expected player economy round number to be %d but got %d", expectedEconomy.RoundNumber, economy.RoundNumber)
		}
		if economy.SteamID64 != expectedEconomy.SteamID64 {
			t.Errorf("expected SteamID for %s to be %d but got %d", expectedEconomy.Name, expectedEconomy.SteamID64, economy.SteamID64)
		}
		if economy.Name != expectedEconomy.Name {
			t.Errorf("expected player name to be %s but got %s", expectedEconomy.Name, economy.Name)
		}
		if economy.Type != expectedEconomy.Type {
			t.Errorf("expected player economy type to be %s but got %s for player %s at round %d", expectedEconomy.Type, economy.Type, economy.Name, expectedEconomy.RoundNumber)
		}
		if economy.StartMoney != expectedEconomy.StartMoney {
			t.Errorf("expected start money to be %d but got %d for player %s at round %d", expectedEconomy.StartMoney, economy.StartMoney, economy.Name, expectedEconomy.RoundNumber)
		}
		if economy.MoneySpent != expectedEconomy.MoneySpent {
			t.Errorf("expected money spent to be %d but got %d for player %s at round %d", expectedEconomy.MoneySpent, economy.MoneySpent, economy.Name, expectedEconomy.RoundNumber)
		}
		if economy.EquipmentValue != expectedEconomy.EquipmentValue {
			t.Errorf("expected equipment value to be %d but got %d for player %s at round %d", expectedEconomy.EquipmentValue, economy.EquipmentValue, economy.Name, expectedEconomy.RoundNumber)
		}
	}
}
