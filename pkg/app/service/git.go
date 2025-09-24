package service

import (
	"errors"
)

var (
	ErrGitRepositoryNotFound = errors.New("repository not found")
	ErrBranchNotFound        = errors.New("branch not found")
	ErrNoChangesToCommit     = errors.New("no changes to commit")
)

type Git interface {
	GetBranch() (string, error)
	GetRepository() (string, error)
	Commit(message string) error
	Merge(from string, to string) error
	CreateBranch(from string, name string) error
}
