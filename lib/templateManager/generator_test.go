package templateManager

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/TheFriendlyCoder/rejigger/lib"
	ao "github.com/TheFriendlyCoder/rejigger/lib/applicationOptions"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_basicGenerator(t *testing.T) {
	r := require.New(t)
	a := assert.New(t)

	gitFS, err := lib.GetGitFilesystem("https://github.com/TheFriendlyCoder/rejigger.git")
	r.NoError(err)

	tests := map[string]struct {
		fileSystem afero.Fs
		sourceDir  string
	}{
		"Local file system template": {
			fileSystem: afero.NewOsFs(),
			sourceDir:  getProjectDir(),
		},
		"Git file system template": {
			fileSystem: gitFS,
			sourceDir:  "testdata/projects/simple",
		},
	}

	for name, data := range tests {
		t.Run(name, func(t *testing.T) {

			options := ao.TemplateOptions{
				Source: data.sourceDir,
				Type:   ao.TstLocal,
				Name:   "MyTemplate",
			}

			// Given an empty temp folder
			tmpDir, err := os.MkdirTemp("", "")
			r.NoError(err)
			defer os.RemoveAll(tmpDir)

			// We attempt to run the generator
			expVersion := "1.6.9"
			expProj := "MyProj"
			context := map[string]any{
				"project_name": expProj,
				"version":      expVersion,
			}
			fs := data.fileSystem
			err = generate(fs, options, tmpDir, context)

			r.NoError(err, "Failed to run generator")

			a.DirExists(filepath.Join(tmpDir, "MyProj"))
			a.NoFileExists(filepath.Join(tmpDir, ".rejig.yml"))

			exp := filepath.Join(data.sourceDir, ".gitignore")
			act := filepath.Join(tmpDir, ".gitignore")
			a.FileExists(act)
			a.True(isUnmodified(fs, r, exp, act))

			act = filepath.Join(tmpDir, "version.txt")
			a.FileExists(act)
			a.True(fileContains(r, act, expVersion))
			a.False(fileContains(r, act, "{{version}}"))

			act = filepath.Join(tmpDir, "MyProj", "main.txt")
			a.FileExists(act)
			a.True(fileContains(r, act, expProj))
			a.False(fileContains(r, act, "{{project_name}}"))
		})
	}
}

func Test_generateWithExclusions(t *testing.T) {
	r := require.New(t)
	a := assert.New(t)

	fileSystem := afero.NewOsFs()
	sourceDir := getProjectDir()

	options := ao.TemplateOptions{
		Source:     sourceDir,
		Type:       ao.TstLocal,
		Name:       "MyTemplate",
		Exclusions: []string{".*/main.txt"},
	}

	// Given an empty temp folder
	tmpDir, err := os.MkdirTemp("", "")
	r.NoError(err)
	defer os.RemoveAll(tmpDir)

	// We attempt to run the generator
	expVersion := "1.6.9"
	expProj := "MyProj"
	context := map[string]any{
		"project_name": expProj,
		"version":      expVersion,
	}
	fs := fileSystem
	err = generate(fs, options, tmpDir, context)

	r.NoError(err, "Failed to run generator")

	a.DirExists(filepath.Join(tmpDir, "MyProj"))
	a.NoFileExists(filepath.Join(tmpDir, ".rejig.yml"))

	exp := filepath.Join(sourceDir, ".gitignore")
	act := filepath.Join(tmpDir, ".gitignore")
	a.FileExists(act)
	a.True(isUnmodified(fs, r, exp, act))

	act = filepath.Join(tmpDir, "version.txt")
	a.FileExists(act)
	a.True(fileContains(r, act, expVersion))
	a.False(fileContains(r, act, "{{version}}"))

	act = filepath.Join(tmpDir, "MyProj", "main.txt")
	a.NoFileExists(act)
}
