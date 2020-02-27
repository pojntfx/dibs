package utils

// CommandFlow is a manageable collection of commands
type CommandFlow struct {
	commands []*ManageableCommand
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
