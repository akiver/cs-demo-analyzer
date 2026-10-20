package api

import (
	"github.com/akiver/cs-demo-analyzer/internal/demo"
	"github.com/akiver/cs-demo-analyzer/pkg/api/constants"
)

func newMatch(source constants.DemoSource, demoInfo *demo.Demo) Match {
	game := constants.CSGO
	if demoInfo.IsSource2() {
		// The build number of CS2 when it was publicly available is 9832, everything below is coming from the limited test.
		if demoInfo.BuildNumber < 9832 {
			game = constants.CS2LT
		} else {
			game = constants.CS2
		}
	}

	match := Match{
		Checksum:                  demoInfo.Checksum,
		Source:                    source,
		Game:                      game,
		Type:                      "GOTV", // By default assume it's a GOTV demo, it will be updated during parsing.
		ShareCode:                 demoInfo.ShareCode,
		TickCount:                 demoInfo.TickCount,
		Date:                      demoInfo.Date,
		DemoFilePath:              demoInfo.FilePath,
		DemoFileName:              demoInfo.FileName,
		TickRate:                  demoInfo.TickRate,
		MapName:                   demoInfo.MapName,
		FrameRate:                 demoInfo.FrameRate,
		Duration:                  demoInfo.Duration,
		ServerName:                demoInfo.ServerName,
		ClientName:                demoInfo.ClientName,
		BuildNumber:               demoInfo.BuildNumber,
		NetworkProtocol:           demoInfo.NetworkProtocol,
		PlayersBySteamID:          make(map[uint64]*Player),
		Rounds:                    []*Round{},
		Kills:                     []*Kill{},
		Clutches:                  []*Clutch{},
		BombsPlanted:              []*BombPlanted{},
		BombsDefused:              []*BombDefused{},
		BombsExploded:             []*BombExploded{},
		BombsPlantStart:           []*BombPlantStart{},
		BombsDefuseStart:          []*BombDefuseStart{},
		Damages:                   []*Damage{},
		Shots:                     []*Shot{},
		PlayersFlashed:            []*PlayerFlashed{},
		PlayersBuy:                []*PlayerBuy{},
		PlayerEconomies:           []*PlayerEconomy{},
		PlayerPositions:           []*PlayerPosition{},
		GrenadePositions:          []*GrenadePosition{},
		GrenadeBounces:            []*GrenadeBounce{},
		InfernoPositions:          []*InfernoPosition{},
		SmokesStart:               []*SmokeStart{},
		DecoysStart:               []*DecoyStart{},
		HeGrenadesExplode:         []*HeGrenadeExplode{},
		FlashbangsExplode:         []*FlashbangExplode{},
		HostagePickUpStart:        []*HostagePickUpStart{},
		HostagePickedUp:           []*HostagePickedUp{},
		HostageRescued:            []*HostageRescued{},
		HostagePositions:          []*HostagePosition{},
		ChickenPositions:          []*ChickenPosition{},
		ChatMessages:              []*ChatMessage{},
		Coaches:                   []*Coach{},
		ChickenDeaths:             []*ChickenDeath{},
		GrenadeProjectilesDestroy: []*GrenadeProjectileDestroy{},
	}

	match.InitTeams()

	return match
}
