package app

import "gitrack/pkg/app/service"

func NewProvider(gitrack service.Gitrack) Provider {
	return Provider{
		Gitrack: gitrack,
	}
}

type Provider struct {
	Gitrack service.Gitrack
}
