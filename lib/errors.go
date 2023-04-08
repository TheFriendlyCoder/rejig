package lib

import (
	"fmt"
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

var AOInvalidSourceTypeError = fmt.Errorf("unsupported template source type")
var AOTemplateOptionsDecodeError = fmt.Errorf("unable to decode template options")
var AOInventoryOptionsDecodeError = fmt.Errorf("unable to decode inventory options")
