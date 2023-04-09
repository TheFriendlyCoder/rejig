package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"testing"

	ao "github.com/TheFriendlyCoder/rejigger/lib/applicationOptions"
	e "github.com/TheFriendlyCoder/rejigger/lib/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_generateUsageLine(t *testing.T) {
	a := assert.New(t)
	result := generateUsageLine()

	a.Contains(result, "create")
	a.Contains(result, "targetPath")
	a.Contains(result, "templateAlias")
}

func Test_ValidateArgsSuccess(t *testing.T) {
	r := require.New(t)

	// Given an empty temp folder
	tmpDir, err := os.MkdirTemp("", "")
	r.NoError(err)
	defer os.RemoveAll(tmpDir)

	// with an empty output folder
	destDir := path.Join(tmpDir, "dest")
	r.NoError(os.Mkdir(destDir, 0700))

	// and a mock set of app options
	templateName := "MyTemplate"
	options := ao.AppOptions{
		Templates: []ao.TemplateOptions{{
			Alias:  templateName,
			Source: ".",
			Type:   ao.TstLocal,
		}},
	}

	// validation should succeed
	args := []string{destDir, templateName}
	r.NoError(validateArgs(options, args))
}

func Test_ValidateArgsTemplateNotExists(t *testing.T) {
	r := require.New(t)

	// Given an empty temp folder
	tmpDir, err := os.MkdirTemp("", "")
	r.NoError(err)
	defer os.RemoveAll(tmpDir)

	// with an empty output folder
	destDir := path.Join(tmpDir, "dest")
	r.NoError(os.Mkdir(destDir, 0700))

	// and a mock set of app options
	templateName := "MyTemplate"
	options := ao.AppOptions{
		Templates: []ao.TemplateOptions{},
	}
	// when we validate our input args
	args := []string{destDir, templateName}
	result := validateArgs(options, args)

	// we expect the proper error to be returned
	r.Error(result)
	r.ErrorIs(result, e.NewUnknownTemplateError(templateName))
}

func Test_ValidateArgsTargetDirNotEmpty(t *testing.T) {
	r := require.New(t)

	// Given an empty temp folder
	tmpDir, err := os.MkdirTemp("", "")
	r.NoError(err)
	defer os.RemoveAll(tmpDir)

	// with a non-empty destination folder
	destDir := path.Join(tmpDir, "dest")
	r.NoError(os.Mkdir(destDir, 0700))
	_, err = os.Create(path.Join(destDir, "fubar.txt"))
	r.NoError(err)

	// and a mock set of app options
	templateName := "MyTemplate"
	options := ao.AppOptions{
		Templates: []ao.TemplateOptions{{
			Alias:  templateName,
			Source: ".",
			Type:   ao.TstLocal,
		}},
	}

	// when we validate our input args
	args := []string{destDir, templateName}
	result := validateArgs(options, args)

	// we expect the proper error to be returned
	r.Error(result)
	r.ErrorIs(result, e.NewPathError(destDir, e.PePathNotEmpty))
}

func Test_CreateCommandSucceeds(t *testing.T) {

	r := require.New(t)
	a := assert.New(t)

	// Given an empty temp folder
	tmpDir, err := os.MkdirTemp("", "")
	r.NoError(err)
	defer os.RemoveAll(tmpDir)

	// an empty output folder
	outputDir := path.Join(tmpDir, "output")
	r.NoError(os.Mkdir(outputDir, 0700))

	// and an app options file with a template pointing to our project
	templateName := "MyTemplate"
	srcDir := getProjectDir("simple")
	optionsText := fmt.Sprintf(`
templates:
  - type: local
    source: %s
    alias: %s`, srcDir, templateName)
	configFile := path.Join(tmpDir, "options.yml")
	err = os.WriteFile(configFile, []byte(optionsText), 0600)
	r.NoError(err)

	// and some fake user input to respond to prompts from the template
	output := new(bytes.Buffer)
	fakeInput := new(bytes.Buffer)
	_, err = fakeInput.WriteString("MyProj\n")
	r.NoError(err)
	_, err = fakeInput.WriteString("1.2.3\n")
	r.NoError(err)

	// When we trigger the create command
	rootCmd.SetOut(output)
	rootCmd.SetErr(output)
	rootCmd.SetIn(fakeInput)
	rootCmd.SetArgs([]string{"--config=" + configFile, "create", outputDir, templateName})
	err = rootCmd.Execute()
	r.NoError(err, "CLI command should have succeeded")

	r.Contains(output.String(), templateName)
	r.Contains(output.String(), outputDir)

	a.DirExists(filepath.Join(outputDir, "MyProj"))
	a.NoFileExists(filepath.Join(outputDir, ".rejig.yml"))

	act := filepath.Join(outputDir, ".gitignore")
	a.FileExists(act)

	act = filepath.Join(outputDir, "version.txt")
	a.FileExists(act)

	act = filepath.Join(outputDir, "MyProj", "main.txt")
	a.FileExists(act)
	// TODO: add stack-trace support throughout error handlers
	// TODO: add error helpers to unit tests to make sure all errors have a stack trace
	// TODO: add helper to generate stack trace without duplicate frames
}

func Test_CreateCommandTooFewArgs(t *testing.T) {

	r := require.New(t)

	// Given an empty temp folder
	tmpDir, err := os.MkdirTemp("", "")
	r.NoError(err)
	defer os.RemoveAll(tmpDir)

	// with 1 empty subfolder
	srcDir := path.Join(tmpDir, "src")
	r.NoError(os.Mkdir(srcDir, 0700))

	// When we execute our root command with a missing positional arg
	actual := new(bytes.Buffer)
	rootCmd.SetOut(actual)
	rootCmd.SetErr(actual)
	rootCmd.SetArgs([]string{"create", srcDir})
	err = rootCmd.Execute()

	// We expect an error to be returned
	r.Error(err)

	// and some status information to be reported
	r.Contains(actual.String(), "Error:")
}

func Test_CreateCommandInvalidTemplateName(t *testing.T) {

	r := require.New(t)

	// Given an empty temp folder
	tmpDir, err := os.MkdirTemp("", "")
	r.NoError(err)
	defer os.RemoveAll(tmpDir)

	// with an absent source folder
	templateName := "MyTemplate"
	srcDir := path.Join(tmpDir, "src")
	optionsText := fmt.Sprintf(`
templates:
  - type: local
    source: %s
    alias: %s`, srcDir, templateName)
	configFile := path.Join(tmpDir, "options.yml")
	err = os.WriteFile(configFile, []byte(optionsText), 0600)
	r.NoError(err)

	// and an empty output folder
	destDir := path.Join(tmpDir, "dest")
	r.NoError(os.Mkdir(destDir, 0700))

	// when we attempt to execute our create command
	actual := new(bytes.Buffer)
	rootCmd.SetOut(actual)
	rootCmd.SetErr(actual)
	rootCmd.SetArgs([]string{"--config=" + configFile, "create", destDir, "DoesNotExist"})
	err = rootCmd.Execute()

	// We expect an error to be returned
	r.Error(err)

	// And some status info to be reported
	r.Contains(actual.String(), "DoesNotExist")
}

func Test_CreateCommandGenerateFailure(t *testing.T) {
	r := require.New(t)
	a := assert.New(t)

	// Given a read-only output folder
	tmpDir, err := os.MkdirTemp("", "")
	r.NoError(err)
	defer os.RemoveAll(tmpDir)
	outputDir := path.Join(tmpDir, "output")
	r.NoError(os.Mkdir(outputDir, 0400))

	// and an app options file with a template pointing to our project
	templateName := "MyTemplate"
	srcDir := getProjectDir("simple")
	optionsText := fmt.Sprintf(`
templates:
  - type: local
    source: %s
    alias: %s`, srcDir, templateName)
	configFile := path.Join(tmpDir, "options.yml")
	err = os.WriteFile(configFile, []byte(optionsText), 0600)
	r.NoError(err)

	// and some fake user input to respond to prompts from the template
	output := new(bytes.Buffer)
	fakeInput := new(bytes.Buffer)
	_, err = fakeInput.WriteString("MyProj\n")
	r.NoError(err)
	_, err = fakeInput.WriteString("1.2.3\n")
	r.NoError(err)

	// When we trigger the create command
	rootCmd.SetOut(output)
	rootCmd.SetErr(output)
	rootCmd.SetIn(fakeInput)
	rootCmd.SetArgs([]string{"--config=" + configFile, "create", outputDir, templateName})
	err = rootCmd.Execute()

	// The operation should fail
	r.Error(err)

	// And there should be a stack trace in the output
	a.Contains(output.String(), "permission denied")
	a.Contains(output.String(), "create.go")
}

func Test_FindTemplate(t *testing.T) {
	r := require.New(t)

	expAlias := "Fubar"
	expTempl := ao.TemplateOptions{
		Alias:  expAlias,
		Source: "/tmp",
		Type:   ao.TstLocal,
	}
	appOptions := ao.AppOptions{
		Templates: []ao.TemplateOptions{expTempl},
	}

	result, err := findTemplate(appOptions, expAlias)
	r.NoError(err)
	r.Equal(expTempl, result)
}

func Test_FindTemplateInvalidName(t *testing.T) {
	r := require.New(t)

	expAlias := "Fubar.Was.Here"
	appOptions := ao.AppOptions{
		Templates: []ao.TemplateOptions{},
	}

	_, err := findTemplate(appOptions, expAlias)
	r.ErrorIs(err, e.AOInvalidTemplateNameError())
}

func Test_FindTemplateFromInventory(t *testing.T) {
	r := require.New(t)

	expAlias := "Fubar"
	//expNamespace := "MyNS"

	expTempl := ao.TemplateOptions{
		Alias:  expAlias,
		Source: "/tmp",
		Type:   ao.TstLocal,
	}

	// TODO: Use dependency injection to add a fake inventory here
	//		 for searching
	appOptions := ao.AppOptions{
		Templates:   []ao.TemplateOptions{},
		Inventories: []ao.InventoryOptions{},
	}

	result, err := findTemplate(appOptions, expAlias)
	r.NoError(err)
	r.Equal(expTempl, result)
}
