package api

import (
	"github.com/akiver/cs-demo-analyzer/internal/math"
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/common"
	events "github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/events"
)

type PlayerFlashed struct {
	Tick                    int         `json:"tick"`
	RoundNumber             int         `json:"roundNumber"`
	Duration                float32     `json:"duration"`
	FlashedSteamID64        uint64      `json:"flashedSteamId"`
	FlashedName             string      `json:"flashedName"`
	FlashedSide             common.Team `json:"flashedSide"`
	IsFlashedControllingBot bool        `json:"isFlashedControllingBot"`
	FlasherSteamID64        uint64      `json:"flasherSteamId"`
	FlasherName             string      `json:"flasherName"`
	FlasherSide             common.Team `json:"flasherSide"`
	IsFlasherControllingBot bool        `json:"isFlasherControllingBot"`
	// The flashbang exploded in front of the flashed player (less than 90° from where they were looking).
	// The in-game scoreboard counts an enemy as flashed only in this case.
	IsFlashedFacingFlashbang bool `json:"isFlashedFacingFlashbang"`
}

func isPlayerFacingFlashbang(player *common.Player, projectile *common.GrenadeProjectile) bool {
	if projectile == nil {
		return true
	}

	viewDirection := math.ViewDirection(float64(player.ViewDirectionX()), float64(player.ViewDirectionY()))
	toFlashbang := projectile.Position().Sub(player.PositionEyes())

	return viewDirection.Dot(toFlashbang) > 0
}

func newPlayerFlashed(analyzer *Analyzer, event events.PlayerFlashed) *PlayerFlashed {
	return &PlayerFlashed{
		Tick:                     analyzer.currentTick(),
		RoundNumber:              analyzer.currentRound.Number,
		Duration:                 event.Player.FlashDuration,
		FlashedName:              event.Player.Name,
		FlashedSteamID64:         event.Player.SteamID64,
		FlashedSide:              event.Player.Team,
		IsFlashedControllingBot:  event.Player.IsControllingBot(),
		FlasherName:              event.Attacker.Name,
		FlasherSteamID64:         event.Attacker.SteamID64,
		FlasherSide:              event.Attacker.Team,
		IsFlasherControllingBot:  event.Attacker.IsControllingBot(),
		IsFlashedFacingFlashbang: isPlayerFacingFlashbang(event.Player, event.Projectile),
	}
}
