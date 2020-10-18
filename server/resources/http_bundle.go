package resources

import (
	"github.com/dmoles/adler/server/util"
	"github.com/rakyll/statik/fs"
	"net/http"
	"path"
	"path/filepath"
	"strings"
)

// ------------------------------
// httpBundle

// A Bundle implementation backed by an http.FileSystem
type httpBundle struct {
	desc       string
	resourceFS http.FileSystem
}

func (r *httpBundle) Open(resourcePath string) (http.File, error) {
	absPath := r.AbsPath(resourcePath)
	return r.resourceFS.Open(absPath)
}

func (r *httpBundle) Walk(walkFn filepath.WalkFunc) error {
	return fs.Walk(r.resourceFS, "/", walkFn)
}

func (r *httpBundle) AbsPath(resourcePath string) string {
	return path.Join("/", resourcePath)
}

func (r *httpBundle) RelativePath(absPath string) string {
	return strings.TrimPrefix(absPath, "/")
}

func (r *httpBundle) Get(resourcePath string) (Resource, error) {
	f, err := r.Open(resourcePath)
	if err != nil {
		return nil, err
	}
	defer util.CloseQuietly(f)
	info, err := f.Stat()
	if err != nil {
		return nil, err
	}
	return &resource{resourcePath, r, info}, nil
}
