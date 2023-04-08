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

// InventorySourceType enum for all supported source locations for template inventories
type InventorySourceType int64

const (
	// IstUndefined No inventory type defined
	IstUndefined InventorySourceType = iota
	// IstLocal Inventory source is stored on the local file system
	IstLocal
	// IstGit Inventory source is stored in a Git repository
	IstGit
)

// InventoryOptions configuration parameters for an inventory
type InventoryOptions struct {
	// Type identifier describing the protocol to use when retrieving template inventories
	Type InventorySourceType
	// Source path or URL to the inventory
	Source string
	// Namespace prefix to add to all templates contained in this inventory
	Namespace string
}

// AppOptions parsed config options supported by the app
type AppOptions struct {
	// Templates 0 or more sources where template projects are to be found
	Templates   []TemplateOptions
	Inventories []InventoryOptions
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
	for i, curInventory := range a.Inventories {
		if len(curInventory.Namespace) == 0 {
			messages = append(messages, fmt.Sprintf("inventory %d namespace is undefined", i))
		}
		if curInventory.Type == IstUndefined {
			messages = append(messages, fmt.Sprintf("inventory %d type is undefined", i))
		}
		if len(curInventory.Source) == 0 {
			messages = append(messages, fmt.Sprintf("inventory %d source is undefined", i))
		}
	}
	if len(messages) == 0 {
		return nil
	}
	return AppOptionsError{Messages: messages}
}
