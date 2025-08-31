package main

import (
	"log"
	"os"

	"gitrack/pkg/app"
	"gitrack/pkg/app/command"
	infraapp "gitrack/pkg/inferastructure/cli/app"
)

func main() {
	cliApp := infraapp.New("test", "1.0.0")

	provider := app.NewProvider()
	command.RegisterCommands(cliApp, provider)

	exitStatus, err := cliApp.Run()
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(exitStatus)
}
