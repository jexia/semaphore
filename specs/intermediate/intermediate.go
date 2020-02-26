package intermediate

import (
	"github.com/hashicorp/hcl/v2"
)

const (
	// Ext represents the intermediate file extention
	Ext = ".hcl"
)

// Manifest intermediate specs
type Manifest struct {
	Flows     []Flow     `hcl:"flow,block"`
	Proxy     []Proxy    `hcl:"proxy,block"`
	Endpoints []Endpoint `hcl:"endpoint,block"`
	Services  []Service  `hcl:"service,block"`
}

// Flow intermediate specification
type Flow struct {
	Name      string             `hcl:"name,label"`
	DependsOn []string           `hcl:"depends_on,optional"`
	Schema    string             `hcl:"schema,optional"`
	Input     *InputParameterMap `hcl:"input,block"`
	Calls     []Call             `hcl:"call,block"`
	Output    *ParameterMap      `hcl:"output,block"`
}

// Endpoint intermediate specification
type Endpoint struct {
	Flow     string   `hcl:"flow,label"`
	Listener string   `hcl:"listener,label"`
	Codec    string   `hcl:"codec"`
	Body     hcl.Body `hcl:",remain"`
}

// Header represents a collection of key values
type Header struct {
	Body hcl.Body `hcl:",remain"`
}

// InputParameterMap is the initial map of parameter names (keys) and their (templated) values (values)
type InputParameterMap struct {
	Options    *Options                    `hcl:"options,block"`
	Header     *Header                     `hcl:"header,block"`
	Nested     []NestedParameterMap        `hcl:"message,block"`
	Repeated   []InputRepeatedParameterMap `hcl:"repeated,block"`
	Properties hcl.Body                    `hcl:",remain"`
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

// InputRepeatedParameterMap is a map of repeated message blocks/values
type InputRepeatedParameterMap struct {
	Name       string                      `hcl:"name,label"`
	Nested     []NestedParameterMap        `hcl:"message,block"`
	Repeated   []InputRepeatedParameterMap `hcl:"repeated,block"`
	Properties hcl.Body                    `hcl:",remain"`
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
	Codec   string   `hcl:"codec"`
	Schema  string   `hcl:"schema"`
}

// Proxy specification
type Proxy struct {
	Name      string       `hcl:"name,label"`
	DependsOn []string     `hcl:"depends_on,optional"`
	Calls     []Call       `hcl:"call,block"`
	Forward   ProxyForward `hcl:"forward,block"`
}

// ProxyForward specification
type ProxyForward struct {
	Endpoint string        `hcl:"endpoint,label"`
	Header   *Header       `hcl:"header,block"`
	Rollback *RollbackCall `hcl:"rollback,block"`
}
