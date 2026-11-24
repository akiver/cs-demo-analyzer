package math

import (
	"math"

	"github.com/golang/geo/r3"
)

// Return the distance between two vectors in meters.
// https://developer.valvesoftware.com/wiki/Dimensions
// 1 unit = 0.75 inch = 19.05mm = 0.01905m
func GetDistanceBetweenVectors(vectorA r3.Vector, vectorB r3.Vector) float64 {
	return math.Sqrt(math.Pow(vectorA.X-vectorB.X, 2)+math.Pow(vectorA.Y-vectorB.Y, 2)+math.Pow(vectorA.Z-vectorB.Z, 2)) * 0.01905
}

// ViewDirection returns the unit vector of a view from its yaw and pitch in degrees, the pitch being positive when
// looking down like in the game.
func ViewDirection(yawDegrees float64, pitchDegrees float64) r3.Vector {
	yaw := yawDegrees * math.Pi / 180
	pitch := pitchDegrees * math.Pi / 180
	cosPitch := math.Cos(pitch)

	return r3.Vector{
		X: cosPitch * math.Cos(yaw),
		Y: cosPitch * math.Sin(yaw),
		Z: -math.Sin(pitch),
	}
}

func Max(a int, b int) int {
	if a > b {
		return a
	}

	return b
}

// Average returns the arithmetic mean of values, 0 when empty.
func Average(values []int) float32 {
	if len(values) == 0 {
		return 0
	}

	sum := 0
	for _, value := range values {
		sum += value
	}

	return float32(sum) / float32(len(values))
}
