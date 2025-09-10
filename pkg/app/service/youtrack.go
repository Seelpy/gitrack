package service

import (
	"context"
	"errors"
)

var (
	ErrIssueNotFound = errors.New("issue not found")
)

type Youtrack interface {
	GetIssue(ctx context.Context, issueID string) (Issue, error)
}

type Issue struct {
	ID          string
	Title       string
	Description string
	Tags        []string
}
