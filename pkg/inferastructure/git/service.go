package git

import (
	stderrors "errors"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/pkg/errors"
	appservice "gitrack/pkg/app/service"
	"os"
	"path/filepath"
)

const (
	origin = "origin"
)

type service struct {
}

func NewService() appservice.Git {
	return &service{}
}

func (g *service) GetBranch() (string, error) {
	repo, err := g.getCurrentRepository()
	if err != nil {
		return "", err
	}

	head, err := repo.Head()
	if err != nil {
		return "", errors.WithStack(err)
	}

	if !head.Name().IsBranch() {
		return "", errors.WithStack(appservice.ErrBranchNotFound)
	}

	return head.Name().Short(), nil
}

func (g *service) GetRepository() (string, error) {
	repo, err := g.getCurrentRepository()
	if err != nil {
		return "", err
	}

	remotes, err := repo.Remotes()
	if err != nil {
		return "", errors.WithStack(err)
	}

	for _, remote := range remotes {
		if remote.Config().Name == origin && len(remote.Config().URLs) > 0 {
			url := remote.Config().URLs[0]
			return url, nil
		}
	}

	wt, err := repo.Worktree()
	if err != nil {
		return "", errors.WithStack(err)
	}

	return filepath.Base(wt.Filesystem.Root()), nil
}

func (g *service) Commit(message string) error {
	repo, err := g.getCurrentRepository()
	if err != nil {
		return err
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return errors.WithStack(err)
	}

	status, err := worktree.Status()
	if err != nil {
		return errors.WithStack(err)
	}

	if status.IsClean() {
		return errors.WithStack(appservice.ErrNoChangesToCommit)
	}

	_, err = worktree.Add(".")
	if err != nil {
		return errors.WithStack(err)
	}

	_, err = worktree.Commit(message, &git.CommitOptions{})
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (g *service) Merge(from string, to string) error {
	repo, err := g.getCurrentRepository()
	if err != nil {
		return err
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return errors.WithStack(err)
	}

	fromRef, err := repo.Reference(plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", from)), true)
	if err != nil {
		return errors.WithStack(appservice.ErrBranchNotFound)
	}

	toRef, err := repo.Reference(plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", to)), true)
	if err != nil {
		return errors.WithStack(appservice.ErrBranchNotFound)
	}

	currentHead, err := repo.Head()
	if err != nil {
		return errors.WithStack(err)
	}

	err = worktree.Checkout(&git.CheckoutOptions{
		Branch: toRef.Name(),
		Force:  false,
	})
	if err != nil {
		return errors.WithStack(err)
	}

	fromCommit, err := repo.CommitObject(fromRef.Hash())
	if err != nil {
		return errors.WithStack(err)
	}

	toCommit, err := repo.CommitObject(toRef.Hash())
	if err != nil {
		return errors.WithStack(err)
	}

	isAncestor, err := fromCommit.IsAncestor(toCommit)
	if err != nil {
		return errors.WithStack(err)
	}

	if isAncestor {
		return nil
	}

	_, err = worktree.Commit(fmt.Sprintf("Merge branch '%s' into '%s'", from, to), &git.CommitOptions{
		Parents: []plumbing.Hash{
			toRef.Hash(),
			fromRef.Hash(),
		},
	})
	if err != nil {
		worktree.Checkout(&git.CheckoutOptions{
			Branch: currentHead.Name(),
			Force:  true,
		})
		return errors.WithStack(err)
	}

	return nil
}

func (g *service) getCurrentRepository() (*git.Repository, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	repo, err := git.PlainOpen(currentDir)
	if stderrors.Is(err, git.ErrRepositoryNotExists) {
		return nil, appservice.ErrGitRepositoryNotFound
	}
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return repo, nil
}
