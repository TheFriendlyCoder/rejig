package templateManager

import (
	"bufio"
	"fmt"
	"github.com/TheFriendlyCoder/rejigger/lib"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"os"
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
	manifestData lib.ManifestData
	// templateContext this is a mapping of user configurable options supported
	// by this template, with the values that are provided by the user. These
	// parameters customize the behavior of the template and are used when
	// generating a new instance of the template
	templateContext map[string]any
}

// New constructs new instances of our template manager, which allows the caller
// to interact with a template in various ways
func New(options lib.TemplateOptions) (templateManager, error) {
	// Initialize empty options
	retval := templateManager{}
	retval.Options = options
	retval.templateContext = map[string]any{}

	// Parse manifest file
	manifestPath := filepath.Join(options.Source, manifestFileName)
	_, err := os.Stat(manifestPath)
	if err != nil {
		return retval, errors.Wrap(err, "Unable to read manifest file")
	}

	retval.manifestData, err = lib.ParseManifest(manifestPath)
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
	fs := afero.NewOsFs()
	return errors.Wrap(lib.Generate(fs, t.Options.Source, targetPath, t.templateContext), "Failed generating project")
}
