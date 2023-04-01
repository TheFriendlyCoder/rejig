package cmd

type pathErrorTypes int64

const (
	pathNotFound pathErrorTypes = iota
	pathNotEmpty
)

type pathError struct {
	path      string
	errorType pathErrorTypes
}

func (p pathError) Error() string {
	switch p.errorType {
	case pathNotFound:
		return "Path not found: " + p.path
	case pathNotEmpty:
		return "Path must be empty: " + p.path
	default:
		panic("Unsupported error")
	}
}
