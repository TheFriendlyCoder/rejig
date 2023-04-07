package lib

import (
	"fmt"
	"strings"
)

// -=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=
//									PathError

type PathErrorTypes int64

const (
	PeUndefined PathErrorTypes = iota
	PePathNotFound
	PePathNotEmpty
	PeFileNotFound
)

type PathError struct {
	Path      string
	ErrorType PathErrorTypes
}

func (p PathError) Error() string {
	var retval string
	switch p.ErrorType {
	case PePathNotFound:
		retval = "Path not found: " + p.Path
	case PePathNotEmpty:
		retval = "Path must be empty: " + p.Path
	case PeFileNotFound:
		retval = "File not found: " + p.Path
	case PeUndefined:
		panic("Unsupported Path error")
	}
	return retval
}

// -=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=
//										AppOptionsError

type AppOptionsError struct {
	Messages []string
}

func (e AppOptionsError) Error() string {
	if len(e.Messages) == 1 {
		return "Failed to parse application options: " + e.Messages[0]
	}
	retval := strings.Join(e.Messages, "\n\t")
	return "Failed to parse application options:\n\t" + retval
}

// -=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=
//										UnknownTemplateError

type UnknownTemplateError struct {
	TemplateAlias string
}

func (e UnknownTemplateError) Error() string {
	return "Template not found in application inventory: " + e.TemplateAlias
}

// -=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=
//										InternalError

type InternalError struct {
	Message string
}

func (e InternalError) Error() string {
	return "Internal implementation error: " + e.Message
}

var CommandContextNotDefined = InternalError{"Command context not properly initialized"}

// -=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=
//										Misc Errors

var AppOptionsInvalidSourceTypeError = fmt.Errorf("unsupported template source type")
var AppOptionsDecodeError = fmt.Errorf("unable to decode template options")
