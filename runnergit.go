package gitter

import (
	"bytes"
	"errors"
	"os/exec"
	"strings"

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

	cmd := exec.Command("git", cmdArgs...)
	stdout, err := cmd.Output()
	if err != nil {
		return err
	}
	log.Debug(stdout)

	return nil
}

// Status implements Gitter.Status using exec git.
func (r *RunnerGitRepo) Status() (git.Status, error) {
	cmdArgs := []string{
		"-C",
		r.WorkingDir,
		"status",
		"-z",
	}

	cmd := exec.Command("git", cmdArgs...)
	stdout, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	log.Debug(stdout)

	return ParseStatusZ(string(stdout)), nil
}

func ParseStatusZ(stdout string) git.Status {
	status := git.Status{}
	for _, line := range strings.Split(stdout, "\000") {
		fileStatus := git.FileStatus{
			Staging:  git.StatusCode(line[0]),
			Worktree: git.StatusCode(line[1]),
		}

		// untracked dont seem to ever be quoted:
		//  ltheisen@ltp52s:/tmp/gitter$ git status --porcelain
		//  ?? asdf fdsa
		// so dont need to parse
		name := line[3:]

		// tracked seem to be quoted
		//  ltheisen@ltp52s:/tmp/gitter$ git status --porcelain
		//  A  "RE DME .md"
		if fileStatus.Staging != git.Untracked {
			var extra string
			name, extra = parseStatusNameExtra(line[3:])
			fileStatus.Extra = extra
		}

		status[name] = &fileStatus
	}

	return status
}

// parseStatusNameExtra parses the nameExtra portion of status following:
// Short Format
// In the short-format, the status of each path is shown as one of these forms
//
// XY PATH
// XY ORIG_PATH -> PATH
//
// where ORIG_PATH is where the renamed/copied contents came from. ORIG_PATH is only shown when the entry is renamed or copied. The XY is a
// two-letter status code.
//
// The fields (including the ->) are separated from each other by a single space. If a filename contains whitespace or other nonprintable
// characters, that field will be quoted in the manner of a C string literal: surrounded by ASCII double quote (34) characters, and with
// interior special characters backslash-escaped.
//
// For paths with merge conflicts, X and Y show the modification states of each side of the merge. For paths that do not have merge conflicts,
// X shows the status of the index, and Y shows the status of the work tree. For untracked paths, XY are ??. Other status codes can be
// interpreted as follows:
//
// ·   ' ' = unmodified
// ·   M = modified
// ·   A = added
// ·   D = deleted
// ·   R = renamed
// ·   C = copied
// ·   U = updated but unmerged
func parseStatusNameExtra(nameExtra string) (string, string) {
	backslash := rune(92)
	quote := rune(34)
	space := rune(32)

	extra := ""

	var buffer bytes.Buffer
	prev := rune(0)
	inQuote := false
	skip := 0
	for _, r := range nameExtra {
		if skip > 0 {
			skip = skip - 1
			continue
		}
		if r == quote {
			if inQuote {
				if prev != backslash {
					inQuote = false
					continue
				}
			} else {
				inQuote = true
				continue
			}
		}
		if !inQuote && r == space {
			extra = buffer.String()
			buffer.Reset()
			skip = len("-> ")
			continue
		}
		buffer.WriteRune(r)
		prev = r
	}

	return buffer.String(), extra
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
