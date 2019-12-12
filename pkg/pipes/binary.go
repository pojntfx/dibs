package pipes

import (
	"github.com/google/uuid"
	"os"
	"os/exec"
	"path/filepath"
)

type Binary struct {
	Build       Build
	PathInImage string
	DistPath    string
	CleanGlob   string
}

const EmptyRunCommand = "echo"

func (binary *Binary) GetBinaryFromDockerImage(platform string) error {
	id := uuid.New().String()
	distPath, _ := filepath.Split(binary.DistPath)
	if err := os.MkdirAll(distPath, 0777); err != nil {
		return err
	}

	out, err := exec.Command("docker", "ps", "-aqf", "name="+id).Output()
	if err != nil {
		return err
	}
	if string(out) != "\n" {
		if err := binary.Build.execDocker("run", "--platform", platform, "--name", id, binary.Build.Tag, EmptyRunCommand); err != nil {
			return err
		}

		if err := binary.Build.execDocker("cp", id+":"+binary.PathInImage, binary.DistPath); err != nil {
			return err
		}

		if err := binary.Build.execDocker("rm", "-f", id); err != nil {
			return err
		}
	} else {
		return binary.GetBinaryFromDockerImage(platform)
	}

	return nil
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
