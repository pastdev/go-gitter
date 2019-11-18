package gitter_test

import (
	"testing"

	"github.com/pastdev/go-gitter"
)

func Test_FindVersion(t *testing.T) {
	f, err := gitter.NewTestFixture(gitter.NewGoGitter)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Cleanup()

	err = f.WriteFileAndAddAndCommit("version.txt", "0.1", "initial")
	if err != nil {
		t.Fatal(err)
	}
	expected := "0.1.0"
	actual, err := f.WorkingGoGitRepo.FindVersion(gitter.FindVersionModeVersionTxt)
	if err != nil {
		t.Fatal(err)
	}
	if expected != actual {
		t.Fatalf("[%s] != [%s]", expected, actual)
	}

	err = f.WriteFileAndAddAndCommit("README.md", "# Test", "add readme")
	if err != nil {
		t.Fatal(err)
	}
	expected = "0.1.1"
	actual, err = f.WorkingGoGitRepo.FindVersion(gitter.FindVersionModeVersionTxt)
	if err != nil {
		t.Fatal(err)
	}
	if expected != actual {
		t.Fatalf("[%s] != [%s]", expected, actual)
	}

	err = f.WriteFileAndAddAndCommit("README.md", "# Test\n\nTest rules", "update readme")
	if err != nil {
		t.Fatal(err)
	}
	expected = "0.1.2"
	actual, err = f.WorkingGoGitRepo.FindVersion(gitter.FindVersionModeVersionTxt)
	if err != nil {
		t.Fatal(err)
	}
	if expected != actual {
		t.Fatalf("[%s] != [%s]", expected, actual)
	}

	err = f.WriteFileAndAddAndCommit("version.txt", "0.2", "version update")
	if err != nil {
		t.Fatal(err)
	}
	expected = "0.2.0"
	actual, err = f.WorkingGoGitRepo.FindVersion(gitter.FindVersionModeVersionTxt)
	if err != nil {
		t.Fatal(err)
	}
	if expected != actual {
		t.Fatalf("[%s] != [%s]", expected, actual)
	}
}
