package lib

import (
	"github.com/hashicorp/go-version"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
	"os"
	"path"
	"testing"
)

// sampleData loads sample data for a test from the test data folder
func sampleData(filename string) (string, error) {
	retval := path.Join("testdata", filename)
	var _, err = os.Stat(retval)
	return retval, errors.Wrap(err, "checking existence of test data file")
}

func Test_parseManifest(t *testing.T) {
	r := require.New(t)

	srcFile, err := sampleData("simple_manifest.yml")
	r.NoError(err, "Failed to locate sample data")

	manifest, err := ParseManifest(srcFile)
	r.NoError(err, "Failed to parse manifest file")

	expSchemaVersion, err := version.NewVersion("1.0")
	r.NoError(err)
	expJiggerVersion, err := version.NewVersion("0.0")
	r.NoError(err)
	r.Equal(*expSchemaVersion, manifest.Schema.Version)
	r.Equal(*expJiggerVersion, manifest.Schema.JiggerVersion)
	r.Equal(1, len(manifest.TemplateParams))
	r.Contains(manifest.TemplateParams, "project_name")
	r.Equal("MyProj", manifest.TemplateParams["project_name"])
}

func Test_parseManifestNotExist(t *testing.T) {
	r := require.New(t)
	a := assert.New(t)

	// Given an empty temp folder
	tmpDir, err := os.MkdirTemp("", "")
	r.NoError(err, "Error creating temp folder")

	// Make sure we always clean up our temp folder
	defer func() {
		r.NoError(os.RemoveAll(tmpDir), "Error deleting temp folder")
	}()

	manifest, err := ParseManifest(path.Join(tmpDir, "fubar.yml"))
	a.Error(err)
	a.Nil(manifest)
}

func Test_parseManifestInvalidYAML(t *testing.T) {
	r := require.New(t)
	a := assert.New(t)

	// Given an empty temp folder
	tmpDir, err := os.MkdirTemp("", "")
	r.NoError(err, "Error creating temp folder")

	// Make sure we always clean up our temp folder
	defer func() {
		r.NoError(os.RemoveAll(tmpDir), "Error deleting temp folder")
	}()

	// and a sample config file that contains non-yaml data
	samplefile := path.Join(tmpDir, "fubar.yml")
	srcfile, err := os.Create(samplefile)
	r.NoError(err, "Failed to create test file")
	_, err = srcfile.WriteString("This is not compatible yaml")
	r.NoError(err, "Failed writing test data to disk")
	r.NoError(srcfile.Close())

	// when we try to parse the file
	manifest, err := ParseManifest(samplefile)

	// then we expect error results
	emptyTypeErr := yaml.TypeError{}
	a.Contains(err.Error(), emptyTypeErr.Error())
	a.Nil(manifest)
}
