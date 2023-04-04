package lib

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_successfulValidation(t *testing.T) {
	r := require.New(t)

	opts := AppOptions{
		Templates: []TemplateOptions{{
			Alias:  "My Template",
			Source: "https://some/location",
			Type:   TST_GIT,
			Folder: "subfolder",
		}},
	}

	err := opts.Validate()
	r.NoError(err, "Validation should have succeeded")
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
			Type:   TST_GIT,
		}},
	}

	err := opts.Validate()
	r.NoError(err, "Validation should have succeeded")
}

func Test_validationTemplateWithoutType(t *testing.T) {
	r := require.New(t)

	opts := AppOptions{
		Templates: []TemplateOptions{{
			Alias:  "My Template",
			Source: "https://some/location",
			Folder: "subfolder",
		}},
	}

	err := opts.Validate()
	r.Error(err, "Validation should have failed")
}

func Test_validationTemplateWithoutAlias(t *testing.T) {
	r := require.New(t)

	opts := AppOptions{
		Templates: []TemplateOptions{{
			Source: "https://some/location",
			Type:   TST_GIT,
			Folder: "subfolder",
		}},
	}

	err := opts.Validate()
	r.Error(err, "Validation should have failed")
}

func Test_validationTemplateWithoutSource(t *testing.T) {
	r := require.New(t)

	opts := AppOptions{
		Templates: []TemplateOptions{{
			Alias:  "My Template",
			Type:   TST_GIT,
			Folder: "subfolder",
		}},
	}

	err := opts.Validate()
	r.Error(err, "Validation should have failed")
}

func Test_validationTemplateCompoundError(t *testing.T) {
	r := require.New(t)
	a := assert.New(t)

	opts := AppOptions{
		Templates: []TemplateOptions{{}},
	}

	err := opts.Validate()
	r.Error(err, "Validation should have failed")

	a.Contains(err.Error(), "source is undefined")
	a.Contains(err.Error(), "type is undefined")
	a.Contains(err.Error(), "alias is undefined")
}
