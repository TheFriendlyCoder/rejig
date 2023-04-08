package templateManager

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// testProjectDir Gets the path to a specific test project
func testProjectDir(projectName string) string {
	retval := path.Join("..", "..", "testProjects", projectName)
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

// sampleDataFile loads the path to a sample data for from the test data folder
func sampleDataFile(filename string) string {
	retval := path.Join("..", "testdata", filename)
	info, err := os.Stat(retval)
	if err != nil {
		panic("Critical test failure: unable to access test data " + filename)
	}
	if info.IsDir() {
		panic("Critical test failure: test file " + filename + " appears to be a directory")
	}
	return retval
}

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
