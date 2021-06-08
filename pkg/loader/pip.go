package loader

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/ferama/bruco/pkg/common"
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
		pythonPath, err := common.FindPython()
		if err != nil {
			log.Fatalln("can't find python executable")
		}
		args := []string{
			pythonPath, "-m", "pip", "install", "-r", reqPath,
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
