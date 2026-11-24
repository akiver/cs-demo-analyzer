package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/akiver/cs-demo-analyzer/internal/visibility"
)

const manifestFileName = "manifest.json"

// manifest describes the map geometry files of a folder, it lets clients download only the files that changed and
// verify them.
type manifest struct {
	// Version of the geometry file format, clients must ignore files of a version they don't support.
	FormatVersion int                    `json:"formatVersion"`
	Maps          map[string]manifestMap `json:"maps"`
}

type manifestMap struct {
	File   string `json:"file"`
	Size   int64  `json:"size"`
	SHA256 string `json:"sha256"`
}

// writeManifest writes the manifest of all the map geometry files of the folder, not only the ones just extracted.
func writeManifest(outputDir string) error {
	entries, err := os.ReadDir(outputDir)
	if err != nil {
		return err
	}

	result := manifest{
		FormatVersion: visibility.GeometryFileVersion,
		Maps:          make(map[string]manifestMap),
	}
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || filepath.Ext(name) != visibility.GeometryFileExtension {
			continue
		}

		mapEntry, err := newManifestMap(filepath.Join(outputDir, name))
		if err != nil {
			return err
		}
		result.Maps[strings.TrimSuffix(name, visibility.GeometryFileExtension)] = mapEntry
	}

	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(outputDir, manifestFileName), append(data, '\n'), 0o644)
}

func newManifestMap(path string) (manifestMap, error) {
	file, err := os.Open(path)
	if err != nil {
		return manifestMap{}, err
	}
	defer file.Close()

	hash := sha256.New()
	size, err := io.Copy(hash, file)
	if err != nil {
		return manifestMap{}, err
	}

	return manifestMap{
		File:   filepath.Base(path),
		Size:   size,
		SHA256: hex.EncodeToString(hash.Sum(nil)),
	}, nil
}
