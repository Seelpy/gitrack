package command

import (
	"errors"
	"fmt"
)

func newLogger(command Command) Command {
	return &logger{
		next: command,
	}
}

type logger struct {
	next Command
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
		fmt.Println(errors.Unwrap(err))
	}
	return nil
}
