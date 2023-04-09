package applicationOptions

import (
	"fmt"

	e "github.com/TheFriendlyCoder/rejigger/lib/errors"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"
)

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

// InventoryData stores data parsed from a template inventory definition file
type InventoryData struct {
	// Templates list of templates defined in the inventory
	Templates []TemplateOptions `yaml:"templates"`
}

// decodeInventoryOptions decodes raw YAMl data into proper parsed inventory options
func decodeInventoryOptions(raw interface{}) (map[string]interface{}, error) {
	// Map the "type" field from a character string format to an enumerated type
	inventoryData, ok := raw.(map[string]interface{})
	if !ok {
		return nil, e.AOTemplateOptionsDecodeError()
	}

	var newVal InventorySourceType
	switch inventoryData["type"] {
	case nil:
		return nil, e.AOInventoryOptionsDecodeError()
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

// parseInventory parses a template manifest file and returns a reference to
// the parsed representation of the contents of the file
func parseInventory(srcFS afero.Fs, path string) (InventoryData, error) {
	// TODO: Write validator for inventory, and probably should not allow "local" type for
	//		 templates defined in a Git inventory
	var retval InventoryData
	buf, err := afero.ReadFile(srcFS, path)
	if err != nil {
		return retval, errors.Wrap(err, "Failed to open inventory file")
	}

	// TODO: Find some way to get "Strict" mode to work properly (aka: KnownFields in v3)
	//		https://github.com/go-yaml/yaml/issues/460
	//		https://github.com/go-yaml/yaml/issues/642
	err = yaml.Unmarshal(buf, &retval)
	if err != nil {
		return retval, errors.Wrap(err, "Failed to parse YAML content from inventory file")
	}

	return retval, nil
}
