package command

import (
	"gitrack/pkg/app"
	"log"
)

func init() {
	commands = append(
		commands,
		func(provider app.Provider) Command {
			return newTest()
		},
	)
}

func newTest() Command {
	return &test{
		subCommands: make([]Command, 0),
	}
}

type test struct {
	subCommands []Command
}

func (c *test) Name() string {
	return "test"
}

func (c *test) Help() string {
	return "test help"
}

func (c *test) Description() string {
	return "test description"
}

func (c *test) SubCommands() []Command {
	return nil
}

func (c *test) Run(args []string) error {
	log.Println(args)
	return nil
}
