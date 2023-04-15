package applicationOptions

import (
	"path"

	"github.com/TheFriendlyCoder/rejigger/lib"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

// -=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=
//																				  Constants

const inventoryFileName = ".rejig.inv.yml"

// -=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=
//														    			InventorySourceType

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

// toString Converts the value from our enumeration to a string representation
func (i *InventorySourceType) toString() string {
	switch *i {
	case IstGit:
		return "git"
	case IstLocal:
		return "local"
	case IstUndefined:
		fallthrough
	case IstUnknown:
		fallthrough
	default:
		return ""
	}
}

// fromString populates our enumeration from an arbitrary character string
func (i *InventorySourceType) fromString(value string) {
	switch value {
	case "git":
		*i = IstGit
	case "local":
		*i = IstLocal
	case "":
		*i = IstUndefined
	default:
		*i = IstUnknown
	}
}

// -=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=
//																		   InventoryOptions

// InventoryOptions configuration parameters for an inventory
type InventoryOptions struct {
	// Type identifier describing the protocol to use when retrieving template inventories
	Type InventorySourceType
	// Source path or URL to the inventory
	Source string
	// Namespace prefix to add to all templates contained in this inventory
	Namespace string
}

// GetNamespace friendly name associated with the namespace. Used when referring to templates
// in this inventory from the command line
func (i *InventoryOptions) GetNamespace() string {
	return i.Namespace
}

// GetSource Path or URL where the source inventory can be found
func (i *InventoryOptions) GetSource() string {
	return i.Source
}

// GetType identifier describing the protocol to use when retrieving inventory content
func (i *InventoryOptions) GetType() InventorySourceType {
	return i.Type
}

// GetFilesystem loads an appropriate virtual file system to allow reading of the
// inventory data based on the source type, in a filesystem agnostic way
func (i *InventoryOptions) GetFilesystem() (afero.Fs, error) {
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
		panic("should never happen: unsupported inventory type: " + i.Type.toString())
	}
}

// GetRoot gets the path to the root folder of the virtual file system associated with
// this template
func (i *InventoryOptions) GetRoot() string {
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
		panic("Unsupported template source type " + i.Type.toString())
	}
}

// GetInventoryFile gets the path, relative to the filesystem root, where the inventory definition
// file is found
func (i *InventoryOptions) GetInventoryFile() string {
	return path.Join(i.GetRoot(), inventoryFileName)
}

// GetTemplateDefinitions gets a list of all templates defined in this inventory
func (i *InventoryOptions) GetTemplateDefinitions() ([]TemplateOptions, error) {
	fileSystem, err := i.GetFilesystem()
	if err != nil {
		return nil, err
	}

	// Read in the inventory file
	inventoryPath := i.GetInventoryFile()
	_, err = fileSystem.Stat(inventoryPath)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// TODO: Cache these results so they can be reused
	inventory, err := parseInventory(fileSystem, inventoryPath)
	if err != nil {
		return nil, err
	}

	// Rework the template metadata, inheriting properties from the
	// parent inventory
	retval := make([]TemplateOptions, 0, len(inventory.Templates))
	for _, curTemplate := range inventory.Templates {
		temp := TemplateOptions{
			Type:   TemplateSourceType(i.Type),
			Root:   curTemplate.Source,
			Source: i.Source,
			// TODO: Consider setting name to i.Namespace + "." + curTemplate.Name
			Name: curTemplate.Name,
		}
		retval = append(retval, temp)
	}

	// TODO: consider how/where I should be validating the contents of the templates
	// TODO: consider how to handle duplicate templates
	// TODO: consider forcing the "name" field to be unique in each inventory
	// TODO: consider pre-pending namespace to name to ensure uniqueness
	// TODO: probably should not allow "local" type for templates defined in a Git inventory
	return retval, nil
}
