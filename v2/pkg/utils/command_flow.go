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

func (f *CommandFlow) recreateCommands() error {
	newCommands := []*ManageableCommand{}

	for _, command := range f.commands {
		manageableCommand := NewManageableCommand(command.GetExecLine(), command.GetStdoutChan(), command.GetStderrChan())

		newCommands = append(newCommands, manageableCommand)
	}

	for _, command := range newCommands {
		if err := command.Start(); err != nil {
			return err
		}
	}

	f.commands = newCommands

	return nil
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
	if f.isRestart {
		return f.Wait()
	}

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

// Restart restarts the flow
func (f *CommandFlow) Restart() error {
	f.isRestart = true

	if err := f.Stop(); err != nil {
		return err
	}

	if err := f.recreateCommands(); err != nil {
		return err
	}

	f.isRestart = false

	return nil
}
