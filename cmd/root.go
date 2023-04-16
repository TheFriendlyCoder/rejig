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

// CobraThemes mapping table that converts one of our theme types from
// a numeric identifier to a set of pre-defined Cobra color definitions
var CobraThemes = map[ao.ThemeType]cc.Config{
	ao.ThtDark: {
		Headings: cc.HiCyan + cc.Bold + cc.Underline,
		Commands: cc.HiYellow + cc.Bold,
		Example:  cc.Italic,
		ExecName: cc.Bold,
		Flags:    cc.Bold,
	},
	ao.ThtLight: {
		Headings: cc.HiBlue + cc.Bold + cc.Underline,
		Commands: cc.HiMagenta + cc.Bold,
		Example:  cc.Italic,
		ExecName: cc.Bold,
		Flags:    cc.Bold,
	},
	ao.ThtUndefined: {},
	ao.ThtUnknown:   {},
}

// RootCmd definition for the main root / entry point command for the app
func RootCmd() cobra.Command {
	retval := cobra.Command{
		Use:   "rejig",
		Short: "Project templating tool",
		Long: `The rejigger app allows you to generate source code projects from
specially formatted files stored on disk or in Git repositories`,
		// By default, cmd will always show the app usage message if the command
		// fails and returns an error. This flag disables that behavior.
		SilenceUsage: true,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// Sanity checks to make sure application is set up properly
			_, ok := cmd.Context().Value(shared.CkViper).(*viper.Viper)
			if !ok {
				panic("Internal configuration error")
			}
			_, ok = cmd.Context().Value(shared.CkOptions).(ao.AppOptions)
			if !ok {
				panic("Internal configuration error")
			}
		},
	}
	// TODO: if we want to pass path to config file on command line, we may be able to use
	// 		 this helper method to pre-parse the flags from the command line before executing
	// 		 the actual command, allowing us to load app options before execution
	// 			retval.ParseFlags()
	retval.AddCommand(create.CreateCmd())
	return retval
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(cmd *cobra.Command) error {

	// Initialize config file
	v := viper.GetViper()
	err := initViper(v)
	if err != nil {
		return err
	}

	appOptions, err := ao.FromViper(v)
	if err != nil {
		return err
	}

	// Setup application context
	ctx := context.Background()
	ctx = context.WithValue(ctx, shared.CkOptions, appOptions)
	ctx = context.WithValue(ctx, shared.CkViper, v)

	// Setup color theme
	cobraConfig := CobraThemes[appOptions.Other.Theme]
	cobraConfig.RootCmd = cmd
	cobraConfig.NoExtraNewlines = true
	cc.Init(&cobraConfig)

	// Run our command
	return errors.WithStack(cmd.ExecuteContext(ctx))
}

// initViper initializes the Viper app configuration framework
func initViper(v *viper.Viper) error {
	// Find home directory.
	home, err := os.UserHomeDir()
	if err != nil {
		panic("Critical application failure: user home folder not found")
	}

	// Search config in home directory with name ".rejig" (without extension).
	v.AddConfigPath(home)
	v.SetConfigType("yaml")
	v.SetConfigName(".rejig")

	// If a config file is found, read it in.
	err = v.ReadInConfig()

	// If there is no config file, we ignore that error and assume
	// there is no app config
	if !errors.As(err, &viper.ConfigFileNotFoundError{}) && !os.IsNotExist(err) {
		return errors.WithStack(err)
	}
	return nil
}
