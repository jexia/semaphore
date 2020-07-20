package core

import (
	"github.com/jexia/semaphore/pkg/core/api"
	"github.com/jexia/semaphore/pkg/core/instance"
	"github.com/jexia/semaphore/pkg/specs"
)

// CollectSpecs calls all defined resolvers and returns a specs collection
func CollectSpecs(ctx instance.Context, options api.Options) (*api.Collection, error) {
	result := &api.Collection{
		Flows:     &specs.FlowsManifest{},
		Endpoints: &specs.EndpointsManifest{},
		Services:  &specs.ServicesManifest{},
		Schema:    &specs.SchemaManifest{},
	}

	for _, resolver := range options.Flows {
		if resolver == nil {
			continue
		}

		manifests, err := resolver(ctx)
		if err != nil {
			return nil, err
		}

		specs.MergeFlowsManifest(result.Flows, manifests...)
	}

	for _, resolver := range options.Endpoints {
		if resolver == nil {
			continue
		}

		manifests, err := resolver(ctx)
		if err != nil {
			return nil, err
		}

		specs.MergeEndpointsManifest(result.Endpoints, manifests...)
	}

	for _, resolver := range options.Services {
		if resolver == nil {
			continue
		}

		manifests, err := resolver(ctx)
		if err != nil {
			return nil, err
		}

		specs.MergeServiceManifest(result.Services, manifests...)
	}

	for _, resolver := range options.Schemas {
		if resolver == nil {
			continue
		}

		manifests, err := resolver(ctx)
		if err != nil {
			return nil, err
		}

		specs.MergeSchemaManifest(result.Schema, manifests...)
	}

	return result, nil
}
