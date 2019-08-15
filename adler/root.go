package adler

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// Deprecated: TODO: use Resolver instead
type Root interface {
	Page(urlPath string) (Page, error)
}

func NewRoot(rootDir string) (Root, error) {
	info, err := os.Stat(rootDir)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("%#v is not a directory", rootDir)
	}
	rootDirAbs, err := filepath.Abs(rootDir)
	return &root{rootDirAbs: rootDirAbs}, nil
}

type root struct {
	rootDirAbs string
}

func (r *root) Page(urlPath string) (Page, error) {
	pathElements := strings.Split(urlPath, "/")
	for _, pathElement := range pathElements {
		if pathElement == ".." {
			return nil, invalidPath(urlPath)
		}
	}
	decodedPath, err := url.PathUnescape(urlPath)
	if err != nil {
		return nil, invalidPath(urlPath)
	}
	filePath := filepath.Join(r.rootDirAbs, decodedPath)
	return NewPage(filePath)
}

