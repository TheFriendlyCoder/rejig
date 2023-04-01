package cmd

import (
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func Test_FilePathInErrorMessage(t *testing.T) {
	type data struct {
		errorType pathErrorTypes
	}

	tests := []data{
		{errorType: pathNotFound},
		{errorType: pathNotEmpty},
		{errorType: fileNotFound},
	}

	for _, tt := range tests {
		r := require.New(t)

		expectedPath, err := os.Getwd()
		r.NoError(err, "Failed to get current working folder")

		val := pathError{
			path:      expectedPath,
			errorType: tt.errorType,
		}

		r.Contains(val.Error(), expectedPath, "Error message should include the name of the path that caused the error")
	}
}

func Test_UnknownPathError(t *testing.T) {
	r := require.New(t)

	expectedPath, err := os.Getwd()
	r.NoError(err, "Failed to get current working folder")

	val := pathError{
		path: expectedPath,
	}

	r.Panics(func() { _ = val.Error() }, "Attempting to reference an error on an unsupported type should panic")
}
