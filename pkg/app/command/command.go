package command

type Command interface {
	Name() string
	Help() string
	Description() string
	Run(args []string) error
}
