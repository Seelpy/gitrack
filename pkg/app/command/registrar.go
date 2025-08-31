package command

type Registrar interface {
	Register(command Command)
}
