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
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Parse options file, if it exists, and register the results
			// in our command context
			v := viper.GetViper()
			err := initViper(v)
			if err != nil {
				return err
			}

			appOptions, err := ao.FromViper(v)
			if err != nil {
				return err
			}
			ctx := context.WithValue(cmd.Context(), shared.CkOptions, appOptions)
			cmd.SetContext(ctx)
			return nil
		},
	}
	retval.AddCommand(create.CreateCmd())
	return retval
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(cmd *cobra.Command) error {
	cc.Init(&cc.Config{
		RootCmd:         cmd,
		Headings:        cc.HiCyan + cc.Bold + cc.Underline,
		Commands:        cc.HiYellow + cc.Bold,
		Example:         cc.Italic,
		ExecName:        cc.Bold,
		Flags:           cc.Bold,
		NoExtraNewlines: true,
	})

	return errors.WithStack(cmd.Execute())
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
