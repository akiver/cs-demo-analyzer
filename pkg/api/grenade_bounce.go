package api

import (
	"fmt"

	"github.com/akiver/cs-demo-analyzer/pkg/api/constants"
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/common"
)

type GrenadeBounce struct {
	Frame            int                  `json:"frame"`
	Tick             int                  `json:"tick"`
	RoundNumber      int                  `json:"roundNumber"`
	GrenadeID        string               `json:"grenadeId"`
	ProjectileID     int64                `json:"projectileId"`
	GrenadeName      constants.WeaponName `json:"grenadeName"`
	X                float64              `json:"x"`
	Y                float64              `json:"y"`
	Z                float64              `json:"z"`
	ThrowerSteamID64 uint64               `json:"throwerSteamId"`
	ThrowerName      string               `json:"throwerName"`
	ThrowerSide      common.Team          `json:"throwerSide"`
	ThrowerTeamName  string               `json:"throwerTeamName"`
	ThrowerVelocityX float64              `json:"throwerVelocityX"`
	ThrowerVelocityY float64              `json:"throwerVelocityY"`
	ThrowerVelocityZ float64              `json:"throwerVelocityZ"`
	ThrowerPitch     float32              `json:"throwerPitch"`
	ThrowerYaw       float32              `json:"throwerYaw"`
}

func newGrenadeBounceFromProjectile(analyzer *Analyzer, projectile *common.GrenadeProjectile) *GrenadeBounce {
	if projectile == nil {
		fmt.Println("Projectile nil in grenade projectile bounce event")
		return nil
	}

	if projectile.WeaponInstance == nil {
		fmt.Println("Projectile weapon instance nil in grenade projectile bounce event")
		return nil
	}

	thrower := projectile.Thrower
	if thrower == nil {
		fmt.Println("Thrower nil in grenade projectile bounce event, falling back to owner")
		thrower = projectile.WeaponInstance.Owner
		if thrower == nil {
			fmt.Println("Owner nil in grenade projectile bounce event")
			return nil
		}
	}

	velocity := thrower.Velocity()

	parser := analyzer.parser
	throwerTeam := thrower.Team
	return &GrenadeBounce{
		Frame:            parser.CurrentFrame(),
		Tick:             analyzer.currentTick(),
		RoundNumber:      analyzer.currentRound.Number,
		GrenadeID:        projectile.WeaponInstance.UniqueID2().String(),
		ProjectileID:     projectile.UniqueID(),
		GrenadeName:      equipmentToWeaponName[projectile.WeaponInstance.Type],
		X:                projectile.Position().X,
		Y:                projectile.Position().Y,
		Z:                projectile.Position().Z,
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
