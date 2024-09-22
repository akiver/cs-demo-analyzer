package api

import (
	"embed"
	"errors"
	"fmt"
	"math"
	"os"
	"strconv"

	"github.com/akiver/cs-demo-analyzer/internal/converters"
	d "github.com/akiver/cs-demo-analyzer/internal/demo"
	"github.com/akiver/cs-demo-analyzer/internal/slice"
	"github.com/akiver/cs-demo-analyzer/internal/strings"
	"github.com/akiver/cs-demo-analyzer/pkg/api/constants"
	"github.com/golang/geo/r3"
	dem "github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs"
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/common"
	events "github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/events"
	st "github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/sendtables"
)

const (
	// Maximum number of seconds between a teammate death and a possible revenge kill to be considered as a trade kill.
	tradeKillDelaySeconds = 5
	// Number of seconds that players equipment value will be computed after the end of freezetimes.
	// It's not done at the end of freezetime because players are still able to buy a few seconds after it.
	// Using the seconds from mp_buytime which is 20 seconds may leads to inaccurate results since during this timelapse
	// players may have throw grenades, been killed...
	equipmentValueDelaySeconds = 7
)

type Analyzer struct {
	parser             dem.Parser
	match              *Match
	currentRound       *Round
	matchStarted       func() bool
	postProcess        func(analyzer *Analyzer)
	isSource2          bool
	isFirstRoundOfHalf bool
	// Flag to handle demos with missing round end events.
	isRoundEndDetected bool
	// How many seconds players are allowed to buy when a round start, this variable is updated if the ConVar
	// "mp_buytime" has been detected. Used to know what players have bought at the beginning of rounds.
	buyTimeSeconds int
	// Store the last freezetime end tick detected, used to compute players equipments value for each round.
	// It's not done at the end of freezetime because players are still able to buy a few seconds after it.
	lastFreezeTimeEndTick int
	clutch1               *Clutch
	clutch2               *Clutch
	// Store bomb plant position to retrieve it when it explodes.
	bombPlantPosition r3.Vector
	// Store the last grenade thrown by player in order to create a link between the grenade unique ID and the
	// projectile unique ID.
	// The grenade ID is not the same as the projectile ID and when a player throws a grenade, its projectile is
	// created only a few ticks later.
	lastGrenadeThrownByPlayer map[uint64]*Shot
	// Used to detect which player is untying an hostage when the hostage's state changed to "Being untied".
	// Because several hostages can be untied at the same time, we keep track of players that started untying hostages
	// to detect which player is untying an hostage in case of consecutive events.
	playersUntyingAnHostage map[uint64]int
	chickenEntities         []st.Entity
}

//go:embed event-list-dump/*.bin
var eventListFolder embed.FS

func getGameEventListBinForProtocol(networkProtocol int) ([]byte, error) {
	switch {
	case networkProtocol < 13992:
		return eventListFolder.ReadFile("event-list-dump/13990.bin")
	case networkProtocol <= 13992:
		return eventListFolder.ReadFile("event-list-dump/13992.bin")
	default:
		return eventListFolder.ReadFile("event-list-dump/14023.bin")
	}
}

type AnalyzeDemoOptions struct {
	IncludePositions bool
	Source           constants.DemoSource
}

func analyzeDemo(demoPath string, options AnalyzeDemoOptions) (*Match, error) {
	if options.Source != "" {
		err := ValidateDemoSource(options.Source)
		if err != nil {
			return nil, err
		}
	}

	demo, err := d.GetDemoFromPath(demoPath)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(demoPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf(fmt.Sprintf("demo file \"%s\" not found", demoPath))
		}
		return nil, err
	}
	defer file.Close()

	parserConfig := dem.DefaultParserConfig
	parserConfig.NetMessageDecryptionKey = demo.NetMessageDecryptionPublicKey
	parserConfig.DisableMimicSource1Events = demo.Type == constants.DemoTypePOV

	gameEventListBin, err := getGameEventListBinForProtocol(demo.NetworkProtocol)
	if err != nil {
		return nil, err
	}
	parserConfig.Source2FallbackGameEventListBin = gameEventListBin

	parser := dem.NewParserWithConfig(file, parserConfig)
	defer parser.Close()

	_, err = parser.ParseHeader()
	if err != nil {
		return nil, err
	}

	source := options.Source
	if source == "" {
		source = d.GetDemoSource(demo)
	}

	if demo.IsSource2() && demo.Type == constants.DemoTypePOV {
		return nil, errors.New("cs2 pov demos are not supported (CS2POVDemosNotSupported)")
	}

	match := newMatch(source, demo)

	analyzer := &Analyzer{
		parser:                    parser,
		match:                     &match,
		isSource2:                 demo.IsSource2(),
		isFirstRoundOfHalf:        true,
		isRoundEndDetected:        false,
		buyTimeSeconds:            20, // Default mp_buytime value is 20
		lastFreezeTimeEndTick:     -1,
		bombPlantPosition:         r3.Vector{},
		lastGrenadeThrownByPlayer: make(map[uint64]*Shot),
		playersUntyingAnHostage:   make(map[uint64]int),
		postProcess:               defaultPostProcess,
	}

	analyzer.currentRound = &Round{
		analyzer:           analyzer,
		Number:             1,
		StartTick:          1,
		StartFrame:         1,
		FreezeTimeEndFrame: -1,
		FreezeTimeEndTick:  -1,
		TeamAName:          match.TeamA.Name,
		TeamBName:          match.TeamB.Name,
		TeamASide:          *match.TeamA.CurrentSide,
		TeamBSide:          *match.TeamB.CurrentSide,
	}

	analyzer.registerCommonHandlers(options.IncludePositions)

	switch source {
	case constants.DemoSourceFaceIt:
		createFaceItAnalyzer(analyzer)
	case constants.DemoSourceESEA:
		createEseaAnalyzer(analyzer)
	case constants.DemoSourceEbot:
		createEbotAnalyzer(analyzer)
	case constants.DemoSourceChallengermode:
		createChallengermodeAnalyzer(analyzer)
	case constants.DemoSourceEsportal:
		createEsportalAnalyzer(analyzer)
	case constants.DemoSourceCEVO:
		return nil, errors.New("cevo demos are not supported (CevoNotSupported)")
	case constants.DemoSourceFastcup:
		createFastcupAnalyzer(analyzer)
	case constants.DemoSourceFiveEPlay:
		createFiveEPlayAnalyzer(analyzer)
	case constants.DemoSourceGamersclub:
		// Looks like they use an eBot fork but rounds are not detected properly.
		return nil, errors.New("gamersclub demos are not supported (GamersClubNotSupported)")
	case constants.DemoSourcePopFlash:
		// TODO notImplemented Unsure how demos from PopFlash are nowadays, need to find some to test.
		// Even latest CSGO demos from PopFlash may not really work because their recording system has probably changed
		// and may not be compatible with the Valve one anymore (used to be the V2).
		if demo.IsSource2() {
			return nil, errors.New("cs2 PopFlash demos are not supported (PopFlashNotSupported)")
		}
		createValveAnalyzer(analyzer)
	case constants.DemoSourceValve, constants.DemoSourcePerfectWorld, constants.DemoSourceESL:
		createValveAnalyzer(analyzer)
	default:
		return nil, errors.New("unknown demo source, please specify the source with the -source flag (UnknownSource)")
	}

	err = parser.ParseToEnd()
	// Do not stop if the demo is corrupted, usually the error occurs at the end of the parsing.
	// Depending on how far we were able to parse the demo we may still have data.
	isCorruptedDemo := errors.Is(err, dem.ErrUnexpectedEndOfDemo)
	if err != nil && !isCorruptedDemo {
		return nil, err
	}

	// Required for CS2 demos, the following data are available only at the end of the parsing in the CDemoFileInfo message.
	header := parser.Header()
	match.TickCount = header.PlaybackTicks
	match.Duration = header.PlaybackTime
	if match.Duration > 0 {
		match.FrameRate = float64(header.PlaybackFrames) / header.PlaybackTime.Seconds()
	}

	analyzer.postProcess(analyzer)
	match.computeResultStats()

	return &match, nil
}

func AnalyzeDemo(demoPath string, options AnalyzeDemoOptions) (*Match, error) {
	match, err := analyzeDemo(demoPath, options)

	return match, err
}

type AnalyzeAndExportDemoOptions struct {
	IncludePositions bool
	Source           constants.DemoSource
	Format           constants.ExportFormat
	MinifyJSON       bool
}

func AnalyzeAndExportDemo(demoPath string, outputPath string, options AnalyzeAndExportDemoOptions) error {
	var err error
	if options.Format != "" {
		err = ValidateExportFormat(options.Format)
		if err != nil {
			return err
		}
	}

	match, err := analyzeDemo(demoPath, AnalyzeDemoOptions{
		IncludePositions: options.IncludePositions,
		Source:           options.Source,
	})

	if err != nil {
		return err
	}

	switch options.Format {
	case "csv":
		err = exportMatchToCSV(match, outputPath)
	case "json":
		err = exportMatchToJSON(match, outputPath, options.MinifyJSON)
	case "csdm":
		err = exportMatchForCSDM(match, outputPath)
	}

	return err
}

func (analyzer *Analyzer) currentTick() int {
	return analyzer.parser.GameState().IngameTick()
}

func (analyzer *Analyzer) reset() {
	analyzer.isFirstRoundOfHalf = true
	analyzer.isRoundEndDetected = false
	analyzer.lastFreezeTimeEndTick = -1
	analyzer.lastGrenadeThrownByPlayer = make(map[uint64]*Shot)
	analyzer.playersUntyingAnHostage = make(map[uint64]int)
	analyzer.chickenEntities = nil
	analyzer.clutch1 = nil
	analyzer.clutch2 = nil
	analyzer.match.reset()
	analyzer.updateTeamNames()
	for _, player := range analyzer.match.PlayersBySteamID {
		player.reset()
	}

	analyzer.currentRound = &Round{
		analyzer:               analyzer,
		Number:                 1,
		StartTick:              1,
		StartFrame:             1,
		FreezeTimeEndFrame:     -1,
		FreezeTimeEndTick:      -1,
		TeamASide:              *analyzer.match.TeamA.CurrentSide,
		TeamBSide:              *analyzer.match.TeamB.CurrentSide,
		TeamAName:              analyzer.match.TeamA.Name,
		TeamBName:              analyzer.match.TeamB.Name,
		weaponsBoughtUniqueIds: nil,
	}
}

func (analyzer *Analyzer) registerPlayer(player *common.Player, teamState *common.TeamState) {
	match := analyzer.match
	if _, alreadyExists := match.PlayersBySteamID[player.SteamID64]; alreadyExists {
		return
	}

	if teamState == nil || player.SteamID64 == 0 {
		return
	}

	if player.Entity == nil || (analyzer.isSource2 && player.PlayerPawnEntity() == nil) {
		return
	}

	// Ignore coaches
	if prop, exists := player.Entity.PropertyValue("m_iCoachingTeam"); exists {
		coachTeam := prop.Int()
		if coachTeam == int(common.TeamCounterTerrorists) || coachTeam == int(common.TeamTerrorists) {
			return
		}
	}

	playerTeam := teamState.Team()
	if playerTeam == common.TeamSpectators || playerTeam == common.TeamUnassigned {
		return
	}

	var team *Team
	if *match.TeamA.CurrentSide == playerTeam {
		team = match.TeamA
	} else {
		team = match.TeamB
	}

	color, _ := player.ColorOrErr()
	rank := player.Rank()
	userID := 0
	if player.UserID <= math.MaxUint16 {
		userID = player.UserID & 0xff
	}

	newPlayer := &Player{
		match:              match,
		SteamID64:          player.SteamID64,
		UserID:             userID,
		Name:               strings.ReplaceUTF8ByteSequences(player.Name),
		Team:               team,
		CrosshairShareCode: player.CrosshairCode(),
		Color:              color,
		RankType:           player.RankType(),
		Rank:               rank,
		OldRank:            rank,
		WinCount:           player.CompetitiveWins(),
	}
	match.PlayersBySteamID[player.SteamID64] = newPlayer

	// TODO PR demoinfocs but m_bIsLookingAtWeapon doesn't updates on CS2 demos :/
	var inspectingWeaponProp st.Property
	if analyzer.isSource2 {
		inspectingWeaponProp = player.PlayerPawnEntity().Property("m_pWeaponServices.m_bIsLookingAtWeapon")
	} else if player.Entity != nil {
		inspectingWeaponProp = player.Entity.Property("m_bIsLookingAtWeapon")
	}

	if inspectingWeaponProp != nil {
		inspectingWeaponProp.OnUpdate(func(val st.PropertyValue) {
			if val.BoolVal() {
				newPlayer.InspectWeaponCount += 1
			}
		})
	}
}

func (analyzer *Analyzer) registerUnknownPlayers() {
	for _, player := range analyzer.parser.GameState().Participants().All() {
		analyzer.registerPlayer(player, player.TeamState)
	}
}

func (analyzer *Analyzer) updateTeamNames() {
	gameState := analyzer.parser.GameState()
	match := analyzer.match
	teamNameA := strings.ReplaceUTF8ByteSequences(gameState.Team(*match.TeamA.CurrentSide).ClanName())
	teamNameB := strings.ReplaceUTF8ByteSequences(gameState.Team(*match.TeamB.CurrentSide).ClanName())
	if teamNameA == teamNameB {
		return
	}

	if teamNameA != "" {
		match.TeamA.Name = teamNameA
		analyzer.currentRound.TeamAName = teamNameA
	}
	if teamNameB != "" {
		match.TeamB.Name = teamNameB
		analyzer.currentRound.TeamBName = teamNameB
	}
}

// updatePlayersScores updates the score of all PlayersBySteamID that aren't spectating or unassigned from the parser information.
func (analyzer *Analyzer) updatePlayersScores() {
	players := analyzer.parser.GameState().Participants().Playing()
	for _, playingPlayer := range players {
		// Ignores 0 values because the initial score is 0 and some sources like ESEA reset it when the match is actually not done.
		if player, exist := analyzer.match.PlayersBySteamID[playingPlayer.SteamID64]; exist && playingPlayer.Score() != 0 {
			player.Score = playingPlayer.Score()
			player.MvpCount = playingPlayer.MVPs()
		}
	}
}

func (analyzer *Analyzer) updatePlayersCurrentTeam() {
	match := analyzer.match
	parser := analyzer.parser
	for _, player := range match.PlayersBySteamID {
		for _, playingPlayer := range parser.GameState().Participants().Playing() {
			if playingPlayer.SteamID64 != player.SteamID64 {
				continue
			}

			playerTeam := playingPlayer.TeamState.Team()
			if playerTeam == common.TeamSpectators || playerTeam == common.TeamUnassigned {
				continue
			}

			if *player.Team.CurrentSide == playerTeam {
				continue
			}

			if *match.TeamA.CurrentSide == playerTeam {
				player.Team = match.TeamA
			} else {
				player.Team = match.TeamB
			}
		}
	}
}

func (analyzer *Analyzer) processMatchStart() {
	parser := analyzer.parser
	match := analyzer.match
	analyzer.updatePlayersCurrentTeam()

	currentRound := analyzer.currentRound
	currentRound.StartFrame = parser.CurrentFrame()
	currentRound.StartTick = analyzer.currentTick()
	currentRound.TeamASide = *match.TeamA.CurrentSide
	currentRound.TeamBSide = *match.TeamB.CurrentSide
	analyzer.updateTeamNames()
	analyzer.createPlayersEconomies()
}

// This creates a Round, reset analyzer data related to the current round and set the new round as the current round.
func (analyzer *Analyzer) createRound() {
	analyzer.clutch1 = nil
	analyzer.clutch2 = nil
	analyzer.lastGrenadeThrownByPlayer = make(map[uint64]*Shot)
	analyzer.playersUntyingAnHostage = map[uint64]int{}
	roundNumber := analyzer.currentRound.Number + 1

	analyzer.currentRound = newRound(roundNumber, analyzer)
	analyzer.createPlayersEconomies()
}

// secondsHasPassedSinceTick indicates if x seconds have passed since a specific tick.
func (analyzer *Analyzer) secondsHasPassedSinceTick(seconds float64, tick int) bool {
	parser := analyzer.parser
	return float64(analyzer.currentTick()-tick)*parser.TickTime().Seconds() >= seconds
}

func (analyzer *Analyzer) isKnifeRound() bool {
	playerWithEmptyWalletCount := 0
	players := analyzer.parser.GameState().Participants()
	for _, player := range players.Playing() {
		if player.Money() == 0 {
			playerWithEmptyWalletCount += 1
		}
	}

	return playerWithEmptyWalletCount == len(players.Playing())
}

func (analyzer *Analyzer) handleRoundEnd(winnerTeam common.Team) {
	analyzer.isRoundEndDetected = true
	analyzer.isFirstRoundOfHalf = false
	analyzer.currentRound.WinnerName = analyzer.match.Team(winnerTeam).Name
	analyzer.currentRound.WinnerSide = *analyzer.match.Team(winnerTeam).CurrentSide
	analyzer.currentRound.EndTick = analyzer.currentTick()
	analyzer.currentRound.EndOfficiallyTick = analyzer.currentTick()
	analyzer.currentRound.EndFrame = analyzer.parser.CurrentFrame()
	analyzer.currentRound.EndOfficiallyFrame = analyzer.parser.CurrentFrame()

	if analyzer.currentRound.FreezeTimeEndFrame == -1 {
		analyzer.currentRound.FreezeTimeEndFrame = analyzer.parser.CurrentFrame()
	}
	if analyzer.currentRound.FreezeTimeEndTick == -1 {
		analyzer.currentRound.FreezeTimeEndTick = analyzer.currentTick()
	}

	if analyzer.clutch1 != nil {
		analyzer.clutch1.HasWon = winnerTeam == analyzer.clutch1.Side
	}
	if analyzer.clutch2 != nil {
		analyzer.clutch2.HasWon = winnerTeam == analyzer.clutch2.Side
	}

	analyzer.updatePlayersScores()
}

func (analyzer *Analyzer) createPlayersEconomies() {
	match := analyzer.match

	match.PlayerEconomies = slice.Filter(match.PlayerEconomies, func(economy *PlayerEconomy, index int) bool {
		return economy.RoundNumber != analyzer.currentRound.Number
	})

	for _, player := range analyzer.parser.GameState().Participants().Playing() {
		economy := newPlayerEconomy(analyzer, player)
		match.PlayerEconomies = append(match.PlayerEconomies, economy)
	}
}

func (analyzer *Analyzer) computePlayersEconomies() {
	for _, player := range analyzer.parser.GameState().Participants().Playing() {
		economy := analyzer.match.GetPlayerEconomyAtRound(player.Name, player.SteamID64, analyzer.currentRound.Number)
		if economy == nil {
			continue
		}

		economy.updateValues(analyzer, player)
	}
}

func (analyzer *Analyzer) defaultRoundStartHandler(event events.RoundStart) {
	if !analyzer.matchStarted() {
		return
	}

	// No Rounds have been added yet, don't create a new one in this case, it's still the first round.
	if len(analyzer.match.Rounds) == 0 {
		return
	}

	analyzer.createRound()
}

func (analyzer *Analyzer) defaultRoundFreezetimeChangedHandler(event events.RoundFreezetimeChanged) {
	// It may not be accurate to create players economy on round start because it's possible to buy
	// a few ticks before the round start event and so may results in incorrect values.
	// Do it when the freeze time starts, it's updated just before round start events.
	if event.NewIsFreezetime {
		analyzer.createPlayersEconomies()
	} else {
		analyzer.currentRound.FreezeTimeEndTick = analyzer.currentTick()
		analyzer.currentRound.FreezeTimeEndFrame = analyzer.parser.CurrentFrame()
		analyzer.lastFreezeTimeEndTick = analyzer.currentTick()
	}
}

func (analyzer *Analyzer) defaultRoundEndOfficiallyHandler(event events.RoundEndOfficial) {
	if !analyzer.matchStarted() {
		return
	}

	analyzer.match.Rounds = append(analyzer.match.Rounds, analyzer.currentRound)
}

func (analyzer *Analyzer) registerCommonHandlers(includePositions bool) {
	parser := analyzer.parser
	match := analyzer.match

	parser.RegisterEventHandler(func(event events.TickRateInfoAvailable) {
		match.TickRate = parser.TickRate()
	})

	parser.RegisterEventHandler(func(event events.POVRecordingPlayerDetected) {
		match.Type = "POV"
	})

	parser.RegisterEventHandler(func(event events.PlayerTeamChange) {
		analyzer.registerPlayer(event.Player, event.NewTeamState)
	})

	parser.RegisterEventHandler(func(event events.RoundStart) {
		analyzer.registerUnknownPlayers()
	})

	parser.RegisterEventHandler(func(event events.ConVarsUpdated) {
		for varName, varValue := range event.UpdatedConVars {
			switch varName {
			case "mp_buytime":
				buyTime, _ := strconv.Atoi(varValue)
				if buyTime > 0 {
					analyzer.buyTimeSeconds = buyTime
				}
			case "mp_maxrounds":
				maxRounds, _ := strconv.Atoi(varValue)
				// Sanity check for ESEA demos. mp_maxrounds is set high values such as 999.
				// The real max rounds will be detected at the end of the analyze based on the scores.
				if maxRounds > 0 && maxRounds < 99 {
					match.MaxRounds = maxRounds
				}
			}
		}
	})

	parser.RegisterEventHandler(func(event events.OtherDeath) {
		if event.OtherType != "chicken" {
			return
		}

		if event.Killer == nil {
			fmt.Println(analyzer.currentTick(), "A chicken has been killed but the killer is nil")
			return
		}

		chickenDeath := newChickenDeath(analyzer.parser.CurrentFrame(), analyzer.currentTick(), analyzer.currentRound.Number, event.Killer.SteamID64, equipmentToWeaponName[event.Weapon.Type])
		analyzer.match.ChickenDeaths = append(analyzer.match.ChickenDeaths, chickenDeath)
	})

	parser.RegisterEventHandler(func(event events.ItemPickup) {
		if !analyzer.matchStarted() || event.Player == nil || event.Player.IsBot || !event.Player.IsInBuyZone() {
			return
		}

		weaponType := event.Weapon.Type
		isNotBuyableWeapon := weaponType == common.EqKnife || weaponType == common.EqBomb || weaponType == common.EqUnknown || weaponType == common.EqWorld
		if isNotBuyableWeapon {
			return
		}

		isDefaultTerroristsPistol := event.Player.Team == common.TeamTerrorists && weaponType == common.EqGlock
		if isDefaultTerroristsPistol {
			return
		}

		isDefaultCounterTerroristsPistol := event.Player.Team == common.TeamCounterTerrorists && (weaponType == common.EqUSP || weaponType == common.EqP2000)
		if isDefaultCounterTerroristsPistol {
			return
		}

		currentRound := analyzer.currentRound
		isDroppedWeapon := slice.Contains(currentRound.weaponsBoughtUniqueIds, event.Weapon.UniqueID2().String())
		if isDroppedWeapon {
			return
		}

		isBuyTimeEnded := currentRound.secondsPassedSinceRoundStart(analyzer.buyTimeSeconds)
		if isBuyTimeEnded {
			return
		}

		currentRound.weaponsBoughtUniqueIds = append(currentRound.weaponsBoughtUniqueIds, event.Weapon.UniqueID2().String())
		match.PlayersBuy = append(match.PlayersBuy, newPlayerBuy(analyzer, event))
	})

	parser.RegisterEventHandler(func(event events.PlayerHurt) {
		if !analyzer.matchStarted() || event.Player == nil {
			return
		}

		damage := newDamageFromGameEvent(analyzer, event)
		if damage != nil {
			match.Damages = append(match.Damages, damage)
		}
	})

	parser.RegisterEventHandler(func(event events.FrameDone) {
		shouldComputeEconomy := analyzer.lastFreezeTimeEndTick != -1 && analyzer.secondsHasPassedSinceTick(equipmentValueDelaySeconds, analyzer.lastFreezeTimeEndTick)
		if shouldComputeEconomy {
			analyzer.computePlayersEconomies()
			analyzer.currentRound.computeTeamsEconomy()
			analyzer.lastFreezeTimeEndTick = -1
		}
	})

	if includePositions {
		parser.RegisterEventHandler(func(event events.FrameDone) {
			if !analyzer.matchStarted() {
				return
			}

			for _, chickenEntity := range analyzer.chickenEntities {
				chickenPosition := newChickenPositionFromEntity(analyzer, chickenEntity)
				match.ChickenPositions = append(match.ChickenPositions, chickenPosition)
			}

			for _, projectile := range parser.GameState().GrenadeProjectiles() {
				position := newGrenadePositionFromProjectile(analyzer, projectile)
				if position != nil {
					match.GrenadePositions = append(match.GrenadePositions, position)
				}
			}

			for _, inferno := range parser.GameState().Infernos() {
				infernoPosition := newInfernoPositionFromInferno(analyzer, inferno)
				if infernoPosition != nil {
					match.InfernoPositions = append(match.InfernoPositions, infernoPosition)
				}
			}

			for _, player := range parser.GameState().Participants().Playing() {
				playerPosition := newPlayerPosition(analyzer, player)
				match.PlayerPositions = append(match.PlayerPositions, playerPosition)
			}

			for _, hostage := range parser.GameState().Hostages() {
				hostagePosition := newHostagePositionFromHostage(analyzer, hostage)
				match.HostagePositions = append(match.HostagePositions, hostagePosition)
			}
		})
	}

	parser.RegisterEventHandler(func(event events.RoundEnd) {
		if !analyzer.matchStarted() {
			return
		}

		analyzer.handleRoundEnd(event.Winner)
		analyzer.currentRound.updateEndReasonFromRoundEndEvent(event, analyzer.match.MapName)
	})

	missingGameEventDescriptorsWarnCount := 0
	parser.RegisterEventHandler(func(event events.ParserWarn) {
		fmt.Println(event)
		if event.Type != events.WarnTypeGameEventBeforeDescriptors {
			return
		}

		// CS2 demos recorded on a server with sv_hibernate_when_empty 0 doesn't contain the msg CMsgSource1LegacyGameEventList.
		// It looks like demos coming from https://dathost.net (Esportal seems to use it) use this setting and so
		// these demos are not analyzable.
		// While the workaround in the following PR work, it has not been merged because it's not a good long term
		// solution and as it's CS2 bug, Valve should fix this issue.
		// https://github.com/markus-wa/demoinfocs-golang/pull/460
		missingGameEventDescriptorsWarnCount += 1
		if missingGameEventDescriptorsWarnCount >= 20 {
			panic(errors.New("missing game event descriptors (ErrMissingGameEventDescriptors)"))
		}
	})

	parser.RegisterEventHandler(func(event events.RoundEndOfficial) {
		if !analyzer.matchStarted() {
			return
		}

		analyzer.currentRound.EndOfficiallyFrame = parser.CurrentFrame()
		analyzer.currentRound.EndOfficiallyTick = analyzer.currentTick()
	})

	parser.RegisterEventHandler(func(event events.Kill) {
		if !analyzer.matchStarted() {
			return
		}

		var killerSteamID64 uint64
		if event.Killer != nil {
			killerSteamID64 = event.Killer.SteamID64
		}
		var victimSteamID64 uint64
		if event.Victim != nil {
			victimSteamID64 = event.Victim.SteamID64
		}
		kill := newKillFromGameEvent(analyzer, event)
		if kill != nil {
			match.Kills = append(match.Kills, kill)
		}

		if analyzer.clutch1 != nil {
			clutcherSteamId := analyzer.clutch1.ClutcherSteamID64
			if clutcherSteamId == victimSteamID64 {
				analyzer.clutch1.ClutcherSurvived = false
			} else if clutcherSteamId == killerSteamID64 {
				analyzer.clutch1.ClutcherKillCount += 1
			}
		}

		if analyzer.clutch2 != nil {
			clutcherSteamId := analyzer.clutch2.ClutcherSteamID64
			if clutcherSteamId == victimSteamID64 {
				analyzer.clutch2.ClutcherSurvived = false
			} else if clutcherSteamId == killerSteamID64 {
				analyzer.clutch2.ClutcherKillCount += 1
			}
		}

		var counterTerroristsAlive []*common.Player
		for _, player := range parser.GameState().TeamCounterTerrorists().Members() {
			if player.IsAlive() && victimSteamID64 != player.SteamID64 {
				counterTerroristsAlive = append(counterTerroristsAlive, player)
			}
		}

		var terroristsAlive []*common.Player
		for _, player := range parser.GameState().TeamTerrorists().Members() {
			if player.IsAlive() && victimSteamID64 != player.SteamID64 {
				terroristsAlive = append(terroristsAlive, player)
			}
		}

		// 1vX detection
		if analyzer.clutch1 == nil && (len(counterTerroristsAlive) == 1 || len(terroristsAlive) == 1) {
			var side common.Team
			var opponentCount int
			var clutcher *common.Player
			if len(counterTerroristsAlive) == 1 {
				side = common.TeamCounterTerrorists
				opponentCount = len(terroristsAlive)
				clutcher = counterTerroristsAlive[0]
			} else {
				side = common.TeamTerrorists
				opponentCount = len(counterTerroristsAlive)
				clutcher = terroristsAlive[0]
			}

			analyzer.clutch1 = newClutch(analyzer, clutcher, side, opponentCount)
			match.Clutches = append(match.Clutches, analyzer.clutch1)
		}

		// 1v1 detection
		if analyzer.clutch1 != nil && len(counterTerroristsAlive) == 1 && len(terroristsAlive) == 1 {
			var side common.Team
			var clutcher *common.Player
			if analyzer.clutch1.Side == common.TeamCounterTerrorists {
				side = common.TeamTerrorists
				clutcher = terroristsAlive[0]
			} else {
				side = common.TeamCounterTerrorists
				clutcher = counterTerroristsAlive[0]
			}

			analyzer.clutch2 = newClutch(analyzer, clutcher, side, 1)
			match.Clutches = append(match.Clutches, analyzer.clutch2)
		}
	})

	parser.RegisterEventHandler(func(event events.HeExplode) {
		if !analyzer.matchStarted() {
			return
		}

		heGrenadeExplode := newHeGrenadeExplodeFromGameEvent(analyzer, event)
		if heGrenadeExplode != nil {
			match.HeGrenadesExplode = append(match.HeGrenadesExplode, heGrenadeExplode)
		}
	})

	parser.RegisterEventHandler(func(event events.WeaponFire) {
		if !analyzer.matchStarted() {
			return
		}

		shot := newShot(analyzer, event)
		if shot == nil {
			return
		}

		if event.Weapon.Class() == common.EqClassGrenade {
			analyzer.lastGrenadeThrownByPlayer[shot.PlayerSteamID64] = shot
		}

		match.Shots = append(match.Shots, shot)
	})

	parser.RegisterEventHandler(func(event events.BombPlanted) {
		if !analyzer.matchStarted() {
			return
		}

		analyzer.bombPlantPosition = event.Player.LastAlivePosition

		bombPlanted := newBombPlanted(analyzer, event)
		match.BombsPlanted = append(match.BombsPlanted, bombPlanted)
	})

	parser.RegisterEventHandler(func(event events.BombDefused) {
		if !analyzer.matchStarted() {
			return
		}

		bombDefused := newBombDefused(analyzer, event)
		match.BombsDefused = append(match.BombsDefused, bombDefused)
		analyzer.currentRound.EndReason = events.RoundEndReasonBombDefused
	})

	parser.RegisterEventHandler(func(event events.BombExplode) {
		if !analyzer.matchStarted() {
			return
		}

		bombExploded := newBombExploded(analyzer, event)
		match.BombsExploded = append(match.BombsExploded, bombExploded)
		// Update the round end reason because sometimes when the bomb exploded, the round end event indicates
		// RoundEndReasonTerroristsWin instead of RoundEndReasonTargetBombed.
		analyzer.currentRound.EndReason = events.RoundEndReasonTargetBombed
	})

	parser.RegisterEventHandler(func(event events.BombPlantBegin) {
		if !analyzer.matchStarted() {
			return
		}

		bombPlantStart := newBombPlantStart(analyzer, event)
		match.BombsPlantStart = append(match.BombsPlantStart, bombPlantStart)
	})

	parser.RegisterEventHandler(func(event events.BombDefuseStart) {
		if !analyzer.matchStarted() {
			return
		}

		bombDefuseStart := newBombDefuseStart(analyzer, event.Player)
		match.BombsDefuseStart = append(match.BombsDefuseStart, bombDefuseStart)
	})

	parser.RegisterEventHandler(func(event events.PlayerFlashed) {
		if !analyzer.matchStarted() || event.Player == nil || event.Attacker == nil || event.Player.IsBot {
			return
		}

		playerFlashed := newPlayerFlashed(analyzer, event)
		match.PlayersFlashed = append(match.PlayersFlashed, playerFlashed)
	})

	parser.RegisterEventHandler(func(event events.FlashExplode) {
		if !analyzer.matchStarted() {
			return
		}

		flashbangExplode := newFlashbangExplodeFromGameEvent(analyzer, event)
		if flashbangExplode != nil {
			match.FlashbangsExplode = append(match.FlashbangsExplode, flashbangExplode)
		}
	})

	parser.RegisterEventHandler(func(event events.GrenadeProjectileThrow) {
		if !analyzer.matchStarted() {
			return
		}

		projectile := event.Projectile
		if projectile == nil {
			fmt.Println("Projectile nil in grenade projectile throw event")
			return
		}

		thrower := projectile.Thrower
		if thrower == nil {
			fmt.Println("Thrower nil in grenade projectile throw event, falling back to owner")
			thrower = projectile.WeaponInstance.Owner
			if thrower == nil {
				fmt.Println("Owner nil in grenade projectile throw event")
				return
			}
		}

		lastGrenadeShotExist := analyzer.lastGrenadeThrownByPlayer[thrower.SteamID64]
		if lastGrenadeShotExist == nil {
			fmt.Println("A projectile throw event occurred whereas its weapon fired event didn't occurred.")
		} else {
			lastGrenadeShotExist.ProjectileID = projectile.UniqueID()
			delete(analyzer.lastGrenadeThrownByPlayer, thrower.SteamID64)
		}
	})

	parser.RegisterEventHandler(func(event events.GrenadeProjectileBounce) {
		if !analyzer.matchStarted() {
			return
		}

		grenadeBounce := newGrenadeBounceFromProjectile(analyzer, event.Projectile)
		if grenadeBounce != nil {
			match.GrenadeBounces = append(match.GrenadeBounces, grenadeBounce)
		}
	})

	parser.RegisterEventHandler(func(event events.GrenadeProjectileDestroy) {
		if !analyzer.matchStarted() {
			return
		}

		grenadeProjectileDestroy := newGrenadeProjectileDestroyFromProjectile(analyzer, event.Projectile)
		if grenadeProjectileDestroy != nil {
			match.GrenadeProjectilesDestroy = append(match.GrenadeProjectilesDestroy, grenadeProjectileDestroy)
		}
	})

	parser.RegisterEventHandler(func(event events.DecoyStart) {
		if !analyzer.matchStarted() {
			return
		}

		decoyStart := newDecoyStartFromGameEvent(analyzer, event)
		if decoyStart != nil {
			match.DecoysStart = append(match.DecoysStart, decoyStart)
		}
	})

	parser.RegisterEventHandler(func(event events.SmokeStart) {
		if !analyzer.matchStarted() {
			return
		}

		smokeStart := newSmokeStartFromGameEvent(analyzer, event)
		if smokeStart != nil {
			match.SmokesStart = append(match.SmokesStart, smokeStart)
		}
	})

	parser.RegisterEventHandler(func(event events.ScoreUpdated) {
		// Update players score when a team's score has been updated to handle the last kill of a match.
		// The AnnouncementWinPanelMatch event is triggered before ScoreUpdated and we can't use the Kill event because
		// the player's score at this point is the current player's score (before the kill).
		analyzer.updatePlayersScores()
	})

	parser.RegisterEventHandler(func(event events.RankUpdate) {
		match.IsRanked = true
		if player, exist := match.PlayersBySteamID[event.SteamID64()]; exist {
			player.Rank = event.RankNew
			player.OldRank = event.RankOld
			player.WinCount = event.WinCount
		}
	})

	parser.RegisterEventHandler(func(event events.HostageStateChanged) {
		if !analyzer.matchStarted() {
			return
		}

		switch event.NewState {
		case common.HostageStateBeingUntied:
			for _, player := range parser.GameState().Participants().Playing() {
				playerHasNoUntyingInProgress := analyzer.playersUntyingAnHostage[player.SteamID64] == 0
				if player.IsGrabbingHostage() && playerHasNoUntyingInProgress {
					analyzer.playersUntyingAnHostage[player.SteamID64] = event.Hostage.Entity.ID()
					hostagePickedUpStart := newHostagePickupStart(analyzer, player, event.Hostage)
					match.HostagePickUpStart = append(match.HostagePickUpStart, hostagePickedUpStart)
					break
				}
			}
		case common.HostageStateGettingPickedUp:
			if event.Hostage.Leader() != nil {
				delete(analyzer.playersUntyingAnHostage, event.Hostage.Leader().SteamID64)
				hostagePickedUp := newHostagePickedUp(analyzer, event.Hostage)
				match.HostagePickedUp = append(match.HostagePickedUp, hostagePickedUp)
			}
		case common.HostageStateIdle:
			// Case when a player started to untie an hostage and cancelled.
			for _, player := range parser.GameState().Participants().Playing() {
				playerWasUntyingTheHostage := analyzer.playersUntyingAnHostage[player.SteamID64] == event.Hostage.Entity.ID()
				if !player.IsGrabbingHostage() && playerWasUntyingTheHostage {
					delete(analyzer.playersUntyingAnHostage, player.SteamID64)
					break
				}
			}
		case common.HostageStateBeingCarried:
			if event.Hostage.Leader() != nil {
				delete(analyzer.playersUntyingAnHostage, event.Hostage.Leader().SteamID64)
			}
		}
	})

	parser.RegisterEventHandler(func(event events.HostageRescued) {
		if !analyzer.matchStarted() || event.Hostage.Leader() == nil {
			return
		}

		hostageRescued := newHostageRescued(analyzer, event.Hostage)
		match.HostageRescued = append(match.HostageRescued, hostageRescued)
	})

	parser.RegisterEventHandler(func(event events.RoundMVPAnnouncement) {
		if !analyzer.matchStarted() {
			return
		}

		analyzer.updatePlayersScores()
	})

	parser.RegisterEventHandler(func(event events.ItemRefund) {
		if !analyzer.matchStarted() {
			return
		}

		// Loop in reverse order to find the last purchase of the player and flag it as refunded.
		for i := len(match.PlayersBuy) - 1; i >= 0; i-- {
			buy := match.PlayersBuy[i]

			if buy.PlayerSteamID64 != event.Player.SteamID64 || buy.WeaponName != equipmentToWeaponName[event.Weapon.Type] {
				continue
			}

			if buy.RoundNumber != analyzer.currentRound.Number {
				break
			}

			buy.HasRefunded = true
		}
	})

	parser.RegisterEventHandler(func(event events.OvertimeNumberChanged) {
		match.OvertimeCount = event.NewCount
	})

	parser.RegisterEventHandler(func(event events.ChatMessage) {
		chatMessage := newChatMessageFromGameEvent(analyzer, event)
		match.ChatMessages = append(match.ChatMessages, chatMessage)
	})

	parser.RegisterEventHandler(func(event events.ConVarsUpdated) {
		for conVarName, conVarValue := range event.UpdatedConVars {
			if conVarName == "game_type" {
				match.GameType = constants.GameType(converters.StringToInt(conVarValue))
			} else if conVarName == "game_mode" {
				match.GameMode = constants.GameMode(converters.StringToInt(conVarValue))
			}
		}
	})

	parser.RegisterEventHandler(func(event events.DataTablesParsed) {
		serverClasses := parser.ServerClasses()
		serverClasses.FindByName("CChicken").OnEntityCreated(func(entity st.Entity) {
			analyzer.chickenEntities = append(analyzer.chickenEntities, entity)
		})

		// We don't use the event TeamSideSwitch to detect teams switch because it's triggered several times at the same tick with POV demos.
		var currentGamePhase common.GamePhase = common.GamePhaseInit
		parser.RegisterEventHandler(func(event events.GamePhaseChanged) {
			if !analyzer.matchStarted() {
				return
			}
			newGamePhase := event.NewGamePhase
			valueChanged := newGamePhase != currentGamePhase
			currentGamePhase = newGamePhase
			if !valueChanged {
				return
			}

			switch newGamePhase {
			case common.GamePhaseTeamSideSwitch:
				{
					match.swapTeams()
					// round_start is triggered before this entity update with CS:GO demos only. Update the current teams side here.
					// It doesn't happen with CS2 demos because the parser dispatch events in a more expected way.
					if !analyzer.isSource2 {
						analyzer.currentRound.TeamASide = *match.TeamA.CurrentSide
						analyzer.currentRound.TeamBSide = *match.TeamB.CurrentSide
					}
					analyzer.isFirstRoundOfHalf = true
				}
			}
		})

		getCurrentRoundBombDefusedEvent := func() *BombDefused {
			bombsDefused := analyzer.match.BombsDefused
			bombDefusedCount := len(bombsDefused)
			if bombDefusedCount > 0 {
				lastBombDefused := bombsDefused[bombDefusedCount-1]
				if lastBombDefused.RoundNumber == analyzer.currentRound.Number {
					return lastBombDefused
				}
			}

			return nil
		}

		if !analyzer.isSource2 {
			// CSGO workaround to detect missing bomb defused events and update the round end reason.
			// Both events may be missing with old CSGO demos.
			// We don't have the problem with CS2 demos because the parser dispatch bomb events using props updates
			// rather than game events.
			// TODO move it to demoinfocs?
			serverClasses.FindByName("CPlantedC4").OnEntityCreated(func(bombEntity st.Entity) {
				// Old CSGO demos don't have these properties (~ before end of 2018).
				siteProp := bombEntity.Property("m_nBombSite")
				defuserProp := bombEntity.Property("m_hBombDefuser")
				if siteProp == nil || defuserProp == nil {
					return
				}

				siteNumber := siteProp.Value().Int()
				site := events.BomsiteUnknown
				if siteNumber == 0 {
					site = events.BombsiteA
				} else if siteNumber == 1 {
					site = events.BombsiteB
				}

				var defuser *common.Player
				defuserProp.OnUpdate(func(val st.PropertyValue) {
					defuser = parser.GameState().Participants().FindByHandle64(uint64(val.Int()))
				})

				bombEntity.Property("m_bBombDefused").OnUpdate(func(val st.PropertyValue) {
					isDefused := val.BoolVal()
					if !isDefused {
						return
					}

					bombsDefused := getCurrentRoundBombDefusedEvent()
					// Don't create a new bomb defused event if it already exists, it means the game event has been
					// triggered.
					if bombsDefused != nil {
						return
					}

					bombDefused := newBombDefused(analyzer, events.BombDefused{
						BombEvent: events.BombEvent{
							Player: defuser,
							Site:   site,
						},
					})
					match.BombsDefused = append(match.BombsDefused, bombDefused)
					analyzer.currentRound.EndReason = events.RoundEndReasonBombDefused
				})
			})
		}

		serverClasses.FindByName("CCSGameRulesProxy").OnEntityCreated(func(entity st.Entity) {
			// Fallback to detect rounds end for demos missing round end events.
			// ! Don't use m_eRoundWinReason because it's not available with old demos (at least 2014)
			roundWinStatusProperty := entity.Property("cs_gamerules_data.m_iRoundWinStatus")
			if roundWinStatusProperty == nil {
				roundWinStatusProperty = entity.Property("m_pGameRules.m_iRoundWinStatus")
			}
			if roundWinStatusProperty != nil {
				roundWinStatusProperty.OnUpdate(func(val st.PropertyValue) {
					roundWinStatus := byte(val.Int())
					if roundWinStatus == byte(constants.RoundWinStatusUnassigned) {
						analyzer.isRoundEndDetected = false
						return
					}

					if !analyzer.matchStarted() || analyzer.isRoundEndDetected || roundWinStatus == byte(constants.RoundWinStatusDraw) {
						return
					}

					winnerSide := common.TeamCounterTerrorists
					if roundWinStatus == byte(constants.RoundWinStatusTWon) {
						winnerSide = common.TeamTerrorists
					}

					analyzer.handleRoundEnd(winnerSide)

					// We don't have access to the round end reason here, try to detect it based on current game state.
					currentRound := analyzer.currentRound
					if winnerSide == common.TeamCounterTerrorists {
						currentRound.EndReason = events.RoundEndReasonCTWin

						// If a T is still alive it means CTs won because Ts ran out of time
						for _, player := range parser.GameState().Participants().TeamMembers(common.TeamTerrorists) {
							if player.IsAlive() {
								currentRound.EndReason = events.RoundEndReasonTargetSaved
								break
							}
						}

						// If there is a bomb defused event, it means CTs won
						currentRoundBombDefusedEvent := getCurrentRoundBombDefusedEvent()
						if currentRoundBombDefusedEvent != nil {
							currentRound.EndReason = events.RoundEndReasonBombDefused
						}
					} else {
						currentRound.EndReason = events.RoundEndReasonTerroristsWin

						bombsExploded := analyzer.match.BombsExploded
						bombExplodedCount := len(bombsExploded)
						if bombExplodedCount > 0 {
							lastBombExploded := bombsExploded[bombExplodedCount-1]
							if lastBombExploded.RoundNumber == currentRound.Number {
								currentRound.EndReason = events.RoundEndReasonTargetBombed
							}
						}
					}
				})
			}

			if analyzer.isSource2 {
				abortedProp := entity.Property("m_pGameRules.m_nMatchAbortedEarlyReason")
				if abortedProp != nil {
					abortedProp.OnUpdate(func(val st.PropertyValue) {
						reason := val.Int()
						if reason == 0 {
							return
						}

						// TODO notImplemented Find a demo with a VAC live ban to get the correct value
						fmt.Println("Match aborted with reason", reason)
						analyzer.match.HasVacLiveBan = true
					})
				} else {
					// Old demos
					entity.Property("m_pGameRules.m_bMatchAbortedDueToPlayerBan").OnUpdate(func(val st.PropertyValue) {
						analyzer.match.HasVacLiveBan = val.BoolVal()
					})
				}
			}

			var lastTimeoutUpdateTick = -1
			onTimeoutUpdate := func(val st.PropertyValue) {
				currentTick := analyzer.currentTick()
				if lastTimeoutUpdateTick != currentTick && val.Float() == 0 && !analyzer.isRoundEndDetected {
					analyzer.currentRound.StartTick = currentTick
					analyzer.currentRound.StartFrame = parser.CurrentFrame()
				}
				lastTimeoutUpdateTick = currentTick
			}
			// Demos before the CSGO update that added teams timeout doesn't contain this prop.
			// In this case round start values will come from round_start event which is fine as timeout didn't exist.
			counterTerroristTimeoutRemainingProperty := entity.Property("cs_gamerules_data.m_flCTTimeOutRemaining")
			if counterTerroristTimeoutRemainingProperty == nil {
				counterTerroristTimeoutRemainingProperty = entity.Property("m_pGameRules.m_flCTTimeOutRemaining")
			}
			if counterTerroristTimeoutRemainingProperty != nil {
				counterTerroristTimeoutRemainingProperty.OnUpdate(onTimeoutUpdate)
			}
			terroristTimeoutRemainingProperty := entity.Property("cs_gamerules_data.m_flTerroristTimeOutRemaining")
			if terroristTimeoutRemainingProperty == nil {
				terroristTimeoutRemainingProperty = entity.Property("m_pGameRules.m_flTerroristTimeOutRemaining")
			}
			if terroristTimeoutRemainingProperty != nil {
				terroristTimeoutRemainingProperty.OnUpdate(onTimeoutUpdate)
			}
		})
	})
}

func defaultPostProcess(analyzer *Analyzer) {
	match := analyzer.match
	currentRound := analyzer.currentRound
	if len(match.Rounds) < currentRound.Number {
		match.Rounds = append(match.Rounds, currentRound)
	}
}
