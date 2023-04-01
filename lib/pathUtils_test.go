package lib

import (
	"github.com/stretchr/testify/require"
	"os"
	"path"
	"testing"
)

func TestDirExists(t *testing.T) {
	result := DirExists(".")
	require.True(t, result)
}

func TestDirNotExists(t *testing.T) {
	r := require.New(t)

	// Given an empty temp folder
	tmpDir, err := os.MkdirTemp("", "")
	r.NoError(err, "Error creating temp folder")

	// Make sure we always clean up our temp folder
	defer func() {
		err = os.RemoveAll(tmpDir)
		r.NoError(err, "Error deleting temp folder")
	}()

	// When we request a check on a file system object that doesn't exist
	result := DirExists(path.Join(tmpDir, "fubar"))
	r.False(result, "Path should not exist")
}

func TestDirExistsButIsFile(t *testing.T) {
	r := require.New(t)

	// Given a test folder with a file in it
	tmpDir, err := os.MkdirTemp("", "")
	r.NoError(err, "Error creating temp folder")

	// Make sure we always clean up our temp folder
	defer func() {
		err = os.RemoveAll(tmpDir)
		r.NoError(err, "Error deleting temp folder")
	}()

	var filename = path.Join(tmpDir, "fubar")
	_, err = os.Create(filename)
	r.NoError(err, "Failed to create test file")

	// When we check to see if that file exists as a directory
	result := DirExists(filename)

	// Then, the response should be False because the dir already exists
	r.False(result, "Path should not exist")

}

func TestInvalidDirExists(t *testing.T) {
	// The stat method doesn't accept strings with nested null characters in it
	// but we expect that error to get swallowed and just return a false result
	result := DirExists(".\x00asdf")
	require.False(t, result, "Path with invalid char should fail dir check")
}
