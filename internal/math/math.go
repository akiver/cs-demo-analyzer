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

func Max(a int, b int) int {
	if a > b {
		return a
	}

	return b
}
