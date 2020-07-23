package providers

import (
	"github.com/jexia/semaphore/pkg/core/instance"
	"github.com/jexia/semaphore/pkg/specs"
)

// FlowsResolver when called collects the available flow(s) with the configured configuration
type FlowsResolver func(instance.Context) ([]*specs.FlowsManifest, error)

// EndpointsResolver when called collects the available endpoint(s) with the configured configuration
type EndpointsResolver func(instance.Context) ([]*specs.EndpointsManifest, error)

// ServiceResolvers represents a collection of service resolvers
type ServiceResolvers []ServicesResolver

// Resolve resolves all the given service resolvers and returns a aggregated service list
func (resolvers ServiceResolvers) Resolve(ctx instance.Context) (specs.ServiceList, error) {
	services := make(specs.ServiceList, 0)

	for _, resolver := range resolvers {
		if resolver == nil {
			continue
		}

		result, err := resolver(ctx)
		if err != nil {
			return nil, err
		}

		services.Append(result)
	}

	return services, nil
}

// ServicesResolver when called collects the available service(s) with the configured configuration
type ServicesResolver func(instance.Context) (specs.ServiceList, error)

// SchemaResolvers represents a collection of schema resolvers
type SchemaResolvers []SchemaResolver

// Resolve resolves all schema resolves and returns a aggregated Object
func (resolvers SchemaResolvers) Resolve(ctx instance.Context) (specs.Objects, error) {
	objects := make(specs.Objects)

	for _, resolver := range resolvers {
		if resolver == nil {
			continue
		}

		result, err := resolver(ctx)
		if err != nil {
			return nil, err
		}

		objects.Append(result)
	}

	return objects, nil
}

// SchemaResolver when called collects the available service(s) with the configured configuration
type SchemaResolver func(instance.Context) (specs.Objects, error)
