package pipes

import (
	"errors"
	"github.com/pojntfx/dibs/pkg/utils"
	"os"
	"path/filepath"
)

type Dibs struct {
	Manifest  Manifest
	Chart     Chart
	Platforms []Platform
}

func (dibs *Dibs) BuildDockerManifest(platform string) (string, error) {
	var manifestsToAdd []string

	for _, platform := range dibs.Platforms {
		manifestsToAdd = append(manifestsToAdd, platform.Assets.Build.Tag)
	}

	if output, err := dibs.Platforms[0].Assets.Build.execDocker(platform, append([]string{"manifest", "create", dibs.Manifest.Tag}, manifestsToAdd...)...); err != nil {
		return output, err
	}

	return "", nil
}

func (dibs *Dibs) PushDockerManifest(platform string) (string, error) {
	return dibs.Platforms[0].Assets.Build.execDocker(platform, "manifest", "push", "--purge", dibs.Manifest.Tag)
}

func (dibs *Dibs) BuildHelmChart(platform string) (string, error) {
	if err := os.MkdirAll(dibs.Chart.DistDir, 0777); err != nil {
		return "", err
	}

	if output, err := dibs.Platforms[0].Assets.Build.execHelm(platform, "dep", "up", dibs.Chart.SrcDir); err != nil {
		return output, err
	}

	return dibs.Platforms[0].Assets.Build.execHelm(platform, "package", "-d", dibs.Chart.DistDir, dibs.Chart.SrcDir)
}

func (dibs *Dibs) CleanHelmChart() error {
	for _, glob := range dibs.Chart.CleanGlobs {
		filesToRemove, _ := filepath.Glob(glob)

		for _, fileToRemove := range filesToRemove {
			if err := os.RemoveAll(fileToRemove); err != nil {
				return err
			}
		}
	}

	return nil
}

func (dibs *Dibs) GetPlatforms(wantedPlatform string, all bool) ([]Platform, error) {
	if all {
		return dibs.Platforms, nil
	} else {
		for _, platform := range dibs.Platforms {
			if platform.Platform == wantedPlatform {
				return []Platform{platform}, nil
			}
		}
	}

	return []Platform{}, errors.New("platform not found")
}

func (dibs *Dibs) PushHelmChart(platform, gitUserName, gitUserEmail, gitCommitMessage, githubUsername, githubToken, githubRepoName, githubRepoURL, githubPagesUrl string, pushDir []string) (string, error) {
	if output, err := dibs.Platforms[0].Assets.Build.execCR(platform, "upload", "-o", githubUsername, "-t", githubToken, "-r", githubRepoName, "-p", dibs.Chart.DistDir); err != nil {
		return output, err
	}

	git := utils.Git{
		UserName:      gitUserName,
		UserEmail:     gitUserEmail,
		CommitMessage: gitCommitMessage,
		WorkDir:       filepath.Join(pushDir...),
		Token:         githubToken,
	}

	if err := git.Clone(githubRepoURL); err != nil {
		return "", err
	}

	output, err := dibs.Platforms[0].Assets.Build.execCR(platform, "index", "-o", githubUsername, "-t", githubToken, "-r", githubRepoName, "-p", dibs.Chart.DistDir, "-i", filepath.Join(append(pushDir, "index.yaml")...), "-c", githubPagesUrl)
	if err != nil {
		return output, err
	}

	return output, git.AddCommitAndPush()
}
