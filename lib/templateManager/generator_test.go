package templateManager

import (
	"fmt"
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

	gitFS, err := getGitTemplate("https://github.com/TheFriendlyCoder/rejiggerTestTemplate.git")
	r.NoError(err)

	tests := map[string]struct {
		fileSystem afero.Fs
		sourceDir  string
	}{
		"Local file system template": {
			fileSystem: afero.NewOsFs(),
			sourceDir:  testProjectDir("simple"),
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
			a.True(unmodified(fs, r, exp, act))

			act = filepath.Join(tmpDir, "version.txt")
			a.FileExists(act)
			a.True(contains(r, act, expVersion))
			a.False(contains(r, act, "{{version}}"))

			act = filepath.Join(tmpDir, "MyProj", "main.txt")
			a.FileExists(act)
			a.True(contains(r, act, expProj))
			a.False(contains(r, act, "{{project_name}}"))
		})
	}
}
