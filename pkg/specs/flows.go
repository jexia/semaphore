package specs

import "github.com/Knetic/govaluate"

// FlowsManifest holds a collection of definitions and resources
type FlowsManifest struct {
	Error *ParameterMap `json:"error"`
	Flows Flows         `json:"flows"`
	Proxy Proxies       `json:"proxies"`
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
		if manifest.Error != nil {
			left.Error = manifest.Error
		}

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
	GetError() *ParameterMap
	GetOnError() *OnError
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
	Name    string        `json:"name"`
	Input   *ParameterMap `json:"input"`
	Nodes   []*Node       `json:"nodes"`
	Output  *ParameterMap `json:"output"`
	Error   *ParameterMap `json:"error"`
	OnError *OnError      `json:"on_error"`
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
func (flow *Flow) GetError() *ParameterMap {
	return flow.Error
}

// GetOnError returns the error handling of the given flow
func (flow *Flow) GetOnError() *OnError {
	return flow.OnError
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
	Error   *ParameterMap `json:"error"`
	OnError *OnError      `json:"on_error"`
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
func (proxy *Proxy) GetError() *ParameterMap {
	return proxy.Error
}

// GetOnError returns the error handling of the given flow
func (proxy *Proxy) GetOnError() *OnError {
	return proxy.OnError
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

// Condition represents a condition which could be true or false
type Condition struct {
	RawExpression string                         `json:"raw_expression"`
	Expression    *govaluate.EvaluableExpression `json:"-"`
	Params        *ParameterMap                  `json:"-"`
}

// Node represents a point inside a given flow where a request or rollback could be preformed.
// Nodes could be executed synchronously or asynchronously.
// All calls are referencing a service method, the service should match the alias defined inside the service.
// The request and response proto messages are used for type definitions.
// A call could contain the request headers, request body, rollback, and the execution type.
type Node struct {
	Name      string           `json:"name"`
	Condition *Condition       `json:"condition"`
	DependsOn map[string]*Node `json:"depends_on"`
	Call      *Call            `json:"call"`
	Rollback  *Call            `json:"rollback"`
	Error     *ParameterMap    `json:"error"`
	OnError   *OnError         `json:"on_error"`
}

// GetError returns the error object for the given node
func (node *Node) GetError() *ParameterMap {
	return node.Error
}

// GetOnError returns the error handling for the given node
func (node *Node) GetOnError() *OnError {
	return node.OnError
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
	Status  *Property
	Message *Property
	Params  map[string]*Property
}
