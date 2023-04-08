package templateManager

import (
	"bufio"
	"fmt"
	"github.com/TheFriendlyCoder/rejigger/lib"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/storage/memory"
	"path/filepath"
	"strings"
)

const manifestFileName = ".rejig.yml"

// templateManager stores all of the state related to the template being generated
// and provides methods to interact with that template
type templateManager struct {
	// Options parsed configuration options that describe the template
	// these values are typically provided by the template inventory
	Options lib.TemplateOptions
	// manifestData data parsed from the template manifest file provided
	// by the template itself. This info describes the content and behavior
	// of the template, user configurable options, and other such information
	manifestData ManifestData
	// templateContext this is a mapping of user configurable options supported
	// by this template, with the values that are provided by the user. These
	// parameters customize the behavior of the template and are used when
	// generating a new instance of the template
	templateContext map[string]any
	// srcFilesystem virtual file system to use when interacting with source files
	// defining the content of the template we are managing by this struct
	srcFilesystem afero.Fs
}

// getGitTemplate loads a template source from a Git repository
func getGitTemplate(gitURL string) (afero.Fs, error) {
	appFS := afero.NewMemMapFs()
	fs := NewWraper(appFS, ".", false)

	_, err := git.Clone(memory.NewStorage(), fs, &git.CloneOptions{
		URL: gitURL,
	})
	if err != nil {
		return appFS, errors.Wrap(err, "Failed querying Git repo")
	}
	return appFS, nil
}

// rootDir gets the root folder where the template source files are stored
// takes into account the virtual file system used by the manager
func (t *templateManager) rootDir() string {
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

// New constructs new instances of our template manager, which allows the caller
// to interact with a template in various ways
func New(options lib.TemplateOptions) (templateManager, error) {
	// Initialize empty options
	retval := templateManager{}
	retval.Options = options
	retval.templateContext = map[string]any{}

	var err error
	switch options.Type {
	case lib.TstGit:
		retval.srcFilesystem, err = getGitTemplate(options.Source)
		if err != nil {
			return retval, errors.Wrap(err, "Failed loading Git template")
		}
	case lib.TstLocal:
		retval.srcFilesystem = afero.NewOsFs()
	case lib.TstUndefined:
		fallthrough
	default:
		panic("should never happen: unsupported template type: " + options.Alias)
	}

	// Parse manifest file
	manifestPath := filepath.Join(retval.rootDir(), manifestFileName)
	_, err = retval.srcFilesystem.Stat(manifestPath)
	if err != nil {
		return retval, errors.Wrap(err, "Unable to read manifest file")
	}

	retval.manifestData, err = parseManifest(retval.srcFilesystem, manifestPath)
	if err != nil {
		return retval, errors.Wrap(err, "Failed parsing manifest file")
	}
	return retval, nil
}

// GatherParams iterates over all user defined options supported by this
// template, and prompts the user for values for them all
func (t *templateManager) GatherParams(cmd *cobra.Command) error {
	// TODO: Consider moving this functionality into calling class
	// TODO: return as a no-op if there aren't any args to gather
	reader := bufio.NewReader(cmd.InOrStdin())
	for _, arg := range t.manifestData.Template.Args {
		lib.SNF(fmt.Fprintf(cmd.OutOrStdout(), "%s(%s): ", arg.Description, arg.Name))

		// NOTE: Scanln method apparently doesn't work for reading input strings that have
		// spaces in them, so we use a read buffer here instead
		value, err := reader.ReadString('\n')
		if err != nil {
			return errors.Wrap(err, "Failed reading user input")
		}
		// Here we need to trim white space from our input value to get rid
		// of the trailing newline characters which are included in the read buffer
		t.templateContext[arg.Name] = strings.TrimSpace(value)
	}
	return nil
}

// Generate produces a new template based on the parameters defined in this
// object, in the specified output folder
func (t *templateManager) Generate(targetPath string) error {
	// TODO: add these support methods to the manager class as private methods
	return errors.Wrap(generate(t.srcFilesystem, t.rootDir(), targetPath, t.templateContext), "Failed generating project")
}
