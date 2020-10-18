package resources

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type dirBundle struct {
	resourcesDir string
}

func (d *dirBundle) Open(resourcePath string) (http.File, error) {
	fullPath := filepath.Join(d.resourcesDir, resourcePath)
	return os.Open(fullPath)
}

func (d *dirBundle) Walk(walkFn filepath.WalkFunc) error {
	return filepath.Walk(d.resourcesDir, walkFn)
}

func (d *dirBundle) AbsPath(resourcePath string) string {
	return filepath.Join(d.resourcesDir, resourcePath)
}

func (d *dirBundle) RelativePath(absPath string) string {
	return strings.TrimPrefix(absPath, d.resourcesDir)
}

func (d *dirBundle) Get(resourcePath string) (Resource, error) {
	fullPath := filepath.Join(d.resourcesDir, resourcePath)
	info, err := os.Stat(fullPath)
	if err != nil {
		return nil, err
	}
	return &resource{resourcePath, d, info}, nil
}

func newDirResources(resourcesDir string) Bundle {
	resourcesDirAbs, err := filepath.Abs(resourcesDir)
	if err != nil {
		panic(err)
	}
	return &dirBundle{resourcesDir: resourcesDirAbs}
}
