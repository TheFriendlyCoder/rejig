package applicationOptions

import (
	"fmt"

	e "github.com/TheFriendlyCoder/rejigger/lib/errors"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"
)

// InventoryData stores data parsed from a template inventory definition file
type InventoryData struct {
	// Templates list of templates defined in the inventory
	Templates []TemplateOptions `yaml:"templates"`

	// TODO: consider if we should allow inventories to be nested
	//		 would have to watch out for circular imports
}

// decodeInventoryOptions decodes raw YAMl data into proper parsed inventory options
func decodeInventoryOptions(raw interface{}) (map[string]interface{}, error) {
	// Map the "type" field from a character string format to an enumerated type
	inventoryData, ok := raw.(map[string]interface{})
	if !ok {
		return nil, e.AOTemplateOptionsDecodeError()
	}

	var newVal InventorySourceType
	temp, ok := inventoryData["type"].(string)
	if !ok {
		return nil, e.AOTemplateOptionsDecodeError()
	}
	newVal.fromString(temp)
	inventoryData["type"] = newVal
	return inventoryData, nil
}

func (a AppOptions) validateInventory() []string {
	var retval []string
	allNames := map[string]int{}
	for i, curInventory := range a.Inventories {
		if len(curInventory.Namespace) == 0 {
			retval = append(retval, fmt.Sprintf("inventory %d namespace is undefined", i))
		}
		allNames[curInventory.Namespace] += 1

		if curInventory.Type == IstUndefined {
			retval = append(retval, fmt.Sprintf("inventory %d type is undefined", i))
		} else if curInventory.Type == IstUnknown {
			retval = append(retval, fmt.Sprintf("inventory %d type is not supported", i))
		}
		if len(curInventory.Source) == 0 {
			retval = append(retval, fmt.Sprintf("inventory %d source is undefined", i))
		}
	}

	// Make sure the inventory names are all unique
	for name, count := range allNames {
		if count > 1 {
			retval = append(retval, fmt.Sprintf("there are %d inventories named %s", count, name))
		}
	}
	return retval
}

// parseInventory parses a template manifest file and returns a reference to
// the parsed representation of the contents of the file
func parseInventory(srcFS afero.Fs, path string) (InventoryData, error) {
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
