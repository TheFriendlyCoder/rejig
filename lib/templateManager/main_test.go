package templateManager

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"testing"

	ao "github.com/TheFriendlyCoder/rejigger/lib/applicationOptions"
	"github.com/hashicorp/go-version"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_templateManagerConstructor(t *testing.T) {
	r := require.New(t)
	a := assert.New(t)

	// Given an empty temp folder
	tmpDir, err := os.MkdirTemp("", "")
	r.NoError(err)
	defer os.RemoveAll(tmpDir)

	// and a valid template config file
	expectedParam1 := "project_name"
	expectedDesc1 := "Name of the source code project"
	expectedParam2 := "version"
	expectedDesc2 := "Initial version number for the project"
	expectedTemplateVersion := "2.0"
	expectedAppVersion := "0.0.1"
	expectedSchemaVersion := "1.0"
	templateConfigText := fmt.Sprintf(`
versions:
  schema: %s
  rejigger: %s
  template: %s
template:
  args:
    - name: %s
      description: %s
    - name: %s
      description: %s`, expectedSchemaVersion, expectedAppVersion, expectedTemplateVersion, expectedParam1, expectedDesc1, expectedParam2, expectedDesc2)
	configFile := path.Join(tmpDir, manifestFileName)
	err = os.WriteFile(configFile, []byte(templateConfigText), 0600)
	r.NoError(err)

	options := ao.TemplateOptions2{
		Source: tmpDir,
		Alias:  "MyAlias",
		Type:   ao.TstLocal,
	}

	// When we try constructing a new instance of our manager
	template, err := New(&options)

	// We expect the new object to be created and initialized properly
	r.NoError(err)
	a.Equal(&options, template.Options)
	a.Equal(2, len(template.manifestData.Template.Args))
	a.Equal(template.manifestData.Template.Args[0].Name, expectedParam1)
	a.Equal(template.manifestData.Template.Args[0].Description, expectedDesc1)
	a.Equal(template.manifestData.Template.Args[1].Name, expectedParam2)
	a.Equal(template.manifestData.Template.Args[1].Description, expectedDesc2)
	v, err := version.NewVersion(expectedTemplateVersion)
	r.NoError(err)
	a.Equal(template.manifestData.Versions.Template, *v)
	v, err = version.NewVersion(expectedSchemaVersion)
	r.NoError(err)
	a.Equal(template.manifestData.Versions.Schema, *v)
	v, err = version.NewVersion(expectedAppVersion)
	r.NoError(err)
	a.Equal(template.manifestData.Versions.Jigger, *v)
}

func Test_templateManagerConstructor_NoManifest(t *testing.T) {
	r := require.New(t)
	a := assert.New(t)

	// Given an empty temp folder
	tmpDir, err := os.MkdirTemp("", "")
	r.NoError(err)
	defer os.RemoveAll(tmpDir)

	options := ao.TemplateOptions2{
		Source: tmpDir,
		Alias:  "MyAlias",
		Type:   ao.TstLocal,
	}

	// When we try constructing a new instance of our manager
	_, err = New(&options)
	r.Error(err)
	a.True(os.IsNotExist(errors.Cause(err)))
	// Make sure the error message includes some type of reference to our manifest file
	a.Contains(err.Error(), ".rejig.yml")
}

func Test_templateManagerConstructor_InvalidManifest(t *testing.T) {
	r := require.New(t)
	a := assert.New(t)

	// Given an empty temp folder
	tmpDir, err := os.MkdirTemp("", "")
	r.NoError(err)
	defer os.RemoveAll(tmpDir)

	// and an invalid template config file
	templateConfigText := "fubar"
	configFile := path.Join(tmpDir, manifestFileName)
	err = os.WriteFile(configFile, []byte(templateConfigText), 0600)
	r.NoError(err)

	options := ao.TemplateOptions2{
		Source: tmpDir,
		Alias:  "MyAlias",
		Type:   ao.TstLocal,
	}

	// When we try constructing a new instance of our manager
	_, err = New(&options)

	// Operation should succeed
	r.Error(err)

	// Make sure the error mentions that we were parsing a template manifest
	a.Contains(err.Error(), "Error parsing template manifest")
}

func Test_templateManagerGatherParams(t *testing.T) {
	tests := map[string]struct {
		projName string
		projVer  string
	}{
		"Project name with no spaces": {
			projName: "MyProj",
			projVer:  "1.2.3",
		},
		"Project name with spaces": {
			projName: "My Proj",
			projVer:  "1.2.3",
		},
	}

	for name, data := range tests {
		t.Run(name, func(t *testing.T) {

			r := require.New(t)
			a := assert.New(t)

			// Given an empty temp folder
			tmpDir, err := os.MkdirTemp("", "")
			r.NoError(err)
			defer os.RemoveAll(tmpDir)

			// and a valid template config file
			expectedParam1 := "project_name"
			expectedDesc1 := "Name of the source code project"
			expectedParam2 := "version"
			expectedDesc2 := "Initial version number for the project"
			templateConfigText := fmt.Sprintf(`
versions:
  schema: 1.0
  rejigger: 0.0.1
  template: 1.0
template:
  args:
    - name: %s
      description: %s
    - name: %s
      description: %s`, expectedParam1, expectedDesc1, expectedParam2, expectedDesc2)
			configFile := path.Join(tmpDir, manifestFileName)
			err = os.WriteFile(configFile, []byte(templateConfigText), 0600)
			r.NoError(err)

			options := ao.TemplateOptions2{
				Source: tmpDir,
				Alias:  "MyAlias",
				Type:   ao.TstLocal,
			}

			// and a fake command with some user input to respond to prompts from the template
			output := new(bytes.Buffer)
			fakeInput := new(bytes.Buffer)
			expectedProjName := data.projName
			expectedVersion := data.projVer
			_, err = fakeInput.WriteString(expectedProjName + "\n")
			r.NoError(err)
			_, err = fakeInput.WriteString(expectedVersion + "\n")
			r.NoError(err)

			cmd := cobra.Command{}
			cmd.SetOut(output)
			cmd.SetErr(output)
			cmd.SetIn(fakeInput)

			// when we process the user input
			tm, err := New(&options)
			r.NoError(err)
			err = tm.GatherParams(&cmd)
			r.NoError(err)

			// We expect the parsed parameters to be populated in the manager
			a.Equal(expectedProjName, tm.templateContext["project_name"])
			a.Equal(expectedVersion, tm.templateContext["version"])

			// And we expect to have been prompted for the various input params
			a.Contains(output.String(), expectedDesc1)
			a.Contains(output.String(), expectedDesc2)
			a.Contains(output.String(), expectedParam1)
			a.Contains(output.String(), expectedParam2)
		})
	}
}

func Test_templateManagerGenerate(t *testing.T) {
	r := require.New(t)

	tests := map[string]struct {
		sourceDir    string
		templateType ao.TemplateSourceType
	}{
		"Local file system template": {
			sourceDir:    getProjectDir("simple"),
			templateType: ao.TstLocal,
		},
		"Git file system template": {
			sourceDir:    "https://github.com/TheFriendlyCoder/rejiggerTestTemplate.git",
			templateType: ao.TstGit,
		},
	}

	for name, data := range tests {
		t.Run(name, func(t *testing.T) {

			// Given an empty temp folder
			tmpDir, err := os.MkdirTemp("", "")
			r.NoError(err)
			defer os.RemoveAll(tmpDir)

			options := ao.TemplateOptions2{
				Source: data.sourceDir,
				Alias:  "MyAlias",
				Type:   data.templateType,
			}

			// and a fake command with some user input to respond to prompts from the template
			output := new(bytes.Buffer)
			fakeInput := new(bytes.Buffer)
			expectedProjName := "MyProj"
			expectedVersion := "1.2.3"
			_, err = fakeInput.WriteString(expectedProjName + "\n")
			r.NoError(err)
			_, err = fakeInput.WriteString(expectedVersion + "\n")
			r.NoError(err)

			cmd := cobra.Command{}
			cmd.SetOut(output)
			cmd.SetErr(output)
			cmd.SetIn(fakeInput)

			// when we process the user input
			tm, err := New(&options)
			r.NoError(err)
			err = tm.GatherParams(&cmd)
			r.NoError(err)
			err = tm.Generate(tmpDir)
			r.NoError(err)
		})
	}
}

func Test_templateManagerFailToGenerate(t *testing.T) {
	r := require.New(t)

	// Given a read-only output folder
	tmpDir, err := os.MkdirTemp("", "")
	r.NoError(err)
	defer os.RemoveAll(tmpDir)
	outputDir := path.Join(tmpDir, "output")
	r.NoError(os.Mkdir(outputDir, 0400))

	options := ao.TemplateOptions2{
		Source: getProjectDir("simple"),
		Alias:  "MyAlias",
		Type:   ao.TstLocal,
	}

	// and a fake command with some user input to respond to prompts from the template
	output := new(bytes.Buffer)
	fakeInput := new(bytes.Buffer)
	expectedProjName := "MyProj"
	expectedVersion := "1.2.3"
	_, err = fakeInput.WriteString(expectedProjName + "\n")
	r.NoError(err)
	_, err = fakeInput.WriteString(expectedVersion + "\n")
	r.NoError(err)

	cmd := cobra.Command{}
	cmd.SetOut(output)
	cmd.SetErr(output)
	cmd.SetIn(fakeInput)

	tm, err := New(&options)
	r.NoError(err)
	err = tm.GatherParams(&cmd)
	r.NoError(err)

	// When we try generating in a path that doesn't exist
	err = tm.Generate(outputDir)

	// The operation should fail
	r.Error(err)

}
