package api

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/akiver/cs-demo-analyzer/internal/converters"
	"github.com/akiver/cs-demo-analyzer/internal/csv"
	"github.com/akiver/cs-demo-analyzer/internal/slice"
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/common"
)

func exportMatchToCSV(match *Match, outputPath string) error {
	if stat, err := os.Stat(outputPath); err != nil || !stat.IsDir() {
		return errors.New("incorrect output provided, make sure it's a folder that exists and you have write access")
	}

	outputPath = outputPath + string(os.PathSeparator) + match.DemoFileName

	var writeMatch = func() {
		header := []string{
			"checksum",
			"game",
			"demo path",
			"demo name",
			"date",
			"source",
			"type",
			"map",
			"server name",
			"client name",
			"tick count",
			"tickrate",
			"framerate",
			"game type",
			"game mode",
			"game mode str",
			"is ranked",
			"duration",
			"network protocol",
			"build number",
			"kill count",
			"assist count",
			"death count",
			"shot count",
			"winner name",
			"winner side",
			"overtime count",
			"max rounds",
			"has vac live ban",
		}

		winnerName := ""
		winnerSide := common.TeamUnassigned
		if match.Winner != nil {
			winnerName = match.Winner.Name
			winnerSide = *match.Winner.CurrentSide
		}
		line := []string{
			match.Checksum,
			match.Game.String(),
			match.DemoFilePath,
			match.DemoFileName,
			match.Date.Format(time.RFC3339),
			match.Source.String(),
			match.Type.String(),
			match.MapName,
			match.ServerName,
			match.ClientName,
			converters.IntToString(match.TickCount),
			converters.Float64ToString(match.TickRate),
			converters.Float64ToString(match.FrameRate),
			converters.Float64ToString(match.Duration.Seconds()),
			converters.IntToString(match.NetworkProtocol),
			converters.IntToString(match.BuildNumber),
			match.GameType.String(),
			match.GameMode.String(),
			match.GameModeStr().String(),
			converters.BoolToString(match.IsRanked),
			converters.IntToString(match.KillCount()),
			converters.IntToString(match.AssistCount()),
			converters.IntToString(match.DeathCount()),
			converters.IntToString(match.ShotCount()),
			winnerName,
			converters.TeamToString(winnerSide),
			converters.IntToString(match.OvertimeCount),
			converters.IntToString(match.MaxRounds),
			converters.BoolToString(match.HasVacLiveBan),
		}

		csv.WriteLinesIntoCsvFile(outputPath+"_match.csv", [][]string{
			header,
			line,
		})
	}

	var buildTeamLine = func(team *Team) []string {
		line := []string{
			team.Name,
			team.Letter.String(),
			converters.IntToString(team.Score),
			converters.IntToString(team.ScoreFirstHalf),
			converters.IntToString(team.ScoreSecondHalf),
			converters.TeamToString(*team.CurrentSide),
			match.Checksum,
		}

		return line
	}

	var writeTeams = func() {
		header := []string{
			"name",
			"letter",
			"score",
			"score first half",
			"score second half",
			"current side",
			"match checksum",
		}

		csv.WriteLinesIntoCsvFile(outputPath+"_teams.csv", [][]string{
			header,
			buildTeamLine(match.TeamA),
			buildTeamLine(match.TeamB),
		})
	}

	var writePlayers = func() {
		header := []string{
			"name",
			"steamid",
			"score",
			"team name",
			"kills",
			"assists",
			"deaths",
			"headshots",
			"hs %",
			"k/d",
			"kast",
			"avg damages per round",
			"avg kills per round",
			"avg death per round",
			"utility_damage_per_round",
			"mvp",
			"rank type",
			"rank",
			"old rank",
			"win count",
			"bomb planted",
			"bomb defused",
			"hostage rescued",
			"health damage",
			"armor damage",
			"utility damage",
			"1v1",
			"1v2",
			"1v3",
			"1v4",
			"1v5",
			"1v1 won",
			"1v2 won",
			"1v3 won",
			"1v4 won",
			"1v5 won",
			"1v1 lost",
			"1v2 lost",
			"1v3 lost",
			"1v4 lost",
			"1v5 lost",
			"first kill",
			"first death",
			"trade kill",
			"trade death",
			"first trade kill",
			"first trade death",
			"1k",
			"2k",
			"3k",
			"4k",
			"5k",
			"htlv 2",
			"htlv",
			"crosshair share code",
			"color",
			"inspect weapon count",
			"match checksum",
		}

		lines := [][]string{header}
		for _, player := range match.Players() {
			line := []string{
				player.Name,
				converters.Uint64ToString(player.SteamID64),
				converters.IntToString(player.Score),
				player.TeamName(),
				converters.IntToString(player.KillCount()),
				converters.IntToString(player.AssistCount()),
				converters.IntToString(player.DeathCount()),
				converters.IntToString(player.HeadshotCount()),
				converters.IntToString(player.HeadshotPercent()),
				converters.Float32ToString(player.KillDeathRatio()),
				converters.Float32ToString(player.KAST()),
				converters.Float32ToString(player.AverageDamagePerRound()),
				converters.Float32ToString(player.AverageKillPerRound()),
				converters.Float32ToString(player.AverageDeathPerRound()),
				converters.Float32ToString(player.UtilityDamagePerRound()),
				converters.IntToString(player.MvpCount),
				converters.IntToString(player.RankType),
				converters.IntToString(player.Rank),
				converters.IntToString(player.OldRank),
				converters.IntToString(player.WinCount),
				converters.IntToString(player.BombPlantedCount()),
				converters.IntToString(player.BombDefusedCount()),
				converters.IntToString(player.HostageRescuedCount()),
				converters.IntToString(player.HealthDamage()),
				converters.IntToString(player.ArmorDamage()),
				converters.IntToString(player.UtilityDamage()),
				converters.IntToString(player.OneVsOneCount()),
				converters.IntToString(player.OneVsTwoCount()),
				converters.IntToString(player.OneVsThreeCount()),
				converters.IntToString(player.OneVsFourCount()),
				converters.IntToString(player.OneVsFiveCount()),
				converters.IntToString(player.OneVsOneWonCount()),
				converters.IntToString(player.OneVsTwoWonCount()),
				converters.IntToString(player.OneVsThreeWonCount()),
				converters.IntToString(player.OneVsFourWonCount()),
				converters.IntToString(player.OneVsFiveWonCount()),
				converters.IntToString(player.OneVsOneLostCount()),
				converters.IntToString(player.OneVsTwoLostCount()),
				converters.IntToString(player.OneVsThreeLostCount()),
				converters.IntToString(player.OneVsFourLostCount()),
				converters.IntToString(player.OneVsFiveLostCount()),
				converters.IntToString(player.FirstKillCount()),
				converters.IntToString(player.FirstDeathCount()),
				converters.IntToString(player.TradeKillCount()),
				converters.IntToString(player.TradeDeathCount()),
				converters.IntToString(player.FirstTradeKillCount()),
				converters.IntToString(player.FirstTradeDeathCount()),
				converters.IntToString(player.OneKillCount()),
				converters.IntToString(player.TwoKillCount()),
				converters.IntToString(player.ThreeKillCount()),
				converters.IntToString(player.FourKillCount()),
				converters.IntToString(player.FiveKillCount()),
				converters.Float32ToString(player.HltvRating2()),
				converters.Float32ToString(player.HltvRating()),
				player.CrosshairShareCode,
				converters.ColorToString(player.Color),
				converters.IntToString(player.InspectWeaponCount),
				match.Checksum,
			}
			lines = append(lines, line)
		}

		csv.WriteLinesIntoCsvFile(outputPath+"_players.csv", lines)
	}

	var writePlayerPositions = func() {
		header := []string{
			"frame",
			"tick",
			"is alive",
			"x",
			"y",
			"z",
			"yaw",
			"flash duration remaining",
			"side",
			"money",
			"health",
			"armor",
			"has helmet",
			"has bomb",
			"has defuse kit",
			"is ducking",
			"is airborne",
			"is scoping",
			"is defusing",
			"is planting",
			"is grabbing hostage",
			"active weapon name",
			"equipments",
			"grenades",
			"pistols",
			"smgs",
			"rifles",
			"heavy",
			"steamid",
			"name",
			"round",
			"match checksum",
		}

		lines := [][]string{header}
		for _, position := range match.PlayerPositions {
			line := []string{
				converters.IntToString(position.Frame),
				converters.IntToString(position.Tick),
				converters.BoolToString(position.IsAlive),
				converters.Float64ToString(position.X),
				converters.Float64ToString(position.Y),
				converters.Float64ToString(position.Z),
				converters.Float32ToString(position.Yaw),
				converters.Float64ToString(position.FlashDurationRemaining),
				converters.TeamToString(position.Side),
				converters.IntToString(position.Money),
				converters.IntToString(position.Health),
				converters.IntToString(position.Armor),
				converters.BoolToString(position.HasHelmet),
				converters.BoolToString(position.HasBomb),
				converters.BoolToString(position.HasDefuseKit),
				converters.BoolToString(position.IsDucking),
				converters.BoolToString(position.IsAirborne),
				converters.BoolToString(position.IsScoping),
				converters.BoolToString(position.IsDefusing),
				converters.BoolToString(position.IsPlanting),
				converters.BoolToString(position.IsGrabbingHostage),
				position.ActiveWeaponName.String(),
				strings.Join(slice.ToStrings(position.Equipments), ","),
				strings.Join(slice.ToStrings(position.Grenades), ","),
				strings.Join(slice.ToStrings(position.Pistols), ","),
				strings.Join(slice.ToStrings(position.SMGs), ","),
				strings.Join(slice.ToStrings(position.Rifles), ","),
				strings.Join(slice.ToStrings(position.Heavy), ","),
				converters.Uint64ToString(position.SteamID64),
				position.Name,
				converters.IntToString(position.RoundNumber),
				match.Checksum,
			}
			lines = append(lines, line)
		}

		csv.WriteLinesIntoCsvFile(outputPath+"_positions.csv", lines)
	}

	var writeShots = func() {
		header := []string{
			"frame",
			"tick",
			"round",
			"weapon name",
			"weapon id",
			"projectile id",
			"x",
			"y",
			"z",
			"player name",
			"player steamid",
			"player team name",
			"player side",
			"is player controlling bot",
			"player velocity x",
			"player velocity y",
			"player velocity z",
			"yaw",
			"pitch",
			"recoil index",
			"aim punch angle x",
			"aim punch angle y",
			"view punch angle x",
			"view punch angle y",
			"match checksum",
		}

		lines := [][]string{header}
		for _, shot := range match.Shots {
			line := []string{
				converters.IntToString(shot.Frame),
				converters.IntToString(shot.Tick),
				converters.IntToString(shot.RoundNumber),
				shot.WeaponName.String(),
				shot.WeaponID,
				converters.Int64ToString(shot.ProjectileID),
				converters.Float64ToString(shot.X),
				converters.Float64ToString(shot.Y),
				converters.Float64ToString(shot.Z),
				shot.PlayerName,
				converters.Uint64ToString(shot.PlayerSteamID64),
				shot.PlayerTeamName,
				converters.TeamToString(shot.PlayerSide),
				converters.BoolToString(shot.IsPlayerControllingBot),
				converters.Float32ToString(shot.Yaw),
				converters.Float32ToString(shot.Pitch),
				converters.Float64ToString(shot.PlayerVelocityX),
				converters.Float64ToString(shot.PlayerVelocityY),
				converters.Float64ToString(shot.PlayerVelocityZ),
				converters.Float32ToString(shot.RecoilIndex),
				converters.Float64ToString(shot.AimPunchAngleX),
				converters.Float64ToString(shot.AimPunchAngleY),
				converters.Float64ToString(shot.ViewPunchAngleX),
				converters.Float64ToString(shot.ViewPunchAngleY),
				match.Checksum,
			}
			lines = append(lines, line)
		}

		csv.WriteLinesIntoCsvFile(outputPath+"_shots.csv", lines)
	}

	var writeRounds = func() {
		header := []string{
			"number",
			"start tick",
			"start frame",
			"freezetime end tick",
			"freezetime end frame",
			"end tick",
			"end frame",
			"end officially tick",
			"end officially frame",
			"team a name",
			"team b name",
			"score team a",
			"score team b",
			"team a side",
			"team b side",
			"team a start money",
			"team b start money",
			"team a equipment value",
			"team b equipment value",
			"team a money spent",
			"team b money spent",
			"team a economy type",
			"team b economy type",
			"duration",
			"end reason",
			"winner name",
			"winner side",
			"overtime number",
			"match checksum",
		}
		lines := [][]string{header}
		for _, round := range match.Rounds {
			line := []string{
				converters.IntToString(round.Number),
				converters.IntToString(round.StartTick),
				converters.IntToString(round.StartFrame),
				converters.IntToString(round.FreezeTimeEndTick),
				converters.IntToString(round.FreezeTimeEndFrame),
				converters.IntToString(round.EndTick),
				converters.IntToString(round.EndFrame),
				converters.IntToString(round.EndOfficiallyTick),
				converters.IntToString(round.EndOfficiallyFrame),
				round.TeamAName,
				round.TeamBName,
				converters.IntToString(round.TeamAScore),
				converters.IntToString(round.TeamBScore),
				converters.TeamToString(round.TeamASide),
				converters.TeamToString(round.TeamBSide),
				converters.IntToString(round.StartMoneyTeamA()),
				converters.IntToString(round.StartMoneyTeamB()),
				converters.IntToString(round.TeamAEquipmentValue),
				converters.IntToString(round.TeamBEquipmentValue),
				converters.IntToString(round.TeamAMoneySpent),
				converters.IntToString(round.TeamBMoneySpent),
				round.TeamAEconomyType.String(),
				round.TeamBEconomyType.String(),
				converters.Int64ToString(round.Duration),
				converters.RoundEndReasonToString(round.EndReason),
				round.WinnerName,
				converters.TeamToString(round.WinnerSide),
				converters.IntToString(round.OvertimeNumber),
				match.Checksum,
			}
			lines = append(lines, line)
		}

		csv.WriteLinesIntoCsvFile(outputPath+"_rounds.csv", lines)
	}

	var writeRoundEconomies = func() {
		header := []string{
			"steamid",
			"name",
			"player side",
			"start money",
			"money spent",
			"equipment value",
			"type",
			"round",
			"match checksum",
		}

		lines := [][]string{header}
		for _, economy := range match.PlayerEconomies {
			line := []string{
				converters.Uint64ToString(economy.SteamID64),
				economy.Name,
				converters.TeamToString(economy.PlayerSide),
				converters.IntToString(economy.StartMoney),
				converters.IntToString(economy.MoneySpent),
				converters.IntToString(economy.EquipmentValue),
				economy.Type.String(),
				converters.IntToString(economy.RoundNumber),
				match.Checksum,
			}
			lines = append(lines, line)
		}

		csv.WriteLinesIntoCsvFile(outputPath+"_players_economy.csv", lines)
	}

	var writeClutches = func() {
		header := []string{
			"frame",
			"tick",
			"round",
			"opponents",
			"side",
			"won",
			"steamid",
			"name",
			"survived",
			"kill_count",
			"match checksum",
		}

		lines := [][]string{header}
		for _, clutch := range match.Clutches {
			line := []string{
				converters.IntToString(clutch.Frame),
				converters.IntToString(clutch.Tick),
				converters.IntToString(clutch.RoundNumber),
				converters.IntToString(clutch.OpponentCount),
				converters.TeamToString(clutch.Side),
				converters.BoolToString(clutch.HasWon),
				converters.Uint64ToString(clutch.ClutcherSteamID64),
				clutch.ClutcherName,
				converters.BoolToString(clutch.ClutcherSurvived),
				converters.IntToString(clutch.ClutcherKillCount),
				match.Checksum,
			}
			lines = append(lines, line)
		}

		csv.WriteLinesIntoCsvFile(outputPath+"_clutches.csv", lines)
	}

	var writeChickenDeaths = func() {
		header := []string{
			"frame",
			"tick",
			"round",
			"killer steam id",
			"weapon name",
			"match checksum",
		}

		lines := [][]string{header}
		for _, chickenDeath := range match.ChickenDeaths {
			line := []string{
				converters.IntToString(chickenDeath.Frame),
				converters.IntToString(chickenDeath.Tick),
				converters.IntToString(chickenDeath.RoundNumber),
				converters.Uint64ToString(chickenDeath.KillerSteamID),
				chickenDeath.WeaponName.String(),
				match.Checksum,
			}
			lines = append(lines, line)
		}

		csv.WriteLinesIntoCsvFile(outputPath+"_chicken_deaths.csv", lines)
	}

	var writeChickenPositions = func() {
		header := []string{
			"frame",
			"tick",
			"round",
			"x",
			"y",
			"z",
			"match checksum",
		}

		lines := [][]string{header}
		for _, position := range match.ChickenPositions {
			line := []string{
				converters.IntToString(position.Frame),
				converters.IntToString(position.Tick),
				converters.IntToString(position.RoundNumber),
				converters.Float64ToString(position.X),
				converters.Float64ToString(position.Y),
				converters.Float64ToString(position.Z),
				match.Checksum,
			}
			lines = append(lines, line)
		}

		csv.WriteLinesIntoCsvFile(outputPath+"_chicken_positions.csv", lines)
	}

	var writeDamages = func() {
		header := []string{
			"frame",
			"tick",
			"round",
			"health",
			"armor",
			"victim health",
			"victim new health",
			"victim armor",
			"victim new armor",
			"attacker steamid",
			"attacker side",
			"attacker team name",
			"is attacker controlling bot",
			"victim steamid",
			"victim side",
			"victim team name",
			"is victim controlling bot",
			"weapon name",
			"weapon class",
			"hitgroup",
			"weapon unique id",
			"match checksum",
		}

		lines := [][]string{header}
		for _, damage := range match.Damages {
			line := []string{
				converters.IntToString(damage.Frame),
				converters.IntToString(damage.Tick),
				converters.IntToString(damage.RoundNumber),
				converters.IntToString(damage.HealthDamage),
				converters.IntToString(damage.ArmorDamage),
				converters.IntToString(damage.VictimHealth),
				converters.IntToString(damage.VictimNewHealth),
				converters.IntToString(damage.VictimArmor),
				converters.IntToString(damage.VictimNewArmor),
				converters.Uint64ToString(damage.AttackerSteamID64),
				converters.TeamToString(damage.AttackerSide),
				damage.AttackerTeamName,
				converters.BoolToString(damage.IsAttackerControllingBot),
				converters.Uint64ToString(damage.VictimSteamID64),
				converters.TeamToString(damage.VictimSide),
				damage.VictimTeamName,
				converters.BoolToString(damage.IsVictimControllingBot),
				damage.WeaponName.String(),
				string(damage.WeaponType),
				converters.HitgroupToString(damage.HitGroup),
				damage.WeaponUniqueID,
				match.Checksum,
			}
			lines = append(lines, line)
		}

		csv.WriteLinesIntoCsvFile(outputPath+"_damages.csv", lines)
	}

	var writeKills = func() {
		header := []string{
			"frame",
			"tick",
			"round",
			"killer name",
			"killer steamid",
			"killer side",
			"killer team name",
			"victim name",
			"victim steamid",
			"victim side",
			"victim team name",
			"assister name",
			"assister steamid",
			"assister side",
			"assister team name",
			"weapon name",
			"weapon_type",
			"headshot",
			"penetrated objects",
			"is flash assist",
			"killer controlling bot",
			"victim controlling bot",
			"assister controlling bot",
			"killer x",
			"killer y",
			"killer z",
			"is killer airborne",
			"is killer blinded",
			"victim x",
			"victim y",
			"victim z",
			"is victim airborne",
			"is victim blinded",
			"is victim inspecting weapon",
			"assister x",
			"assister y",
			"assister z",
			"is trade kill",
			"is trade death",
			"is through smoke",
			"is no scope",
			"distance",
			"match checksum",
		}

		lines := [][]string{header}
		for _, kill := range match.Kills {
			line := []string{
				converters.IntToString(kill.Frame),
				converters.IntToString(kill.Tick),
				converters.IntToString(kill.RoundNumber),
				kill.KillerName,
				converters.Uint64ToString(kill.KillerSteamID64),
				converters.TeamToString(kill.KillerSide),
				kill.KillerTeamName,
				kill.VictimName,
				converters.Uint64ToString(kill.VictimSteamID64),
				converters.TeamToString(kill.VictimSide),
				kill.VictimTeamName,
				kill.AssisterName,
				converters.Uint64ToString(kill.AssisterSteamID64),
				converters.TeamToString(kill.AssisterSide),
				kill.AssisterTeamName,
				kill.WeaponName.String(),
				string(kill.WeaponType),
				converters.BoolToString(kill.IsHeadshot),
				converters.IntToString(kill.PenetratedObjects),
				converters.BoolToString(kill.IsAssistedFlash),
				converters.BoolToString(kill.IsKillerControllingBot),
				converters.BoolToString(kill.IsVictimControllingBot),
				converters.BoolToString(kill.IsAssisterControllingBot),
				converters.Float64ToString(kill.KillerX),
				converters.Float64ToString(kill.KillerY),
				converters.Float64ToString(kill.KillerZ),
				converters.BoolToString(kill.IsKillerAirborne),
				converters.BoolToString(kill.IsKillerBlinded),
				converters.Float64ToString(kill.VictimX),
				converters.Float64ToString(kill.VictimY),
				converters.Float64ToString(kill.VictimZ),
				converters.BoolToString(kill.IsVictimAirborne),
				converters.BoolToString(kill.IsVictimBlinded),
				converters.BoolToString(kill.IsVictimInspectingWeapon),
				converters.Float64ToString(kill.AssisterX),
				converters.Float64ToString(kill.AssisterY),
				converters.Float64ToString(kill.AssisterZ),
				converters.BoolToString(kill.IsTradeKill),
				converters.BoolToString(kill.IsTradeDeath),
				converters.BoolToString(kill.IsThroughSmoke),
				converters.BoolToString(kill.IsNoScope),
				converters.Float32ToString(kill.Distance),
				match.Checksum,
			}
			lines = append(lines, line)
		}

		csv.WriteLinesIntoCsvFile(outputPath+"_kills.csv", lines)
	}

	var writeBombsPlanted = func() {
		header := []string{
			"frame",
			"tick",
			"round",
			"site",
			"planter steamid",
			"planter name",
			"is player controlling bot",
			"x",
			"y",
			"z",
			"match checksum",
		}

		lines := [][]string{header}
		for _, bombPlanted := range match.BombsPlanted {
			line := []string{
				converters.IntToString(bombPlanted.Frame),
				converters.IntToString(bombPlanted.Tick),
				converters.IntToString(bombPlanted.RoundNumber),
				bombPlanted.Site,
				converters.Uint64ToString(bombPlanted.PlanterSteamID64),
				bombPlanted.PlanterName,
				converters.BoolToString(bombPlanted.IsPlayerControllingBot),
				converters.Float64ToString(bombPlanted.X),
				converters.Float64ToString(bombPlanted.Y),
				converters.Float64ToString(bombPlanted.Z),
				match.Checksum,
			}
			lines = append(lines, line)
		}

		csv.WriteLinesIntoCsvFile(outputPath+"_bombs_planted.csv", lines)
	}

	var writeBombsDefuseStart = func() {
		header := []string{
			"frame",
			"tick",
			"round",
			"defuser steamid",
			"defuser name",
			"is defuser controlling bot",
			"x",
			"y",
			"z",
			"match checksum",
		}

		lines := [][]string{header}
		for _, bombDefuseStart := range match.BombsDefuseStart {
			line := []string{
				converters.IntToString(bombDefuseStart.Frame),
				converters.IntToString(bombDefuseStart.Tick),
				converters.IntToString(bombDefuseStart.RoundNumber),
				converters.Uint64ToString(bombDefuseStart.PlanterSteamID64),
				bombDefuseStart.PlanterName,
				converters.BoolToString(bombDefuseStart.IsPlayerControllingBot),
				converters.Float64ToString(bombDefuseStart.X),
				converters.Float64ToString(bombDefuseStart.Y),
				converters.Float64ToString(bombDefuseStart.Z),
				match.Checksum,
			}
			lines = append(lines, line)
		}

		csv.WriteLinesIntoCsvFile(outputPath+"_bombs_defuse_start.csv", lines)
	}

	var writeBombsDefused = func() {
		header := []string{
			"frame",
			"tick",
			"round",
			"site",
			"defuser steamid",
			"defuser name",
			"is player controlling bot",
			"x",
			"y",
			"z",
			"counter terrorist alive count",
			"terrorist alive count",
			"match checksum",
		}

		lines := [][]string{header}
		for _, bombDefused := range match.BombsDefused {
			line := []string{
				converters.IntToString(bombDefused.Frame),
				converters.IntToString(bombDefused.Tick),
				converters.IntToString(bombDefused.RoundNumber),
				bombDefused.Site,
				converters.Uint64ToString(bombDefused.DefuserSteamID64),
				bombDefused.DefuserName,
				converters.BoolToString(bombDefused.IsPlayerControllingBot),
				converters.Float64ToString(bombDefused.X),
				converters.Float64ToString(bombDefused.Y),
				converters.Float64ToString(bombDefused.Z),
				converters.IntToString(bombDefused.CounterTerroristAliveCount),
				converters.IntToString(bombDefused.TerroristAliveCount),
				match.Checksum,
			}
			lines = append(lines, line)
		}

		csv.WriteLinesIntoCsvFile(outputPath+"_bombs_defused.csv", lines)
	}

	var writeBombsExploded = func() {
		header := []string{
			"frame",
			"tick",
			"round",
			"site",
			"planter steamid",
			"planter name",
			"is player controlling bot",
			"x",
			"y",
			"z",
			"match checksum",
		}

		lines := [][]string{header}
		for _, bombExploded := range match.BombsExploded {
			line := []string{
				converters.IntToString(bombExploded.Frame),
				converters.IntToString(bombExploded.Tick),
				converters.IntToString(bombExploded.RoundNumber),
				bombExploded.Site,
				converters.Uint64ToString(bombExploded.PlanterSteamID64),
				bombExploded.PlanterName,
				converters.BoolToString(bombExploded.IsPlayerControllingBot),
				converters.Float64ToString(bombExploded.X),
				converters.Float64ToString(bombExploded.Y),
				converters.Float64ToString(bombExploded.Z),
				match.Checksum,
			}
			lines = append(lines, line)
		}

		csv.WriteLinesIntoCsvFile(outputPath+"_bombs_exploded.csv", lines)
	}

	var writeBombsPlantStart = func() {
		header := []string{
			"frame",
			"tick",
			"round",
			"site",
			"planter steamid",
			"planter name",
			"is player controlling bot",
			"x",
			"y",
			"z",
			"match checksum",
		}

		lines := [][]string{header}
		for _, bombPlantStart := range match.BombsPlantStart {
			line := []string{
				converters.IntToString(bombPlantStart.Frame),
				converters.IntToString(bombPlantStart.Tick),
				converters.IntToString(bombPlantStart.RoundNumber),
				bombPlantStart.Site,
				converters.Uint64ToString(bombPlantStart.PlanterSteamID64),
				bombPlantStart.PlanterName,
				converters.BoolToString(bombPlantStart.IsPlayerControllingBot),
				converters.Float64ToString(bombPlantStart.X),
				converters.Float64ToString(bombPlantStart.Y),
				converters.Float64ToString(bombPlantStart.Z),
				match.Checksum,
			}
			lines = append(lines, line)
		}

		csv.WriteLinesIntoCsvFile(outputPath+"_bombs_plant_start.csv", lines)
	}

	var writePlayersFlashed = func() {
		header := []string{
			"frame",
			"tick",
			"round",
			"duration",
			"flashed steamid",
			"flashed name",
			"flashed side",
			"is flashed controlling bot",
			"flasher steamid",
			"flasher name",
			"flasher side",
			"is flasher controlling bot",
			"match checksum",
		}

		lines := [][]string{header}
		for _, playerFlashed := range match.PlayersFlashed {
			line := []string{
				converters.IntToString(playerFlashed.Frame),
				converters.IntToString(playerFlashed.Tick),
				converters.IntToString(playerFlashed.RoundNumber),
				converters.Float32ToString(playerFlashed.Duration),
				converters.Uint64ToString(playerFlashed.FlashedSteamID64),
				playerFlashed.FlashedName,
				converters.TeamToString(playerFlashed.FlashedSide),
				converters.BoolToString(playerFlashed.IsFlashedControllingBot),
				converters.Uint64ToString(playerFlashed.FlasherSteamID64),
				playerFlashed.FlasherName,
				converters.TeamToString(playerFlashed.FlasherSide),
				converters.BoolToString(playerFlashed.IsFlasherControllingBot),
				match.Checksum,
			}
			lines = append(lines, line)
		}

		csv.WriteLinesIntoCsvFile(outputPath+"_players_flashed.csv", lines)
	}

	var writePlayersBuy = func() {
		header := []string{
			"frame",
			"tick",
			"round",
			"player steamid",
			"player side",
			"player name",
			"weapon name",
			"weapon type",
			"weapon unique id",
			"has refunded",
			"match checksum",
		}

		lines := [][]string{header}
		for _, playerBuy := range match.PlayersBuy {
			line := []string{
				converters.IntToString(playerBuy.Frame),
				converters.IntToString(playerBuy.Tick),
				converters.IntToString(playerBuy.RoundNumber),
				converters.Uint64ToString(playerBuy.PlayerSteamID64),
				converters.TeamToString(playerBuy.PlayerSide),
				playerBuy.PlayerName,
				playerBuy.WeaponName.String(),
				playerBuy.WeaponType.String(),
				playerBuy.WeaponUniqueID,
				converters.BoolToString(playerBuy.HasRefunded),
				match.Checksum,
			}
			lines = append(lines, line)
		}

		csv.WriteLinesIntoCsvFile(outputPath+"_players_buy.csv", lines)
	}

	var writeGrenadePositions = func() {
		header := []string{
			"frame",
			"tick",
			"round",
			"grenade id",
			"projectile id",
			"grenade name",
			"x",
			"y",
			"z",
			"thrower steamid",
			"thrower name",
			"thrower side",
			"thrower team name",
			"thrower velocity x",
			"thrower velocity y",
			"thrower velocity z",
			"thrower yaw",
			"thrower pitch",
			"match checksum",
		}

		lines := [][]string{header}
		for _, position := range match.GrenadePositions {
			line := []string{
				converters.IntToString(position.Frame),
				converters.IntToString(position.Tick),
				converters.IntToString(position.RoundNumber),
				position.GrenadeID,
				converters.Int64ToString(position.ProjectileID),
				position.GrenadeName.String(),
				converters.Float64ToString(position.X),
				converters.Float64ToString(position.Y),
				converters.Float64ToString(position.Z),
				converters.Uint64ToString(position.ThrowerSteamID64),
				position.ThrowerName,
				converters.TeamToString(position.ThrowerSide),
				position.ThrowerTeamName,
				converters.Float64ToString(position.ThrowerVelocityX),
				converters.Float64ToString(position.ThrowerVelocityY),
				converters.Float64ToString(position.ThrowerVelocityZ),
				converters.Float32ToString(position.ThrowerYaw),
				converters.Float32ToString(position.ThrowerPitch),
				match.Checksum,
			}
			lines = append(lines, line)
		}

		csv.WriteLinesIntoCsvFile(outputPath+"_grenade_positions.csv", lines)
	}

	var writeGrenadeBounces = func() {
		header := []string{
			"frame",
			"tick",
			"round",
			"grenade id",
			"projectile id",
			"grenade name",
			"x",
			"y",
			"z",
			"thrower steamid",
			"thrower name",
			"thrower side",
			"thrower team name",
			"thrower velocity x",
			"thrower velocity y",
			"thrower velocity z",
			"thrower yaw",
			"thrower pitch",
			"match checksum",
		}

		lines := [][]string{header}
		for _, position := range match.GrenadeBounces {
			line := []string{
				converters.IntToString(position.Frame),
				converters.IntToString(position.Tick),
				converters.IntToString(position.RoundNumber),
				position.GrenadeID,
				converters.Int64ToString(position.ProjectileID),
				position.GrenadeName.String(),
				converters.Float64ToString(position.X),
				converters.Float64ToString(position.Y),
				converters.Float64ToString(position.Z),
				converters.Uint64ToString(position.ThrowerSteamID64),
				position.ThrowerName,
				converters.TeamToString(position.ThrowerSide),
				position.ThrowerTeamName,
				converters.Float64ToString(position.ThrowerVelocityX),
				converters.Float64ToString(position.ThrowerVelocityY),
				converters.Float64ToString(position.ThrowerVelocityZ),
				converters.Float32ToString(position.ThrowerYaw),
				converters.Float32ToString(position.ThrowerPitch),
				match.Checksum,
			}
			lines = append(lines, line)
		}

		csv.WriteLinesIntoCsvFile(outputPath+"_grenade_bounces.csv", lines)
	}

	var writeGrenadeProjectilesDestroy = func() {
		header := []string{
			"frame",
			"tick",
			"round",
			"grenade id",
			"projectile id",
			"grenade name",
			"x",
			"y",
			"z",
			"thrower steamid",
			"thrower name",
			"thrower side",
			"thrower team name",
			"thrower velocity x",
			"thrower velocity y",
			"thrower velocity z",
			"thrower yaw",
			"thrower pitch",
			"match checksum",
		}

		lines := [][]string{header}
		for _, event := range match.GrenadeProjectilesDestroy {
			line := []string{
				converters.IntToString(event.Frame),
				converters.IntToString(event.Tick),
				converters.IntToString(event.RoundNumber),
				event.GrenadeID,
				converters.Int64ToString(event.ProjectileID),
				event.GrenadeName.String(),
				converters.Float64ToString(event.X),
				converters.Float64ToString(event.Y),
				converters.Float64ToString(event.Z),
				converters.Uint64ToString(event.ThrowerSteamID64),
				event.ThrowerName,
				converters.TeamToString(event.ThrowerSide),
				event.ThrowerTeamName,
				converters.Float64ToString(event.ThrowerVelocityX),
				converters.Float64ToString(event.ThrowerVelocityY),
				converters.Float64ToString(event.ThrowerVelocityZ),
				converters.Float32ToString(event.ThrowerYaw),
				converters.Float32ToString(event.ThrowerPitch),
				match.Checksum,
			}
			lines = append(lines, line)
		}

		csv.WriteLinesIntoCsvFile(outputPath+"_grenade_projectiles_destroy.csv", lines)
	}

	var writeSmokesStart = func() {
		header := []string{
			"frame",
			"tick",
			"round",
			"grenade id",
			"projectile id",
			"x",
			"y",
			"z",
			"thrower steamid",
			"thrower name",
			"thrower side",
			"thrower team name",
			"thrower velocity x",
			"thrower velocity y",
			"thrower velocity z",
			"thrower yaw",
			"thrower pitch",
			"match checksum",
		}

		lines := [][]string{header}
		for _, event := range match.SmokesStart {
			line := []string{
				converters.IntToString(event.Frame),
				converters.IntToString(event.Tick),
				converters.IntToString(event.RoundNumber),
				event.GrenadeID,
				converters.Int64ToString(event.ProjectileID),
				converters.Float64ToString(event.X),
				converters.Float64ToString(event.Y),
				converters.Float64ToString(event.Z),
				converters.Uint64ToString(event.ThrowerSteamID64),
				event.ThrowerName,
				converters.TeamToString(event.ThrowerSide),
				event.ThrowerTeamName,
				converters.Float64ToString(event.ThrowerVelocityX),
				converters.Float64ToString(event.ThrowerVelocityY),
				converters.Float64ToString(event.ThrowerVelocityZ),
				converters.Float32ToString(event.ThrowerYaw),
				converters.Float32ToString(event.ThrowerPitch),
				match.Checksum,
			}
			lines = append(lines, line)
		}

		csv.WriteLinesIntoCsvFile(outputPath+"_smokes_start.csv", lines)
	}

	var writeHeGrenadesExplode = func() {
		header := []string{
			"frame",
			"tick",
			"round",
			"grenade id",
			"projectile id",
			"x",
			"y",
			"z",
			"thrower steamid",
			"thrower name",
			"thrower side",
			"thrower team name",
			"thrower velocity x",
			"thrower velocity y",
			"thrower velocity z",
			"thrower yaw",
			"thrower pitch",
			"match checksum",
		}

		lines := [][]string{header}
		for _, event := range match.HeGrenadesExplode {
			line := []string{
				converters.IntToString(event.Frame),
				converters.IntToString(event.Tick),
				converters.IntToString(event.RoundNumber),
				event.GrenadeID,
				converters.Int64ToString(event.ProjectileID),
				converters.Float64ToString(event.X),
				converters.Float64ToString(event.Y),
				converters.Float64ToString(event.Z),
				converters.Uint64ToString(event.ThrowerSteamID64),
				event.ThrowerName,
				converters.TeamToString(event.ThrowerSide),
				event.ThrowerTeamName,
				converters.Float64ToString(event.ThrowerVelocityX),
				converters.Float64ToString(event.ThrowerVelocityY),
				converters.Float64ToString(event.ThrowerVelocityZ),
				converters.Float32ToString(event.ThrowerYaw),
				converters.Float32ToString(event.ThrowerPitch),
				match.Checksum,
			}
			lines = append(lines, line)
		}

		csv.WriteLinesIntoCsvFile(outputPath+"_he_grenades_explode.csv", lines)
	}

	var writeFlashbangsExplode = func() {
		header := []string{
			"frame",
			"tick",
			"round",
			"grenade id",
			"projectile id",
			"x",
			"y",
			"z",
			"thrower steamid",
			"thrower name",
			"thrower side",
			"thrower team name",
			"thrower velocity x",
			"thrower velocity y",
			"thrower velocity z",
			"thrower yaw",
			"thrower pitch",
			"match checksum",
		}

		lines := [][]string{header}
		for _, event := range match.FlashbangsExplode {
			line := []string{
				converters.IntToString(event.Frame),
				converters.IntToString(event.Tick),
				converters.IntToString(event.RoundNumber),
				event.GrenadeID,
				converters.Int64ToString(event.ProjectileID),
				converters.Float64ToString(event.X),
				converters.Float64ToString(event.Y),
				converters.Float64ToString(event.Z),
				converters.Uint64ToString(event.ThrowerSteamID64),
				event.ThrowerName,
				converters.TeamToString(event.ThrowerSide),
				event.ThrowerTeamName,
				converters.Float64ToString(event.ThrowerVelocityX),
				converters.Float64ToString(event.ThrowerVelocityY),
				converters.Float64ToString(event.ThrowerVelocityZ),
				converters.Float32ToString(event.ThrowerYaw),
				converters.Float32ToString(event.ThrowerPitch),
				match.Checksum,
			}
			lines = append(lines, line)
		}

		csv.WriteLinesIntoCsvFile(outputPath+"_flashbangs_explode.csv", lines)
	}

	var writeDecoysStart = func() {
		header := []string{
			"frame",
			"tick",
			"round",
			"grenade id",
			"projectile id",
			"x",
			"y",
			"z",
			"thrower steamid",
			"thrower name",
			"thrower side",
			"thrower team name",
			"thrower velocity x",
			"thrower velocity y",
			"thrower velocity z",
			"thrower yaw",
			"thrower pitch",
			"match checksum",
		}

		lines := [][]string{header}
		for _, event := range match.DecoysStart {
			line := []string{
				converters.IntToString(event.Frame),
				converters.IntToString(event.Tick),
				converters.IntToString(event.RoundNumber),
				event.GrenadeID,
				converters.Int64ToString(event.ProjectileID),
				converters.Float64ToString(event.X),
				converters.Float64ToString(event.Y),
				converters.Float64ToString(event.Z),
				converters.Uint64ToString(event.ThrowerSteamID64),
				event.ThrowerName,
				converters.TeamToString(event.ThrowerSide),
				event.ThrowerTeamName,
				converters.Float64ToString(event.ThrowerVelocityX),
				converters.Float64ToString(event.ThrowerVelocityY),
				converters.Float64ToString(event.ThrowerVelocityZ),
				converters.Float32ToString(event.ThrowerYaw),
				converters.Float32ToString(event.ThrowerPitch),
				match.Checksum,
			}
			lines = append(lines, line)
		}

		csv.WriteLinesIntoCsvFile(outputPath+"_decoys_start.csv", lines)
	}

	var writeInfernoPositions = func() {
		header := []string{
			"frame",
			"tick",
			"round",
			"thrower steamid",
			"thrower name",
			"unique id",
			"x",
			"y",
			"z",
			"convex hull 2d",
			"match checksum",
		}

		lines := [][]string{header}
		for _, position := range match.InfernoPositions {
			var convexHull2D string
			for index, point := range position.ConvexHull2D {
				startCharacter := ""
				if index > 0 {
					startCharacter = ","
				}
				convexHull2D += fmt.Sprintf("%s%f,%f", startCharacter, point.X, point.Y)
			}
			if convexHull2D == "" {
				convexHull2D = "0,0"
			}
			line := []string{
				converters.IntToString(position.Frame),
				converters.IntToString(position.Tick),
				converters.IntToString(position.RoundNumber),
				converters.Uint64ToString(position.ThrowerSteamID64),
				position.ThrowerName,
				converters.Int64ToString(position.UniqueID),
				converters.Float64ToString(position.X),
				converters.Float64ToString(position.Y),
				converters.Float64ToString(position.Z),
				convexHull2D,
				match.Checksum,
			}
			lines = append(lines, line)
		}

		csv.WriteLinesIntoCsvFile(outputPath+"_inferno_positions.csv", lines)
	}

	var writeChatMessages = func() {
		header := []string{
			"frame",
			"tick",
			"round",
			"steamid",
			"name",
			"message",
			"is alive",
			"side",
			"match checksum",
		}

		lines := [][]string{header}
		for _, chatMessage := range match.ChatMessages {
			line := []string{
				converters.IntToString(chatMessage.Frame),
				converters.IntToString(chatMessage.Tick),
				converters.IntToString(chatMessage.RoundNumber),
				converters.Uint64ToString(chatMessage.SenderSteamID64),
				chatMessage.SenderName,
				chatMessage.Message,
				converters.BoolToString(chatMessage.IsSenderAlive),
				converters.TeamToString(chatMessage.SenderSide),
				match.Checksum,
			}
			lines = append(lines, line)
		}

		csv.WriteLinesIntoCsvFile(outputPath+"_chat_messages.csv", lines)
	}

	var writeHostagePositions = func() {
		header := []string{
			"frame",
			"tick",
			"round",
			"x",
			"y",
			"z",
			"state",
			"match checksum",
		}

		lines := [][]string{header}
		for _, position := range match.HostagePositions {
			line := []string{
				converters.IntToString(position.Frame),
				converters.IntToString(position.Tick),
				converters.IntToString(position.RoundNumber),
				converters.Float64ToString(position.X),
				converters.Float64ToString(position.Y),
				converters.Float64ToString(position.Z),
				converters.ByteToString(byte(position.State)),
				match.Checksum,
			}
			lines = append(lines, line)
		}

		csv.WriteLinesIntoCsvFile(outputPath+"_hostage_positions.csv", lines)
	}

	var writeHostagePickUpStart = func() {
		header := []string{
			"frame",
			"tick",
			"round",
			"steamid",
			"is player controlling bot",
			"hostage entity id",
			"x",
			"y",
			"z",
			"match checksum",
		}

		lines := [][]string{header}
		for _, event := range match.HostagePickUpStart {
			line := []string{
				converters.IntToString(event.Frame),
				converters.IntToString(event.Tick),
				converters.IntToString(event.RoundNumber),
				converters.Uint64ToString(event.PlayerSteamID64),
				converters.BoolToString(event.IsPlayerControllingBot),
				converters.IntToString(event.HostageEntityId),
				converters.Float64ToString(event.X),
				converters.Float64ToString(event.Y),
				converters.Float64ToString(event.Z),
				match.Checksum,
			}
			lines = append(lines, line)
		}

		csv.WriteLinesIntoCsvFile(outputPath+"_hostage_pick_up_start.csv", lines)
	}

	var writeHostagePickedUp = func() {
		header := []string{
			"frame",
			"tick",
			"round",
			"steamid",
			"is player controlling bot",
			"hostage entity id",
			"x",
			"y",
			"z",
			"match checksum",
		}

		lines := [][]string{header}
		for _, event := range match.HostagePickedUp {
			line := []string{
				converters.IntToString(event.Frame),
				converters.IntToString(event.Tick),
				converters.IntToString(event.RoundNumber),
				converters.Uint64ToString(event.PlayerSteamID64),
				converters.BoolToString(event.IsPlayerControllingBot),
				converters.IntToString(event.HostageEntityId),
				converters.Float64ToString(event.X),
				converters.Float64ToString(event.Y),
				converters.Float64ToString(event.Z),
				match.Checksum,
			}
			lines = append(lines, line)
		}

		csv.WriteLinesIntoCsvFile(outputPath+"_hostage_picked_up.csv", lines)
	}

	var writeHostageRescued = func() {
		header := []string{
			"frame",
			"tick",
			"round",
			"steamid",
			"is player controlling bot",
			"hostage entity id",
			"x",
			"y",
			"z",
			"match checksum",
		}

		lines := [][]string{header}
		for _, event := range match.HostageRescued {
			line := []string{
				converters.IntToString(event.Frame),
				converters.IntToString(event.Tick),
				converters.IntToString(event.RoundNumber),
				converters.Uint64ToString(event.PlayerSteamID64),
				converters.BoolToString(event.IsPlayerControllingBot),
				converters.IntToString(event.HostageEntityId),
				converters.Float64ToString(event.X),
				converters.Float64ToString(event.Y),
				converters.Float64ToString(event.Z),
				match.Checksum,
			}
			lines = append(lines, line)
		}

		csv.WriteLinesIntoCsvFile(outputPath+"_hostage_rescued.csv", lines)
	}

	var functions = []func(){
		writeMatch,
		writeTeams,
		writePlayers,
		writePlayerPositions,
		writeShots,
		writeRounds,
		writeRoundEconomies,
		writeClutches,
		writeChickenDeaths,
		writeChickenPositions,
		writeDamages,
		writeKills,
		writeBombsPlanted,
		writeBombsDefuseStart,
		writeBombsDefused,
		writeBombsExploded,
		writeBombsPlantStart,
		writePlayersFlashed,
		writePlayersBuy,
		writeGrenadePositions,
		writeGrenadeBounces,
		writeGrenadeProjectilesDestroy,
		writeSmokesStart,
		writeHeGrenadesExplode,
		writeFlashbangsExplode,
		writeDecoysStart,
		writeInfernoPositions,
		writeChatMessages,
		writeHostagePositions,
		writeHostagePickUpStart,
		writeHostagePickedUp,
		writeHostageRescued,
	}
	var wg sync.WaitGroup

	for _, function := range functions {
		wg.Add(1)
		go func(function func()) {
			defer wg.Done()
			function()
		}(function)
	}

	wg.Wait()

	return nil
}
