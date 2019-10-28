package gitter

import (
	"gopkg.in/src-d/go-git.v4"
)

// GoGitRepo is a go-git implemementation of a git repo.
type GoGitRepo struct {
	Repo
	Repository *git.Repository
}

func (r GoGitRepo) worktree() (*git.Worktree, error) {
	var err error
	if r.Repository == nil {
		r.Repository, err = git.PlainOpen(r.WorkingDir)
		if err != nil {
			return nil, err
		}
	}

	worktree, err := r.Repository.Worktree()
	if err != nil {
		return nil, err
	}

	return worktree, nil
}

// Add implements gitter.Add using go-git.
func (r GoGitRepo) Add(args *AddArgs) error {
	worktree, err := r.worktree()
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

	r.Repository = repo

	return nil
}

// Status implements gitter.Status
func (r GoGitRepo) Status() (git.Status, error) {
	worktree, err := r.worktree()
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
