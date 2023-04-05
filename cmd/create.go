package cmd

import (
	"fmt"
	"github.com/TheFriendlyCoder/rejigger/lib"
	"github.com/pkg/errors"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/spf13/cobra"
)

// rootArgs parsed command line arguments
type rootArgs struct {
	// targetPath path to the folder where the new project is to be created
	targetPath string
	// templateInfo name of the template from the app inventory to use for generating the new project
	templateInfo lib.TemplateOptions
}

// run Primary entry point function for our generator
func run(cmd *cobra.Command, args *rootArgs) error {
	// We have to use cmd.OutOrStdout() to ensure output is redirected to Cobra
	// stream handler, to facilitate testing (ie: it allows us to capture output
	// during unit testing to validate results of CLI operations)
	lib.SNF(fmt.Fprintf(cmd.OutOrStdout(), "Loading template %s...\n", args.templateInfo.Alias))

	// TODO: rework gettemplate to download the source repo to disk
	// reason: go-git in memory construct uses a custom memory file system
	//			api called go-billy which isn't 100% compatible with the os.File
	//			interface, making it hard to write processing code that works with
	//			both in-memory constructs AND file based ones
	templateSrc, err := lib.GetTemplate(args.templateInfo.Source)
	if err != nil {
		return errors.Wrap(err, "Failed to load Git template")
	}
	fs := *templateSrc
	_, err = fs.Stat(".rejig.yml")
	if err != nil {
		// TODO: See if there's any need to check for the IsNotExist error
		return errors.Wrap(err, "Unable to read manifest file")
	}

	manifest, err := fs.Open(".rejig.yml")
	if err != nil {
		return errors.Wrap(err, "Failed reading manifest file from template")
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
	lib.SNF(fmt.Fprintf(cmd.OutOrStdout(), "Generating project %s from template %s...\n", args.targetPath, args.templateInfo.Alias))

	return errors.Wrap(lib.Generate(fs, args.targetPath, args.templateInfo, context), "Failed generating project")
	// TODO: after generating, put a manifest file in the root folder summarizing what we did so we
	//		 can regenerate or update the project later
	// TODO: make terminology consistent (ie: config file for the app, manifest file for the template,
	//		 and something else for storing status of generated project - maybe audit file?
}

// validateArgs checks to see if the command line args provided to the app are valid
func validateArgs(options lib.AppOptions, args []string) error {
	// The target folder must be empty, if it exists
	if lib.DirExists(args[0]) {
		contents, err := os.ReadDir(args[0])
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

	// see if a template with the given name exists
	// TODO: Hook up this validator to Cobra using an enum parameter that
	//		can maybe work with command line expansion as well
	found := false
	for _, curTemplate := range options.Templates {
		if curTemplate.Alias == args[1] {
			found = true
			break
		}
	}
	if !found {
		return lib.TemplateNotFoundError{
			TemplateAlias: args[1],
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
		if err := validateArgs(appOptions, args); err != nil {
			return err
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		parsedArgs := rootArgs{
			targetPath:   args[0],
			templateInfo: appOptions.GetTemplate(args[1]),
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
