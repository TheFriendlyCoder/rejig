package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

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
