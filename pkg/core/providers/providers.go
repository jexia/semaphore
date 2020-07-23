package providers

import (
	"github.com/jexia/semaphore/pkg/core/api"
	"github.com/jexia/semaphore/pkg/core/instance"
	"github.com/jexia/semaphore/pkg/specs"
)

// Resolve calls all defined resolvers and returns a specs collection
func Resolve(ctx instance.Context, options api.Options) (*specs.Collection, error) {
	flows := &specs.FlowsManifest{}
	endpoints := &specs.EndpointsManifest{}
	services := &specs.ServicesManifest{}
	schemas := &specs.SchemaManifest{}

	for _, resolver := range options.Flows {
		if resolver == nil {
			continue
		}

		manifests, err := resolver(ctx)
		if err != nil {
			return nil, err
		}

		flows.Append(manifests...)
	}

	for _, resolver := range options.Endpoints {
		if resolver == nil {
			continue
		}

		manifests, err := resolver(ctx)
		if err != nil {
			return nil, err
		}

		endpoints.Append(manifests...)
	}

	for _, resolver := range options.Services {
		if resolver == nil {
			continue
		}

		manifests, err := resolver(ctx)
		if err != nil {
			return nil, err
		}

		services.Append(manifests...)
	}

	for _, resolver := range options.Schemas {
		if resolver == nil {
			continue
		}

		manifests, err := resolver(ctx)
		if err != nil {
			return nil, err
		}

		schemas.Append(manifests...)
	}

	result := &specs.Collection{
		FlowsManifest:     flows,
		EndpointsManifest: endpoints,
		ServicesManifest:  services,
		SchemaManifest:    schemas,
	}

	return result, nil
}
