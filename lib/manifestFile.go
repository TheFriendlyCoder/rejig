package lib

import (
	"github.com/hashicorp/go-version"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	"os"
)

// SchemaData stores metadata about a template manifest used by the rejigger tool
type SchemaData struct {
	// Version schema version describing the format of the manifest file and its contents
	Version version.Version `yaml:"version"`
	// JiggerVersion the minimum version of the Rejigger app needed to process to the template
	JiggerVersion version.Version `yaml:"rejigger_version"`
}

// ManifestData parsed content of the manifest file associated with a template
type ManifestData struct {
	Schema         SchemaData             `yaml:"schema"`
	TemplateParams map[string]interface{} `yaml:"-,flow"`
}

// UnmarshalYAML custom YAML decoding method compatible with the YAML parsing library
func (m *ManifestData) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var schemaFields struct {
		Schema SchemaData `yaml:"schema"`
	}
	if err := unmarshal(&schemaFields); err != nil {
		return err
	}
	var remaining map[string]interface{}
	if err := unmarshal(&remaining); err != nil {
		return err
	}
	m.Schema = schemaFields.Schema
	delete(remaining, "schema")
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
