package providers

import (
	"github.com/jexia/semaphore/pkg/core/api"
	"github.com/jexia/semaphore/pkg/core/instance"
	"github.com/jexia/semaphore/pkg/specs"
)

// Resolve calls all defined resolvers and returns a specs collection
func Resolve(ctx instance.Context, options api.Options) (*specs.Collection, error) {
	result := &specs.Collection{
		FlowsManifest:     &specs.FlowsManifest{},
		EndpointsManifest: &specs.EndpointsManifest{},
		ServicesManifest:  &specs.ServicesManifest{},
		SchemaManifest:    &specs.SchemaManifest{},
	}

	for _, resolver := range options.Flows {
		if resolver == nil {
			continue
		}

		manifests, err := resolver(ctx)
		if err != nil {
			return nil, err
		}

		specs.MergeFlowsManifest(result.FlowsManifest, manifests...)
	}

	for _, resolver := range options.Endpoints {
		if resolver == nil {
			continue
		}

		manifests, err := resolver(ctx)
		if err != nil {
			return nil, err
		}

		specs.MergeEndpointsManifest(result.EndpointsManifest, manifests...)
	}

	for _, resolver := range options.Services {
		if resolver == nil {
			continue
		}

		manifests, err := resolver(ctx)
		if err != nil {
			return nil, err
		}

		specs.MergeServiceManifest(result.ServicesManifest, manifests...)
	}

	for _, resolver := range options.Schemas {
		if resolver == nil {
			continue
		}

		manifests, err := resolver(ctx)
		if err != nil {
			return nil, err
		}

		specs.MergeSchemaManifest(result.SchemaManifest, manifests...)
	}

	return result, nil
}
