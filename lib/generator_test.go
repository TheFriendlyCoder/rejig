package lib

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"path"
	"path/filepath"
	"testing"
)

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

func Test_basicGenerator(t *testing.T) {
	r := require.New(t)
	a := assert.New(t)

	srcPath, err := sampleProj("simple")
	r.NoError(err, "Failed to locate sample project")

	// Given an empty temp folder
	tmpDir, err := os.MkdirTemp("", "")
	r.NoError(err, "Error creating temp folder")

	// Make sure we always clean up our temp folder
	defer func() {
		r.NoError(os.RemoveAll(tmpDir), "Error deleting temp folder")
	}()

	// We attempt to run the generator
	context := map[string]any{
		"project_name": "MyProj",
		"version":      "1.6.9",
	}
	err = Generate(*srcPath, tmpDir, context)

	r.NoError(err, "Failed to run generator")
	fmt.Println(tmpDir)
	a.DirExists(filepath.Join(tmpDir, "MyProj"))
	a.NoFileExists(filepath.Join(tmpDir, ".rejig.yml"))
	a.FileExists(filepath.Join(tmpDir, ".gitignore"))
	a.FileExists(filepath.Join(tmpDir, "version.txt"))
	a.FileExists(filepath.Join(tmpDir, "MyProj", "main.txt"))
}
