package cmd

import (
	"bytes"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"os"
	"path"
	"path/filepath"
	"testing"
)

// sampleProj loads path to a specific sample project to use for testing the generator logic
func sampleProj(projName string) (string, error) {
	retval, err := filepath.Abs(path.Join("..", "testProjects", projName))
	if err != nil {
		return retval, errors.Wrap(err, "Failed to generate absolute path")
	}
	_, err = os.Stat(retval)
	if err != nil {
		return retval, errors.Wrap(err, "checking existence of test data file")
	}
	return retval, nil
}

func Test_helpCommand(t *testing.T) {
	r := require.New(t)

	// When we run the root command with a custom config file
	actual := new(bytes.Buffer)
	rootCmd.SetOut(actual)
	rootCmd.SetErr(actual)

	rootCmd.SetArgs([]string{"help"})
	err := rootCmd.Execute()

	// We expect no error from our command
	r.NoError(err, "CLI command should have succeeded")
	r.Contains(actual.String(), "rejigger")
}
