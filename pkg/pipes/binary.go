package pipes

import (
	"github.com/google/uuid"
	"os"
	"path/filepath"
)

type Binary struct {
	Build       Build
	PathInImage string
	DistPath    string
	CleanGlob   string
}

const EmptyRunCommand = "echo"

func (binary *Binary) GetBinaryFromDockerImage(platform string) (string, error) {
	id := uuid.New().String()
	distPath, _ := filepath.Split(binary.DistPath)
	if err := os.MkdirAll(distPath, 0777); err != nil {
		return "", err
	}

	output, err := binary.Build.execDocker(platform, "ps", "-aqf", "name="+id)
	if err != nil {
		return output, err
	}
	if output != "\n" {
		if output, err := binary.Build.execDocker(platform, "run", "--platform", platform, "-e", "TARGETPLATFORM="+platform, "--name", id, binary.Build.Tag, EmptyRunCommand); err != nil {
			return output, err
		}

		if output, err := binary.Build.execDocker(platform, "cp", id+":"+binary.PathInImage, binary.DistPath); err != nil {
			return output, err
		}

		if output, err := binary.Build.execDocker(platform, "rm", "-f", id); err != nil {
			return output, err
		}
	} else {
		return binary.GetBinaryFromDockerImage(platform)
	}

	return "", nil
}

func (binary *Binary) Clean() error {
	filesToRemove, _ := filepath.Glob(binary.CleanGlob)

	for _, fileToRemove := range filesToRemove {
		if err := os.Remove(fileToRemove); err != nil {
			return err
		}
	}

	return nil
}
