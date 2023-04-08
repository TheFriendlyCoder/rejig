package lib

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_successfulValidation(t *testing.T) {
	a := assert.New(t)

	tests := map[string]struct {
		templateOptions  []TemplateOptions
		inventoryOptions []InventoryOptions
	}{
		"Local template file system type": {
			templateOptions: []TemplateOptions{{
				Alias:  "My Template",
				Source: "https://some/location",
				Type:   TstLocal,
			}},
		},
		"Git template file system type": {
			templateOptions: []TemplateOptions{{
				Alias:  "My Template",
				Source: "https://some/location",
				Type:   TstGit,
			}},
		},
		"Local inventory file system type": {
			inventoryOptions: []InventoryOptions{{
				Namespace: "Fubar",
				Source:    "https://some/location",
				Type:      IstLocal,
			}},
		},
		"Git inventory file system type": {
			inventoryOptions: []InventoryOptions{{
				Namespace: "Fubar",
				Source:    "https://some/location",
				Type:      IstGit,
			}},
		},
	}

	for name, data := range tests {
		t.Run(name, func(t *testing.T) {

			opts := AppOptions{
				Templates:   data.templateOptions,
				Inventories: data.inventoryOptions,
			}
			err := opts.Validate()
			a.NoError(err)
		})
	}
}

func Test_successfulValidationEmptyConfig(t *testing.T) {
	r := require.New(t)

	opts := AppOptions{}

	err := opts.Validate()
	r.NoError(err, "Validation should have succeeded")
}

func Test_validationFailures(t *testing.T) {
	r := require.New(t)

	tests := map[string]struct {
		templateOptions  []TemplateOptions
		inventoryOptions []InventoryOptions
	}{
		"Template missing type": {
			templateOptions: []TemplateOptions{{
				Alias:  "My Template",
				Source: "https://some/location",
			}},
		},
		"Template missing alias": {
			templateOptions: []TemplateOptions{{
				Source: "https://some/location",
				Type:   TstLocal,
			}},
		},
		"Template missing source": {
			templateOptions: []TemplateOptions{{
				Alias: "My Template",
				Type:  TstGit,
			}},
		},
		"Inventory missing type": {
			inventoryOptions: []InventoryOptions{{
				Namespace: "Fubar",
				Source:    "https://some/location",
			}},
		},
		"Inventory missing source": {
			inventoryOptions: []InventoryOptions{{
				Namespace: "Fubar",
				Type:      IstLocal,
			}},
		},
		"Inventory missing namespsce": {
			inventoryOptions: []InventoryOptions{{
				Source: "https://some/location",
				Type:   IstLocal,
			}},
		},
	}

	for name, data := range tests {
		t.Run(name, func(t *testing.T) {
			opts := AppOptions{
				Templates:   data.templateOptions,
				Inventories: data.inventoryOptions,
			}
			err := opts.Validate()
			r.Error(err)
		})
	}
}

func Test_validationTemplateCompoundError(t *testing.T) {
	r := require.New(t)
	a := assert.New(t)

	opts := AppOptions{
		Templates:   []TemplateOptions{{}},
		Inventories: []InventoryOptions{{}},
	}

	err := opts.Validate()
	r.Error(err)

	a.Contains(err.Error(), "template 0 source is undefined")
	a.Contains(err.Error(), "template 0 type is undefined")
	a.Contains(err.Error(), "template 0 alias is undefined")

	a.Contains(err.Error(), "inventory 0 source is undefined")
	a.Contains(err.Error(), "inventory 0 type is undefined")
	a.Contains(err.Error(), "inventory 0 namespace is undefined")
}
