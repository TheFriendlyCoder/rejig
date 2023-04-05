package lib

import (
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_getTemplate(t *testing.T) {
	r := require.New(t)

	tmp, err := GetTemplate("https://github.com/TheFriendlyCoder/rejigger.git")
	fs := *tmp
	r.NoError(err, "Should have worked")
	res, err := afero.ReadDir(fs, ".")
	r.NoError(err, "Can't read from in memory file system")
	r.True(len(res) > 0)
}
