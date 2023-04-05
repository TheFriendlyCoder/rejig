package cmd

import (
	"bytes"
	"github.com/TheFriendlyCoder/rejigger/lib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"path"
	"path/filepath"
	"testing"
)

func Test_generateUsageLine(t *testing.T) {
	a := assert.New(t)
	result := generateUsageLine()

	a.Contains(result, "create")
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
	r.ErrorAs(result, &lib.PathError{Path: srcDir, ErrorType: lib.PE_PATH_NOT_FOUND})
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
	r.ErrorAs(result, &lib.PathError{Path: destDir, ErrorType: lib.PE_PATH_NOT_EMPTY})
}

func Test_CreateCommandSucceeds(t *testing.T) {

	r := require.New(t)
	a := assert.New(t)

	// Given an empty temp folder
	tmpDir, err := os.MkdirTemp("", "")
	r.NoError(err, "Error creating temp folder")

	// Make sure we always clean up our temp folder
	defer func() {
		r.NoError(os.RemoveAll(tmpDir), "Error deleting temp folder")
	}()

	srcDir, err := sampleProj("simple")
	r.NoError(err, "Locating sample project should always succeed")

	output := new(bytes.Buffer)
	fakeInput := new(bytes.Buffer)
	_, err = fakeInput.WriteString("MyProj\n")
	r.NoError(err, "Failed generating sample input")
	_, err = fakeInput.WriteString("1.2.3\n")
	r.NoError(err, "Failed generating sample input")

	rootCmd.SetOut(output)
	rootCmd.SetErr(output)
	rootCmd.SetIn(fakeInput)
	rootCmd.SetArgs([]string{"create", *srcDir, tmpDir})
	err = rootCmd.Execute()
	r.NoError(err, "CLI command should have succeeded")

	r.Contains(output.String(), *srcDir)
	r.Contains(output.String(), tmpDir)

	a.DirExists(filepath.Join(tmpDir, "MyProj"))
	a.NoFileExists(filepath.Join(tmpDir, ".rejig.yml"))

	//exp := filepath.Join(*srcDir, ".gitignore")
	act := filepath.Join(tmpDir, ".gitignore")
	a.FileExists(act)
	//a.True(unmodified(r, exp, act))

	act = filepath.Join(tmpDir, "version.txt")
	a.FileExists(act)
	//a.True(contains(r, act, expVersion))
	//a.False(contains(r, act, "{{version}}"))

	act = filepath.Join(tmpDir, "MyProj", "main.txt")
	a.FileExists(act)
	//a.True(contains(r, act, expProj))
	//a.False(contains(r, act, "{{project_name}}"))
}

func Test_CreateCommandTooFewArgs(t *testing.T) {

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
	rootCmd.SetArgs([]string{"create", srcDir})
	err = rootCmd.Execute()

	// We expect an error to be returned
	r.Error(err, "CLI command should have failed")

	// and some status information to be reported
	r.Contains(actual.String(), "Error:")
}

func Test_CreateCommandInvalidArgs(t *testing.T) {

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
	rootCmd.SetArgs([]string{"create", srcDir, destDir})
	err = rootCmd.Execute()

	// We expect an error to be returned
	r.Error(err, "CLI command should have succeeded")

	// And some status info to be reported
	r.Contains(actual.String(), srcDir)
}
