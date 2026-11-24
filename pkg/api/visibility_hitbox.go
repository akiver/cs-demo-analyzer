package api

import (
	"math"

	csmath "github.com/akiver/cs-demo-analyzer/internal/math"
	"github.com/golang/geo/r3"
)

// Demos don't contain players' bones, the body of a player is approximated with points placed from its origin (feet),
// eyes height (lower when crouching) and yaw. A player is visible when any of these points can be seen.

type hitboxPoint struct {
	forward float64 // Units along the yaw direction
	right   float64 // Units to the right of the yaw direction
	height  float64 // Fraction of the eyes height
}

// Ordered from the most to the least likely visible point when peeking to exit early.
var hitboxPoints = [...]hitboxPoint{
	{0, 0, 1},      // Head (eyes)
	{2, 0, 0.72},   // Chest
	{0, 0, 1.08},   // Top of the head
	{0, 10, 0.82},  // Right shoulder
	{0, -10, 0.82}, // Left shoulder
	{26, 4, 0.82},  // Gun muzzle
	{12, 4, 0.78},  // Hands
	{0, 0, 0.88},   // Neck
	{0, 0, 0.6},    // Stomach
	{0, 5, 0.48},   // Right hip
	{0, -5, 0.48},  // Left hip
	{0, 5, 0.25},   // Right knee
	{0, -5, 0.25},  // Left knee
	{0, 6, 0.04},   // Right foot
	{0, -6, 0.04},  // Left foot
	{0, 0, 0.4},    // Pelvis
	{-4, 0, 0.7},   // Back
	{0, 12, 0.66},  // Right elbow
	{0, -12, 0.66}, // Left elbow
	{18, -2, 0.8},  // Left hand on the gun
	{0, 0, 0.14},   // Shins
	{6, 0, 0.94},   // Face
	{-3, 0, 0.98},  // Back of the head
	{0, 7, 0.55},   // Right flank
	{0, -7, 0.55},  // Left flank
	{3, 0, 0.32},   // Thighs
	{0, 0, 0.75},   // Upper chest
	{-2, 6, 0.36},  // Right buttock
	{-2, -6, 0.36}, // Left buttock
}

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
}

// fovTangents returns the horizontal and vertical tangents of half the field of view for a 4:3 fov in degrees.
func fovTangents(fovDegrees float64) (float64, float64) {
	tan := math.Tan(fovDegrees / 2 * math.Pi / 180)

	return tan * defaultFovTanHorizontal, tan * defaultFovTanVertical
}

func (snapshot *playerSnapshot) hitboxPoint(point hitboxPoint) r3.Vector {
	eyesHeight := snapshot.eyes.Z - snapshot.origin.Z
	if eyesHeight <= 0 {
		eyesHeight = defaultEyesHeight
	}

	yaw := snapshot.yaw * math.Pi / 180
	cos, sin := math.Cos(yaw), math.Sin(yaw)

	return r3.Vector{
		X: snapshot.origin.X + point.forward*cos + point.right*sin,
		Y: snapshot.origin.Y + point.forward*sin - point.right*cos,
		Z: snapshot.origin.Z + point.height*eyesHeight,
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
