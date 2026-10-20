package constants

// Side is the side (team) a player plays on. The values are the game's team numbers.
type Side byte

const (
	SideUnassigned        Side = 0
	SideSpectators        Side = 1
	SideTerrorists        Side = 2
	SideCounterTerrorists Side = 3
)
