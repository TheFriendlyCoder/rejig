package lib

import "fmt"

// TemplateSourceType enum for all supported source locations for loading templates
type TemplateSourceType int64

const (
	// TstUndefined No template type defined in template config
	TstUndefined TemplateSourceType = iota
	// TstLocal Template source is stored on the local file system
	TstLocal
	// TstGit Template source is stored in a Git repository
	TstGit
)

// TemplateOptions metadata describing the source location for a source template
type TemplateOptions struct {
	// Type identifier describing the protocol to use when retrieving template content
	Type TemplateSourceType
	// Source Path or URL where the source template can be found
	Source string
	// Alias friendly name associated with the template. Used when referring to the template
	// from the command line
	Alias string
}

// AppOptions parsed config options supported by the app
type AppOptions struct {
	// Templates 0 or more sources where template projects are to be found
	Templates []TemplateOptions
}

// Validate checks the contents of the parsed application options to make sure they
// meet the requirements for the application
func (a AppOptions) Validate() error {
	var messages []string
	for i, curTemplate := range a.Templates {
		if len(curTemplate.Alias) == 0 {
			messages = append(messages, fmt.Sprintf("template %d alias is undefined", i))
		}
		if curTemplate.Type == TstUndefined {
			messages = append(messages, fmt.Sprintf("template %d type is undefined", i))
		}
		if len(curTemplate.Source) == 0 {
			messages = append(messages, fmt.Sprintf("template %d source is undefined", i))
		}
	}
	if len(messages) == 0 {
		return nil
	}
	return AppOptionsError{Messages: messages}
}
