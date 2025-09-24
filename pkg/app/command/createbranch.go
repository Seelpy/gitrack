package command

import (
	"context"
	stderrors "errors"
	"github.com/pkg/errors"
	"gitrack/pkg/app/service"
)

var (
	ErrIssueIDIsRequired = stderrors.New("<issueID> is required params")
)

func newCreateBranch(gitrack service.Gitrack) Command {
	return &createBranch{
		gitrack: gitrack,
	}
}

type createBranch struct {
	gitrack service.Gitrack
}

func (c *createBranch) Name() string {
	return "cb"
}

func (c *createBranch) Help() string {
	return "cb <issueID>"
}

func (c *createBranch) Description() string {
	return "Create branch from youtrack issue"
}

func (c *createBranch) Run(args []string) error {
	ctx := context.Background()
	if len(args) == 0 {
		return errors.WithStack(ErrIssueIDIsRequired)
	}

	err := c.gitrack.CreateBranch(ctx, args[0])
	if err != nil {
		return err
	}

	return nil
}
