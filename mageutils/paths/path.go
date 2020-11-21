// +build mage

package paths

import (
	"github.com/dmoles/adler/mageutils"
	"github.com/dmoles/adler/mageutils/repo"
	"github.com/go-git/go-git/v5"
	"io"
	"os"
	"path/filepath"
	"time"
)

// ------------------------------------------------------------
// Exported

type Path interface {
	Path() string
	ModTime() (*time.Time, error)
	AsNewAs(path string) (*bool, error)
	Ignored() (*bool, error)
}

func New(relPath string) (Path, error) {
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

// ------------------------------------------------------------
// Unexported

type path struct {
	path     string
	repoPath string
}

func (p *path) Path() string {
	return p.path
}

func (p *path) ModTime() (*time.Time, error) {
	if status := repo.Status(p.repoPath); status == git.Unmodified {
		return p.gitModTime()
	}
	return p.fileModTime()
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
	return repo.Ignored(p.repoPath)
}

func (p *path) asNewAs(path string) (*bool, error) {
	p2, err := New(path)
	if err != nil {
		return nil, err
	}
	order, err := p.compareTo(p2)
	if err != nil {
		return nil, err
	}
	result := *order >= 0
	return &result, nil
}

func (p *path) asNewAsAny(dirPath string) (*bool, error) {
	result := true
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		// TODO: ignored?
		r, err := p.asNewAs(path)
		if err != nil {
			return err
		}
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

func (p *path) compareTo(p2 Path) (*int, error) {
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
	return repo.CommitTime(p.repoPath)
}

// ------------------------------
// Unexported functions

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
	rpath, err := filepath.Rel(mageutils.ProjectRoot(), absPath)
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