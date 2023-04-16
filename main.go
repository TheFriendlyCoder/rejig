package main

import (
	"os"

	"github.com/TheFriendlyCoder/rejigger/cmd"
)

func main() {
	rootCmd := cmd.RootCmd()
	err := cmd.Execute(&rootCmd)
	if err != nil {
		os.Exit(1)
	}
}
