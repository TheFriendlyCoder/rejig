package lib

import (
	"os"
	"path"
	"testing"
)

func TestDirExists(t *testing.T) {
	result := DirExists(".")
	if !result {
		t.Fatal("Path should exist")
	}
}

func TestDirNotExists(t *testing.T) {
	// Given an empty temp folder
	tmpDir, err := os.MkdirTemp("", "")
	if err != nil {
		t.Fatalf("Error creating temp folder: %s", err)
	}

	// When we request a check on a file system object that doesn't exist
	result := DirExists(path.Join(tmpDir, "fubar"))

	// Then we expect our directory query to return false
	if result {
		t.Error("Path should not exist")
	}

	// Clean up test data
	err = os.Remove(tmpDir)
	if err != nil {
		t.Fatalf("Error deleting temp folder: %s", err)
	}
}

func TestDirExistsButIsFile(t *testing.T) {
	// Given a test folder with a file in it
	tmpDir, err := os.MkdirTemp("", "")
	if err != nil {
		t.Fatalf("Error creating temp folder: %s", err)
	}
	var filename = path.Join(tmpDir, "fubar")

	file, err := os.Create(filename)
	if err != nil {
		t.Fatalf("Failed to create test file: %s", err)
	}

	// When we check to see if that file exists as a directory
	result := DirExists(filename)

	// Then we should get a false return value
	if result {
		t.Error("Path should not exist")
	}

	// Cleanup test data
	err = file.Close()
	if err != nil {
		t.Errorf("Error closing test file %s", err)
	}
	err = os.RemoveAll(tmpDir)
	if err != nil {
		t.Fatalf("Error deleting temp folder: %s", err)
	}
}

func TestInvalidDirExists(t *testing.T) {
	// The stat method doesn't accept strings with nested null characters in it
	// but we expect that error to get swallowed and just return a false result
	result := DirExists(".\x00asdf")
	if result {
		t.Fatal("Path should not exist")
	}
}
