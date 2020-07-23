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

	for _, resolver := range options.FlowResolvers {
		if resolver == nil {
			continue
		}

		manifests, err := resolver(ctx)
		if err != nil {
			return nil, err
		}

		flows.Append(manifests...)
	}

	for _, resolver := range options.EndpointResolvers {
		if resolver == nil {
			continue
		}

		manifests, err := resolver(ctx)
		if err != nil {
			return nil, err
		}

		endpoints.Append(manifests...)
	}

	result := &specs.Collection{
		FlowsManifest:     flows,
		EndpointsManifest: endpoints,
	}

	return result, nil
}
