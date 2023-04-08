package templateManager

import (
	"github.com/hashicorp/go-version"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
	"os"
	"path"
	"testing"
)

func Test_parseManifest(t *testing.T) {
	r := require.New(t)
	a := assert.New(t)

	// Given a sample manifest file
	srcFile := sampleDataFile("simple_manifest.yml")

	// When we parse it
	manifest, err := parseManifest(afero.NewOsFs(), srcFile)

	// We expect no errors
	r.NoError(err)

	// and the parsed data should look as expected
	expSchemaVersion, err := version.NewVersion("1.0")
	r.NoError(err)
	expJiggerVersion, err := version.NewVersion("0.0.1")
	r.NoError(err)
	expTemplateVersion, err := version.NewVersion("2.0")
	r.NoError(err)

	// Validate version info
	a.Equal(*expSchemaVersion, manifest.Versions.Schema)
	a.Equal(*expJiggerVersion, manifest.Versions.Jigger)
	a.Equal(*expTemplateVersion, manifest.Versions.Template)

	// Validate template metadata
	a.Equal(2, len(manifest.Template.Args))
	a.Equal("project_name", manifest.Template.Args[0].Name)
	a.Equal("Name of the source code project", manifest.Template.Args[0].Description)
	a.Equal("version", manifest.Template.Args[1].Name)
	a.Equal("Initial version number for the project", manifest.Template.Args[1].Description)

	// Validate unparsed values
	a.Equal(1, len(manifest.MiscParams))
	a.Contains(manifest.MiscParams, "misc")
	a.Equal("fubar", manifest.MiscParams["misc"])
}

func Test_parseManifestNotExist(t *testing.T) {
	r := require.New(t)
	a := assert.New(t)

	// Given an empty temp folder
	tmpDir, err := os.MkdirTemp("", "")
	r.NoError(err)
	defer os.RemoveAll(tmpDir)

	_, err = parseManifest(afero.NewOsFs(), path.Join(tmpDir, "fubar.yml"))
	a.Error(err)
}

func Test_parseManifestInvalidYAML(t *testing.T) {
	r := require.New(t)
	a := assert.New(t)

	// Given an empty temp folder
	tmpDir, err := os.MkdirTemp("", "")
	r.NoError(err)
	defer os.RemoveAll(tmpDir)

	// and a sample config file that contains non-yaml data
	samplefile := path.Join(tmpDir, "fubar.yml")
	srcfile, err := os.Create(samplefile)
	r.NoError(err)
	_, err = srcfile.WriteString("This is not compatible yaml")
	r.NoError(err)
	r.NoError(srcfile.Close())

	// when we try to parse the file
	_, err = parseManifest(afero.NewOsFs(), samplefile)

	// then we expect error results
	emptyTypeErr := yaml.TypeError{}
	a.Contains(err.Error(), emptyTypeErr.Error())
}

func Test_parseManifestInvalidTemplateArgs(t *testing.T) {
	a := assert.New(t)

	srcFile := sampleDataFile("simple_manifest_with_invalid_args.yml")

	_, err := parseManifest(afero.NewOsFs(), srcFile)
	// TODO: Find some way to make error reporting here more user friendly
	//		 may require a different YAML parsing library
	// https://github.com/go-yaml/yaml/pull/901
	a.Error(err)
}
