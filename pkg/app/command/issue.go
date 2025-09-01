package command

import (
	"context"
	"gitrack/pkg/app/service"
	"log"
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
	return "issue"
}

func (c *branch) Help() string {
	return "Get issue information"
}

func (c *branch) Description() string {
	return "Get issue information"
}

func (c *branch) Run(_ []string) error {
	ctx := context.Background()
	issue, err := c.gitrack.GetBranchInfo(ctx)
	if err != nil {
		return err
	}

	log.Printf("Issue ID: %s", issue.ID)
	log.Printf("Title: %s", issue.Title)
	log.Printf("Description: %s", issue.Description)
	log.Printf("State: %v", issue.State)
	log.Printf("Tags: %v", issue.Tags)

	return nil
}
