package applicationOptions

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_appOptionsErrorSingleMessage(t *testing.T) {
	a := assert.New(t)
	expectedMessage := "Something went wrong"
	err := AppOptionsError{Messages: []string{expectedMessage}}
	a.Contains(err.Error(), "Failed to parse application option:")
	a.Contains(err.Error(), expectedMessage)
}

func Test_appOptionsErrorMultiMessage(t *testing.T) {
	a := assert.New(t)
	expectedMessage1 := "Something went wrong"
	expectedMessage2 := "Something else wrong"
	err := AppOptionsError{Messages: []string{expectedMessage1, expectedMessage2}}
	a.Contains(err.Error(), "Failed to parse application options:")
	a.Contains(err.Error(), expectedMessage1)
	a.Contains(err.Error(), expectedMessage2)
}
