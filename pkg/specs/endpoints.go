package specs

import "github.com/jexia/semaphore/v2/pkg/specs/metadata"

// EndpointList represents a collection of endpoints
type EndpointList []*Endpoint

// Append merges the incoming manifest to the existing (left) manifest
func (endpoints *EndpointList) Append(list EndpointList) {
	*endpoints = append(*endpoints, list...)
}

// Get returns all endpoints for the given flow
func (endpoints EndpointList) Get(flow string) []*Endpoint {
	result := make([]*Endpoint, 0)
	for _, endpoint := range endpoints {
		if endpoint.Flow == flow {
			result = append(result, endpoint)
		}
	}

	return result
}

// Endpoint exposes a flow. Endpoints are not parsed by Semaphore and have custom implementations in each caller.
// The name of the endpoint represents the flow which should be executed.
type Endpoint struct {
	*metadata.Meta
	Flow     string  `json:"flow,omitempty"`
	Listener string  `json:"listener,omitempty"`
	Options  Options `json:"options,omitempty"`
}
