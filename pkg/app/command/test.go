package command

import (
	"log"
)

func newTest() Command {
	return &test{}
}

type test struct {
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

func (c *test) Run(args []string) error {
	log.Println(args)
	return nil
}
