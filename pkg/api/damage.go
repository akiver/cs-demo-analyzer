package api

import (
	"fmt"

	"github.com/akiver/cs-demo-analyzer/internal/math"
	"github.com/akiver/cs-demo-analyzer/pkg/api/constants"
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/common"
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/events"
)

type Damage struct {
	Frame                    int                  `json:"frame"`
	Tick                     int                  `json:"tick"`
	RoundNumber              int                  `json:"roundNumber"`
	HealthDamage             int                  `json:"healthDamage"`
	ArmorDamage              int                  `json:"armorDamage"`
	AttackerSteamID64        uint64               `json:"attackerSteamId"`
	AttackerSide             common.Team          `json:"attackerSide"`
	AttackerTeamName         string               `json:"attackerTeamName"`
	IsAttackerControllingBot bool                 `json:"isAttackerControllingBot"`
	VictimHealth             int                  `json:"victimHealth"`
	VictimNewHealth          int                  `json:"victimNewHealth"`
	VictimArmor              int                  `json:"victimArmor"`
	VictimNewArmor           int                  `json:"victimNewArmor"`
	VictimSteamID64          uint64               `json:"victimSteamId"`
	VictimSide               common.Team          `json:"victimSide"`
	VictimTeamName           string               `json:"victimTeamName"`
	IsVictimControllingBot   bool                 `json:"isVictimControllingBot"`
	HitGroup                 events.HitGroup      `json:"hitgroup"`
	WeaponName               constants.WeaponName `json:"weaponName"`
	WeaponType               constants.WeaponType `json:"weaponType"`
	WeaponUniqueID           string               `json:"weaponUniqueId"`
}

func (damage *Damage) IsGrenadeWeapon() bool {
	return damage.WeaponType == constants.WeaponTypeGrenade
}

func (damage *Damage) isValidPlayerDamageEvent(player *Player) bool {
	if damage.AttackerSteamID64 != player.SteamID64 {
		return false
	}
	if damage.AttackerSteamID64 == damage.VictimSteamID64 {
		return false
	}
	if damage.VictimSteamID64 == 0 {
		return false
	}
	if damage.IsAttackerControllingBot {
		return false
	}
	if damage.AttackerSide == damage.VictimSide {
		return false
	}

	return true
}

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
		attackerTeamName = match.Team(event.Attacker.Team).Name
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
		AttackerSide:             attackerSide,
		AttackerTeamName:         attackerTeamName,
		IsAttackerControllingBot: isAttackerControllingBot,
		VictimSteamID64:          event.Player.SteamID64,
		VictimSide:               event.Player.Team,
		VictimTeamName:           match.Team(event.Player.Team).Name,
		WeaponName:               equipmentToWeaponName[event.Weapon.Type],
		WeaponType:               getEquipmentWeaponType(*event.Weapon),
		HitGroup:                 event.HitGroup,
		WeaponUniqueID:           event.Weapon.UniqueID2().String(),
	}
}
