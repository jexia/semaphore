package specs

// EndpointsManifest holds a collection of flow endpoints
type EndpointsManifest struct {
	Endpoints EndpointList `json:"endpoints,omitempty"`
}

// Append merges the incoming manifest to the existing (left) manifest
func (manifest *EndpointsManifest) Append(incoming ...*EndpointsManifest) {
	if manifest == nil {
		return
	}

	for _, right := range incoming {
		manifest.Endpoints = append(manifest.Endpoints, right.Endpoints...)
	}
}

// EndpointList represents a collection of endpoints
type EndpointList []*Endpoint

// Get attempts to find a endpoint for the given flow
func (collection EndpointList) Get(flow string) []*Endpoint {
	result := make([]*Endpoint, 0)
	for _, endpoint := range collection {
		if endpoint.Flow == flow {
			result = append(result, endpoint)
		}
	}

	return result
}

// Endpoint exposes a flow. Endpoints are not parsed by Semaphore and have custom implementations in each caller.
// The name of the endpoint represents the flow which should be executed.
type Endpoint struct {
	Flow     string  `json:"flow,omitempty"`
	Listener string  `json:"listener,omitempty"`
	Options  Options `json:"options,omitempty"`
}
