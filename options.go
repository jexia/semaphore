package maestro

import "github.com/jexia/maestro/constructor"

// WithDefinitions defines the HCL definitions path to be used
var WithDefinitions = constructor.WithDefinitions

// WithCodec appends the given codec to the collection of available codecs
var WithCodec = constructor.WithCodec

// WithCaller appends the given caller to the collection of available callers
var WithCaller = constructor.WithCaller

// WithListener appends the given listener to the collection of available listeners
var WithListener = constructor.WithListener

// WithSchema appends the schema collection to the schema store
var WithSchema = constructor.WithSchema

// WithLogLevel sets the log level for the given module
var WithLogLevel = constructor.WithLogLevel

// WithFunctions defines the custom defined functions to be used
var WithFunctions = constructor.WithFunctions
