package math

import (
	"math"
	"testing"

	"github.com/golang/geo/r3"
)

func TestViewDirection(t *testing.T) {
	tests := []struct {
		name  string
		yaw   float64
		pitch float64
		want  r3.Vector
	}{
		{"east", 0, 0, r3.Vector{X: 1}},
		{"north", 90, 0, r3.Vector{Y: 1}},
		{"west", 180, 0, r3.Vector{X: -1}},
		{"south", -90, 0, r3.Vector{Y: -1}},
		{"down", 0, 90, r3.Vector{Z: -1}},
		{"up", 0, -90, r3.Vector{Z: 1}},
		{"up as stored in demos", 0, 270, r3.Vector{Z: 1}},
		{"45 degrees down", 0, 45, r3.Vector{X: math.Sqrt2 / 2, Z: -math.Sqrt2 / 2}},
	}
	for _, test := range tests {
		got := ViewDirection(test.yaw, test.pitch)
		if got.Sub(test.want).Norm() > 1e-9 {
			t.Errorf("%s: ViewDirection(%v, %v) = %v, want %v", test.name, test.yaw, test.pitch, got, test.want)
		}
	}
}

func TestAverage(t *testing.T) {
	tests := []struct {
		values []int
		want   float32
	}{
		{nil, 0},
		{[]int{100}, 100},
		{[]int{100, 200}, 150},
		{[]int{100, 200, 400}, float32(700) / 3},
	}
	for _, test := range tests {
		if got := Average(test.values); got != test.want {
			t.Errorf("Average(%v) = %v, want %v", test.values, got, test.want)
		}
	}
}
