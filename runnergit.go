package gitter

import (
	"bytes"
	"errors"
	"os/exec"

	log "github.com/sirupsen/logrus"
	"gopkg.in/src-d/go-git.v4"
)

// RunnerGitRepo is an exec git implemementation of a git repo.
type RunnerGitRepo struct {
	Repo
}

// Add implements Gitter.Add using exec git.
func (r *RunnerGitRepo) Add(args *AddArgs) error {
	if len(args.Paths) <= 0 {
		return errors.New("args.Paths requires at least one path")
	}
	cmdArgs := []string{
		"-C",
		r.WorkingDir,
		"add",
	}

	cmdArgs = append(cmdArgs, "--")
	cmdArgs = append(cmdArgs, args.Paths...)

	cmd := exec.Command("git", cmdArgs...)
	stdout, err := cmd.Output()
	if err != nil {
		return err
	}
	log.Debug(stdout)

	return nil
}

func runForStdout(args ...string) ([]byte, error) {
	cmd := exec.Command("git", args...)
	return cmd.Output()
}

func run(args ...string) error {
	stdout, err := runForStdout(args...)
	if err != nil {
		return err
	}
	log.Debug(stdout)

	return nil
}

func (r *RunnerGitRepo) Commit(args *CommitArgs) error {
	cmdArgs := []string{
		"-C",
		r.WorkingDir,
		"commit",
		"--message",
		args.Message,
	}

	if args.All {
		cmdArgs = append(cmdArgs, "--all")
	}

	return run(cmdArgs...)
}

// GetWorkingDir returns the Repo.WorkingDir.
func (r *RunnerGitRepo) GetWorkingDir() string {
	return r.WorkingDir
}

// Init implements Gitter.Init using exec git.
func (r *RunnerGitRepo) Init(args *InitArgs) error {
	cmdArgs := []string{
		"-C",
		r.WorkingDir,
		"init",
	}

	if args.Bare {
		cmdArgs = append(cmdArgs, "--bare")
	}

	return run(cmdArgs...)
}

// Status implements Gitter.Status using exec git.
func (r *RunnerGitRepo) Status() (git.Status, error) {
	cmdArgs := []string{
		"-C",
		r.WorkingDir,
		"status",
		"-z",
	}

	stdout, err := runForStdout(cmdArgs...)
	if err != nil {
		return nil, err
	}
	log.Debug(stdout)

	return ParseStatusZ(string(stdout)), nil
}

// ParseStatusZ parses git status -z output.  From the doucmentation:
//
// There is also an alternate -z format recommended for machine parsing. In
// that format, the status field is the same, but some other things change.
// First, the -> is omitted from rename entries and the field order is
// reversed (e.g from -> to becomes to from). Second, a NUL (ASCII 0)
// follows each filename, replacing space as a field separator and the
// terminating newline (but a space still separates the status field from
// the first filename). Third, filenames containing special characters are
// not specially formatted; no quoting or backslash-escaping is performed.
func ParseStatusZ(stdout string) git.Status {
	null := rune(0)

	status := git.Status{}

	current := stdout
	for len(current) > 0 {
		fileStatus := &git.FileStatus{
			Staging:  git.StatusCode(current[0]),
			Worktree: git.StatusCode(current[1]),
		}

		var buffer bytes.Buffer
		var i int
		var r rune
		name := ""
		for i, r = range current[3:] {
			if r == null {
				if name == "" {
					name = buffer.String()
					if fileStatus.Staging == git.Renamed || fileStatus.Worktree == git.Renamed {
						buffer.Reset()
						continue
					}
				} else {
					fileStatus.Extra = buffer.String()
				}
				break
			}

			buffer.WriteRune(r)
		}

		status[name] = fileStatus
		current = current[4+i:]
	}

	return status
}

// NewRunnerGitter returns a new RunnerGitRepo as a Gitter.
func NewRunnerGitter(workingDir string) Gitter {
	return NewRunnerGitRepo(workingDir)
}

// NewRunnerGitRepo returns a new RunnerGitRepo.
func NewRunnerGitRepo(workingDir string) *RunnerGitRepo {
	repo := RunnerGitRepo{}
	repo.WorkingDir = workingDir
	return &repo
}
