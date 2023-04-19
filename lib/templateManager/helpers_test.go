package templateManager

import (
	"bytes"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

// getProjectDir Gets the path to a specific test project
func getProjectDir() string {
	projectName := "simple"
	retval := path.Join("..", "..", "testdata", "projects", projectName)
	var info, err = os.Stat(retval)
	if err != nil {
		panic("Critical test failure: unable to access test project " + projectName)
	}
	if !info.IsDir() {
		panic("Critical test failure: test project " + projectName + " appears to be a file")
	}
	retval, err = filepath.Abs(retval)
	if err != nil {
		panic("Critical test failure: unable to generate absolute path for test project " + projectName)
	}
	return retval
}

// getManifestFile loads the path to a sample data for from the test data folder
func getManifestFile(filename string) string {
	retval := path.Join("..", "..", "testdata", "manifests", filename)
	info, err := os.Stat(retval)
	if err != nil {
		panic("Critical test failure: unable to access test manifest " + filename)
	}
	if info.IsDir() {
		panic("Critical test failure: test manifest " + filename + " appears to be a directory")
	}
	return retval
}

// isUnmodified compares the contents of 2 files and returns true if they are
// the identical
func isUnmodified(fs afero.Fs, r *require.Assertions, file1 string, file2 string) bool {
	f1, err := afero.ReadFile(fs, file1)
	r.NoError(err)

	f2, err := os.ReadFile(file2)
	r.NoError(err)

	return bytes.Equal(f1, f2)
}

// fileContains checks for a certain character string in a file and returns
// true if it is found
func fileContains(r *require.Assertions, file string, pattern string) bool {
	contents, err := os.ReadFile(file)
	r.NoError(err)

	return strings.Contains(string(contents), pattern)
}
