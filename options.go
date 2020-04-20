package maestro

import (
	"github.com/jexia/maestro/internal/constructor"
	"github.com/jexia/maestro/internal/logger"
	"github.com/jexia/maestro/pkg/codec"
	"github.com/jexia/maestro/pkg/definitions"
	"github.com/jexia/maestro/pkg/functions"
	"github.com/jexia/maestro/pkg/instance"
	"github.com/jexia/maestro/pkg/transport"
)

// NewOptions constructs a constructor.Options object from the given constructor.Option constructors
func NewOptions(ctx instance.Context, options ...constructor.Option) constructor.Options {
	result := constructor.NewOptions(ctx)

	for _, option := range options {
		option(&result)
	}

	return result
}

// WithFlows appends the given flows resolver to the available flow resolvers
func WithFlows(definition definitions.FlowsResolver) constructor.Option {
	return func(options *constructor.Options) {
		options.Flows = append(options.Flows, definition)
	}
}

// WithServices appends the given service resolver to the available service resolvers
func WithServices(definition definitions.ServicesResolver) constructor.Option {
	return func(options *constructor.Options) {
		options.Services = append(options.Services, definition)
	}
}

// WithEndpoints appends the given endpoint resolver to the available endpoint resolvers
func WithEndpoints(definition definitions.EndpointsResolver) constructor.Option {
	return func(options *constructor.Options) {
		options.Endpoints = append(options.Endpoints, definition)
	}
}

// WithSchema appends the schema collection to the schema store
func WithSchema(resolver definitions.SchemaResolver) constructor.Option {
	return func(options *constructor.Options) {
		options.Schemas = append(options.Schemas, resolver)
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

// WithFunctions defines the custom defined functions to be used
func WithFunctions(custom functions.Custom) constructor.Option {
	return func(options *constructor.Options) {
		if options.Functions == nil {
			options.Functions = functions.Custom{}
		}

		for key, fn := range custom {
			options.Functions[key] = fn
		}
	}
}

// WithLogLevel sets the log level for the given module
func WithLogLevel(module logger.Module, level string) constructor.Option {
	return func(options *constructor.Options) {
		err := options.Ctx.SetLevel(module, level)
		if err != nil {
			options.Ctx.Logger(logger.Core).Warnf("unable to set the logging level, %s", err)
		}
	}
}
