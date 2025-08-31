package command

import (
	"github.com/mitchellh/cli"
	"gitrack/pkg/app/command"
)

func NewCommandFactory(command command.Command) cli.CommandFactory {
	return func() (cli.Command, error) {
		return newAdapter(command), nil
	}
}

func newAdapter(command command.Command) cli.Command {
	return &adapter{
		command: command,
	}
}

type adapter struct {
	command command.Command
}

func (a *adapter) Help() string {
	return a.command.Help()
}

func (a *adapter) Synopsis() string {
	return a.command.Description()
}

func (a *adapter) Run(args []string) int {
	if a.command.Run(args) == nil {
		return 0
	}
	return 1
}
