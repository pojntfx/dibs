package utils

import (
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	"path/filepath"
	"time"
)

type GitAdvanced struct {
	UserName string
	WorkDir  string
	Token    string
}

func (gitAdvanced GitAdvanced) Clone(url string) error {
	_, err := git.PlainClone(gitAdvanced.WorkDir, false, &git.CloneOptions{
		URL:      url,
		Auth:     &http.BasicAuth{Username: gitAdvanced.UserName, Password: gitAdvanced.Token},
		Progress: nil,
	})

	return err
}

func (gitAdvanced GitAdvanced) AddCommitAndPush(email, commitMessage string) error {
	g, err := git.PlainOpen(filepath.Join(gitAdvanced.WorkDir))
	if err != nil {
		return err
	}

	wt, err := g.Worktree()
	if err != nil {
		return err
	}

	if _, err = wt.Add("."); err != nil {
		return err
	}

	if _, err = wt.Commit(commitMessage, &git.CommitOptions{
		Author: &object.Signature{
			Name:  gitAdvanced.UserName,
			Email: email,
			When:  time.Now(),
		},
	}); err != nil {
		return err
	}

	if err = g.Push(&git.PushOptions{
		Auth: &http.BasicAuth{Username: gitAdvanced.UserName, Password: gitAdvanced.Token},
	}); err != nil {
		return err
	}

	return nil
}
