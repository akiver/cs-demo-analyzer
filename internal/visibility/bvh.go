// Package visibility answers line of sight queries against a map collision mesh.
package visibility

import (
	"math"
	"slices"
)

const (
	maxLeafTriangles = 4
	sahBinCount      = 16
	// Cost of traversing a node relative to intersecting a triangle, used by the surface area heuristic (SAH).
	sahTraversalCost = 1.0
)

// Vec3 is a 3D vector stored with float32 precision to keep the mesh compact.
type Vec3 [3]float32

// Node is a BVH node.
// For a leaf (Count > 0) LeftOrFirst is the index of its first triangle, otherwise it's the index of its left child, the
// right child being stored right after it.
type Node struct {
	Min         Vec3
	Max         Vec3
	LeftOrFirst uint32
	Count       uint32
}

// BVH is a Bounding Volume Hierarchy over the triangles of a mesh.
// Triangles are stored in leaf order: triangle i is made of Vertices[Indices[i*3]], Vertices[Indices[i*3+1]] and
// Vertices[Indices[i*3+2]].
type BVH struct {
	Vertices []Vec3
	Indices  []uint32
	Nodes    []Node
}

type buildTriangle struct {
	min      Vec3
	max      Vec3
	centroid Vec3
	index    uint32
}

type aabb struct {
	min Vec3
	max Vec3
}

func emptyAABB() aabb {
	inf := float32(math.Inf(1))

	return aabb{min: Vec3{inf, inf, inf}, max: Vec3{-inf, -inf, -inf}}
}

func (box *aabb) grow(lower Vec3, upper Vec3) {
	for axis := range 3 {
		box.min[axis] = min(box.min[axis], lower[axis])
		box.max[axis] = max(box.max[axis], upper[axis])
	}
}

func (box *aabb) growPoint(point Vec3) {
	box.grow(point, point)
}

func (box aabb) area() float32 {
	dx := box.max[0] - box.min[0]
	dy := box.max[1] - box.min[1]
	dz := box.max[2] - box.min[2]
	if dx < 0 || dy < 0 || dz < 0 {
		return 0
	}

	return dx*dy + dy*dz + dz*dx
}

// BuildBVH builds a BVH over the given indexed triangles using a binned surface area heuristic.
// Degenerate triangles and triangles with non-finite vertices are dropped.
func BuildBVH(vertices []Vec3, indices []uint32) *BVH {
	triangles := make([]buildTriangle, 0, len(indices)/3)
	for i := 0; i+2 < len(indices); i += 3 {
		a, b, c := indices[i], indices[i+1], indices[i+2]
		if int(a) >= len(vertices) || int(b) >= len(vertices) || int(c) >= len(vertices) {
			continue
		}
		v0, v1, v2 := vertices[a], vertices[b], vertices[c]
		if !isFinite(v0) || !isFinite(v1) || !isFinite(v2) || isDegenerate(v0, v1, v2) {
			continue
		}

		box := emptyAABB()
		box.growPoint(v0)
		box.growPoint(v1)
		box.growPoint(v2)
		triangles = append(triangles, buildTriangle{
			min:   box.min,
			max:   box.max,
			index: uint32(i / 3),
			centroid: Vec3{
				(v0[0] + v1[0] + v2[0]) / 3,
				(v0[1] + v1[1] + v2[1]) / 3,
				(v0[2] + v1[2] + v2[2]) / 3,
			},
		})
	}

	bvh := &BVH{
		Vertices: vertices,
		Nodes:    make([]Node, 0, max(1, 2*len(triangles)/maxLeafTriangles)),
	}

	bvh.Nodes = append(bvh.Nodes, Node{})
	if len(triangles) > 0 {
		bvh.subdivide(0, triangles, 0)
	}

	bvh.Indices = make([]uint32, 0, len(triangles)*3)
	for _, triangle := range triangles {
		offset := triangle.index * 3
		bvh.Indices = append(bvh.Indices, indices[offset], indices[offset+1], indices[offset+2])
	}

	return bvh
}

func isFinite(vertex Vec3) bool {
	for _, value := range vertex {
		if math.IsNaN(float64(value)) || math.IsInf(float64(value), 0) {
			return false
		}
	}

	return true
}

func isDegenerate(v0 Vec3, v1 Vec3, v2 Vec3) bool {
	e1 := sub(v1, v0)
	e2 := sub(v2, v0)
	cross := Vec3{
		e1[1]*e2[2] - e1[2]*e2[1],
		e1[2]*e2[0] - e1[0]*e2[2],
		e1[0]*e2[1] - e1[1]*e2[0],
	}

	return cross[0] == 0 && cross[1] == 0 && cross[2] == 0
}

func sub(a Vec3, b Vec3) Vec3 {
	return Vec3{a[0] - b[0], a[1] - b[1], a[2] - b[2]}
}

// subdivide fills the node at nodeIndex with the given triangles, which start at firstTriangle in the final order.
func (bvh *BVH) subdivide(nodeIndex int, triangles []buildTriangle, firstTriangle int) {
	bounds := emptyAABB()
	centroidBounds := emptyAABB()
	for _, triangle := range triangles {
		bounds.grow(triangle.min, triangle.max)
		centroidBounds.growPoint(triangle.centroid)
	}

	makeLeaf := func() {
		bvh.Nodes[nodeIndex] = Node{
			Min:         bounds.min,
			Max:         bounds.max,
			LeftOrFirst: uint32(firstTriangle),
			Count:       uint32(len(triangles)),
		}
	}

	if len(triangles) <= maxLeafTriangles {
		makeLeaf()
		return
	}

	axis, splitPosition, splitCost, ok := findBestSplit(triangles, centroidBounds)
	leafCost := float32(len(triangles))
	if parentArea := bounds.area(); ok && parentArea > 0 {
		splitCost = sahTraversalCost + splitCost/parentArea
	}

	if ok && splitCost >= leafCost && len(triangles) <= maxLeafTriangles*4 {
		makeLeaf()
		return
	}

	leftCount := 0
	if ok {
		leftCount = partition(triangles, func(triangle buildTriangle) bool {
			return triangle.centroid[axis] < splitPosition
		})
	}

	if leftCount == 0 || leftCount == len(triangles) {
		// The heuristic could not split the triangles (e.g. identical centroids), fallback to a median split.
		axis = longestAxis(centroidBounds)
		slices.SortFunc(triangles, func(a, b buildTriangle) int {
			if a.centroid[axis] < b.centroid[axis] {
				return -1
			}
			if a.centroid[axis] > b.centroid[axis] {
				return 1
			}
			return 0
		})
		leftCount = len(triangles) / 2
	}

	leftIndex := len(bvh.Nodes)
	bvh.Nodes = append(bvh.Nodes, Node{}, Node{})
	bvh.Nodes[nodeIndex] = Node{
		Min:         bounds.min,
		Max:         bounds.max,
		LeftOrFirst: uint32(leftIndex),
		Count:       0,
	}

	bvh.subdivide(leftIndex, triangles[:leftCount], firstTriangle)
	bvh.subdivide(leftIndex+1, triangles[leftCount:], firstTriangle+leftCount)
}

func longestAxis(box aabb) int {
	extent := sub(box.max, box.min)
	axis := 0
	if extent[1] > extent[axis] {
		axis = 1
	}
	if extent[2] > extent[axis] {
		axis = 2
	}

	return axis
}

func partition(triangles []buildTriangle, isLeft func(buildTriangle) bool) int {
	left := 0
	for i := range triangles {
		if isLeft(triangles[i]) {
			triangles[left], triangles[i] = triangles[i], triangles[left]
			left++
		}
	}

	return left
}

type sahBin struct {
	bounds aabb
	count  int
}

// findBestSplit returns the axis and position minimizing the SAH cost (sum of child area * triangle count).
func findBestSplit(triangles []buildTriangle, centroidBounds aabb) (int, float32, float32, bool) {
	bestCost := float32(math.Inf(1))
	bestAxis := -1
	var bestPosition float32

	for axis := range 3 {
		minBound := centroidBounds.min[axis]
		maxBound := centroidBounds.max[axis]
		if maxBound <= minBound {
			continue
		}

		var bins [sahBinCount]sahBin
		for i := range bins {
			bins[i].bounds = emptyAABB()
		}

		scale := float32(sahBinCount) / (maxBound - minBound)
		for _, triangle := range triangles {
			binIndex := min(sahBinCount-1, int((triangle.centroid[axis]-minBound)*scale))
			bins[binIndex].count++
			bins[binIndex].bounds.grow(triangle.min, triangle.max)
		}

		var leftArea, rightArea [sahBinCount - 1]float32
		var leftCount, rightCount [sahBinCount - 1]int
		leftBox, rightBox := emptyAABB(), emptyAABB()
		leftSum, rightSum := 0, 0
		for i := range sahBinCount - 1 {
			leftSum += bins[i].count
			leftCount[i] = leftSum
			leftBox.grow(bins[i].bounds.min, bins[i].bounds.max)
			leftArea[i] = leftBox.area()

			rightSum += bins[sahBinCount-1-i].count
			rightCount[sahBinCount-2-i] = rightSum
			rightBox.grow(bins[sahBinCount-1-i].bounds.min, bins[sahBinCount-1-i].bounds.max)
			rightArea[sahBinCount-2-i] = rightBox.area()
		}

		binWidth := (maxBound - minBound) / float32(sahBinCount)
		for i := range sahBinCount - 1 {
			if leftCount[i] == 0 || rightCount[i] == 0 {
				continue
			}
			cost := float32(leftCount[i])*leftArea[i] + float32(rightCount[i])*rightArea[i]
			if cost < bestCost {
				bestCost = cost
				bestAxis = axis
				bestPosition = minBound + binWidth*float32(i+1)
			}
		}
	}

	return bestAxis, bestPosition, bestCost, bestAxis != -1
}

// TriangleCount returns the number of triangles in the BVH.
func (bvh *BVH) TriangleCount() int {
	return len(bvh.Indices) / 3
}
