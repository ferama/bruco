package loader

import (
	"fmt"
	"net/url"
	"strings"
)

// Loader object will parse the handlerURL and prepare the function
// to be executed from processor
type Loader struct {
	rawUrl string
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
	parsed, err := url.Parse(l.rawUrl)
	if err != nil {
		return "", err
	}
	var path string
	switch lower := strings.ToLower(parsed.Scheme); lower {
	case "local":
		path = fmt.Sprintf("%s%s", parsed.Host, parsed.Path)
	default:
		return "", fmt.Errorf("unsupported scheme: %s", parsed.Scheme)
	}

	return path, nil
}
