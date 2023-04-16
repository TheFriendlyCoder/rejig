package applicationOptions

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ThemeTypeStringConversion(t *testing.T) {
	r := require.New(t)

	tests := map[string]struct {
		source string
		target ThemeType
	}{
		"Light": {
			source: "light",
			target: ThtLight,
		},
		"Dark": {
			source: "dark",
			target: ThtDark,
		},
		"Empty string": {
			source: "",
			target: ThtUndefined,
		},
		"Unknown": {
			source: "fubar",
			target: ThtUnknown,
		},
	}
	// TODO: Make a list of every enum value used in every test,
	//       and make sure every enum has at least 1 test for it
	r.Equal(int(ThtDark)+1, len(tests))
	for name, data := range tests {
		t.Run(name, func(t *testing.T) {
			var temp ThemeType
			temp.FromString(data.source)
			r.Equal(data.target, temp)
		})
	}
}

func Test_ThemeTypeStringCasting(t *testing.T) {
	r := require.New(t)

	tests := map[string]struct {
		target string
		source ThemeType
	}{
		"Light": {
			target: "light",
			source: ThtLight,
		},
		"Dark": {
			target: "dark",
			source: ThtDark,
		},
		"Undefined": {
			target: "",
			source: ThtUndefined,
		},
		"Unknown": {
			target: "",
			source: ThtUnknown,
		},
	}
	// TODO: Find a better way to make sure we've tested all enumerations
	r.Equal(int(ThtDark)+1, len(tests))
	for name, data := range tests {
		t.Run(name, func(t *testing.T) {
			r.Equal(data.target, data.source.ToString())
		})
	}
}
