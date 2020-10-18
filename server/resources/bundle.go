package resources

import (
	"net/http"
	"path/filepath"
)

// A bundle of resources
type Bundle interface {
	Open(resourcePath string) (http.File, error)
	Walk(walkFn filepath.WalkFunc) error
	AbsPath(resourcePath string) string
	RelativePath(absPath string) string
	Get(resourcePath string) (Resource, error)
}
