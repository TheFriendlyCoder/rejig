package lib

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_getTemplate(t *testing.T) {
	r := require.New(t)

	tmp, err := getTemplate("https://github.com/TheFriendlyCoder/rejigger.git")
	fs := *tmp
	r.NoError(err, "Should have worked")
	res, err := fs.ReadDir(".")
	r.NoError(err, "Can't read from in memory file system")
	r.True(len(res) > 0)
}
