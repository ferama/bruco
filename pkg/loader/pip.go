package loader

import (
	"log"
	"os"
	"path/filepath"
)

func fileExists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	} else {
		return false
	}
}

func runPip(workingDir string) error {
	reqPath := filepath.Join(workingDir, "requirements.txt")
	if fileExists(reqPath) {
		log.Println("found: ", reqPath)
	}

	return nil
}
