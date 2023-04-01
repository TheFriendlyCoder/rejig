package cmd

import (
	"fmt"
	"github.com/TheFriendlyCoder/rejigger/lib"
	cc "github.com/ivanpirog/coloredcobra"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
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
func run(args *rootArgs) {
	fmt.Printf("Generating from %s to %s...\n", args.sourcePath, args.targetPath)
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
			log.Fatal(err)
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
		run(&parsedArgs)
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

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.rejig.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
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
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
