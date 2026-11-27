package api

import (
	"math"

	csmath "github.com/akiver/cs-demo-analyzer/internal/math"
	"github.com/akiver/cs-demo-analyzer/internal/visibility"
	"github.com/golang/geo/r3"
)

// Demos don't contain players' bones, the body of a player is approximated with the points of a hitbox template placed
// from its origin (feet), eyes height (lower when crouching) and yaw, see internal/visibility/hitbox.go. A player is
// visible when any of these points can be seen.

const (
	// Eyes height when standing, used when the view offset is not available.
	defaultEyesHeight = 64.0
	// Horizontal and vertical tangents of half the field of view with the default fov (90° for a 4:3 screen) on a 16:9
	// screen, 106.26° horizontally and 73.74° vertically.
	defaultFovTanHorizontal = 4.0 / 3.0
	defaultFovTanVertical   = 0.75
	// Scoped weapons field of view, first and second zoom levels.
	scopedFovDegrees     = 40.0
	fullScopedFovDegrees = 15.0
	// A smoke is approximated by an ellipsoid centered above its detonation position.
	smokeHorizontalRadius = 144.0
	smokeVerticalRadius   = 112.0
	smokeCenterHeight     = 64.0
)

type playerSnapshot struct {
	recorded bool
	isAlive  bool
	origin   r3.Vector
	eyes     r3.Vector
	yaw      float64 // Degrees
	pitch    float64 // Degrees, positive when looking down
	// Tangents of half the field of view, depending on the zoom level.
	fovTanHorizontal float64
	fovTanVertical   float64
	// The player is too blinded by a flashbang to see.
	isBlinded bool
	// Body of the player for its pose, the default standing pose when nil.
	hitboxTemplate *visibility.HitboxTemplate
}

var defaultHitboxTemplate = visibility.GetHitboxTemplate(visibility.StanceStanding, false, visibility.WeaponKindRifle)

func (snapshot *playerSnapshot) hitboxPoints() []visibility.HitboxPoint {
	return snapshot.template().Points
}

func (snapshot *playerSnapshot) template() *visibility.HitboxTemplate {
	if snapshot.hitboxTemplate == nil {
		return defaultHitboxTemplate
	}

	return snapshot.hitboxTemplate
}

// fovTangents returns the horizontal and vertical tangents of half the field of view for a 4:3 fov in degrees.
func fovTangents(fovDegrees float64) (float64, float64) {
	tan := math.Tan(fovDegrees / 2 * math.Pi / 180)

	return tan * defaultFovTanHorizontal, tan * defaultFovTanVertical
}

func (snapshot *playerSnapshot) eyesHeight() float64 {
	eyesHeight := snapshot.eyes.Z - snapshot.origin.Z
	if eyesHeight <= 0 {
		return defaultEyesHeight
	}

	return eyesHeight
}

// hitboxPoint returns the world position of a point of the player's template, scaled vertically by the eyes height.
func (snapshot *playerSnapshot) hitboxPoint(point visibility.HitboxPoint) r3.Vector {
	yaw := snapshot.yaw * math.Pi / 180
	cos, sin := math.Cos(yaw), math.Sin(yaw)

	return r3.Vector{
		X: snapshot.origin.X + point.Forward*cos + point.Right*sin,
		Y: snapshot.origin.Y + point.Forward*sin - point.Right*cos,
		Z: snapshot.origin.Z + point.Up*snapshot.eyesHeight()/snapshot.template().EyesHeight,
	}
}

// viewBasis returns the forward, right and up vectors of the player's view.
func (snapshot *playerSnapshot) viewBasis() (r3.Vector, r3.Vector, r3.Vector) {
	yaw := snapshot.yaw * math.Pi / 180
	forward := csmath.ViewDirection(snapshot.yaw, snapshot.pitch)
	right := r3.Vector{X: math.Sin(yaw), Y: -math.Cos(yaw), Z: 0}
	up := right.Cross(forward)

	return forward, right, up
}

// isInFieldOfView returns true if the point is visible on the player's screen, ignoring obstacles.
func (snapshot *playerSnapshot) isInFieldOfView(point r3.Vector, forward r3.Vector, right r3.Vector, up r3.Vector) bool {
	direction := point.Sub(snapshot.eyes)
	depth := direction.Dot(forward)
	if depth <= 0 {
		return false
	}

	return math.Abs(direction.Dot(right)) <= depth*snapshot.fovTanHorizontal &&
		math.Abs(direction.Dot(up)) <= depth*snapshot.fovTanVertical
}

// segmentIntersectsSmoke returns true if the segment from a to b goes through the smoke detonated at the given
// position.
func segmentIntersectsSmoke(a r3.Vector, b r3.Vector, smokePosition r3.Vector) bool {
	// Scale the Z axis to turn the ellipsoid into a sphere.
	scale := smokeHorizontalRadius / smokeVerticalRadius
	center := r3.Vector{X: smokePosition.X, Y: smokePosition.Y, Z: (smokePosition.Z + smokeCenterHeight) * scale}
	start := r3.Vector{X: a.X, Y: a.Y, Z: a.Z * scale}
	end := r3.Vector{X: b.X, Y: b.Y, Z: b.Z * scale}

	segment := end.Sub(start)
	lengthSquared := segment.Norm2()
	t := 0.0
	if lengthSquared > 0 {
		t = math.Max(0, math.Min(1, center.Sub(start).Dot(segment)/lengthSquared))
	}
	closest := start.Add(segment.Mul(t))

	return closest.Sub(center).Norm2() <= smokeHorizontalRadius*smokeHorizontalRadius
}
