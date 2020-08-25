package providers

import (
	"github.com/jexia/semaphore"
	"github.com/jexia/semaphore/pkg/broker"
	"github.com/jexia/semaphore/pkg/providers"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/transport"
)

// NewOptions constructs a Options object from the given Option constructors
func NewOptions(ctx *broker.Context, core semaphore.Options, options ...Option) (Options, error) {
	result := Options{
		Options:           core,
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
	AfterConstructor  AfterConstructor
}

// AfterConstructor is called after the specifications is constructored
type AfterConstructor func(*broker.Context, specs.FlowListInterface, specs.EndpointList, specs.ServiceList, specs.Schemas) error

// AfterConstructorHandler wraps the after constructed function to allow middleware to be chained
type AfterConstructorHandler func(AfterConstructor) AfterConstructor

// WithAfterConstructor the passed function gets called once all options have been applied
func WithAfterConstructor(wrapper AfterConstructorHandler) Option {
	return func(ctx *broker.Context, options *Options) {
		if options.AfterConstructor == nil {
			options.AfterConstructor = wrapper(func(*broker.Context, specs.FlowListInterface, specs.EndpointList, specs.ServiceList, specs.Schemas) error {
				return nil
			})
			return
		}

		options.AfterConstructor = wrapper(options.AfterConstructor)
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
