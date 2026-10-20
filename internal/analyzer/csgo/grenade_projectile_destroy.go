package csgo

import (
	"fmt"

	"github.com/akiver/cs-demo-analyzer/pkg/api/constants"
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/common"
)

func newGrenadeProjectileDestroyFromProjectile(analyzer *Analyzer, projectile *common.GrenadeProjectile) *GrenadeProjectileDestroy {
	if projectile == nil {
		fmt.Println("Projectile nil in grenade projectile destroy creation")
		return nil
	}

	thrower := projectile.Thrower
	if thrower == nil {
		fmt.Println("Thrower nil in grenade projectile destroy creation, falling back to owner")
		if projectile.WeaponInstance == nil {
			fmt.Println("Projectile weapon instance nil in grenade projectile destroy creation")
			return nil
		}

		thrower = projectile.WeaponInstance.Owner
		if thrower == nil {
			fmt.Println("Owner nil in grenade projectile destroy creation")
			return nil
		}
	}

	velocity := thrower.Velocity()

	parser := analyzer.parser
	throwerTeam := thrower.Team
	return &GrenadeProjectileDestroy{
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
