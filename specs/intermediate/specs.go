package intermediate

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/jexia/maestro/specs"
	"github.com/zclconf/go-cty/cty"
)

// ParseManifest parses the given intermediate manifest to a specs manifest
func ParseManifest(manifest Manifest, functions specs.CustomDefinedFunctions) (*specs.Manifest, error) {
	result := &specs.Manifest{
		Callers:   make([]*specs.Caller, len(manifest.Callers)),
		Endpoints: make([]*specs.Endpoint, len(manifest.Endpoints)),
		Services:  make([]*specs.Service, len(manifest.Services)),
		Flows:     make([]*specs.Flow, len(manifest.Flows)),
		Proxy:     make([]*specs.Proxy, len(manifest.Proxy)),
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
		flow, err := ParseIntermediateFlow(flow, functions)
		if err != nil {
			return nil, err
		}

		result.Flows[index] = flow
	}

	for index, proxy := range manifest.Proxy {
		proxy, err := ParseIntermediateProxy(proxy, functions)
		if err != nil {
			return nil, err
		}

		result.Proxy[index] = proxy
	}

	return result, nil
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
		Codec:   service.Codec,
		Schema:  service.Schema,
	}

	return &result
}

// ParseIntermediateFlow parses the given intermediate flow to a specs flow
func ParseIntermediateFlow(flow Flow, functions specs.CustomDefinedFunctions) (*specs.Flow, error) {
	input, err := ParseIntermediateParameterMap(ParseIntermediateInputParameterMap(flow.Input), functions)
	if err != nil {
		return nil, err
	}

	output, err := ParseIntermediateParameterMap(flow.Output, functions)
	if err != nil {
		return nil, err
	}

	result := specs.Flow{
		Name:      flow.Name,
		DependsOn: make(map[string]*specs.Flow, len(flow.DependsOn)),
		Schema:    flow.Schema,
		Input:     input,
		Calls:     make([]*specs.Call, len(flow.Calls)),
		Output:    output,
	}

	for _, dependency := range flow.DependsOn {
		result.DependsOn[dependency] = nil
	}

	for index, call := range flow.Calls {
		call, err := ParseIntermediateCall(call, functions)
		if err != nil {
			return nil, err
		}

		result.Calls[index] = call
	}

	return &result, nil
}

// ParseIntermediateProxy parses the given intermediate proxy to a specs proxy
func ParseIntermediateProxy(proxy Proxy, functions specs.CustomDefinedFunctions) (*specs.Proxy, error) {
	forward, err := ParseIntermediateProxyForward(proxy.Forward, functions)
	if err != nil {
		return nil, err
	}

	result := specs.Proxy{
		Name:      proxy.Name,
		DependsOn: make(map[string]*specs.Flow, len(proxy.DependsOn)),
		Calls:     make([]*specs.Call, len(proxy.Calls)),
		Forward:   forward,
	}

	for _, dependency := range proxy.DependsOn {
		result.DependsOn[dependency] = nil
	}

	for index, call := range proxy.Calls {
		call, err := ParseIntermediateCall(call, functions)
		if err != nil {
			return nil, err
		}

		result.Calls[index] = call
	}

	return &result, nil
}

// ParseIntermediateProxyForward parses the given intermediate proxy forward to a specs proxy forward
func ParseIntermediateProxyForward(proxy ProxyForward, functions specs.CustomDefinedFunctions) (*specs.ProxyForward, error) {
	result := specs.ProxyForward{
		Endpoint: proxy.Endpoint,
	}

	if proxy.Header != nil {
		header, err := ParseIntermediateHeader(proxy.Header, functions)
		if err != nil {
			return nil, err
		}

		result.Header = header
	}

	if proxy.Rollback != nil {
		rollback, err := ParseIntermediateRollbackCall(proxy.Rollback, functions)
		if err != nil {
			return nil, err
		}

		result.Rollback = rollback
	}

	return &result, nil
}

// ParseIntermediateInputParameterMap parses the given input parameter map
func ParseIntermediateInputParameterMap(params *InputParameterMap) *ParameterMap {
	if params == nil {
		return nil
	}

	result := &ParameterMap{
		Options:    params.Options,
		Header:     params.Header,
		Nested:     params.Nested,
		Repeated:   make([]RepeatedParameterMap, len(params.Repeated)),
		Properties: params.Properties,
	}

	for index, repeated := range params.Repeated {
		result.Repeated[index] = ParseIntermediateInputRepeatedParameterMap(repeated)
	}

	return result
}

// ParseIntermediateInputRepeatedParameterMap parses the given input repeated parameter map
func ParseIntermediateInputRepeatedParameterMap(repeated InputRepeatedParameterMap) RepeatedParameterMap {
	result := RepeatedParameterMap{
		Name:       repeated.Name,
		Nested:     repeated.Nested,
		Repeated:   make([]RepeatedParameterMap, len(repeated.Repeated)),
		Properties: repeated.Properties,
	}

	for index, repeated := range repeated.Repeated {
		result.Repeated[index] = ParseIntermediateInputRepeatedParameterMap(repeated)
	}

	return result
}

// ParseIntermediateParameterMap parses the given intermediate parameter map to a spec parameter map
func ParseIntermediateParameterMap(params *ParameterMap, functions specs.CustomDefinedFunctions) (*specs.ParameterMap, error) {
	if params == nil {
		return nil, nil
	}

	properties, _ := params.Properties.JustAttributes()

	header, err := ParseIntermediateHeader(params.Header, functions)
	if err != nil {
		return nil, err
	}

	result := specs.ParameterMap{
		Options:    ParseIntermediateOptions(params.Options),
		Header:     header,
		Nested:     make(map[string]*specs.NestedParameterMap, len(params.Nested)),
		Repeated:   make(map[string]*specs.RepeatedParameterMap, len(params.Repeated)),
		Properties: make(map[string]*specs.Property, len(properties)),
	}

	for _, attr := range properties {
		results, err := ParseIntermediateProperty(attr.Name, functions, attr)
		if err != nil {
			return nil, err
		}

		result.Properties[attr.Name] = results
	}

	for _, nested := range params.Nested {
		results, err := ParseIntermediateNestedParameterMap(nested, functions, nested.Name)
		if err != nil {
			return nil, err
		}

		result.Nested[nested.Name] = results
	}

	for _, repeated := range params.Repeated {
		results, err := ParseIntermediateRepeatedParameterMap(repeated, functions, repeated.Name)
		if err != nil {
			return nil, err
		}

		result.Repeated[repeated.Name] = results
	}

	return &result, nil
}

// ParseIntermediateNestedParameterMap parses the given intermediate parameter map to a spec parameter map
func ParseIntermediateNestedParameterMap(params NestedParameterMap, functions specs.CustomDefinedFunctions, path string) (*specs.NestedParameterMap, error) {
	properties, _ := params.Properties.JustAttributes()
	result := specs.NestedParameterMap{
		Name:       params.Name,
		Path:       path,
		Nested:     make(map[string]*specs.NestedParameterMap, len(params.Nested)),
		Repeated:   make(map[string]*specs.RepeatedParameterMap, len(params.Repeated)),
		Properties: make(map[string]*specs.Property, len(properties)),
	}

	for _, nested := range params.Nested {
		returns, err := ParseIntermediateNestedParameterMap(nested, functions, specs.JoinPath(path, nested.Name))
		if err != nil {
			return nil, err
		}

		result.Nested[nested.Name] = returns
	}

	for _, repeated := range params.Repeated {
		returns, err := ParseIntermediateRepeatedParameterMap(repeated, functions, specs.JoinPath(path, repeated.Name))
		if err != nil {
			return nil, err
		}

		result.Repeated[repeated.Name] = returns
	}

	for _, attr := range properties {
		returns, err := ParseIntermediateProperty(specs.JoinPath(path, attr.Name), functions, attr)
		if err != nil {
			return nil, err
		}

		result.Properties[attr.Name] = returns
	}

	return &result, nil
}

// ParseIntermediateRepeatedParameterMap parses the given intermediate repeated parameter map to a spec repeated parameter map
func ParseIntermediateRepeatedParameterMap(params RepeatedParameterMap, functions specs.CustomDefinedFunctions, path string) (*specs.RepeatedParameterMap, error) {
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
		returns, err := ParseIntermediateNestedParameterMap(nested, functions, specs.JoinPath(path, nested.Name))
		if err != nil {
			return nil, err
		}

		result.Nested[nested.Name] = returns
	}

	for _, repeated := range params.Repeated {
		returns, err := ParseIntermediateRepeatedParameterMap(repeated, functions, specs.JoinPath(path, repeated.Name))
		if err != nil {
			return nil, err
		}

		result.Repeated[repeated.Name] = returns
	}

	for _, attr := range properties {
		returns, err := ParseIntermediateProperty(specs.JoinPath(path, attr.Name), functions, attr)
		if err != nil {
			return nil, err
		}

		result.Properties[attr.Name] = returns
	}

	return &result, nil
}

// ParseIntermediateHeader parses the given intermediate header to a spec header
func ParseIntermediateHeader(header *Header, functions specs.CustomDefinedFunctions) (specs.Header, error) {
	if header == nil {
		return nil, nil
	}

	attributes, _ := header.Body.JustAttributes()
	result := make(specs.Header, len(attributes))

	for _, attr := range attributes {
		results, err := ParseIntermediateProperty(attr.Name, functions, attr)
		if err != nil {
			return nil, err
		}

		result[attr.Name] = results
	}

	return result, nil
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
func ParseIntermediateCall(call Call, functions specs.CustomDefinedFunctions) (*specs.Call, error) {
	request, err := ParseIntermediateParameterMap(call.Request, functions)
	if err != nil {
		return nil, err
	}

	rollback, err := ParseIntermediateRollbackCall(call.Rollback, functions)
	if err != nil {
		return nil, err
	}

	result := specs.Call{
		DependsOn: make(map[string]*specs.Call, len(call.DependsOn)),
		Name:      call.Name,
		Endpoint:  call.Endpoint,
		Type:      call.Type,
		Request:   request,
		Rollback:  rollback,
	}

	for _, dependency := range call.DependsOn {
		result.DependsOn[dependency] = nil
	}

	return &result, nil
}

// ParseIntermediateRollbackCall parses the given intermediate rollback call to a spec rollback call
func ParseIntermediateRollbackCall(call *RollbackCall, functions specs.CustomDefinedFunctions) (*specs.RollbackCall, error) {
	if call == nil {
		return nil, nil
	}

	results, err := ParseIntermediateParameterMap(call.Request, functions)
	if err != nil {
		return nil, err
	}

	result := specs.RollbackCall{
		Endpoint: call.Endpoint,
		Request:  results,
	}

	return &result, nil
}

// ParseIntermediateProperty parses the given intermediate property to a spec property
func ParseIntermediateProperty(path string, functions specs.CustomDefinedFunctions, property *hcl.Attribute) (*specs.Property, error) {
	if property == nil {
		return nil, nil
	}

	value, _ := property.Expr.Value(nil)
	result := &specs.Property{
		Path: path,
		Expr: property.Expr,
	}

	// Template definitions could be improved to be more consistent
	if value.Type() == cty.String && specs.IsType(value.AsString()) {
		specs.SetType(result, value)
		return result, nil
	}

	if value.Type() != cty.String || !specs.IsTemplate(value.AsString()) {
		specs.SetDefaultValue(result, value)
		return result, nil
	}

	result, err := specs.ParseTemplate(path, functions, value.AsString())
	if err != nil {
		return nil, err
	}

	return result, nil
}
