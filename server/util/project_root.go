package util

import (
	"fmt"
	"path/filepath"
	"runtime"
)

const packageName = "adler"

var projectRoot *string

func ProjectRoot() string {
	if projectRoot == nil {
		pr := findProjectRoot()
		projectRoot = &pr
	}
	return *projectRoot
}

func findProjectRoot() string {
	_, start, _, _ := runtime.Caller(0)
	for p := start; p != ""; p, _ = splitCleanly(p) {
		if filepath.Base(p) == packageName {
			return p
		}
	}
	errMsg := fmt.Sprintf("Unable to locate package '%s' starting from %s", packageName, start)
	panic(errMsg)
}

// filepath.Split() keeps trailing slashes and so can't be
// usefully called recursively
func splitCleanly(path string) (dir, base string) {
	dir = filepath.Dir(path)
	base = filepath.Base(path)
	return
}
