package cmd

import (
	"context"
	"github.com/TheFriendlyCoder/rejigger/lib"
	cc "github.com/ivanpirog/coloredcobra"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"reflect"
)

// cfgFile path to the file containing config options for the app
var cfgFile string

// ContextKey enumerated type defining keys in the Cobra context manager used to store
// and retrieve common command properties
type ContextKey int64

const (
	// CkOptions Parsed application options loaded from the environment or app config file
	// should be managed exclusively by the root command
	CkOptions ContextKey = iota
	// CkArgs Command line args, parsed into an internal struct format
	// Type of this context object is unique for each command
	CkArgs
)

// checkErr replacement for the cobra method of the same name, which unfortunately calls os.exit
// under the hood, making it impossible to write unit tests for it. This helper calls out to panic()
// which allows us to intercept the termination signal during testing
func checkErr(err error) {
	if err != nil {
		panic("Critical application failure")
	}
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Short: "Project templating tool",
	Long: `The rejigger app allows you to generate source code projects from
specially formatted files stored on disk or in Git repositories`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// Parse options file, if it exists, and register the results
		// in our command context
		if err := cobra.ExactArgs(2)(cmd, args); err != nil {
			return err
		}
		appOptions, err := initConfig()
		if err != nil {
			return err
		}
		ctx := context.WithValue(cmd.Context(), CkOptions, appOptions)
		cmd.SetContext(ctx)
		return nil
	},
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
	rootCmd.Use = "rejig"

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// init function called by GO when the module is first loaded, to initialize the application state
func init() {
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
			return nil, lib.AppOptionsDecodeError
		}

		var newVal lib.TemplateSourceType
		switch templateData["type"] {
		case "":
			newVal = lib.TstUndefined
		case "git":
			newVal = lib.TstGit
		case "local":
			newVal = lib.TstLocal
		default:
			return nil, lib.AppOptionsInvalidSourceTypeError
		}
		templateData["type"] = newVal
		return templateData, nil
	}
}

// initConfig reads app config info from a config file if provided
func initConfig() (lib.AppOptions, error) {
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

	var appOptions lib.AppOptions

	// If a config file is found, read it in.
	err := viper.ReadInConfig()
	// If there is no config file, we ignore that error and assume
	// there is no app config
	if errors.As(err, &viper.ConfigFileNotFoundError{}) {
		return appOptions, nil
	} else if err != nil {
		return appOptions, errors.Wrap(err, "Failure reading options file")
	}

	// If we have a config file, try parsing it
	if err = viper.Unmarshal(&appOptions, viper.DecodeHook(appOptionsDecoder())); err != nil {
		return appOptions, errors.Wrap(err, "Failed parsing app options file")
	}

	// Then validate the results to make sure they meet the application requirements
	if err = appOptions.Validate(); err != nil {
		return appOptions, errors.Wrap(err, "App options file failed validation")
	}
	return appOptions, nil
}
