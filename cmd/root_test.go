package cmd

import (
	"bytes"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"path"
	"testing"
)

// sampleData loads sample data for a test from the test data folder
func sampleData(filename string) (string, error) {
	retval := path.Join("testdata", filename)
	var _, err = os.Stat(retval)
	return retval, errors.Wrap(err, "checking existence of test data file")
}

func Test_generateUsageLine(t *testing.T) {
	a := assert.New(t)
	result := generateUsageLine()

	a.Contains(result, "rejig")
	a.Contains(result, "sourcePath")
	a.Contains(result, "targetPath")
}

func Test_ValidateArgsSuccess(t *testing.T) {
	r := require.New(t)

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

	// when we validate our input args
	args := []string{srcDir, destDir}
	r.NoError(validateArgs(args))
}

func Test_ValidateArgsSourceDirNotExists(t *testing.T) {
	r := require.New(t)

	// Given an empty temp folder
	tmpDir, err := os.MkdirTemp("", "")
	r.NoError(err, "Error creating temp folder")

	// Make sure we always clean up our temp folder
	defer func() {
		r.NoError(os.RemoveAll(tmpDir), "Error deleting temp folder")
	}()

	// with 1 subfolder that doesn't exist
	destDir := path.Join(tmpDir, "dest")
	r.NoError(os.Mkdir(destDir, 0700), "Error creating destination folder")
	srcDir := path.Join(tmpDir, "src")

	// when we validate our input args
	args := []string{srcDir, destDir}
	result := validateArgs(args)

	// we expect the proper error to be returned
	r.Error(result)
	r.ErrorAs(result, &pathError{path: srcDir, errorType: pathNotFound})
}

func Test_ValidateArgsTargetDirNotEmpty(t *testing.T) {
	r := require.New(t)

	// Given an empty temp folder
	tmpDir, err := os.MkdirTemp("", "")
	r.NoError(err, "Error creating temp folder")

	// Make sure we always clean up our temp folder
	defer func() {
		r.NoError(os.RemoveAll(tmpDir), "Error deleting temp folder")
	}()

	// with 2 subfolders
	destDir := path.Join(tmpDir, "dest")
	srcDir := path.Join(tmpDir, "src")
	r.NoError(os.Mkdir(destDir, 0700), "Error creating destination folder")
	r.NoError(os.Mkdir(srcDir, 0700), "Error creating destination folder")

	// and a non-empty destination folder
	_, err = os.Create(path.Join(destDir, "fubar.txt"))
	r.NoError(err, "Failed to create test file")

	// when we validate our input args
	args := []string{srcDir, destDir}
	result := validateArgs(args)

	// we expect the proper error to be returned
	r.Error(result)
	r.ErrorAs(result, &pathError{path: destDir, errorType: pathNotEmpty})
}

func Test_RootCommandSucceeds(t *testing.T) {

	r := require.New(t)

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

	actual := new(bytes.Buffer)
	rootCmd.SetOut(actual)
	rootCmd.SetErr(actual)
	rootCmd.SetArgs([]string{srcDir, destDir})
	err = rootCmd.Execute()
	r.NoError(err, "CLI command should have succeeded")

	r.Contains(actual.String(), srcDir)
	r.Contains(actual.String(), destDir)
}

func Test_RootCommandTooFewArgs(t *testing.T) {

	r := require.New(t)

	// Given an empty temp folder
	tmpDir, err := os.MkdirTemp("", "")
	r.NoError(err, "Error creating temp folder")

	// Make sure we always clean up our temp folder
	defer func() {
		r.NoError(os.RemoveAll(tmpDir), "Error deleting temp folder")
	}()

	// with 1 empty subfolder
	srcDir := path.Join(tmpDir, "src")
	r.NoError(os.Mkdir(srcDir, 0700), "Error creating source folder")

	// When we execute our root command with a missing positional arg
	actual := new(bytes.Buffer)
	rootCmd.SetOut(actual)
	rootCmd.SetErr(actual)
	rootCmd.SetArgs([]string{srcDir})
	err = rootCmd.Execute()

	// We expect an error to be returned
	r.Error(err, "CLI command should have failed")

	// and some status information to be reported
	r.Contains(actual.String(), "Error:")
}

func Test_RootCommandInvalidArgs(t *testing.T) {

	r := require.New(t)

	// Given an empty temp folder
	tmpDir, err := os.MkdirTemp("", "")
	r.NoError(err, "Error creating temp folder")

	// Make sure we always clean up our temp folder
	defer func() {
		r.NoError(os.RemoveAll(tmpDir), "Error deleting temp folder")
	}()

	// with an absent source folder
	srcDir := path.Join(tmpDir, "src")

	// with a valid dest folder
	destDir := path.Join(tmpDir, "dest")
	r.NoError(os.Mkdir(destDir, 0700), "Error creating destination folder")

	// when we attempt to execute our root command
	actual := new(bytes.Buffer)
	rootCmd.SetOut(actual)
	rootCmd.SetErr(actual)
	rootCmd.SetArgs([]string{srcDir, destDir})
	err = rootCmd.Execute()

	// We expect an error to be returned
	r.Error(err, "CLI command should have succeeded")

	// And some status info to be reported
	r.Contains(actual.String(), srcDir)
}

func Test_loadManifestFile(t *testing.T) {
	r := require.New(t)
	sampleManifest, err := sampleData("simple_manifest.yml")
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
	rootCmd.SetArgs([]string{srcDir, destDir, "--config=" + sampleManifest})
	err = rootCmd.Execute()

	// We expect no error from our command
	r.NoError(err, "CLI command should have succeeded")

	// And some status info to be reported
	//r.Contains(actual.String(), sampleManifest)
}
