package template

import (
	"bufio"
	"fmt"
	"github.com/TheFriendlyCoder/rejigger/lib"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"path/filepath"
)

type Template struct {
	Options         lib.TemplateOptions
	manifestData    lib.ManifestData
	filesys         afero.Fs
	templateContext map[string]any
}

func (t *Template) Validate() error {
	// TODO: Validate file system
	// TODO: Validate manifest file

	// TODO: Maybe make this a New method to run it on creation
	return nil
}

const manifestFileName = ".rejig.yml"

func (t *Template) LoadManifest(cmd *cobra.Command) error {
	// TODO: Break dependency between cmd and this class - we shouldn't be doing direct IO here
	err := errors.Wrap(t.loadFilesystem(), "Error loading template file system")
	if err != nil {
		return err
	}

	manifestFile := filepath.Join(t.getRootFolder(), manifestFileName)
	t.manifestData, err = lib.ParseManifest(t.filesys, manifestFile)
	if err != nil {
		return errors.Wrap(err, "Failed parsing manifest file")
	}
	t.templateContext = map[string]any{}
	for _, arg := range t.manifestData.Template.Args {
		lib.SNF(fmt.Fprintf(cmd.OutOrStdout(), "%s(%s): ", arg.Description, arg.Name))
		var temp string
		// TODO: Find out why scan doesn't like input with spaces in it
		in := bufio.NewReader(cmd.InOrStdin())
		temp, err := in.ReadString('\n')
		if err != nil {
			return errors.Wrap(err, "Failed reading input")
		}
		//lib.SNF(fmt.Fscanln(cmd.InOrStdin(), &temp))
		t.templateContext[arg.Name] = temp
	}
	return nil
}

func (t *Template) getRootFolder() string {
	// TODO: Docstrings
	// TODO: Test coverage
	switch t.Options.Type {
	case lib.TstGit:
		return "."
	case lib.TstLocal:
		return t.Options.Source
	case lib.TstUndefined:
		fallthrough
	default:
		panic("should never happen: unsupported template type: " + t.Options.Alias)
	}
}
func (t *Template) Generate(targetPath string) error {
	rootDir := t.getRootFolder()

	// TODO: see if there are times we don't need to use wrap
	// TODO: move Generate function into this method
	return errors.Wrap(lib.Generate(t.filesys, rootDir, targetPath, t.templateContext), "Failed generating project")
}

func (t *Template) loadFilesystem() error {
	switch t.Options.Type {
	case lib.TstGit:
		var err error
		t.filesys, err = lib.GetGitTemplate(t.Options.Source)
		if err != nil {
			return errors.Wrap(err, "Unable to read template from Git repo")
		}
	case lib.TstLocal:
		t.filesys = afero.NewOsFs()
	case lib.TstUndefined:
		fallthrough
	default:
		panic("should never happen: unsupported template type: " + t.Options.Alias)
	}
	return nil
}
