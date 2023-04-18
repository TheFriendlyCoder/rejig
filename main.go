package main

import (
	"fmt"
	"os"

	"github.com/TheFriendlyCoder/rejigger/cmd"
	"github.com/TheFriendlyCoder/rejigger/lib"
)

func main() {
	rootCmd := cmd.RootCmd()
	err := cmd.Execute(&rootCmd)
	if err != nil {
		lib.SNF(fmt.Println(err.Error()))
		os.Exit(1)
	}
}
