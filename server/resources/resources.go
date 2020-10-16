package resources

import (
	"github.com/dmoles/adler/server/util"
	_ "github.com/dmoles/adler/statik"
	"github.com/rakyll/statik/fs"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

// ------------------------------------------------------------
// Exported

// ------------------------------
// An individual resource

type Resource interface {
	Path() string
	Bundle() Bundle
	Stat() os.FileInfo
	Open() (http.File, error)
	Read() ([]byte, error)
	Copy(w io.Writer) (int64, error)
	ContentLength() int64
	ContentType() string
}

type resource struct {
	path string
	bundle Bundle
	info os.FileInfo
}

func (r *resource) Path() string {
	return r.path
}

func (r *resource) Bundle() Bundle {
	return r.bundle
}

func (r *resource) Stat() os.FileInfo {
	return r.info
}

func (r *resource) Open() (http.File, error) {
	return r.bundle.Open(r.path)
}

func (r *resource) Read() ([]byte, error) {
	f, err := r.Open()
	if err != nil {
		return nil, err
	}
	defer util.CloseQuietly(f)
	return ioutil.ReadAll(f)
}

func (r *resource) Copy(w io.Writer) (int64, error) {
	f, err := r.Open()
	if err != nil {
		return 0, err
	}
	defer util.CloseQuietly(f)
	return io.Copy(w, f)
}

func (r *resource) ContentLength() int64 {
	return r.info.Size()
}

func (r *resource) ContentType() string {
	return util.ContentType(r.path)
}

// ------------------------------
// A bundle of resources

type Bundle interface {
	Open(resourcePath string) (http.File, error)
	Walk(walkFn filepath.WalkFunc) error
	AbsPath(resourcePath string) string
	RelativePath(absPath string) string
	Get(resourcePath string) (Resource, error)
}

func Get(resourcePath string) (Resource, error) {
	return defaultBundle.Get(resourcePath)
}

func Open(resourcePath string) (http.File, error) {
	return defaultBundle.Open(resourcePath)
}

// ------------------------------------------------------------
// Unexported

var defaultBundle Bundle = statikFSBundle("adler")

func statikFSBundle(namespace string) *httpBundle {
	statikFS, err := fs.NewWithNamespace(namespace)
	if err != nil {
		panic(err)
	}
	rr := httpBundle{desc: namespace, resourceFS: statikFS}
	return &rr
}
