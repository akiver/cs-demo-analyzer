package api

import (
	"github.com/akiver/cs-demo-analyzer/pkg/api/constants"
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/common"
)

func getValidPlayerCount(players []*common.Player) int {
	playerCount := 0
	for _, player := range players {
		if !player.IsBot && !player.IsUnknown {
			playerCount++
		}
	}

	return playerCount
}

func computePlayerEconomyType(analyzer *Analyzer, player *common.Player) constants.EconomyType {
	if analyzer.isFirstRoundOfHalf && analyzer.match.OvertimeCount == 0 {
		return constants.EconomyTypePistol
	}

	equipmentValue := player.EquipmentValueCurrent()
	if equipmentValue <= 1000 {
		return constants.EconomyTypeEco
	}

	playerSide := player.Team
	minFullEquipmentValue := 4500
	if playerSide == common.TeamTerrorists {
		minFullEquipmentValue = 4000
	}

	if equipmentValue >= minFullEquipmentValue {
		return constants.EconomyTypeFull
	}

	if len(analyzer.match.Rounds) > 0 {
		previousRound := analyzer.match.Rounds[len(analyzer.match.Rounds)-1]
		if previousRound.WinnerSide != playerSide && player.Money() <= 400 {
			return constants.EconomyTypeForceBuy
		}
	}

	return constants.EconomyTypeSemi
}

func computeTeamEconomyType(analyzer *Analyzer, team *common.TeamState) constants.EconomyType {
	if analyzer.isFirstRoundOfHalf && analyzer.match.OvertimeCount == 0 {
		return constants.EconomyTypePistol
	}

	teamSide := team.Team()
	playerCount := getValidPlayerCount(team.Members())
	equipmentValue := team.CurrentEquipmentValue()
	if equipmentValue <= 1000*playerCount {
		return constants.EconomyTypeEco
	}

	minFullEquipmentValue := 4500 * playerCount
	if teamSide == common.TeamTerrorists {
		minFullEquipmentValue = 4000 * playerCount
	}

	if equipmentValue >= minFullEquipmentValue {
		return constants.EconomyTypeFull
	}

	if len(analyzer.match.Rounds) > 0 {
		previousRound := analyzer.match.Rounds[len(analyzer.match.Rounds)-1]
		money := 0
		for _, player := range team.Members() {
			money += player.Money()
		}

		if previousRound.WinnerSide != teamSide && money < 400*playerCount {
			return constants.EconomyTypeForceBuy
		}
	}

	return constants.EconomyTypeSemi
}
