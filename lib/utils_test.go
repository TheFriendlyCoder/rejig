package lib

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_SNFshouldPanic(t *testing.T) {
	a := assert.New(t)

	a.Panics(func() { SNF(fmt.Errorf("Some Failure")) })
	a.Panics(func() { SNF(fmt.Errorf("Some Failure"), "hello") })
	a.Panics(func() { SNF("hello", fmt.Errorf("Some Failure")) })
}

func Test_SNFshouldNotPanic(t *testing.T) {
	a := assert.New(t)
	a.NotPanics(func() { SNF("hello") })
	a.NotPanics(func() { SNF("hello", 16) })
}
