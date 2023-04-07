package lib

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_getTemplate(t *testing.T) {
	r := require.New(t)

	tmp, err := getTemplate("https://github.com/TheFriendlyCoder/rejigger.git")
	r.NoError(err)

	res, err := tmp.ReadDir(".")
	r.NoError(err)
	r.True(len(res) > 0)
}
