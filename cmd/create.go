package cmd

import (
	"fmt"
	"github.com/TheFriendlyCoder/rejigger/lib"
	"github.com/TheFriendlyCoder/rejigger/lib/template"
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

func getTemplate(appOptions lib.AppOptions, alias string) (template.Template, error) {
	// Lookup our template information
	var retval template.Template
	var curTemplate lib.TemplateOptions
	found := false
	for _, t := range appOptions.Templates {
		if t.Alias == alias {
			curTemplate = t
			found = true
			break
		}
	}
	if !found {
		return retval, lib.UnknownTemplateError{TemplateAlias: alias}
	}
	retval.Options = curTemplate
	return retval, nil
}

// run Primary entry point function for our generator
func run(cmd *cobra.Command, args rootArgs) error {
	// We have to use cmd.OutOrStdout() to ensure output is redirected to Cobra
	// stream handler, to facilitate testing (ie: it allows us to capture output
	// during unit testing to validate results of CLI operations)
	lib.SNF(fmt.Fprintf(cmd.OutOrStdout(), "Loading template %s...\n", args.templateAlias))

	appOptions, ok := cmd.Context().Value(CkOptions).(lib.AppOptions)
	if !ok {
		return lib.CommandContextNotDefined
	}
	curTemplate, err := getTemplate(appOptions, args.templateAlias)
	if err != nil {
		// TODO: Audit code and see if WithStack may be preferable to Wrap in many places
		return errors.WithStack(err)
	}

	err = errors.WithStack(curTemplate.Validate())
	if err != nil {
		return err
	}

	err = errors.WithStack(curTemplate.LoadManifest(cmd))
	if err != nil {
		return err
	}
	lib.SNF(fmt.Fprintf(cmd.OutOrStdout(), "Generating project %s from template %s...\n", args.targetPath, curTemplate.Options.Alias))

	return errors.Wrap(curTemplate.Generate(args.targetPath), "Failed generating template")
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
	Run: func(cmd *cobra.Command, args []string) {
		parsedArgs := rootArgs{
			targetPath:    args[0],
			templateAlias: args[1],
		}
		err := run(cmd, parsedArgs)
		if err != nil {
			// TODO: In verbose mode output stack trace
			// TODO: Test stack track without using wrap to see how it looks
			// https://pkg.go.dev/github.com/pkg/errors#hdr-Retrieving_the_stack_trace_of_an_error_or_wrapper
			type stackTracer interface {
				StackTrace() errors.StackTrace
			}
			err2, _ := err.(stackTracer)
			for _, f := range err2.StackTrace() {
				fmt.Printf("%+s:%d\n", f, f)
			}
			lib.SNF(fmt.Fprintf(cmd.ErrOrStderr(), "Failed to generate project: \n\t%s\n", errors.Cause(err).Error()))
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
