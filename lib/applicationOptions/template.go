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
	allNames := map[string]int{}
	for i, curTemplate := range a.Templates {
		if len(curTemplate.GetName()) == 0 {
			retval = append(retval, fmt.Sprintf("template %d name is undefined", i))
		}
		allNames[curTemplate.GetName()] += 1

		if curTemplate.Type == TstUndefined {
			retval = append(retval, fmt.Sprintf("template %d type is undefined", i))
		} else if curTemplate.Type == TstUnknown {
			retval = append(retval, fmt.Sprintf("template %d type is not supported", i))
		}
		if len(curTemplate.GetSource()) == 0 {
			retval = append(retval, fmt.Sprintf("template %d source is undefined", i))
		}
	}

	// See if any template names are duplicated
	for name, count := range allNames {
		if count > 1 {
			retval = append(retval, fmt.Sprintf("there are %d templates with the name %s", count, name))
		}
	}
	return retval
}
