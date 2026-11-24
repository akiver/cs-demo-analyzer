package source2

import (
	"encoding/binary"
	"fmt"
	"math"
	"slices"
	"strings"
)

// CollisionAttribute describes how a physics shape interacts with the world.
type CollisionAttribute struct {
	CollisionGroup  string
	InteractAs      []string
	InteractExclude []string
}

// CollisionMesh is an indexed triangle mesh.
type CollisionMesh struct {
	Vertices [][3]float32
	Indices  []uint32
}

// CollisionStats counts the extracted and skipped shapes.
type CollisionStats struct {
	Hulls            int
	MeshTriangles    int
	Capsules         int
	Spheres          int
	SkippedShapes    int
	SkippedTriangles int
	Compounds        int
}

const (
	capsuleSegments = 8
	sphereRings     = 4
)

// IsBlockingLineOfSight returns true if the collision attribute is for geometry that blocks the view of players.
// It keeps the default collision group, including its variants that players or grenades don't collide with (e.g.
// detailed stairs smoothed by a clip), and explicit line of sight blockers.
// It drops clips, ladders, sky, sound blockers, nav clips and see-through geometry such as windows and fences
// (tagged passbullets).
func IsBlockingLineOfSight(attribute CollisionAttribute) bool {
	if slices.Contains(attribute.InteractAs, "blocklos") {
		return true
	}

	if strings.EqualFold(attribute.CollisionGroup, "default") {
		for _, tag := range attribute.InteractAs {
			switch tag {
			case "passbullets", "window", "sky", "playerclip", "npcclip", "csgo_grenadeclip", "ladder":
				return false
			}
		}

		return true
	}

	return slices.Contains(attribute.InteractAs, "solid") && !slices.Contains(attribute.InteractAs, "window")
}

// ExtractCollisionMesh returns the triangles of the shapes of a physics aggregate (the decoded PHYS block of a
// world_physics.vmdl_c file) whose collision attribute passes the filter.
func ExtractCollisionMesh(physics any, filter func(CollisionAttribute) bool) (*CollisionMesh, CollisionStats, error) {
	var stats CollisionStats
	root, ok := physics.(map[string]any)
	if !ok {
		return nil, stats, fmt.Errorf("physics data is not an object")
	}

	attributes := parseCollisionAttributes(root)
	isIncluded := func(shape map[string]any) bool {
		index := toInt(shape["m_nCollisionAttributeIndex"])
		if index < 0 || index >= len(attributes) {
			return false
		}

		return filter(attributes[index])
	}

	mesh := &CollisionMesh{}
	parts, _ := root["m_parts"].([]any)
	for _, part := range parts {
		partObject, ok := part.(map[string]any)
		if !ok {
			continue
		}
		shape, ok := partObject["m_rnShape"].(map[string]any)
		if !ok {
			continue
		}

		if compounds, ok := shape["m_compounds"].([]any); ok {
			stats.Compounds += len(compounds)
		}

		for _, item := range asObjects(shape["m_hulls"]) {
			hull, ok := item["m_Hull"].(map[string]any)
			if !ok {
				continue
			}
			if !isIncluded(item) {
				stats.SkippedShapes++
				continue
			}
			if err := mesh.addHull(hull); err != nil {
				return nil, stats, fmt.Errorf("hull: %w", err)
			}
			stats.Hulls++
		}

		for _, item := range asObjects(shape["m_meshes"]) {
			rnMesh, ok := item["m_Mesh"].(map[string]any)
			if !ok {
				continue
			}
			triangles := len(asBytes(rnMesh["m_Triangles"])) / 12
			if !isIncluded(item) {
				stats.SkippedShapes++
				stats.SkippedTriangles += triangles
				continue
			}
			if err := mesh.addMesh(rnMesh); err != nil {
				return nil, stats, fmt.Errorf("mesh: %w", err)
			}
			stats.MeshTriangles += triangles
		}

		for _, item := range asObjects(shape["m_capsules"]) {
			capsule, ok := item["m_Capsule"].(map[string]any)
			if !ok {
				continue
			}
			if !isIncluded(item) {
				stats.SkippedShapes++
				continue
			}
			centers, _ := capsule["m_vCenter"].([]any)
			if len(centers) != 2 {
				return nil, stats, fmt.Errorf("capsule must have 2 centers")
			}
			mesh.addCapsule(toVec3(centers[0]), toVec3(centers[1]), float32(toFloat(capsule["m_flRadius"])))
			stats.Capsules++
		}

		for _, item := range asObjects(shape["m_spheres"]) {
			sphere, ok := item["m_Sphere"].(map[string]any)
			if !ok {
				continue
			}
			if !isIncluded(item) {
				stats.SkippedShapes++
				continue
			}
			center := toVec3(sphere["m_vCenter"])
			mesh.addCapsule(center, center, float32(toFloat(sphere["m_flRadius"])))
			stats.Spheres++
		}
	}

	return mesh, stats, nil
}

func parseCollisionAttributes(root map[string]any) []CollisionAttribute {
	var attributes []CollisionAttribute
	for _, item := range asObjects(root["m_collisionAttributes"]) {
		attribute := CollisionAttribute{
			CollisionGroup:  toString(item["m_CollisionGroupString"]),
			InteractAs:      toLowerStrings(item["m_InteractAsStrings"]),
			InteractExclude: toLowerStrings(item["m_InteractExcludeStrings"]),
		}
		if len(attribute.InteractAs) == 0 {
			// Older assets
			attribute.InteractAs = toLowerStrings(item["m_PhysicsTagStrings"])
		}
		attributes = append(attributes, attribute)
	}

	return attributes
}

func readVertices(blob []byte) ([][3]float32, error) {
	if len(blob)%12 != 0 {
		return nil, fmt.Errorf("invalid vertices blob size %d", len(blob))
	}

	vertices := make([][3]float32, len(blob)/12)
	for i := range vertices {
		for axis := range 3 {
			vertices[i][axis] = math.Float32frombits(binary.LittleEndian.Uint32(blob[i*12+axis*4:]))
		}
	}

	return vertices, nil
}

func (mesh *CollisionMesh) addHull(hull map[string]any) error {
	verticesKey := "m_Vertices"
	if _, ok := hull["m_VertexPositions"]; ok {
		// Since 2023-11, m_Vertices holds per vertex data and positions are stored separately.
		verticesKey = "m_VertexPositions"
	}

	vertices, err := readVertices(asBytes(hull[verticesKey]))
	if err != nil {
		return err
	}

	// Half edge: next, twin, origin, face (1 byte each). Face: edge (1 byte).
	edges := asBytes(hull["m_Edges"])
	faces := asBytes(hull["m_Faces"])
	if len(edges)%4 != 0 {
		return fmt.Errorf("invalid edges blob size %d", len(edges))
	}
	edgeCount := len(edges) / 4
	next := func(edge int) int { return int(edges[edge*4]) }
	origin := func(edge int) int { return int(edges[edge*4+2]) }

	base := uint32(len(mesh.Vertices))
	mesh.Vertices = append(mesh.Vertices, vertices...)
	for _, faceEdge := range faces {
		start := int(faceEdge)
		if start >= edgeCount {
			return fmt.Errorf("face edge %d out of range", start)
		}

		// Fan triangulation of the convex face.
		edge := next(start)
		if edge >= edgeCount {
			return fmt.Errorf("edge %d out of range", edge)
		}
		for steps := 0; edge != start && steps < edgeCount; steps++ {
			nextEdge := next(edge)
			if nextEdge >= edgeCount {
				return fmt.Errorf("edge %d out of range", nextEdge)
			}
			if nextEdge == start {
				break
			}
			a, b, c := origin(start), origin(edge), origin(nextEdge)
			if a >= len(vertices) || b >= len(vertices) || c >= len(vertices) {
				return fmt.Errorf("vertex out of range")
			}
			mesh.Indices = append(mesh.Indices, base+uint32(a), base+uint32(b), base+uint32(c))
			edge = nextEdge
		}
	}

	return nil
}

func (mesh *CollisionMesh) addMesh(rnMesh map[string]any) error {
	vertices, err := readVertices(asBytes(rnMesh["m_Vertices"]))
	if err != nil {
		return err
	}

	triangles := asBytes(rnMesh["m_Triangles"])
	if len(triangles)%12 != 0 {
		return fmt.Errorf("invalid triangles blob size %d", len(triangles))
	}

	base := uint32(len(mesh.Vertices))
	mesh.Vertices = append(mesh.Vertices, vertices...)
	for i := 0; i < len(triangles); i += 4 {
		index := binary.LittleEndian.Uint32(triangles[i:])
		if int(index) >= len(vertices) {
			return fmt.Errorf("triangle vertex %d out of range", index)
		}
		mesh.Indices = append(mesh.Indices, base+index)
	}

	return nil
}

// addCapsule tessellates a capsule (a sphere when both centers are equal) as a low poly convex hull.
func (mesh *CollisionMesh) addCapsule(a [3]float32, b [3]float32, radius float32) {
	if radius <= 0 {
		return
	}

	axis := [3]float64{float64(b[0] - a[0]), float64(b[1] - a[1]), float64(b[2] - a[2])}
	length := math.Sqrt(axis[0]*axis[0] + axis[1]*axis[1] + axis[2]*axis[2])
	if length < 1e-6 {
		axis = [3]float64{0, 0, 1}
	} else {
		axis = [3]float64{axis[0] / length, axis[1] / length, axis[2] / length}
	}

	// Orthonormal basis around the capsule axis.
	reference := [3]float64{1, 0, 0}
	if math.Abs(axis[0]) > 0.9 {
		reference = [3]float64{0, 1, 0}
	}
	u := normalize(crossProduct(axis, reference))
	v := crossProduct(axis, u)

	r := float64(radius)
	// Rings from the bottom pole (around a) to the top pole (around b): each hemisphere has sphereRings rings.
	type ring struct {
		center [3]float32
		radius float64
		offset float64
	}
	var rings []ring
	for i := 1; i <= sphereRings; i++ {
		angle := -math.Pi/2 + math.Pi/2*float64(i)/float64(sphereRings)
		rings = append(rings, ring{center: a, radius: r * math.Cos(angle), offset: r * math.Sin(angle)})
	}
	for i := 0; i < sphereRings; i++ {
		angle := math.Pi / 2 * float64(i) / float64(sphereRings)
		rings = append(rings, ring{center: b, radius: r * math.Cos(angle), offset: r * math.Sin(angle)})
	}

	point := func(center [3]float32, offset float64, radius float64, theta float64) [3]float32 {
		cos, sin := math.Cos(theta), math.Sin(theta)
		var result [3]float32
		for i := range 3 {
			result[i] = center[i] + float32(axis[i]*offset+(u[i]*cos+v[i]*sin)*radius)
		}

		return result
	}

	base := uint32(len(mesh.Vertices))
	mesh.Vertices = append(mesh.Vertices, point(a, -r, 0, 0)) // bottom pole
	for _, ring := range rings {
		for segment := range capsuleSegments {
			theta := 2 * math.Pi * float64(segment) / capsuleSegments
			mesh.Vertices = append(mesh.Vertices, point(ring.center, ring.offset, ring.radius, theta))
		}
	}
	top := uint32(len(mesh.Vertices))
	mesh.Vertices = append(mesh.Vertices, point(b, r, 0, 0)) // top pole

	ringStart := func(ring int) uint32 { return base + 1 + uint32(ring*capsuleSegments) }
	for segment := range capsuleSegments {
		nextSegment := (segment + 1) % capsuleSegments
		mesh.Indices = append(mesh.Indices, base, ringStart(0)+uint32(nextSegment), ringStart(0)+uint32(segment))
		last := ringStart(len(rings) - 1)
		mesh.Indices = append(mesh.Indices, top, last+uint32(segment), last+uint32(nextSegment))
		for ringIndex := 0; ringIndex+1 < len(rings); ringIndex++ {
			current := ringStart(ringIndex)
			upper := ringStart(ringIndex + 1)
			mesh.Indices = append(mesh.Indices,
				current+uint32(segment), current+uint32(nextSegment), upper+uint32(nextSegment),
				current+uint32(segment), upper+uint32(nextSegment), upper+uint32(segment),
			)
		}
	}
}

func crossProduct(a [3]float64, b [3]float64) [3]float64 {
	return [3]float64{a[1]*b[2] - a[2]*b[1], a[2]*b[0] - a[0]*b[2], a[0]*b[1] - a[1]*b[0]}
}

func normalize(vector [3]float64) [3]float64 {
	length := math.Sqrt(vector[0]*vector[0] + vector[1]*vector[1] + vector[2]*vector[2])

	return [3]float64{vector[0] / length, vector[1] / length, vector[2] / length}
}

func asObjects(value any) []map[string]any {
	array, _ := value.([]any)
	objects := make([]map[string]any, 0, len(array))
	for _, item := range array {
		if object, ok := item.(map[string]any); ok {
			objects = append(objects, object)
		}
	}

	return objects
}

func asBytes(value any) []byte {
	bytes, _ := value.([]byte)

	return bytes
}

func toInt(value any) int {
	switch typed := value.(type) {
	case int64:
		return int(typed)
	case uint64:
		return int(typed)
	case float64:
		return int(typed)
	default:
		return -1
	}
}

func toFloat(value any) float64 {
	switch typed := value.(type) {
	case float64:
		return typed
	case int64:
		return float64(typed)
	case uint64:
		return float64(typed)
	default:
		return 0
	}
}

func toString(value any) string {
	text, _ := value.(string)

	return text
}

func toLowerStrings(value any) []string {
	array, _ := value.([]any)
	var result []string
	for _, item := range array {
		if text, ok := item.(string); ok {
			result = append(result, strings.ToLower(text))
		}
	}

	return result
}

func toVec3(value any) [3]float32 {
	var result [3]float32
	switch typed := value.(type) {
	case []float64:
		for i := 0; i < 3 && i < len(typed); i++ {
			result[i] = float32(typed[i])
		}
	case []float32:
		for i := 0; i < 3 && i < len(typed); i++ {
			result[i] = typed[i]
		}
	case []any:
		for i := 0; i < 3 && i < len(typed); i++ {
			result[i] = float32(toFloat(typed[i]))
		}
	}

	return result
}
