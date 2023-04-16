package applicationOptions

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_InventorySourceTypeStringConversion(t *testing.T) {
	r := require.New(t)

	tests := map[string]struct {
		source string
		target InventorySourceType
	}{
		"Git": {
			source: "git",
			target: IstGit,
		},
		"local": {
			source: "local",
			target: IstLocal,
		},
		"Empty string": {
			source: "",
			target: IstUndefined,
		},
		"Unknown": {
			source: "fubar",
			target: IstUnknown,
		},
	}
	// TODO: Make a list of every enum value used in every test,
	//       and make sure every enum has at least 1 test for it
	r.Equal(int(IstGit)+1, len(tests))
	for name, data := range tests {
		t.Run(name, func(t *testing.T) {
			var temp InventorySourceType
			temp.fromString(data.source)
			r.Equal(data.target, temp)
		})
	}
}

func Test_InventorySourceTypeStringCasting(t *testing.T) {
	r := require.New(t)

	tests := map[string]struct {
		target string
		source InventorySourceType
	}{
		"Git": {
			target: "git",
			source: IstGit,
		},
		"Local": {
			target: "local",
			source: IstLocal,
		},
		"Undefined": {
			target: "",
			source: IstUndefined,
		},
		"Unknown": {
			target: "",
			source: IstUnknown,
		},
	}
	// TODO: Find a better way to make sure we've tested all enumerations
	r.Equal(int(IstGit)+1, len(tests))
	for name, data := range tests {
		t.Run(name, func(t *testing.T) {
			r.Equal(data.target, data.source.toString())
		})
	}
}

func Test_InventoryOptionsBasicGetters(t *testing.T) {
	a := assert.New(t)

	expSource := "/path/to/template"
	expType := IstLocal
	expNamespace := "MyNamespace"
	opts := InventoryOptions{
		Source:    expSource,
		Type:      expType,
		Namespace: expNamespace,
	}

	a.Equal(expType, opts.GetType())
	a.Equal(expSource, opts.GetSource())
	a.Equal(expNamespace, opts.GetNamespace())
}

func Test_InventoryOptionsGetFilesystem(t *testing.T) {
	r := require.New(t)
	a := assert.New(t)

	tests := map[string]struct {
		expSource string
		expType   InventorySourceType
		expFSName string
	}{
		"Git filesystem": {
			expSource: "https://github.com/TheFriendlyCoder/rejigger.git",
			expType:   IstGit,
			expFSName: "MemMapFS",
		},
		"Local filesystem": {
			expSource: os.TempDir(),
			expType:   IstLocal,
			expFSName: "OsFs",
		},
	}

	for name, data := range tests {
		t.Run(name, func(t *testing.T) {
			opts := InventoryOptions{
				Source:    data.expSource,
				Type:      data.expType,
				Namespace: "MyNamespace",
			}

			fs, err := opts.GetFilesystem()
			r.NoError(err)
			a.Equal(data.expFSName, fs.Name())
		})
	}
}

func Test_InventoryOptionsGetInventoryFile(t *testing.T) {
	a := assert.New(t)

	tests := map[string]struct {
		expSource string
		expPath   string
		expType   InventorySourceType
	}{
		"Git manifest path": {
			expSource: "https://url/to/repo",
			expPath:   ".rejig.inv.yml",
			expType:   IstGit,
		},
		"Local manifest path": {
			expSource: "/path/to/template",
			expPath:   "/path/to/template/.rejig.inv.yml",
			expType:   IstLocal,
		},
	}

	for name, data := range tests {
		t.Run(name, func(t *testing.T) {

			opts := InventoryOptions{
				Source:    data.expSource,
				Type:      data.expType,
				Namespace: "MyNamespace",
			}

			a.Equal(data.expPath, opts.GetInventoryFile())
		})
	}
}

func Test_InventoryOptionsGettersPanic(t *testing.T) {
	a := assert.New(t)

	tests := map[string]struct {
		expType InventorySourceType
	}{
		"Undefined template type": {
			expType: IstUndefined,
		},
		"Unknown template type": {
			expType: IstUnknown,
		},
	}

	for name, data := range tests {
		t.Run(name, func(t *testing.T) {

			opts := InventoryOptions{
				Source:    "/some/path",
				Type:      data.expType,
				Namespace: "MyNamespace",
			}
			a.Panics(func() { opts.GetRoot() })
			a.Panics(func() {
				_, err := opts.GetFilesystem()
				a.NoError(err)
			})
			a.Panics(func() { opts.GetInventoryFile() })
		})
	}
}

func Test_getLocalTemplateDefinitions(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)

	// Given an empty temp folder
	tmpDir, err := os.MkdirTemp("", "")
	r.NoError(err)
	defer os.RemoveAll(tmpDir)

	// And a sample config file
	outputFile := path.Join(tmpDir, inventoryFileName)
	expName := "test1"
	expNamespace := "FuBar"
	expType := TstLocal
	expSource := "my/subdir"
	invData := fmt.Sprintf(`
templates:
  - name: %s
    source: %s
`, expName, expSource)
	fh, err := os.Create(outputFile)
	r.NoError(err)
	_, err = fh.WriteString(invData)
	r.NoError(err)
	r.NoError(fh.Close())

	inv := InventoryOptions{
		Type:      IstLocal,
		Namespace: expNamespace,
		Source:    tmpDir,
	}
	opts, err := inv.GetTemplateDefinitions()
	r.NoError(err)
	a.Equal(1, len(opts))
	a.Equal(expType, opts[0].GetType())
	a.Equal(path.Join(tmpDir, expSource), opts[0].GetRoot())
	a.Equal(expName, opts[0].GetName())
}

func Test_getLocalTemplateDefinitionsFailures(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)

	// Given an empty temp folder
	tmpDir, err := os.MkdirTemp("", "")
	r.NoError(err)
	defer os.RemoveAll(tmpDir)

	subdir := path.Join(tmpDir, "sample")
	r.NoError(os.Mkdir(subdir, 0700))
	outputFile := path.Join(subdir, inventoryFileName)
	r.NoError(os.WriteFile(outputFile, []byte("this isn't valid YAML !!?!"), 0600))

	tests := map[string]struct {
		expType   InventorySourceType
		expSource string
		expError  string
	}{
		"Invalid inventory source": {
			expType:   IstGit,
			expError:  "Failed to load remote Git repository",
			expSource: "https://github.com/TheFriendlyCoder/doesnotexist",
		},
		"Missing inventory file": {
			expType:   IstLocal,
			expError:  "no such file or directory",
			expSource: tmpDir,
		},
		"Invalid inventory file": {
			expType:   IstLocal,
			expError:  "Failed to parse YAML content",
			expSource: subdir,
		},
	}

	for name, data := range tests {
		t.Run(name, func(t *testing.T) {

			inv := InventoryOptions{
				Type:      data.expType,
				Namespace: "FuBar",
				Source:    data.expSource,
			}
			_, err = inv.GetTemplateDefinitions()
			r.Error(err)
			a.Contains(err.Error(), data.expError)
		})
	}
}
