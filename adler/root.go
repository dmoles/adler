package adler

import (
	"fmt"
	"io/ioutil"
	"os"
)

type Directory interface {
	Links() (link Links, err error)
}

func NewDirectory(path string) (dir Directory, err error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("%v is not a directory", path)
	}
	return &directory{path: path}, nil
}

type directory struct {
	path string
}

// TODO: figure out relative paths
func (d *directory) Links() (link Links, err error) {
	info, err := ioutil.ReadDir(d.path)
	if err != nil {
		return nil, err
	}
	links := make(Links, len(info))
	for i, l := range info {
		links[i] = NewLink(l)
	}
	return links, nil
}
