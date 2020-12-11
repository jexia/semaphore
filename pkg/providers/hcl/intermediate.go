package hcl

import "github.com/hashicorp/hcl/v2"

// Manifest intermediate specs
type Manifest struct {
	LogLevel         string        `hcl:"log_level,optional"`
	GraphQL          *GraphQL      `hcl:"graphql,block"`
	HTTP             *HTTP         `hcl:"http,block"`
	GRPC             *GRPC         `hcl:"grpc,block"`
	Prometheus       *Prometheus   `hcl:"prometheus,block"`
	Protobuffers     []string      `hcl:"protobuffers,optional"`
	Avro             []string      `hcl:"avro,optional"`
	Openapi3         []string      `hcl:"openapi3,optional"`
	Include          []string      `hcl:"include,optional"`
	Error            *ParameterMap `hcl:"error,block"`
	Flows            []Flow        `hcl:"flow,block"`
	Proxy            []Proxy       `hcl:"proxy,block"`
	Endpoints        []Endpoint    `hcl:"endpoint,block"`
	Services         []Service     `hcl:"service,block"`
	ServiceSelector  []Services    `hcl:"services,block"`
	DiscoveryServers []Discovery   `hcl:"discovery,block"`
}

// GraphQL represents the GraphQL option definitions
type GraphQL struct {
	Address string `hcl:"address"`
}

// HTTP represent the HTTP option definitions.
type HTTP struct {
	Address      string   `hcl:"address"`
	Origin       []string `hcl:"origin,optional"`
	CertFile     string   `hcl:"cert_file,optional"`
	KeyFile      string   `hcl:"key_file,optional"`
	ReadTimeout  string   `hcl:"read_timeout,optional"`
	WriteTimeout string   `hcl:"write_timeout,optional"`
}

// GRPC represent the gRPC option definitions.
type GRPC struct {
	Address string `hcl:"address"`
}

// Prometheus represent the prometheus option definitions.
type Prometheus struct {
	Address string `hcl:"address"`
}

// ResourceBlock contains resources and references.
type ResourceBlock struct {
	References []Resources `hcl:"resources,block"`
	Resources  []Resource  `hcl:"resource,block"`
}

// Before intermediate specification.
type Before ResourceBlock

// Condition represents a condition on which the scoped resources should be executed if true.
type Condition struct {
	ResourceBlock `hcl:",remain"`

	Expression string      `hcl:"expression,label"`
	Conditions []Condition `hcl:"if,block"`
}

// Flow intermediate specification.
type Flow struct {
	Condition `hcl:",remain"`

	Name    string             `hcl:"name,label"`
	Error   *ParameterMap      `hcl:"error,block"`
	Input   *InputParameterMap `hcl:"input,block"`
	OnError *OnError           `hcl:"on_error,block"`
	Before  *Before            `hcl:"before,block"`
	Output  *ParameterMap      `hcl:"output,block"`
}

// BaseParameterMap contains a set of generic fields.
type BaseParameterMap struct {
	Nested     []NestedParameterMap   `hcl:"message,block"`
	Repeated   []RepeatedParameterMap `hcl:"repeated,block"`
	Properties hcl.Body               `hcl:",remain"`
}

type OutputParameterMap struct {
	BaseParameterMap `hcl:",remain"`

	Schema string `hcl:"schema,label"`
}

// ParameterMap is the initial map of parameter names (keys) and their (templated) values (values).
type ParameterMap struct {
	Options *BlockOptions       `hcl:"options,block"`
	Status  *int64              `hcl:"status,optional"`
	Header  *Header             `hcl:"header,block"`
	Payload *OutputParameterMap `hcl:"payload,block"`
}

// Resources represent a collection of resources which are references or custom defined functions.
type Resources struct {
	Properties hcl.Body `hcl:",remain"`
}

// Endpoint intermediate specification.
type Endpoint struct {
	Flow     string   `hcl:"flow,label"`
	Listener string   `hcl:"listener,label"`
	Options  hcl.Body `hcl:",remain"`
}

// Header represents a collection of key values.
type Header struct {
	Body hcl.Body `hcl:",remain"`
}

// InputParameterMap is the initial map of parameter names (keys) and their (templated) values (values).
type InputParameterMap struct {
	Schema  string        `hcl:"schema,label"`
	Options *BlockOptions `hcl:"options,block"`
	Header  []string      `hcl:"header,optional"`
}

// BlockOptions holds the raw options.
type BlockOptions struct {
	Body hcl.Body `hcl:",remain"`
}

// NestedParameterMap is a map of parameter names (keys) and their (templated) values (values).
type NestedParameterMap struct {
	BaseParameterMap `hcl:",remain"`

	Name string `hcl:"name,label"`
}

// InputRepeatedParameterMap is a map of repeated message blocks/values.
type InputRepeatedParameterMap struct {
	Name       string                      `hcl:"name,label"`
	Nested     []NestedParameterMap        `hcl:"message,block"`
	Repeated   []InputRepeatedParameterMap `hcl:"repeated,block"`
	Properties hcl.Body                    `hcl:",remain"`
}

// RepeatedParameterMap is a map of repeated message blocks/values.
type RepeatedParameterMap struct {
	BaseParameterMap `hcl:",remain"`

	Name     string `hcl:"name,label"`
	Template string `hcl:"template,label"`
}

// Resource intermediate specification.
type Resource struct {
	Name         string        `hcl:"name,label"`
	DependsOn    []string      `hcl:"depends_on,optional"`
	Request      *Call         `hcl:"request,block"`
	Rollback     *Call         `hcl:"rollback,block"`
	OnError      *OnError      `hcl:"on_error,block"`
	ExpectStatus []int         `hcl:"expect_status,optional"`
	Error        *ParameterMap `hcl:"error,block"`
}

// OnError intermediate specification.
type OnError struct {
	Schema string        `hcl:"schema,optional"`
	Params *BlockOptions `hcl:"params,block"`
	Body   hcl.Body      `hcl:",remain"`
}

// Call intermediate specification.
type Call struct {
	BaseParameterMap `hcl:",remain"`

	Service    string        `hcl:"service,label"`
	Method     string        `hcl:"method,label"`
	Parameters *BlockOptions `hcl:"params,block"`
	Options    *BlockOptions `hcl:"options,block"`
	Header     *Header       `hcl:"header,block"`
}

// Service specification.
type Service struct {
	Package   string        `hcl:"package,label"`
	Name      string        `hcl:"name,label"`
	Transport string        `hcl:"transport,optional"`
	Codec     string        `hcl:"codec,optional"`
	Host      string        `hcl:"host,optional"`
	Methods   []Method      `hcl:"method,block"`
	Options   *BlockOptions `hcl:"options,block"`
	Resolver  string        `hcl:"resolver,optional"`
}

// ServiceSelector targets any service match in the given service selector.
type ServiceSelector struct {
	Pattern       string   `hcl:"pattern,label"`
	Host          string   `hcl:"host,optional"`
	Transport     string   `hcl:"transport,optional"`
	Codec         string   `hcl:"codec,optional"`
	RequestCodec  string   `hcl:"request_codec,optional"`
	ResponseCodec string   `hcl:"response_codec,optional"`
	Resolver      string   `hcl:"resolver,optional"`
	Options       hcl.Body `hcl:",remain"`
}

// Services specification.
type Services struct {
	Selectors []ServiceSelector `hcl:"select,block"`
}

// Method represents a service method
type Method struct {
	Name    string        `hcl:"name,label"`
	Input   string        `hcl:"request,optional"`
	Output  string        `hcl:"response,optional"`
	Options *BlockOptions `hcl:"options,block"`
}

// Proxy specification.
type Proxy struct {
	Condition `hcl:",remain"`

	Name    string         `hcl:"name,label"`
	Error   *ParameterMap  `hcl:"error,block"`
	Input   *ProxyInput    `hcl:"input,block"`
	OnError *OnError       `hcl:"on_error,block"`
	Before  *Before        `hcl:"before,block"`
	Forward ProxyForward   `hcl:"forward,block"`
	Rewrite []ProxyRewrite `hcl:"rewrite,block"`
}

// ProxyInput represents the proxy input block.
type ProxyInput struct {
	Options *BlockOptions `hcl:"options,block"`
	Header  []string      `hcl:"header,optional"`
	Params  string        `hcl:"params,optional"`
}

// ProxyForward specification.
type ProxyForward struct {
	Service string  `hcl:"service,label"`
	Header  *Header `hcl:"header,block"`
}

// ProxyRewrite describes rewrite rules. Allows rewriting/replacement URL segments
//
// Example:
//
//	endpoint "test" "http" {
//		method = "GET"
//		endpoint = "/prefix/*tail"
//	}
//
//	proxy "test" {
//		forward "com.semaphore.TestService" {}
//		rewrite "/prefix/(?P<first>\w+)/(?P<second>\w+)" "/<second>/<first>" {} // swap URL segments
//		rewrite "/prefix/(?P<tail>.*)" "/<tail>" {} // strip static `/prefix`
//	}
//
// NOTE: multiple rewrite rules can be specified but only one can be applied (the first one which
// matches the pattern) so put strict rules first (by priority) and keep wildcards in the bottom
type ProxyRewrite struct {
	Pattern  string   `hcl:"pattern,label"`
	Template string   `hcl:"template,label"`
	Options  hcl.Body `hcl:",remain"`
}

// Discovery describes a service discovery client configuration
//
// Examples:
//
//	discovery "consul" {
//		address = "http://localhost:8500"
//	}
//
//	discovery "foobar" {
//		provider = "consul"
//		address = "http://localhost:8500"
//	}
type Discovery struct {
	Name     string `hcl:"name,label"`
	Provider string `hcl:"provider,optional"`
	Address  string `hcl:"address"`
}
