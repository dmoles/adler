package resources

import (
	"github.com/rakyll/statik/fs"
	"net/http"
	"path"
	"path/filepath"
	"strings"
)

// ------------------------------
// httpResources

// A Resources implementation backed by an http.FileSystem
type httpResources struct {
	desc       string
	resourceFS http.FileSystem
}

func (r *httpResources) Open(resourcePath string) (http.File, error) {
	absPath := r.AbsPath(resourcePath)
	return r.resourceFS.Open(absPath)
}

func (r *httpResources) Walk(walkFn filepath.WalkFunc) error {
	return fs.Walk(r.resourceFS, "/", walkFn)
}

func (r *httpResources) AbsPath(resourcePath string) string {
	return path.Join("/", resourcePath)
}

func (r *httpResources) RelativePath(absPath string) string {
	return strings.TrimPrefix(absPath, "/")
}
