package app

type Provider interface {
}

func NewProvider() Provider {
	return &provider{}
}

type provider struct {
}
