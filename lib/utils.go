package lib

import (
	"fmt"
	"os"
	"strings"

	"github.com/TheFriendlyCoder/rejigger/lib/thirdparty"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"gopkg.in/src-d/go-git.v4"
	ssh2 "gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
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
	fs := thirdparty.NewBillyWraper(appFS, ".", false)

	opts := git.CloneOptions{
		URL: gitURL,
	}

	if !strings.HasPrefix(gitURL, "http") {
		// TODO: Figure out some way to unit test this block
		sshFile := fmt.Sprintf("%s/.ssh/id_rsa", os.Getenv("HOME"))
		_, err := os.Stat(sshFile)
		if os.IsNotExist(err) {
			return appFS, errors.Wrap(err, fmt.Sprintf("Can not find SSH key %s. Run ssh-keygen first.", sshFile))
		} else if err != nil {
			return appFS, errors.WithStack(err)
		}

		// TODO: add support for encrypted SSH key
		authKey, err := ssh2.NewPublicKeysFromFile("git", sshFile, "")
		if err != nil {
			return appFS, errors.WithStack(err)
		}

		opts.Auth = authKey
	}

	_, err := git.Clone(memory.NewStorage(), fs, &opts)
	if err != nil {
		return appFS, errors.Wrap(err, "Failed to load remote Git repository: "+gitURL)
	}
	return appFS, nil
}
