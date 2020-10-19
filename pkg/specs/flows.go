package specs

import "github.com/jexia/semaphore/pkg/specs/metadata"

// FlowListInterface represents a collection of flow interfaces
type FlowListInterface []FlowInterface

// Append appends the given flow list to the collection
func (collection *FlowListInterface) Append(list FlowListInterface) {
	*collection = append(*collection, list...)
}

// Get attempts to find a flow matching the given name
func (collection FlowListInterface) Get(name string) FlowInterface {
	for _, flow := range collection {
		if flow.GetName() == name {
			return flow
		}
	}

	return nil
}

// FlowList represents a collection of flows
type FlowList []*Flow

// Get attempts to find a flow matching the given name
func (collection FlowList) Get(name string) *Flow {
	for _, flow := range collection {
		if flow.Name == name {
			return flow
		}
	}

	return nil
}

// FlowInterface represents a proxy or flow manager.
type FlowInterface interface {
	GetMeta() *metadata.Meta
	GetName() string
	GetNodes() NodeList
	SetNodes(NodeList)
	GetInput() *ParameterMap
	GetOutput() *ParameterMap
	GetOnError() *OnError
	GetForward() *Call
}

// Flow defines a set of calls that should be called chronologically and produces an output message.
// Calls could reference other resources when constructing messages.
// All references are strictly typed and fetched from the configured schemas.
//
// All flows should contain a unique name.
// Calls are nested inside of flows and contain two labels, a unique name within the flow and the service and method to be called.
// A dependency reference structure is generated within the flow which allows Semaphore to figure out which calls could be called parallel to improve performance.
type Flow struct {
	*metadata.Meta
	Name    string        `json:"name,omitempty"`
	Input   *ParameterMap `json:"input,omitempty"`
	Nodes   NodeList      `json:"nodes,omitempty"`
	Output  *ParameterMap `json:"output,omitempty"`
	OnError *OnError      `json:"on_error,omitempty"`
}

// GetMeta returns the metadata value object of the given flow
func (flow *Flow) GetMeta() *metadata.Meta {
	return flow.Meta
}

// GetName returns the flow name
func (flow *Flow) GetName() string {
	return flow.Name
}

// GetNodes returns the calls of the given flow
func (flow *Flow) GetNodes() NodeList {
	return flow.Nodes
}

// SetNodes sets the given node list
func (flow *Flow) SetNodes(nodes NodeList) {
	flow.Nodes = nodes
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

// ProxyList represents a collection of proxies
type ProxyList []*Proxy

// Get attempts to find a proxy matching the given name
func (collection ProxyList) Get(name string) *Proxy {
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
	*metadata.Meta
	Input   *ParameterMap `json:"input,omitempty"`
	Name    string        `json:"name,omitempty"`
	Nodes   NodeList      `json:"nodes,omitempty"`
	Forward *Call         `json:"forward,omitempty"`
	OnError *OnError      `json:"on_error,omitempty"`
}

// GetMeta returns the metadata value object of the given proxy
func (proxy *Proxy) GetMeta() *metadata.Meta {
	return proxy.Meta
}

// GetName returns the flow name
func (proxy *Proxy) GetName() string {
	return proxy.Name
}

// GetNodes returns the calls of the given flow
func (proxy *Proxy) GetNodes() NodeList {
	return proxy.Nodes
}

// SetNodes sets the given node list
func (proxy *Proxy) SetNodes(nodes NodeList) {
	proxy.Nodes = nodes
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

// Evaluable can be evaluated (runs the expression) using provided parameters
type Evaluable interface {
	Evaluate(parameters map[string]interface{}) (interface{}, error)
}

// Condition represents a condition which could be true or false
type Condition struct {
	*metadata.Meta
	RawExpression string `json:"raw_expression,omitempty"`
	Evaluable     `json:"-"`
	Params        *ParameterMap `json:"-"`
}

// GetParameters returns the list of parameters
func (c Condition) GetParameters() *ParameterMap { return c.Params }

// NodeList represents a collection of nodes
type NodeList []*Node

// Get returns a node with the given ID
func (nodes NodeList) Get(name string) *Node {
	for _, node := range nodes {
		if node.ID == name {
			return node
		}
	}

	return nil
}

// NodeType represents the type of the given node.
// A type determines the purpose of the node but not the implementation.
type NodeType string

var (
	// NodeCall is assigned when the given node is used to call a external service
	NodeCall NodeType = "call"
	// NodeCondition is assigned when the given node is used to execute a
	// conditional expression.
	NodeCondition NodeType = "condition"
	// NodeIntermediate is assigned when the node is used as a coming in
	// between nodes.
	NodeIntermediate NodeType = "intermediate"
)

// Dependencies represents a collection of node dependencies
type Dependencies map[string]*Node

// Append appends the given input to the already existing dependencies and returns the result
func (dependencies Dependencies) Append(input Dependencies) Dependencies {
	result := make(Dependencies, len(dependencies)+len(input))

	for key, val := range dependencies {
		result[key] = val
	}

	for key, val := range input {
		result[key] = val
	}

	return result
}

// Node represents a point inside a given flow where a request or rollback could be preformed.
// Nodes could be executed synchronously or asynchronously.
// All calls are referencing a service method, the service should match the alias defined inside the service.
// The request and response proto messages are used for type definitions.
// A call could contain the request headers, request body, rollback, and the execution type.
type Node struct {
	*metadata.Meta
	Type         NodeType      `json:"type,omitempty"`
	ID           string        `json:"id,omitempty"`
	Name         string        `json:"name,omitempty"`
	Intermediate *ParameterMap `json:"intermediate,omitempty"`
	Condition    *Condition    `json:"condition,omitempty"`
	DependsOn    Dependencies  `json:"depends_on,omitempty"`
	Call         *Call         `json:"call,omitempty"`
	Rollback     *Call         `json:"rollback,omitempty"`
	ExpectStatus []int         `json:"expect_status,omitempty"`
	OnError      *OnError      `json:"on_error,omitempty"`
}

// GetOnError returns the error handling for the given node
func (node *Node) GetOnError() *OnError {
	return node.OnError
}

// Call represents a call which is executed during runtime
type Call struct {
	*metadata.Meta
	Service    string        `json:"service,omitempty"`
	Method     string        `json:"method,omitempty"`
	Request    *ParameterMap `json:"request,omitempty"`
	Response   *ParameterMap `json:"response,omitempty"`
	Descriptor *Method       `json:"-"`
}

// ErrorHandle represents a error handle object
type ErrorHandle interface {
	GetResponse() *ParameterMap
	GetStatusCode() *Property
	GetMessage() *Property
}

// OnError represents the variables that have to be returned if a unexpected error is returned
type OnError struct {
	*metadata.Meta
	Response *ParameterMap `json:"response,omitempty"`
	// Question: does it make sense to use full property here or just a Template is enough?
	Status  *Property            `json:"status,omitempty"`
	Message *Property            `json:"message,omitempty"`
	Params  map[string]*Property `json:"params,omitempty"`
}

// Clone clones the given error
func (err *OnError) Clone() *OnError {
	if err == nil {
		return nil
	}

	result := &OnError{
		Meta:     err.Meta,
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
