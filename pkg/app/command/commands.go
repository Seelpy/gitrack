package command

import (
	"gitrack/pkg/app"
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
		func(provider app.Provider) Command { return newLogger(newTest()) },
		func(provider app.Provider) Command { return newLogger(newBranch(provider.Gitrack)) },
		func(provider app.Provider) Command { return newLogger(newSetConfigPath()) },
		func(provider app.Provider) Command { return newLogger(newGetConfigPath()) },
		func(provider app.Provider) Command { return newLogger(newMerge(provider.Gitrack)) },
	)
}
