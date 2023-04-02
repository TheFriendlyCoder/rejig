package lib

import (
	"github.com/hashicorp/go-version"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	"os"
)

// VersionData stores version information about the template
type VersionData struct {
	// Schema version describing the format of the manifest file and its contents
	Schema version.Version `yaml:"schema"`
	// Jigger the minimum version of the Rejigger app needed to process to the template
	Jigger version.Version `yaml:"rejigger"`
	// Template the version number associated with the current template
	Template version.Version `yaml:"template"`
}

// ManifestData parsed content of the manifest file associated with a template
type ManifestData struct {
	Versions       VersionData            `yaml:"versions"`
	TemplateParams map[string]interface{} `yaml:"-,flow"`
}

// UnmarshalYAML custom YAML decoding method compatible with the YAML parsing library
func (m *ManifestData) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var versionFields struct {
		Versions VersionData `yaml:"versions"`
	}
	if err := unmarshal(&versionFields); err != nil {
		return err
	}
	var remaining map[string]interface{}
	if err := unmarshal(&remaining); err != nil {
		return err
	}
	m.Versions = versionFields.Versions
	delete(remaining, "versions")
	m.TemplateParams = remaining
	return nil
}

// ParseManifest parses a template manifest file and returns a reference to
// the parsed representation of the contents of the file
func ParseManifest(path string) (*ManifestData, error) {
	buf, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to open manifest file")
	}

	retval := &ManifestData{}
	err = yaml.Unmarshal(buf, retval)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to parse YAML content from manifest file")
	}

	return retval, nil
}
