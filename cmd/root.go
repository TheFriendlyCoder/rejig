package cmd

import (
	"context"
	"os"

	"github.com/TheFriendlyCoder/rejigger/cmd/create"
	"github.com/TheFriendlyCoder/rejigger/cmd/shared"
	ao "github.com/TheFriendlyCoder/rejigger/lib/applicationOptions"
	cc "github.com/ivanpirog/coloredcobra"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// cfgFile path to the file containing config options for the app
var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "rejig",
	Short: "Project templating tool",
	Long: `The rejigger app allows you to generate source code projects from
specially formatted files stored on disk or in Git repositories`,
	// By default, cmd will always show the app usage message if the command
	// fails and returns an error. This flag disables that behavior.
	SilenceUsage: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Parse options file, if it exists, and register the results
		// in our command context

		appOptions, err := initConfig(viper.GetViper(), cfgFile)
		if err != nil {
			return err
		}
		ctx := context.WithValue(cmd.Context(), shared.CkOptions, appOptions)
		cmd.SetContext(ctx)
		return nil
	},
}

// checkErr replacement for the cobra method of the same name, which unfortunately calls os.exit
// under the hood, making it impossible to write unit tests for it. This helper calls out to panic()
// which allows us to intercept the termination signal during testing
func checkErr(err error) {
	if err != nil {
		panic("Critical application failure")
	}
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

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// initConfig reads app config info from a config file if provided
func initConfig(v *viper.Viper, configFilePath string) (ao.AppOptions, error) {
	if configFilePath != "" {
		// Use config file from the command line flag.
		v.SetConfigFile(configFilePath)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		checkErr(err)

		// Search config in home directory with name ".rejig" (without extension).
		v.AddConfigPath(home)
		v.SetConfigType("yaml")
		v.SetConfigName(".rejig")
	}

	appOptions := ao.New()

	// If a config file is found, read it in.
	err := v.ReadInConfig()
	// If there is no config file, we ignore that error and assume
	// there is no app config
	if os.IsNotExist(err) {
		return appOptions, nil
	} else if err != nil {
		return appOptions, errors.WithStack(err)
	}

	appOptions, err = ao.FromViper(v)
	if err != nil {
		return appOptions, err
	}

	// Then validate the results to make sure they meet the application requirements
	err = appOptions.Validate()
	return appOptions, err
}

// init function called by GO when the module is first loaded, to initialize the application state
func init() {
	// Global application flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.rejig.yaml)")
	rootCmd.AddCommand(create.CreateCmd())
}
