package applicationOptions

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_TemplateSourceTypeStringConverstion(t *testing.T) {
	r := require.New(t)

	tests := map[string]struct {
		source string
		target TemplateSourceType
	}{
		"Git": {
			source: "git",
			target: TstGit,
		},
		"local": {
			source: "local",
			target: TstLocal,
		},
		"Empty string": {
			source: "",
			target: TstUndefined,
		},
		"Unknown": {
			source: "fubar",
			target: TstUnknown,
		},
	}
	// TODO: Make a list of every enum value used in every test,
	//       and make sure every enum has at least 1 test for it
	r.Equal(int(TstGit)+1, len(tests))
	for name, data := range tests {
		t.Run(name, func(t *testing.T) {
			var temp TemplateSourceType
			temp.fromString(data.source)
			r.Equal(data.target, temp)
		})
	}
}

func Test_TemplateSourceTypeStringCasting(t *testing.T) {
	r := require.New(t)

	tests := map[string]struct {
		target string
		source TemplateSourceType
	}{
		"Git": {
			target: "git",
			source: TstGit,
		},
		"Local": {
			target: "local",
			source: TstLocal,
		},
		"Undefined": {
			target: "",
			source: TstUndefined,
		},
		"Unknown": {
			target: "",
			source: TstUnknown,
		},
	}
	// TODO: Find a better way to make sure we've tested all enumerations
	r.Equal(int(TstGit)+1, len(tests))
	for name, data := range tests {
		t.Run(name, func(t *testing.T) {
			r.Equal(data.target, data.source.toString())
		})
	}
}
