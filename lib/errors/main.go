package errors

import (
	"strings"

	"github.com/pkg/errors"
)

// -=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=
//										AppOptionsError

type appOptionsError struct {
	Messages []string
}

func (e appOptionsError) Error() string {
	if len(e.Messages) == 1 {
		return "Failed to parse application option: " + e.Messages[0]
	}
	retval := strings.Join(e.Messages, "\n\t")
	return "Failed to parse application options:\n\t" + retval
}

func (e appOptionsError) Is(other error) bool {
	var newVal appOptionsError
	if errors.As(other, &newVal) {
		if len(e.Messages) != len(newVal.Messages) {
			return false
		}
		for i, curMessage := range newVal.Messages {
			if e.Messages[i] != curMessage {
				return false
			}
		}
		return true
	}
	return false
}

func NewAppOptionsError(messages []string) error {
	// TODO: Every place i this file we use WithStack we should find a way
	//       to unroll the stack frames and remove the call to the construction
	//		 methods since they are useless to the stack debugging
	return errors.WithStack(appOptionsError{messages})
}

// -=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=
//									PathError

type PathErrorTypes int64

const (
	PeUndefined PathErrorTypes = iota
	PePathNotFound
	PePathNotEmpty
	PeFileNotFound
)

type pathError struct {
	Path      string
	ErrorType PathErrorTypes
}

func (e pathError) Error() string {
	var retval string
	switch e.ErrorType {
	case PePathNotFound:
		retval = "Path not found: " + e.Path
	case PePathNotEmpty:
		retval = "Path must be empty: " + e.Path
	case PeFileNotFound:
		retval = "File not found: " + e.Path
	case PeUndefined:
		panic("Unsupported Path error")
	}
	return retval
}

func (e pathError) Is(other error) bool {
	var newVal pathError
	if errors.As(other, &newVal) {
		if e.Path != newVal.Path {
			return false
		}
		if e.ErrorType != newVal.ErrorType {
			return false
		}
		return true
	}
	return false
}

func NewPathError(path string, errType PathErrorTypes) error {
	return errors.WithStack(pathError{path, errType})
}

// -=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=
//										UnknownTemplateError

type unknownTemplateError struct {
	TemplateAlias string
}

func (e unknownTemplateError) Error() string {
	return "Template not found in application inventory: " + e.TemplateAlias
}

func (e unknownTemplateError) Is(other error) bool {
	var newVal unknownTemplateError
	if errors.As(other, &newVal) {
		return e.TemplateAlias == newVal.TemplateAlias
	}
	return false
}

func NewUnknownTemplateError(alias string) error {
	return errors.WithStack(unknownTemplateError{alias})
}

// -=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=
//										InternalError

type internalError struct {
	Message string
}

func (e internalError) Error() string {
	return "Internal implementation error: " + e.Message
}

func (e internalError) Is(other error) bool {
	var newVal internalError
	if errors.As(other, &newVal) {
		return e.Message == newVal.Message
	}
	return false
}

func NewInternalError(message string) error {
	return errors.WithStack(internalError{message})
}

func CommandContextNotDefined() error {
	return NewInternalError("Command context not properly initialized")
}

// -=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=
//														SimpleError

type simpleError struct {
	Message string
}

func NewSimpleError(message string) error {
	return errors.WithStack(simpleError{message})
}
func (e simpleError) Error() string {
	return e.Message
}

func (e simpleError) Is(other error) bool {
	var newVal simpleError
	if errors.As(other, &newVal) {
		return e.Message == newVal.Message
	}
	return false
}

// -=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=
//										Misc Errors

func AOInvalidSourceTypeError() error {
	return NewSimpleError("unsupported template source type")
}

func AOTemplateOptionsDecodeError() error {
	return NewSimpleError("unable to decode template options")
}

func AOInventoryOptionsDecodeError() error {
	return NewSimpleError("unable to decode inventory options")
}
