package cmd

import (
	"bytes"
	"github.com/TheFriendlyCoder/rejigger/lib"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"path"
	"path/filepath"
	"testing"
)

// sampleData loads sample data for a test from the test data folder
func sampleData(filename string) (string, error) {
	retval := path.Join("testdata", filename)
	var _, err = os.Stat(retval)
	return retval, errors.Wrap(err, "checking existence of test data file")
}

// sampleProj loads path to a specific sample project to use for testing the generator logic
func sampleProj(projName string) (*string, error) {
	retval, err := filepath.Abs(path.Join("..", "testProjects", projName))
	if err != nil {
		return nil, errors.Wrap(err, "Failed to generate absolute path")
	}
	_, err = os.Stat(retval)
	if err != nil {
		return nil, errors.Wrap(err, "checking existence of test data file")
	}
	return &retval, nil
}

func Test_loadAppOptionsFile(t *testing.T) {
	r := require.New(t)
	a := assert.New(t)

	sampleConfig, err := sampleData("app_options.yml")
	r.NoError(err, "sample config file not found")

	// Given an empty temp folder
	tmpDir, err := os.MkdirTemp("", "")
	r.NoError(err, "Error creating temp folder")

	// Make sure we always clean up our temp folder
	defer func() {
		r.NoError(os.RemoveAll(tmpDir), "Error deleting temp folder")
	}()

	// with 2 empty subfolders
	srcDir := path.Join(tmpDir, "src")
	destDir := path.Join(tmpDir, "dest")
	r.NoError(os.Mkdir(srcDir, 0700), "Error creating source folder")
	r.NoError(os.Mkdir(destDir, 0700), "Error creating destination folder")

	// When we run the root command with a custom config file
	actual := new(bytes.Buffer)
	rootCmd.SetOut(actual)
	rootCmd.SetErr(actual)
	// NOTE: We have to pass a command, 'help' in this case, to trigger the app init
	rootCmd.SetArgs([]string{"--config=" + sampleConfig, "help"})
	err = rootCmd.Execute()

	// We expect no error from our command
	r.NoError(err, "CLI command should have succeeded")

	// Validate parsed options
	a.Equal(1, len(appOptions.Templates))
	a.Equal("simple1", appOptions.Templates[0].Alias)
	a.Equal("testProjects/simple", appOptions.Templates[0].Folder)
	a.Equal("https://github.com/TheFriendlyCoder/rejigger", appOptions.Templates[0].Source)
	a.Equal(lib.TST_GIT, appOptions.Templates[0].Type)
}

func Test_loadAppOptionsFileInvalidTemplateType(t *testing.T) {
	r := require.New(t)

	sampleConfig, err := sampleData("app_options_invalid_template_type.yml")
	r.NoError(err, "sample config file not found")

	// Given an empty temp folder
	tmpDir, err := os.MkdirTemp("", "")
	r.NoError(err, "Error creating temp folder")

	// Make sure we always clean up our temp folder
	defer func() {
		r.NoError(os.RemoveAll(tmpDir), "Error deleting temp folder")
	}()

	// with 2 empty subfolders
	srcDir := path.Join(tmpDir, "src")
	destDir := path.Join(tmpDir, "dest")
	r.NoError(os.Mkdir(srcDir, 0700), "Error creating source folder")
	r.NoError(os.Mkdir(destDir, 0700), "Error creating destination folder")

	// When we run the root command with a custom config file
	actual := new(bytes.Buffer)
	rootCmd.SetOut(actual)
	rootCmd.SetErr(actual)
	rootCmd.SetArgs([]string{"--config=" + sampleConfig, "help"})
	r.PanicsWithValue("Critical application failure", func() { err = rootCmd.Execute() })
	r.NoError(err)
}

// TODO: have suite auto-reset viper singleton on each test loop
