// Command map-geometry extracts the collision geometry blocking players' line of sight from CS2 maps and writes one
// <map>.tri file per map, used by the analyzer to compute visibility based stats (e.g. time to damage).
// It also writes a manifest.json listing the size and SHA-256 of every .tri file of the output folder.
//
// Usage:
//
//	go run ./cmd/map-geometry -maps-dir "<CS2>/game/csgo/maps" -output ./map-geometry [-maps de_dust2,de_mirage]
//	go run ./cmd/map-geometry -vpk "<CS2>/game/csgo/maps/de_dust2.vpk" -output ./map-geometry
package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"time"

	"github.com/akiver/cs-demo-analyzer/internal/source2"
	"github.com/akiver/cs-demo-analyzer/internal/visibility"
)

func defaultMapsDir() string {
	var candidates []string
	switch runtime.GOOS {
	case "windows":
		candidates = []string{`C:\Program Files (x86)\Steam\steamapps\common\Counter-Strike Global Offensive\game\csgo\maps`}
	case "linux":
		home, _ := os.UserHomeDir()
		candidates = []string{
			filepath.Join(home, ".steam/steam/steamapps/common/Counter-Strike Global Offensive/game/csgo/maps"),
			filepath.Join(home, ".local/share/Steam/steamapps/common/Counter-Strike Global Offensive/game/csgo/maps"),
		}
	}

	for _, candidate := range candidates {
		if info, err := os.Stat(candidate); err == nil && info.IsDir() {
			return candidate
		}
	}

	return ""
}

func main() {
	mapsDir := flag.String("maps-dir", defaultMapsDir(), "CS2 maps folder (<CS2>/game/csgo/maps)")
	outputDir := flag.String("output", "map-geometry", "Output folder")
	mapsList := flag.String("maps", "", "Comma separated list of maps to extract (default all maps with collision geometry)")
	vpkPath := flag.String("vpk", "", "Path of a map VPK file (<map>.vpk) to extract only this map, -maps-dir is ignored")
	flag.Parse()

	var mapNames []string
	if *vpkPath != "" {
		if *mapsList != "" {
			fmt.Fprintln(os.Stderr, "-vpk and -maps can't be used together")
			os.Exit(2)
		}
		if filepath.Ext(*vpkPath) != ".vpk" {
			fmt.Fprintf(os.Stderr, "%s is not a .vpk file\n", *vpkPath)
			os.Exit(2)
		}
		// Fail before creating the output folder.
		if _, err := os.Stat(*vpkPath); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		// The map VPK is named after the map like in the maps folder.
		*mapsDir = filepath.Dir(*vpkPath)
		mapNames = []string{strings.TrimSuffix(filepath.Base(*vpkPath), ".vpk")}
	} else {
		if *mapsDir == "" {
			fmt.Fprintln(os.Stderr, "-maps-dir is required")
			flag.Usage()
			os.Exit(2)
		}

		var err error
		mapNames, err = listMaps(*mapsDir, *mapsList)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

	if err := os.MkdirAll(*outputDir, 0o755); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	failed := false
	for _, mapName := range mapNames {
		if err := extractMap(*mapsDir, *outputDir, mapName); err != nil {
			if errors.Is(err, source2.ErrVPKFileNotFound) {
				fmt.Printf("%s: skipped, no world physics\n", mapName)
				continue
			}
			fmt.Fprintf(os.Stderr, "%s: %v\n", mapName, err)
			failed = true
		}
	}

	if err := writeManifest(*outputDir); err != nil {
		fmt.Fprintf(os.Stderr, "failed to write the manifest: %v\n", err)
		failed = true
	}

	if failed {
		os.Exit(1)
	}
}

func listMaps(mapsDir string, mapsList string) ([]string, error) {
	if mapsList != "" {
		var names []string
		for name := range strings.SplitSeq(mapsList, ",") {
			if name = strings.TrimSpace(name); name != "" {
				names = append(names, name)
			}
		}

		return names, nil
	}

	entries, err := os.ReadDir(mapsDir)
	if err != nil {
		return nil, err
	}

	var names []string
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || filepath.Ext(name) != ".vpk" {
			continue
		}
		name = strings.TrimSuffix(name, ".vpk")
		if shouldExtractMap(name) {
			names = append(names, name)
		}
	}
	slices.Sort(names)

	return names, nil
}

// shouldExtractMap returns false for the maps not worth extracting when no -maps list is given: main menu lobby,
// workshop previews, vanity (weapon/agent preview), arms race...
func shouldExtractMap(name string) bool {
	for _, prefix := range []string{"ar_", "lobby_", "rush_", "workshop_preview_", "graphics_settings"} {
		if strings.HasPrefix(name, prefix) {
			return false
		}
	}

	return !strings.HasSuffix(name, "_night") && !strings.HasSuffix(name, "_vanity")
}

func extractMap(mapsDir string, outputDir string, mapName string) error {
	start := time.Now()
	vpk, err := source2.OpenVPK(filepath.Join(mapsDir, mapName+".vpk"))
	if err != nil {
		return err
	}

	resource, err := vpk.ReadFile("maps/" + mapName + "/world_physics.vmdl_c")
	if err != nil {
		return err
	}

	block, err := source2.ResourceBlock(resource, "PHYS")
	if err != nil {
		return err
	}

	physics, err := source2.DecodeKV3(block)
	if err != nil {
		return err
	}

	mesh, stats, err := source2.ExtractCollisionMesh(physics, source2.IsBlockingLineOfSight)
	if err != nil {
		return err
	}

	vertices := make([]visibility.Vec3, len(mesh.Vertices))
	for i, vertex := range mesh.Vertices {
		vertices[i] = visibility.Vec3(vertex)
	}
	bvh := visibility.BuildBVH(vertices, mesh.Indices)

	outputPath := filepath.Join(outputDir, mapName+visibility.GeometryFileExtension)
	if err := writeGeometryFile(outputPath, bvh); err != nil {
		return err
	}

	info, err := os.Stat(outputPath)
	if err != nil {
		return err
	}

	fmt.Printf(
		"%s: %d triangles (%d hulls, %d mesh triangles, %d capsules, %d spheres), skipped %d shapes (%d mesh triangles), %d BVH nodes, %.1f MB, %s\n",
		mapName, bvh.TriangleCount(), stats.Hulls, stats.MeshTriangles, stats.Capsules, stats.Spheres, stats.SkippedShapes,
		stats.SkippedTriangles, len(bvh.Nodes), float64(info.Size())/1024/1024, time.Since(start).Round(time.Millisecond),
	)
	if stats.Compounds > 0 {
		fmt.Printf("%s: warning: %d compound shapes are not supported and were ignored\n", mapName, stats.Compounds)
	}

	return nil
}

// writeGeometryFile writes the geometry to a temporary file renamed to path on success, so that a failed write never
// leaves a truncated geometry file that would end up in the manifest. A previous file at path is kept on failure.
func writeGeometryFile(path string, bvh *visibility.BVH) error {
	// The extension must not be the geometry one, the manifest lists all the geometry files of the folder.
	tmpPath := path + ".tmp"
	file, err := os.Create(tmpPath)
	if err != nil {
		return err
	}

	err = visibility.WriteGeometry(file, bvh)
	if closeErr := file.Close(); err == nil {
		err = closeErr
	}
	if err == nil {
		err = os.Rename(tmpPath, path)
	}
	if err != nil {
		os.Remove(tmpPath)
	}

	return err
}
