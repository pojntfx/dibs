package utils

type Manifest struct {
	Tag string
}

func (dibs *Dibs) BuildDockerManifest() error {
	for _, platform := range dibs.Platforms {
		if err := platform.execDocker("manifest", "create", "--amend", platform.Image.Tag, dibs.Manifest.Tag); err != nil {
			return err
		}
	}

	return nil
}

func (dibs *Dibs) PushDockerManifest() error {
	return dibs.Platforms[0].execDocker("manifest", "push", dibs.Manifest.Tag)
}
