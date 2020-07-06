package maestro

import (
	"github.com/jexia/maestro/internal/codec"
	"github.com/jexia/maestro/pkg/core/api"
	"github.com/jexia/maestro/pkg/core/instance"
	"github.com/jexia/maestro/pkg/core/logger"
	"github.com/jexia/maestro/pkg/functions"
	"github.com/jexia/maestro/pkg/providers"
	"github.com/jexia/maestro/pkg/transport"
)

// NewOptions constructs a api.Options object from the given api.Option constructors
func NewOptions(ctx instance.Context, options ...api.Option) (api.Options, error) {
	result := api.NewOptions(ctx)

	for _, option := range options {
		if option == nil {
			continue
		}

		option(&result)
	}

	for _, middleware := range result.Middleware {
		options, err := middleware(ctx)
		if err != nil {
			return result, err
		}

		for _, option := range options {
			option(&result)
		}
	}

	return result, nil
}

// NewCollection constructs a new options collection
func NewCollection(options ...api.Option) []api.Option {
	return options
}

// WithFlows appends the given flows resolver to the available flow resolvers
func WithFlows(definition providers.FlowsResolver) api.Option {
	return func(options *api.Options) {
		options.Flows = append(options.Flows, definition)
	}
}

// WithServices appends the given service resolver to the available service resolvers
func WithServices(definition providers.ServicesResolver) api.Option {
	return func(options *api.Options) {
		options.Services = append(options.Services, definition)
	}
}

// WithEndpoints appends the given endpoint resolver to the available endpoint resolvers
func WithEndpoints(definition providers.EndpointsResolver) api.Option {
	return func(options *api.Options) {
		options.Endpoints = append(options.Endpoints, definition)
	}
}

// WithSchema appends the schema collection to the schema store
func WithSchema(resolver providers.SchemaResolver) api.Option {
	return func(options *api.Options) {
		options.Schemas = append(options.Schemas, resolver)
	}
}

// WithCodec appends the given codec to the collection of available codecs
func WithCodec(codec codec.Constructor) api.Option {
	return func(options *api.Options) {
		options.Codec[codec.Name()] = codec
	}
}

// WithCaller appends the given caller to the collection of available callers
func WithCaller(caller transport.NewCaller) api.Option {
	return func(options *api.Options) {
		options.Callers = append(options.Callers, caller(options.Ctx))
	}
}

// WithListener appends the given listener to the collection of available listeners
func WithListener(listener transport.NewListener) api.Option {
	return func(options *api.Options) {
		options.Listeners = append(options.Listeners, listener(options.Ctx))
	}
}

// WithFunctions defines the custom defined functions to be used
func WithFunctions(custom functions.Custom) api.Option {
	return func(options *api.Options) {
		if options.Functions == nil {
			options.Functions = functions.Custom{}
		}

		for key, fn := range custom {
			options.Functions[key] = fn
		}
	}
}

// WithLogLevel sets the log level for the given module
func WithLogLevel(module logger.Module, level string) api.Option {
	return func(options *api.Options) {
		err := options.Ctx.SetLevel(module, level)
		if err != nil {
			options.Ctx.Logger(logger.Core).Warnf("unable to set the logging level, %s", err)
		}
	}
}
