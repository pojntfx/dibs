package utils

import (
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	"os"
	"path/filepath"
	"time"
)

// HelmManager manages Helm
type HelmManager struct {
	dir                    string
	stdoutChan, stderrChan chan string
}

// NewHelmManager creates a new HelmManager
func NewHelmManager(dir string, stdoutChan, stderrChan chan string) *HelmManager {
	return &HelmManager{
		dir:        dir,
		stdoutChan: stdoutChan,
		stderrChan: stderrChan,
	}
}

// Build builds a Helm chart
func (h *HelmManager) Build(src, dist string) error {
	depUpCommand := NewManageableCommand("helm dep up "+src, h.dir, h.stdoutChan, h.stderrChan)

	if err := depUpCommand.Start(); err != nil {
		return err
	}

	if err := depUpCommand.Wait(); err != nil {
		return err
	}

	buildCommand := NewManageableCommand("helm package -d "+dist+" "+src, h.dir, h.stdoutChan, h.stderrChan)

	if err := buildCommand.Start(); err != nil {
		return err
	}

	return buildCommand.Wait()
}

// Push releases a Helm chart using GitHub, GitHub releases and GitHub pages
func (h *HelmManager) Push(gitUserName, gitUserEmail, gitCommitMessage, githubUserName, githubToken, githubRepositoryName, githubRepositoryUrl, githubPagesUrl, chartDist, cloneDir string) error {
	uploadCommand := NewManageableCommand("cr upload -o "+githubUserName+" -t "+githubToken+" -r "+githubRepositoryName+" -p "+chartDist, h.dir, h.stdoutChan, h.stderrChan)

	if err := uploadCommand.Start(); err != nil {
		return err
	}

	if err := uploadCommand.Wait(); err != nil {
		return err
	}

	if err := os.RemoveAll(cloneDir); err != nil {
		return err
	}

	if _, err := git.PlainClone(cloneDir, false, &git.CloneOptions{
		URL:      githubRepositoryUrl,
		Auth:     &http.BasicAuth{Username: githubUserName, Password: githubToken},
		Progress: nil,
	}); err != nil {
		return err
	}

	updateIndexCommand := NewManageableCommand("cr index -o "+githubUserName+" -t "+githubToken+" -r "+githubRepositoryName+" -p "+chartDist+" -i "+filepath.Join(cloneDir, "index.yaml")+" -c "+githubPagesUrl, h.dir, h.stdoutChan, h.stderrChan)

	if err := updateIndexCommand.Start(); err != nil {
		return err
	}

	if err := updateIndexCommand.Wait(); err != nil {
		return err
	}

	g, err := git.PlainOpen(cloneDir)
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

	if _, err = wt.Commit(gitCommitMessage, &git.CommitOptions{
		Author: &object.Signature{
			Name:  gitUserName,
			Email: gitUserEmail,
			When:  time.Now(),
		},
	}); err != nil {
		return err
	}

	return g.Push(&git.PushOptions{
		Auth: &http.BasicAuth{Username: githubUserName, Password: githubToken},
	})
}
