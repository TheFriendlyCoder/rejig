package lib

import (
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

// SNF ShouldNotFail helper method that looks through all return values from a function call,
// detects any error objects that are returned, and if any of them are failure objects then
// it triggers a panic operation, since errors should not typically happen in practice
// Example use:
//
//	lib.SNF(fmt.FPrintf(...))
func SNF(args ...any) {
	// loop through all input args
	for _, c := range args {
		// see if any can be cast to an error type
		if e, ok := c.(error); ok {
			// and if so, we know it is a non-nil error result
			// so we want to trigger a panic
			panic("Unexpected error: " + e.Error())
		}
	}
}

// GetGitFilesystem loads a remote Git repository into an in-memory virtual file system
func GetGitFilesystem(gitURL string) (afero.Fs, error) {
	appFS := afero.NewMemMapFs()
	fs := NewBillyWraper(appFS, ".", false)

	_, err := git.Clone(memory.NewStorage(), fs, &git.CloneOptions{
		URL: gitURL,
	})
	if err != nil {
		return appFS, errors.Wrap(err, "Failed querying Git repo")
	}
	return appFS, nil
}
