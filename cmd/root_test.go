package cmd

import (
	"bytes"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_helpCommand(t *testing.T) {
	r := require.New(t)

	tmpDir, err := os.MkdirTemp("", "")
	r.NoError(err)
	defer os.RemoveAll(tmpDir)

	// Point our command to an empty config file so the test suite
	// doesn't try and read config options from the users home folder
	configFile := path.Join(tmpDir, "config.yml")
	fh, err := os.Create(configFile)
	r.NoError(err)
	_, err = fh.WriteString("")
	r.NoError(err)
	r.NoError(fh.Close())

	// When we run the root command with a custom config file
	actual := new(bytes.Buffer)
	rootCmd.SetOut(actual)
	rootCmd.SetErr(actual)

	rootCmd.SetArgs([]string{"--config", configFile, "help"})
	err = rootCmd.Execute()

	// We expect no error from our command
	r.NoError(err, "CLI command should have succeeded")
	r.Contains(actual.String(), "rejigger")
}
