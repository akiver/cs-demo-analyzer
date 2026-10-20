package constants

// HitGroup is the body part a player was hit in. The values are the game's hit group numbers.
type HitGroup byte

const (
	HitGroupGeneric  HitGroup = 0
	HitGroupHead     HitGroup = 1
	HitGroupChest    HitGroup = 2
	HitGroupStomach  HitGroup = 3
	HitGroupLeftArm  HitGroup = 4
	HitGroupRightArm HitGroup = 5
	HitGroupLeftLeg  HitGroup = 6
	HitGroupRightLeg HitGroup = 7
	HitGroupNeck     HitGroup = 8
	HitGroupGear     HitGroup = 10
)
