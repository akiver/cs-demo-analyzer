package api

import (
	"fmt"

	"github.com/akiver/cs-demo-analyzer/internal/math"
	"github.com/akiver/cs-demo-analyzer/pkg/api/constants"
	"github.com/golang/geo/r3"
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/common"
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/events"
)

type Kill struct {
	Frame                    int                  `json:"frame"`
	Tick                     int                  `json:"tick"`
	RoundNumber              int                  `json:"roundNumber"`
	WeaponType               constants.WeaponType `json:"weaponType"`
	WeaponName               constants.WeaponName `json:"weaponName"`
	KillerName               string               `json:"killerName"`
	KillerSteamID64          uint64               `json:"killerSteamId"`
	KillerSide               common.Team          `json:"killerSide"`
	KillerTeamName           string               `json:"killerTeamName"`
	KillerX                  float64              `json:"killerX"`
	KillerY                  float64              `json:"killerY"`
	KillerZ                  float64              `json:"killerZ"`
	IsKillerAirborne         bool                 `json:"is_killer_airborne"`
	IsKillerBlinded          bool                 `json:"is_killer_blinded"`
	IsKillerControllingBot   bool                 `json:"isKillerControllingBot"`
	VictimName               string               `json:"victimName"`
	VictimSteamID64          uint64               `json:"victimSteamId"`
	VictimSide               common.Team          `json:"victimSide"`
	VictimTeamName           string               `json:"victimTeamName"`
	VictimX                  float64              `json:"victimX"`
	VictimY                  float64              `json:"victimY"`
	VictimZ                  float64              `json:"victimZ"`
	IsVictimAirborne         bool                 `json:"is_victim_airborne"`
	IsVictimBlinded          bool                 `json:"is_victim_blinded"`
	IsVictimControllingBot   bool                 `json:"isVictimControllingBot"`
	IsVictimInspectingWeapon bool                 `json:"isVictimInspectingWeapon"`
	AssisterName             string               `json:"assisterName"`
	AssisterSteamID64        uint64               `json:"assisterSteamId"`
	AssisterSide             common.Team          `json:"assisterSide"`
	AssisterTeamName         string               `json:"assisterTeamName"`
	AssisterX                float64              `json:"assisterX"`
	AssisterY                float64              `json:"assisterY"`
	AssisterZ                float64              `json:"assisterZ"`
	IsAssisterControllingBot bool                 `json:"isAssisterControllingBot"`
	IsHeadshot               bool                 `json:"isHeadshot"`
	PenetratedObjects        int                  `json:"penetratedObjects"`
	IsAssistedFlash          bool                 `json:"isAssistedFlash"`
	IsThroughSmoke           bool                 `json:"isThroughSmoke"`
	IsNoScope                bool                 `json:"isNoScope"`
	IsTradeKill              bool                 `json:"isTradeKill"`  // The attacker did a trade kill
	IsTradeDeath             bool                 `json:"isTradeDeath"` // The victim did a trade death
	Distance                 float32              `json:"distance"`
}

func (kill *Kill) IsSuicide() bool {
	return kill.KillerSteamID64 == kill.VictimSteamID64
}

func (kill *Kill) IsTeamKill() bool {
	return kill.KillerSide == kill.VictimSide
}

func newKillFromGameEvent(analyzer *Analyzer, event events.Kill) *Kill {
	if event.Weapon == nil {
		fmt.Println("Player kill event without weapon occurred")
		return nil
	}
	if event.Victim == nil {
		fmt.Println("Player kill event without victim occurred")
		return nil
	}
	parser := analyzer.parser
	match := analyzer.match
	var killerName string = "World"
	var killerSteamID uint64
	var killerSide common.Team
	var killerTeamName string
	var isKillerControllingBot bool
	var isKillerAirborne bool
	var isKillerBlinded bool
	var killerX float64
	var killerY float64
	var killerZ float64
	if event.Killer != nil {
		killerName = event.Killer.Name
		killerSteamID = event.Killer.SteamID64
		killerSide = event.Killer.Team
		killerTeamName = match.Team(event.Killer.Team).Name
		isKillerControllingBot = event.Killer.IsControllingBot()
		killerX = event.Killer.Position().X
		killerY = event.Killer.Position().Y
		killerZ = event.Killer.Position().Z
		isKillerAirborne = event.Killer.IsAirborne()
		isKillerBlinded = event.Killer.IsBlinded()
	}

	var isVictimInspectingWeapon bool
	if analyzer.isSource2 {
		isVictimInspectingWeapon = event.Victim.PlayerPawnEntity().PropertyValueMust("m_pWeaponServices.m_bIsLookingAtWeapon").BoolVal()
	} else if event.Victim.Entity != nil {
		isVictimInspectingWeapon = event.Victim.Entity.PropertyValueMust("m_bIsLookingAtWeapon").BoolVal()
	}

	var isTradeKill bool
	for _, kill := range analyzer.match.Kills {
		if kill.RoundNumber == analyzer.currentRound.Number && killerSteamID != 0 {
			if kill.KillerSteamID64 == event.Victim.SteamID64 && analyzer.secondsHasPassedSinceTick(tradeKillDelaySeconds, kill.Tick) {
				isTradeKill = true
				kill.IsTradeDeath = true
			}
		}
	}

	var assisterName string
	var assisterSteamID uint64
	var assisterSide common.Team
	var assisterTeamName string
	var isAssisterControllingBot bool
	var assisterX float64
	var assisterY float64
	var assisterZ float64
	if event.Assister != nil {
		assisterName = event.Assister.Name
		assisterSteamID = event.Assister.SteamID64
		assisterSide = event.Assister.Team
		assisterTeamName = match.Team(event.Assister.Team).Name
		isAssisterControllingBot = event.Assister.IsControllingBot()
		assisterX = event.Assister.Position().X
		assisterY = event.Assister.Position().Y
		assisterZ = event.Assister.Position().Z
	}

	distance := event.Distance
	if distance == 0 && event.Killer != nil && event.Victim != nil {
		var killerPosition r3.Vector
		var victimPosition r3.Vector
		if analyzer.isSource2 {
			killerPosition = event.Killer.Position()
			victimPosition = event.Victim.Position()
		} else {
			killerPosition = event.Killer.PositionEyes()
			victimPosition = event.Victim.PositionEyes()
		}

		distance = float32(math.GetDistanceBetweenVectors(killerPosition, victimPosition))
	}

	victimPosition := event.Victim.Position()

	return &Kill{
		Frame:                    parser.CurrentFrame(),
		Tick:                     analyzer.currentTick(),
		RoundNumber:              analyzer.currentRound.Number,
		KillerName:               killerName,
		KillerSteamID64:          killerSteamID,
		KillerSide:               killerSide,
		KillerTeamName:           killerTeamName,
		VictimName:               event.Victim.Name,
		VictimSteamID64:          event.Victim.SteamID64,
		VictimSide:               event.Victim.Team,
		VictimTeamName:           match.Team(event.Victim.Team).Name,
		AssisterName:             assisterName,
		AssisterSteamID64:        assisterSteamID,
		AssisterSide:             assisterSide,
		AssisterTeamName:         assisterTeamName,
		IsHeadshot:               event.IsHeadshot,
		PenetratedObjects:        event.PenetratedObjects,
		WeaponName:               equipmentToWeaponName[event.Weapon.Type],
		WeaponType:               getEquipmentWeaponType(*event.Weapon),
		IsKillerControllingBot:   isKillerControllingBot,
		IsVictimControllingBot:   event.Victim.IsControllingBot(),
		IsAssisterControllingBot: isAssisterControllingBot,
		KillerX:                  killerX,
		KillerY:                  killerY,
		KillerZ:                  killerZ,
		IsKillerAirborne:         isKillerAirborne,
		IsKillerBlinded:          isKillerBlinded,
		VictimX:                  victimPosition.X,
		VictimY:                  victimPosition.Y,
		VictimZ:                  victimPosition.Z,
		IsVictimAirborne:         event.Victim.IsAirborne(),
		IsVictimBlinded:          event.Victim.IsBlinded(),
		IsVictimInspectingWeapon: isVictimInspectingWeapon,
		AssisterX:                assisterX,
		AssisterY:                assisterY,
		AssisterZ:                assisterZ,
		IsTradeKill:              isTradeKill,
		IsAssistedFlash:          event.AssistedFlash,
		IsThroughSmoke:           event.ThroughSmoke,
		IsNoScope:                event.NoScope,
		Distance:                 distance,
	}
}
