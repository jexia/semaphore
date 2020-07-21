package specs

import (
	"github.com/Knetic/govaluate"
)

// FlowsManifest holds a collection of definitions and resources
type FlowsManifest struct {
	Error *ParameterMap `json:"error,omitempty"`
	Flows Flows         `json:"flows,omitempty"`
	Proxy Proxies       `json:"proxies,omitempty"`
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
// A dependency reference structure is generated within the flow which allows Semaphore to figure out which calls could be called parallel to improve performance.
type Flow struct {
	Name    string        `json:"name,omitempty"`
	Input   *ParameterMap `json:"input,omitempty"`
	Nodes   []*Node       `json:"nodes,omitempty"`
	Output  *ParameterMap `json:"output,omitempty"`
	OnError *OnError      `json:"on_error,omitempty"`
}

// GetName returns the flow name
func (flow *Flow) GetName() string {
	return flow.Name
}

// GetNodes returns the calls of the given flow
func (flow *Flow) GetNodes() []*Node {
	return flow.Nodes
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
	Input   *ParameterMap `json:"input,omitempty"`
	Name    string        `json:"name,omitempty"`
	Nodes   []*Node       `json:"nodes,omitempty"`
	Forward *Call         `json:"forward,omitempty"`
	OnError *OnError      `json:"on_error,omitempty"`
}

// GetName returns the flow name
func (proxy *Proxy) GetName() string {
	return proxy.Name
}

// GetNodes returns the calls of the given flow
func (proxy *Proxy) GetNodes() []*Node {
	return proxy.Nodes
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
	RawExpression string                         `json:"raw_expression,omitempty"`
	Expression    *govaluate.EvaluableExpression `json:"-"`
	Params        *ParameterMap                  `json:"-"`
}

// Node represents a point inside a given flow where a request or rollback could be preformed.
// Nodes could be executed synchronously or asynchronously.
// All calls are referencing a service method, the service should match the alias defined inside the service.
// The request and response proto messages are used for type definitions.
// A call could contain the request headers, request body, rollback, and the execution type.
type Node struct {
	ID           string           `json:"id,omitempty"`
	Name         string           `json:"name,omitempty"`
	Condition    *Condition       `json:"condition,omitempty"`
	DependsOn    map[string]*Node `json:"depends_on,omitempty"`
	Call         *Call            `json:"call,omitempty"`
	Rollback     *Call            `json:"rollback,omitempty"`
	ExpectStatus []int            `json:"expect_status,omitempty"`
	OnError      *OnError         `json:"on_error,omitempty"`
}

// GetOnError returns the error handling for the given node
func (node *Node) GetOnError() *OnError {
	return node.OnError
}

// Call represents a call which is executed during runtime
type Call struct {
	Service    string        `json:"service,omitempty"`
	Method     string        `json:"method,omitempty"`
	Request    *ParameterMap `json:"request,omitempty"`
	Response   *ParameterMap `json:"response,omitempty"`
	Descriptor *Method       `json:"-"`
}

// OnError represents the variables that have to be returned if a unexpected error is returned
type OnError struct {
	Response *ParameterMap        `json:"response,omitempty"`
	Status   *Property            `json:"status,omitempty"`
	Message  *Property            `json:"message,omitempty"`
	Params   map[string]*Property `json:"params,omitempty"`
}

// Clone clones the given error
func (err *OnError) Clone() *OnError {
	if err == nil {
		return nil
	}

	result := &OnError{
		Response: err.Response.Clone(),
		Status:   err.Status.Clone(),
		Message:  err.Message.Clone(),
		Params:   make(map[string]*Property, len(err.Params)),
	}

	for key, param := range err.Params {
		result.Params[key] = param.Clone()
	}

	return result
}

// GetResponse returns the error response
func (err *OnError) GetResponse() *ParameterMap {
	if err == nil {
		return nil
	}

	return err.Response
}

// GetStatusCode returns the status code property
func (err *OnError) GetStatusCode() *Property {
	if err == nil {
		return nil
	}

	return err.Status
}

// GetMessage returns the message property
func (err *OnError) GetMessage() *Property {
	if err == nil {
		return nil
	}

	return err.Message
}
