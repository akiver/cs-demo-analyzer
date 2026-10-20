package constants

// HostageState is a hostage's state. The values are the game's hostage state numbers.
type HostageState byte

const (
	HostageStateIdle            HostageState = 0
	HostageStateBeingUntied     HostageState = 1
	HostageStateGettingPickedUp HostageState = 2
	HostageStateBeingCarried    HostageState = 3
	HostageStateFollowingPlayer HostageState = 4
	HostageStateGettingDropped  HostageState = 5
	HostageStateRescued         HostageState = 6
	HostageStateDead            HostageState = 7
)
