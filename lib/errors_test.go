package lib

import (
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func Test_FilePathInErrorMessage(t *testing.T) {
	type data struct {
		errorType PathErrorTypes
	}

	// TODO: Convert to table test
	tests := []data{
		{errorType: PePathNotFound},
		{errorType: PePathNotEmpty},
		{errorType: PeFileNotFound},
	}

	for _, tt := range tests {
		r := require.New(t)

		expectedPath, err := os.Getwd()
		r.NoError(err)

		val := PathError{
			Path:      expectedPath,
			ErrorType: tt.errorType,
		}

		r.Contains(val.Error(), expectedPath)
	}
}

func Test_UnknownPathError(t *testing.T) {
	r := require.New(t)

	expectedPath, err := os.Getwd()
	r.NoError(err)

	val := PathError{
		Path: expectedPath,
	}

	r.Panics(func() { _ = val.Error() })
}
