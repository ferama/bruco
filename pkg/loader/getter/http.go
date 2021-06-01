package getter

import (
	"fmt"
	"io"
	"net/http"
)

// HttpGetter
// It allows to load a function from http server
// Example:
// $ bruco https://github.com/ferama/bruco/raw/main/hack/examples/zipped/sentiment.zip
type HttpGetter struct {
	getterCommon
}

// NewHttpGetter
func NewHttpGetter() *HttpGetter {
	return &HttpGetter{}
}

// Download
func (g *HttpGetter) Download(resourceURL string) (string, error) {
	file, err := g.getTmpFile(resourceURL)
	if err != nil {
		return "", fmt.Errorf("can't store file %s", err)
	}
	defer file.Close()
	// do not forget to set this one for correct cleanup by the parent class
	// getterCommon
	g.PayloadPath = file.Name()

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
