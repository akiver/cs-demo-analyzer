package cs2

import (
	"fmt"

	"github.com/akiver/cs-demo-analyzer/internal/math"
	"github.com/akiver/cs-demo-analyzer/pkg/api/constants"
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/common"
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/events"
)

func newDamageFromGameEvent(analyzer *Analyzer, event events.PlayerHurt) *Damage {
	if event.Weapon == nil {
		fmt.Println("Player hurt event without weapon occurred")
		return nil
	}
	parser := analyzer.parser
	match := analyzer.match
	attackerSteamID := uint64(0)
	attackerSide := common.TeamUnassigned
	attackerTeamName := "World"
	isAttackerControllingBot := false
	if event.Attacker != nil {
		attackerSteamID = event.Attacker.SteamID64
		attackerSide = event.Attacker.Team
		attackerTeamName = match.Team(constants.Side(event.Attacker.Team)).Name
		isAttackerControllingBot = event.Attacker.IsControllingBot()
	}

	return &Damage{
		RoundNumber:              analyzer.currentRound.Number,
		Frame:                    parser.CurrentFrame(),
		Tick:                     analyzer.currentTick(),
		HealthDamage:             math.Max(0, event.HealthDamageTaken),
		ArmorDamage:              math.Max(0, event.ArmorDamageTaken),
		VictimHealth:             math.Max(0, event.Player.Health()),
		VictimArmor:              math.Max(0, event.Player.Armor()),
		VictimNewHealth:          math.Max(0, event.Health),
		VictimNewArmor:           math.Max(0, event.Armor),
		IsVictimControllingBot:   event.Player.IsControllingBot(),
		AttackerSteamID64:        attackerSteamID,
		AttackerSide:             constants.Side(attackerSide),
		AttackerTeamName:         attackerTeamName,
		IsAttackerControllingBot: isAttackerControllingBot,
		VictimSteamID64:          event.Player.SteamID64,
		VictimSide:               constants.Side(event.Player.Team),
		VictimTeamName:           match.Team(constants.Side(event.Player.Team)).Name,
		WeaponName:               equipmentToWeaponName[event.Weapon.Type],
		WeaponType:               getEquipmentWeaponType(*event.Weapon),
		HitGroup:                 constants.HitGroup(event.HitGroup),
		WeaponUniqueID:           event.Weapon.UniqueID().String(),
	}
}
