package gitter_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/pastdev/go-gitter"
	"gopkg.in/src-d/go-git.v4"
)

type TestDir struct {
	Dir           string
	WorkingDir    string
	WorkingGitter gitter.Gitter
	OriginDir     string
	OriginGitter  gitter.Gitter
}

func newTestDir(new func(workingDir string) gitter.Gitter) (*TestDir, error) {
	temp, err := ioutil.TempDir("", "gitter_test_")
	if err != nil {
		return nil, fmt.Errorf("unable to create temp dir: %v", err)
	}

	dir := TestDir{
		Dir:        temp,
		WorkingDir: path.Join(temp, "working"),
		OriginDir:  path.Join(temp, "origin.git"),
	}

	err = os.MkdirAll(dir.WorkingDir, 0700)
	if err != nil {
		return nil, fmt.Errorf("unable to create temp working dir: %v", err)
	}
	dir.WorkingGitter = new(dir.WorkingDir)
	err = dir.WorkingGitter.Init(&gitter.InitArgs{})
	if err != nil {
		return nil, fmt.Errorf("unable to init temp working dir: %v", err)
	}

	err = os.MkdirAll(dir.OriginDir, 0700)
	if err != nil {
		return nil, fmt.Errorf("unable to create temp origin dir: %v", err)
	}
	dir.OriginGitter = new(dir.OriginDir)
	err = dir.OriginGitter.Init(&gitter.InitArgs{Bare: true})
	if err != nil {
		return nil, fmt.Errorf("unable to init temp origin dir: %v", err)
	}

	return &dir, nil
}

func (d *TestDir) cleanup() {
	os.RemoveAll(d.Dir)
}

func testAddSucceeds(
	t *testing.T,
	name string,
	new func(workingDir string) gitter.Gitter,
	args *gitter.InitArgs) {

	dir, err := newTestDir(new)
	if err != nil {
		t.Errorf("unable to create temp dir: %v", err)
	}
	defer dir.cleanup()

	workingDir := dir.WorkingDir
	g := dir.WorkingGitter

	readme := path.Join(workingDir, "README.md")
	err = ioutil.WriteFile(readme, []byte("# Gitter"), 0600)
	if err != nil {
		t.Errorf("%s write readme failed: %v", name, err)
	}
	g.Add(&gitter.AddArgs{Paths: []string{readme}})

	status, err := g.Status()
	if err != nil {
		t.Errorf("%s git status failed: %v", name, err)
	}
	if status.IsClean() {
		t.Errorf("%s git status should not have been clean", name)
	}
}

func Test_AddSucceeds(t *testing.T) {
	// testAddSucceeds(t,
	// 	"runnergit simple add",
	// 	gitter.NewRunnerGitter,
	// 	&gitter.InitArgs{})

	testAddSucceeds(t,
		"gogit simple add",
		gitter.NewGoGitter,
		&gitter.InitArgs{})
}

func testInitSucceeds(
	t *testing.T,
	name string,
	new func(workingDir string) gitter.Gitter,
	args *gitter.InitArgs) {

	// newTestDir init's working regular, and origin with --bare
	dir, err := newTestDir(new)
	if err != nil {
		t.Errorf("%s .git missing: %v", name, err)
	}

	if info, err := os.Stat(path.Join(dir.WorkingDir, ".git")); !(err == nil && info.IsDir()) {
		t.Errorf("%s .git missing: %v", name, err)
	}
}

func Test_InitSucceeds(t *testing.T) {
	testInitSucceeds(t,
		"runnergit simple init",
		gitter.NewRunnerGitter,
		&gitter.InitArgs{})

	testInitSucceeds(t,
		"gogit simple init",
		gitter.NewGoGitter,
		&gitter.InitArgs{})
}

func Test_ParseStatusZ(t *testing.T) {
	stdout := strings.Join(
		[]string{
			"A  \"RE D\\\"ME .md\" -> README.md",
			" D README.md",
			"A  \"foo bar\"",
			"MD foobar",
			"?? abc/",
			"?? asdf fdsa",
		},
		"\000")
	status := gitter.ParseStatusZ(stdout)
	if status.File("RE DME .md").Staging != git.Added {
		t.Errorf("expected added")
	}
}
