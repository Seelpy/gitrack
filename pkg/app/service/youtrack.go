package service

import (
	"context"
	"errors"
)

type IssueState int

const (
	IssueStateCodeReview = IssueState(iota)
	IssueStateOther      = IssueState(iota)
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
	State       IssueState
	Tags        []string
}
