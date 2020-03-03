package utils

// DockerManager manages Docker
type DockerManager struct {
	stdoutChan, stderrChan chan string
}

// NewDockerManager creates a new DockerManager
func NewDockerManager(stdoutChan, stderrChan chan string) *DockerManager {
	return &DockerManager{
		stdoutChan: stdoutChan,
		stderrChan: stderrChan,
	}
}

// Build builds and tags a Docker image
func (d *DockerManager) Build(file, context, tag string) error {
	command := NewManageableCommand("docker build -f "+file+" -t "+tag+" "+context, d.stdoutChan, d.stderrChan)

	if err := command.Start(); err != nil {
		return err
	}

	return command.Wait()
}
