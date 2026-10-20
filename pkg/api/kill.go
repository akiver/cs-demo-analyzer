package api

import (
	"fmt"

	"github.com/akiver/cs-demo-analyzer/internal/math"
	"github.com/akiver/cs-demo-analyzer/pkg/api/constants"
	"github.com/markus-wa/demoinfocs-golang/v6/pkg/demoinfocs/common"
	"github.com/markus-wa/demoinfocs-golang/v6/pkg/demoinfocs/events"
)

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
		killerTeamName = match.Team(constants.Side(event.Killer.Team)).Name
		isKillerControllingBot = event.Killer.IsControllingBot()
		killerX = event.Killer.Position().X
		killerY = event.Killer.Position().Y
		killerZ = event.Killer.Position().Z
		isKillerAirborne = event.Killer.IsAirborne()
		isKillerBlinded = event.Killer.IsBlinded()
	}

	var isVictimInspectingWeapon bool
	victim := analyzer.match.PlayersBySteamID[event.Victim.SteamID64]
	if victim != nil {
		isVictimInspectingWeapon = analyzer.isInspectingWeapon(victim)
		// if inspection was cancelled and 2+ seconds have passed, consider that the victim is no longer inspecting the weapon
		inspection := analyzer.weaponInspection(victim)
		if isVictimInspectingWeapon && inspection.cancelledTick > -1 && analyzer.secondsHasPassedSinceTick(2, inspection.cancelledTick) {
			isVictimInspectingWeapon = false
		}
	}

	var isTradeKill bool
	for _, kill := range analyzer.match.Kills {
		if kill.RoundNumber == analyzer.currentRound.Number && killerSteamID != 0 {
			if kill.KillerSteamID64 == event.Victim.SteamID64 && !analyzer.secondsHasPassedSinceTick(tradeKillDelaySeconds, kill.Tick) {
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
		assisterTeamName = match.Team(constants.Side(event.Assister.Team)).Name
		isAssisterControllingBot = event.Assister.IsControllingBot()
		assisterX = event.Assister.Position().X
		assisterY = event.Assister.Position().Y
		assisterZ = event.Assister.Position().Z
	}

	distance := event.Distance
	if distance == 0 && event.Killer != nil && event.Victim != nil {
		distance = float32(math.GetDistanceBetweenVectors(event.Killer.Position(), event.Victim.Position()))
	}

	victimPosition := event.Victim.Position()

	return &Kill{
		Frame:                    parser.CurrentFrame(),
		Tick:                     analyzer.currentTick(),
		RoundNumber:              analyzer.currentRound.Number,
		KillerName:               killerName,
		KillerSteamID64:          killerSteamID,
		KillerSide:               constants.Side(killerSide),
		KillerTeamName:           killerTeamName,
		VictimName:               event.Victim.Name,
		VictimSteamID64:          event.Victim.SteamID64,
		VictimSide:               constants.Side(event.Victim.Team),
		VictimTeamName:           match.Team(constants.Side(event.Victim.Team)).Name,
		AssisterName:             assisterName,
		AssisterSteamID64:        assisterSteamID,
		AssisterSide:             constants.Side(assisterSide),
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
