package log

import (
	"errors"
	"gitrack/pkg/app/command"
	"log"
)

func NewLogger(command command.Command) command.Command {
	return &logger{
		next: command,
	}
}

type logger struct {
	next command.Command
}

func (c *logger) Name() string {
	return c.next.Name()
}

func (c *logger) Help() string {
	return c.next.Help()
}

func (c *logger) Description() string {
	return c.next.Description()
}

func (c *logger) Run(args []string) error {
	err := c.next.Run(args)
	if err != nil {
		log.Println(errors.Unwrap(err))
	}
	return nil
}
