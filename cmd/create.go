package cmd

import (
	"context"
	"fmt"
	"github.com/TheFriendlyCoder/rejigger/lib"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/spf13/cobra"
)

// rootArgs parsed command line arguments
type rootArgs struct {
	// targetPath path to the folder where the new project is to be created
	targetPath string
	// templateAlias name of the template to use to create the new project from
	templateAlias string
}

// run Primary entry point function for our generator
func run(cmd *cobra.Command, args rootArgs) error {
	// We have to use cmd.OutOrStdout() to ensure output is redirected to Cobra
	// stream handler, to facilitate testing (ie: it allows us to capture output
	// during unit testing to validate results of CLI operations)
	lib.SNF(fmt.Fprintf(cmd.OutOrStdout(), "Loading template %s...\n", args.templateAlias))
	// TODO: Consider not using context objects - the main purpose of them is to share config
	// options between commands, but it seems like the pre-run / run functions for each
	// command are run independently from all other commands so they aren't run in a
	// hierarchy ... so there's not much use here (To Be Confirmed)
	appOptions, ok := cmd.Context().Value(CK_OPTIONS).(lib.AppOptions)
	if !ok {
		return lib.InternalError{Message: "Failed to retrieve app options"}
	}

	// Lookup our template information
	var curTemplate lib.TemplateOptions
	found := false
	for _, t := range appOptions.Templates {
		if t.Alias == args.templateAlias {
			curTemplate = t
			found = true
			break
		}
	}
	if !found {
		return lib.UnknownTemplateError{TemplateAlias: args.templateAlias}
	}

	manifest := filepath.Join(curTemplate.Source, ".rejig.yml")
	_, err := os.Stat(manifest)
	if err != nil {
		// TODO: See if there's any need to check for the IsNotExist error
		return errors.Wrap(err, "Unable to read manifest file")
	}

	manifestData, err := lib.ParseManifest(manifest)
	if err != nil {
		return errors.Wrap(err, "Failed parsing manifest file")
	}

	context := map[string]any{}
	for _, arg := range manifestData.Template.Args {
		lib.SNF(fmt.Fprintf(cmd.OutOrStdout(), "%s(%s): ", arg.Description, arg.Name))
		var temp string
		lib.SNF(fmt.Fscanln(cmd.InOrStdin(), &temp))
		context[arg.Name] = temp
	}
	lib.SNF(fmt.Fprintf(cmd.OutOrStdout(), "Generating project %s from template %s...\n", args.targetPath, curTemplate.Alias))

	fs := afero.NewOsFs()
	return errors.Wrap(lib.Generate(fs, curTemplate.Source, args.targetPath, context), "Failed generating project")
	// TODO: after generating, put a manifest file in the root folder summarizing what we did so we
	//		 can regenerate or update the project later
	// TODO: make terminology consistent (ie: config file for the app, manifest file for the template,
	//		 and something else for storing status of generated project - maybe audit file?
}

// validateArgs checks to see if the command line args provided to the app are valid
func validateArgs(options lib.AppOptions, args []string) error {
	if lib.DirExists(args[0]) {
		contents, err := os.ReadDir(args[0])
		if err != nil {
			log.Panic(err)
		}
		if len(contents) != 0 {
			return lib.PathError{
				Path:      args[0],
				ErrorType: lib.PE_PATH_NOT_EMPTY,
			}
		}
	}

	//Validate template name
	found := false
	for _, t := range options.Templates {
		if t.Alias == args[1] {
			found = true
			break
		}
	}
	if !found {
		return lib.UnknownTemplateError{TemplateAlias: args[1]}
	}
	return nil
}

// createCmd defines the structure for the "create" subcommand
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "create a new project from a template",
	Long:  `Creates a new project in an empty folder using content defined in a template`,
	Args: func(cmd *cobra.Command, args []string) error {
		// TODO: Put config file parsing in base command prerun and call it from
		// prerun in all sub-commands:
		// cmd.Parent().PreRun(cmd, args)
		// NOTE: context object only works within a single command and is used
		// exclusively for passing info between pre-run, run, post-run operations
		// TODO: Is there a benefit of using this technique and not just calling the
		// arg parser directly
		// Optionally run one of the validators provided by cobra
		if err := cobra.ExactArgs(2)(cmd, args); err != nil {
			return err
		}
		appOptions, err := initConfig()
		if err != nil {
			return err
		}
		ctx := context.WithValue(cmd.Context(), CK_OPTIONS, appOptions)
		cmd.SetContext(ctx)

		if err := validateArgs(appOptions, args); err != nil {
			return err
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		parsedArgs := rootArgs{
			targetPath:    args[0],
			templateAlias: args[1],
		}
		err := run(cmd, parsedArgs)
		if err != nil {
			lib.SNF(fmt.Fprintf(cmd.ErrOrStderr(), "Failed to generate project"))
		}
	},
}

// generateUsageLine dynamically generates a usage line for the app based on the contents
// of the rootArgs struct, using reflection
func generateUsageLine() string {
	var retval string
	retval = "create "
	temp := reflect.ValueOf(&rootArgs{}).Elem()

	for i := 0; i < temp.NumField(); i++ {
		varName := temp.Type().Field(i).Name
		retval += varName + " "
	}
	return strings.TrimSpace(retval)
}

func init() {
	createCmd.Use = generateUsageLine()
	rootCmd.AddCommand(createCmd)
}
