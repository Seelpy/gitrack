package command

import (
	"context"
	"fmt"
	"gitrack/pkg/app/service"
)

func newBranch(gitrack service.Gitrack) Command {
	return &branch{
		gitrack: gitrack,
	}
}

type branch struct {
	gitrack service.Gitrack
}

func (c *branch) Name() string {
	return "b"
}

func (c *branch) Help() string {
	return "b"
}

func (c *branch) Description() string {
	return "Get issue information from youtrack issue"
}

func (c *branch) Run(_ []string) error {
	ctx := context.Background()
	issue, err := c.gitrack.GetBranchInfo(ctx)
	if err != nil {
		return err
	}

	fmt.Printf("Issue ID: 	%s\n", issue.ID)
	fmt.Printf("Title: 		%s\n", issue.Title)
	fmt.Printf("Description: 	%s\n", issue.Description)
	fmt.Printf("Tags: 		%s\n", issue.Tags)

	return nil
}
