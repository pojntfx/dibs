package pipes

import (
	"github.com/google/uuid"
	"os"
	"path/filepath"
)

type Assets struct {
	Build       Build
	PathInImage string
	DistPath    string
	CleanGlob   string
}

const EmptyRunCommand = "echo"

func (assets *Assets) GetAssetsFromDockerImage(platform string) (string, error) {
	id := uuid.New().String()
	distPath, _ := filepath.Split(assets.DistPath)
	if err := os.MkdirAll(distPath, 0777); err != nil {
		return "", err
	}

	output, err := assets.Build.execDocker(platform, "ps", "-aqf", "name="+id)
	if err != nil {
		return output, err
	}
	if output != "\n" {
		if output, err := assets.Build.execDocker(platform, "run", "--platform", platform, "-e", "TARGETPLATFORM="+platform, "--name", id, assets.Build.Tag, EmptyRunCommand); err != nil {
			return output, err
		}

		if output, err := assets.Build.execDocker(platform, "cp", id+":"+assets.PathInImage, assets.DistPath); err != nil {
			return output, err
		}

		if output, err := assets.Build.execDocker(platform, "rm", "-f", id); err != nil {
			return output, err
		}
	} else {
		return assets.GetAssetsFromDockerImage(platform)
	}

	return "", nil
}

func (assets *Assets) Clean() error {
	filesToRemove, _ := filepath.Glob(assets.CleanGlob)

	for _, fileToRemove := range filesToRemove {
		if err := os.Remove(fileToRemove); err != nil {
			return err
		}
	}

	return nil
}