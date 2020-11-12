package providers

import (
	"github.com/jexia/semaphore/pkg/broker"
	"github.com/jexia/semaphore/pkg/checks"
	"github.com/jexia/semaphore/pkg/compare"
	"github.com/jexia/semaphore/pkg/dependencies"
	"github.com/jexia/semaphore/pkg/functions"
	"github.com/jexia/semaphore/pkg/providers"
	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/references/forwarding"
	"github.com/jexia/semaphore/pkg/specs"
)

// Collection represents a collection of specification lists and objects.
// These objects could be used to initialize a Semaphore broker.
type Collection struct {
	specs.FlowListInterface
	specs.EndpointList
	specs.ServiceList
	specs.Schemas
	specs.ServiceDiscoveryClients
}

// Resolve collects and constructs the a specs from the given options.
// The specifications are received from the providers. The property types are
// defined and functions are prepared. Once done is a specs collection returned
// that could be used to update the listeners.
func Resolve(ctx *broker.Context, mem functions.Collection, options Options) (Collection, error) {
	if options.BeforeConstructor != nil {
		err := options.BeforeConstructor(ctx, mem, options.Options)
		if err != nil {
			return Collection{}, err
		}
	}

	flows, err := options.FlowResolvers.Resolve(ctx)
	if err != nil {
		return Collection{}, err
	}

	endpoints, err := options.EndpointResolvers.Resolve(ctx)
	if err != nil {
		return Collection{}, err
	}

	schemas, err := options.SchemaResolvers.Resolve(ctx)
	if err != nil {
		return Collection{}, err
	}

	serviceDiscoveryClients, err := options.DiscoveryServiceResolvers.Resolve(ctx)
	if err != nil {
		return Collection{}, err
	}

	services, err := options.ServiceResolvers.Resolve(ctx)
	if err != nil {
		return Collection{}, err
	}

	err = checks.FlowDuplicates(ctx, flows)
	if err != nil {
		return Collection{}, err
	}

	err = providers.ResolveSchemaDefinitions(ctx, services, schemas, flows)
	if err != nil {
		return Collection{}, err
	}

	err = functions.PrepareFunctions(ctx, mem, options.Functions, flows)
	if err != nil {
		return Collection{}, err
	}

	err = references.Resolve(ctx, flows)
	if err != nil {
		return Collection{}, err
	}

	forwarding.ResolveReferences(ctx, flows, mem)

	err = dependencies.ResolveFlows(ctx, flows)
	if err != nil {
		return Collection{}, err
	}

	err = compare.Types(ctx, services, schemas, flows)
	if err != nil {
		return Collection{}, err
	}

	if options.AfterConstructor != nil {
		err = options.AfterConstructor(ctx, flows, endpoints, services, schemas)
		if err != nil {
			return Collection{}, err
		}
	}

	result := Collection{
		FlowListInterface:       flows,
		EndpointList:            endpoints,
		ServiceList:             services,
		Schemas:                 schemas,
		ServiceDiscoveryClients: serviceDiscoveryClients,
	}

	return result, nil
}
