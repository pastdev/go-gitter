package gitter

import (
	"fmt"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

// GoGitRepo is a go-git implemementation of a git repo.
type GoGitRepo struct {
	Repo
	goGitRepository *git.Repository
}

// Add implements gitter.Add using go-git.
func (r GoGitRepo) Add(args *AddArgs) error {
	_, worktree, err := r.Worktree()
	if err != nil {
		return err
	}

	for _, path := range args.Paths {
		_, err := worktree.Add(path)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r GoGitRepo) Author() (*object.Signature, error) {
	repo, err := r.Repository()
	if err != nil {
		return nil, err
	}

	c, err := repo.Config()
	if err != nil {
		return nil, fmt.Errorf("unable to read config: %v", err)
	}

	user := c.Raw.Section("user").Options
	return &object.Signature{
		Name:  user.Get("name"),
		Email: user.Get("email"),
	}, nil
}

func (r GoGitRepo) Commit(args *CommitArgs) error {
	_, worktree, err := r.Worktree()
	if err != nil {
		return err
	}

	author, err := r.Author()
	if err != nil {
		return err
	}

	_, err = worktree.Commit(
		args.Message,
		&git.CommitOptions{
			All:    args.All,
			Author: author,
		})
	if err != nil {
		return fmt.Errorf("commit failed: %v", err)
	}

	return nil
}

// GetWorkingDir returns the Repo.WorkingDir.
func (r GoGitRepo) GetWorkingDir() string {
	return r.WorkingDir
}

// Init implements gitter.Init using go-git.
func (r GoGitRepo) Init(args *InitArgs) error {
	repo, err := git.PlainInit(r.WorkingDir, args.Bare)
	if err != nil {
		return err
	}

	r.goGitRepository = repo

	return nil
}

// Status implements gitter.Status
func (r GoGitRepo) Status() (git.Status, error) {
	_, worktree, err := r.Worktree()
	if err != nil {
		return nil, err
	}

	status, err := worktree.Status()
	if err != nil {
		return nil, err
	}

	return status, nil
}

// NewGoGitter returns a new GoGitRepo as a Gitter.
func NewGoGitter(workingDir string) Gitter {
	return NewGoGitRepo(workingDir)
}

// NewGoGitRepo returns a new GoGitRepo.
func NewGoGitRepo(workingDir string) *GoGitRepo {
	repo := GoGitRepo{}
	repo.WorkingDir = workingDir
	return &repo
}

func (r *GoGitRepo) Repository() (*git.Repository, error) {
	if r.goGitRepository != nil {
		return r.goGitRepository, nil
	}

	repo, err := git.PlainOpen(r.WorkingDir)
	if err != nil {
		return nil, err
	}

	r.goGitRepository = repo
	return repo, nil
}

func (r GoGitRepo) Worktree() (*git.Repository, *git.Worktree, error) {
	repo, err := r.Repository()
	if err != nil {
		return nil, nil, err
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return nil, nil, err
	}

	return repo, worktree, nil
}
