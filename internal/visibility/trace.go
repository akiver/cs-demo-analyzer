package visibility

import (
	"math"

	"github.com/golang/geo/r3"
)

const (
	// Intersections closer than this fraction of the segment length to its ends are ignored, so that a point lying
	// exactly on a surface (e.g. a foot on the ground) is not considered occluded by it.
	segmentEndEpsilon = 1e-4
)

// Tracer runs segment queries against a BVH.
// It's not safe for concurrent use, create one Tracer per goroutine.
type Tracer struct {
	bvh   *BVH
	stack []uint32
}

// NewTracer returns a tracer for the given BVH.
func NewTracer(bvh *BVH) *Tracer {
	return &Tracer{bvh: bvh, stack: make([]uint32, 0, 64)}
}

// IsOccluded returns true if any triangle intersects the segment from a to b.
func (tracer *Tracer) IsOccluded(a r3.Vector, b r3.Vector) bool {
	bvh := tracer.bvh
	if len(bvh.Nodes) == 0 || len(bvh.Indices) == 0 {
		return false
	}

	origin := [3]float64{a.X, a.Y, a.Z}
	direction := [3]float64{b.X - a.X, b.Y - a.Y, b.Z - a.Z}
	var inverseDirection [3]float64
	for axis := range 3 {
		inverseDirection[axis] = 1 / direction[axis] // ±Inf when the direction is 0 on this axis, handled by the slab test
	}

	stack := tracer.stack[:0]
	defer func() { tracer.stack = stack[:0] }()

	nodeIndex := uint32(0)
	for {
		node := &bvh.Nodes[nodeIndex]
		if intersectsAABB(node, origin, inverseDirection) {
			if node.Count == 0 {
				stack = append(stack, node.LeftOrFirst+1)
				nodeIndex = node.LeftOrFirst
				continue
			}

			first := node.LeftOrFirst
			for triangle := first; triangle < first+node.Count; triangle++ {
				if tracer.intersectsTriangle(triangle, origin, direction) {
					return true
				}
			}
		}

		if len(stack) == 0 {
			return false
		}
		nodeIndex = stack[len(stack)-1]
		stack = stack[:len(stack)-1]
	}
}

// intersectsAABB returns true if the segment origin + t * direction, t in [0, 1], intersects the node bounds.
func intersectsAABB(node *Node, origin [3]float64, inverseDirection [3]float64) bool {
	tMin := 0.0
	tMax := 1.0
	for axis := range 3 {
		t1 := (float64(node.Min[axis]) - origin[axis]) * inverseDirection[axis]
		t2 := (float64(node.Max[axis]) - origin[axis]) * inverseDirection[axis]
		if math.IsNaN(t1) || math.IsNaN(t2) {
			// The segment is parallel to the slab and starts on one of its planes.
			if origin[axis] < float64(node.Min[axis]) || origin[axis] > float64(node.Max[axis]) {
				return false
			}
			continue
		}
		if t1 > t2 {
			t1, t2 = t2, t1
		}
		tMin = max(tMin, t1)
		tMax = min(tMax, t2)
		if tMin > tMax {
			return false
		}
	}

	return true
}

// intersectsTriangle runs the Möller–Trumbore algorithm, culling disabled.
func (tracer *Tracer) intersectsTriangle(triangle uint32, origin [3]float64, direction [3]float64) bool {
	bvh := tracer.bvh
	offset := triangle * 3
	v0 := bvh.Vertices[bvh.Indices[offset]]
	v1 := bvh.Vertices[bvh.Indices[offset+1]]
	v2 := bvh.Vertices[bvh.Indices[offset+2]]

	e1 := [3]float64{float64(v1[0] - v0[0]), float64(v1[1] - v0[1]), float64(v1[2] - v0[2])}
	e2 := [3]float64{float64(v2[0] - v0[0]), float64(v2[1] - v0[1]), float64(v2[2] - v0[2])}

	p := cross(direction, e2)
	determinant := dot(e1, p)
	if math.Abs(determinant) < 1e-12 {
		return false
	}
	inverseDeterminant := 1 / determinant

	s := [3]float64{origin[0] - float64(v0[0]), origin[1] - float64(v0[1]), origin[2] - float64(v0[2])}
	u := dot(s, p) * inverseDeterminant
	if u < 0 || u > 1 {
		return false
	}

	q := cross(s, e1)
	v := dot(direction, q) * inverseDeterminant
	if v < 0 || u+v > 1 {
		return false
	}

	t := dot(e2, q) * inverseDeterminant

	return t > segmentEndEpsilon && t < 1-segmentEndEpsilon
}

func cross(a [3]float64, b [3]float64) [3]float64 {
	return [3]float64{
		a[1]*b[2] - a[2]*b[1],
		a[2]*b[0] - a[0]*b[2],
		a[0]*b[1] - a[1]*b[0],
	}
}

func dot(a [3]float64, b [3]float64) float64 {
	return a[0]*b[0] + a[1]*b[1] + a[2]*b[2]
}
