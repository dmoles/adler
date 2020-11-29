// +build mage

package main

import (
	"github.com/dmoles/adler/server/util"
	ignore "github.com/get-woke/go-gitignore"
	"github.com/go-git/go-git/v5"
	"os/exec"
	"time"
)

var repo *git.Repository
var wt *git.Worktree
var gi *ignore.GitIgnore

func gitStatus(path string) git.StatusCode {
	// cheap workaround for https://github.com/go-git/go-git/issues/119
	cmd := exec.Command("git", "status", "-s", path)
	output, err := cmd.Output()
	if err != nil {
		return git.Untracked
	}
	if len(output) >= 2 {
		for i := 0; i < 2; i++ {
			s := git.StatusCode(output[i])
			if s != git.Unmodified {
				return s
			}
		}
	}
	return git.Unmodified
}

func gitCommitTime(path string) (*time.Time, error) {
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

func gitIgnored(path string) (*bool, error) {
	gitIgnore, err := gitIgnore()
	if err != nil {
		return nil, err
	}
	ignored := gitIgnore.MatchesPath(path)
	return &ignored, nil
}

func repository() (*git.Repository, error) {
	if repo == nil {
		projectRoot := util.ProjectRoot()
		r, err := git.PlainOpen(projectRoot)
		if err != nil {
			return nil, err
		}
		repo = r
	}
	return repo, nil
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
