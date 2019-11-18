package gitter_test

import (
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"strings"
	"testing"

	"github.com/pastdev/go-gitter"
	"gopkg.in/src-d/go-git.v4"
)

func testAddSucceeds(
	t *testing.T,
	name string,
	new func(workingDir string) gitter.Gitter,
	args *gitter.InitArgs) {

	dir, err := gitter.NewTestFixture(new)
	if err != nil {
		t.Errorf("unable to create temp dir: %v", err)
	}
	defer dir.Cleanup()

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
	testAddSucceeds(t,
		"runnergit simple add",
		gitter.NewRunnerGitter,
		&gitter.InitArgs{})

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

	fixture, err := gitter.NewTestFixture(new)
	if err != nil {
		t.Errorf("%s .git missing: %v", name, err)
	}

	if info, err := os.Stat(path.Join(fixture.WorkingDir, ".git")); !(err == nil && info.IsDir()) {
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
			"R  README.md\000RE D\\\"ME .md\"",
			"A  \"foo bar\"",
			"MD foobar",
			"?? abc/",
			"?? asdf fdsa",
			"",
		},
		"\000")
	status := gitter.ParseStatusZ(stdout)

	if len(status) != 5 {
		t.Errorf("Expected 5 changes, found %d", len(status))
	}

	name := "README.md"
	expected := &git.FileStatus{
		Staging:  git.Renamed,
		Worktree: git.Unmodified,
		Extra:    "RE D\\\"ME .md\"",
	}
	actual := status.File(name)
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("%s: %v != %v", name, expected, actual)
	}

	name = "\"foo bar\""
	expected = &git.FileStatus{
		Staging:  git.Added,
		Worktree: git.Unmodified,
	}
	actual = status.File(name)
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("%s: %v != %v", name, expected, actual)
	}

	name = "foobar"
	expected = &git.FileStatus{
		Staging:  git.Modified,
		Worktree: git.Deleted,
	}
	actual = status.File(name)
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("%s: %v != %v", name, expected, actual)
	}

	name = "abc/"
	expected = &git.FileStatus{
		Staging:  git.Untracked,
		Worktree: git.Untracked,
	}
	actual = status.File(name)
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("%s: %v != %v", name, expected, actual)
	}

	name = "asdf fdsa"
	expected = &git.FileStatus{
		Staging:  git.Untracked,
		Worktree: git.Untracked,
	}
	actual = status.File(name)
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("%s: %v != %v", name, expected, actual)
	}
}
