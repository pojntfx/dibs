package utils

import (
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

// BinaryManager manages binaries
type BinaryManager struct {
	dir                    string
	stdoutChan, stderrChan chan string
}

// NewBinaryManager creates a new BinaryManager
func NewBinaryManager(dir string, stdoutChan, stderrChan chan string) *BinaryManager {
	return &BinaryManager{
		dir:        dir,
		stdoutChan: stdoutChan,
		stderrChan: stderrChan,
	}
}

func getLatestGitTag(dir string) (string, error) {
	// Based on https://github.com/src-d/go-git/issues/1030#issuecomment-443679681
	repository, err := git.PlainOpen(dir)

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

// Push releases a binary to GitHub releases
func (b *BinaryManager) Push(githubUserName, githubToken, githubRepository, dir, assetOut string) error {
	version, err := getLatestGitTag(dir)
	if err != nil {
		return err
	}

	command := NewManageableCommand("ghr -replace -t "+githubToken+" -u "+githubUserName+" -r "+githubRepository+" "+version+" "+assetOut, b.dir, b.stdoutChan, b.stderrChan)

	if err := command.Start(); err != nil {
		return err
	}

	return command.Wait()
}
