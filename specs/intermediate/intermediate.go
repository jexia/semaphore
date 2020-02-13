package intermediate

import (
	"github.com/hashicorp/hcl/v2"
)

// Manifest intermediate specs
type Manifest struct {
	Flows     []Flow     `hcl:"flow,block"`
	Endpoints []Endpoint `hcl:"endpoint,block"`
	Services  []Service  `hcl:"service,block"`
	Callers   []Caller   `hcl:"caller,block"`
}

// Flow intermediate specification
type Flow struct {
	Name      string        `hcl:"name,label"`
	DependsOn []string      `hcl:"depends_on,optional"`
	Input     *ParameterMap `hcl:"input,block"`
	Calls     []Call        `hcl:"call,block"`
	Output    *ParameterMap `hcl:"output,block"`
}

// Endpoint intermediate specification
type Endpoint struct {
	Flow string   `hcl:"flow,label"`
	Body hcl.Body `hcl:",remain"`
}

// Header represents a collection of key values
type Header struct {
	Body hcl.Body `hcl:",remain"`
}

// ParameterMap is the initial map of parameter names (keys) and their (templated) values (values)
type ParameterMap struct {
	Options    *Options               `hcl:"options,block"`
	Header     *Header                `hcl:"header,block"`
	Nested     []NestedParameterMap   `hcl:"message,block"`
	Repeated   []RepeatedParameterMap `hcl:"repeated,block"`
	Properties hcl.Body               `hcl:",remain"`
}

// Options holds the raw options
type Options struct {
	Body hcl.Body `hcl:",remain"`
}

// NestedParameterMap is a map of parameter names (keys) and their (templated) values (values)
type NestedParameterMap struct {
	Name       string                 `hcl:"name,label"`
	Nested     []NestedParameterMap   `hcl:"message,block"`
	Repeated   []RepeatedParameterMap `hcl:"repeated,block"`
	Properties hcl.Body               `hcl:",remain"`
}

// RepeatedParameterMap is a map of repeated message blocks/values
type RepeatedParameterMap struct {
	Name       string                 `hcl:"name,label"`
	Template   string                 `hcl:"template,label"`
	Nested     []NestedParameterMap   `hcl:"message,block"`
	Repeated   []RepeatedParameterMap `hcl:"repeated,block"`
	Properties hcl.Body               `hcl:",remain"`
}

// Call intermediate specification
type Call struct {
	DependsOn []string      `hcl:"depends_on,optional"`
	Name      string        `hcl:"name,label"`
	Endpoint  string        `hcl:"endpoint,label"`
	Type      string        `hcl:"type,optional"`
	Request   *ParameterMap `hcl:"request,block"`
	Rollback  *RollbackCall `hcl:"rollback,block"`
}

// RollbackCall intermediate specification
type RollbackCall struct {
	Endpoint string        `hcl:"endpoint,label"`
	Request  *ParameterMap `hcl:"request,block"`
}

// Service specification
type Service struct {
	Options *Options `hcl:"options,block"`
	Alias   string   `hcl:"alias,label"`
	Caller  string   `hcl:"caller,label"`
	Host    string   `hcl:"host"`
	Schema  string   `hcl:"schema"`
}

// Caller specification
type Caller struct {
	Name string   `hcl:"name,label"`
	Body hcl.Body `hcl:",remain"`
}

// Proxy specification
type Proxy struct {
	Name    string        `hcl:"name,label"`
	Calls   []Call        `hcl:"call,block"`
	Forward ProxyForward  `hcl:"forward,block"`
	Output  *ParameterMap `hcl:"output,block"`
}

// ProxyForward specification
type ProxyForward struct {
	Name     string        `hcl:"name,label"`
	Endpoint string        `hcl:"endpoint,label"`
	Rollback *RollbackCall `hcl:"rollback,block"`
}
