package applicationOptions

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func Test_TemplateSourceTypeStringConverstion(t *testing.T) {
	r := require.New(t)

	tests := map[string]struct {
		source string
		target TemplateSourceType
	}{
		"Git": {
			source: "git",
			target: TstGit,
		},
		"local": {
			source: "local",
			target: TstLocal,
		},
		"Empty string": {
			source: "",
			target: TstUndefined,
		},
		"Unknown": {
			source: "fubar",
			target: TstUnknown,
		},
	}
	// TODO: Make a list of every enum value used in every test,
	//       and make sure every enum has at least 1 test for it
	r.Equal(int(TstGit)+1, len(tests))
	for name, data := range tests {
		t.Run(name, func(t *testing.T) {
			var temp TemplateSourceType
			temp.fromString(data.source)
			r.Equal(data.target, temp)
		})
	}
}

func Test_TemplateSourceTypeStringCasting(t *testing.T) {
	r := require.New(t)

	tests := map[string]struct {
		target string
		source TemplateSourceType
	}{
		"Git": {
			target: "git",
			source: TstGit,
		},
		"Local": {
			target: "local",
			source: TstLocal,
		},
		"Undefined": {
			target: "",
			source: TstUndefined,
		},
		"Unknown": {
			target: "",
			source: TstUnknown,
		},
	}
	// TODO: Find a better way to make sure we've tested all enumerations
	r.Equal(int(TstGit)+1, len(tests))
	for name, data := range tests {
		t.Run(name, func(t *testing.T) {
			r.Equal(data.target, data.source.toString())
		})
	}
}

func Test_TemplateOptionsBasicGetters(t *testing.T) {
	a := assert.New(t)

	expSource := "/path/to/template"
	expType := TstLocal
	expAlias := "MyTemplate"
	opts := TemplateOptions{
		Source: expSource,
		Type:   expType,
		Alias:  expAlias,
	}

	a.Equal(expType, opts.GetType())
	a.Equal(expSource, opts.GetSource())
	a.Equal(expAlias, opts.GetName())
}

func Test_TemplateOptionsGetFilesystem(t *testing.T) {
	r := require.New(t)
	a := assert.New(t)

	tests := map[string]struct {
		expSource string
		expType   TemplateSourceType
		expFSName string
	}{
		"Git filesystem": {
			expSource: "https://github.com/TheFriendlyCoder/rejigger.git",
			expType:   TstGit,
			expFSName: "MemMapFS",
		},
		"Local filesystem": {
			expSource: os.TempDir(),
			expType:   TstLocal,
			expFSName: "OsFs",
		},
	}

	for name, data := range tests {
		t.Run(name, func(t *testing.T) {
			opts := TemplateOptions{
				Source: data.expSource,
				Type:   data.expType,
				Alias:  "MyTemplate",
			}

			fs, err := opts.GetFilesystem()
			r.NoError(err)
			a.Equal(data.expFSName, fs.Name())
		})
	}
}

func Test_TemplateOptionsGetManifestPath(t *testing.T) {
	a := assert.New(t)

	tests := map[string]struct {
		expSource string
		expPath   string
		expType   TemplateSourceType
	}{
		"Git manifest path": {
			expSource: "https://url/to/repo",
			expPath:   ".rejig.yml",
			expType:   TstGit,
		},
		"Local manifest path": {
			expSource: "/path/to/template",
			expPath:   "/path/to/template/.rejig.yml",
			expType:   TstLocal,
		},
	}

	for name, data := range tests {
		t.Run(name, func(t *testing.T) {

			opts := TemplateOptions{
				Source: data.expSource,
				Type:   data.expType,
				Alias:  "MyTemplate",
			}

			a.Equal(data.expPath, opts.GetManifestPath())
		})
	}
}

func Test_TemplateOptionsGettersPanic(t *testing.T) {
	a := assert.New(t)

	tests := map[string]struct {
		expType TemplateSourceType
	}{
		"Undefined template type": {
			expType: TstUndefined,
		},
		"Unknown template type": {
			expType: TstUnknown,
		},
	}

	for name, data := range tests {
		t.Run(name, func(t *testing.T) {

			opts := TemplateOptions{
				Source: "/some/path",
				Type:   data.expType,
				Alias:  "MyTemplate",
			}
			a.Panics(func() { opts.GetRoot() })
			a.Panics(func() {
				_, err := opts.GetFilesystem()
				a.NoError(err)
			})
			a.Panics(func() { opts.GetManifestPath() })
		})
	}
}

func Test_parseYaml(t *testing.T) {
	r := require.New(t)

	tests := map[string]struct {
		expType  TemplateSourceType
		yamlData string
	}{
		"Git template type": {
			expType:  TstGit,
			yamlData: "git",
		},
		"Local template type": {
			expType:  TstLocal,
			yamlData: "local",
		},
		"Unknown template type": {
			expType:  TstUnknown,
			yamlData: "fubar",
		},
		"Undefined template type": {
			expType:  TstUndefined,
			yamlData: "",
		},
	}

	for name, data := range tests {
		t.Run(name, func(t *testing.T) {

			var actType TemplateSourceType
			r.NoError(yaml.Unmarshal([]byte(data.yamlData), &actType))

			r.Equal(data.expType, actType)
		})
	}
}

func Test_parseInvalidYaml(t *testing.T) {
	r := require.New(t)
	yamlData := "fu: 123"

	var actType TemplateSourceType
	err := yaml.Unmarshal([]byte(yamlData), &actType)
	r.Error(err)
	r.Contains(err.Error(), "Unable to parse template")
}
