package api

import (
	"encoding/json"

	"github.com/akiver/cs-demo-analyzer/pkg/api/constants"
	common "github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/common"
)

type Player struct {
	match              *Match
	SteamID64          uint64       `json:"steamId"`
	Index              int          `json:"-"`
	Name               string       `json:"name"`
	Score              int          `json:"score"`
	Team               *Team        `json:"team"`
	MvpCount           int          `json:"mvpCount"`
	RankType           int          `json:"rankType"`
	Rank               int          `json:"rank"`
	OldRank            int          `json:"oldRank"`
	WinCount           int          `json:"winCount"`
	CrosshairShareCode string       `json:"crosshairShareCode"`
	Color              common.Color `json:"color"`
	InspectWeaponCount int          `json:"inspectWeaponCount"`
}

type PlayerAlias Player

type PlayerJSON struct {
	*PlayerAlias
	KillCount             int     `json:"killCount"`
	DeathCount            int     `json:"deathCount"`
	AssistCount           int     `json:"assistCount"`
	KillDeathRatio        float32 `json:"killDeathRatio"`
	KAST                  float32 `json:"kast"`
	BombDefusedCount      int     `json:"bombDefusedCount"`
	BombPlantedCount      int     `json:"bombPlantedCount"`
	HealthDamage          int     `json:"healthDamage"`
	ArmorDamage           int     `json:"armorDamage"`
	UtilityDamage         int     `json:"utilityDamage"`
	HeadshotCount         int     `json:"headshotCount"`
	HeadshotPercent       int     `json:"headshotPercent"`
	OneVsOneCount         int     `json:"oneVsOneCount"`
	OneVsOneWonCount      int     `json:"oneVsOneWonCount"`
	OneVsOneLostCount     int     `json:"oneVsOneLostCount"`
	OneVsTwoCount         int     `json:"oneVsTwoCount"`
	OneVsTwoWonCount      int     `json:"oneVsTwoWonCount"`
	OneVsTwoLostCount     int     `json:"oneVsTwoLostCount"`
	OneVsThreeCount       int     `json:"oneVsThreeCount"`
	OneVsThreeWonCount    int     `json:"oneVsThreeWonCount"`
	OneVsThreeLostCount   int     `json:"oneVsThreeLostCount"`
	OneVsFourCount        int     `json:"oneVsFourCount"`
	OneVsFourWonCount     int     `json:"oneVsFourWonCount"`
	OneVsFourLostCount    int     `json:"oneVsFourLostCount"`
	OneVsFiveCount        int     `json:"oneVsFiveCount"`
	OneVsFiveWonCount     int     `json:"oneVsFiveWonCount"`
	OneVsFiveLostCount    int     `json:"oneVsFiveLostCount"`
	HostageRescuedCount   int     `json:"hostageRescuedCount"`
	AverageKillPerRound   float32 `json:"averageKillPerRound"`
	AverageDeathPerRound  float32 `json:"averageDeathPerRound"`
	AverageDamagePerRound float32 `json:"averageDamagePerRound"`
	UtilityDamagePerRound float32 `json:"utilityDamagePerRound"`
	FirstKillCount        int     `json:"firstKillCount"`
	FirstDeathCount       int     `json:"firstDeathCount"`
	FirstTradeDeathCount  int     `json:"firstTradeDeathCount"`
	TradeDeathCount       int     `json:"tradeDeathCount"`
	TradeKillCount        int     `json:"tradeKillCount"`
	FirstTradeKillCount   int     `json:"firstTradeKillCount"`
	OneKillCount          int     `json:"oneKillCount"`
	TwoKillCount          int     `json:"twoKillCount"`
	ThreeKillCount        int     `json:"threeKillCount"`
	FourKillCount         int     `json:"fourKillCount"`
	FiveKillCount         int     `json:"fiveKillCount"`
	HltvRating            float32 `json:"hltvRating"`
	HltvRating2           float32 `json:"hltvRating2"`
}

func (player *Player) MarshalJSON() ([]byte, error) {
	return json.Marshal(PlayerJSON{
		PlayerAlias:           (*PlayerAlias)(player),
		KillCount:             player.KillCount(),
		DeathCount:            player.DeathCount(),
		AssistCount:           player.AssistCount(),
		KillDeathRatio:        player.KillDeathRatio(),
		KAST:                  player.KAST(),
		BombDefusedCount:      player.BombDefusedCount(),
		BombPlantedCount:      player.BombPlantedCount(),
		HealthDamage:          player.HealthDamage(),
		ArmorDamage:           player.ArmorDamage(),
		UtilityDamage:         player.UtilityDamage(),
		HeadshotCount:         player.HeadshotCount(),
		HeadshotPercent:       player.HeadshotPercent(),
		OneVsOneCount:         player.OneVsOneCount(),
		OneVsOneWonCount:      player.OneVsOneWonCount(),
		OneVsOneLostCount:     player.OneVsOneLostCount(),
		OneVsTwoCount:         player.OneVsTwoCount(),
		OneVsTwoWonCount:      player.OneVsTwoWonCount(),
		OneVsTwoLostCount:     player.OneVsTwoLostCount(),
		OneVsThreeCount:       player.OneVsThreeCount(),
		OneVsThreeWonCount:    player.OneVsThreeWonCount(),
		OneVsThreeLostCount:   player.OneVsThreeLostCount(),
		OneVsFourCount:        player.OneVsFourCount(),
		OneVsFourWonCount:     player.OneVsFourWonCount(),
		OneVsFourLostCount:    player.OneVsFourLostCount(),
		OneVsFiveCount:        player.OneVsFiveCount(),
		OneVsFiveWonCount:     player.OneVsFiveWonCount(),
		OneVsFiveLostCount:    player.OneVsFiveLostCount(),
		HostageRescuedCount:   player.HostageRescuedCount(),
		AverageKillPerRound:   player.AverageKillPerRound(),
		AverageDeathPerRound:  player.AverageDeathPerRound(),
		AverageDamagePerRound: player.AverageDamagePerRound(),
		UtilityDamagePerRound: player.UtilityDamagePerRound(),
		FirstKillCount:        player.FirstKillCount(),
		FirstDeathCount:       player.FirstDeathCount(),
		FirstTradeDeathCount:  player.FirstTradeDeathCount(),
		TradeDeathCount:       player.TradeDeathCount(),
		TradeKillCount:        player.TradeKillCount(),
		FirstTradeKillCount:   player.FirstTradeKillCount(),
		OneKillCount:          player.OneKillCount(),
		TwoKillCount:          player.TwoKillCount(),
		ThreeKillCount:        player.ThreeKillCount(),
		FourKillCount:         player.FourKillCount(),
		FiveKillCount:         player.FiveKillCount(),
		HltvRating2:           player.HltvRating2(),
		HltvRating:            player.HltvRating(),
	})
}

func (player *Player) TeamName() string {
	return player.Team.Name
}

func (player *Player) String() string {
	return player.Name
}

// This returns the percentage of rounds in which the player either had a kill, assist, survived or was traded.
func (player *Player) KAST() float32 {
	kastPerRound := make(map[int]bool)
	for _, round := range player.match.Rounds {
		kastPerRound[round.Number] = false
		playerSurvived := true

		for _, kill := range player.match.Kills {
			if round.Number != kill.RoundNumber {
				continue
			}

			isTeamKill := kill.KillerSide == kill.VictimSide
			if isTeamKill {
				continue
			}

			if kill.AssisterSteamID64 == player.SteamID64 {
				kastPerRound[round.Number] = true
				continue
			}

			if kill.KillerSteamID64 == player.SteamID64 && kill.VictimSteamID64 != player.SteamID64 {
				kastPerRound[round.Number] = true
				continue
			}

			if kill.VictimSteamID64 == player.SteamID64 {
				playerSurvived = false
				if kill.IsTradeDeath {
					kastPerRound[round.Number] = true
				}
			}
		}
		if playerSurvived {
			kastPerRound[round.Number] = true
		}
	}

	kastEventCount := 0
	for _, hasKASTEvent := range kastPerRound {
		if hasKASTEvent {
			kastEventCount++
		}
	}

	if len(kastPerRound) > 0 {
		return float32(kastEventCount) / float32(len(kastPerRound)) * 100
	}

	return 0
}

func (player *Player) BombPlantedCount() int {
	var bombPlantedCount int
	for _, bombPlanted := range player.match.BombsPlanted {
		if bombPlanted.PlanterSteamID64 == player.SteamID64 && !bombPlanted.IsPlayerControllingBot {
			bombPlantedCount++
		}
	}

	return bombPlantedCount
}

func (player *Player) BombDefusedCount() int {
	var bombDefusedCount int
	for _, bombDefused := range player.match.BombsDefused {
		if bombDefused.DefuserSteamID64 == player.SteamID64 && !bombDefused.IsPlayerControllingBot {
			bombDefusedCount++
		}
	}

	return bombDefusedCount
}

func (player *Player) Clutches() []*Clutch {
	clutches := []*Clutch{}
	for _, clutch := range player.match.Clutches {
		if clutch.ClutcherSteamID64 == player.SteamID64 {
			clutches = append(clutches, clutch)
		}
	}

	return clutches
}

func (player *Player) OneVsOneWonCount() int {
	return player.oneVsXClutchesWonCount(1)
}

func (player *Player) OneVsOneLostCount() int {
	return player.oneVsXClutchesLostCount(1)
}

func (player *Player) OneVsOneCount() int {
	clutches := player.oneVsXClutches(1)
	return len(clutches)
}

func (player *Player) OneVsTwoWonCount() int {
	return player.oneVsXClutchesWonCount(2)
}

func (player *Player) OneVsTwoLostCount() int {
	return player.oneVsXClutchesLostCount(2)
}

func (player *Player) OneVsTwoCount() int {
	clutches := player.oneVsXClutches(2)
	return len(clutches)
}

func (player *Player) OneVsThreeWonCount() int {
	return player.oneVsXClutchesWonCount(3)
}

func (player *Player) OneVsThreeLostCount() int {
	return player.oneVsXClutchesLostCount(3)
}

func (player *Player) OneVsThreeCount() int {
	clutches := player.oneVsXClutches(3)
	return len(clutches)
}

func (player *Player) OneVsFourWonCount() int {
	return player.oneVsXClutchesWonCount(4)
}

func (player *Player) OneVsFourLostCount() int {
	return player.oneVsXClutchesLostCount(4)
}

func (player *Player) OneVsFourCount() int {
	clutches := player.oneVsXClutches(4)
	return len(clutches)
}

func (player *Player) OneVsFiveWonCount() int {
	return player.oneVsXClutchesWonCount(5)
}

func (player *Player) OneVsFiveLostCount() int {
	return player.oneVsXClutchesLostCount(5)
}

func (player *Player) OneVsFiveCount() int {
	clutches := player.oneVsXClutches(5)
	return len(clutches)
}

func (player *Player) HostageRescuedCount() int {
	var hostageRescuedCount int
	for _, hostageRescued := range player.match.HostageRescued {
		if hostageRescued.PlayerSteamID64 == player.SteamID64 && !hostageRescued.IsPlayerControllingBot {
			hostageRescuedCount++
		}
	}

	return hostageRescuedCount
}

func (player *Player) kills() []*Kill {
	var kills []*Kill = make([]*Kill, 0)
	for _, kill := range player.match.Kills {
		if kill.KillerSteamID64 == player.SteamID64 && !kill.IsKillerControllingBot {
			kills = append(kills, kill)
		}
	}

	return kills
}

func (player *Player) Deaths() []*Kill {
	var deaths []*Kill = make([]*Kill, 0)
	for _, death := range player.match.Kills {
		if death.VictimSteamID64 == player.SteamID64 && !death.IsVictimControllingBot {
			deaths = append(deaths, death)
		}
	}

	return deaths
}

func (player *Player) KillCount() int {
	var killCount int
	for _, kill := range player.match.Kills {
		if kill.KillerSteamID64 == player.SteamID64 {
			if kill.IsKillerControllingBot {
				continue
			}

			if kill.IsSuicide() {
				// CSGO decreases player kill count on disconnection caused by network issue or by a vote kick, we don't
				isClientDisconnection := kill.WeaponName == constants.WeaponWorld
				if !isClientDisconnection {
					killCount--
				}
				continue
			}

			if kill.IsTeamKill() {
				killCount--
				continue
			}

			killCount++
		} else if kill.VictimSteamID64 == player.SteamID64 {
			if kill.IsVictimControllingBot {
				continue
			}

			isSuicide := kill.KillerSteamID64 == 0 && kill.WeaponName == constants.WeaponWorld
			if isSuicide {
				killCount--
			}
		}
	}

	return killCount
}

func (player *Player) DeathCount() int {
	var deathCount int
	for _, kill := range player.Deaths() {
		if kill.IsSuicide() {
			isClientDisconnection := kill.WeaponName == constants.WeaponWorld
			if isClientDisconnection {
				continue
			}
		}

		deathCount++
	}

	return deathCount
}

func (player *Player) AssistCount() int {
	var assistCount int
	for _, kill := range player.match.Kills {
		if kill.AssisterSteamID64 == player.SteamID64 && !kill.IsAssisterControllingBot && kill.AssisterSide != kill.VictimSide {
			assistCount++
		}
	}

	return assistCount
}

func (player *Player) HeadshotCount() int {
	var headshotCount int
	for _, kill := range player.kills() {
		if !kill.IsHeadshot || kill.IsSuicide() || kill.IsTeamKill() {
			continue
		}

		headshotCount++
	}

	return headshotCount
}

func (player *Player) FirstKillCount() int {
	var firstKillCount int
	for _, round := range player.match.Rounds {
		var killsInRound []*Kill = make([]*Kill, 0)
		for _, kill := range player.match.Kills {
			if kill.RoundNumber != round.Number {
				continue
			}
			killsInRound = append(killsInRound, kill)
		}

		for _, kill := range killsInRound {
			if kill.IsKillerControllingBot {
				continue
			}

			isSuicide := kill.KillerSteamID64 == kill.VictimSteamID64
			if isSuicide {
				continue
			}

			isTeamKill := kill.KillerSide == kill.VictimSide
			if isTeamKill {
				continue
			}

			if kill.KillerSteamID64 == player.SteamID64 {
				firstKillCount++
			}
			break
		}
	}

	return firstKillCount
}

func (player *Player) FirstDeathCount() int {
	var firstDeathCount int
	for _, round := range player.match.Rounds {
		var killsInRound []*Kill = make([]*Kill, 0)
		for _, kill := range player.match.Kills {
			if kill.RoundNumber != round.Number {
				continue
			}
			killsInRound = append(killsInRound, kill)
		}

		for _, kill := range killsInRound {
			if kill.IsKillerControllingBot || kill.IsSuicide() || kill.IsTeamKill() {
				continue
			}

			if kill.VictimSteamID64 == player.SteamID64 {
				firstDeathCount++
			}
			break
		}
	}

	return firstDeathCount
}

func (player *Player) FirstTradeDeathCount() int {
	var firstTradeDeathCount int
	for _, kills := range player.match.KillsByRound() {
		for _, kill := range kills {
			if kill.IsVictimControllingBot || kill.IsSuicide() || kill.IsTeamKill() {
				continue
			}

			if kill.VictimSteamID64 == player.SteamID64 && kill.IsTradeDeath {
				firstTradeDeathCount++
			}

			break
		}
	}

	return firstTradeDeathCount
}

func (player *Player) FirstTradeKillCount() int {
	var firstTradeKillCount int
	for _, kills := range player.match.KillsByRound() {
		for _, kill := range kills {
			if kill.IsKillerControllingBot || kill.IsSuicide() || kill.IsTeamKill() {
				continue
			}

			if kill.KillerSteamID64 == player.SteamID64 && kill.IsTradeKill {
				firstTradeKillCount++
			}

			break
		}
	}

	return firstTradeKillCount
}

func (player *Player) TradeDeathCount() int {
	var tradeDeathCount int
	for _, kill := range player.Deaths() {
		if !kill.IsTradeDeath || kill.IsSuicide() || kill.IsTeamKill() {
			continue
		}

		tradeDeathCount++
	}

	return tradeDeathCount
}

func (player *Player) TradeKillCount() int {
	var tradeKillCount int
	for _, kill := range player.kills() {
		if !kill.IsTradeKill || kill.IsSuicide() || kill.IsTeamKill() {
			continue
		}

		tradeKillCount++
	}

	return tradeKillCount
}

func (player *Player) HealthDamage() int {
	var healthDamage int
	for _, damage := range player.match.Damages {
		if damage.isValidPlayerDamageEvent(player) {
			healthDamage += damage.HealthDamage
		}
	}

	return healthDamage
}

func (player *Player) ArmorDamage() int {
	var armorDamage int
	for _, damage := range player.match.Damages {
		if damage.isValidPlayerDamageEvent(player) {
			armorDamage += damage.ArmorDamage
		}
	}

	return armorDamage
}

func (player *Player) UtilityDamage() int {
	var utilityDamage int
	for _, damage := range player.match.Damages {
		if damage.isValidPlayerDamageEvent(player) && damage.IsGrenadeWeapon() {
			utilityDamage += damage.HealthDamage
		}
	}

	return utilityDamage
}

func (player *Player) HeadshotPercent() int {
	killCount := player.KillCount()
	if killCount > 0 {
		return 100 * player.HeadshotCount() / killCount
	}

	return 0
}

func (player *Player) KillDeathRatio() float32 {
	killCount := player.KillCount()
	if killCount <= 0 {
		return 0
	}

	deathCount := player.DeathCount()
	if deathCount > 0 {
		return float32(killCount) / float32(deathCount)
	}

	return float32(killCount)
}

func (player *Player) AverageKillPerRound() float32 {
	killCount := player.KillCount()
	roundCount := player.roundCount()
	if killCount <= 0 || roundCount <= 0 {
		return 0
	}

	return float32(killCount) / float32(roundCount)
}

func (player *Player) AverageAssistPerRound() float32 {
	assistCount := player.AssistCount()
	roundCount := player.roundCount()
	if assistCount <= 0 || roundCount <= 0 {
		return 0
	}

	return float32(assistCount) / float32(roundCount)
}

func (player *Player) AverageDeathPerRound() float32 {
	deathCount := player.DeathCount()
	roundCount := player.roundCount()
	if deathCount <= 0 || roundCount <= 0 {
		return 0
	}

	return float32(deathCount) / float32(roundCount)
}

func (player *Player) AverageDamagePerRound() float32 {
	roundCount := player.roundCount()
	if roundCount > 0 {
		return float32(player.HealthDamage()) / float32(roundCount)
	}

	return 0
}

func (player *Player) UtilityDamagePerRound() float32 {
	roundCount := player.roundCount()
	if roundCount > 0 {
		return float32(player.UtilityDamage()) / float32(roundCount)
	}

	return 0
}

func (player *Player) OneKillCount() int {
	return player.getXKillCount(1)
}

func (player *Player) TwoKillCount() int {
	return player.getXKillCount(2)
}

func (player *Player) ThreeKillCount() int {
	return player.getXKillCount(3)
}

func (player *Player) FourKillCount() int {
	return player.getXKillCount(4)
}

func (player *Player) FiveKillCount() int {
	return player.getXKillCount(5)
}

// This returns the "impact" as described in the following blog post.
// https://flashed.gg/posts/reverse-engineering-hltv-rating/
// 2.13*KPR + 0.42*Assist per Round -0.41 ≈ impact
func (player *Player) impact() float32 {
	return 2.13*player.AverageKillPerRound() + 0.42*player.AverageAssistPerRound() + -0.41
}

// This returns the player's HLTV rating 2.0.
// https://flashed.gg/posts/reverse-engineering-hltv-rating/
// 0.0073*KAST + 0.3591*KPR + -0.5329*DPR + 0.2372*Impact + 0.0032*ADR + 0.1587 ≈ Rating 2.0
func (player *Player) HltvRating2() float32 {
	rating := 0.0073*player.KAST() + 0.3591*player.AverageKillPerRound() + -0.5329*player.AverageDeathPerRound() + 0.2372*player.impact() + 0.0032*float32(player.AverageDamagePerRound()) + 0.1587

	if rating < 0 {
		return 0
	}

	return rating
}

// This returns the player's HLTV rating 1.0.
// Formula: https://web.archive.org/web/20170427062206/http://www.hltv.org/?pageid=242&eventid=0
func (player *Player) HltvRating() float32 {
	roundCount := float32(player.roundCount())
	if roundCount == 0 {
		return 0
	}

	killRating := player.AverageKillPerRound() / 0.679
	survivalRating := (roundCount - float32(player.DeathCount())) / roundCount / 0.317
	roundsWithMultipleKillsRating := (float32(player.OneKillCount()) + 4*float32(player.TwoKillCount()) + 9*float32(player.ThreeKillCount()) + 16*float32(player.FourKillCount()) + 25*float32(player.FiveKillCount())) / roundCount / 1.277
	rating := (killRating + 0.7*survivalRating + roundsWithMultipleKillsRating) / 2.7

	return rating
}

func (player *Player) roundCount() int {
	return len(player.match.Rounds)
}

func (player *Player) oneVsXClutches(opponentCount int) []*Clutch {
	clutches := player.Clutches()
	for _, clutch := range player.match.Clutches {
		if clutch.OpponentCount == opponentCount {
			clutches = append(clutches, clutch)
		}
	}

	return clutches
}

func (player *Player) oneVsXClutchesWonCount(opponentCount int) int {
	clutches := player.oneVsXClutches(opponentCount)
	var count int
	for _, clutch := range clutches {
		if clutch.HasWon {
			count++
		}
	}

	return count
}

func (player *Player) oneVsXClutchesLostCount(opponentCount int) int {
	clutches := player.oneVsXClutches(opponentCount)
	var count int
	for _, clutch := range clutches {
		if !clutch.HasWon {
			count++
		}
	}

	return count
}

func (player *Player) getXKillCount(count int) int {
	var xKillCount int
	for _, kills := range player.match.KillsByRound() {
		playerKillInRoundCount := 0
		for _, kill := range kills {
			if kill.KillerSteamID64 != player.SteamID64 || kill.IsKillerControllingBot || kill.IsSuicide() || kill.IsTeamKill() {
				continue
			}

			playerKillInRoundCount++
		}

		if playerKillInRoundCount == count {
			xKillCount++
		}
	}

	return xKillCount
}

func (player *Player) reset() {
	player.Score = 0
	player.MvpCount = 0
}
