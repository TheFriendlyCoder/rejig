package lib

import (
	"fmt"
	"os"
)

// DirExists checks to see if the path given points to a folder that exists
func DirExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	if err != nil {
		fmt.Println(err.Error())
		// TODO: consider what other errors can happen here
	}
	return info.IsDir()

}
