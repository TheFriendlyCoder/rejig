package lib

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_successfulValidation(t *testing.T) {
	a := assert.New(t)

	// TODO: Convert this to a table test
	allTypes := [...]TemplateSourceType{
		TstLocal,
		TstGit,
	}

	for _, t := range allTypes {
		opts := AppOptions{
			Templates: []TemplateOptions{{
				Alias:  "My Template",
				Source: "https://some/location",
				Type:   t,
			}},
		}

		err := opts.Validate()
		a.NoError(err)
	}
}

func Test_successfulValidationEmptyConfig(t *testing.T) {
	r := require.New(t)

	opts := AppOptions{}

	err := opts.Validate()
	r.NoError(err, "Validation should have succeeded")
}

func Test_successfulValidationWithoutOptionals(t *testing.T) {
	r := require.New(t)

	opts := AppOptions{
		Templates: []TemplateOptions{{
			Alias:  "My Template",
			Source: "https://some/location",
			Type:   TstGit,
		}},
	}

	err := opts.Validate()
	r.NoError(err)
}

func Test_validationTemplateWithoutType(t *testing.T) {
	r := require.New(t)

	opts := AppOptions{
		Templates: []TemplateOptions{{
			Alias:  "My Template",
			Source: "https://some/location",
		}},
	}

	err := opts.Validate()
	r.Error(err)
}

func Test_validationTemplateWithoutAlias(t *testing.T) {
	r := require.New(t)

	opts := AppOptions{
		Templates: []TemplateOptions{{
			Source: "https://some/location",
			Type:   TstGit,
		}},
	}

	err := opts.Validate()
	r.Error(err)
}

func Test_validationTemplateWithoutSource(t *testing.T) {
	r := require.New(t)

	opts := AppOptions{
		Templates: []TemplateOptions{{
			Alias: "My Template",
			Type:  TstGit,
		}},
	}

	err := opts.Validate()
	r.Error(err)
}

func Test_validationTemplateCompoundError(t *testing.T) {
	r := require.New(t)
	a := assert.New(t)

	opts := AppOptions{
		Templates: []TemplateOptions{{}},
	}

	err := opts.Validate()
	r.Error(err)

	a.Contains(err.Error(), "source is undefined")
	a.Contains(err.Error(), "type is undefined")
	a.Contains(err.Error(), "alias is undefined")
}
