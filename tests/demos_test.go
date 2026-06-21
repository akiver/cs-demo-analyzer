package tests

import (
	"testing"

	"github.com/akiver/cs-demo-analyzer/pkg/api"
	"github.com/akiver/cs-demo-analyzer/pkg/api/constants"
	"github.com/akiver/cs-demo-analyzer/tests/assertion"
	"github.com/akiver/cs-demo-analyzer/tests/testsutils"
)

// demoTestCase describes a demo to analyze and snapshot.
type demoTestCase struct {
	testName string
	demoName string
	source   constants.DemoSource
}

var cs2DemoTestCases = []demoTestCase{
	// Contains a knife round
	{"5EPlay_g161_20231231135244670959707_2023_Mirage", "5eplay_g161_20231231135244670959707_2023_mirage", constants.DemoSourceFiveEPlay},
	// Contains players disconnection + BOTs
	{"5EPlay_g161_20240107155319726585099_2023_Nuke", "5eplay_g161-20240107155319726585099_2023_nuke", constants.DemoSourceFiveEPlay},
	// https://www.challengermode.com/s/CsgoAllstars/games/b4d31195-bae0-42c0-bbcc-08dbf6b17863
	{"ChallengerMode_6c306e56_8170_4092_b402_08dbf813e452_2023_Anubis", "challengermode_6c306e56-8170-4092-b402-08dbf813e452", constants.DemoSourceFaceIt},
	// https://www.hltv.org/matches/2367504/monte-vs-og-roobet-cup-2023
	{"Ebot_Monte_VS_OG_Roobet_Cup_2023_Anubis", "ebot_monte_vs_og_roobet_cup_2023_anubis", constants.DemoSourceEbot},
	// 2v2
	// https://esplay.com/m/nnQccdWWJtkc/midnight-vs-haha123xd
	{"Esplay_nnQccdWWJtkc_2025_Vertigo", "esplay_nnQccdWWJtkc_2025_vertigo", constants.DemoSourceEsplay},
	// 5v5 with overtimes.
	// https://esplay.com/m/ntfNCNcmKCQc/team-lina-vs-team-qara
	{"Esplay_ntfNCNcmKCQc_2025_Mirage", "esplay_ntfNCNcmKCQc_2025_mirage", constants.DemoSourceEsplay},
	// 5v5 with surrender.
	// https://esplay.com/m/nvBBvqNCfFHV/team-pytonorm-vs-team-shawty
	{"Esplay_nvBBvqNCfFHV_2025_Train", "esplay_nvBBvqNCfFHV_2025_train", constants.DemoSourceEsplay},
	// https://esportal.com/en/match/6008132
	{"Esportal_6008132_2023_Mirage", "esportal_6008132_2023_mirage", constants.DemoSourceEbot},
	// https://esportal.com/en/match/6045888
	{"Esportal_6045888_2024_Mirage", "esportal_6045888_2024_mirage", constants.DemoSourceEbot},
	{"Esportligaen_6a173af37c12e17acad1185e_2026_Mirage", "esportligaen_6a173af37c12e17acad1185e_2026_mirage", constants.DemoSourceEsportligaen},
	// https://cs2.fastcup.net/matches/11851975
	{"Fastcup_11851975_11876310_202312171749_Competitive_2023_Mirage", "fastcup_11851975_11876310_202312171749_2023_mirage", constants.DemoSourceFastcup},
	// https://www.hltv.org/stats/matches/mapstatsid/179910/aurora-vs-3dmax
	// - Contains a round backup restore at round 15
	// - Teams stay after knife round
	{"MatchZy_Aurora_vs_3dmax_m3_2024_Anubis", "matchzy_aurora_vs_3dmax_m3_anubis", constants.DemoSourceMatchZy},
	// https://www.hltv.org/matches/2373609/bleed-vs-parivision-skyesports-championship-2024
	// - 2 overtimes
	{"MatchZy_Bleed_vs_Parivision_2024_Mirage", "matchzy_bleed_vs_parivision_2024_mirage", constants.DemoSourceMatchZy},
	// - Recording started after the end of the knife round
	{"MatchZy_Iskandear_vs_Kirill_2024_Train", "matchzy_iskandear_vs_kirill_2024_train", constants.DemoSourceMatchZy},
	// The first round is restored 2 times
	{"MatchZy_PressurE_vs_cyphin_2024_Nuke", "matchzy_pressure_vs_cyphin_2024_nuke", constants.DemoSourceMatchZy},
	{"Pracc_CHHyStHPwEnGVsFmrZLWsDQwBX_2026_Nuke", "pracc_CHHyStHPwEnGVsFmrZLWsDQwBX_2026_nuke", constants.DemoSourcePracc},
	// Match with 1 overtime.
	// The server freezes during round 22, players are disconnected but the round still ends by bomb explosion.
	// The match is paused and resumed when all players are reconnected.
	// https://renown.gg/match/1363
	{"Renown_Match_1363_2025_Ancient", "renown_match_1363_2025_ancient", constants.DemoSourceRenown},
	// https://renown.gg/match/8
	{"Renown_Match_8_2025_Mirage", "renown_match_8_2025_mirage", constants.DemoSourceRenown},
}

var csgoDemoTestCases = []demoTestCase{
	// https://www.hltv.org/stats/matches/mapstatsid/151559/astralis-vs-saw
	{"ChallengerMode_Saw_VS_Astralis_M3_Paris_2023_CQ_Inferno", "challengermode_saw_vs_astralis_m3_Paris_2023_CQ_inferno", constants.DemoSourceFaceIt},
	// https://www.hltv.org/stats/matches/mapstatsid/151549/astralis-vs-saw
	{"ChallengerMode_Saw_VS_Astralis_M2_Paris_2023_CQ_Vertigo", "challengermode_saw_vs_astralis_m2_paris_2023_cq_vertigo", constants.DemoSourceFaceIt},
	// https://www.hltv.org/stats/matches/mapstatsid/40296/astralis-vs-envy
	{"Ebot_Astralis_VS_Envyus_Game_Show_Global_eSports_Cup_2016_Cache", "ebot_astralis_vs_envyus_game_show_global_esports_cup_2016_cache", constants.DemoSourceEbot},
	// Pause requested at the end of the round 3 and effective at the beginning of the round 4 (tick 32229).
	// The eBot backup is restored at tick 43560, that's when the round 4 really starts.
	// The round 6 is cancelled few seconds after it started, a backup is restored and the round 6 really start.
	// https://www.hltv.org/stats/matches/mapstatsid/38754/cloud9-vs-ninjas-in-pyjamas
	{"Ebot_Cloud9_VS_NIP_IEM_Oakland_2016_Train", "ebot_cloud9_vs_nip_iem_oakland_2016_train", constants.DemoSourceEbot},
	// E-frag starts as T, won the knife round and switch to CT.
	// There is then a warmup and the first round is cancelled during freezetime (!stop command).
	// The first round really starts after a restart by the eBot.
	// At the beginning of the second round, a pause is requested. The game is paused and the 2nd round starts at tick 74004.
	// https://www.hltv.org/stats/matches/mapstatsid/28406/e-fragnet-vs-faze
	{"Ebot_Efrag_Net_VS_Faze_IEM_Oakland_2016_Cache", "ebot_efrag_net_vs_faze_iem_oakland_2016_cache", constants.DemoSourceEbot},
	// Knife round with teams switch.
	// The game is stopped a few seconds after the beginning of the 1st round.
	// https://www.hltv.org/stats/matches/mapstatsid/28613/nerdrage-vs-galatics
	{"Ebot_Galatics_VS_Nerdrage_AlienTech_CSGO_League_Season1_2016_Cache", "ebot_galatics_vs_nerdrage_alientech_csgo_league_season1_2016_cache", constants.DemoSourceEbot},
	// Demo with 1 overtime.
	// Recording start after the end of the first round's freezetime.
	// Pauses at the end of the round 6, 8, 13, 15, 24, 26, 29.
	// After a few seconds at the end of the round 20, the round is cancelled (!stop) and really starts at tick 330201.
	// https://www.hltv.org/stats/matches/mapstatsid/42335/immortals-vs-north
	{"Ebot_Immortals_VS_North_IEM_Katowice_2017_Overpass", "ebot_immortals_vs_north_iem_katowice_2017_overpass", constants.DemoSourceEbot},
	// Knife round at the beginning and teams are swapped.
	// https://www.hltv.org/matches/2317273/optic-vs-faze-iem-oakland-2017
	{"Ebot_Optic_VS_Faze_IEM_Oakland_2017_Overpass", "ebot_optic_vs_faze_iem_oakland_2017_overpass", constants.DemoSourceEbot},
	// ESEA demo with warmup, half time break, coaches and pauses at freeze time.
	// https://www.hltv.org/stats/matches/mapstatsid/46797/clg-vs-liquid
	{"Esea_CLG_VS_Liquid_IEM_Cologne_2017_Cbble", "esea_clg_vs_liquid_iem_cologne_2017_cbble", constants.DemoSourceESEA},
	{"Esea_12283595_Cache", "esea_match_12283595_cache", constants.DemoSourceESEA},
	{"Faceit_2d199377f1794b488ac7e5714c754639_Mirage", "faceit_2d199377-f179-4b48-8ac7-e5714c754639_mirage", constants.DemoSourceFaceIt},
	// https://csgo.fastcup.net/matches/11853539
	{"Fastcup_11853539_11877805_2312172018_Competitive_2023_Dust2", "fastcup_11853539_11877805_2312172018_2023_dust2", constants.DemoSourceFastcup},
	// Demo containing 1 tactical timeout, multiple players disconnection, BOT overtake, suicides.
	{"Valve_Match730_003402256765125919145_0103110035_190_Nuke", "valve_match730_003402256765125919145_0103110035_190_nuke", constants.DemoSourceValve},
	// Demo with a surrender, bots and contains 1 tactical timeout.
	{"Valve_Match730_003408404295698088038_1541485657_202_Mirage", "valve_match730_003408404295698088038_1541485657_202_mirage", constants.DemoSourceValve},
	// Valve short match (MR 8) demo.
	{"Valve_Match730_003598554255364980910_1802085029_272_Ancient", "valve_match730_003598554255364980910_1802085029_272_ancient", constants.DemoSourceValve},
}

// runDemoTestCases runs every test case as a parallel subtest, reading the demos from the given game folder.
func runDemoTestCases(t *testing.T, game string, testCases []demoTestCase) {
	for _, tc := range testCases {
		// Prefix the subtest with the game folder so the suite can be filtered by game,
		// e.g. go test ./tests -run TestDemos/cs2.
		t.Run(game+"/"+tc.testName, func(t *testing.T) {
			t.Parallel()

			demoPath := testsutils.GetDemoPath(game, tc.demoName)
			match, err := api.AnalyzeDemo(demoPath, api.AnalyzeDemoOptions{
				Source: tc.source,
			})
			if err != nil {
				t.Fatal(err)
			}

			assertion.AssertMatchSnapshot(t, match, tc.demoName)
		})
	}
}

// TestDemos analyzes all demos and compares the result against its snapshot.
// Pass the -update flag to update the snapshots.
func TestDemos(t *testing.T) {
	runDemoTestCases(t, "cs2", cs2DemoTestCases)
	runDemoTestCases(t, "csgo", csgoDemoTestCases)
}
