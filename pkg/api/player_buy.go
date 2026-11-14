package api

import (
	"github.com/akiver/cs-demo-analyzer/pkg/api/constants"
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/common"
)

type PlayerBuy struct {
	Tick            int                  `json:"tick"`
	RoundNumber     int                  `json:"roundNumber"`
	PlayerSteamID64 uint64               `json:"playerSteamId"`
	PlayerSide      common.Team          `json:"playerSide"`
	PlayerName      string               `json:"playerName"`
	WeaponName      constants.WeaponName `json:"weaponName"`
	WeaponType      constants.WeaponType `json:"weaponType"`
	WeaponUniqueID  string               `json:"weaponUniqueId"`
	HasRefunded     bool                 `json:"hasRefunded"`
}

type playerFlashbangs struct {
	count int
	money int
}

func newPlayerBuy(analyzer *Analyzer, player *common.Player, weapon *common.Equipment) *PlayerBuy {
	return &PlayerBuy{
		Tick:            analyzer.currentTick(),
		RoundNumber:     analyzer.currentRound.Number,
		PlayerSteamID64: player.SteamID64,
		PlayerName:      player.Name,
		PlayerSide:      player.Team,
		WeaponName:      equipmentToWeaponName[weapon.Type],
		WeaponType:      getEquipmentWeaponType(*weapon),
		WeaponUniqueID:  weapon.UniqueID().String(),
	}
}

func (analyzer *Analyzer) flagLastBuyAsRefunded(player *common.Player, weapon *common.Equipment) {
	match := analyzer.match
	// Loop in reverse order to find the last purchase of the player and flag it as refunded.
	for i := len(match.PlayersBuy) - 1; i >= 0; i-- {
		buy := match.PlayersBuy[i]

		if buy.PlayerSteamID64 != player.SteamID64 || buy.WeaponName != equipmentToWeaponName[weapon.Type] || buy.HasRefunded {
			continue
		}

		if buy.RoundNumber != analyzer.currentRound.Number {
			break
		}

		buy.HasRefunded = true
		break
	}
}

// Buying a flashbang while holding one adds an ammo to the same weapon, there is no "item_pickup" event in this case.
// When one of these flashbangs is refunded, the weapon is not destroyed and there is no ItemRefund event either.
// To detect it, we compare players' flashbang count and money with the previous frame.
// The money decreases when a flashbang is bought but doesn't change when a flashbang dropped by a teammate is picked up.
func (analyzer *Analyzer) detectFlashbangBuys() {
	round := analyzer.currentRound
	if !analyzer.matchStarted() || round.secondsPassedSinceRoundStart(analyzer.buyTimeSeconds) {
		return
	}

	if round.playersFlashbangs == nil {
		round.playersFlashbangs = make(map[uint64]playerFlashbangs)
	}

	for _, player := range analyzer.parser.GameState().Participants().Playing() {
		current := playerFlashbangs{
			count: int(player.FlashbangCount()),
			money: player.Money(),
		}
		previous, exists := round.playersFlashbangs[player.SteamID64]
		round.playersFlashbangs[player.SteamID64] = current
		if !exists || current.count == previous.count || player.IsBot || !player.IsInBuyZone() {
			continue
		}

		flashbang := getPlayerFlashbang(player)
		if flashbang == nil {
			continue
		}

		if current.count > previous.count && current.money < previous.money {
			// The first flashbang is a new weapon, its purchase is detected with the "item_pickup" event.
			boughtCount := current.count - max(previous.count, 1)
			for range boughtCount {
				analyzer.match.PlayersBuy = append(analyzer.match.PlayersBuy, newPlayerBuy(analyzer, player, flashbang))
			}
		} else if current.count < previous.count && current.count > 0 && current.money > previous.money {
			for range previous.count - current.count {
				analyzer.flagLastBuyAsRefunded(player, flashbang)
			}
		}
	}
}

func getPlayerFlashbang(player *common.Player) *common.Equipment {
	for _, weapon := range player.Weapons() {
		if weapon.Type == common.EqFlash {
			return weapon
		}
	}

	return nil
}
