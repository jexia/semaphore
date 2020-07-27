package semaphore

import (
	"github.com/jexia/semaphore/pkg/codec"
	"github.com/jexia/semaphore/pkg/core"
	"github.com/jexia/semaphore/pkg/core/api"
	"github.com/jexia/semaphore/pkg/core/instance"
	"github.com/jexia/semaphore/pkg/core/logger"
	"github.com/jexia/semaphore/pkg/functions"
	"github.com/jexia/semaphore/pkg/providers"
	"github.com/jexia/semaphore/pkg/transport"
)

// DefaultOptions sets the default options for the given options object
func DefaultOptions(options *api.Options) {
	options.Constructor = core.Construct
}

// NewOptions constructs a api.Options object from the given api.Option constructors
func NewOptions(ctx instance.Context, options ...api.Option) (api.Options, error) {
	result := api.NewOptions(ctx)
	DefaultOptions(&result)

	if options == nil {
		return result, nil
	}

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

// WithConstructor appends the given flows resolver to the available flow resolvers
func WithConstructor(constructor api.Constructor) api.Option {
	return func(options *api.Options) {
		options.Constructor = constructor
	}
}

// WithFlows appends the given flows resolver to the available flow resolvers
func WithFlows(definition providers.FlowsResolver) api.Option {
	return func(options *api.Options) {
		options.FlowResolvers = append(options.FlowResolvers, definition)
	}
}

// WithServices appends the given service resolver to the available service resolvers
func WithServices(definition providers.ServicesResolver) api.Option {
	return func(options *api.Options) {
		options.ServiceResolvers = append(options.ServiceResolvers, definition)
	}
}

// WithEndpoints appends the given endpoint resolver to the available endpoint resolvers
func WithEndpoints(definition providers.EndpointsResolver) api.Option {
	return func(options *api.Options) {
		options.EndpointResolvers = append(options.EndpointResolvers, definition)
	}
}

// WithSchema appends the schema collection to the schema store
func WithSchema(resolver providers.SchemaResolver) api.Option {
	return func(options *api.Options) {
		options.SchemaResolvers = append(options.SchemaResolvers, resolver)
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
