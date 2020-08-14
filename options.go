package semaphore

import (
	"github.com/jexia/semaphore/pkg/broker"
	"github.com/jexia/semaphore/pkg/broker/config"
	"github.com/jexia/semaphore/pkg/broker/logger"
	"github.com/jexia/semaphore/pkg/codec"
	"github.com/jexia/semaphore/pkg/functions"
	"github.com/jexia/semaphore/pkg/providers"
	"github.com/jexia/semaphore/pkg/transport"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewOptions constructs a config.Options object from the given config.Option constructors
func NewOptions(ctx *broker.Context, options ...config.Option) (config.Options, error) {
	result := config.NewOptions()

	if options == nil {
		return result, nil
	}

	for _, option := range options {
		if option == nil {
			continue
		}

		option(ctx, &result)
	}

	for _, middleware := range result.Middleware {
		options, err := middleware(ctx)
		if err != nil {
			return result, err
		}

		for _, option := range options {
			option(ctx, &result)
		}
	}

	return result, nil
}

// NewCollection constructs a new options collection
func NewCollection(options ...config.Option) []config.Option {
	return options
}

// WithFlows appends the given flows resolver to the available flow resolvers
func WithFlows(definition providers.FlowsResolver) config.Option {
	return func(ctx *broker.Context, options *config.Options) {
		options.FlowResolvers = append(options.FlowResolvers, definition)
	}
}

// WithServices appends the given service resolver to the available service resolvers
func WithServices(definition providers.ServicesResolver) config.Option {
	return func(ctx *broker.Context, options *config.Options) {
		options.ServiceResolvers = append(options.ServiceResolvers, definition)
	}
}

// WithEndpoints appends the given endpoint resolver to the available endpoint resolvers
func WithEndpoints(definition providers.EndpointsResolver) config.Option {
	return func(ctx *broker.Context, options *config.Options) {
		options.EndpointResolvers = append(options.EndpointResolvers, definition)
	}
}

// WithSchema appends the schema collection to the schema store
func WithSchema(resolver providers.SchemaResolver) config.Option {
	return func(ctx *broker.Context, options *config.Options) {
		options.SchemaResolvers = append(options.SchemaResolvers, resolver)
	}
}

// WithCodec appends the given codec to the collection of available codecs
func WithCodec(codec codec.Constructor) config.Option {
	return func(ctx *broker.Context, options *config.Options) {
		options.Codec[codec.Name()] = codec
	}
}

// WithCaller appends the given caller to the collection of available callers
func WithCaller(caller transport.NewCaller) config.Option {
	return func(ctx *broker.Context, options *config.Options) {
		options.Callers = append(options.Callers, caller(ctx))
	}
}

// WithListener appends the given listener to the collection of available listeners
func WithListener(listener transport.NewListener) config.Option {
	return func(ctx *broker.Context, options *config.Options) {
		options.Listeners = append(options.Listeners, listener(ctx))
	}
}

// WithFunctions defines the custom defined functions to be used
func WithFunctions(custom functions.Custom) config.Option {
	return func(ctx *broker.Context, options *config.Options) {
		if options.Functions == nil {
			options.Functions = functions.Custom{}
		}

		for key, fn := range custom {
			options.Functions[key] = fn
		}
	}
}

// WithLogLevel sets the log level for the given module
func WithLogLevel(pattern string, value string) config.Option {
	return func(ctx *broker.Context, options *config.Options) {
		level := zapcore.InfoLevel
		err := level.UnmarshalText([]byte(value))
		if err != nil {
			logger.Error(ctx, "unable to unmarshal log level", zap.String("level", value))
		}

		logger.SetLevel(ctx, pattern, level)
	}
}
