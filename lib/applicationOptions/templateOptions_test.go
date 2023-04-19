package applicationOptions

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func Test_TemplateSourceTypeStringConversion(t *testing.T) {
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
	expName := "MyName"
	opts := TemplateOptions{
		Source: expSource,
		Type:   expType,
		Name:   expName,
	}

	a.Equal(expType, opts.GetType())
	a.Equal(expName, opts.GetName())
	a.Equal(expSource, opts.GetSource())
}

func Test_TemplateOptionsIsFileExcluded(t *testing.T) {
	r := require.New(t)

	tests := map[string]struct {
		exclusions []string
		result     bool
		filename   string
	}{
		"No exclusions": {
			exclusions: []string{},
			result:     false,
			filename:   "something.txt",
		},
		"Exclusions don't match": {
			exclusions: []string{".*/other.txt"},
			result:     false,
			filename:   "something.txt",
		},
		"Regex matches": {
			exclusions: []string{".*/other.txt$"},
			result:     true,
			filename:   "some/path/other.txt",
		},
		"Regex partial match fails": {
			exclusions: []string{".*/other.txt$"},
			result:     false,
			filename:   "some/path/other.txt/folder",
		},
	}

	for name, data := range tests {
		t.Run(name, func(t *testing.T) {
			opts := TemplateOptions{
				Source:     "mypath",
				Name:       "MyTemplate",
				Type:       TstLocal,
				Exclusions: data.exclusions,
			}
			r.Equal(data.result, opts.IsFileExcluded(data.filename))
		})
	}
}

func Test_TemplateOptionsIsFileExcludedInvalidRegex(t *testing.T) {
	r := require.New(t)

	opts := TemplateOptions{
		Source:     "mypath",
		Name:       "MyTemplate",
		Type:       TstLocal,
		Exclusions: []string{"**/notvalid/*?."},
	}
	r.Panics(func() { opts.IsFileExcluded("main.txt") })
}

func Test_TemplateOptionsGetSourceHomeFolder(t *testing.T) {
	r := require.New(t)
	a := assert.New(t)

	home, err := os.UserHomeDir()
	r.NoError(err)

	tests := map[string]struct {
		inSource  string
		inType    TemplateSourceType
		expSource string
	}{
		"Local template with home folder": {
			inSource:  "~/template",
			inType:    TstLocal,
			expSource: path.Join(home, "template"),
		},
		"Git template with home folder": {
			inSource:  "~/template",
			inType:    TstGit,
			expSource: "~/template",
		},
	}

	for name, data := range tests {
		t.Run(name, func(t *testing.T) {

			opts := TemplateOptions{
				Source: data.inSource,
				Type:   data.inType,
				Name:   "MyName",
			}

			a.Equal(data.expSource, opts.GetSource())
		})
	}
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
				Name:   "MyTemplate",
			}

			fs, err := opts.GetFilesystem()
			r.NoError(err)
			a.Equal(data.expFSName, fs.Name())
		})
	}
}

func Test_TemplateOptionsGetManifestFile(t *testing.T) {
	a := assert.New(t)

	tests := map[string]struct {
		expSource string
		expPath   string
		expSubdir string
		expType   TemplateSourceType
	}{
		"Git manifest path no subfolder": {
			expSource: "https://url/to/repo",
			expPath:   ".rejig.yml",
			expType:   TstGit,
			expSubdir: "",
		},
		"Git manifest path with subfolder": {
			expSource: "https://url/to/repo",
			expPath:   "fubar/.rejig.yml",
			expType:   TstGit,
			expSubdir: "fubar",
		},
		"Local manifest path no subdir": {
			expSource: "/path/to/template",
			expPath:   "/path/to/template/.rejig.yml",
			expType:   TstLocal,
			expSubdir: "",
		},
		"Local manifest path with subdir": {
			expSource: "/path/to/template",
			expPath:   "/path/to/template/fubar/.rejig.yml",
			expType:   TstLocal,
			expSubdir: "fubar",
		},
	}

	for name, data := range tests {
		t.Run(name, func(t *testing.T) {

			opts := TemplateOptions{
				Source: data.expSource,
				Type:   data.expType,
				Name:   "MyTemplate",
				Root:   data.expSubdir,
			}

			a.Equal(data.expPath, opts.GetManifestFile())
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
				Name:   "MyTemplate",
			}
			a.Panics(func() { opts.GetRoot() })
			a.Panics(func() {
				_, err := opts.GetFilesystem()
				a.NoError(err)
			})
			a.Panics(func() { opts.GetManifestFile() })
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
