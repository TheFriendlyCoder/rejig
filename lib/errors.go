package lib

import (
	"fmt"
	"strings"
)

// -=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=
//									PathError

type PathErrorTypes int64

const (
	PE_UNDEFINED PathErrorTypes = iota
	PE_PATH_NOT_FOUND
	PE_PATH_NOT_EMPTY
	PE_FILE_NOT_FOUND
)

type PathError struct {
	Path      string
	ErrorType PathErrorTypes
}

func (p PathError) Error() string {
	var retval string
	switch p.ErrorType {
	case PE_PATH_NOT_FOUND:
		retval = "Path not found: " + p.Path
	case PE_PATH_NOT_EMPTY:
		retval = "Path must be empty: " + p.Path
	case PE_FILE_NOT_FOUND:
		retval = "File not found: " + p.Path
	case PE_UNDEFINED:
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
//										TemplateNotFoundError

type TemplateNotFoundError struct {
	TemplateAlias string
}

func (e TemplateNotFoundError) Error() string {
	return "Template not found in inventory: " + e.TemplateAlias
}

// -=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=
//										Misc Errors

var APP_OPTIONS_INVALID_SOURCE_TYPE_ERROR = fmt.Errorf("Unsupported template source type")
var APP_OPTIONS_DECODE_ERROR = fmt.Errorf("Unable to decode template options")
