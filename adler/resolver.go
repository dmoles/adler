package adler

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

type Resolver struct {
	rootDir string
}

func NewResolver(rootDir string) (*Resolver, error) {
	rootDirAbs, err := filepath.Abs(rootDir)
	if err != nil {
		return nil, err
	}
	info, err := os.Stat(rootDirAbs)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("not a directory: %s", rootDirAbs)
	}
	return &Resolver{rootDir: rootDirAbs}, nil
}

func (r *Resolver) Resolve(urlPath string) (string, error) {
	pathElements := strings.Split(urlPath, "/")
	for _, pathElement := range pathElements {
		if pathElement == ".." {
			return "", invalidPath(urlPath)
		}
	}
	decodedPath, err := url.PathUnescape(urlPath)
	if err != nil {
		return "", invalidPath(urlPath)
	}
	filePath := filepath.Join(r.rootDir, decodedPath)
	return filePath, nil
}