package intermediate

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/jexia/maestro/specs"
	"github.com/zclconf/go-cty/cty"
)

// ParseManifest parses the given intermediate manifest to a specs manifest
func ParseManifest(manifest Manifest, functions specs.CustomDefinedFunctions) (*specs.Manifest, error) {
	result := &specs.Manifest{
		Endpoints: make([]*specs.Endpoint, len(manifest.Endpoints)),
		Services:  make([]*specs.Service, len(manifest.Services)),
		Flows:     make([]*specs.Flow, len(manifest.Flows)),
		Proxy:     make([]*specs.Proxy, len(manifest.Proxy)),
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

// ParseIntermediateEndpoint parses the given intermediate endpoint to a specs endpoint
func ParseIntermediateEndpoint(endpoint Endpoint) *specs.Endpoint {
	result := specs.Endpoint{
		Options:  ParseIntermediateOptions(endpoint.Options),
		Flow:     endpoint.Flow,
		Listener: endpoint.Listener,
		Codec:    endpoint.Codec,
	}

	return &result
}

// ParseIntermediateService parses the given intermediate service to a specs service
func ParseIntermediateService(service Service) *specs.Service {
	result := specs.Service{
		Options: ParseIntermediateOptions(service.Options),
		Name:    service.Name,
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
		Nodes:     make([]*specs.Node, len(flow.Calls)),
		Output:    output,
	}

	for _, dependency := range flow.DependsOn {
		result.DependsOn[dependency] = nil
	}

	for index, call := range flow.Calls {
		node, err := ParseIntermediateNode(call, functions)
		if err != nil {
			return nil, err
		}

		result.Nodes[index] = node
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

// ParseIntermediateProxy parses the given intermediate proxy to a specs proxy
func ParseIntermediateProxy(proxy Proxy, functions specs.CustomDefinedFunctions) (*specs.Proxy, error) {
	forward, err := ParseIntermediateProxyForward(proxy.Forward, functions)
	if err != nil {
		return nil, err
	}

	result := specs.Proxy{
		Name:      proxy.Name,
		DependsOn: make(map[string]*specs.Flow, len(proxy.DependsOn)),
		Nodes:     make([]*specs.Node, len(proxy.Calls)),
		Forward:   forward,
	}

	for _, dependency := range proxy.DependsOn {
		result.DependsOn[dependency] = nil
	}

	for index, node := range proxy.Calls {
		node, err := ParseIntermediateNode(node, functions)
		if err != nil {
			return nil, err
		}

		result.Nodes[index] = node
	}

	return &result, nil
}

// ParseIntermediateProxyForward parses the given intermediate proxy forward to a specs proxy forward
func ParseIntermediateProxyForward(proxy ProxyForward, functions specs.CustomDefinedFunctions) (*specs.ProxyForward, error) {
	result := specs.ProxyForward{
		Service:  proxy.Service,
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
		rollback, err := ParseIntermediateCall(proxy.Rollback, functions)
		if err != nil {
			return nil, err
		}

		result.Rollback = rollback
	}

	return &result, nil
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
		Options:    make(specs.Options),
		Header:     header,
		Nested:     make(map[string]*specs.NestedParameterMap, len(params.Nested)),
		Repeated:   make(map[string]*specs.RepeatedParameterMap, len(params.Repeated)),
		Properties: make(map[string]*specs.Property, len(properties)),
	}

	if params.Options != nil {
		result.Options = ParseIntermediateOptions(params.Options.Body)
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
		Template:   specs.ParsePropertyReference(params.Template),
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
func ParseIntermediateOptions(options hcl.Body) specs.Options {
	if options == nil {
		return specs.Options{}
	}

	result := specs.Options{}
	attrs, _ := options.JustAttributes()

	for key, val := range attrs {
		val, _ := val.Expr.Value(nil)
		if val.Type() != cty.String {
			continue
		}

		result[key] = val.AsString()
	}

	return result
}

// ParseIntermediateNode parses the given intermediate call to a spec call
func ParseIntermediateNode(node Node, functions specs.CustomDefinedFunctions) (*specs.Node, error) {
	call, err := ParseIntermediateCall(node.Request, functions)
	if err != nil {
		return nil, err
	}

	rollback, err := ParseIntermediateCall(node.Rollback, functions)
	if err != nil {
		return nil, err
	}

	result := specs.Node{
		DependsOn: make(map[string]*specs.Node, len(node.DependsOn)),
		Name:      node.Name,
		Type:      node.Type,
		Call:      call,
		Rollback:  rollback,
	}

	for _, dependency := range node.DependsOn {
		result.DependsOn[dependency] = nil
	}

	return &result, nil
}

// ParseIntermediateCall parses the given intermediate call to a spec call
func ParseIntermediateCall(call *Call, functions specs.CustomDefinedFunctions) (*specs.Call, error) {
	if call == nil {
		return nil, nil
	}

	results, err := ParseIntermediateCallParameterMap(call, functions)
	if err != nil {
		return nil, err
	}

	result := specs.Call{
		Service:  call.Service,
		Endpoint: call.Endpoint,
		Request:  results,
	}

	return &result, nil
}

// ParseIntermediateCallParameterMap parses the given intermediate parameter map to a spec parameter map
func ParseIntermediateCallParameterMap(params *Call, functions specs.CustomDefinedFunctions) (*specs.ParameterMap, error) {
	if params == nil {
		return nil, nil
	}

	properties, _ := params.Properties.JustAttributes()

	header, err := ParseIntermediateHeader(params.Header, functions)
	if err != nil {
		return nil, err
	}

	result := specs.ParameterMap{
		Options:    make(specs.Options),
		Header:     header,
		Nested:     make(map[string]*specs.NestedParameterMap, len(params.Nested)),
		Repeated:   make(map[string]*specs.RepeatedParameterMap, len(params.Repeated)),
		Properties: make(map[string]*specs.Property, len(properties)),
	}

	if params.Options != nil {
		result.Options = ParseIntermediateOptions(params.Options.Body)
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

// ParseIntermediateProperty parses the given intermediate property to a spec property
func ParseIntermediateProperty(path string, functions specs.CustomDefinedFunctions, property *hcl.Attribute) (*specs.Property, error) {
	if property == nil {
		return nil, nil
	}

	value, _ := property.Expr.Value(nil)
	result := &specs.Property{
		Name: property.Name,
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
