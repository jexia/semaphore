package hcl

import (
	"github.com/hashicorp/hcl/v2"
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
	Input     *InputParameterMap `hcl:"input,block"`
	Resources []Node             `hcl:"resource,block"`
	Output    *ParameterMap      `hcl:"output,block"`
}

// ParameterMap is the initial map of parameter names (keys) and their (templated) values (values)
type ParameterMap struct {
	Schema     string                 `hcl:"schema,label"`
	Options    *Options               `hcl:"options,block"`
	Header     *Header                `hcl:"header,block"`
	Nested     []NestedParameterMap   `hcl:"message,block"`
	Repeated   []RepeatedParameterMap `hcl:"repeated,block"`
	Properties hcl.Body               `hcl:",remain"`
}

// Endpoint intermediate specification
type Endpoint struct {
	Flow     string   `hcl:"flow,label"`
	Listener string   `hcl:"listener,label"`
	Options  hcl.Body `hcl:",remain"`
}

// Header represents a collection of key values
type Header struct {
	Body hcl.Body `hcl:",remain"`
}

// InputParameterMap is the initial map of parameter names (keys) and their (templated) values (values)
type InputParameterMap struct {
	Schema     string                      `hcl:"schema,label"`
	Options    *Options                    `hcl:"options,block"`
	Header     []string                    `hcl:"header,optional"`
	Nested     []NestedParameterMap        `hcl:"message,block"`
	Repeated   []InputRepeatedParameterMap `hcl:"repeated,block"`
	Properties hcl.Body                    `hcl:",remain"`
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

// Node intermediate specification
type Node struct {
	Name      string   `hcl:"name,label"`
	DependsOn []string `hcl:"depends_on,optional"`
	Type      string   `hcl:"type,optional"`
	Request   *Call    `hcl:"request,block"`
	Rollback  *Call    `hcl:"rollback,block"`
}

// Call intermediate specification
type Call struct {
	Service    string                 `hcl:"service,label"`
	Method     string                 `hcl:"method,label"`
	Options    *Options               `hcl:"options,block"`
	Header     *Header                `hcl:"header,block"`
	Nested     []NestedParameterMap   `hcl:"message,block"`
	Repeated   []RepeatedParameterMap `hcl:"repeated,block"`
	Properties hcl.Body               `hcl:",remain"`
}

// Service specification
type Service struct {
	Package   string   `hcl:"package,label"`
	Name      string   `hcl:"name,label"`
	Transport string   `hcl:"transport,label"`
	Codec     string   `hcl:"codec,label"`
	Host      string   `hcl:"host"`
	Methods   []Method `hcl:"method,block"`
	Options   *Options `hcl:"options,block"`
}

// Method represents a service method
type Method struct {
	Name     string   `hcl:"name,label"`
	Request  string   `hcl:"request,optional"`
	Response string   `hcl:"response,optional"`
	Options  *Options `hcl:"options,block"`
}

// Proxy specification
type Proxy struct {
	Name      string       `hcl:"name,label"`
	DependsOn []string     `hcl:"depends_on,optional"`
	Resources []Node       `hcl:"resource,block"`
	Forward   ProxyForward `hcl:"forward,block"`
}

// ProxyForward specification
type ProxyForward struct {
	Service string  `hcl:"service,label"`
	Header  *Header `hcl:"header,block"`
}
