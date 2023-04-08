package cmd

import (
	"fmt"
	"github.com/TheFriendlyCoder/rejigger/lib"
	"github.com/TheFriendlyCoder/rejigger/lib/templateManager"
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
	// templateAlias name of the template to use to create the new project from
	templateAlias string
}

// findTemplate looks up a specific template in the template inventory
func findTemplate(appOptions lib.AppOptions, alias string) (lib.TemplateOptions, error) {
	for _, t := range appOptions.Templates {
		if t.Alias == alias {
			return t, nil
		}
	}
	return lib.TemplateOptions{}, lib.UnknownTemplateError{TemplateAlias: alias}
}

// run Primary entry point function for our generator
func run(cmd *cobra.Command, args rootArgs) error {
	// We have to use cmd.OutOrStdout() to ensure output is redirected to Cobra
	// stream handler, to facilitate testing (ie: it allows us to capture output
	// during unit testing to validate results of CLI operations)
	lib.SNF(fmt.Fprintf(cmd.OutOrStdout(), "Loading template %s...\n", args.templateAlias))

	appOptions, ok := cmd.Context().Value(CkOptions).(lib.AppOptions)
	if !ok {
		return lib.InternalError{Message: "Failed to retrieve app options"}
	}

	curTemplate, err := findTemplate(appOptions, args.templateAlias)
	if err != nil {
		return errors.WithStack(err)
	}
	tm, err := templateManager.New(curTemplate)
	if err != nil {
		return errors.Wrap(err, "Error initializing template manager")
	}

	if err = tm.GatherParams(cmd); err != nil {
		return errors.Wrap(err, "Error gathering template parameters")
	}
	lib.SNF(fmt.Fprintf(cmd.OutOrStdout(), "Generating project %s from template %s...\n", args.targetPath, curTemplate.Alias))

	return errors.Wrap(tm.Generate(args.targetPath), "Failed generating project")
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
				ErrorType: lib.PePathNotEmpty,
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
	Args:  cobra.MinimumNArgs(2),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// Initialize our command context from the root command
		err := cmd.Parent().PreRunE(cmd, args)
		if err != nil {
			return errors.Wrap(err, "Initialization error")
		}
		appOptions, ok := cmd.Context().Value(CkOptions).(lib.AppOptions)
		if !ok {
			return lib.CommandContextNotDefined
		}
		return validateArgs(appOptions, args)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		parsedArgs := rootArgs{
			targetPath:    args[0],
			templateAlias: args[1],
		}
		err := run(cmd, parsedArgs)
		if err != nil {
			// https://pkg.go.dev/github.com/pkg/errors#hdr-Retrieving_the_stack_trace_of_an_error_or_wrapper
			type stackTracer interface {
				StackTrace() errors.StackTrace
			}
			if temp, ok := interface{}(err).(stackTracer); ok {
				for _, f := range temp.StackTrace() {
					lib.SNF(fmt.Fprintf(cmd.ErrOrStderr(), "%+s:%d\n", f, f))
				}
			}
			lib.SNF(fmt.Fprintln(cmd.ErrOrStderr(), "Failed to generate project"))
			lib.SNF(fmt.Fprintln(cmd.ErrOrStderr(), err.Error()))
		}
		return err
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
