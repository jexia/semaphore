package openapi3

import (
	"fmt"
	"net/http"
	"sort"
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

	params := transport.NamedParameters(options.Endpoint)
	sort.Strings(params)

	for _, key := range params {
		result.Parameters = append(result.Parameters, GenerateParameter(key, true, ParameterPath, nil))
	}

	if input != nil {
		// ensure header order
		headers := make([]string, 0, len(flow.GetInput().Header))
		for key := range flow.GetInput().Header {
			headers = append(headers, key)
		}

		sort.Strings(headers)

		for _, key := range headers {
			prop := flow.GetInput().Header[key]
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
		result.Description = property.Description
		result.Schema = &Schema{
			Type: types.Open(property.Type()),
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

	object.Components.Schemas[params.Schema] = GenerateSchema(params.Property.Description, params.Property.Template)
}

// GenerateSchema generates a new schema for the given property
func GenerateSchema(description string, property specs.Template) *Schema {
	result := &Schema{
		Description: description,
		Type:        types.Open(property.Type()),
	}

	switch {
	case property.Scalar != nil:
		result.Default = property.Scalar.Default

		break
	case property.Message != nil:
		result.Properties = make(map[string]*Schema, len(property.Message))

		for _, nested := range property.Message {
			result.Properties[nested.Name] = GenerateSchema(nested.Description, nested.Template)

			if nested.Label == labels.Required {
				result.Required = append(result.Required, nested.Name)
			}
		}

		break
	case property.Enum != nil:
		// ensure property enum order
		result.Enum = make([]interface{}, len(property.Enum.Keys))
		keys := make([]int, 0, len(property.Enum.Positions))

		for key := range property.Enum.Positions {
			keys = append(keys, int(key))
		}

		sort.Ints(keys)

		for pos, key := range keys {
			result.Enum[pos] = property.Enum.Positions[int32(key)].Key
		}

		break
	case property.Repeated != nil:
		template, err := property.Repeated.Template()
		if err != nil {
			panic(err)
		}

		return &Schema{
			Description: description,
			Items:       GenerateSchema("", template),
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
