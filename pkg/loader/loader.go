package loader

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
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
func (l *Loader) Load(resourceURL string) (string, error) {
	var err error
	var path string

	parsed, err := url.Parse(resourceURL)
	if err != nil {
		return "", err
	}
	switch lower := strings.ToLower(parsed.Scheme); lower {
	case "":
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
	path, err = archive.getResourcePath()
	l.archive = archive
	return path, err
}

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

func (l *Loader) Cleanup() {
	os.Remove(l.payloadPath)
	l.archive.cleanup()
}
