package cmd

import (
	"fmt"
	"github.com/TheFriendlyCoder/rejigger/lib"
	"github.com/pkg/errors"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/spf13/cobra"
)

// rootArgs parsed command line arguments
type rootArgs struct {
	// sourcePath path to the folder containing the template to be loaded
	sourcePath string
	// targetPath path to the folder where the new project is to be created
	targetPath string
}

// run Primary entry point function for our generator
func run(cmd *cobra.Command, args *rootArgs) error {
	// We have to use cmd.OutOrStdout() to ensure output is redirected to Cobra
	// stream handler, to facilitate testing (ie: it allows us to capture output
	// during unit testing to validate results of CLI operations)
	lib.SNF(fmt.Fprintf(cmd.OutOrStdout(), "Loading template from %s...\n", args.sourcePath))

	manifest := filepath.Join(args.sourcePath, ".rejig.yml")
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
	lib.SNF(fmt.Fprintf(cmd.OutOrStdout(), "Generating project %s from %s...\n", args.targetPath, args.sourcePath))

	return errors.Wrap(lib.Generate(args.sourcePath, args.targetPath, context), "Failed generating project")
	// TODO: after generating, put a manifest file in the root folder summarizing what we did so we
	//		 can regenerate or update the project later
	// TODO: make terminology consistent (ie: config file for the app, manifest file for the template,
	//		 and something else for storing status of generated project - maybe audit file?
}

// validateArgs checks to see if the command line args provided to the app are valid
func validateArgs(args []string) error {
	if !lib.DirExists(args[0]) {
		return lib.PathError{
			Path:      args[0],
			ErrorType: lib.PE_PATH_NOT_FOUND,
		}
	}
	if lib.DirExists(args[1]) {
		contents, err := os.ReadDir(args[1])
		if err != nil {
			log.Panic(err)
		}
		if len(contents) != 0 {
			return lib.PathError{
				Path:      args[1],
				ErrorType: lib.PE_PATH_NOT_EMPTY,
			}
		}
	}
	return nil
}

// createCmd defines the structure for the "create" subcommand
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "create a new project from a template",
	Long:  `Creates a new project in an empty folder using content defined in a template`,
	Args: func(cmd *cobra.Command, args []string) error {
		// Optionally run one of the validators provided by cobra
		if err := cobra.ExactArgs(2)(cmd, args); err != nil {
			return err
		}
		if err := validateArgs(args); err != nil {
			return err
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		parsedArgs := rootArgs{
			sourcePath: args[0],
			targetPath: args[1],
		}
		err := run(cmd, &parsedArgs)
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
