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
	rawUrl string

	archiveFilePath string
}

// NewLoader creates a new Loader object
func NewLoader(url string) *Loader {
	loader := &Loader{
		rawUrl: url,
	}

	return loader
}

// Load actually loads and prepare the function to be executed
func (l *Loader) Load() (string, error) {
	var err error
	var path string

	parsed, err := url.Parse(l.rawUrl)
	if err != nil {
		return "", err
	}
	switch lower := strings.ToLower(parsed.Scheme); lower {
	case "local":
		path = fmt.Sprintf("%s%s", parsed.Host, parsed.Path)
	case "http", "https":
		path, err = l.loadFromHttp()
		// log.Println(p)
		// path = "./hack/examples/basic"
	default:
		return "", fmt.Errorf("unsupported scheme: %s", parsed.Scheme)
	}

	return path, err
}

func (l *Loader) loadFromHttp() (string, error) {
	parsed, _ := url.Parse(l.rawUrl)
	path := parsed.Path
	segments := strings.Split(path, "/")
	fileName := segments[len(segments)-1]

	file, err := ioutil.TempFile(os.TempDir(), fmt.Sprintf("bruco_*_%s", fileName))
	if err != nil {
		return "", fmt.Errorf("can't store file %s", err)
	}
	defer file.Close()

	l.archiveFilePath = file.Name()
	client := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}
	resp, err := client.Get(l.rawUrl)
	if err != nil {
		return "", fmt.Errorf("can't store file %s", err)
	}
	defer resp.Body.Close()

	// Put content on file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return "", fmt.Errorf("can't store file %s", err)
	}

	archive := newArchive(file.Name())
	return archive.getHandlerPath()
}

func (l *Loader) Cleanup() {
	os.Remove(l.archiveFilePath)
}
