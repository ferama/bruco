package loader

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// Loader object will parse the handlerURL and prepare the function
// to be executed from processor
type Loader struct {
	archive     *archive
	payloadPath string
}

// NewLoader creates a new Loader object
func NewLoader() *Loader {
	loader := &Loader{}

	return loader
}

// Load actually loads and prepare the function to be executed
func (l *Loader) load(resourceURL string) (string, error) {
	var err error
	var path string

	parsed, err := url.Parse(resourceURL)
	if err != nil {
		return "", err
	}
	switch lower := strings.ToLower(parsed.Scheme); lower {
	case "":
		// no scheme. Its a local file
		path = fmt.Sprintf("%s%s", parsed.Host, parsed.Path)
	case "http", "https":
		path, err = l.httpDownload(resourceURL)
		l.payloadPath = path
	default:
		return "", fmt.Errorf("unsupported scheme: %s", parsed.Scheme)
	}
	if err != nil {
		return "", err
	}
	archive := newArchive(path)
	// if its an archive, extract it. If not get the orginal path
	path, err = archive.getResourcePath()
	l.archive = archive

	return path, err
}

// search for a config file
func (l *Loader) GetConfig(fileURL string) (*os.File, error) {
	var fileHandler *os.File
	var err error

	filePath, err := l.load(fileURL)
	if err != nil {
		return nil, err
	}

	fileHandler, err = os.Open(filePath)
	if err != nil {
		return nil, err
	}
	fi, err := fileHandler.Stat()
	if err != nil {
		return nil, err
	}
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
				path := filepath.Join(filePath, entries[0].Name(), "config.yaml")
				fileHandler.Close()
				fileHandler, err = os.Open(path)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	runPip(filepath.Dir(fileHandler.Name()))

	// fileURL is not a directory.
	// Assuming that I'm running brugo against a config.yaml directly
	return fileHandler, nil
}

// downloads the resource from an http server
func (l *Loader) httpDownload(resourceURL string) (string, error) {
	parsed, _ := url.Parse(resourceURL)
	path := parsed.Path
	segments := strings.Split(path, "/")
	fileName := segments[len(segments)-1]

	file, err := ioutil.TempFile(os.TempDir(), fmt.Sprintf("bruco_*_%s", fileName))
	if err != nil {
		return "", fmt.Errorf("can't store file %s", err)
	}
	defer file.Close()

	client := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}
	resp, err := client.Get(resourceURL)
	if err != nil {
		return "", fmt.Errorf("can't store file %s", err)
	}
	defer resp.Body.Close()

	// Put content on file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return "", fmt.Errorf("can't store file %s", err)
	}

	return file.Name(), nil
}

// Cleanup removes temporary files used during the loading process
func (l *Loader) Cleanup() {
	os.Remove(l.payloadPath)
	l.archive.cleanup()
}
