package applicationOptions

import (
	"fmt"

	e "github.com/TheFriendlyCoder/rejigger/lib/errors"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
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

func (t TemplateSourceType) ToString() string {
	switch t {
	case TstUndefined:
		return "Undefined"
	case TstLocal:
		return "Local"
	case TstGit:
		return "Git"
	case TstUnknown:
		fallthrough
	default:
		return "Unknown"
	}
}

// toString Converts the value from our enumeration to a string representation
func (t *TemplateSourceType) toString() string {
	switch *t {
	case TstGit:
		return "git"
	case TstLocal:
		return "local"
	case TstUndefined:
		// TODO: See if I should handle these values differently
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

// decodeTemplateOptions decodes raw YAMl data into proper parsed template options
func decodeTemplateOptions(raw interface{}) (map[string]interface{}, error) {
	// Map the "type" field from a character string format to an enumerated type
	templateData, ok := raw.(map[string]interface{})
	if !ok {
		return nil, e.AOTemplateOptionsDecodeError()
	}

	var newVal TemplateSourceType
	switch templateData["type"] {
	case nil:
		return nil, e.AOInvalidSourceTypeError()
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

// validateTemplates checks to make sure the template options in our application
// options meets the application requirements
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
