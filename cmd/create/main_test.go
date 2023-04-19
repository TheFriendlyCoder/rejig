package create

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/TheFriendlyCoder/rejigger/cmd/internal"
	"github.com/TheFriendlyCoder/rejigger/cmd/shared"
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
	a.Contains(result, "templateName")
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
			Name:   templateName,
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
			Name:   templateName,
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

	outputDir := path.Join(tmpDir, "output")

	// and an app options file with a template pointing to our project
	templateName := "MyTemplate"
	srcDir := internal.GetProjectDir()

	appOptions := ao.AppOptions{
		Templates: []ao.TemplateOptions{{
			Type:   ao.TstLocal,
			Source: srcDir,
			Name:   templateName,
		}},
	}

	// and some fake user input to respond to prompts from the template
	output := new(bytes.Buffer)
	fakeInput := new(bytes.Buffer)
	_, err = fakeInput.WriteString("MyProj\n")
	r.NoError(err)
	_, err = fakeInput.WriteString("1.2.3\n")
	r.NoError(err)

	// When we trigger the create command
	createCmd := CreateCmd()
	createCmd.SetOut(output)
	createCmd.SetErr(output)
	createCmd.SetIn(fakeInput)
	ctx := context.TODO()
	ctx = context.WithValue(ctx, shared.CkOptions, appOptions)
	createCmd.SetArgs([]string{outputDir, templateName})
	err = createCmd.ExecuteContext(ctx)
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

	// When we execute our root command with a missing positional arg
	ctx := context.TODO()
	ctx = context.WithValue(ctx, shared.CkOptions, ao.AppOptions{})
	actual := new(bytes.Buffer)
	createCmd := CreateCmd()
	createCmd.SetOut(actual)
	createCmd.SetErr(actual)
	createCmd.SetArgs([]string{tmpDir})
	err = createCmd.ExecuteContext(ctx)

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

	// And a config with no template defined in it
	ctx := context.TODO()
	ctx = context.WithValue(ctx, shared.CkOptions, ao.AppOptions{})

	// when we attempt to create a project using a template name that doesn't exit
	actual := new(bytes.Buffer)
	createCmd := CreateCmd()
	createCmd.SetOut(actual)
	createCmd.SetErr(actual)
	createCmd.SetArgs([]string{tmpDir, "DoesNotExist"})
	err = createCmd.ExecuteContext(ctx)

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
	srcDir := internal.GetProjectDir()
	appOptions := ao.AppOptions{
		Templates: []ao.TemplateOptions{{
			Type:   ao.TstLocal,
			Source: srcDir,
			Name:   templateName,
		}},
	}
	ctx := context.TODO()
	ctx = context.WithValue(ctx, shared.CkOptions, appOptions)

	// and some fake user input to respond to prompts from the template
	output := new(bytes.Buffer)
	fakeInput := new(bytes.Buffer)
	_, err = fakeInput.WriteString("MyProj\n")
	r.NoError(err)
	_, err = fakeInput.WriteString("1.2.3\n")
	r.NoError(err)

	// When we trigger the create command
	createCmd := CreateCmd()
	createCmd.SetOut(output)
	createCmd.SetErr(output)
	createCmd.SetIn(fakeInput)
	createCmd.SetArgs([]string{outputDir, templateName})
	err = createCmd.ExecuteContext(ctx)

	// The operation should fail
	r.Error(err)

	// And there should be a stack trace in the output
	a.Contains(output.String(), "permission denied")
	a.Contains(output.String(), "main.go")
}

func Test_FindTemplate(t *testing.T) {
	r := require.New(t)

	expName := "Fubar"
	expTempl := ao.TemplateOptions{
		Name:   expName,
		Source: "/tmp",
		Type:   ao.TstLocal,
	}
	appOptions := ao.AppOptions{
		Templates: []ao.TemplateOptions{expTempl},
	}

	result, err := findTemplate(appOptions, expName)
	r.NoError(err)
	r.Equal(expTempl, result)
}

func Test_FindTemplateInvalidName(t *testing.T) {
	r := require.New(t)

	expName := "Fubar.Was.Here"
	appOptions := ao.AppOptions{
		Templates: []ao.TemplateOptions{},
	}

	_, err := findTemplate(appOptions, expName)
	r.ErrorIs(err, e.AOInvalidTemplateNameError())
}

func Test_FindTemplateFromInventory(t *testing.T) {
	r := require.New(t)

	// Given a couple of working folders
	tmpDir, err := os.MkdirTemp("", "")
	r.NoError(err)
	defer os.RemoveAll(tmpDir)

	workDirName := path.Join(tmpDir, "mytemplate")
	err = os.Mkdir(workDirName, 0700)
	r.NoError(err)

	// And a sample inventory file
	outputFile := path.Join(tmpDir, ".rejig.inv.yml")
	expName := "test1"
	invData := fmt.Sprintf(`
templates:
  - name: %s
    source: %s
`, expName, workDirName)
	fh, err := os.Create(outputFile)
	r.NoError(err)
	_, err = fh.WriteString(invData)
	r.NoError(err)
	r.NoError(fh.Close())

	expTempl := ao.TemplateOptions{
		Name:   expName,
		Source: tmpDir,
		SubDir: path.Join(tmpDir, "mytemplate"),
		Type:   ao.TstLocal,
	}
	expNamespace := "MyNS"
	tempInvOpts := ao.InventoryOptions{
		Type:      ao.IstLocal,
		Source:    tmpDir,
		Namespace: expNamespace,
	}

	// TODO: Use dependency injection to add a fake inventory here
	//		 for searching, so we don't have to construct all the test data
	appOptions := ao.AppOptions{
		Templates:   []ao.TemplateOptions{},
		Inventories: []ao.InventoryOptions{tempInvOpts},
	}

	result, err := findTemplate(appOptions, expNamespace+"."+expName)
	r.NoError(err)
	r.Equal(expTempl, result)
}
