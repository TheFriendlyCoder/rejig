package lib

import (
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_getTemplate(t *testing.T) {
	r := require.New(t)

	tmp, err := GetGitTemplate("https://github.com/TheFriendlyCoder/rejigger.git")
	r.NoError(err)

	res, err := afero.ReadDir(tmp, ".")
	r.NoError(err)
	r.True(len(res) > 0)
}
