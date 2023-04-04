package cmd

import (
	"fmt"
	"github.com/TheFriendlyCoder/rejigger/lib"
	cc "github.com/ivanpirog/coloredcobra"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

// cfgFile path to the file containing config options for the app
var cfgFile string

// appOptions global config options for the app
var appOptions *lib.AppOptions

// rootArgs parsed command line arguments
type rootArgs struct {
	// sourcePath path to the folder containing the template to be loaded
	sourcePath string
	// targetPath path to the folder where the new project is to be created
	targetPath string
}

// checkErr replacement for the cobra method of the same name, which unfortunately calls os.exit
// under the hood, making it impossible to write unit tests for it. This helper calls out to panic()
// which allows us to intercept the termination signal during testing
func checkErr(err error) {
	if err != nil {
		panic("Critical application failure")
	}
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

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Short: "Project templating tool",
	Long: `The rejigger app allows you to generate source code projects from
specially formatted files stored on disk or in Git repositories`,
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
	retval = "rejig "
	temp := reflect.ValueOf(&rootArgs{}).Elem()

	for i := 0; i < temp.NumField(); i++ {
		varName := temp.Type().Field(i).Name
		retval += varName + " "
	}
	return strings.TrimSpace(retval)
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cc.Init(&cc.Config{
		RootCmd:         rootCmd,
		Headings:        cc.HiCyan + cc.Bold + cc.Underline,
		Commands:        cc.HiYellow + cc.Bold,
		Example:         cc.Italic,
		ExecName:        cc.Bold,
		Flags:           cc.Bold,
		NoExtraNewlines: true,
	})
	rootCmd.Use = generateUsageLine()
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// init function called by GO when the module is first loaded, to initialize the application state
func init() {
	cobra.OnInitialize(initConfig)

	// Global application flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.rejig.yaml)")

}

// appOptionsDecoder custom hook method used to translate raw config data into a structure
// that is easier to leverage in the application code
func appOptionsDecoder() mapstructure.DecodeHookFuncType {
	// Based on example found here:
	//		https://sagikazarmark.hu/blog/decoding-custom-formats-with-viper/
	return func(
		src reflect.Type,
		target reflect.Type,
		raw interface{},
	) (interface{}, error) {

		// For now we only need to customize the "type" field of the TemplateOptions
		if (target != reflect.TypeOf(lib.TemplateOptions{})) {
			return raw, nil
		}

		// Map the "type" field from a character string format to an enumerated type
		templateData, ok := raw.(map[string]interface{})
		if !ok {
			return nil, lib.APP_OPTIONS_DECODE_ERROR
		}

		var newVal lib.TemplateSourceType
		switch templateData["type"] {
		case "":
			newVal = lib.TST_UNDEFINED
		case "git":
			newVal = lib.TST_GIT
		default:
			return nil, lib.APP_OPTIONS_INVALID_SOURCE_TYPE_ERROR
		}
		templateData["type"] = newVal
		return templateData, nil
	}
}

// initConfig reads app config info from a config file if provided
func initConfig() {
	if cfgFile != "" {
		// Use config file from the command line flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		checkErr(err)

		// Search config in home directory with name ".rejig" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".rejig")
	}

	// If a config file is found, read it in.
	err := viper.ReadInConfig()
	// If there is no config file, we ignore that error and assume
	// there is no app config
	if errors.As(err, &viper.ConfigFileNotFoundError{}) {
		return
	}
	checkErr(err)

	// Parse the config data
	err = viper.Unmarshal(&appOptions, viper.DecodeHook(appOptionsDecoder()))
	checkErr(err)

	// Then validate the results to make sure they meet the application requirements
	err = appOptions.Validate()
	checkErr(err)
}
