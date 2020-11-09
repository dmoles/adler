package resources

import (
	_ "github.com/dmoles/adler/statik"
	"github.com/rakyll/statik/fs"
)

// ------------------------------------------------------------
// Exported

func Get(resourcePath string) (Resource, error) {
	return defaultBundle.Get(resourcePath)
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
