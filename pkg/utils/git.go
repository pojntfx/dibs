package utils

import (
	"gopkg.in/src-d/go-git.v4"
	gitconf "gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"path/filepath"
	"time"
)

// Git is the configuration for the Git version control system
type Git struct {
	RemoteName    string // Name that the sync remote should use
	RemoteURL     string // Base URL of the sync remote
	UserName      string // Name to use for commits
	UserEmail     string // Email to use for commits
	CommitMessage string // Message to use for commits
}

// GetGitURL returns the URL of a git repo
func GetGitURL(baseURL, module string) string {
	completeURL := baseURL + "/" + module

	return completeURL
}

// PushModule adds all files to a git repo, commits and pushes them to a remote
func (metadata *Git) PushModule(pushDir string) error {
	g, err := git.PlainOpen(filepath.Join(pushDir))
	if err != nil {
		return err
	}

	if _, err = g.CreateRemote(&gitconf.RemoteConfig{
		Name: metadata.RemoteName,
		URLs: []string{metadata.RemoteURL},
	}); err != nil {
		return err
	}

	wt, err := g.Worktree()
	if err != nil {
		return err
	}

	if _, err = wt.Add("."); err != nil {
		return err
	}

	if _, err = wt.Commit(WithTimestamp(metadata.CommitMessage), &git.CommitOptions{
		Author: &object.Signature{
			Name:  metadata.UserName,
			Email: metadata.UserEmail,
			When:  time.Now(),
		},
	}); err != nil {
		return err
	}

	if err = g.Push(&git.PushOptions{
		RemoteName: metadata.RemoteName,
		RefSpecs:   []gitconf.RefSpec{"+refs/heads/master:refs/heads/master"},
	}); err != nil {
		return err
	}

	return nil
}
