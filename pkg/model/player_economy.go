package model

import (
	"github.com/akiver/cs-demo-analyzer/pkg/api/constants"
)

type PlayerEconomy struct {
	RoundNumber    int                   `json:"roundNumber"`
	Name           string                `json:"name"`
	SteamID64      uint64                `json:"steamId"`
	StartMoney     int                   `json:"startMoney"`
	MoneySpent     int                   `json:"moneySpent"`
	EquipmentValue int                   `json:"equipmentValue"`
	Type           constants.EconomyType `json:"type"`
	PlayerSide     constants.Side        `json:"playerSide"`
}
