package command

import (
	"fmt"
	"gitrack/data/config"
)

func newGetConfigPath() Command {
	return &getConfigPath{}
}

type getConfigPath struct {
}

func (c *getConfigPath) Name() string {
	return "get-config-path"
}

func (c *getConfigPath) Help() string {
	return "get-config-path"
}

func (c *getConfigPath) Description() string {
	return "get config path"
}

func (c *getConfigPath) Run(_ []string) error {
	path, err := config.GetConfigPath(ConfigPath)
	if err != nil {
		return err
	}
	fmt.Printf("Config paht: %s\n", path)
	return nil
}
