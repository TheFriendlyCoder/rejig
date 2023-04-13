package applicationOptions

import (
	"fmt"
	"path"
	"strconv"

	"github.com/TheFriendlyCoder/rejigger/lib"
	e "github.com/TheFriendlyCoder/rejigger/lib/errors"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"
)

const inventoryFileName = ".rejig.inv.yml"

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
	// TODO: Write validator for inventory
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

// getFilesystem loads an appropriate virtual file system to allow reading of the
// inventory data based on the source type, in a filesystem agnostic way
func (i *InventoryOptions) getFilesystem() (afero.Fs, error) {
	switch i.Type {
	case IstLocal:
		return afero.NewOsFs(), nil
	case IstGit:
		return lib.GetGitFilesystem(i.Source)
	case IstUnknown:
		fallthrough
	case IstUndefined:
		fallthrough
	default:
		panic("should never happen: unsupported inventory type: " + strconv.FormatInt(int64(i.Type), 10))
	}
}

// getPath gets the path to the inventory file, taking into account the file system type
// used by the inventory
func (i *InventoryOptions) getPath() string {
	switch i.Type {
	case IstLocal:
		return i.Source
	case IstGit:
		return "."
	case IstUnknown:
		fallthrough
	case IstUndefined:
		fallthrough
	default:
		panic("should never happen: unsupported inventory type: " + strconv.FormatInt(int64(i.Type), 10))
	}
}

// GetTemplateDefinitions gets a list of all templates defined in this inventory
func (i *InventoryOptions) GetTemplateDefinitions() ([]TemplateOptions, error) {
	fileSystem, err := i.getFilesystem()
	if err != nil {
		return nil, err
	}
	rootDir := i.getPath()

	// Read in the inventory file
	inventoryPath := path.Join(rootDir, inventoryFileName)
	_, err = fileSystem.Stat(inventoryPath)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// TODO: Cache these results so they can be reused
	inventory, err := parseInventory(fileSystem, inventoryPath)
	if err != nil {
		return nil, err
	}

	// TODO: consider how/where I should be validating the contents of the templates
	// TODO: consider how to handle duplicate templates
	// TODO: consider forcing the "name" field to be unique in each inventory
	// TODO: consider pre-pending namespace to name to ensure uniqueness
	// TODO: probably should not allow "local" type for templates defined in a Git inventory
	return inventory.Templates, nil
}
