package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/cache"
	"github.com/go-git/go-git/v5/storage/filesystem"
)

func openRepo() (string, *git.Repository, error) {
	root, err := findRepoRoot()
	if err != nil {
		return "", nil, err
	}
	r, err := git.Open(
		filesystem.NewStorage(
			osfs.New(filepath.Join(root, ".git")),
			cache.NewObjectLRUDefault(),
		),
		osfs.New(root),
	)
	return root, r, err
}

func findRepoRoot() (string, error) {
	p, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		_, err := os.Stat(filepath.Join(p, ".git"))
		if err == nil {
			return p, nil
		}
		if !os.IsNotExist(err) {
			return "", err
		}
		newp := filepath.Join(p, "..")
		if newp == p {
			break
		}
		p = newp
	}
	return "", fmt.Errorf("fatal: not a git repository")
}

func repoRelPath(root, filename string) (string, error) {
	filename, err := filepath.Abs(filename)
	if err != nil {
		return "", err
	}
	return filepath.Rel(root, filename)
}
