package cmd

import (
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

// appOptions global config options for the app
var appOptions *lib.AppOptions

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
