package utils

import (
	"gopkg.in/src-d/go-git.v4"
	gitConfiguration "gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	"path/filepath"
	"time"
)

// Git is the configuration for the Git version control system
type Git struct {
	RemoteName    string // Name that the sync remote should use (PushToRemote only)
	RemoteURL     string // Base URL of the sync remote (PushToRemote only)
	UserName      string // Name to use for commits
	UserEmail     string // Email to use for commits
	CommitMessage string // Message to use for commits
	WorkDir       string // Directory in which to work
	Token         string // Token to use to pull and push (HTTP basic auth)
}

// GetGitURL returns the URL of a git repo
func GetGitURL(baseURL, module string) string {
	completeURL := baseURL + "/" + module

	return completeURL
}

// PushToRemote adds all files to a git repo, commits and pushes them to a remote
func (metaGit *Git) PushToRemote(pushDir string) error {
	g, err := git.PlainOpen(filepath.Join(pushDir))
	if err != nil {
		return err
	}

	if _, err = g.CreateRemote(&gitConfiguration.RemoteConfig{
		Name: metaGit.RemoteName,
		URLs: []string{metaGit.RemoteURL},
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

	if _, err = wt.Commit(metaGit.CommitMessage, &git.CommitOptions{
		Author: &object.Signature{
			Name:  metaGit.UserName,
			Email: metaGit.UserEmail,
			When:  time.Now(),
		},
	}); err != nil {
		return err
	}

	if err = g.Push(&git.PushOptions{
		RemoteName: metaGit.RemoteName,
		RefSpecs:   []gitConfiguration.RefSpec{"+refs/heads/master:refs/heads/master"},
	}); err != nil {
		return err
	}

	return nil
}

// Clone clones a git repository
func (metaGit Git) Clone(url string) error {
	_, err := git.PlainClone(metaGit.WorkDir, false, &git.CloneOptions{
		URL:      url,
		Auth:     &http.BasicAuth{Username: metaGit.UserName, Password: metaGit.Token},
		Progress: nil,
	})

	return err
}

// GetLatestTag returns the latest tag from a Git repository
func (metaGit Git) GetLatestTag() (string, error) {
	// Based on https://github.com/src-d/go-git/issues/1030#issuecomment-443679681
	repository, err := git.PlainOpen(filepath.Join(metaGit.WorkDir))

	if repository != nil {
		tagRefs, err := repository.Tags()
		if err != nil {
			return "", err
		}

		var latestTagCommit *object.Commit
		var latestTagName string
		if err := tagRefs.ForEach(func(tagRef *plumbing.Reference) error {
			revision := plumbing.Revision(tagRef.Name().String())
			tagCommitHash, err := repository.ResolveRevision(revision)
			if err != nil {
				return err
			}

			commit, err := repository.CommitObject(*tagCommitHash)
			if err != nil {
				return err
			}

			if latestTagCommit == nil {
				latestTagCommit = commit
				latestTagName = tagRef.Name().Short()
			}

			if commit.Committer.When.After(latestTagCommit.Committer.When) {
				latestTagCommit = commit
				latestTagName = tagRef.Name().Short()
			}

			return nil
		}); err != nil {
			return "", err
		}

		return latestTagName, nil
	}

	return "", err
}

// AddCommitAndPush stages all files in a git repository, then commits and pushes them
func (metaGit Git) AddCommitAndPush() error {
	g, err := git.PlainOpen(filepath.Join(metaGit.WorkDir))
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

	if _, err = wt.Commit(metaGit.CommitMessage, &git.CommitOptions{
		Author: &object.Signature{
			Name:  metaGit.UserName,
			Email: metaGit.UserEmail,
			When:  time.Now(),
		},
	}); err != nil {
		return err
	}

	if err = g.Push(&git.PushOptions{
		Auth: &http.BasicAuth{Username: metaGit.UserName, Password: metaGit.Token},
	}); err != nil {
		return err
	}

	return nil
}
