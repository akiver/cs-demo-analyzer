package cs2

import (
	"fmt"

	"github.com/akiver/cs-demo-analyzer/pkg/api/constants"
	"github.com/markus-wa/demoinfocs-golang/v6/pkg/demoinfocs/common"
)

func newGrenadePositionFromProjectile(analyzer *Analyzer, projectile *common.GrenadeProjectile) *GrenadePosition {
	if projectile.WeaponInstance == nil {
		fmt.Println("Projectile weapon instance nil in grenade projectile position")
		return nil
	}

	thrower := projectile.Thrower
	if thrower == nil {
		fmt.Println("Thrower nil in grenade projectile position, falling back to owner")
		thrower = projectile.WeaponInstance.Owner
		if thrower == nil {
			fmt.Println("Owner nil in grenade projectile position")
			return nil
		}
	}

	velocity := thrower.Velocity()

	parser := analyzer.parser
	throwerTeam := thrower.Team
	return &GrenadePosition{
		Frame:            parser.CurrentFrame(),
		Tick:             analyzer.currentTick(),
		RoundNumber:      analyzer.currentRound.Number,
		GrenadeID:        projectile.WeaponInstance.UniqueID().String(),
		ProjectileID:     projectile.UniqueID().String(),
		GrenadeName:      equipmentToWeaponName[projectile.WeaponInstance.Type],
		X:                projectile.Position().X,
		Y:                projectile.Position().Y,
		Z:                projectile.Position().Z,
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
