// Package hooks holds the unexported operations of pkg/model, which pkg/model registers when it's initialised.
// They are stored untyped because pkg/model can't import internal/modelaccess, which refers to pkg/model's types.
// Use internal/modelaccess instead of this package.
package hooks

var (
	ResetMatch             any
	ResetRound             any
	InitTeams              any
	SwapTeams              any
	ComputeResultStats     any
	DeleteIncompleteRounds any
	SetGameModeStr         any
	MatchJSON              any
	ResetPlayer            any
	SetPlayerMatch         any
	SetRoundMatch          any
	ChickenID              any
	SetChickenID           any
	HostageEntity          any
	SetHostageEntity       any
)
