package specs

// Manifest holds a collection of definitions and resources
type Manifest struct {
	Flows     []Flow
	Endpoints []Endpoint
	Services  []Service
	Callers   []Caller
}

// Flow defines a set of calls that should be called chronologically and produces an output message.
// Calls could reference other resources when constructing messages.
// All references are strictly typed. Properties are fetched from the given proto buffers or inputs.
//
// All flows should contain a unique name.
// Calls are nested inside of flows and contain two labels, a unique name within the flow and the service and method to be called.
// A dependency reference structure is generated within the flow which allows Maestro to figure out which calls could be called parallel to improve performance.
type Flow struct {
	Name   string
	Input  ParameterMap
	Calls  []Call
	Output ParameterMap
}

// Endpoint exposes a flow. Endpoints are not parsed by Maestro and have custom implementations in each caller.
// The name of the endpoint represents the flow which should be executed.
type Endpoint struct {
	Flow string
	Body map[string]interface{}
}

// Header represents a collection of key values
type Header map[string]string

// Reference represents a mustach template reference
type Reference struct {
	Resource string
	Path     string
}

// Property represents a value property.
// A value property could contain a constant value or a value reference.
type Property struct {
	Default   string
	Constant  bool
	Reference Reference
}

// ParameterMap is the initial map of parameter names (keys) and their (templated) values (values)
type ParameterMap struct {
	Options    map[string]interface{}
	Header     Header
	Message    []NestedParameterMap
	Repeated   []RepeatedParameterMap
	Properties map[string]Property
}

// NestedParameterMap is a map of parameter names (keys) and their (templated) values (values)
type NestedParameterMap struct {
	Name       string
	Message    []NestedParameterMap
	Repeated   []RepeatedParameterMap
	Properties map[string]Property
}

// RepeatedParameterMap is a map of repeated message blocks/values
type RepeatedParameterMap struct {
	Name       string
	Template   string
	Message    []NestedParameterMap
	Repeated   []RepeatedParameterMap
	Properties map[string]Property
}

// Call calls the given service and method.
// Calls could be executed synchronously or asynchronously.
// All calls are referencing a service method, the service should match the alias defined inside the service.
// The request and response proto messages are used for type definitions.
// A call could contain the request headers, request body, rollback, and the execution type.
type Call struct {
	Name     string
	Endpoint string
	Type     string
	Request  ParameterMap
	Rollback RollbackCall
}

// RollbackCall represents the rollback call which is executed when a call inside a flow failed.
type RollbackCall struct {
	Endpoint string
	Request  ParameterMap
}

// Service represent external service which could be called inside the flows.
// The service name is an alias which could be referenced inside calls.
// The host of the service and proto service method should be defined for each service.
//
// The request and response message defined inside the proto buffers are used for type definitions.
// The FQN (fully qualified name) of the proto method should be used.
// Each service references a caller implementation to be used.
type Service struct {
	Options map[string]interface{}
	Alias   string
	Caller  string
	Host    string
	Proto   string
}

// Caller Each implementation has to be configured and defined before running the service.
// All values are passed as attributes to the callers to be unmarshalled.
// These attributes could be used for configuration purposes
type Caller struct {
	Name       string
	Properties map[string]interface{}
}

// Proxy streams the incoming request to the given service.
// Proxies could define calls that are executed before the request body is forwarded.
// A proxy forward could ideally be used for file uploads or large messages which could not be stored in memory.
type Proxy struct {
	Name    string
	Calls   []Call
	Forward ProxyForward
	Output  ParameterMap
}

// ProxyForward represents the service endpoint where the proxy should forward the stream to when all calls succeed.
type ProxyForward struct {
	Name     string
	Endpoint string
	Rollback RollbackCall
}
