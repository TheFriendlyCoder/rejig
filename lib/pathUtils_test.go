package lib

import (
	"testing"
)

func TestDirExists(t *testing.T) {
	result := DirExists(".")
	if !result {
		t.Fatal("Path should exist")
	}
}
