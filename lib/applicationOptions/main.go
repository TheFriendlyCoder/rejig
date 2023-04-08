package applicationOptions

import (
	"fmt"
	"github.com/TheFriendlyCoder/rejigger/lib"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"reflect"
)

// TemplateSourceType enum for all supported source locations for loading templates
type TemplateSourceType int64

const (
	// TstUndefined No template type defined in template config
	TstUndefined TemplateSourceType = iota
	// TstUnknown Template type provided but not currently supported
	TstUnknown
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
	// IstUnknown Inventory type provided but not currently supported
	IstUnknown
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

// FromViper constructs a new set of application options from a Viper config file
func FromViper(v *viper.Viper) (AppOptions, error) {
	var retval AppOptions
	err := errors.Wrap(v.Unmarshal(&retval, viper.DecodeHook(appOptionsDecoder())), "Failed parsing app options file")
	return retval, err
}

// New constructor for a new set of application options
func New() (AppOptions, error) {
	return AppOptions{}, nil
}

func (a AppOptions) validateTemplates() []string {
	var retval []string
	for i, curTemplate := range a.Templates {
		if len(curTemplate.Alias) == 0 {
			retval = append(retval, fmt.Sprintf("template %d alias is undefined", i))
		}
		if curTemplate.Type == TstUndefined {
			retval = append(retval, fmt.Sprintf("template %d type is undefined", i))
		} else if curTemplate.Type == TstUnknown {
			retval = append(retval, fmt.Sprintf("template %d type is not supported", i))
		}
		if len(curTemplate.Source) == 0 {
			retval = append(retval, fmt.Sprintf("template %d source is undefined", i))
		}
	}
	return retval
}
func (a AppOptions) validateInventory() []string {
	var retval []string
	for i, curInventory := range a.Inventories {
		if len(curInventory.Namespace) == 0 {
			retval = append(retval, fmt.Sprintf("inventory %d namespace is undefined", i))
		}
		if curInventory.Type == IstUndefined {
			retval = append(retval, fmt.Sprintf("inventory %d type is undefined", i))
		} else if curInventory.Type == IstUnknown {
			retval = append(retval, fmt.Sprintf("inventory %d type is not supported", i))
		}
		if len(curInventory.Source) == 0 {
			retval = append(retval, fmt.Sprintf("inventory %d source is undefined", i))
		}
	}
	return retval
}

// Validate checks the contents of the parsed application options to make sure they
// meet the requirements for the application
func (a AppOptions) Validate() error {
	var messages []string
	messages = append(messages, a.validateTemplates()...)
	messages = append(messages, a.validateInventory()...)
	if len(messages) == 0 {
		return nil
	}
	return AppOptionsError{Messages: messages}
}

// decodeTemplateOptions decodes raw YAMl data into proper parsed template options
func decodeTemplateOptions(raw interface{}) (map[string]interface{}, error) {
	// Map the "type" field from a character string format to an enumerated type
	templateData, ok := raw.(map[string]interface{})
	if !ok {
		return nil, lib.AOTemplateOptionsDecodeError
	}

	var newVal TemplateSourceType
	switch templateData["type"] {
	case nil:
		return nil, lib.AOInvalidSourceTypeError
	case "git":
		newVal = TstGit
	case "local":
		newVal = TstLocal
	default:
		newVal = TstUnknown
	}
	templateData["type"] = newVal
	return templateData, nil
}

// decodeInventoryOptions decodes raw YAMl data into proper parsed inventory options
func decodeInventoryOptions(raw interface{}) (map[string]interface{}, error) {
	// Map the "type" field from a character string format to an enumerated type
	inventoryData, ok := raw.(map[string]interface{})
	if !ok {
		return nil, lib.AOTemplateOptionsDecodeError
	}

	var newVal InventorySourceType
	switch inventoryData["type"] {
	case nil:
		return nil, lib.AOInventoryOptionsDecodeError
	case "git":
		newVal = IstGit
	case "local":
		newVal = IstLocal
	default:
		newVal = IstUnknown
	}
	inventoryData["type"] = newVal
	return inventoryData, nil
}

// appOptionsDecoder custom hook method used to translate raw config data into a structure
// that is easier to leverage in the application code
func appOptionsDecoder() mapstructure.DecodeHookFuncType {
	// Based on example found here:
	//		https://sagikazarmark.hu/blog/decoding-custom-formats-with-viper/
	return func(
		src reflect.Type,
		target reflect.Type,
		raw interface{},
	) (interface{}, error) {

		// TODO: Find a way to detect partial / incomplete parse matches
		// ie: if a template option is missing one field, viper won't map
		//		it to the correct type and it just gets ignored
		// TODO: Find a way to enable strict mode decoding here
		//		 that might work better
		if (target == reflect.TypeOf(TemplateOptions{})) {
			newData, err := decodeTemplateOptions(raw)
			return newData, errors.Wrap(err, "Error decoding template options")
		}

		if (target == reflect.TypeOf(InventoryOptions{})) {
			newData, err := decodeInventoryOptions(raw)
			return newData, errors.Wrap(err, "Error decoding inventory options")
		}
		return raw, nil

	}
}
