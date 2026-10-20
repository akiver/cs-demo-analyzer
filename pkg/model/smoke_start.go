package model

import (
	"github.com/akiver/cs-demo-analyzer/pkg/api/constants"
)

type SmokeStart struct {
	Frame            int            `json:"frame"`
	Tick             int            `json:"tick"`
	RoundNumber      int            `json:"roundNumber"`
	GrenadeID        string         `json:"grenadeId"`
	ProjectileID     string         `json:"projectileId"`
	X                float64        `json:"x"`
	Y                float64        `json:"y"`
	Z                float64        `json:"z"`
	ThrowerSteamID64 uint64         `json:"throwerSteamId"`
	ThrowerName      string         `json:"throwerName"`
	ThrowerSide      constants.Side `json:"throwerSide"`
	ThrowerTeamName  string         `json:"throwerTeamName"`
	ThrowerVelocityX float64        `json:"throwerVelocityX"`
	ThrowerVelocityY float64        `json:"throwerVelocityY"`
	ThrowerVelocityZ float64        `json:"throwerVelocityZ"`
	ThrowerPitch     float32        `json:"throwerPitch"`
	ThrowerYaw       float32        `json:"throwerYaw"`
}
