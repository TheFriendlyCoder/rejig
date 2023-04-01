package cmd

type pathErrorTypes int64

const (
	undefined pathErrorTypes = iota
	pathNotFound
	pathNotEmpty
	fileNotFound
)

type pathError struct {
	path      string
	errorType pathErrorTypes
}

func (p pathError) Error() string {
	var retval string
	switch p.errorType {
	case pathNotFound:
		retval = "Path not found: " + p.path
	case pathNotEmpty:
		retval = "Path must be empty: " + p.path
	case fileNotFound:
		retval = "File not found: " + p.path
	case undefined:
		panic("Unsupported path error")
	}
	return retval
}
