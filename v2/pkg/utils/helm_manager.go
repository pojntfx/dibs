package utils

// HelmManager manages Helm
type HelmManager struct {
	dir                    string
	stdoutChan, stderrChan chan string
}

// NewHelmManager creates a new HelmManager
func NewHelmManager(dir string, stdoutChan, stderrChan chan string) *HelmManager {
	return &HelmManager{
		dir:        dir,
		stdoutChan: stdoutChan,
		stderrChan: stderrChan,
	}
}

// Build builds a Helm chart
func (h *HelmManager) Build(src, dist string) error {
	depUpCommand := NewManageableCommand("helm dep up "+src, h.dir, h.stdoutChan, h.stderrChan)

	if err := depUpCommand.Start(); err != nil {
		return err
	}

	if err := depUpCommand.Wait(); err != nil {
		return err
	}

	buildCommand := NewManageableCommand("helm package -d "+dist+" "+src, h.dir, h.stdoutChan, h.stderrChan)

	if err := buildCommand.Start(); err != nil {
		return err
	}

	return buildCommand.Wait()
}
