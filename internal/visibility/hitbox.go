package visibility

import (
	"math"

	"github.com/golang/geo/r3"
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/common"
)

// Demos don't contain players' bones. The body of a player is approximated with points placed from their origin
// (feet), yaw and eyes height. The points come from the real hitbox capsules of the game, captured in game for many
// poses and averaged per template, see cmd/hitbox-templates. A template depends on the stance, whether the player is
// moving and the kind of weapon held, which are the parts of the pose a demo tells.

// HitboxPoint is a point of the body in the frame of the player.
type HitboxPoint struct {
	Forward float64 // Units along the yaw direction
	Right   float64 // Units to the right of the yaw direction
	Up      float64 // Units above the feet, for the eyes height of the template
}

// HitboxTemplate is the body of a player in one kind of pose.
type HitboxTemplate struct {
	// Mean eyes height of the poses of the template, the points are scaled vertically by the actual eyes height.
	EyesHeight float64
	// Ordered from the most to the least likely visible part (head, torso, arms, legs) to exit early when peeking.
	Points []HitboxPoint
}

type Stance int

const (
	StanceStanding Stance = iota
	StanceCrouching
	StanceAirborne
	stanceCount
)

type WeaponKind int

const (
	WeaponKindOther WeaponKind = iota // Knife, C4, zeus, no weapon
	WeaponKindPistol
	WeaponKindRifle
	WeaponKindScoped // A rifle or sniper while scoped
	WeaponKindSMG    // SMGs and heavy weapons
	WeaponKindGrenade
	weaponKindCount
)

const (
	// Below that eyes height a player is crouching, the eyes are at 64 standing and 46 fully crouched.
	crouchingEyesHeight = 60.0
	// Above that horizontal speed a player is moving, the legs and arms swing.
	movingSpeed = 40.0
)

// HitboxTemplateCount is the number of templates, one per stance, moving state and weapon kind.
const HitboxTemplateCount = int(stanceCount) * 2 * int(weaponKindCount)

func StanceOf(eyesHeight float64, isAirborne bool) Stance {
	if isAirborne {
		return StanceAirborne
	}
	if eyesHeight < crouchingEyesHeight {
		return StanceCrouching
	}

	return StanceStanding
}

// IsMoving returns true when the player moves horizontally, jumping or falling in place doesn't swing the legs and arms.
func IsMoving(velocity r3.Vector) bool {
	return math.Hypot(velocity.X, velocity.Y) > movingSpeed
}

func WeaponKindOf(weapon *common.Equipment, isScoped bool) WeaponKind {
	if weapon == nil {
		return WeaponKindOther
	}

	switch weapon.Class() {
	case common.EqClassPistols:
		return WeaponKindPistol
	case common.EqClassRifle:
		if isScoped {
			return WeaponKindScoped
		}
		return WeaponKindRifle
	case common.EqClassSMG, common.EqClassHeavy:
		return WeaponKindSMG
	case common.EqClassGrenade:
		return WeaponKindGrenade
	default:
		return WeaponKindOther
	}
}

// HitboxTemplateIndex returns the index of the template of a pose in the generated templates.
func HitboxTemplateIndex(stance Stance, isMoving bool, weaponKind WeaponKind) int {
	moving := 0
	if isMoving {
		moving = 1
	}

	return (int(stance)*2+moving)*int(weaponKindCount) + int(weaponKind)
}

// SplitHitboxTemplateIndex returns the pose of the template at an index, the inverse of HitboxTemplateIndex.
func SplitHitboxTemplateIndex(index int) (Stance, bool, WeaponKind) {
	weaponKind := WeaponKind(index % int(weaponKindCount))
	isMoving := index/int(weaponKindCount)%2 == 1
	stance := Stance(index / int(weaponKindCount) / 2)

	return stance, isMoving, weaponKind
}

// HitboxTemplateName describes the pose of a template, for the generated code.
func HitboxTemplateName(index int) string {
	// Sized by the counts so that a stance or weapon kind added without a name is caught by the tests.
	stances := [stanceCount]string{"standing", "crouching", "airborne"}
	weaponKinds := [weaponKindCount]string{"other", "pistol", "rifle", "scoped", "SMG", "grenade"}
	stance, isMoving, weaponKind := SplitHitboxTemplateIndex(index)
	movement := "still"
	if isMoving {
		movement = "moving"
	}

	return stances[stance] + ", " + movement + ", " + weaponKinds[weaponKind]
}

// GetHitboxTemplate returns the template of the pose.
func GetHitboxTemplate(stance Stance, isMoving bool, weaponKind WeaponKind) *HitboxTemplate {
	return &hitboxTemplates[HitboxTemplateIndex(stance, isMoving, weaponKind)]
}
