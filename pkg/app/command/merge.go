package command

import (
	"context"
	"gitrack/pkg/app/service"
)

func newMerge(gitrack service.Gitrack) Command {
	return &merge{
		gitrack: gitrack,
	}
}

type merge struct {
	gitrack service.Gitrack
}

func (c *merge) Name() string {
	return "merge"
}

func (c *merge) Help() string {
	return "merge"
}

func (c *merge) Description() string {
	return "Merge branch automatically"
}

func (c *merge) Run(_ []string) error {
	ctx := context.Background()
	err := c.gitrack.Merge(ctx)
	if err != nil {
		return err
	}

	return nil
}
