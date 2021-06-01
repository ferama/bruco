package getter

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// HttpGetter
type HttpGetter struct {
	payloadPath string
}

// NewHttpGetter
func NewHttpGetter() *HttpGetter {
	return &HttpGetter{}
}

// Download
func (g *HttpGetter) Download(resourceURL string) (string, error) {
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

	g.payloadPath = file.Name()
	return file.Name(), nil
}

// Cleanup
func (g *HttpGetter) Cleanup() {
	os.Remove(g.payloadPath)
}
