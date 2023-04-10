package cmd

import (
	"fmt"
	"github.com/TheFriendlyCoder/rejigger/lib"
	ao "github.com/TheFriendlyCoder/rejigger/lib/applicationOptions"
	e "github.com/TheFriendlyCoder/rejigger/lib/errors"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// listCmd defines the structure for the "create" subcommand
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list all available templates",
	Long:  `list all available templates`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// Initialize our command context from the root command
		err := cmd.Parent().PreRunE(cmd, args)
		if err != nil {
			return errors.WithStack(err)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: See if we can move the error handling / stack trace stuff into
		// a PostRun method
		err := runList(cmd)
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

func runList(cmd *cobra.Command) error {
	appOptions, ok := cmd.Context().Value(CkOptions).(ao.AppOptions)
	if !ok {
		return e.NewInternalError("Failed to retrieve app options")
	}
	for _, t := range appOptions.Templates {
		fmt.Printf("local\t%s\t%s\n", t.Alias, t.Type.ToString())
	}
	for _, i := range appOptions.Inventories {
		inv, err := i.GetTemplateDefinitions()
		if err != nil {
			return err
		}
		for _, t := range inv {
			fmt.Printf("%s\t%s\t%s\n", i.Namespace, t.Alias, t.Type.ToString())
		}
	}
	return nil
}

func init() {
	listCmd.Use = "list"
	rootCmd.AddCommand(listCmd)
}
