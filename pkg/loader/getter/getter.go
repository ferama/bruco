package getter

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"strings"
)

type Getter interface {
	Download(string) (string, error)
	Cleanup()
}

type getterCommon struct {
	PayloadPath string
}

func (g *getterCommon) Cleanup() {
	os.Remove(g.PayloadPath)
}

func (g *getterCommon) getTmpFile(resourceURL string) (*os.File, error) {
	parsed, _ := url.Parse(resourceURL)
	path := parsed.Path
	segments := strings.Split(path, "/")
	fileName := segments[len(segments)-1]

	file, err := ioutil.TempFile(os.TempDir(), fmt.Sprintf("bruco_*_%s", fileName))
	if err != nil {
		return nil, fmt.Errorf("can't store file %s", err)
	}
	return file, nil
}
