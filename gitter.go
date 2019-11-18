package gitter

import (
	"gopkg.in/src-d/go-git.v4"
)

type Repo struct {
	WorkingDir string
}

type AddArgs struct {
	Paths []string
}

type CommitArgs struct {
	All     bool
	Message string
}

type InitArgs struct {
	Bare bool
}

type Gitter interface {
	Add(*AddArgs) error
	Commit(*CommitArgs) error
	GetWorkingDir() string
	Init(*InitArgs) error
	Status() (git.Status, error)
}
