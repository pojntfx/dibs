package utils

// CommandFlow is a manageable collection of commands
type CommandFlow struct {
	isRestart bool
	commands  []*ManageableCommand
}

// NewCommandFlow creates a new CommandFlow
func NewCommandFlow(commands []string, stdoutChan, stderrChan chan string) *CommandFlow {
	commandFlow := &CommandFlow{
		isRestart: false,
	}

	for _, command := range commands {
		manageableCommand := NewManageableCommand(command, stdoutChan, stderrChan)

		commandFlow.commands = append(commandFlow.commands, manageableCommand)
	}

	return commandFlow
}

// Start starts the command flow
func (f *CommandFlow) Start() error {
	for _, command := range f.commands {
		if err := command.Start(); err != nil {
			return err
		}
	}

	return nil
}

// Wait waits for the command flow to complete
func (f *CommandFlow) Wait() error {
	for _, command := range f.commands {
		if err := command.Wait(); err != nil {
			return err
		}
	}

	if f.isRestart {
		return f.Wait()
	}

	return nil
}

// Stop stops the flow
func (f *CommandFlow) Stop() error {
	for i := len(f.commands) - 1; i >= 0; i-- {
		command := f.commands[i]

		if !command.IsStopped() {
			if err := command.Stop(); err != nil {
				return err
			}
		}
	}

	return nil
}
