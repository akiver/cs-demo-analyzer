package api

import (
	"fmt"

	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/common"
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/events"
)

type SmokeStart struct {
	Frame            int         `json:"frame"`
	Tick             int         `json:"tick"`
	RoundNumber      int         `json:"roundNumber"`
	GrenadeID        string      `json:"grenadeId"`
	ProjectileID     int64       `json:"projectileId"`
	X                float64     `json:"x"`
	Y                float64     `json:"y"`
	Z                float64     `json:"z"`
	ThrowerSteamID64 uint64      `json:"throwerSteamId"`
	ThrowerName      string      `json:"throwerName"`
	ThrowerSide      common.Team `json:"throwerSide"`
	ThrowerTeamName  string      `json:"throwerTeamName"`
	ThrowerVelocityX float64     `json:"throwerVelocityX"`
	ThrowerVelocityY float64     `json:"throwerVelocityY"`
	ThrowerVelocityZ float64     `json:"throwerVelocityZ"`
	ThrowerPitch     float32     `json:"throwerPitch"`
	ThrowerYaw       float32     `json:"throwerYaw"`
}

func newSmokeStartFromGameEvent(analyzer *Analyzer, event events.SmokeStart) *SmokeStart {
	grenade := event.Grenade
	if grenade == nil {
		fmt.Println("Grenade nil in smoke start event")
		return nil
	}

	thrower := event.Thrower
	if thrower == nil {
		fmt.Println("Thrower nil in smoke start event")
		return nil
	}

	throwerTeam := thrower.Team
	parser := analyzer.parser
	var projectileID int64
	for _, projectile := range parser.GameState().GrenadeProjectiles() {
		if projectile.WeaponInstance.UniqueID2() == grenade.UniqueID2() {
			projectileID = projectile.UniqueID()
			break
		}
	}

	velocity := thrower.Velocity()

	return &SmokeStart{
		Frame:            parser.CurrentFrame(),
		Tick:             analyzer.currentTick(),
		RoundNumber:      analyzer.currentRound.Number,
		GrenadeID:        grenade.UniqueID2().String(),
		ProjectileID:     projectileID,
		X:                event.Position.X,
		Y:                event.Position.Y,
		Z:                event.Position.Z,
		ThrowerSteamID64: thrower.SteamID64,
		ThrowerName:      thrower.Name,
		ThrowerSide:      throwerTeam,
		ThrowerTeamName:  analyzer.match.Team(throwerTeam).Name,
		ThrowerVelocityX: velocity.X,
		ThrowerVelocityY: velocity.Y,
		ThrowerVelocityZ: velocity.Z,
		ThrowerYaw:       thrower.ViewDirectionX(),
		ThrowerPitch:     thrower.ViewDirectionY(),
	}

}
