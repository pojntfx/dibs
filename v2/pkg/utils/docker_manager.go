package utils

// DockerManager manages Docker
type DockerManager struct {
	dir                    string
	stdoutChan, stderrChan chan string
}

// NewDockerManager creates a new DockerManager
func NewDockerManager(dir string, stdoutChan, stderrChan chan string) *DockerManager {
	return &DockerManager{
		dir:        dir,
		stdoutChan: stdoutChan,
		stderrChan: stderrChan,
	}
}

// Build builds and tags a Docker image
func (d *DockerManager) Build(file, context, tag string) error {
	command := NewManageableCommand("docker build -f "+file+" -t "+tag+" "+context, d.dir, d.stdoutChan, d.stderrChan)

	if err := command.Start(); err != nil {
		return err
	}

	return command.Wait()
}

// Push pushes a Docker image
func (d *DockerManager) Push(tag string) error {
	command := NewManageableCommand("docker push "+tag, d.dir, d.stdoutChan, d.stderrChan)

	if err := command.Start(); err != nil {
		return err
	}

	return command.Wait()
}

// Run runs a command in a Docker image
func (d *DockerManager) Run(tag, execLine string) error {
	command := NewManageableCommand("docker run "+tag+" "+execLine, d.dir, d.stdoutChan, d.stderrChan)

	if err := command.Start(); err != nil {
		return err
	}

	return command.Wait()
}
