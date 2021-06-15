package loader

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/ferama/bruco/pkg/common"
	"github.com/ferama/bruco/pkg/loader/getter"
)

const containerConfigPath = "/bruco/config.yaml"

// Loader object will parse the handlerURL and prepare the function
// to be executed from processor
type Loader struct {
	archive *archive
	// payloadPath string
	getter getter.Getter
}

// NewLoader creates a new Loader object
func NewLoader() *Loader {
	loader := &Loader{}

	return loader
}

// Load actually loads and prepare the function to be executed
func (l *Loader) runGetter(resourceURL string) (string, error) {
	var err error
	var path string

	parsed, err := url.Parse(resourceURL)
	if err != nil {
		return "", err
	}
	switch lower := strings.ToLower(parsed.Scheme); lower {
	case "":
		// no scheme. It's a local file
		path = fmt.Sprintf("%s%s", parsed.Host, parsed.Path)
	case "http", "https":
		l.getter = getter.NewHttpGetter()
		path, err = l.getter.Download(resourceURL)
	case "s3", "s3s":
		l.getter = getter.NewS3Getter()
		path, err = l.getter.Download(resourceURL)
	default:
		return "", fmt.Errorf("unsupported scheme: %s", parsed.Scheme)
	}
	if err != nil {
		return "", err
	}
	archive := newArchive(path)
	// if it's an archive, extract it. If not get the orginal path
	path, err = archive.getResourcePath()
	l.archive = archive

	return path, err
}

// LoadFunction loads the function and returns a pointer to the config yaml
// file
func (l *Loader) LoadFunction(fileURL string) (*os.File, string, error) {
	var fileHandler *os.File
	var err error
	workingDir := ""

	filePath, err := l.runGetter(fileURL)
	if err != nil {
		return nil, workingDir, err
	}

	fileHandler, err = os.Open(filePath)
	if err != nil {
		return nil, workingDir, err
	}

	fi, err := fileHandler.Stat()
	if err != nil {
		return nil, workingDir, err
	}

	workingDir = filepath.Dir(fileHandler.Name())
	foundInPackage := true
	if fi.IsDir() {
		// it's a directory. Search for a config.yaml inside
		path := filepath.Join(filePath, "config.yaml")
		fileHandler.Close()
		fileHandler, err = os.Open(path)

		if err != nil {
			// config.yaml not found. Try to get it from a subdir.
			// |
			// . functiondir
			// 		|
			//		. config.yaml
			//		. handler.py
			//
			entries, _ := ioutil.ReadDir(filePath)
			if len(entries) > 0 {
				workingDir = filepath.Join(filePath, entries[0].Name())
				path := filepath.Join(workingDir, "config.yaml")
				fileHandler.Close()
				fileHandler, err = os.Open(path)
				if err != nil {
					foundInPackage = false
				}
			} else {
				return nil, workingDir, err
			}
		}
	}

	// if container conf exists, take precedence
	containerConf, err := os.Open(containerConfigPath)
	if err == nil {
		fileHandler = containerConf
	} else {
		if !foundInPackage {
			return nil, workingDir, err
		}
	}

	disablePip, _ := common.GetenvBool("BRUCO_DISABLE_PIP")

	if !disablePip {
		if err := runPip(workingDir); err != nil {
			return nil, "", err
		}
	}
	// If fileURL is not a directory I'm assuming that I'm running bruco
	// against a config.yaml directly
	return fileHandler, workingDir, nil
}

// Cleanup removes temporary files used during the loading process
func (l *Loader) Cleanup() {
	if l.getter != nil {
		l.getter.Cleanup()
	}
	if l.archive != nil {
		l.archive.cleanup()
	}
}
