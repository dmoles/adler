// +build mage

package main

import (
	"fmt"
	"github.com/dmoles/adler/server/util"
	"github.com/go-git/go-git/v5"
	"github.com/magefile/mage/mg"
	"io"
	"os"
	"path/filepath"
	"time"
)

func newPath(relPath string) (*path, error) {
	dir, err := isDir(relPath)
	if err != nil {
		return nil, err
	}
	if *dir {
		return nil, os.ErrInvalid
	}
	repoPath, err := toRepoPath(relPath)
	if err != nil {
		return nil, err
	}
	return &path{relPath, repoPath}, nil
}

type path struct {
	path     string
	repoPath string
}

func (p *path) Path() string {
	return p.path
}

func (p *path) ModTime() (*time.Time, error) {
	if status := p.gitStatus(); status == git.Unmodified {
		return p.gitModTime()
	}
	return p.fileModTime()
}

func (p *path) gitStatus() git.StatusCode {
	return gitStatus(p.repoPath)
}

func (p *path) ModTimeString() string {
	t, err := p.ModTime()
	if err != nil {
		return err.Error()
	}
	return t.Format(time.RFC3339Nano)
}

func (p *path) AsNewAs(path string) (*bool, error) {
	absPath, err := toAbsPath(path)
	if err != nil {
		return nil, err
	}
	dir, err := isDir(absPath)
	if err != nil {
		return nil, err
	}
	if *dir {
		return p.asNewAsAny(absPath)
	}
	return p.asNewAs(absPath)
}

func (p *path) Ignored() (*bool, error) {
	return gitIgnored(p.repoPath)
}

func (p *path) asNewAs(path string) (*bool, error) {
	p2, err := newPath(path)
	if err != nil {
		return nil, err
	}
	order, err := p.compareByModTime(p2)
	if err != nil {
		return nil, err
	}
	o := *order
	result := o >= 0

	if mg.Verbose() {
		var op string
		switch {
		case o < 0:
			op = "<"
		case o > 0:
			op = ">"
		default:
			op = "="
		}
		s1 := string([]byte { byte(p.gitStatus()) })
		s2 := string([]byte { byte(p2.gitStatus()) })
		msg := fmt.Sprintf("%s (%v: %s) %s %s (%v: %s)", p.repoPath, s1, p.ModTimeString(), op, p2.repoPath, s2, p2.ModTimeString())
		println(msg)
	}

	return &result, nil
}

func (p *path) asNewAsAny(dirPath string) (*bool, error) {
	result := true
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		ignored, err := gitIgnored(path)
		if err != nil {
			return err
		}
		if *ignored {
			return nil
		}

		r, err := p.asNewAs(path)
		if err != nil {
			return err
		}
		// TODO: is this really correct?
		result = result && *r
		if !result {
			return io.EOF
		}
		return nil
	})
	if err == nil || err == io.EOF {
		return &result, nil
	}
	return nil, err
}

func (p *path) compareByModTime(p2 *path) (*int, error) {
	mt1, err := p.ModTime()
	if err != nil {
		return nil, err
	}
	mt2, err := p2.ModTime()
	if err != nil {
		return nil, err
	}
	order := 0
	if mt1.Before(*mt2) {
		order = -1
	} else if mt1.After(*mt2) {
		order = 1
	}
	return &order, nil
}

func (p *path) fileModTime() (*time.Time, error) {
	info, err := os.Stat(p.repoPath)
	if err != nil {
		return nil, err
	}
	mtime := info.ModTime()
	return &mtime, nil
}

func (p *path) gitModTime() (*time.Time, error) {
	return gitCommitTime(p.repoPath)
}

func toAbsPath(path string) (string, error) {
	if filepath.IsAbs(path) {
		return path, nil
	}
	return filepath.Abs(path)
}

func toRepoPath(path string) (string, error) {
	absPath, err := toAbsPath(path)
	if err != nil {
		return "", err
	}
	rpath, err := filepath.Rel(util.ProjectRoot(), absPath)
	if err != nil {
		return "", err
	}
	return rpath, nil
}

func isDir(path string) (*bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	isDir := info.IsDir()
	return &isDir, nil
}
