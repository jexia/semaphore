package specs

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/jexia/maestro/instance"
	"github.com/jexia/maestro/schema"
	"github.com/jexia/maestro/specs/labels"
	"github.com/jexia/maestro/specs/types"
)

// Resolver when called collects the available manifest(s) with the configured configuration
type Resolver func(instance.Context, CustomDefinedFunctions) (*Manifest, error)

// FlowManager represents a flow manager
type FlowManager interface {
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

// Manifest holds a collection of definitions and resources
type Manifest struct {
	Flows     Flows     `json:"flows"`
	Proxy     Proxies   `json:"proxies"`
	Endpoints Endpoints `json:"endpoints"`
}

// GetFlow attempts to find a flow or proxy matching the given name
func (manifest *Manifest) GetFlow(name string) FlowManager {
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
func (manifest *Manifest) Merge(incoming *Manifest) {
	manifest.Flows = append(manifest.Flows, incoming.Flows...)
	manifest.Proxy = append(manifest.Proxy, incoming.Proxy...)
	manifest.Endpoints = append(manifest.Endpoints, incoming.Endpoints...)
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

// Options represents a collection of options
type Options map[string]string

// Header represents a collection of key values
type Header map[string]*Property

// CustomDefinedFunctions represents a collection of custom defined functions that could be called inside a template
type CustomDefinedFunctions map[string]PrepareFunction

// PrepareFunction prepares the custom defined function.
// The given arguments represent the exprected types that are passed when called.
// Properties returned should be absolute.
type PrepareFunction func(args ...*Property) (*Property, FunctionExec, error)

// FunctionExec executes the function and passes the expected types as stores
// A store should be returned which could be used to encode the function property
type FunctionExec func(store Store) error

// Functions represents a collection of functions
type Functions map[string]*Function

// Function represents a custom defined function
type Function struct {
	Arguments []*Property
	Fn        FunctionExec
	Returns   *Property
}

// PropertyReference represents a mustach template reference
type PropertyReference struct {
	Resource string
	Path     string
	Property *Property
}

func (reference *PropertyReference) String() string {
	return reference.Resource + ReferenceDelimiter + reference.Path
}

// Clone returns a clone of the given property reference
func (reference *PropertyReference) Clone() *PropertyReference {
	return &PropertyReference{
		Resource: reference.Resource,
		Path:     reference.Path,
		Property: reference.Property,
	}
}

// Property represents a value property.
// A value property could contain a constant value or a value reference.
type Property struct {
	Name      string
	Path      string
	Default   interface{}
	Type      types.Type
	Label     labels.Label
	Reference *PropertyReference
	Nested    map[string]*Property
	Expr      hcl.Expression // TODO: replace this with a custom solution
	Desciptor schema.Property
}

// ParameterMap is the initial map of parameter names (keys) and their (templated) values (values)
type ParameterMap struct {
	Schema    string    `json:"schema"`
	Options   Options   `json:"options"`
	Header    Header    `json:"header"`
	Property  *Property `json:"property"`
	Functions Functions
}

// Node represents a point inside a given flow where a request or rollback could be preformed.
// Nodes could be executed synchronously or asynchronously.
// All calls are referencing a service method, the service should match the alias defined inside the service.
// The request and response proto messages are used for type definitions.
// A call could contain the request headers, request body, rollback, and the execution type.
type Node struct {
	Name       string           `json:"name"`
	DependsOn  map[string]*Node `json:"depends_on"`
	Call       *Call            `json:"call"`
	Rollback   *Call            `json:"rollback"`
	Descriptor schema.Method    `json:"-"`
}

// GetName returns the call name
func (call *Node) GetName() string {
	return call.Name
}

// GetDescriptor returns the call descriptor
func (call *Node) GetDescriptor() schema.Method {
	return call.Descriptor
}

// Call represents a call which is executed during runtime
type Call struct {
	Service    string        `json:"service"`
	Method     string        `json:"method"`
	Request    *ParameterMap `json:"request"`
	Response   *ParameterMap `json:"response"`
	Descriptor schema.Method `json:"-"`
}

// GetRequest returns the call request parameter map
func (call *Call) GetRequest() *ParameterMap {
	return call.Request
}

// GetResponse returns the call response parameter map
func (call *Call) GetResponse() *ParameterMap {
	return call.Response
}

// GetService returns the call service
func (call *Call) GetService() string {
	return call.Service
}

// GetMethod returns the call endpoint
func (call *Call) GetMethod() string {
	return call.Method
}

// GetDescriptor returns the call descriptor
func (call *Call) GetDescriptor() schema.Method {
	return call.Descriptor
}

// SetResponse sets the given parameter map as response
func (call *Call) SetResponse(params *ParameterMap) {
	call.Response = params
}

// SetDescriptor sets the call descriptor
func (call *Call) SetDescriptor(descriptor schema.Method) {
	call.Descriptor = descriptor
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
