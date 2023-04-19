package internal

import (
	"os"
	"path"
	"path/filepath"
)

// GetProjectDir Gets the path to a specific test project
func GetProjectDir() string {
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
