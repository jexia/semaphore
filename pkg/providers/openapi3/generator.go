package openapi3

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/jexia/semaphore/pkg/providers/openapi3/types"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/labels"
	transport "github.com/jexia/semaphore/pkg/transport/http"
)

// Generate generates a openapi v3.0 specification object
func Generate(endpoints specs.EndpointList, flows specs.FlowListInterface) (*Object, error) {
	result := Object{
		Version: "3.0.0",
	}

	for _, endpoint := range endpoints {
		// OpenAPI specs are ment for HTTP endpoints
		if endpoint.Listener != "http" {
			continue
		}

		flow := flows.Get(endpoint.Flow)
		if flow == nil {
			continue
		}

		err := IncludeEndpoint(&result, endpoint, flow)
		if err != nil {
			return nil, err
		}
	}

	return &result, nil
}

// IncludeEndpoint includes the given endpoint into the object paths
func IncludeEndpoint(object *Object, endpoint *specs.Endpoint, flow specs.FlowInterface) error {
	options, err := transport.ParseEndpointOptions(endpoint.Options)
	if err != nil {
		return err
	}

	if object.Paths == nil {
		object.Paths = make(map[string]*PathItem)
	}

	path := options.Endpoint
	params := transport.RawNamedParameters(options.Endpoint)
	for _, param := range params {
		single := transport.NamedParameters(param)
		template := fmt.Sprintf("{%s}", single[0])
		path = strings.Replace(path, param, template, 1)
	}

	if object.Paths[path] == nil {
		object.Paths[path] = &PathItem{}
	}

	operation := GenerateOperation(object, endpoint, options, flow)
	result := object.Paths[path]

	switch options.Method {
	case http.MethodGet:
		result.Get = operation
	case http.MethodPut:
		result.Put = operation
	case http.MethodPost:
		result.Post = operation
	case http.MethodDelete:
		result.Delete = operation
	case http.MethodOptions:
		result.Options = operation
	case http.MethodHead:
		result.Head = operation
	case http.MethodPatch:
		result.Patch = operation
	case http.MethodTrace:
		result.Trace = operation
	}

	return nil
}

// GenerateOperation generates a operation object from the given endpoint and options
func GenerateOperation(object *Object, endpoint *specs.Endpoint, options *transport.EndpointOptions, flow specs.FlowInterface) *Operation {
	input := flow.GetInput()
	output := flow.GetOutput()
	result := &Operation{}

	for _, key := range transport.NamedParameters(options.Endpoint) {
		result.Parameters = append(result.Parameters, GenerateParameter(key, true, ParameterPath, nil))
	}

	if input != nil {
		for key, prop := range flow.GetInput().Header {
			result.Parameters = append(result.Parameters, GenerateParameter(key, false, ParameterHeader, prop))
		}

		if input.Property != nil && flow.GetForward() == nil && RequiresRequestBody(options.Method) {
			result.RequestBody = &RequestBody{
				Content: map[string]MediaType{
					string(transport.ApplicationJSON): {
						Schema: Schema{
							Reference: fmt.Sprintf("#/components/schemas/%s", input.Schema),
						},
					},
				},
			}

			IncludeParameterMap(object, input)
		}
	}

	if output != nil {
		if input.Property != nil {
			result.Responses = map[string]*Response{
				"default": {
					Content: map[string]MediaType{
						string(transport.ApplicationJSON): {
							Schema: Schema{
								Reference: fmt.Sprintf("#/components/schemas/%s", output.Schema),
							},
						},
					},
				},
			}
		}

		IncludeParameterMap(object, output)
	}

	return result
}

// GenerateParameter includes the given parameter to the available parameters
func GenerateParameter(key string, required bool, in ParameterIn, property *specs.Property) *Parameter {
	result := &Parameter{
		Name:     key,
		In:       in,
		Required: required,
		Schema: &Schema{
			Type: types.String,
		},
	}

	if property != nil {
		result.Description = property.Comment
		result.Schema = &Schema{
			Type: types.Open(property.Type),
		}
	}

	return result
}

// IncludeParameterMap includes the given parameters into the object schema components
func IncludeParameterMap(object *Object, params *specs.ParameterMap) {
	if params == nil {
		return
	}

	if object.Components == nil {
		object.Components = &Components{}
	}

	if object.Components.Schemas == nil {
		object.Components.Schemas = map[string]*Schema{}
	}

	object.Components.Schemas[params.Schema] = GenerateSchema(params.Property)
}

// GenerateSchema generates a new schema for the given property
func GenerateSchema(property *specs.Property) *Schema {
	if property == nil {
		return nil
	}

	result := &Schema{
		Description: property.Comment,
		Default:     property.Default,
		Type:        types.Open(property.Type),
	}

	if property.Nested != nil {
		result.Properties = make(map[string]*Schema, len(property.Nested))
	}

	for key, nested := range property.Nested {
		result.Properties[key] = GenerateSchema(nested)

		if nested.Label == labels.Required {
			result.Required = append(result.Required, key)
		}
	}

	if property.Enum != nil {
		result.Enum = make([]interface{}, 0, len(property.Enum.Keys))

		for key := range property.Enum.Keys {
			result.Enum = append(result.Enum, key)
		}
	}

	return result
}

// RequiresRequestBody checks whether the given method requires a request body.
// This method does not validate the given value and expects to receive a valid HTTP method.
func RequiresRequestBody(method string) bool {
	if method == http.MethodGet {
		return false
	}

	if method == http.MethodDelete {
		return false
	}

	return true
}
