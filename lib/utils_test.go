package lib

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
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
