package pipes

import "errors"

type Dibs struct {
	Manifest  Manifest
	Platforms []Platform
}

func (dibs *Dibs) BuildDockerManifest(platform string) (string, error) {
	var manifestsToAdd []string

	for _, platform := range dibs.Platforms {
		manifestsToAdd = append(manifestsToAdd, platform.Binary.Build.Tag)
	}

	if output, err := dibs.Platforms[0].Binary.Build.execDocker(platform, append([]string{"manifest", "create", "--amend", dibs.Manifest.Tag}, manifestsToAdd...)...); err != nil {
		return output, err
	}

	return "", nil
}

func (dibs *Dibs) PushDockerManifest(platform string) (string, error) {
	return dibs.Platforms[0].Binary.Build.execDocker(platform, "manifest", "push", dibs.Manifest.Tag)
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
