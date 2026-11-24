package visibility

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"math/rand/v2"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"testing"

	"github.com/golang/geo/r3"
)

// box returns the 12 triangles of an axis aligned box.
func box(vertices []Vec3, indices []uint32, lower Vec3, upper Vec3) ([]Vec3, []uint32) {
	base := uint32(len(vertices))
	for i := range 8 {
		vertex := lower
		if i&1 != 0 {
			vertex[0] = upper[0]
		}
		if i&2 != 0 {
			vertex[1] = upper[1]
		}
		if i&4 != 0 {
			vertex[2] = upper[2]
		}
		vertices = append(vertices, vertex)
	}

	faces := [][4]uint32{
		{0, 1, 3, 2}, {4, 5, 7, 6}, // z min, z max
		{0, 1, 5, 4}, {2, 3, 7, 6}, // y min, y max
		{0, 2, 6, 4}, {1, 3, 7, 5}, // x min, x max
	}
	for _, face := range faces {
		indices = append(indices,
			base+face[0], base+face[1], base+face[2],
			base+face[0], base+face[2], base+face[3],
		)
	}

	return vertices, indices
}

func newTestScene() *BVH {
	var vertices []Vec3
	var indices []uint32
	// A ground plane, a wall and a crate.
	vertices, indices = box(vertices, indices, Vec3{-1000, -1000, -10}, Vec3{1000, 1000, 0})
	vertices, indices = box(vertices, indices, Vec3{-10, -200, 0}, Vec3{10, 200, 128})
	vertices, indices = box(vertices, indices, Vec3{300, 300, 0}, Vec3{350, 350, 50})

	return BuildBVH(vertices, indices)
}

func TestIsOccluded(t *testing.T) {
	tracer := NewTracer(newTestScene())

	tests := []struct {
		name     string
		a        r3.Vector
		b        r3.Vector
		occluded bool
	}{
		{"through the wall", r3.Vector{X: -100, Y: 0, Z: 64}, r3.Vector{X: 100, Y: 0, Z: 64}, true},
		{"over the wall", r3.Vector{X: -100, Y: 0, Z: 200}, r3.Vector{X: 100, Y: 0, Z: 200}, false},
		{"around the wall", r3.Vector{X: -100, Y: 300, Z: 64}, r3.Vector{X: 100, Y: 300, Z: 64}, false},
		{"parallel to the wall", r3.Vector{X: -50, Y: -500, Z: 64}, r3.Vector{X: -50, Y: 500, Z: 64}, false},
		{"crate in the way", r3.Vector{X: 250, Y: 325, Z: 20}, r3.Vector{X: 400, Y: 325, Z: 20}, true},
		{"over the crate", r3.Vector{X: 250, Y: 325, Z: 60}, r3.Vector{X: 400, Y: 325, Z: 60}, false},
		{"into the ground", r3.Vector{X: 500, Y: 0, Z: 64}, r3.Vector{X: 500, Y: 0, Z: -64}, true},
		{"ends on the ground", r3.Vector{X: 500, Y: 0, Z: 64}, r3.Vector{X: 600, Y: 0, Z: 0}, false},
		{"vertical segment", r3.Vector{X: 500, Y: 500, Z: 64}, r3.Vector{X: 500, Y: 500, Z: 10}, false},
		{"zero length", r3.Vector{X: 500, Y: 500, Z: 64}, r3.Vector{X: 500, Y: 500, Z: 64}, false},
		{"ends before the wall", r3.Vector{X: -100, Y: 0, Z: 64}, r3.Vector{X: -11, Y: 0, Z: 64}, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := tracer.IsOccluded(test.a, test.b); got != test.occluded {
				t.Errorf("IsOccluded(%v, %v) = %v, want %v", test.a, test.b, got, test.occluded)
			}
			if got := tracer.IsOccluded(test.b, test.a); got != test.occluded {
				t.Errorf("IsOccluded(%v, %v) = %v, want %v (reversed)", test.b, test.a, got, test.occluded)
			}
		})
	}
}

func bruteForceIsOccluded(bvh *BVH, a r3.Vector, b r3.Vector) bool {
	tracer := NewTracer(bvh)
	origin := [3]float64{a.X, a.Y, a.Z}
	direction := [3]float64{b.X - a.X, b.Y - a.Y, b.Z - a.Z}
	for triangle := range uint32(bvh.TriangleCount()) {
		if tracer.intersectsTriangle(triangle, origin, direction) {
			return true
		}
	}

	return false
}

func TestBVHMatchesBruteForce(t *testing.T) {
	random := rand.New(rand.NewPCG(1, 2))
	randomPoint := func() Vec3 {
		return Vec3{random.Float32()*2000 - 1000, random.Float32()*2000 - 1000, random.Float32() * 500}
	}

	var vertices []Vec3
	var indices []uint32
	for range 5000 {
		center := randomPoint()
		base := uint32(len(vertices))
		for range 3 {
			vertices = append(vertices, Vec3{
				center[0] + random.Float32()*60 - 30,
				center[1] + random.Float32()*60 - 30,
				center[2] + random.Float32()*60 - 30,
			})
		}
		indices = append(indices, base, base+1, base+2)
	}

	bvh := BuildBVH(vertices, indices)
	if bvh.TriangleCount() != 5000 {
		t.Fatalf("expected 5000 triangles, got %d", bvh.TriangleCount())
	}

	tracer := NewTracer(bvh)
	occludedCount := 0
	for range 3000 {
		a, b := randomPoint(), randomPoint()
		from := r3.Vector{X: float64(a[0]), Y: float64(a[1]), Z: float64(a[2])}
		to := r3.Vector{X: float64(b[0]), Y: float64(b[1]), Z: float64(b[2])}
		want := bruteForceIsOccluded(bvh, from, to)
		if got := tracer.IsOccluded(from, to); got != want {
			t.Fatalf("IsOccluded(%v, %v) = %v, brute force = %v", from, to, got, want)
		}
		if want {
			occludedCount++
		}
	}

	if occludedCount == 0 || occludedCount == 3000 {
		t.Fatalf("test scene is not meaningful, %d occluded segments", occludedCount)
	}
}

func TestBuildBVHDropsInvalidTriangles(t *testing.T) {
	vertices := []Vec3{{0, 0, 0}, {1, 0, 0}, {0, 1, 0}, {2, 0, 0}}
	indices := []uint32{
		0, 1, 2, // valid
		0, 1, 3, // degenerate (collinear)
		0, 0, 2, // degenerate (repeated vertex)
		0, 1, 9, // out of range
	}

	bvh := BuildBVH(vertices, indices)
	if bvh.TriangleCount() != 1 {
		t.Fatalf("expected 1 triangle, got %d", bvh.TriangleCount())
	}
}

func TestEmptyBVH(t *testing.T) {
	tracer := NewTracer(BuildBVH(nil, nil))
	if tracer.IsOccluded(r3.Vector{}, r3.Vector{X: 1}) {
		t.Fatal("empty BVH must not occlude")
	}
}

func TestGeometryFileRoundTrip(t *testing.T) {
	bvh := newTestScene()

	var buffer bytes.Buffer
	if err := WriteGeometry(&buffer, bvh); err != nil {
		t.Fatal(err)
	}

	got, err := ReadGeometry(&buffer)
	if err != nil {
		t.Fatal(err)
	}

	if len(got.Vertices) != len(bvh.Vertices) || len(got.Indices) != len(bvh.Indices) || len(got.Nodes) != len(bvh.Nodes) {
		t.Fatalf("round trip mismatch: %d/%d/%d vs %d/%d/%d", len(got.Vertices), len(got.Indices), len(got.Nodes), len(bvh.Vertices), len(bvh.Indices), len(bvh.Nodes))
	}
	for i := range bvh.Nodes {
		if got.Nodes[i] != bvh.Nodes[i] {
			t.Fatalf("node %d mismatch", i)
		}
	}

	tracer := NewTracer(got)
	if !tracer.IsOccluded(r3.Vector{X: -100, Y: 0, Z: 64}, r3.Vector{X: 100, Y: 0, Z: 64}) {
		t.Fatal("expected the wall to occlude after a round trip")
	}
}

func TestLoadGeometryFile(t *testing.T) {
	bvh := newTestScene()
	var buffer bytes.Buffer
	if err := WriteGeometry(&buffer, bvh); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "de_test"+GeometryFileExtension)
	if err := os.WriteFile(path, buffer.Bytes(), 0o644); err != nil {
		t.Fatal(err)
	}

	got, checksum, err := LoadGeometryFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if got.TriangleCount() != bvh.TriangleCount() {
		t.Fatalf("expected %d triangles, got %d", bvh.TriangleCount(), got.TriangleCount())
	}
	// The checksum is the one of the file, like in the manifest.
	want := sha256.Sum256(buffer.Bytes())
	if checksum != hex.EncodeToString(want[:]) {
		t.Fatalf("expected checksum %x, got %s", want, checksum)
	}

	if _, _, err := LoadGeometryFile(filepath.Join(t.TempDir(), "missing"+GeometryFileExtension)); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("expected os.ErrNotExist, got %v", err)
	}
}

// Sections bigger than a chunk are decoded by several reads.
func TestGeometryFileRoundTripMultipleChunks(t *testing.T) {
	var vertices []Vec3
	var indices []uint32
	for i := range 2 * geometryChunkElements {
		x := float32(i * 10)
		base := uint32(len(vertices))
		vertices = append(vertices, Vec3{x, 0, 0}, Vec3{x + 5, 0, 0}, Vec3{x, 5, float32(i % 7)})
		indices = append(indices, base, base+1, base+2)
	}
	bvh := BuildBVH(vertices, indices)

	var buffer bytes.Buffer
	if err := WriteGeometry(&buffer, bvh); err != nil {
		t.Fatal(err)
	}

	got, err := ReadGeometry(&buffer)
	if err != nil {
		t.Fatal(err)
	}

	if !slices.Equal(got.Vertices, bvh.Vertices) {
		t.Fatal("vertices mismatch")
	}
	if !slices.Equal(got.Indices, bvh.Indices) {
		t.Fatal("indices mismatch")
	}
	if !slices.Equal(got.Nodes, bvh.Nodes) {
		t.Fatal("nodes mismatch")
	}
}

func TestReadGeometryRejectsInvalidData(t *testing.T) {
	if _, err := ReadGeometry(bytes.NewReader([]byte("not gzip"))); err == nil {
		t.Fatal("expected an error")
	}
}

// newGeometryHeader returns a geometry file made of a header only.
func newGeometryHeader(t *testing.T, vertexCount uint32, indexCount uint32, nodeCount uint32) *bytes.Buffer {
	t.Helper()

	header := geometryFileMagic[:]
	for _, value := range []uint32{GeometryFileVersion, vertexCount, indexCount, nodeCount} {
		header = binary.LittleEndian.AppendUint32(header, value)
	}

	var buffer bytes.Buffer
	gzipWriter := gzip.NewWriter(&buffer)
	if _, err := gzipWriter.Write(header); err != nil {
		t.Fatal(err)
	}
	if err := gzipWriter.Close(); err != nil {
		t.Fatal(err)
	}

	return &buffer
}

func TestReadGeometryRejectsTooLargeCounts(t *testing.T) {
	tooLarge := uint32(maxGeometryElements + 3)
	for _, counts := range [][3]uint32{{tooLarge, 0, 0}, {0, tooLarge, 0}, {0, 0, tooLarge}} {
		_, err := ReadGeometry(newGeometryHeader(t, counts[0], counts[1], counts[2]))
		if !errors.Is(err, ErrInvalidGeometryFile) {
			t.Fatalf("counts %v: expected ErrInvalidGeometryFile, got %v", counts, err)
		}
	}
}

// A header declaring the maximum counts without the data must not allocate the memory of the declared sections.
func TestReadGeometryTruncatedFileAllocations(t *testing.T) {
	maxIndexCount := uint32(maxGeometryElements - maxGeometryElements%3)
	file := newGeometryHeader(t, maxGeometryElements, maxIndexCount, maxGeometryElements)

	var before, after runtime.MemStats
	runtime.ReadMemStats(&before)
	_, err := ReadGeometry(file)
	runtime.ReadMemStats(&after)

	if !errors.Is(err, ErrInvalidGeometryFile) {
		t.Fatalf("expected ErrInvalidGeometryFile, got %v", err)
	}

	const maxAllocatedBytes = 64 << 20
	if allocated := after.TotalAlloc - before.TotalAlloc; allocated > maxAllocatedBytes {
		t.Fatalf("allocated %d MB reading a truncated file, want at most %d MB", allocated>>20, maxAllocatedBytes>>20)
	}
}
