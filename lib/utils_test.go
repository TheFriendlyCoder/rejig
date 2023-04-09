package lib

import (
	"fmt"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var someerr = fmt.Errorf("Some Failure")

func Test_SNFshouldPanic(t *testing.T) {
	a := assert.New(t)

	a.Panics(func() { SNF(someerr) })
	a.Panics(func() { SNF(someerr, "hello") })
	a.Panics(func() { SNF("hello", someerr) })
}

func Test_SNFshouldNotPanic(t *testing.T) {
	a := assert.New(t)
	a.NotPanics(func() { SNF("hello") })
	a.NotPanics(func() { SNF("hello", 16) })
}

func Test_getGitFilesystem(t *testing.T) {
	r := require.New(t)

	tmp, err := GetGitFilesystem("https://github.com/TheFriendlyCoder/rejiggerTestTemplate.git")
	r.NoError(err)

	res, err := afero.ReadDir(tmp, ".")
	r.NoError(err)
	r.True(len(res) > 0)
}
