package applicationOptions

import (
	"fmt"

	e "github.com/TheFriendlyCoder/rejigger/lib/errors"
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
