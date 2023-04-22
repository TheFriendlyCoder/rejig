package applicationOptions

import (
	"os"
	"path"
	"regexp"
	"strings"

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
	Source string `yaml:"source"` // TODO: make this member private
	// Name friendly name associated with the template. Used when referring to the template
	// from the command line
	Name string `yaml:"name"`
	// SubDir optional sub-directory under the template Source location where the template
	// definition is found. If not provided, the template is expected to exist in the root
	// folder of the Source location
	SubDir string `yaml:"subdir"`
	// Exclusions set of 0 or more regular expressions defining files to be excluded from
	// template processing
	Exclusions []string `yaml:"exclusions"`

	// regexExclusions cache of pre-compiled regular expressions built from the Exclusions list
	regexExclusions []*regexp.Regexp
}

// buildRegex populates the regexExclusions cache in the TemplateOptions struct
// performs a no-op if the cache has already been populated
func (t *TemplateOptions) buildRegex() error {
	if t.regexExclusions != nil {
		return nil
	}

	ignoreList := make([]*regexp.Regexp, 0, len(t.Exclusions))
	for _, curExpr := range t.Exclusions {
		curIgnore, err := regexp.Compile(curExpr)
		if err != nil {
			return errors.WithStack(err)
		}
		ignoreList = append(ignoreList, curIgnore)
	}
	t.regexExclusions = ignoreList
	return nil
}

// IsFileExcluded returns true if the given file path should be excluded based on the
// exclusion rules provided by the template options, false if not
func (t *TemplateOptions) IsFileExcluded(curFile string) bool {
	// TODO: move this regex builder into object creation to avoid having a panic here
	err := t.buildRegex()
	if err != nil {
		panic(err.Error())
	}

	for _, curIgnore := range t.regexExclusions {
		if curIgnore.Match([]byte(curFile)) {
			return true
		}
	}
	return false
}

// GetSource gets the path to the source folder where the template definition lives
func (t *TemplateOptions) GetSource() string {
	if t.Type != TstLocal {
		return t.Source
	}
	if !strings.HasPrefix(t.Source, "~/") {
		return t.Source
	}
	// TODO: Put the home dir in the global settings / context, and validate it on application start,
	// 		 that way I don't need to worry about it anywhere else
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic("Users home directory not defined")
	}
	// TODO: Consider adding support for env var expansion using os.ExpandEnv()
	return path.Join(homeDir, t.Source[1:])
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

// GetProjectRoot gets the path to the root folder of the virtual file system associated with
// this template
func (t *TemplateOptions) GetProjectRoot() string {
	switch t.Type {
	case TstLocal:
		if t.SubDir == "" {
			return t.GetSource()
		} else {
			return path.Join(t.GetSource(), t.SubDir)
		}

	case TstGit:
		if t.SubDir == "" {
			return "."
		} else {
			return t.SubDir
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
	return path.Join(t.GetProjectRoot(), manifestFileName)
}
