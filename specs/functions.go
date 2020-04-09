package specs

// CustomDefinedFunctions represents a collection of custom defined functions that could be called inside a template
type CustomDefinedFunctions map[string]PrepareFunction

// PrepareFunction prepares the custom defined function.
// The given arguments represent the exprected types that are passed when called.
// Properties returned should be absolute.
type PrepareFunction func(args ...*Property) (*Property, FunctionExec, error)

// FunctionExec executes the function and passes the expected types as stores
// A store should be returned which could be used to encode the function property
type FunctionExec func(store Store) error

// Functions represents a collection of functions
type Functions map[string]*Function

// Function represents a custom defined function
type Function struct {
	Arguments []*Property
	Fn        FunctionExec
	Returns   *Property
}
