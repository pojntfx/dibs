package utils

// CommandFlow is a manageable collection of commands
type CommandFlow struct {
	isRestart bool
	commands  []*ManageableCommand
}

// NewCommandFlow creates a new CommandFlow
func NewCommandFlow(commands []string, stdoutChan, stderrChan chan string) *CommandFlow {
	commandFlow := &CommandFlow{}

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
