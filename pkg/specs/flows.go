package specs

// FlowsManifest holds a collection of definitions and resources
type FlowsManifest struct {
	Error *Error  `json:"error"`
	Flows Flows   `json:"flows"`
	Proxy Proxies `json:"proxies"`
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

// MergeFlowsManifest merges the incoming manifest to the existing (left) manifest
func MergeFlowsManifest(left *FlowsManifest, incoming ...*FlowsManifest) {
	for _, manifest := range incoming {
		left.Flows = append(left.Flows, manifest.Flows...)
		left.Proxy = append(left.Proxy, manifest.Proxy...)
	}
}

// FlowResourceManager represents a proxy or flow manager.
type FlowResourceManager interface {
	GetName() string
	GetNodes() []*Node
	GetInput() *ParameterMap
	GetOutput() *ParameterMap
	GetError() *Error
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
	Error  *Error        `json:"error"`
}

// GetName returns the flow name
func (flow *Flow) GetName() string {
	return flow.Name
}

// GetNodes returns the calls of the given flow
func (flow *Flow) GetNodes() []*Node {
	return flow.Nodes
}

// GetError returns the error of the given flow
func (flow *Flow) GetError() *Error {
	return flow.Error
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
	Input   *ParameterMap `json:"input"`
	Name    string        `json:"name"`
	Nodes   []*Node       `json:"nodes"`
	Forward *Call         `json:"forward"`
	Error   *Error        `json:"error"`
}

// GetName returns the flow name
func (proxy *Proxy) GetName() string {
	return proxy.Name
}

// GetNodes returns the calls of the given flow
func (proxy *Proxy) GetNodes() []*Node {
	return proxy.Nodes
}

// GetError returns the error of the given flow
func (proxy *Proxy) GetError() *Error {
	return proxy.Error
}

// GetInput returns the input of the given flow
func (proxy *Proxy) GetInput() *ParameterMap {
	return proxy.Input
}

// GetOutput returns the output of the given flow
func (proxy *Proxy) GetOutput() *ParameterMap {
	return nil
}

// GetForward returns the proxy forward of the given flow
func (proxy *Proxy) GetForward() *Call {
	return proxy.Forward
}

// Error represents a error message returned to the user once a unexpected error is returned
type Error struct {
	Schema   string    `json:"schema"`
	Header   Header    `json:"header"`
	Property *Property `json:"property"`
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
	OnError   *OnError         `json:"on_error"`
	Error     *Error           `json:"error"`
}

// Call represents a call which is executed during runtime
type Call struct {
	Service    string        `json:"service"`
	Method     string        `json:"method"`
	Request    *ParameterMap `json:"request"`
	Response   *ParameterMap `json:"response"`
	Descriptor *Method       `json:"-"`
}

// OnError represents the variables that have to be returned if a unexpected error is returned
type OnError struct {
	Schema  string
	Status  string
	Message string
	Params  map[string]*PropertyReference
}
