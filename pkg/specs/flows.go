package specs

// NewFlowsManifest constructs a new empty flows manifest
func NewFlowsManifest() *FlowsManifest {
	return &FlowsManifest{
		Flows:     make(Flows, 0),
		Proxy:     make(Proxies, 0),
		Endpoints: make(Endpoints, 0),
	}
}

// FlowsManifest holds a collection of definitions and resources
type FlowsManifest struct {
	Flows     Flows     `json:"flows"`
	Proxy     Proxies   `json:"proxies"`
	Endpoints Endpoints `json:"endpoints"`
}

// GetFlow attempts to find a flow or proxy matching the given name
func (manifest *FlowsManifest) GetFlow(name string) FlowResourceManager {
	flow := manifest.Flows.Get(name)
	if flow != nil {
		return flow
	}

	proxy := manifest.Proxy.Get(name)
	if proxy != nil {
		return proxy
	}

	return nil
}

// Merge merges the incoming manifest to the existing (left) manifest
func (manifest *FlowsManifest) Merge(incoming *FlowsManifest) {
	manifest.Flows = append(manifest.Flows, incoming.Flows...)
	manifest.Proxy = append(manifest.Proxy, incoming.Proxy...)
	manifest.Endpoints = append(manifest.Endpoints, incoming.Endpoints...)
}

// FlowResourceManager represents a proxy or flow manager.
type FlowResourceManager interface {
	GetName() string
	GetNodes() []*Node
	GetInput() *ParameterMap
	GetOutput() *ParameterMap
	GetForward() *Call
}

// Flows represents a collection of flows
type Flows []*Flow

// Get attempts to find a flow matching the given name
func (collection Flows) Get(name string) *Flow {
	for _, flow := range collection {
		if flow.Name == name {
			return flow
		}
	}

	return nil
}

// Flow defines a set of calls that should be called chronologically and produces an output message.
// Calls could reference other resources when constructing messages.
// All references are strictly typed and fetched from the configured schemas.
//
// All flows should contain a unique name.
// Calls are nested inside of flows and contain two labels, a unique name within the flow and the service and method to be called.
// A dependency reference structure is generated within the flow which allows Maestro to figure out which calls could be called parallel to improve performance.
type Flow struct {
	Name   string        `json:"name"`
	Input  *ParameterMap `json:"input"`
	Nodes  []*Node       `json:"nodes"`
	Output *ParameterMap `json:"output"`
}

// GetName returns the flow name
func (flow *Flow) GetName() string {
	return flow.Name
}

// GetNodes returns the calls of the given flow
func (flow *Flow) GetNodes() []*Node {
	return flow.Nodes
}

// GetInput returns the input of the given flow
func (flow *Flow) GetInput() *ParameterMap {
	return flow.Input
}

// GetOutput returns the output of the given flow
func (flow *Flow) GetOutput() *ParameterMap {
	return flow.Output
}

// GetForward returns the proxy forward of the given flow
func (flow *Flow) GetForward() *Call {
	return nil
}

// Endpoint exposes a flow. Endpoints are not parsed by Maestro and have custom implementations in each caller.
// The name of the endpoint represents the flow which should be executed.
type Endpoint struct {
	Flow     string  `json:"flow"`
	Listener string  `json:"listener"`
	Options  Options `json:"options"`
}

// Proxies represents a collection of proxies
type Proxies []*Proxy

// Get attempts to find a proxy matching the given name
func (collection Proxies) Get(name string) *Proxy {
	for _, proxy := range collection {
		if proxy.Name == name {
			return proxy
		}
	}

	return nil
}

// Proxy streams the incoming request to the given service.
// Proxies could define calls that are executed before the request body is forwarded.
// A proxy forward could ideally be used for file uploads or large messages which could not be stored in memory.
type Proxy struct {
	Name    string  `json:"name"`
	Nodes   []*Node `json:"nodes"`
	Forward *Call   `json:"forward"`
}

// GetName returns the flow name
func (proxy *Proxy) GetName() string {
	return proxy.Name
}

// GetNodes returns the calls of the given flow
func (proxy *Proxy) GetNodes() []*Node {
	return proxy.Nodes
}

// GetInput returns the input of the given flow
func (proxy *Proxy) GetInput() *ParameterMap {
	return nil
}

// GetOutput returns the output of the given flow
func (proxy *Proxy) GetOutput() *ParameterMap {
	return nil
}

// GetForward returns the proxy forward of the given flow
func (proxy *Proxy) GetForward() *Call {
	return proxy.Forward
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

// Node represents a point inside a given flow where a request or rollback could be preformed.
// Nodes could be executed synchronously or asynchronously.
// All calls are referencing a service method, the service should match the alias defined inside the service.
// The request and response proto messages are used for type definitions.
// A call could contain the request headers, request body, rollback, and the execution type.
type Node struct {
	Name      string           `json:"name"`
	DependsOn map[string]*Node `json:"depends_on"`
	Call      *Call            `json:"call"`
	Rollback  *Call            `json:"rollback"`
}

// Call represents a call which is executed during runtime
type Call struct {
	Service    string        `json:"service"`
	Method     string        `json:"method"`
	Request    *ParameterMap `json:"request"`
	Response   *ParameterMap `json:"response"`
	Descriptor *Method       `json:"-"` // TODO: check the usage of the descriptor
}
