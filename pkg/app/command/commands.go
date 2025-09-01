package command

import (
	"gitrack/pkg/app"
	"gitrack/pkg/app/command/log"
)

type (
	registerCommandFunc func(provider app.Provider) Command
)

var (
	commands = make([]registerCommandFunc, 0)
)

func RegisterCommands(registrar Registrar, provider app.Provider) {
	for _, registerFunc := range commands {
		registrar.Register(registerFunc(provider))
	}
}

func init() {
	commands = append(
		commands,
		func(provider app.Provider) Command { return log.NewLogger(newTest()) },
		func(provider app.Provider) Command { return log.NewLogger(newBranch(provider.Gitrack)) },
	)
}
