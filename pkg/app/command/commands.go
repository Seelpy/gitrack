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
