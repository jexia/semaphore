package specs

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/jexia/maestro/schema"
	"github.com/jexia/maestro/specs/types"
)

// FlowManager represents a flow manager
type FlowManager interface {
	GetName() string
	GetDependencies() map[string]*Flow
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
	Flows     Flows
	Proxy     Proxies
	Endpoints Endpoints
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

// MergeLeft merges the incoming manifest to the existing (left) manifest
func (manifest *Manifest) MergeLeft(incoming *Manifest) {
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
	Name      string
	DependsOn map[string]*Flow
	Input     *ParameterMap
	Nodes     []*Node
	Output    *ParameterMap
}

// GetName returns the flow name
func (flow *Flow) GetName() string {
	return flow.Name
}

// GetDependencies returns the dependencies of the given flow
func (flow *Flow) GetDependencies() map[string]*Flow {
	return flow.DependsOn
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
	Flow     string
	Listener string
	Codec    string
	Options  Options
}

// Options represents a collection of options
type Options map[string]string

// Header represents a collection of key values
type Header map[string]*Property

// CustomDefinedFunctions represents a collection of custom defined functions that could be called inside a template
type CustomDefinedFunctions map[string]PrepareCustomFunction

// PrepareCustomFunction prepares the custom defined function.
// The given arguments represent the exprected types that are passed when called.
type PrepareCustomFunction func(path string, args ...*Property) (*Property, error)

// HandleCustomFunction executes the function and passes the expected types as interface{}.
// The expected property type should always be returned.
type HandleCustomFunction func(args ...interface{}) interface{}

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
	Label     types.Label
	Reference *PropertyReference
	Nested    map[string]*Property
	Expr      hcl.Expression // TODO: marked for removal
	Function  HandleCustomFunction
	Desciptor schema.Property
}

// Clone returns a hollow clone of the property
func (property *Property) Clone() *Property {
	return &Property{
		Name:      property.Name,
		Path:      property.Path,
		Default:   property.Default,
		Type:      property.Type,
		Label:     property.Label,
		Reference: property.Reference,
		Expr:      property.Expr,
		Function:  property.Function,
		Desciptor: property.Desciptor,
	}
}

// ParameterMap is the initial map of parameter names (keys) and their (templated) values (values)
type ParameterMap struct {
	Schema   string
	Options  Options
	Header   Header
	Property *Property
}

// Node represents a point inside a given flow where a request or rollback could be preformed.
// Nodes could be executed synchronously or asynchronously.
// All calls are referencing a service method, the service should match the alias defined inside the service.
// The request and response proto messages are used for type definitions.
// A call could contain the request headers, request body, rollback, and the execution type.
type Node struct {
	Name       string
	DependsOn  map[string]*Node
	Type       string
	Call       *Call
	Rollback   *Call
	Descriptor schema.Method
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
	Service    string
	Method     string
	Request    *ParameterMap
	Response   *ParameterMap
	Descriptor schema.Method
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

// SetDescriptor sets the call descriptor
func (call *Call) SetDescriptor(descriptor schema.Method) {
	if descriptor != nil {
		call.Response = ToParameterMap(nil, "", descriptor.GetOutput())
	}

	call.Descriptor = descriptor
}

// Proxy streams the incoming request to the given service.
// Proxies could define calls that are executed before the request body is forwarded.
// A proxy forward could ideally be used for file uploads or large messages which could not be stored in memory.
type Proxy struct {
	Name      string
	DependsOn map[string]*Flow
	Nodes     []*Node
	Forward   *Call
}

// GetName returns the flow name
func (proxy *Proxy) GetName() string {
	return proxy.Name
}

// GetDependencies returns the dependencies of the given flow
func (proxy *Proxy) GetDependencies() map[string]*Flow {
	return proxy.DependsOn
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
