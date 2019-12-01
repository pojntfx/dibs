package utils

import (
	git "gopkg.in/src-d/go-git.v4"
	gitconf "gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"path/filepath"
	"time"
)

type Git struct {
	RemoteName    string
	RemoteURL     string
	UserName      string
	UserEmail     string
	CommitMessage string
}

// GetGitURL returns the URL of a git repo
func GetGitURL(baseURL, m string) string {
	return baseURL + "/" + m
}

// PushModule adds all files to a git repo, commits and finally pushes them to a remote
func (metadata *Git) PushModule(module, pushDir string) error {
	g, err := git.PlainOpen(filepath.Join(pushDir))
	if err != nil {
		return err
	}

	_, err = g.CreateRemote(&gitconf.RemoteConfig{
		Name: metadata.RemoteName,
		URLs: []string{metadata.RemoteURL},
	})
	if err != nil {
		return err
	}

	wt, err := g.Worktree()
	if err != nil {
		return err
	}

	_, err = wt.Add(".")
	if err != nil {
		return err
	}

	_, err = wt.Commit(WithTimestamp(metadata.CommitMessage), &git.CommitOptions{
		Author: &object.Signature{
			Name:  metadata.UserName,
			Email: metadata.UserEmail,
			When:  time.Now(),
		},
	})
	if err != nil {
		return err
	}

	err = g.Push(&git.PushOptions{
		RemoteName: metadata.RemoteName,
		RefSpecs:   []gitconf.RefSpec{"+refs/heads/master:refs/heads/master"},
	})
	if err != nil {
		return err
	}

	return nil
}
