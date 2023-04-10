package errors

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_UnknownPathError(t *testing.T) {
	r := require.New(t)

	expectedPath, err := os.Getwd()
	r.NoError(err)

	val := pathError{
		Path: expectedPath,
	}

	r.Panics(func() { _ = val.Error() })
}

func Test_CheckErrorIsWorking(t *testing.T) {
	r := require.New(t)

	tests := map[string]struct {
		srcType  error
		destType error
	}{
		"Check appOptionsError": {
			srcType:  NewAppOptionsError([]string{"My Error", "Other Error"}),
			destType: NewAppOptionsError([]string{"My Error", "Other Error"}),
		},
		"Check pathError": {
			srcType:  NewPathError("My Path", PePathNotFound),
			destType: NewPathError("My Path", PePathNotFound),
		},
		"Check unknownTemplateError": {
			srcType:  NewUnknownTemplateError("My Template"),
			destType: NewUnknownTemplateError("My Template"),
		},
		"Check internalError": {
			srcType:  NewInternalError("My Error"),
			destType: NewInternalError("My Error"),
		},
		"Check CommandContextNotDefined": {
			srcType:  CommandContextNotDefined(),
			destType: CommandContextNotDefined(),
		},
		"Check simpleError": {
			srcType:  NewSimpleError("My Error"),
			destType: NewSimpleError("My Error"),
		},
		"Check aoInvalidSourceTypeError": {
			srcType:  AOInvalidSourceTypeError(),
			destType: AOInvalidSourceTypeError(),
		},
		"Check aoTemplateOptionsDecodeError": {
			srcType:  AOTemplateOptionsDecodeError(),
			destType: AOTemplateOptionsDecodeError(),
		},
		"Check aoInventoryOptionsDecodeError": {
			srcType:  AOInventoryOptionsDecodeError(),
			destType: AOInventoryOptionsDecodeError(),
		},
	}

	for name, data := range tests {
		t.Run(name, func(t *testing.T) {
			r.ErrorIs(data.srcType, data.destType)
		})
	}
}

func Test_CheckErrorMessages(t *testing.T) {
	r := require.New(t)

	tests := map[string]struct {
		srcType    error
		expMessage string
	}{
		"Check appOptionsError multiple messages": {
			srcType:    NewAppOptionsError([]string{"My Error", "Other Error"}),
			expMessage: "Failed to parse application options:\n\tMy Error\n\tOther Error",
		},
		"Check appOptionsError single message": {
			srcType:    NewAppOptionsError([]string{"My Error"}),
			expMessage: "Failed to parse application option: My Error",
		},
		"Check pathError path not found": {
			srcType:    NewPathError("My Path", PePathNotFound),
			expMessage: "Path not found: My Path",
		},
		"Check pathError file not found": {
			srcType:    NewPathError("My Path", PeFileNotFound),
			expMessage: "File not found: My Path",
		},
		"Check pathError path not empty": {
			srcType:    NewPathError("My Path", PePathNotEmpty),
			expMessage: "Path must be empty: My Path",
		},
		"Check unknownTemplateError": {
			srcType:    NewUnknownTemplateError("My Template"),
			expMessage: "Template not found in application inventory: My Template",
		},
		"Check internalError": {
			srcType:    NewInternalError("My Error"),
			expMessage: "Internal implementation error: My Error",
		},
		"Check CommandContextNotDefined": {
			srcType:    CommandContextNotDefined(),
			expMessage: "Internal implementation error: Command context not properly initialized",
		},
		"Check simpleError": {
			srcType:    NewSimpleError("My Error"),
			expMessage: "My Error",
		},
		"Check aoInvalidSourceTypeError": {
			srcType:    AOInvalidSourceTypeError(),
			expMessage: "unsupported template source type",
		},
		"Check aoTemplateOptionsDecodeError": {
			srcType:    AOTemplateOptionsDecodeError(),
			expMessage: "unable to decode template options",
		},
		"Check aoInventoryOptionsDecodeError": {
			srcType:    AOInventoryOptionsDecodeError(),
			expMessage: "unable to decode inventory options",
		},
		"Check AOInvalidTemplateNameError": {
			srcType:    AOInvalidTemplateNameError(),
			expMessage: "invalid template name",
		},
	}

	for name, data := range tests {
		t.Run(name, func(t *testing.T) {
			r.Equal(data.expMessage, data.srcType.Error())
		})
	}
}

var fakeErr = errors.New("My fake error")

func Test_CheckNotErrorIsWorking(t *testing.T) {
	r := require.New(t)

	tests := map[string]struct {
		srcType  error
		destType error
	}{
		"Compare appOptionsError to fake error": {
			srcType:  NewAppOptionsError([]string{"My Error", "Other Error"}),
			destType: fakeErr,
		},
		"Compare appOptionsError to different messages": {
			srcType:  NewAppOptionsError([]string{"My Error", "Other Error"}),
			destType: NewAppOptionsError([]string{"My Error again", "Other Error again"}),
		},
		"Compare appOptionsError to subset of messages": {
			srcType:  NewAppOptionsError([]string{"My Error", "Other Error"}),
			destType: NewAppOptionsError([]string{"My Error"}),
		},
		"Compare pathError to fake error": {
			srcType:  NewPathError("My Path", PePathNotFound),
			destType: fakeErr,
		},
		"Compare pathError different error type": {
			srcType:  NewPathError("My Path", PePathNotFound),
			destType: NewPathError("My Path", PePathNotEmpty),
		},
		"Compare pathError different path": {
			srcType:  NewPathError("My Path", PePathNotFound),
			destType: NewPathError("Other Path", PePathNotFound),
		},
		"Compare unknownTemplateError to fake error": {
			srcType:  NewUnknownTemplateError("My Template"),
			destType: fakeErr,
		},
		"Compare internalError to fake error": {
			srcType:  NewInternalError("My Error"),
			destType: fakeErr,
		},
		"Compare simpleError to fake error": {
			srcType:  NewSimpleError("My Error"),
			destType: fakeErr,
		},
		"Compare aoInvalidSourceTypeError to fake error": {
			srcType:  AOInvalidSourceTypeError(),
			destType: fakeErr,
		},
		"Compare aoTemplateOptionsDecodeError to fake error": {
			srcType:  AOTemplateOptionsDecodeError(),
			destType: fakeErr,
		},
		"Compare aoInventoryOptionsDecodeError to fake error": {
			srcType:  AOInventoryOptionsDecodeError(),
			destType: fakeErr,
		},
		"Compare AOInvalidTemplateNameError to fake error": {
			srcType:  AOInvalidTemplateNameError(),
			destType: fakeErr,
		},
	}

	for name, data := range tests {
		t.Run(name, func(t *testing.T) {
			r.NotErrorIs(data.srcType, data.destType)
		})
	}
}
