package cs2

import (
	"fmt"

	"github.com/akiver/cs-demo-analyzer/pkg/api/constants"

	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/events"
)

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
	projectileID := findGrenadeProjectileID(analyzer, grenade.UniqueID(), event.Position)

	velocity := thrower.Velocity()

	return &SmokeStart{
		Frame:            parser.CurrentFrame(),
		Tick:             analyzer.currentTick(),
		RoundNumber:      analyzer.currentRound.Number,
		GrenadeID:        grenade.UniqueID().String(),
		ProjectileID:     projectileID,
		X:                event.Position.X,
		Y:                event.Position.Y,
		Z:                event.Position.Z,
		ThrowerSteamID64: thrower.SteamID64,
		ThrowerName:      thrower.Name,
		ThrowerSide:      constants.Side(throwerTeam),
		ThrowerTeamName:  analyzer.match.Team(constants.Side(throwerTeam)).Name,
		ThrowerVelocityX: velocity.X,
		ThrowerVelocityY: velocity.Y,
		ThrowerVelocityZ: velocity.Z,
		ThrowerYaw:       thrower.ViewDirectionX(),
		ThrowerPitch:     thrower.ViewDirectionY(),
	}

}
