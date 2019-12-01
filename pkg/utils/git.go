package utils

import (
	redis "github.com/go-redis/redis/v7"
	git "gopkg.in/src-d/go-git.v4"
	gitconf "gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"path/filepath"
	"time"
)

// GetGitURL returns the URL of a git repo
func GetGitURL(baseURL, m string) string {
	return baseURL + "/" + m
}

// PushModule adds all files to a git repo, commits and finally pushes them to a remote
func PushModule(r *redis.Client, prefix, suffix, m, pushDir, gitRemoteName, pushURL, gitName, gitEmail, gitCommitMessage string) error {
	g, err := git.PlainOpen(filepath.Join(pushDir))
	if err != nil {
		return err
	}

	_, err = g.CreateRemote(&gitconf.RemoteConfig{
		Name: gitRemoteName,
		URLs: []string{pushURL},
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

	_, err = wt.Commit(WithTimestamp(gitCommitMessage), &git.CommitOptions{
		Author: &object.Signature{
			Name:  gitName,
			Email: gitEmail,
			When:  time.Now(),
		},
	})
	if err != nil {
		return err
	}

	err = g.Push(&git.PushOptions{
		RemoteName: gitRemoteName,
		RefSpecs:   []gitconf.RefSpec{"+refs/heads/master:refs/heads/master"},
	})
	if err != nil {
		return err
	}

	r.Publish(prefix+":"+suffix, WithTimestamp(m))

	return nil
}
