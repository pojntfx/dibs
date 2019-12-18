package pipes

import (
	"errors"
	"github.com/pojntfx/dibs/pkg/utils"
	"os"
	"path/filepath"
)

// Dibs is a distributed build system
type Dibs struct {
	Manifest  Manifest
	Chart     Chart
	Platforms []Platform
}

// BuildDockerManifest builds the Docker manifest
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

// PushDockerManifest pushes the Docker manifest
func (dibs *Dibs) PushDockerManifest(platform string) (string, error) {
	return dibs.Platforms[0].Assets.Build.execDocker(platform, "manifest", "push", "--purge", dibs.Manifest.Tag)
}

// BuildHelmChart builds the Helm chart
func (dibs *Dibs) BuildHelmChart(platform string) (string, error) {
	if err := os.MkdirAll(dibs.Chart.DistDir, 0777); err != nil {
		return "", err
	}

	if output, err := dibs.Platforms[0].Assets.Build.execHelm(platform, "dep", "up", dibs.Chart.SrcDir); err != nil {
		return output, err
	}

	return dibs.Platforms[0].Assets.Build.execHelm(platform, "package", "-d", dibs.Chart.DistDir, dibs.Chart.SrcDir)
}

// CleanHelmChart cleans the Helm chart
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

// GetPlatforms returns the platforms for a platform search string or all platforms if specified
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

// PushHelmChart pushes the Helm chart to GitHub releases and creates an index in a GitHub git repository for GitHub pages
func (dibs *Dibs) PushHelmChart(platform, gitUserName, gitUserEmail, gitCommitMessage, githubUsername, githubToken, githubRepositoryName, githubRepoURL, githubPagesUrl string, pushDir []string) (string, error) {
	if output, err := dibs.Platforms[0].Assets.Build.execChartReleaser(platform, "upload", "-o", githubUsername, "-t", githubToken, "-r", githubRepositoryName, "-p", dibs.Chart.DistDir); err != nil {
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

	output, err := dibs.Platforms[0].Assets.Build.execChartReleaser(platform, "index", "-o", githubUsername, "-t", githubToken, "-r", githubRepositoryName, "-p", dibs.Chart.DistDir, "-i", filepath.Join(append(pushDir, "index.yaml")...), "-c", githubPagesUrl)
	if err != nil {
		return output, err
	}

	return output, git.AddCommitAndPush()
}
