package shared

// ContextKey enumerated type defining keys in the Cobra context manager used to store
// and retrieve common command properties
type ContextKey int64

const (
	// CkOptions Parsed application options loaded from the environment or app config file
	// should be managed exclusively by the root command
	CkOptions ContextKey = iota
	// CkViper Gets the shared Viper config file manager used by all subcommands
	CkViper
)
