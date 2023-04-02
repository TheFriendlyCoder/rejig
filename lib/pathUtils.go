package lib

import (
	"log"
	"os"
)

// DirExists checks to see if the path given points to a folder that exists
func DirExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	if err != nil {
		log.Println(err.Error())
		return false
	}
	return info.IsDir()

}
