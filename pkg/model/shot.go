package model

import (
	"github.com/akiver/cs-demo-analyzer/pkg/api/constants"
)

type Shot struct {
	Frame                  int                  `json:"frame"`
	Tick                   int                  `json:"tick"`
	RoundNumber            int                  `json:"roundNumber"`
	WeaponName             constants.WeaponName `json:"weaponName"`
	WeaponID               string               `json:"weaponId"`
	ProjectileID           string               `json:"projectileId"` // Available only for grenades
	X                      float64              `json:"x"`
	Y                      float64              `json:"y"`
	Z                      float64              `json:"z"`
	PlayerName             string               `json:"playerName"`
	PlayerSteamID64        uint64               `json:"playerSteamId"`
	PlayerTeamName         string               `json:"playerTeamName"`
	PlayerSide             constants.Side       `json:"playerSide"`
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
