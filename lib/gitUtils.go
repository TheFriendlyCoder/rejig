package lib

import (
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

// GetGitTemplate loads a template source from a Git repository
func GetGitTemplate(gitURL string) (afero.Fs, error) {
	appFS := afero.NewMemMapFs()
	fs := NewWraper(appFS, ".", false)

	_, err := git.Clone(memory.NewStorage(), fs, &git.CloneOptions{
		URL: gitURL,
	})
	if err != nil {
		return appFS, errors.Wrap(err, "Failed querying Git repo")
	}
	return appFS, nil
}
