package core

import (
	"github.com/jexia/semaphore/pkg/checks"
	"github.com/jexia/semaphore/pkg/compare"
	"github.com/jexia/semaphore/pkg/core/api"
	"github.com/jexia/semaphore/pkg/core/instance"
	"github.com/jexia/semaphore/pkg/dependencies"
	"github.com/jexia/semaphore/pkg/functions"
	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/references/forwarding"
	"github.com/jexia/semaphore/pkg/schema"
	"github.com/jexia/semaphore/pkg/specs"
)

// Construct construct a specs manifest from the given options.
// The specifications are received from the providers. The property types are defined and functions are prepared.
// Once done is a specs collection returned that could be used to update the listeners.
func Construct(ctx instance.Context, mem functions.Collection, options api.Options) (specs.FlowListInterface, specs.EndpointList, specs.ServiceList, specs.Objects, error) {
	if options.BeforeConstructor != nil {
		err := options.BeforeConstructor(ctx, mem, options)
		if err != nil {
			return nil, nil, nil, nil, err
		}
	}

	flows, err := options.FlowResolvers.Resolve(ctx)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	endpoints, err := options.EndpointResolvers.Resolve(ctx)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	schemas, err := options.SchemaResolvers.Resolve(ctx)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	services, err := options.ServiceResolvers.Resolve(ctx)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	err = checks.FlowDuplicates(ctx, flows)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	err = schema.Define(ctx, services, schemas, flows)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	err = references.Resolve(ctx, flows)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	err = functions.PrepareManifestFunctions(ctx, mem, options.Functions, flows)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	forwarding.ResolveReferences(ctx, flows)

	err = dependencies.ResolveFlows(ctx, flows)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	err = compare.Types(ctx, services, schemas, flows)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	if options.AfterConstructor != nil {
		err = options.AfterConstructor(ctx, flows, endpoints, services, schemas)
		if err != nil {
			return nil, nil, nil, nil, err
		}
	}

	return flows, endpoints, services, schemas, nil
}
