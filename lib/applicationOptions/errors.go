package applicationOptions

import "strings"

// -=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=
//										AppOptionsError

type AppOptionsError struct {
	Messages []string
}

func (e AppOptionsError) Error() string {
	if len(e.Messages) == 1 {
		return "Failed to parse application option: " + e.Messages[0]
	}
	retval := strings.Join(e.Messages, "\n\t")
	return "Failed to parse application options:\n\t" + retval
}
