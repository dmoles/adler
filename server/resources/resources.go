package resources

import (
	"github.com/dmoles/adler/server/util"
	_ "github.com/dmoles/adler/statik"
	"github.com/rakyll/statik/fs"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

// ------------------------------------------------------------
// Exported

type Resources interface {
	Open(resourcePath string) (http.File, error)
	Walk(walkFn filepath.WalkFunc) error
	AbsPath(resourcePath string) string
	RelativePath(absPath string) string
}

func Open(resourcePath string) (http.File, error) {
	return defaultResources.Open(resourcePath)
}

func Stat(resources Resources, resourcePath string) (os.FileInfo, error) {
	f, err := resources.Open(resourcePath)
	if err != nil {
		return nil, err
	}
	defer util.CloseQuietly(f)
	return f.Stat()
}

func Read(resources Resources, resourcePath string) ([]byte, error) {
	f, err := resources.Open(resourcePath)
	if err != nil {
		return nil, err
	}
	defer util.CloseQuietly(f)
	return ioutil.ReadAll(f)
}

// ------------------------------------------------------------
// Unexported

var defaultResources Resources = fromStatikNamespace("adler")

func fromStatikNamespace(namespace string) *httpResources {
	staticFS, err := fs.NewWithNamespace(namespace)
	if err != nil {
		panic(err)
	}
	rr := httpResources{desc: namespace, resourceFS: staticFS}
	return &rr
}
