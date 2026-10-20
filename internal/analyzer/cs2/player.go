package cs2

import (
	"github.com/akiver/cs-demo-analyzer/internal/modelaccess"
	"math"

	"github.com/akiver/cs-demo-analyzer/internal/strings"
	"github.com/akiver/cs-demo-analyzer/pkg/api/constants"
	common "github.com/markus-wa/demoinfocs-golang/v6/pkg/demoinfocs/common"
	st "github.com/markus-wa/demoinfocs-golang/v6/pkg/demoinfocs/sendtables"
)

type weaponInspectionData struct {
	tick          int
	cancelledTick int
}

// weaponInspection returns the player's last weapon inspection, tick -1 if none.
func (analyzer *Analyzer) weaponInspection(player *Player) *weaponInspectionData {
	inspection, ok := analyzer.weaponInspections[player]
	if !ok {
		inspection = &weaponInspectionData{
			tick:          -1,
			cancelledTick: -1,
		}
		analyzer.weaponInspections[player] = inspection
	}

	return inspection
}

func (analyzer *Analyzer) isInspectingWeapon(player *Player) bool {
	inspection := analyzer.weaponInspection(player)

	// weapon inspection animation lasts approx. 5 seconds
	return inspection.tick > -1 && !analyzer.secondsHasPassedSinceTick(5, inspection.tick)
}

func (analyzer *Analyzer) startWeaponInspection(player *Player) {
	player.InspectWeaponCount++
	*analyzer.weaponInspection(player) = weaponInspectionData{
		tick:          analyzer.currentTick(),
		cancelledTick: -1,
	}
}

func (analyzer *Analyzer) stopWeaponInspection(player *Player) {
	inspection := analyzer.weaponInspection(player)
	// update the cancellation tick only if we didn't already detected a cancellation otherwise we may have multiple
	// cancellations for a single inspection which would mess up the time-based logic in isInspectingWeapon.
	if inspection.cancelledTick == -1 {
		inspection.cancelledTick = analyzer.currentTick()
	}
}

func newPlayer(analyzer *Analyzer, currentTeam common.Team, player common.Player) *Player {
	var team *Team
	if *analyzer.match.TeamA.CurrentSide == constants.Side(currentTeam) {
		team = analyzer.match.TeamA
	} else {
		team = analyzer.match.TeamB
	}

	color, _ := player.ColorOrErr()
	rank := player.Rank()
	userID := 0
	if player.UserID <= math.MaxUint16 {
		userID = player.UserID & 0xff
	}

	newPlayer := &Player{
		SteamID64:          player.SteamID64,
		UserID:             userID,
		Name:               strings.ReplaceUTF8ByteSequences(player.Name),
		Team:               team,
		CrosshairShareCode: player.CrosshairCode(),
		Color:              constants.Color(color),
		RankType:           player.RankType(),
		Rank:               rank,
		OldRank:            rank,
		WinCount:           player.CompetitiveWins(),
	}
	modelaccess.SetPlayerMatch(newPlayer, analyzer.match)

	var activeWeaponProp *st.Property
	if analyzer.isSource2 {
		if pawnEntity := player.PlayerPawnEntity(); pawnEntity != nil {
			activeWeaponProp = pawnEntity.Property("m_pWeaponServices.m_hActiveWeapon")
		}
	} else {
		activeWeaponProp = player.Entity.Property("m_hActiveWeapon")
	}
	if activeWeaponProp != nil {
		activeWeaponProp.OnUpdate(func(pv st.PropertyValue) {
			// switching weapon stops the weapon inspection animation
			analyzer.stopWeaponInspection(newPlayer)
		})
	}

	return newPlayer
}
