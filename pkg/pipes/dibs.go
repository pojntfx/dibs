package pipes

import "errors"

type Dibs struct {
	Manifest  Manifest
	Platforms []Platform
}

func (dibs *Dibs) BuildDockerManifest() (string, error) {
	for _, platform := range dibs.Platforms {
		if output, err := platform.Binary.Build.execDocker("manifest", "create", "--amend", platform.Binary.Build.Tag, dibs.Manifest.Tag); err != nil {
			return output, err
		}
	}

	return "", nil
}

func (dibs *Dibs) PushDockerManifest() (string, error) {
	return dibs.Platforms[0].Binary.Build.execDocker("manifest", "push", dibs.Manifest.Tag)
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
