package source2

import (
	"slices"
	"testing"
)

// newTestHull returns a hull made of a single triangular face, edges are (next, twin, origin, face) half edges.
func newTestHull(edges []byte, faces []byte) map[string]any {
	return map[string]any{
		// 3 vertices, their position doesn't matter.
		"m_VertexPositions": make([]byte, 3*12),
		"m_Edges":           edges,
		"m_Faces":           faces,
	}
}

func TestAddHull(t *testing.T) {
	mesh := &CollisionMesh{}
	hull := newTestHull([]byte{
		1, 0, 0, 0,
		2, 0, 1, 0,
		0, 0, 2, 0,
	}, []byte{0})

	if err := mesh.addHull(hull); err != nil {
		t.Fatal(err)
	}
	if len(mesh.Vertices) != 3 {
		t.Fatalf("expected 3 vertices, got %d", len(mesh.Vertices))
	}
	if want := []uint32{0, 1, 2}; !slices.Equal(mesh.Indices, want) {
		t.Fatalf("expected indices %v, got %v", want, mesh.Indices)
	}
}

func TestAddHullRejectsOutOfRangeIndices(t *testing.T) {
	tests := []struct {
		name  string
		edges []byte
		faces []byte
	}{
		{"face edge", []byte{1, 0, 0, 0, 2, 0, 1, 0, 0, 0, 2, 0}, []byte{3}},
		{"next of the first edge", []byte{200, 0, 0, 0, 2, 0, 1, 0, 0, 0, 2, 0}, []byte{0}},
		{"next of a following edge", []byte{1, 0, 0, 0, 200, 0, 1, 0, 0, 0, 2, 0}, []byte{0}},
		{"vertex", []byte{1, 0, 0, 0, 2, 0, 9, 0, 0, 0, 2, 0}, []byte{0}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mesh := &CollisionMesh{}
			if err := mesh.addHull(newTestHull(test.edges, test.faces)); err == nil {
				t.Fatal("expected an error")
			}
		})
	}
}
