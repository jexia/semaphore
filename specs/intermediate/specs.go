package intermediate

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/jexia/maestro/specs"
	"github.com/zclconf/go-cty/cty"
)

// ParseManifest parses the given intermediate manifest to a specs manifest
func ParseManifest(manifest Manifest) specs.Manifest {
	result := specs.Manifest{
		Callers:   make([]*specs.Caller, len(manifest.Callers)),
		Endpoints: make([]*specs.Endpoint, len(manifest.Endpoints)),
		Services:  make([]*specs.Service, len(manifest.Services)),
		Flows:     make([]*specs.Flow, len(manifest.Flows)),
	}

	for index, caller := range manifest.Callers {
		result.Callers[index] = ParseIntermediateCaller(caller)
	}

	for index, endpoint := range manifest.Endpoints {
		result.Endpoints[index] = ParseIntermediateEndpoint(endpoint)
	}

	for index, service := range manifest.Services {
		result.Services[index] = ParseIntermediateService(service)
	}

	for index, flow := range manifest.Flows {
		result.Flows[index] = ParseIntermediateFlow(flow)
	}

	return result
}

// ParseIntermediateCaller parses the given intermediate caller to a specs caller
func ParseIntermediateCaller(caller Caller) *specs.Caller {
	result := specs.Caller{
		Name: caller.Name,
		Body: make(map[string]interface{}),
	}

	gohcl.DecodeBody(caller.Body, nil, &result.Body)
	return &result
}

// ParseIntermediateEndpoint parses the given intermediate endpoint to a specs endpoint
func ParseIntermediateEndpoint(endpoint Endpoint) *specs.Endpoint {
	result := specs.Endpoint{
		Flow: endpoint.Flow,
		Body: make(map[string]interface{}),
	}

	gohcl.DecodeBody(endpoint.Body, nil, &result.Body)
	return &result
}

// ParseIntermediateService parses the given intermediate service to a specs service
func ParseIntermediateService(service Service) *specs.Service {
	result := specs.Service{
		Options: ParseIntermediateOptions(service.Options),
		Alias:   service.Alias,
		Caller:  service.Caller,
		Host:    service.Host,
		Proto:   service.Proto,
	}

	return &result
}

// ParseIntermediateFlow parses the given intermediate flow to a specs flow
func ParseIntermediateFlow(flow Flow) *specs.Flow {
	result := specs.Flow{
		Name:      flow.Name,
		DependsOn: make(map[string]*specs.Flow, len(flow.DependsOn)),
		Input:     ParseIntermediateParameterMap(flow.Input),
		Calls:     make([]*specs.Call, len(flow.Calls)),
		Output:    ParseIntermediateParameterMap(flow.Output),
	}

	for _, dependency := range flow.DependsOn {
		result.DependsOn[dependency] = nil
	}

	for index, call := range flow.Calls {
		result.Calls[index] = ParseIntermediateCall(call)
	}

	return &result
}

// ParseIntermediateParameterMap parses the given intermediate parameter map to a spec parameter map
func ParseIntermediateParameterMap(params *ParameterMap) *specs.ParameterMap {
	if params == nil {
		return nil
	}

	properties, _ := params.Properties.JustAttributes()

	result := specs.ParameterMap{
		Options:    ParseIntermediateOptions(params.Options),
		Header:     ParseIntermediateHeader(params.Header),
		Nested:     make(map[string]*specs.NestedParameterMap, len(params.Nested)),
		Repeated:   make(map[string]*specs.RepeatedParameterMap, len(params.Repeated)),
		Properties: make(map[string]*specs.Property, len(properties)),
	}

	for _, attr := range properties {
		result.Properties[attr.Name] = ParseIntermediateProperty(attr.Name, attr)
	}

	for _, nested := range params.Nested {
		result.Nested[nested.Name] = ParseIntermediateNestedParameterMap(nested, nested.Name)
	}

	for _, repeated := range params.Repeated {
		result.Repeated[repeated.Name] = ParseIntermediateRepeatedParameterMap(repeated, repeated.Name)
	}

	return &result
}

// ParseIntermediateNestedParameterMap parses the given intermediate parameter map to a spec parameter map
func ParseIntermediateNestedParameterMap(params NestedParameterMap, path string) *specs.NestedParameterMap {
	properties, _ := params.Properties.JustAttributes()
	result := specs.NestedParameterMap{
		Name:       params.Name,
		Path:       path,
		Nested:     make(map[string]*specs.NestedParameterMap, len(params.Nested)),
		Repeated:   make(map[string]*specs.RepeatedParameterMap, len(params.Repeated)),
		Properties: make(map[string]*specs.Property, len(properties)),
	}

	for _, nested := range params.Nested {
		result.Nested[nested.Name] = ParseIntermediateNestedParameterMap(nested, specs.JoinPath(path, nested.Name))
	}

	for _, repeated := range params.Repeated {
		result.Repeated[repeated.Name] = ParseIntermediateRepeatedParameterMap(repeated, specs.JoinPath(path, repeated.Name))
	}

	for _, attr := range properties {
		result.Properties[attr.Name] = ParseIntermediateProperty(specs.JoinPath(path, attr.Name), attr)
	}

	return &result
}

// ParseIntermediateRepeatedParameterMap parses the given intermediate repeated parameter map to a spec repeated parameter map
func ParseIntermediateRepeatedParameterMap(params RepeatedParameterMap, path string) *specs.RepeatedParameterMap {
	properties, _ := params.Properties.JustAttributes()
	result := specs.RepeatedParameterMap{
		Name:       params.Name,
		Path:       path,
		Template:   params.Template,
		Nested:     make(map[string]*specs.NestedParameterMap, len(params.Nested)),
		Repeated:   make(map[string]*specs.RepeatedParameterMap, len(params.Repeated)),
		Properties: make(map[string]*specs.Property, len(properties)),
	}

	for _, nested := range params.Nested {
		result.Nested[nested.Name] = ParseIntermediateNestedParameterMap(nested, specs.JoinPath(path, nested.Name))
	}

	for _, repeated := range params.Repeated {
		result.Repeated[repeated.Name] = ParseIntermediateRepeatedParameterMap(repeated, specs.JoinPath(path, repeated.Name))
	}

	for _, attr := range properties {
		result.Properties[attr.Name] = ParseIntermediateProperty(specs.JoinPath(path, attr.Name), attr)
	}

	return &result
}

// ParseIntermediateHeader parses the given intermediate header to a spec header
func ParseIntermediateHeader(header *Header) specs.Header {
	if header == nil {
		return nil
	}

	attributes, _ := header.Body.JustAttributes()
	result := make(specs.Header, len(attributes))

	for _, attr := range attributes {
		result[attr.Name] = ParseIntermediateProperty(attr.Name, attr)
	}

	return result
}

// ParseIntermediateOptions parses the given intermediate options to a spec options
func ParseIntermediateOptions(options *Options) specs.Options {
	if options == nil {
		return specs.Options{}
	}

	result := specs.Options{}
	gohcl.DecodeBody(options.Body, nil, &result)

	return result
}

// ParseIntermediateCall parses the given intermediate call to a spec call
func ParseIntermediateCall(call Call) *specs.Call {
	result := specs.Call{
		DependsOn: make(map[string]*specs.Call, len(call.DependsOn)),
		Name:      call.Name,
		Endpoint:  call.Endpoint,
		Type:      call.Type,
		Request:   ParseIntermediateParameterMap(call.Request),
		Rollback:  ParseIntermediateRollbackCall(call.Rollback),
	}

	for _, dependency := range call.DependsOn {
		result.DependsOn[dependency] = nil
	}

	return &result
}

// ParseIntermediateRollbackCall parses the given intermediate rollback call to a spec rollback call
func ParseIntermediateRollbackCall(call *RollbackCall) *specs.RollbackCall {
	if call == nil {
		return nil
	}

	result := specs.RollbackCall{
		Endpoint: call.Endpoint,
		Request:  ParseIntermediateParameterMap(call.Request),
	}
	return &result
}

// ParseIntermediateProperty parses the given intermediate property to a spec property
func ParseIntermediateProperty(path string, property *hcl.Attribute) *specs.Property {
	if property == nil {
		return nil
	}

	value, _ := property.Expr.Value(nil)
	result := specs.Property{
		Path: path,
		Expr: property.Expr,
	}

	// Template definitions could be improved to be more consistent
	if value.Type() == cty.String && specs.IsType(value.AsString()) {
		specs.SetType(&result, value)
		return &result
	}

	if value.Type() != cty.String || !specs.IsTemplate(value.AsString()) {
		specs.SetDefaultValue(&result, value)
		return &result
	}

	specs.SetTemplate(&result, value)
	return &result
}
