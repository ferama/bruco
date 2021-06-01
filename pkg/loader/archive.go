package loader

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type archive struct {
	sourcePath  string
	toCleanPath string
}

func newArchive(path string) *archive {
	a := &archive{
		sourcePath:  path,
		toCleanPath: "",
	}
	return a
}

// if its an archive, extract it. If not get the orginal path
func (a *archive) getResourcePath() (string, error) {
	ext := filepath.Ext(a.sourcePath)
	var path string
	var err error

	switch lower := strings.ToLower(ext); lower {
	case ".zip":
		path, err = a.unzip(a.sourcePath)
		a.toCleanPath = path
	default:
		path = a.sourcePath

	}
	return path, err
}

func (a *archive) cleanup() {
	if a.toCleanPath != "" {
		os.RemoveAll(a.toCleanPath)
	}
}

func (a *archive) unzip(src string) (string, error) {
	dest, err := ioutil.TempDir(os.TempDir(), "bruco-resource-")
	if err != nil {
		log.Fatal(err)
	}
	r, err := zip.OpenReader(src)
	if err != nil {
		return "", err
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)
		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return "", fmt.Errorf("%s: illegal file path", fpath)
		}

		if f.FileInfo().IsDir() {
			// Make Folder
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}
		// Make File
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return "", err
		}
		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return "", err
		}
		rc, err := f.Open()
		if err != nil {
			return "", err
		}

		_, err = io.Copy(outFile, rc)

		// Close the file without defer to close before next iteration of loop
		outFile.Close()
		rc.Close()

		if err != nil {
			return "", err

		}

	}

	return dest, nil
}
