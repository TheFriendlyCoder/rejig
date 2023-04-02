package cmd

import (
	"fmt"
	"github.com/TheFriendlyCoder/rejigger/lib"
	cc "github.com/ivanpirog/coloredcobra"
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
		return pathError{
			path:      args[0],
			errorType: pathNotFound,
		}
	}
	if lib.DirExists(args[1]) {
		contents, err := os.ReadDir(args[1])
		if err != nil {
			log.Panic(err)
		}
		if len(contents) != 0 {
			return pathError{
				path:      args[1],
				errorType: pathNotEmpty,
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

// initConfig reads app config info from a config file if provided
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".rejig" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".rejig")
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		lib.SNF(fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed()))
		if err != nil {
			panic("Failed to write error output: " + err.Error())
		}
	}
}
