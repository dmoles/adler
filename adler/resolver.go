package adler

import (
	"fmt"
	"os"
	"path/filepath"
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
