package testsutils

import (
	"path/filepath"
)

// GetDemoPath returns the path to the demo for testing.
func GetDemoPath(gameFolder string, name string) string {
	demosFolderPath := "../cs-demos/" + gameFolder + "/"
	return filepath.Join(demosFolderPath + name + ".dem")
}
