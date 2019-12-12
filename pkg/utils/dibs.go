package utils

type Dibs struct {
	Manifest  Manifest
	Platforms []Platform
}

func (dibs *Dibs) BuildDockerManifest(imageTags []string) error {
	for _, platform := range dibs.Platforms {
		if err := platform.Binary.Build.execDocker("manifest", "create", "--amend", platform.Binary.Build.Tag, dibs.Manifest.Tag); err != nil {
			return err
		}
	}

	return nil
}

func (dibs *Dibs) PushDockerManifest() error {
	return dibs.Platforms[0].Binary.Build.execDocker("manifest", "push", dibs.Manifest.Tag)
}
