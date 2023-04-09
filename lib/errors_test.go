package lib

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_FilePathInErrorMessage(t *testing.T) {
	r := require.New(t)

	tests := map[string]struct {
		errorType PathErrorTypes
	}{
		"Path not found": {
			errorType: PePathNotFound,
		},
		"Path not empty": {
			errorType: PePathNotEmpty,
		},
		"Path ": {
			errorType: PeFileNotFound,
		},
	}

	for name, data := range tests {
		t.Run(name, func(t *testing.T) {

			expectedPath, err := os.Getwd()
			r.NoError(err)

			val := PathError{
				Path:      expectedPath,
				ErrorType: data.errorType,
			}

			r.Contains(val.Error(), expectedPath)
		})
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

func Test_ExceptionErrors(t *testing.T) {
	r := require.New(t)

	tests := map[string]struct {
		curErr     error
		expMessage string
	}{
		"UnknownTemplateError": {
			curErr:     UnknownTemplateError{TemplateAlias: "MyAlias"},
			expMessage: "Template not found in application inventory: MyAlias",
		},
		"InternalError": {
			curErr:     InternalError{Message: "Cool Error"},
			expMessage: "Internal implementation error: Cool Error",
		},
	}

	for name, data := range tests {
		t.Run(name, func(t *testing.T) {
			r.Equal(data.expMessage, data.curErr.Error())
		})
	}
}
