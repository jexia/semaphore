package graphql

import (
	"github.com/jexia/semaphore/pkg/transport"
)

const (
	// PathOption represents the object name option key
	PathOption = "path"
	// BaseOption represents the object base option key
	BaseOption = "base"
	// NameOption represents the object name option key
	NameOption = "name"
)

// EndpointOptions represents the available HTTP options
type EndpointOptions struct {
	Name string
	Path string
	Base string
}

// ParseEndpointOptions parses the given specs options into HTTP options
func ParseEndpointOptions(endpoint *transport.Endpoint) (EndpointOptions, error) {
	if endpoint == nil || endpoint.Flow == nil {
		return EndpointOptions{}, nil
	}

	result := EndpointOptions{
		Name: endpoint.Flow.GetName(),
		Path: endpoint.Flow.GetName(),
		Base: "query",
	}

	path, has := endpoint.Options[PathOption]
	if has {
		result.Path = path
	}

	base, has := endpoint.Options[BaseOption]
	if has {
		if base != QueryObject && base != MutationObject {
			return EndpointOptions{}, ErrUnknownBase{
				Base: base,
			}
		}

		result.Base = base
	}

	name, has := endpoint.Options[NameOption]
	if has {
		result.Name = name
	}

	return result, nil
}
