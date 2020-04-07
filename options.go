package maestro

import (
	"github.com/jexia/maestro/codec"
	"github.com/jexia/maestro/internal/constructor"
	"github.com/jexia/maestro/internal/instance"
	"github.com/jexia/maestro/internal/logger"
	"github.com/jexia/maestro/schema"
	"github.com/jexia/maestro/specs"
	"github.com/jexia/maestro/transport"
)

// NewOptions constructs a constructor.Options object from the given constructor.Option constructors
func NewOptions(ctx instance.Context, options ...constructor.Option) constructor.Options {
	result := constructor.Options{
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
func WithDefinitions(definition specs.Resolver) constructor.Option {
	return func(options *constructor.Options) {
		options.Definitions = append(options.Definitions, definition)
	}
}

// WithCodec appends the given codec to the collection of available codecs
func WithCodec(codec codec.Constructor) constructor.Option {
	return func(options *constructor.Options) {
		options.Codec[codec.Name()] = codec
	}
}

// WithCaller appends the given caller to the collection of available callers
func WithCaller(caller transport.NewCaller) constructor.Option {
	return func(options *constructor.Options) {
		options.Callers = append(options.Callers, caller(options.Ctx))
	}
}

// WithListener appends the given listener to the collection of available listeners
func WithListener(listener transport.NewListener) constructor.Option {
	return func(options *constructor.Options) {
		options.Listeners = append(options.Listeners, listener(options.Ctx))
	}
}

// WithSchema appends the schema collection to the schema store
func WithSchema(resolver schema.Resolver) constructor.Option {
	return func(options *constructor.Options) {
		options.Schemas = append(options.Schemas, resolver)
	}
}

// WithFunctions defines the custom defined functions to be used
func WithFunctions(functions specs.CustomDefinedFunctions) constructor.Option {
	return func(options *constructor.Options) {
		if options.Functions == nil {
			options.Functions = specs.CustomDefinedFunctions{}
		}

		for key, fn := range functions {
			options.Functions[key] = fn
		}
	}
}

// WithLogLevel sets the log level for the given module
func WithLogLevel(module logger.Module, level string) constructor.Option {
	return func(options *constructor.Options) {
		err := options.Ctx.SetLevel(module, level)
		if err != nil {
			// TODO: handle error
		}
	}
}
