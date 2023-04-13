package applicationOptions

import (
	"fmt"

	e "github.com/TheFriendlyCoder/rejigger/lib/errors"
)

// decodeTemplateOptions decodes raw YAMl data into proper parsed template options
func decodeTemplateOptions(raw interface{}) (map[string]interface{}, error) {
	// Map the "type" field from a character string format to an enumerated type
	templateData, ok := raw.(map[string]interface{})
	if !ok {
		return nil, e.AOTemplateOptionsDecodeError()
	}

	var newVal TemplateSourceType
	temp, ok := templateData["type"].(string)
	if !ok {
		return nil, e.AOTemplateOptionsDecodeError()
	}
	newVal.fromString(temp)
	templateData["type"] = newVal
	return templateData, nil
}

// validateTemplates checks to make sure the template options in our application
// options meets the application requirements
func (a AppOptions) validateTemplates() []string {
	var retval []string
	for i, curTemplate := range a.Templates {
		if len(curTemplate.GetName()) == 0 {
			retval = append(retval, fmt.Sprintf("template %d name is undefined", i))
		}
		if curTemplate.Type == TstUndefined {
			retval = append(retval, fmt.Sprintf("template %d type is undefined", i))
		} else if curTemplate.Type == TstUnknown {
			retval = append(retval, fmt.Sprintf("template %d type is not supported", i))
		}
		if len(curTemplate.Source) == 0 {
			retval = append(retval, fmt.Sprintf("template %d source is undefined", i))
		}
	}
	return retval
}
