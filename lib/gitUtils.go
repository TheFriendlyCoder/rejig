package lib

import (
	"gopkg.in/src-d/go-billy.v4"
	"gopkg.in/src-d/go-billy.v4/memfs"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

// getTemplate loads a template source from a Git repository
func getTemplate(gitURL string) (*billy.Filesystem, error) {
	fs := memfs.New()
	_, err := git.Clone(memory.NewStorage(), fs, &git.CloneOptions{
		URL: gitURL,
	})
	if err != nil {
		return nil, err
	}
	return &fs, nil
}
