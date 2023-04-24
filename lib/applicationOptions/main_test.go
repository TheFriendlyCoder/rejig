package applicationOptions

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_successfulValidation(t *testing.T) {
	a := assert.New(t)

	tests := map[string]struct {
		templateOptions  []TemplateOptions
		inventoryOptions []InventoryOptions
		otherOptions     OtherOptions
	}{
		"Local template file system type": {
			templateOptions: []TemplateOptions{{
				Name:   "My Template",
				Source: "https://some/location",
				Type:   TstLocal,
			}},
		},
		"Git template file system type": {
			templateOptions: []TemplateOptions{{
				Name:   "My Template",
				Source: "https://some/location",
				Type:   TstGit,
			}},
		},
		"Local inventory file system type": {
			inventoryOptions: []InventoryOptions{{
				Namespace: "Fubar",
				Source:    "https://some/location",
				Type:      IstLocal,
			}},
		},
		"Git inventory file system type": {
			inventoryOptions: []InventoryOptions{{
				Namespace: "Fubar",
				Source:    "https://some/location",
				Type:      IstGit,
			}},
		},
		"Dark theme": {
			otherOptions: OtherOptions{Theme: ThtDark},
		},
	}

	for name, data := range tests {
		t.Run(name, func(t *testing.T) {

			opts := AppOptions{
				Templates:   data.templateOptions,
				Inventories: data.inventoryOptions,
			}
			err := opts.Validate()
			a.NoError(err)
		})
	}
}

func Test_successfulValidationEmptyConfig(t *testing.T) {
	r := require.New(t)

	opts := AppOptions{}

	err := opts.Validate()
	r.NoError(err, "Validation should have succeeded")
}

func Test_validationFailures(t *testing.T) {
	r := require.New(t)

	tests := map[string]struct {
		templateOptions  []TemplateOptions
		inventoryOptions []InventoryOptions
	}{
		"Template missing type": {
			templateOptions: []TemplateOptions{{
				Name:   "My Template",
				Source: "https://some/location",
			}},
		},
		"Template missing name": {
			templateOptions: []TemplateOptions{{
				Source: "https://some/location",
				Type:   TstLocal,
			}},
		},
		"Template missing source": {
			templateOptions: []TemplateOptions{{
				Name: "My Template",
				Type: TstGit,
			}},
		},
		"Template duplicate names": {
			templateOptions: []TemplateOptions{
				{
					Name:   "My Template",
					Source: "https://some/location",
					Type:   TstGit,
				},
				{
					Name:   "My Template",
					Source: "/tmp/location",
					Type:   TstLocal,
				}},
		},
		"Inventory missing type": {
			inventoryOptions: []InventoryOptions{{
				Namespace: "Fubar",
				Source:    "https://some/location",
			}},
		},
		"Inventory missing source": {
			inventoryOptions: []InventoryOptions{{
				Namespace: "Fubar",
				Type:      IstLocal,
			}},
		},
		"Inventory missing namespace": {
			inventoryOptions: []InventoryOptions{{
				Source: "https://some/location",
				Type:   IstLocal,
			}},
		},
		"Inventory duplicate namespace": {
			inventoryOptions: []InventoryOptions{
				{
					Namespace: "MyNamespace",
					Source:    "/tmp/location",
					Type:      IstLocal,
				},
				{
					Namespace: "MyNamespace",
					Source:    "https://some/location",
					Type:      IstGit,
				}},
		},
	}

	for name, data := range tests {
		t.Run(name, func(t *testing.T) {
			opts := AppOptions{
				Templates:   data.templateOptions,
				Inventories: data.inventoryOptions,
			}
			err := opts.Validate()
			r.Error(err)
		})
	}
}

func Test_validationTemplateCompoundError(t *testing.T) {
	r := require.New(t)
	a := assert.New(t)

	opts := AppOptions{
		Templates:   []TemplateOptions{{}},
		Inventories: []InventoryOptions{{}},
	}

	err := opts.Validate()
	r.Error(err)

	a.Contains(err.Error(), "template 0 source is undefined")
	a.Contains(err.Error(), "template 0 type is undefined")
	a.Contains(err.Error(), "template 0 name is undefined")

	a.Contains(err.Error(), "inventory 0 source is undefined")
	a.Contains(err.Error(), "inventory 0 type is undefined")
	a.Contains(err.Error(), "inventory 0 namespace is undefined")
}

func Test_fromViperParseTemplate(t *testing.T) {
	r := require.New(t)
	a := assert.New(t)

	tests := map[string]struct {
		TypeStr  string
		TypeEnum TemplateSourceType
		Source   string
		Name     string
	}{
		"Default local template": {
			TypeStr:  "local",
			TypeEnum: TstLocal,
			Source:   "/path/to/template",
			Name:     "test1",
		},
		"Default Git template": {
			TypeStr:  "git",
			TypeEnum: TstGit,
			Source:   "https://some/url",
			Name:     "test1",
		},
	}

	for name, data := range tests {
		t.Run(name, func(t *testing.T) {
			// Given an empty temp folder
			tmpDir, err := os.MkdirTemp("", "")
			r.NoError(err)
			defer os.RemoveAll(tmpDir)

			// And a sample application options file
			cfgFilePath := path.Join(tmpDir, "sample.yml")
			fh, err := os.Create(cfgFilePath)
			r.NoError(err)
			cfgData := fmt.Sprintf(`
templates:
  - type: %s
    source: %s
    name: %s
`, data.TypeStr, data.Source, data.Name)
			_, err = fh.WriteString(cfgData)
			r.NoError(err)

			// Point viper to our config file
			v := viper.New()
			v.SetConfigFile(cfgFilePath)
			r.NoError(v.ReadInConfig())

			// When we try instantiating our app options from Viper
			options, err := FromViper(v)

			// We expect the operation to succeed
			r.NoError(err)
			a.Equal(1, len(options.Templates))
			a.Equal(data.Name, options.Templates[0].GetName())
			a.Equal(data.Source, options.Templates[0].GetSource())
			a.Equal(data.TypeEnum, options.Templates[0].Type)
		})
	}
}

func Test_fromViperParseInventory(t *testing.T) {
	r := require.New(t)
	a := assert.New(t)

	tests := map[string]struct {
		TypeStr   string
		TypeEnum  InventorySourceType
		Source    string
		Namespace string
	}{
		"Default local template": {
			TypeStr:   "local",
			TypeEnum:  IstLocal,
			Source:    "/path/to/template",
			Namespace: "test",
		},
		"Default Git template": {
			TypeStr:   "git",
			TypeEnum:  IstGit,
			Source:    "https://some/url",
			Namespace: "test",
		},
	}

	for name, data := range tests {
		t.Run(name, func(t *testing.T) {
			// Given an empty temp folder
			tmpDir, err := os.MkdirTemp("", "")
			r.NoError(err)
			defer os.RemoveAll(tmpDir)

			// And a sample application options file
			cfgFilePath := path.Join(tmpDir, "sample.yml")
			fh, err := os.Create(cfgFilePath)
			r.NoError(err)
			cfgData := fmt.Sprintf(`
inventories:
  - type: %s
    source: %s
    namespace: %s
`, data.TypeStr, data.Source, data.Namespace)
			_, err = fh.WriteString(cfgData)
			r.NoError(err)

			// Point viper to our config file
			v := viper.New()
			v.SetConfigFile(cfgFilePath)
			r.NoError(v.ReadInConfig())

			// When we try instantiating our app options from Viper
			options, err := FromViper(v)

			// We expect the operation to succeed
			r.NoError(err)
			a.Equal(1, len(options.Inventories))
			a.Equal(data.Namespace, options.Inventories[0].Namespace)
			a.Equal(data.Source, options.Inventories[0].Source)
			a.Equal(data.TypeEnum, options.Inventories[0].Type)
		})
	}
}

func Test_fromViperParseOtherOptions(t *testing.T) {
	r := require.New(t)
	a := assert.New(t)

	tests := map[string]struct {
		ThemeStr  string
		ThemeEnum ThemeType
	}{
		"Dark theme": {
			ThemeStr:  "dark",
			ThemeEnum: ThtDark,
		},
	}

	for name, data := range tests {
		t.Run(name, func(t *testing.T) {
			// Given an empty temp folder
			tmpDir, err := os.MkdirTemp("", "")
			r.NoError(err)
			defer os.RemoveAll(tmpDir)

			// And a sample application options file
			cfgFilePath := path.Join(tmpDir, "sample.yml")
			fh, err := os.Create(cfgFilePath)
			r.NoError(err)
			cfgData := fmt.Sprintf(`
options:
   theme: %s
`, data.ThemeStr)
			_, err = fh.WriteString(cfgData)
			r.NoError(err)

			// Point viper to our config file
			v := viper.New()
			v.SetConfigFile(cfgFilePath)
			r.NoError(v.ReadInConfig())

			// When we try instantiating our app options from Viper
			options, err := FromViper(v)

			// We expect the operation to succeed
			r.NoError(err)
			a.Equal(data.ThemeEnum, options.Other.Theme)
		})
	}
}

func Test_fromViperParseFailThemeName(t *testing.T) {
	r := require.New(t)

	// Given an empty temp folder
	tmpDir, err := os.MkdirTemp("", "")
	r.NoError(err)
	defer os.RemoveAll(tmpDir)

	// config file where theme is a list type instead of a string type
	cfgFilePath := path.Join(tmpDir, "sample.yml")
	fh, err := os.Create(cfgFilePath)
	r.NoError(err)

	cfgData := `
options:
   theme:
     - MyTheme
`
	_, err = fh.WriteString(cfgData)
	r.NoError(err)

	// Point viper to our config file
	v := viper.New()
	v.SetConfigFile(cfgFilePath)
	r.NoError(v.ReadInConfig())

	// When we try instantiating our app options from Viper
	_, err = FromViper(v)

	// We expect the operation to succeed
	r.Error(err)
}

func Test_fromViperParseFail(t *testing.T) {
	r := require.New(t)

	tests := map[string]struct {
		config string
	}{
		"Valid yaml with missing template type": {
			config: `
templates:
  - source: /some/path2
    name: test2
`,
		},
	}

	for name, data := range tests {
		t.Run(name, func(t *testing.T) {
			// Given an empty temp folder
			tmpDir, err := os.MkdirTemp("", "")
			r.NoError(err)
			defer os.RemoveAll(tmpDir)

			// And a sample application options file
			cfgFilePath := path.Join(tmpDir, "sample.yml")
			fh, err := os.Create(cfgFilePath)
			r.NoError(err)
			_, err = fh.WriteString(data.config)
			r.NoError(err)

			// Point viper to our config file
			v := viper.New()
			v.SetConfigFile(cfgFilePath)
			r.NoError(v.ReadInConfig())

			// When we try instantiating our app options from Viper
			_, err = FromViper(v)

			// The operation should fail
			r.Error(err)

		})
	}
}

func Test_findInventoryByNameNoMatch(t *testing.T) {
	appoptions := AppOptions{}
	result := appoptions.FindInventory("fubar")
	require.Nil(t, result)
}

func Test_findInventoryByName(t *testing.T) {
	expectedNamespace := "FuBar"
	expectedInv := InventoryOptions{
		Namespace: expectedNamespace,
		Source:    "http://fubar/repo",
		Type:      IstGit,
	}
	appoptions := AppOptions{
		Inventories: []InventoryOptions{expectedInv},
	}

	result := appoptions.FindInventory(expectedNamespace)
	require.Equal(t, expectedInv, *result)
}
