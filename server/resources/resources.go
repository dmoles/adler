package resources

import (
	_ "github.com/dmoles/adler/statik"
	"github.com/rakyll/statik/fs"
	"net/http"
	"path"
)

func Open(resourcePath string) (http.File, error) {
	if !path.IsAbs(resourcePath) {
		resourcePath = path.Join("/", resourcePath)
	}
	return resourceFS.Open(resourcePath)
}

var resourceFS http.FileSystem

func init() {
	staticFS, err := fs.NewWithNamespace("adler")
	if err != nil {
		panic(err)
	}
	resourceFS = staticFS
}
