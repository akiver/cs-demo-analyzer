package assertion

import (
	"testing"

	"github.com/akiver/cs-demo-analyzer/pkg/api"
	"github.com/akiver/cs-demo-analyzer/tests/fake"
)

func AssertPlayers(t *testing.T, match *api.Match, players []fake.FakePlayer) {
	for _, expectedPlayer := range players {
		player := match.PlayersBySteamID[expectedPlayer.SteamID64]
		if player == nil {
			t.Errorf("expected player %s with SteamID %d to be found but got nil", expectedPlayer.Name, expectedPlayer.SteamID64)
			continue
		}
		if player.Name != expectedPlayer.Name {
			t.Errorf("expected player name to be %s but got %s for SteamID %d", expectedPlayer.Name, player.Name, player.SteamID64)
		}
		if player.KillCount() != expectedPlayer.KillCount {
			t.Errorf("expected player %s kill count to be %d but got %d for player %s", player.Name, expectedPlayer.KillCount, player.KillCount(), player.Name)
		}
		if player.AssistCount() != expectedPlayer.AssistCount {
			t.Errorf("expected player %s assist count to be %d but got %d for player %s", player.Name, expectedPlayer.AssistCount, player.AssistCount(), player.Name)
		}
		if player.DeathCount() != expectedPlayer.DeathCount {
			t.Errorf("expected player %s death count to be %d but got %d for player %s", player.Name, expectedPlayer.DeathCount, player.DeathCount(), player.Name)
		}
		if player.Score != expectedPlayer.Score {
			t.Errorf("expected player %s score to be %d but got %d for player %s", player.Name, expectedPlayer.Score, player.Score, player.Name)
		}
		if player.Team != expectedPlayer.Team {
			t.Errorf("expected player %s team to be %s but got %s for player %s", player.Name, expectedPlayer.Team.Letter, player.Team.Letter, player.Name)
		}
		if player.MvpCount != expectedPlayer.MvpCount {
			t.Errorf("expected player MVP to be %d but got %d for player %s", expectedPlayer.MvpCount, player.MvpCount, player.Name)
		}
		if player.HeadshotCount() != expectedPlayer.HeadshotCount {
			t.Errorf("expected player headshot count to be %d but got %d for player %s", expectedPlayer.HeadshotCount, player.HeadshotCount(), player.Name)
		}
		if player.UtilityDamage() != expectedPlayer.UtilityDamage {
			t.Errorf("expected player %s utility damage to be %d but got %d for player %s", player.Name, expectedPlayer.UtilityDamage, player.UtilityDamage(), player.Name)
		}
	}
}
