package constructor

import (
	"context"

	"github.com/jexia/maestro/codec"
	"github.com/jexia/maestro/logger"
	"github.com/jexia/maestro/protocol"
	"github.com/jexia/maestro/schema"
	"github.com/jexia/maestro/specs"
)

// Option represents a constructor func which sets a given option
type Option func(*Options)

// Options represents all the available options
type Options struct {
	Ctx         context.Context
	Definitions []specs.Resolver
	Codec       map[string]codec.Constructor
	Callers     protocol.Callers
	Listeners   protocol.Listeners
	Schemas     []schema.Resolver
	Schema      *schema.Store
	Functions   specs.CustomDefinedFunctions
}

// NewOptions constructs a options object from the given option constructors
func NewOptions(ctx context.Context, options ...Option) Options {
	result := Options{
		Ctx:         ctx,
		Definitions: make([]specs.Resolver, 0),
		Codec:       make(map[string]codec.Constructor),
		Schema:      schema.NewStore(ctx),
	}

	for _, option := range options {
		option(&result)
	}

	return result
}

// WithDefinitions defines the HCL definitions path to be used
func WithDefinitions(definition specs.Resolver) Option {
	return func(options *Options) {
		options.Definitions = append(options.Definitions, definition)
	}
}

// WithCodec appends the given codec to the collection of available codecs
func WithCodec(constructor codec.Constructor) Option {
	return func(options *Options) {
		options.Codec[constructor.Name()] = constructor
	}
}

// WithCaller appends the given caller to the collection of available callers
func WithCaller(caller protocol.Caller) Option {
	return func(options *Options) {
		caller.Context(options.Ctx)
		options.Callers = append(options.Callers, caller)
	}
}

// WithListener appends the given listener to the collection of available listeners
func WithListener(listener protocol.Listener) Option {
	return func(options *Options) {
		listener.Context(options.Ctx)
		options.Listeners = append(options.Listeners, listener)
	}
}

// WithSchema appends the schema collection to the schema store
func WithSchema(resolver schema.Resolver) Option {
	return func(options *Options) {
		options.Schemas = append(options.Schemas, resolver)
	}
}

// WithFunctions defines the custom defined functions to be used
func WithFunctions(functions specs.CustomDefinedFunctions) Option {
	return func(options *Options) {
		options.Functions = functions
	}
}

// WithLogLevel sets the log level for the given module
func WithLogLevel(module logger.Module, level string) Option {
	return func(options *Options) {
		err := logger.SetLevel(options.Ctx, module, level)
		if err != nil {
			// TODO: handle error
		}
	}
}
