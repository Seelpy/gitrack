package app

import (
	"github.com/mitchellh/cli"
	"gitrack/pkg/app/command"
	infracommand "gitrack/pkg/inferastructure/cli/command"
	"os"
)

func New(appName string, version string) *App {
	c := cli.NewCLI(appName, version)
	c.Args = os.Args[1:]
	c.Commands = make(map[string]cli.CommandFactory)

	return &App{
		cli: c,
	}
}

type App struct {
	cli *cli.CLI
}

func (a *App) Register(command command.Command) {
	a.cli.Commands[command.Name()] = infracommand.NewCommandFactory(command)
}

func (a *App) Run() (int, error) {
	return a.cli.Run()
}
