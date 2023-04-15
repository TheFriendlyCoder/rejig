package applicationOptions

import (
	"path"

	"github.com/TheFriendlyCoder/rejigger/lib"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"
)

// -=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=
//																				  Constants

const manifestFileName = ".rejig.yml"

// -=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=
//														    			 TemplateSourceType

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

// toString Converts the value from our enumeration to a string representation
func (t *TemplateSourceType) toString() string {
	switch *t {
	case TstGit:
		return "git"
	case TstLocal:
		return "local"
	case TstUndefined:
		fallthrough
	case TstUnknown:
		fallthrough
	default:
		return ""
	}
}

// fromString populates our enumeration from an arbitrary character string
func (t *TemplateSourceType) fromString(value string) {
	switch value {
	case "git":
		*t = TstGit
	case "local":
		*t = TstLocal
	case "":
		*t = TstUndefined
	default:
		*t = TstUnknown
	}
}

// UnmarshalYAML decodes values for our enumeration from YAML content
func (t *TemplateSourceType) UnmarshalYAML(value *yaml.Node) error {
	var temp string
	if err := value.Decode(&temp); err != nil {
		return errors.Wrap(err, "Unable to parse template source type: "+value.Value)
	}
	t.fromString(temp)
	return nil
}

// -=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=
//																			TemplateOptions

// TemplateOptions metadata describing the source location for a source template
type TemplateOptions struct {
	// Type identifier describing the protocol to use when retrieving template content
	Type TemplateSourceType
	// Source Path or URL where the source template can be found
	Source string // TODO: make this member private
	// Name friendly name associated with the template. Used when referring to the template
	// from the command line
	Name string
	// Root optional sub-folder from the source root where the inventory resides
	// If not provided, the template root is assumed to be the root of the source folder
	// Typically used by Git type templates that are stored in a subfolder of a remote repo
	Root string
}

// GetName friendly name associated with the template. Used when referring to the template
// from the command line
func (t *TemplateOptions) GetName() string {
	return t.Name
}

// GetType identifier describing the protocol to use when retrieving template content
func (t *TemplateOptions) GetType() TemplateSourceType {
	return t.Type
}

// GetFilesystem Gets a virtual filesystem pre-loaded to point to the file system for the template
func (t *TemplateOptions) GetFilesystem() (afero.Fs, error) {
	switch t.Type {
	case TstLocal:
		return afero.NewOsFs(), nil
	case TstGit:
		return lib.GetGitFilesystem(t.Source)
	case TstUnknown:
		fallthrough
	case TstUndefined:
		fallthrough
	default:
		panic("Unsupported template source type " + t.Type.toString())
	}
}

// GetRoot gets the path to the root folder of the virtual file system associated with
// this template
func (t *TemplateOptions) GetRoot() string {
	switch t.Type {
	case TstLocal:
		if t.Root == "" {
			return t.Source
		} else {
			return path.Join(t.Source, t.Root)
		}

	case TstGit:
		if t.Root == "" {
			return "."
		} else {
			return t.Root
		}
	case TstUnknown:
		fallthrough
	case TstUndefined:
		fallthrough
	default:
		panic("Unsupported template source type " + t.Type.toString())
	}
}

// GetManifestFile gets the path, relative to the filesystem root, where the template manifest
// file for this template is found
func (t *TemplateOptions) GetManifestFile() string {
	return path.Join(t.GetRoot(), manifestFileName)
}
