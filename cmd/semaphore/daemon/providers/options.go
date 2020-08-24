package providers

import (
	"github.com/jexia/semaphore"
	"github.com/jexia/semaphore/pkg/broker"
	"github.com/jexia/semaphore/pkg/providers"
	"github.com/jexia/semaphore/pkg/transport"
)

// NewOptions constructs a Options object from the given Option constructors
func NewOptions(ctx *broker.Context, options ...Option) (Options, error) {
	root, _ := semaphore.NewOptions(ctx)
	result := Options{
		Options:           root,
		Listeners:         transport.ListenerList{},
		EndpointResolvers: providers.EndpointResolvers{},
		ServiceResolvers:  providers.ServiceResolvers{},
		SchemaResolvers:   providers.SchemaResolvers{},
	}

	if options == nil {
		return result, nil
	}

	for _, option := range options {
		if option == nil {
			continue
		}

		option(ctx, &result)
	}

	return result, nil
}

// Option represents a constructor func which sets a given option
type Option func(*broker.Context, *Options)

// Options represents the available options to resolve the given providers
type Options struct {
	semaphore.Options
	Listeners         transport.ListenerList
	EndpointResolvers providers.EndpointResolvers
	ServiceResolvers  providers.ServiceResolvers
	SchemaResolvers   providers.SchemaResolvers
}

// WithCore sets the given semaphore Options
func WithCore(root semaphore.Options) Option {
	return func(ctx *broker.Context, options *Options) {
		options.Options = root
	}
}

// WithServices appends the given service resolver to the available service resolvers
func WithServices(definition providers.ServicesResolver) Option {
	return func(ctx *broker.Context, options *Options) {
		options.ServiceResolvers = append(options.ServiceResolvers, definition)
	}
}

// WithEndpoints appends the given endpoint resolver to the available endpoint resolvers
func WithEndpoints(definition providers.EndpointsResolver) Option {
	return func(ctx *broker.Context, options *Options) {
		options.EndpointResolvers = append(options.EndpointResolvers, definition)
	}
}

// WithSchema appends the schema collection to the schema store
func WithSchema(resolver providers.SchemaResolver) Option {
	return func(ctx *broker.Context, options *Options) {
		options.SchemaResolvers = append(options.SchemaResolvers, resolver)
	}
}

// WithListener appends the given listener to the collection of available listeners
func WithListener(listener transport.NewListener) Option {
	return func(ctx *broker.Context, options *Options) {
		options.Listeners = append(options.Listeners, listener(ctx))
	}
}
