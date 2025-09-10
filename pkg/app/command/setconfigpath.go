package command

import (
	stderrors "errors"
	"github.com/pkg/errors"
	"gitrack/data/config"
)

const (
	ConfigPath = "pathconfig.yaml"
)

var (
	ErrPathIsRequiredParams = stderrors.New("<path to config> is required params")
)

func newSetConfigPath() Command {
	return &setConfigPath{}
}

type setConfigPath struct {
}

func (c *setConfigPath) Name() string {
	return "set-config-path"
}

func (c *setConfigPath) Help() string {
	return "set-config-path <path to config>"
}

func (c *setConfigPath) Description() string {
	return "set config path"
}

func (c *setConfigPath) Run(args []string) error {
	if len(args) == 0 {
		return errors.WithStack(ErrPathIsRequiredParams)
	}
	return config.SetConfigPath(ConfigPath, args[0])
}
