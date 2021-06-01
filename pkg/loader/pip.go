package loader

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func fileExists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	} else {
		return false
	}
}

func findPip() (string, error) {
	var path string
	var err error
	path, err = exec.LookPath("pip3")
	if err != nil {
		path, err = exec.LookPath("pip")
		if err != nil {
			return "", err
		}
	}
	return path, nil
}

func runPip(workingDir string) error {
	reqPath := filepath.Join(workingDir, "requirements.txt")
	if fileExists(reqPath) {
		pipPath, err := findPip()
		if err != nil {
			log.Println("could not find pip: ", err)
		}
		args := []string{
			pipPath, "install", "-r", reqPath,
		}
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return err
		}
	}

	return nil
}
