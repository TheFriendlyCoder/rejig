package template

import (
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

	return nil
}

const manifestFileName = ".rejig.yml"

func (t *Template) LoadManifest(cmd *cobra.Command) error {
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
		lib.SNF(fmt.Fscanln(cmd.InOrStdin(), &temp))
		t.templateContext[arg.Name] = temp
	}
	return nil
}

func (t *Template) getRootFolder() string {
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

	// TODO: See if there are times when we don't need a wrap
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
