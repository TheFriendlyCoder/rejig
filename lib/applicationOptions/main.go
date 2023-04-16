package applicationOptions

import (
	"reflect"

	e "github.com/TheFriendlyCoder/rejigger/lib/errors"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// TODO: replace this exported symbol with an interface of the same name

// AppOptions parsed config options supported by the app
type AppOptions struct {
	// Templates 0 or more sources where template projects are to be found
	Templates   []TemplateOptions
	Inventories []InventoryOptions
}

// FromViper constructs a new set of application options from a Viper config file
func FromViper(v *viper.Viper) (AppOptions, error) {
	// TODO: think about how to handle "default options" here
	//		 should Viper control all default options?
	var retval AppOptions
	err := v.Unmarshal(&retval, viper.DecodeHook(appOptionsDecoder()))
	if err != nil {
		return retval, errors.WithStack(err)
	}

	// Then validate the results to make sure they meet the application requirements
	return retval, errors.Wrap(retval.Validate(), "Failed decoding application options")
}

// Validate checks the contents of the parsed application options to make sure they
// meet the requirements for the application
func (a AppOptions) Validate() error {
	var messages []string
	messages = append(messages, a.validateTemplates()...)
	messages = append(messages, a.validateInventory()...)
	if len(messages) == 0 {
		return nil
	}
	return e.NewAppOptionsError(messages)
}

// FindInventory locates a template inventory given the namespace name
func (a AppOptions) FindInventory(namespace string) *InventoryOptions {
	for _, curInventory := range a.Inventories {
		if curInventory.Namespace == namespace {
			return &curInventory
		}
	}
	return nil
}

// appOptionsDecoder custom hook method used to translate raw config data into a structure
// that is easier to leverage in the application code
func appOptionsDecoder() mapstructure.DecodeHookFuncType {
	// Based on example found here:
	//		https://sagikazarmark.hu/blog/decoding-custom-formats-with-viper/
	return func(
		src reflect.Type,
		target reflect.Type,
		raw interface{},
	) (interface{}, error) {

		// TODO: Find a way to detect partial / incomplete parse matches
		// ie: if a template option is missing one field, viper won't map
		//		it to the correct type and it just gets ignored
		// TODO: Find a way to enable strict mode decoding here
		//		 that might work better
		// TODO: Replace application config parser with simple YAML parsing
		//		 because it seems simpler
		//		https://github.com/spf13/viper/issues/338
		if (target == reflect.TypeOf(TemplateOptions{})) {
			newData, err := decodeTemplateOptions(raw)
			return newData, err
		}

		if (target == reflect.TypeOf(InventoryOptions{})) {
			newData, err := decodeInventoryOptions(raw)
			return newData, err
		}
		return raw, nil

	}
}
