package service

import (
	"context"
	stderrors "errors"
	"fmt"
	"github.com/pkg/errors"
)

const (
	commitMessageTemplate = "%s %s"
)

var (
	ErrFeatureConfigNotFound = stderrors.New("feature config not found")
)

type GitrackConfigProvider interface {
	GetFeatureConfig(repository string, tag string) (FeatureConfig, error)
}

type Features []FeatureConfig

type FeatureConfig struct {
	FeatureTag string
	Releases   []ReleaseConfig
}

type ReleaseConfig struct {
	ReleaseBranch string
	YoutrackTag   string
}

func NewGitrack(git Git, youtrack Youtrack, configProvider GitrackConfigProvider) Gitrack {
	return Gitrack{
		git:      git,
		youtrack: youtrack,
		config:   configProvider,
	}
}

type Gitrack struct {
	git      Git
	youtrack Youtrack
	config   GitrackConfigProvider
}

type BranchInfo struct {
}

func (g *Gitrack) GetBranchInfo(ctx context.Context) (Issue, error) {
	branch, err := g.git.GetBranch(ctx)
	if err != nil {
		return Issue{}, err
	}
	issue, err := g.youtrack.GetIssue(ctx, branch)
	if err != nil {
		return Issue{}, err
	}
	return issue, nil
}

func (g *Gitrack) Commit(ctx context.Context, message string) error {
	branch, err := g.git.GetBranch(ctx)
	if err != nil {
		return err
	}
	return g.git.Commit(ctx, fmt.Sprintf(commitMessageTemplate, branch, message))
}

func (g *Gitrack) Merge(ctx context.Context) error {
	branch, err := g.git.GetBranch(ctx)
	if err != nil {
		return err
	}

	issue, err := g.youtrack.GetIssue(ctx, branch)
	if err != nil {
		return err
	}

	mergePipeline, err := g.getMergePipeline(ctx, issue)
	if err != nil {
		return err
	}

	for i := 1; i < len(mergePipeline); i++ {
		err = g.git.Merge(ctx, mergePipeline[i-1], mergePipeline[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *Gitrack) getMergePipeline(ctx context.Context, issue Issue) ([]string, error) {
	repository, err := g.git.GetRepository(ctx)
	if err != nil {
		return nil, err
	}

	var config FeatureConfig
	for _, tag := range issue.Tags {
		config, err = g.config.GetFeatureConfig(repository, tag)
		if errors.Is(err, ErrFeatureConfigNotFound) {
			continue
		}
		if err != nil {
			return nil, err
		}
		break
	}

	releasePipeline := make([]string, 0)
	for _, tag := range issue.Tags {
		releasePipeline = g.getMergePipelineByTag(tag, config.Releases)
		if len(releasePipeline) > 0 {
			return append([]string{issue.ID}, releasePipeline...), nil
		}
	}
	return releasePipeline, nil
}

func (g *Gitrack) getMergePipelineByTag(issueTag string, releases []ReleaseConfig) []string {
	pipeline := make([]string, 0)
	find := false
	for _, release := range releases {
		if release.YoutrackTag == issueTag {
			find = true
		}
		if find {
			pipeline = append(pipeline, release.ReleaseBranch)
		}
	}
	return pipeline
}
