package lib

// SNF ShouldNotFail helper method that looks through all return values from a function call,
// detects any error objects that are returned, and if any of them are failure objects then
// it triggers a panic operation, since errors should not typically happen in practice
// Example use:
//
//	lib.SNF(fmt.FPrintf(...))
func SNF(args ...any) {
	// loop through all input args
	for _, c := range args {
		// see if any can be cast to an error type
		if e, ok := c.(error); ok {
			// and if so, we know it is a non-nil error result
			// so we want to trigger a panic
			panic("Unexpected error: " + e.Error())
		}
	}
}
