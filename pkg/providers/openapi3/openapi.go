package openapi3

import "github.com/jexia/semaphore/pkg/providers/openapi3/types"

// Object represents the version 3 open api specification
type Object struct {
	Version    string               `json:"openapi" yaml:"openapi"`
	Servers    []*ServerRef         `json:"servers,omitempty" yaml:"servers,omitempty"`
	Info       Info                 `json:"info" yaml:"info"`
	Paths      map[string]*PathItem `json:"paths,omitempty" yaml:"paths,omitempty"`
	Components *Components          `json:"components,omitempty" yaml:"components,omitempty"`
}

// ServerRef represents a server reference
type ServerRef struct {
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	URL         string `json:"url,omitempty" yaml:"url,omitempty"`
}

// Info provides metadata about the API.
// The metadata MAY be used by the clients if needed,
// and MAY be presented in editing or documentation
// generation tools for convenience.
type Info struct {
	Title       string   `json:"title" yaml:"title"`
	Version     string   `json:"version" yaml:"version"`
	Description string   `json:"description,omitempty" yaml:"description,omitempty"`
	Contact     *Contact `json:"contact,omitempty" yaml:"contact,omitempty"`
}

// Contact information for the exposed API.
type Contact struct {
	Name  string `json:"name,omitempty" yaml:"name,omitempty"`
	URL   string `json:"url,omitempty" yaml:"url,omitempty"`
	Email string `json:"email,omitempty" yaml:"email,omitempty"`
}

// Components Holds a set of reusable objects for different
// aspects of the OAS. All objects defined within the
// components object will have no effect on the API unless
// they are explicitly referenced from properties outside
// the components object.
type Components struct {
	Schemas map[string]*Schema `json:"schemas,omitempty" yaml:"schemas,omitempty"`
}

// Schema Object allows the definition of input and output
// data types. These types can be objects, but also
// primitives and arrays.
type Schema struct {
	Reference   string             `json:"$ref,omitempty" yaml:"$ref,omitempty"`
	Type        types.Type         `json:"type,omitempty" yaml:"type,omitempty"`
	Description string             `json:"description,omitempty" yaml:"description,omitempty"`
	Required    []string           `json:"required,omitempty" yaml:"required,omitempty"`
	Items       *Schema            `json:"items,omitempty" yaml:"items,omitempty"`
	Properties  map[string]*Schema `json:"properties,omitempty" yaml:"properties,omitempty"`
	Default     interface{}        `json:"default,omitempty" yaml:"default,omitempty"`
	Enum        []interface{}      `json:"enum,omitempty" yaml:"enum,omitempty"`
}

// PathItem describes the operations available on a single path.
// A Path Item MAY be empty, due to ACL constraints.
// The path itself is still exposed to the documentation
// viewer but they will not know which operations and
// parameters are available.
type PathItem struct {
	Ref         string     `json:"$ref,omitempty" yaml:"$ref,omitempty"`
	Description string     `json:"description,omitempty" yaml:"description,omitempty"`
	Get         *Operation `json:"get,omitempty" yaml:"get,omitempty"`
	Put         *Operation `json:"put,omitempty" yaml:"put,omitempty"`
	Post        *Operation `json:"post,omitempty" yaml:"post,omitempty"`
	Delete      *Operation `json:"delete,omitempty" yaml:"delete,omitempty"`
	Options     *Operation `json:"options,omitempty" yaml:"options,omitempty"`
	Head        *Operation `json:"head,omitempty" yaml:"head,omitempty"`
	Patch       *Operation `json:"patch,omitempty" yaml:"patch,omitempty"`
	Trace       *Operation `json:"trace,omitempty" yaml:"trace,omitempty"`
}

// Operation describes a single API operation on a path.
type Operation struct {
	Tags        []string             `json:"tags,omitempty" yaml:"tags,omitempty"`
	Description string               `json:"description,omitempty" yaml:"description,omitempty"`
	OperationID string               `json:"operationId,omitempty" yaml:"operationId,omitempty"`
	Parameters  []*Parameter         `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	RequestBody *RequestBody         `json:"requestBody,omitempty" yaml:"requestBody,omitempty"`
	Responses   map[string]*Response `json:"responses,omitempty" yaml:"responses,omitempty"`
}

// ParameterIn represents the parameter in options
type ParameterIn string

// Available parameter in options
const (
	ParameterQuery  ParameterIn = "query"
	ParameterHeader ParameterIn = "header"
	ParameterPath   ParameterIn = "path"
	ParameterCookie ParameterIn = "cookie"
)

// Parameter a list of parameters that are applicable for this operation.
type Parameter struct {
	Name        string      `json:"name,omitempty" yaml:"name,omitempty"`
	In          ParameterIn `json:"in,omitempty" yaml:"in,omitempty"`
	Description string      `json:"description,omitempty" yaml:"description,omitempty"`
	Required    bool        `json:"required" yaml:"required"`
	Schema      *Schema     `json:"schema,omitempty" yaml:"schema,omitempty"`
}

// RequestBody the request body applicable for this operation.
// The requestBody is only supported in HTTP methods where
// the HTTP 1.1 specification RFC7231 has explicitly defined
// semantics for request bodies
type RequestBody struct {
	Description string               `json:"description,omitempty" yaml:"description,omitempty"`
	Content     map[string]MediaType `json:"content,omitempty" yaml:"content,omitempty"`
	Required    bool                 `json:"required" yaml:"required"`
}

// MediaType each Media Type Object provides schema and examples
// for the media type identified by its key.
type MediaType struct {
	Schema Schema `json:"schema,omitempty" yaml:"schema,omitempty"`
}

// Response the list of possible responses as they are returned
// from executing this operation.
type Response struct {
	Description string               `json:"description" yaml:"description"`
	Content     map[string]MediaType `json:"content,omitempty" yaml:"content,omitempty"`
}
