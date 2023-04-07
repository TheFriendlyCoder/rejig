package cmd

import (
	"bytes"
	"fmt"
	"github.com/TheFriendlyCoder/rejigger/lib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"path"
	"path/filepath"
	"testing"
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
	options := lib.AppOptions{
		Templates: []lib.TemplateOptions{{
			Alias:  templateName,
			Source: ".",
			Type:   lib.TstLocal,
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
	options := lib.AppOptions{
		Templates: []lib.TemplateOptions{},
	}
	// when we validate our input args
	args := []string{destDir, templateName}
	result := validateArgs(options, args)

	// we expect the proper error to be returned
	r.Error(result)
	r.ErrorAs(result, &lib.UnknownTemplateError{TemplateAlias: templateName})
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
	options := lib.AppOptions{
		Templates: []lib.TemplateOptions{{
			Alias:  templateName,
			Source: ".",
			Type:   lib.TstLocal,
		}},
	}

	// when we validate our input args
	args := []string{destDir, templateName}
	result := validateArgs(options, args)

	// we expect the proper error to be returned
	r.Error(result)
	r.ErrorAs(result, &lib.PathError{Path: destDir, ErrorType: lib.PePathNotEmpty})
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
	srcDir, err := sampleProj("simple")
	r.NoError(err)
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

	// r.Contains(output.String(), *srcDir)
	r.Contains(output.String(), outputDir)

	a.DirExists(filepath.Join(outputDir, "MyProj"))
	a.NoFileExists(filepath.Join(outputDir, ".rejig.yml"))

	//exp := filepath.Join(*srcDir, ".gitignore")
	act := filepath.Join(outputDir, ".gitignore")
	a.FileExists(act)
	//a.True(unmodified(r, exp, act))

	act = filepath.Join(outputDir, "version.txt")
	a.FileExists(act)
	//a.True(contains(r, act, expVersion))
	//a.False(contains(r, act, "{{version}}"))

	act = filepath.Join(outputDir, "MyProj", "main.txt")
	a.FileExists(act)
	//a.True(contains(r, act, expProj))
	//a.False(contains(r, act, "{{project_name}}"))
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
