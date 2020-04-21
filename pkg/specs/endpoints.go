package specs

// EndpointsManifest holds a collection of flow endpoints
type EndpointsManifest struct {
	Endpoints Endpoints `json:"endpoints"`
}

// MergeEndpointsManifest merges the incoming manifest to the existing (left) manifest
func MergeEndpointsManifest(left *EndpointsManifest, incoming ...*EndpointsManifest) {
	for _, manifest := range incoming {
		left.Endpoints = append(left.Endpoints, manifest.Endpoints...)
	}
}

// Endpoints represents a collection of endpoints
type Endpoints []*Endpoint

// Get attempts to find a endpoint for the given flow
func (collection Endpoints) Get(flow string) []*Endpoint {
	result := make([]*Endpoint, 0)
	for _, endpoint := range collection {
		if endpoint.Flow == flow {
			result = append(result, endpoint)
		}
	}

	return result
}

// Endpoint exposes a flow. Endpoints are not parsed by Maestro and have custom implementations in each caller.
// The name of the endpoint represents the flow which should be executed.
type Endpoint struct {
	Flow     string  `json:"flow"`
	Listener string  `json:"listener"`
	Options  Options `json:"options"`
}
