package gitter

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	log "github.com/sirupsen/logrus"
)

type TestFixture struct {
	Dir              string
	OriginDir        string
	OriginGitter     Gitter
	OriginGoGitRepo  *GoGitRepo
	WorkingDir       string
	WorkingGitter    Gitter
	WorkingGoGitRepo *GoGitRepo
}

func (f *TestFixture) Cleanup() {
	os.RemoveAll(f.Dir)
}

func (f *TestFixture) Commit(message string) error {
	return f.WorkingGitter.Commit(&CommitArgs{
		Message: message,
	})
}

func NewTestFixture(new func(workingDir string) Gitter) (*TestFixture, error) {
	temp, err := ioutil.TempDir("", "gitter_test_")
	if err != nil {
		return nil, fmt.Errorf("unable to create temp dir: %v", err)
	}

	dir := &TestFixture{
		Dir:        temp,
		WorkingDir: path.Join(temp, "working"),
		OriginDir:  path.Join(temp, "origin.git"),
	}

	err = dir.initialize(new)
	if err != nil {
		dir.Cleanup()
		return nil, err
	}

	return dir, nil
}

func (f *TestFixture) initialize(new func(workingDir string) Gitter) error {
	err := os.MkdirAll(f.WorkingDir, 0700)
	if err != nil {
		return fmt.Errorf("unable to create temp working dir: %v", err)
	}
	f.WorkingGitter = new(f.WorkingDir)
	err = f.WorkingGitter.Init(&InitArgs{})
	if err != nil {
		return fmt.Errorf("unable to init temp working dir: %v", err)
	}
	f.WorkingGoGitRepo = NewGoGitRepo(f.WorkingDir)

	f.setAuthor("testfixture", "testfixture@example.com")

	err = os.MkdirAll(f.OriginDir, 0700)
	if err != nil {
		return fmt.Errorf("unable to create temp origin dir: %v", err)
	}
	f.OriginGitter = new(f.OriginDir)
	err = f.OriginGitter.Init(&InitArgs{Bare: true})
	if err != nil {
		return fmt.Errorf("unable to init temp origin dir: %v", err)
	}
	f.OriginGoGitRepo = NewGoGitRepo(f.OriginDir)

	return nil
}

func (f *TestFixture) setAuthor(name, email string) error {
	repo, err := f.WorkingGoGitRepo.Repository()
	if err != nil {
		return fmt.Errorf("unable to get repository: %v", err)
	}
	config, err := repo.Config()
	if err != nil {
		return fmt.Errorf("unable to load config: %v", err)
	}
	user := config.Raw.Section("user")
	user.SetOption("name", name)
	user.SetOption("email", email)

	return nil
}

func (f *TestFixture) workingPath(name string) string {
	return path.Join(f.WorkingDir, name)
}

func (f *TestFixture) WriteFile(name, content string) error {
	file := f.workingPath(name)
	log.Debugf("writing file to: %s", file)
	err := os.MkdirAll(path.Dir(file), 0700)
	if err != nil {
		return fmt.Errorf("unable to create dir path for %s: %v", file, err)
	}
	return ioutil.WriteFile(file, []byte(content), 0600)
}

func (f *TestFixture) WriteFileAndAdd(name, content string) error {
	err := f.WriteFile(name, content)
	if err != nil {
		return err
	}

	return f.WorkingGoGitRepo.Add(
		&AddArgs{
			Paths: []string{
				name,
			},
		})
}

func (f *TestFixture) WriteFileAndAddAndCommit(name, content, message string) error {
	err := f.WriteFileAndAdd(name, content)
	if err != nil {
		return err
	}

	return f.Commit(message)
}
