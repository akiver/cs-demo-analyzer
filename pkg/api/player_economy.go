package api

import (
	"github.com/akiver/cs-demo-analyzer/pkg/api/constants"
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/common"
)

type PlayerEconomy struct {
	RoundNumber    int                   `json:"roundNumber"`
	Name           string                `json:"name"`
	SteamID64      uint64                `json:"steamId"`
	StartMoney     int                   `json:"startMoney"`
	MoneySpent     int                   `json:"moneySpent"`
	EquipmentValue int                   `json:"equipmentValue"`
	Type           constants.EconomyType `json:"type"`
	PlayerSide     common.Team           `json:"playerSide"`
}

func newPlayerEconomy(analyzer *Analyzer, player *common.Player) *PlayerEconomy {
	startMoney := player.Money()
	// eBot demos may start just after the end of the 1st round freeze time.
	// As a result to get the correct start money, we need to sum the current money + the money spent during the round.
	if analyzer.currentRound.Number == 1 {
		startMoney += player.MoneySpentThisRound()
	}

	economy := &PlayerEconomy{
		RoundNumber:    analyzer.currentRound.Number,
		SteamID64:      player.SteamID64,
		Name:           player.Name,
		StartMoney:     startMoney,
		EquipmentValue: player.EquipmentValueCurrent(),
		MoneySpent:     player.MoneySpentThisRound(),
		Type:           computePlayerEconomyType(analyzer, player),
		PlayerSide:     player.Team,
	}

	return economy
}

func (economy *PlayerEconomy) updateValues(analyzer *Analyzer, player *common.Player) {
	economy.PlayerSide = player.Team
	economy.EquipmentValue = player.EquipmentValueCurrent()
	economy.MoneySpent = player.MoneySpentThisRound()
	economy.Type = computePlayerEconomyType(analyzer, player)
}
