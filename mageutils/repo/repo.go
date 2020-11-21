// +build mage

package repo

import (
	"github.com/dmoles/adler/mageutils"
	ignore "github.com/get-woke/go-gitignore"
	"github.com/go-git/go-git/v5"
	"time"
)

// ------------------------------------------------------------
// Exported

func Status(path string) git.StatusCode {
	status, err := status()
	if err != nil {
		return git.Untracked
	}
	fs := status.File(path)
	return fs.Worktree
}

func CommitTime(path string) (*time.Time, error) {
	repo, err := repository()
	if err != nil {
		return nil, err
	}
	lo := git.LogOptions{
		FileName: &path,
		Order:    git.LogOrderCommitterTime,
	}
	commits, err := repo.Log(&lo)
	if err != nil {
		return nil, err
	}
	lastCommit, err := commits.Next()
	if err != nil {
		return nil, err
	}
	commitTime := lastCommit.Committer.When
	return &commitTime, nil
}

func Ignored(path string) (*bool, error) {
	gitIgnore, err := gitIgnore()
	if err != nil {
		return nil, err
	}
	ignored := gitIgnore.MatchesPath(path)
	return &ignored, nil
}

// ------------------------------------------------------------
// Unexported

var repo *git.Repository
var wt *git.Worktree
var gi *ignore.GitIgnore

func repository() (*git.Repository, error) {
	if repo == nil {
		projectRoot := mageutils.ProjectRoot()
		r, err := git.PlainOpen(projectRoot)
		if err != nil {
			return nil, err
		}
		repo = r
	}
	return repo, nil
}

func worktree() (*git.Worktree, error) {
	if wt == nil {
		r, err := repository()
		if err != nil {
			return nil, err
		}
		w, err := r.Worktree()
		if err != nil {
			return nil, err
		}
		wt = w
	}
	return wt, nil
}

func status() (git.Status, error) {
	wt, err := worktree()
	if err != nil {
		return nil, err
	}
	s, err := wt.Status()
	if err != nil {
		return nil, err
	}
	return s, nil
}

func gitIgnore() (*ignore.GitIgnore, error) {
	if gi == nil {
		gitIgnore, err := ignore.CompileIgnoreFile(".gitignore")
		if err != nil {
			return nil, err
		}
		gi = gitIgnore
	}
	return gi, nil
}
