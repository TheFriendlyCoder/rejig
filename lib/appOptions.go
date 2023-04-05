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
func (a *AppOptions) Validate() error {
	var messages []string
	// TODO: make sure no 2 templates have the same name
	// TODO: parse templates into a map instead of a list, using the alias as the key
	for i, curTemplate := range a.Templates {
		if len(curTemplate.Alias) == 0 {
			messages = append(messages, fmt.Sprintf("template %d alias is undefined", i))
		}
		if curTemplate.Type == TST_UNDEFINED {
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

// GetTemplate looks up a specific template in the app inventory based on its name
// This method assumes the requested alias exists in the inventory, and will panic
// if this assumption is broken
func (a *AppOptions) GetTemplate(alias string) TemplateOptions {
	for _, curTemplate := range a.Templates {
		if curTemplate.Alias == alias {
			return curTemplate
		}
	}
	panic("Template " + alias + " not found")
}
