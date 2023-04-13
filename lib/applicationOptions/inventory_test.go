package applicationOptions

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_parseInventoryFile(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)

	// Given an empty temp folder
	tmpDir, err := os.MkdirTemp("", "")
	r.NoError(err)
	defer os.RemoveAll(tmpDir)

	// And a sample config file
	outputFile := path.Join(tmpDir, ".rejig.inv.yml")
	expAlias := "test1"
	expType := TstGit
	expTypeStr := "git"
	expSource := "http://some/repo"
	invData := fmt.Sprintf(`
templates:
  - alias: %s
    source: %s
    type: %s
`, expAlias, expSource, expTypeStr)
	fh, err := os.Create(outputFile)
	r.NoError(err)
	_, err = fh.WriteString(invData)
	r.NoError(err)
	r.NoError(fh.Close())

	// When we try parsing in the inventory
	fs := afero.NewOsFs()
	result, err := parseInventory(fs, outputFile)

	// The inventory should be successfully parsed
	r.NoError(err)
	r.Equal(1, len(result.Templates))
	a.Equal(expType, result.Templates[0].Type)
	a.Equal(expSource, result.Templates[0].Source)
	a.Equal(expAlias, result.Templates[0].GetName())
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
	expAlias := "test1"
	expType := TstGit
	expTypeStr := "git"
	expSource := "http://some/repo"
	invData := fmt.Sprintf(`
templates:
  - alias: %s
    source: %s
    type: %s
`, expAlias, expSource, expTypeStr)
	fh, err := os.Create(outputFile)
	r.NoError(err)
	_, err = fh.WriteString(invData)
	r.NoError(err)
	r.NoError(fh.Close())

	inv := InventoryOptions{
		Type:      IstLocal,
		Namespace: "FuBar",
		Source:    tmpDir,
	}
	opts, err := inv.GetTemplateDefinitions()
	r.NoError(err)
	a.Equal(1, len(opts))
	a.Equal(expType, opts[0].Type)
	a.Equal(expSource, opts[0].Source)
	a.Equal(expAlias, opts[0].GetName())
}
