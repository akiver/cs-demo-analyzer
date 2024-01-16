package fake

import "github.com/akiver/cs-demo-analyzer/pkg/api"

type FakePlayer struct {
	SteamID64     uint64
	Name          string
	Score         int
	Team          *api.Team
	KillCount     int
	AssistCount   int
	DeathCount    int
	MvpCount      int
	HeadshotCount int
	UtilityDamage int
}
