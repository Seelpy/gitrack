package main

import (
	"errors"
	config2 "gitrack/data/config"
	"gitrack/pkg/app/service"
	"gitrack/pkg/inferastructure/git"
	"gitrack/pkg/inferastructure/yt"
	"log"
	"os"

	"gitrack/pkg/app"
	"gitrack/pkg/app/command"
	infraapp "gitrack/pkg/inferastructure/cli/app"
)

func main() {
	cliApp := infraapp.New("test", "1.0.0")

	gitService := git.NewService()
	ytService := yt.NewService("http://localhost:9090/", "perm-YWRtaW4=.NDMtMA==.TaQEjpQ5wLs19NJDhTlKjSfSbkQwws")

	path, err := config2.GetConfigPath(command.ConfigPath)
	if err != nil && !errors.Is(err, config2.ErrPathNotExist) {
		log.Fatal(err)
	}
	config, err := config2.ParseConfig(path)
	if err != nil {
		log.Fatal(err)
	}

	gitrack := service.NewGitrack(gitService, ytService, config)
	provider := app.NewProvider(gitrack)
	command.RegisterCommands(cliApp, provider)

	exitStatus, err := cliApp.Run()
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(exitStatus)
}
