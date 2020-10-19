package providers

import (
	"github.com/jexia/semaphore/pkg/broker"
	"github.com/jexia/semaphore/pkg/specs"
)

// FlowsResolvers represents a collection of flows resolvers
type FlowsResolvers []FlowsResolver

// Resolve resolvers the flows and returns a aggregated response
func (resolvers FlowsResolvers) Resolve(ctx *broker.Context) (specs.FlowListInterface, error) {
	flows := specs.FlowListInterface{}

	for _, resolver := range resolvers {
		if resolver == nil {
			continue
		}

		result, err := resolver(ctx)
		if err != nil {
			return nil, err
		}

		flows.Append(result)
	}

	return flows, nil
}

// FlowsResolver when called collects the available flow(s) with the configured configuration
type FlowsResolver func(*broker.Context) (specs.FlowListInterface, error)

// EndpointResolvers represents a collection of endpoint resolvers
type EndpointResolvers []EndpointsResolver

// Resolve resolves the endpoint resolvers collection and returns a aggregated response
func (resolvers EndpointResolvers) Resolve(ctx *broker.Context) (specs.EndpointList, error) {
	endpoints := specs.EndpointList{}

	for _, resolver := range resolvers {
		if resolver == nil {
			continue
		}

		result, err := resolver(ctx)
		if err != nil {
			return nil, err
		}

		endpoints.Append(result)
	}

	return endpoints, nil
}

// EndpointsResolver when called collects the available endpoint(s) with the configured configuration
type EndpointsResolver func(*broker.Context) (specs.EndpointList, error)

// ServiceResolvers represents a collection of service resolvers
type ServiceResolvers []ServicesResolver

// Resolve resolves all the given service resolvers and returns a aggregated service list
func (resolvers ServiceResolvers) Resolve(ctx *broker.Context) (specs.ServiceList, error) {
	services := specs.ServiceList{}

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
type ServicesResolver func(*broker.Context) (specs.ServiceList, error)

// SchemaResolvers represents a collection of schema resolvers
type SchemaResolvers []SchemaResolver

// Resolve resolves all schema resolves and returns a aggregated Object
func (resolvers SchemaResolvers) Resolve(ctx *broker.Context) (specs.Schemas, error) {
	objects := specs.Schemas{}

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
type SchemaResolver func(*broker.Context) (specs.Schemas, error)
