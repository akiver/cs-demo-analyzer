package api

import (
	"testing"

	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/common"
)

func TestGetPlayerUserIDUsesSource2PlayerInfoSlot(t *testing.T) {
	analyzer := &Analyzer{
		isSource2: true,
		playerSlotBySteamID64: map[uint64]int{
			76561198336883042: 9,
		},
	}
	player := common.Player{
		SteamID64: 76561198336883042,
		UserID:    0,
	}

	userID := getPlayerUserID(analyzer, player)
	if userID != 9 {
		t.Fatalf("Expected Source2 PlayerInfo slot 9, got %d", userID)
	}
}

func TestGetPlayerUserIDFallsBackToMaskedUserID(t *testing.T) {
	analyzer := &Analyzer{
		isSource2:             true,
		playerSlotBySteamID64: make(map[uint64]int),
	}
	player := common.Player{
		SteamID64: 76561198336883042,
		UserID:    0xff09,
	}

	userID := getPlayerUserID(analyzer, player)
	if userID != 9 {
		t.Fatalf("Expected masked user ID 9, got %d", userID)
	}
}
