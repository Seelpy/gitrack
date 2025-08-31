package service

import (
	"context"
	"errors"
)

var (
	ErrGitRepositoryNotFound = errors.New("repository not found")
)

type Git interface {
	GetBranch(ctx context.Context) (string, error)
	GetRepository(ctx context.Context) (string, error)
	Commit(ctx context.Context, message string) error
	Merge(ctx context.Context, from string, to string) error
}
