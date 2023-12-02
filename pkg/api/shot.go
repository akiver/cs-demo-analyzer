package api

import (
	"github.com/akiver/cs-demo-analyzer/pkg/api/constants"
	"github.com/golang/geo/r3"
	common "github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/common"
	events "github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/events"
	st "github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/sendtables"
)

type Shot struct {
	Frame                  int                  `json:"frame"`
	Tick                   int                  `json:"tick"`
	RoundNumber            int                  `json:"roundNumber"`
	WeaponName             constants.WeaponName `json:"weaponName"`
	WeaponID               string               `json:"weaponId"`
	ProjectileID           int64                `json:"projectileId"` // Available only for grenades
	X                      float64              `json:"x"`
	Y                      float64              `json:"y"`
	Z                      float64              `json:"z"`
	PlayerName             string               `json:"playerName"`
	PlayerSteamID64        uint64               `json:"playerSteamId"`
	PlayerTeamName         string               `json:"playerTeamName"`
	PlayerSide             common.Team          `json:"playerSide"`
	IsPlayerControllingBot bool                 `json:"isPlayerControllingBot"`
	PlayerVelocityX        float64              `json:"playerVelocityX"`
	PlayerVelocityY        float64              `json:"playerVelocityY"`
	PlayerVelocityZ        float64              `json:"playerVelocityZ"`
	Yaw                    float32              `json:"yaw"`
	Pitch                  float32              `json:"pitch"`
	RecoilIndex            float32              `json:"recoilIndex"`
	AimPunchAngleX         float64              `json:"aimPunchAngleX"`
	AimPunchAngleY         float64              `json:"aimPunchAngleY"`
	ViewPunchAngleX        float64              `json:"viewPunchAngleX"`
	ViewPunchAngleY        float64              `json:"viewPunchAngleY"`
}

func newShot(analyzer *Analyzer, event events.WeaponFire) *Shot {
	shooter := event.Shooter
	if shooter == nil {
		return nil
	}

	var weaponEntity st.Entity
	activeWeapon := shooter.ActiveWeapon()
	if activeWeapon != nil {
		weaponEntity = activeWeapon.Entity
	} else if event.Weapon.Entity != nil {
		weaponEntity = event.Weapon.Entity
	}

	var recoilIndex float32
	if weaponEntity != nil {
		if prop, exists := weaponEntity.PropertyValue("m_flRecoilIndex"); exists {
			recoilIndex = prop.Float()
		}
	}

	var aimPunchAngle r3.Vector
	var viewPunchAngle r3.Vector
	if analyzer.isSource2 {
		pawnEntity := shooter.PlayerPawnEntity()
		aimPunchAngle = pawnEntity.PropertyValueMust("m_aimPunchAngle").R3Vec()
		// This prop may not exist with demos coming from the early CS2 limited test
		if prop, exists := pawnEntity.PropertyValue("m_pCameraServices.m_vecCsViewPunchAngle"); exists {
			viewPunchAngle = prop.R3Vec()
		}
	} else {
		aimPunchAngle = shooter.Entity.PropertyValueMust("localdata.m_Local.m_aimPunchAngle").R3Vec()
		viewPunchAngle = shooter.Entity.PropertyValueMust("localdata.m_Local.m_viewPunchAngle").R3Vec()
	}

	velocity := shooter.Velocity()

	return &Shot{
		Frame:                  analyzer.parser.CurrentFrame(),
		Tick:                   analyzer.currentTick(),
		RoundNumber:            analyzer.currentRound.Number,
		WeaponName:             equipmentToWeaponName[event.Weapon.Type],
		WeaponID:               event.Weapon.UniqueID2().String(),
		X:                      shooter.Position().X,
		Y:                      shooter.Position().Y,
		Z:                      shooter.Position().Z,
		PlayerName:             shooter.Name,
		PlayerSteamID64:        shooter.SteamID64,
		PlayerTeamName:         analyzer.match.Team(shooter.Team).Name,
		PlayerSide:             shooter.Team,
		IsPlayerControllingBot: shooter.IsControllingBot(),
		PlayerVelocityX:        velocity.X,
		PlayerVelocityY:        velocity.Y,
		PlayerVelocityZ:        velocity.Z,
		Yaw:                    shooter.ViewDirectionX(),
		Pitch:                  shooter.ViewDirectionY(),
		RecoilIndex:            recoilIndex,
		AimPunchAngleX:         aimPunchAngle.X,
		AimPunchAngleY:         aimPunchAngle.Y,
		ViewPunchAngleX:        viewPunchAngle.X,
		ViewPunchAngleY:        viewPunchAngle.Y,
	}
}
