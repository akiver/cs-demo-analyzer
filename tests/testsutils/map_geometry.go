package testsutils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/akiver/cs-demo-analyzer/internal/visibility"
)

// mapGeometryDir is the folder where download-map-geometry.sh downloads the CS2 map geometry.
const mapGeometryDir = "../cs2-map-geometry"

// GetMapGeometryDir returns the folder containing the CS2 map geometry used to compute the time to damage/kill.
// It fails the test if the folder doesn't exist since the snapshots would not match.
func GetMapGeometryDir(t *testing.T) string {
	t.Helper()

	if _, err := os.Stat(mapGeometryDir); err != nil {
		t.Fatalf("map geometry folder %s not found, run ./download-map-geometry.sh: %v", mapGeometryDir, err)
	}

	return mapGeometryDir
}

// AssertMapGeometry fails the test if the geometry of the map is not in the map geometry folder.
// The analyzer only logs a missing geometry and doesn't compute the time to damage/kill, the snapshot would then not
// match or would be updated without these stats.
func AssertMapGeometry(t *testing.T, mapName string) {
	t.Helper()

	path := filepath.Join(mapGeometryDir, mapName+visibility.GeometryFileExtension)
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("map geometry %s not found, make sure %s is in MAPS of download-map-geometry.sh and run it: %v", path, mapName, err)
	}
}
