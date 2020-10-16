package resources

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type dirResources struct {
	resourcesDir string
}

func (d *dirResources) Open(resourcePath string) (http.File, error) {
	fullPath := filepath.Join(d.resourcesDir, resourcePath)
	return os.Open(fullPath)
}

func (d *dirResources) Walk(walkFn filepath.WalkFunc) error {
	return filepath.Walk(d.resourcesDir, walkFn)
}

func (d *dirResources) AbsPath(resourcePath string) string {
	return filepath.Join(d.resourcesDir, resourcePath)
}

func (d *dirResources) RelativePath(absPath string) string {
	return strings.TrimPrefix(absPath, d.resourcesDir)
}

func newDirResources(resourcesDir string) Resources {
	resourcesDirAbs, err := filepath.Abs(resourcesDir)
	if err != nil {
		panic(err)
	}
	dr := dirResources{resourcesDir: resourcesDirAbs}
	return &dr
}
