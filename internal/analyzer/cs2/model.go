package cs2

import "github.com/akiver/cs-demo-analyzer/pkg/model"

// The data model lives in package model, shared by the CS2 and CS:GO analyzers.
// These aliases let the analyzer refer to the model types unqualified.
type (
	Match                    = model.Match
	MatchJSON                = model.MatchJSON
	MatchAlias               = model.MatchAlias
	Player                   = model.Player
	PlayerJSON               = model.PlayerJSON
	PlayerAlias              = model.PlayerAlias
	Round                    = model.Round
	RoundAlias               = model.RoundAlias
	RoundJSON                = model.RoundJSON
	Team                     = model.Team
	Kill                     = model.Kill
	Damage                   = model.Damage
	Shot                     = model.Shot
	Clutch                   = model.Clutch
	Coach                    = model.Coach
	ChatMessage              = model.ChatMessage
	PlayerEconomy            = model.PlayerEconomy
	PlayerBuy                = model.PlayerBuy
	PlayerFlashed            = model.PlayerFlashed
	PlayerPosition           = model.PlayerPosition
	GrenadeBounce            = model.GrenadeBounce
	GrenadePosition          = model.GrenadePosition
	GrenadeProjectileDestroy = model.GrenadeProjectileDestroy
	DecoyStart               = model.DecoyStart
	FlashbangExplode         = model.FlashbangExplode
	HeGrenadeExplode         = model.HeGrenadeExplode
	SmokeStart               = model.SmokeStart
	InfernoPosition          = model.InfernoPosition
	HostagePosition          = model.HostagePosition
	HostagePickUpStart       = model.HostagePickUpStart
	HostagePickedUp          = model.HostagePickedUp
	HostageRescued           = model.HostageRescued
	BombDefuseStart          = model.BombDefuseStart
	BombDefused              = model.BombDefused
	BombExploded             = model.BombExploded
	BombPlantStart           = model.BombPlantStart
	BombPlanted              = model.BombPlanted
	ChickenPosition          = model.ChickenPosition
	ChickenDeath             = model.ChickenDeath
)
