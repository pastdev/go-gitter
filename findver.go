package gitter

import (
	"errors"
	"fmt"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

type findVersionMode string

const (
	FindVersionModeVersionTxt = findVersionMode("version.txt")
	FindVersionModePomXml     = findVersionMode("pom.xml")
)

func (r *GoGitRepo) findVersionUsing(baseParser func(string) (string, error)) (string, error) {
	repo, err := r.Repository()
	if err != nil {
		return "", fmt.Errorf("unable to get repository: %v", err)
	}

	// should probably use git rev-list instad of log
	//   https://github.com/src-d/go-git/issues/757
	entries, err := repo.Log(&git.LogOptions{})
	if err != nil {
		return "", err
	}

	versionChanged := errors.New("")
	versionFile := "version.txt"
	current := ""
	depth := 0

	err = entries.ForEach(func(c *object.Commit) error {
		tree, err := c.Tree()
		if err != nil {
			return fmt.Errorf("cannot retrive tree for %v: %v", c, err)
		}

		file, err := tree.File(versionFile)
		if err != nil {
			return fmt.Errorf("unable to retrieve %s: %v", versionFile, err)
		}

		versionContent, err := file.Contents()
		if err != nil {
			return fmt.Errorf("unable to read %s: %v", versionFile, err)
		}

		versionBase, err := baseParser(versionContent)

		if current == "" {
			current = versionBase
			depth = 0
			return nil
		}

		if versionBase != current {
			return versionChanged
		}

		depth = depth + 1

		return nil
	})
	if err != nil && err != versionChanged {
		return "", err
	}

	return fmt.Sprintf("%s.%d", current, depth), nil
}

func (r *GoGitRepo) FindVersion(mode findVersionMode) (string, error) {
	switch mode {
	case FindVersionModePomXml:
		return "", errors.New("not yet implemented")
	case FindVersionModeVersionTxt:
		return r.findVersionUsing(func(content string) (string, error) {
			return content, nil
		})
	default:
		return "", fmt.Errorf("unsupported mode %s", mode)
	}
}
