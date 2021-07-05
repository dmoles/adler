package resources

import (
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/dmoles/adler/server/util"
)

// ------------------------------------------------------------
// Exported

// Bundle A bundle of resources
type Bundle interface {
	Open(resourcePath string) (fs.File, error)
	Walk(walkFn fs.WalkDirFunc) error
	AbsPath(resourcePath string) string
	RelativePath(absPath string) string
	Get(resourcePath string) (Resource, error)
}

func newDirBundle(resourcesDir string) Bundle {
	resourcesDirAbs, err := filepath.Abs(resourcesDir)
	if err != nil {
		panic(err)
	}
	return &bundle{resourcesDir, os.DirFS(resourcesDirAbs)}
}

// ------------------------------------------------------------
// Unexported

type bundle struct {
	desc       string
	resourceFS fs.FS
}

func (r *bundle) Open(resourcePath string) (fs.File, error) {
	absPath := r.AbsPath(resourcePath)
	return r.resourceFS.Open(absPath)
}

func (r *bundle) Walk(walkFn fs.WalkDirFunc) error {
	return fs.WalkDir(r.resourceFS, ".", walkFn)
}

func (r *bundle) AbsPath(resourcePath string) string {
	return path.Join(".", resourcePath)
}

func (r *bundle) RelativePath(absPath string) string {
	return strings.TrimPrefix(absPath, "./")
}

func (r *bundle) Get(resourcePath string) (Resource, error) {
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
