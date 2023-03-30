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
		// TODO: convert this print to a logger
		fmt.Println(err.Error())
		return false
	}
	return info.IsDir()

}
