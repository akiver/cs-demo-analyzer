package shared

import (
	"github.com/akiver/cs-demo-analyzer/internal/demo"
	"github.com/akiver/cs-demo-analyzer/pkg/api/constants"
	"github.com/akiver/cs-demo-analyzer/pkg/model"
)

// NewMatch returns the match of a demo, filled with what is known before parsing it.
func NewMatch(source constants.DemoSource, demoInfo *demo.Demo) model.Match {
	game := constants.CSGO
	if demoInfo.IsSource2() {
		// The build number of CS2 when it was publicly available is 9832, everything below is coming from the limited test.
		if demoInfo.BuildNumber < 9832 {
			game = constants.CS2LT
		} else {
			game = constants.CS2
		}
	}

	match := model.Match{
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
		PlayersBySteamID:          make(map[uint64]*model.Player),
		Rounds:                    []*model.Round{},
		Kills:                     []*model.Kill{},
		Clutches:                  []*model.Clutch{},
		BombsPlanted:              []*model.BombPlanted{},
		BombsDefused:              []*model.BombDefused{},
		BombsExploded:             []*model.BombExploded{},
		BombsPlantStart:           []*model.BombPlantStart{},
		BombsDefuseStart:          []*model.BombDefuseStart{},
		Damages:                   []*model.Damage{},
		Shots:                     []*model.Shot{},
		PlayersFlashed:            []*model.PlayerFlashed{},
		PlayersBuy:                []*model.PlayerBuy{},
		PlayerEconomies:           []*model.PlayerEconomy{},
		PlayerPositions:           []*model.PlayerPosition{},
		GrenadePositions:          []*model.GrenadePosition{},
		GrenadeBounces:            []*model.GrenadeBounce{},
		InfernoPositions:          []*model.InfernoPosition{},
		SmokesStart:               []*model.SmokeStart{},
		DecoysStart:               []*model.DecoyStart{},
		HeGrenadesExplode:         []*model.HeGrenadeExplode{},
		FlashbangsExplode:         []*model.FlashbangExplode{},
		HostagePickUpStart:        []*model.HostagePickUpStart{},
		HostagePickedUp:           []*model.HostagePickedUp{},
		HostageRescued:            []*model.HostageRescued{},
		HostagePositions:          []*model.HostagePosition{},
		ChickenPositions:          []*model.ChickenPosition{},
		ChatMessages:              []*model.ChatMessage{},
		Coaches:                   []*model.Coach{},
		ChickenDeaths:             []*model.ChickenDeath{},
		GrenadeProjectilesDestroy: []*model.GrenadeProjectileDestroy{},
	}

	match.InitTeams()

	return match
}
