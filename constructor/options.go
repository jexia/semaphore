package constructor

import (
	"github.com/jexia/maestro/codec"
	"github.com/jexia/maestro/instance"
	"github.com/jexia/maestro/logger"
	"github.com/jexia/maestro/schema"
	"github.com/jexia/maestro/specs"
	"github.com/jexia/maestro/transport"
)

// Option represents a constructor func which sets a given option
type Option func(*Options)

// Options represents all the available options
type Options struct {
	Ctx         instance.Context
	Definitions []specs.Resolver
	Codec       codec.Constructors
	Callers     transport.Callers
	Listeners   transport.Listeners
	Schemas     []schema.Resolver
	Schema      *schema.Store
	Functions   specs.CustomDefinedFunctions
}

// NewOptions constructs a options object from the given option constructors
func NewOptions(ctx instance.Context, options ...Option) Options {
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
func WithCaller(caller transport.NewCaller) Option {
	return func(options *Options) {
		options.Callers = append(options.Callers, caller(options.Ctx))
	}
}

// WithListener appends the given listener to the collection of available listeners
func WithListener(listener transport.NewListener) Option {
	return func(options *Options) {
		options.Listeners = append(options.Listeners, listener(options.Ctx))
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
		err := options.Ctx.SetLevel(module, level)
		if err != nil {
			// TODO: handle error
		}
	}
}
