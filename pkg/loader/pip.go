package loader

import (
	"io/ioutil"
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

func findRequirementsFile(path string) (string, bool) {
	var fileHandler *os.File
	var err error

	fileHandler, err = os.Open(path)
	if err != nil {
		return "", false
	}
	fi, err := fileHandler.Stat()
	if err != nil {
		return "", false
	}
	if fi.IsDir() {
		entries, _ := ioutil.ReadDir(path)
		if len(entries) > 0 {
			reqPath := filepath.Join(path, entries[0].Name(), "requirements.txt")
			if fileExists(reqPath) {
				return reqPath, true
			}
		}
	}
	return "", false
}

func runPip(path string) error {
	if reqFilePath, found := findRequirementsFile(path); found {
		// TODO: run pip install
		log.Println("found: ", reqFilePath)
	}

	return nil
}
