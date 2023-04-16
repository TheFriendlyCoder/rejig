package cmd

import (
	"bytes"
	"os"
	"path"
	"testing"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_helpCommand(t *testing.T) {
	r := require.New(t)

	tmpDir, err := os.MkdirTemp("", "")
	r.NoError(err)
	defer os.RemoveAll(tmpDir)

	// Point our command to an empty config file so the test suite
	// doesn't try and read config options from the users home folder
	configFile := path.Join(tmpDir, "config.yml")
	fh, err := os.Create(configFile)
	r.NoError(err)
	_, err = fh.WriteString("")
	r.NoError(err)
	r.NoError(fh.Close())

	// When we run the root command with a custom config file
	actual := new(bytes.Buffer)
	rootCmd.SetOut(actual)
	rootCmd.SetErr(actual)

	rootCmd.SetArgs([]string{"--config", configFile, "help"})
	err = rootCmd.Execute()

	// We expect no error from our command
	r.NoError(err, "CLI command should have succeeded")
	r.Contains(actual.String(), "rejigger")
}

func Test_executeConfigFileInvalidYAML(t *testing.T) {
	r := require.New(t)

	tmpDir, err := os.MkdirTemp("", "")
	r.NoError(err)
	defer os.RemoveAll(tmpDir)

	configFile := path.Join(tmpDir, "sample.yml")
	fh, err := os.Create(configFile)
	r.NoError(err)
	_, err = fh.WriteString("templates: fubar")
	r.NoError(err)
	r.NoError(fh.Close())

	// When we run the root command with a custom config file
	actual := new(bytes.Buffer)
	rootCmd.SetOut(actual)
	rootCmd.SetErr(actual)

	rootCmd.SetArgs([]string{"--config", configFile, "help"})
	err = rootCmd.Execute()
	r.Error(err)
}

func Test_initConfigDefaults(t *testing.T) {
	r := require.New(t)

	v := viper.New()

	_, err := initConfig(v, "")
	_, err2 := os.Stat(v.ConfigFileUsed())
	// On CI the default config file won't be found
	if os.IsNotExist(err2) {
		// and in this case an error is triggered
		r.Error(err)
	} else {
		r.NoError(err)
	}

	// NOTE: we can't validate the results of the parse operation here
	// because the results of the default parsing will depend on the
	// users config file in the current home folder
}

func Test_initConfigFileNotExist(t *testing.T) {
	r := require.New(t)
	a := assert.New(t)

	v := viper.New()

	tmpDir, err := os.MkdirTemp("", "")
	r.NoError(err)
	defer os.RemoveAll(tmpDir)

	opts, err := initConfig(v, path.Join(tmpDir, "fubar.yml"))
	r.NoError(err)
	a.Equal(0, len(opts.Templates))
	a.Equal(0, len(opts.Inventories))
}

func Test_initConfigFileNoPermission(t *testing.T) {
	r := require.New(t)

	v := viper.New()

	tmpDir, err := os.MkdirTemp("", "")
	r.NoError(err)
	defer os.RemoveAll(tmpDir)

	outputDir := path.Join(tmpDir, "test")
	r.NoError(os.Mkdir(outputDir, 0600))

	_, err = initConfig(v, path.Join(outputDir, "fubar.yml"))
	r.Error(err)
}

func Test_initConfigFileInvalidYAML(t *testing.T) {
	r := require.New(t)

	v := viper.New()

	tmpDir, err := os.MkdirTemp("", "")
	r.NoError(err)
	defer os.RemoveAll(tmpDir)

	configFile := path.Join(tmpDir, "sample.yml")
	fh, err := os.Create(configFile)
	r.NoError(err)
	_, err = fh.WriteString("templates: fubar")
	r.NoError(err)
	r.NoError(fh.Close())

	_, err = initConfig(v, configFile)
	r.Error(err)
}

func Test_checkErr(t *testing.T) {
	r := require.New(t)

	r.Panics(func() { checkErr(errors.New("Error")) })
	r.NotPanics(func() { checkErr(nil) })
}
