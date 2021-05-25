package loader

import (
	"fmt"
	"path/filepath"
	"strings"
)

type archive struct {
	path string
}

func newArchive(path string) *archive {
	a := &archive{
		path: path,
	}
	return a
}

func (a *archive) getHandlerPath() (string, error) {
	ext := filepath.Ext(a.path)
	var path string
	var err error

	switch lower := strings.ToLower(ext); lower {
	case ".py":
		path = a.path
	default:
		err = fmt.Errorf("unsupported extension: %s", ext)

	}
	return path, err
}
