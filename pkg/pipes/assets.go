package pipes

import (
	"github.com/google/uuid"
	"os"
	"path/filepath"
)

// Assets are non-image, non-manifest and non-chart artifacts
type Assets struct {
	Build       Build
	PathInImage string   // The path of the asset in the Docker image
	DistPath    string   // The path to which the asset from the Docker image should be put
	CleanGlobs  []string // Array of globs to clean
}

// Command to execute when running the container to get the image; should `exit 0`
const EmptyRunCommand = "echo"

// GetAssetsFromDockerImage gets the assets from a Docker image
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

// Clean removes the assets locally
func (assets *Assets) Clean() error {
	for _, glob := range assets.CleanGlobs {
		filesToRemove, _ := filepath.Glob(glob)

		for _, fileToRemove := range filesToRemove {
			if err := os.RemoveAll(fileToRemove); err != nil {
				return err
			}
		}
	}

	return nil
}

// Push pushes the assets to GitHub releases
func (assets *Assets) Push(platform string, version []string, token string) (string, error) {
	return assets.Build.execGHR(platform, append([]string{"-replace", "-t", token}, append(version, assets.DistPath)...)...)
}
