package utils

// CommandFlow is a manageable collection of commands
type CommandFlow struct {
	isRestart bool
	commands  []*ManageableCommand
}

// NewCommandFlow creates a new CommandFlow
func NewCommandFlow(commands []string, dir string, stdoutChan, stderrChan chan string) *CommandFlow {
	commandFlow := &CommandFlow{
		isRestart: false,
	}

	for _, command := range commands {
		manageableCommand := NewManageableCommand(command, dir, stdoutChan, stderrChan)

		commandFlow.commands = append(commandFlow.commands, manageableCommand)
	}

	return commandFlow
}

func (f *CommandFlow) recreateCommands() error {
	var newCommands []*ManageableCommand

	for _, command := range f.commands {
		manageableCommand := NewManageableCommand(command.GetExecLine(), command.GetDir(), command.GetStdoutChan(), command.GetStderrChan())

		newCommands = append(newCommands, manageableCommand)
	}

	for i, command := range newCommands {
		if err := command.Start(); err != nil {
			return err
		}

		// We don't have to wait for the last one to ensure serial execution
		if i != len(newCommands)-1 {
			_ = command.Wait()
		}
	}

	f.commands = newCommands

	return nil
}

// Start starts the command flow
func (f *CommandFlow) Start() error {
	// TODO: Add test that ensures serial execution of commands
	for i, command := range f.commands {
		if err := command.Start(); err != nil {
			return err
		}

		// We don't have to wait for the last one to ensure serial execution
		if i != len(f.commands)-1 {
			_ = command.Wait()
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
	// TODO: Add test that ensures serial execution of commands

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
