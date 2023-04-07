package lib

import (
	"bytes"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"
)

// unmodified compares the contents of 2 files and returns true if they are
// the identical
func unmodified(r *require.Assertions, file1 string, file2 string) bool {
	f1, err := os.ReadFile(file1)
	r.NoError(err)

	f2, err := os.ReadFile(file2)
	r.NoError(err)

	return bytes.Equal(f1, f2)
}

// contains checks for a certain character string in a file and returns
// true if it is found
func contains(r *require.Assertions, file string, pattern string) bool {
	contents, err := os.ReadFile(file)
	r.NoError(err)

	return strings.Contains(string(contents), pattern)
}

// sampleProj loads Path to a specific sample project to use for testing the generator logic
func sampleProj(projName string) (string, error) {
	retval, err := filepath.Abs(path.Join("..", "testProjects", projName))
	if err != nil {
		return retval, errors.Wrap(err, "Failed to generate absolute Path")
	}
	_, err = os.Stat(retval)
	if err != nil {
		return retval, errors.Wrap(err, "checking existence of test data file")
	}
	return retval, nil
}

func Test_basicGenerator(t *testing.T) {
	r := require.New(t)
	a := assert.New(t)

	srcPath, err := sampleProj("simple")
	r.NoError(err)

	// Given an empty temp folder
	tmpDir, err := os.MkdirTemp("", "")
	r.NoError(err)
	defer os.RemoveAll(tmpDir)

	// We attempt to run the generator
	expVersion := "1.6.9"
	expProj := "MyProj"
	context := map[string]any{
		"project_name": expProj,
		"version":      expVersion,
	}
	fs := afero.NewOsFs()
	err = Generate(fs, srcPath, tmpDir, context)

	r.NoError(err, "Failed to run generator")

	a.DirExists(filepath.Join(tmpDir, "MyProj"))
	a.NoFileExists(filepath.Join(tmpDir, ".rejig.yml"))

	exp := filepath.Join(srcPath, ".gitignore")
	act := filepath.Join(tmpDir, ".gitignore")
	a.FileExists(act)
	a.True(unmodified(r, exp, act))

	act = filepath.Join(tmpDir, "version.txt")
	a.FileExists(act)
	a.True(contains(r, act, expVersion))
	a.False(contains(r, act, "{{version}}"))

	act = filepath.Join(tmpDir, "MyProj", "main.txt")
	a.FileExists(act)
	a.True(contains(r, act, expProj))
	a.False(contains(r, act, "{{project_name}}"))

}
