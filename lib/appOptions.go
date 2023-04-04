package lib

import "fmt"

// TemplateSourceType enum for all supported source locations for loading templates
type TemplateSourceType int64

const (
	TST_UNDEFINED TemplateSourceType = iota
	TST_GIT
)

// TemplateOptions metadata describing the source location for a source template
type TemplateOptions struct {
	// Type identifier describing the protocol to use when retrieving template content
	Type TemplateSourceType
	// Source Path or URL where the source template can be found
	Source string
	// Folder sub-folder within the source location where the template source is defined
	// defaults to the root folder of the source location
	Folder string
	// Name friendly name associated with the template. Used when referring to the template
	// from the command line
	Name string
}

// AppOptions parsed config options supported by the app
type AppOptions struct {
	// Templates 0 or more sources where template projects are to be found
	Templates []TemplateOptions
}

// Validate checks the contents of the parsed application options to make sure they
// meet the requirements for the application
func (a *AppOptions) Validate() error {
	var messages []string
	for i, curTemplate := range a.Templates {
		if len(curTemplate.Name) == 0 {
			messages = append(messages, fmt.Sprintf("template %d name is PE_UNDEFINED", i))
		}
		if curTemplate.Type == TST_UNDEFINED {
			messages = append(messages, fmt.Sprintf("template %d type is PE_UNDEFINED", i))
		}
		if len(curTemplate.Source) == 0 {
			messages = append(messages, fmt.Sprintf("template %d source is PE_UNDEFINED", i))
		}
	}
	if len(messages) == 0 {
		return nil
	}
	return AppOptionsError{Messages: messages}
}
