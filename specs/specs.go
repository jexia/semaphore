package specs

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/jexia/maestro/schema"
	"github.com/jexia/maestro/specs/types"
)

// Object represents a parameter collection
type Object interface {
	GetProperties() map[string]*Property
	GetNestedProperties() map[string]*NestedParameterMap
	GetRepeatedProperties() map[string]*RepeatedParameterMap
	GetLabel() types.Label
	GetHeader() Header
}

// FlowCaller represents a flow caller
type FlowCaller interface {
	GetName() string
	GetEndpoint() string
	GetRequest() Object
	GetResponse() Object
	SetDescriptor(schema.Method)
}

// FlowManager represents a flow manager
type FlowManager interface {
	GetName() string
	GetDependencies() map[string]*Flow
	GetCalls() []*Call
	GetInput() *ParameterMap
	GetOutput() *ParameterMap
}

// Manifest holds a collection of definitions and resources
type Manifest struct {
	Flows     []*Flow
	Proxy     []*Proxy
	Endpoints []*Endpoint
	Services  []*Service
	Callers   []*Caller
}

// MergeLeft merges the incoming manifest to the existing (left) manifest
func (manifest *Manifest) MergeLeft(incoming *Manifest) {
	manifest.Flows = append(manifest.Flows, incoming.Flows...)
	manifest.Proxy = append(manifest.Proxy, incoming.Proxy...)
	manifest.Endpoints = append(manifest.Endpoints, incoming.Endpoints...)
	manifest.Services = append(manifest.Services, incoming.Services...)
	manifest.Callers = append(manifest.Callers, incoming.Callers...)
}

// Flow defines a set of calls that should be called chronologically and produces an output message.
// Calls could reference other resources when constructing messages.
// All references are strictly typed. Properties are fetched from the given proto buffers or inputs.
//
// All flows should contain a unique name.
// Calls are nested inside of flows and contain two labels, a unique name within the flow and the service and method to be called.
// A dependency reference structure is generated within the flow which allows Maestro to figure out which calls could be called parallel to improve performance.
type Flow struct {
	Name       string
	DependsOn  map[string]*Flow
	Schema     string
	Input      *ParameterMap
	Calls      []*Call
	Output     *ParameterMap
	Descriptor schema.Method
}

// GetName returns the flow name
func (flow *Flow) GetName() string {
	return flow.Name
}

// GetDependencies returns the dependencies of the given flow
func (flow *Flow) GetDependencies() map[string]*Flow {
	return flow.DependsOn
}

// GetCalls returns the calls of the given flow
func (flow *Flow) GetCalls() []*Call {
	return flow.Calls
}

// GetInput returns the input of the given flow
func (flow *Flow) GetInput() *ParameterMap {
	return flow.Input
}

// GetOutput returns the output of the given flow
func (flow *Flow) GetOutput() *ParameterMap {
	return flow.Output
}

// Endpoint exposes a flow. Endpoints are not parsed by Maestro and have custom implementations in each caller.
// The name of the endpoint represents the flow which should be executed.
type Endpoint struct {
	Flow string
	Body map[string]interface{}
}

// Options represents a collection of options
type Options map[string]interface{}

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
	Label    types.Label
	Object   Object
}

func (reference *PropertyReference) String() string {
	return reference.Resource + ReferenceDelimiter + reference.Path
}

// Property represents a value property.
// A value property could contain a constant value or a value reference.
type Property struct {
	Path      string
	Name      string
	Default   interface{}
	Type      types.Type
	Reference *PropertyReference
	Expr      hcl.Expression
	Function  HandleCustomFunction
}

// GetPath returns the property path
func (property *Property) GetPath() string {
	return property.Path
}

// GetDefault returns the property default type
func (property *Property) GetDefault() interface{} {
	return property.Default
}

// GetType returns the property type
func (property *Property) GetType() types.Type {
	return property.Type
}

// GetObject returns the property object
func (property *Property) GetObject() Object {
	return nil
}

// Clone returns a clone of the property
func (property *Property) Clone() *Property {
	return &Property{
		Path:      property.Path,
		Default:   property.Default,
		Type:      property.Type,
		Reference: property.Reference,
		Expr:      property.Expr,
	}
}

// ParameterMap is the initial map of parameter names (keys) and their (templated) values (values)
type ParameterMap struct {
	Options    Options
	Header     Header
	Nested     map[string]*NestedParameterMap
	Repeated   map[string]*RepeatedParameterMap
	Properties map[string]*Property
}

// GetProperties returns the properties inside the given parameter map
func (parameters *ParameterMap) GetProperties() map[string]*Property {
	return parameters.Properties
}

// GetNestedProperties returns the nested parameter map inside the given parameter map
func (parameters *ParameterMap) GetNestedProperties() map[string]*NestedParameterMap {
	return parameters.Nested
}

// GetRepeatedProperties returns the repeated parameter map inside the given parameter map
func (parameters *ParameterMap) GetRepeatedProperties() map[string]*RepeatedParameterMap {
	return parameters.Repeated
}

// GetHeader returns the parameter map header
func (parameters *ParameterMap) GetHeader() Header {
	return parameters.Header
}

// GetLabel returns the parameter map label
func (parameters *ParameterMap) GetLabel() types.Label {
	return types.LabelOptional
}

// NestedParameterMap is a map of parameter names (keys) and their (templated) values (values)
type NestedParameterMap struct {
	Path       string
	Name       string
	Nested     map[string]*NestedParameterMap
	Repeated   map[string]*RepeatedParameterMap
	Properties map[string]*Property
}

// GetPath returns the parameter map path
func (nested *NestedParameterMap) GetPath() string {
	return nested.Path
}

// GetProperties returns the properties inside the given parameter map
func (nested *NestedParameterMap) GetProperties() map[string]*Property {
	return nested.Properties
}

// GetNestedProperties returns the nested parameter map inside the given parameter map
func (nested *NestedParameterMap) GetNestedProperties() map[string]*NestedParameterMap {
	return nested.Nested
}

// GetRepeatedProperties returns the repeated parameter map inside the given parameter map
func (nested *NestedParameterMap) GetRepeatedProperties() map[string]*RepeatedParameterMap {
	return nested.Repeated
}

// GetDefault returns nil
func (nested *NestedParameterMap) GetDefault() interface{} {
	return nil
}

// GetType returns the nested parameter map type
func (nested *NestedParameterMap) GetType() types.Type {
	return types.TypeMessage
}

// GetObject returns the nested parameter map type
func (nested *NestedParameterMap) GetObject() Object {
	return nested
}

// GetHeader returns the nested parameter map header
func (nested *NestedParameterMap) GetHeader() Header {
	return nil
}

// GetLabel returns the nested parameter map label
func (nested *NestedParameterMap) GetLabel() types.Label {
	return types.LabelOptional
}

// Clone returns a clone of the nested parameter map
func (nested *NestedParameterMap) Clone(name string, path string) *NestedParameterMap {
	returns := &NestedParameterMap{
		Name:       name,
		Path:       path,
		Properties: make(map[string]*Property, len(nested.Properties)),
		Nested:     make(map[string]*NestedParameterMap, len(nested.Nested)),
		Repeated:   make(map[string]*RepeatedParameterMap, len(nested.Repeated)),
	}

	for name, property := range nested.Properties {
		returns.Properties[name] = property.Clone()
	}

	for name, nested := range nested.Nested {
		returns.Nested[name] = nested.Clone(name, JoinPath(returns.Path, name))
	}

	for name, repeated := range nested.Repeated {
		returns.Repeated[name] = repeated.Clone(name, JoinPath(returns.Path, name))
	}

	return returns
}

// RepeatedParameterMap is a map of repeated message blocks/values
type RepeatedParameterMap struct {
	Path       string
	Name       string
	Template   *PropertyReference
	Nested     map[string]*NestedParameterMap
	Repeated   map[string]*RepeatedParameterMap
	Properties map[string]*Property
}

// GetPath returns the repeated path
func (repeated *RepeatedParameterMap) GetPath() string {
	return repeated.Path
}

// GetProperties returns the properties inside the given parameter map
func (repeated *RepeatedParameterMap) GetProperties() map[string]*Property {
	return repeated.Properties
}

// GetNestedProperties returns the nested parameter map inside the given parameter map
func (repeated *RepeatedParameterMap) GetNestedProperties() map[string]*NestedParameterMap {
	return repeated.Nested
}

// GetRepeatedProperties returns the repeated parameter map inside the given parameter map
func (repeated *RepeatedParameterMap) GetRepeatedProperties() map[string]*RepeatedParameterMap {
	return repeated.Repeated
}

// GetDefault returns nil
func (repeated *RepeatedParameterMap) GetDefault() interface{} {
	return nil
}

// GetType returns the parameter map type
func (repeated *RepeatedParameterMap) GetType() types.Type {
	return types.TypeMessage
}

// GetObject returns the repeated parameter map type
func (repeated *RepeatedParameterMap) GetObject() Object {
	return repeated
}

// GetHeader returns the repeated parameter map header
func (repeated *RepeatedParameterMap) GetHeader() Header {
	return nil
}

// GetLabel returns the repeated parameter map label
func (repeated *RepeatedParameterMap) GetLabel() types.Label {
	return types.LabelRepeated
}

// Clone returns a clone of the nested parameter map
func (repeated *RepeatedParameterMap) Clone(name string, path string) *RepeatedParameterMap {
	returns := &RepeatedParameterMap{
		Name:       name,
		Template:   repeated.Template,
		Path:       path,
		Properties: make(map[string]*Property, len(repeated.Properties)),
		Nested:     make(map[string]*NestedParameterMap, len(repeated.Nested)),
		Repeated:   make(map[string]*RepeatedParameterMap, len(repeated.Repeated)),
	}

	for name, property := range repeated.Properties {
		returns.Properties[name] = property.Clone()
	}

	for name, nested := range repeated.Nested {
		returns.Nested[name] = nested.Clone(name, JoinPath(returns.Path, name))
	}

	for name, repeated := range repeated.Repeated {
		returns.Repeated[name] = repeated.Clone(name, JoinPath(returns.Path, name))
	}

	return returns
}

// Call calls the given service and method.
// Calls could be executed synchronously or asynchronously.
// All calls are referencing a service method, the service should match the alias defined inside the service.
// The request and response proto messages are used for type definitions.
// A call could contain the request headers, request body, rollback, and the execution type.
type Call struct {
	Name       string
	DependsOn  map[string]*Call
	Endpoint   string
	Type       string
	Request    *ParameterMap
	Response   *ParameterMap
	Rollback   *RollbackCall
	Descriptor schema.Method
}

// GetName returns the call name
func (call *Call) GetName() string {
	return call.Name
}

// GetRequest returns the call request parameter map
func (call *Call) GetRequest() Object {
	return call.Request
}

// GetResponse returns the call response parameter map
func (call *Call) GetResponse() Object {
	return call.Response
}

// SetDescriptor sets the call method descriptor
func (call *Call) SetDescriptor(descriptor schema.Method) {
	if descriptor != nil {
		call.Response = ToParameterMap(nil, "", descriptor.GetOutput())
	}

	call.Descriptor = descriptor
}

// GetEndpoint returns the call endpoint
func (call *Call) GetEndpoint() string {
	return call.Endpoint
}

// RollbackCall represents the rollback call which is executed when a call inside a flow failed.
type RollbackCall struct {
	Parent     *Call
	Endpoint   string
	Request    *ParameterMap
	Descriptor schema.Method
}

// GetName returns the call name
func (call *RollbackCall) GetName() string {
	if call.Parent == nil {
		return ""
	}

	return call.Parent.Name
}

// GetRequest returns the call request parameter map
func (call *RollbackCall) GetRequest() Object {
	return call.Request
}

// GetResponse returns the call response parameter map
func (call *RollbackCall) GetResponse() Object {
	return nil
}

// SetDescriptor sets the call method descriptor
func (call *RollbackCall) SetDescriptor(descriptor schema.Method) {
	call.Descriptor = descriptor
}

// GetEndpoint returns the call endpoint
func (call *RollbackCall) GetEndpoint() string {
	return call.Endpoint
}

// Service represent external service which could be called inside the flows.
// The service name is an alias which could be referenced inside calls.
// The host of the service and proto service method should be defined for each service.
//
// The request and response message defined inside the proto buffers are used for type definitions.
// The FQN (fully qualified name) of the proto method should be used.
// Each service references a caller implementation to be used.
type Service struct {
	Options Options
	Alias   string
	Caller  string
	Host    string
	Codec   string
	Schema  string
}

// Caller Each implementation has to be configured and defined before running the service.
// All values are passed as attributes to the callers to be unmarshalled.
// These attributes could be used for configuration purposes
type Caller struct {
	Name string
	Body map[string]interface{}
}

// Proxy streams the incoming request to the given service.
// Proxies could define calls that are executed before the request body is forwarded.
// A proxy forward could ideally be used for file uploads or large messages which could not be stored in memory.
type Proxy struct {
	Name      string
	DependsOn map[string]*Flow
	Calls     []*Call
	Forward   *ProxyForward
}

// GetName returns the flow name
func (proxy *Proxy) GetName() string {
	return proxy.Name
}

// GetDependencies returns the dependencies of the given flow
func (proxy *Proxy) GetDependencies() map[string]*Flow {
	return proxy.DependsOn
}

// GetCalls returns the calls of the given flow
func (proxy *Proxy) GetCalls() []*Call {
	return proxy.Calls
}

// GetInput returns the input of the given flow
func (proxy *Proxy) GetInput() *ParameterMap {
	return nil
}

// GetOutput returns the output of the given flow
func (proxy *Proxy) GetOutput() *ParameterMap {
	return nil
}

// ProxyForward represents the service endpoint where the proxy should forward the stream to when all calls succeed.
type ProxyForward struct {
	Endpoint string
	Header   Header
	Rollback *RollbackCall
}
