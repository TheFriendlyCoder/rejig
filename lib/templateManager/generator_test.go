package templateManager

import (
	"fmt"
	"github.com/TheFriendlyCoder/rejigger/lib"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
)

func Test_basicGenerator(t *testing.T) {
	r := require.New(t)
	a := assert.New(t)

	gitFS, err := lib.GetGitFilesystem("https://github.com/TheFriendlyCoder/rejiggerTestTemplate.git")
	r.NoError(err)

	tests := map[string]struct {
		fileSystem afero.Fs
		sourceDir  string
	}{
		"Local file system template": {
			fileSystem: afero.NewOsFs(),
			sourceDir:  getProjectDir("simple"),
		},
		"Git file system template": {
			fileSystem: gitFS,
			sourceDir:  ".",
		},
	}

	for name, data := range tests {
		t.Run(name, func(t *testing.T) {

			srcPath := data.sourceDir

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
			err = generate(fs, srcPath, tmpDir, context)

			r.NoError(err, "Failed to run generator")

			a.DirExists(filepath.Join(tmpDir, "MyProj"))
			a.NoFileExists(filepath.Join(tmpDir, ".rejig.yml"))

			exp := filepath.Join(srcPath, ".gitignore")
			act := filepath.Join(tmpDir, ".gitignore")
			fmt.Println(act)
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
