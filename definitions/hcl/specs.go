package hcl

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/jexia/maestro/instance"
	"github.com/jexia/maestro/logger"
	"github.com/jexia/maestro/specs"
	"github.com/jexia/maestro/specs/labels"
	"github.com/jexia/maestro/specs/types"
	"github.com/zclconf/go-cty/cty"
)

// ParseSpecs parses the given intermediate manifest to a specs manifest
func ParseSpecs(ctx instance.Context, manifest Manifest, functions specs.CustomDefinedFunctions) (*specs.Manifest, error) {
	ctx.Logger(logger.Core).Info("Parsing intermediate manifest to specs")

	result := &specs.Manifest{
		Endpoints: make([]*specs.Endpoint, len(manifest.Endpoints)),
		Flows:     make([]*specs.Flow, len(manifest.Flows)),
		Proxy:     make([]*specs.Proxy, len(manifest.Proxy)),
	}

	for index, endpoint := range manifest.Endpoints {
		result.Endpoints[index] = ParseIntermediateEndpoint(ctx, endpoint)
	}

	for index, flow := range manifest.Flows {
		flow, err := ParseIntermediateFlow(ctx, flow, functions)
		if err != nil {
			return nil, err
		}

		result.Flows[index] = flow
	}

	for index, proxy := range manifest.Proxy {
		proxy, err := ParseIntermediateProxy(ctx, proxy, functions)
		if err != nil {
			return nil, err
		}

		result.Proxy[index] = proxy
	}

	return result, nil
}

// ParseIntermediateEndpoint parses the given intermediate endpoint to a specs endpoint
func ParseIntermediateEndpoint(ctx instance.Context, endpoint Endpoint) *specs.Endpoint {
	ctx.Logger(logger.Core).WithField("flow", endpoint.Flow).Debug("Parsing intermediate endpoint to specs")

	result := specs.Endpoint{
		Options:  ParseIntermediateSpecOptions(endpoint.Options),
		Flow:     endpoint.Flow,
		Listener: endpoint.Listener,
	}

	return &result
}

// ParseIntermediateFlow parses the given intermediate flow to a specs flow
func ParseIntermediateFlow(ctx instance.Context, flow Flow, functions specs.CustomDefinedFunctions) (*specs.Flow, error) {
	ctx.Logger(logger.Core).WithField("flow", flow.Name).Debug("Parsing intermediate flow to specs")

	input, err := ParseIntermediateInputParameterMap(ctx, flow.Input, functions)
	if err != nil {
		return nil, err
	}

	output, err := ParseIntermediateParameterMap(ctx, flow.Output, functions)
	if err != nil {
		return nil, err
	}

	result := specs.Flow{
		Name:      flow.Name,
		DependsOn: make(map[string]*specs.Flow, len(flow.DependsOn)),
		Input:     input,
		Nodes:     make([]*specs.Node, len(flow.Resources)),
		Output:    output,
	}

	for _, dependency := range flow.DependsOn {
		result.DependsOn[dependency] = nil
	}

	for index, call := range flow.Resources {
		node, err := ParseIntermediateNode(ctx, call, functions)
		if err != nil {
			return nil, err
		}

		result.Nodes[index] = node
	}

	return &result, nil
}

// ParseIntermediateInputParameterMap parses the given input parameter map
func ParseIntermediateInputParameterMap(ctx instance.Context, params *InputParameterMap, functions specs.CustomDefinedFunctions) (*specs.ParameterMap, error) {
	if params == nil {
		return nil, nil
	}

	properties, _ := params.Properties.JustAttributes()
	result := &specs.ParameterMap{
		Schema:  params.Schema,
		Options: make(specs.Options),
		Header:  make(specs.Header, len(params.Header)),
		Property: &specs.Property{
			Type:   types.Message,
			Label:  labels.Optional,
			Nested: map[string]*specs.Property{},
		},
	}

	for _, key := range params.Header {
		result.Header[key] = &specs.Property{
			Path:  key,
			Name:  key,
			Type:  types.String,
			Label: labels.Optional,
		}
	}

	if params.Options != nil {
		result.Options = ParseIntermediateSpecOptions(params.Options.Body)
	}

	for _, attr := range properties {
		results, err := ParseIntermediateProperty(ctx, attr.Name, functions, attr)
		if err != nil {
			return nil, err
		}

		result.Property.Nested[attr.Name] = results
	}

	for _, nested := range params.Nested {
		results, err := ParseIntermediateNestedParameterMap(ctx, nested, functions, nested.Name)
		if err != nil {
			return nil, err
		}

		result.Property.Nested[nested.Name] = results
	}

	for _, intermediate := range params.Repeated {
		repeated := ParseIntermediateInputRepeatedParameterMap(intermediate)
		results, err := ParseIntermediateRepeatedParameterMap(ctx, repeated, functions, repeated.Name)
		if err != nil {
			return nil, err
		}

		result.Property.Nested[repeated.Name] = results
	}

	return result, nil
}

// ParseIntermediateProxy parses the given intermediate proxy to a specs proxy
func ParseIntermediateProxy(ctx instance.Context, proxy Proxy, functions specs.CustomDefinedFunctions) (*specs.Proxy, error) {
	forward, err := ParseIntermediateProxyForward(ctx, proxy.Forward, functions)
	if err != nil {
		return nil, err
	}

	result := specs.Proxy{
		Name:      proxy.Name,
		DependsOn: make(map[string]*specs.Flow, len(proxy.DependsOn)),
		Nodes:     make([]*specs.Node, len(proxy.Resources)),
		Forward:   forward,
	}

	for _, dependency := range proxy.DependsOn {
		result.DependsOn[dependency] = nil
	}

	for index, node := range proxy.Resources {
		node, err := ParseIntermediateNode(ctx, node, functions)
		if err != nil {
			return nil, err
		}

		result.Nodes[index] = node
	}

	return &result, nil
}

// ParseIntermediateProxyForward parses the given intermediate proxy forward to a specs proxy forward
func ParseIntermediateProxyForward(ctx instance.Context, proxy ProxyForward, functions specs.CustomDefinedFunctions) (*specs.Call, error) {
	result := specs.Call{
		Service: proxy.Service,
		Request: &specs.ParameterMap{},
	}

	if proxy.Header != nil {
		header, err := ParseIntermediateHeader(ctx, proxy.Header, functions)
		if err != nil {
			return nil, err
		}

		result.Request.Header = header
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
func ParseIntermediateParameterMap(ctx instance.Context, params *ParameterMap, functions specs.CustomDefinedFunctions) (*specs.ParameterMap, error) {
	if params == nil {
		return nil, nil
	}

	properties, _ := params.Properties.JustAttributes()
	header, err := ParseIntermediateHeader(ctx, params.Header, functions)
	if err != nil {
		return nil, err
	}

	result := specs.ParameterMap{
		Schema:  params.Schema,
		Options: make(specs.Options),
		Header:  header,
		Property: &specs.Property{
			Type:   types.Message,
			Label:  labels.Optional,
			Nested: map[string]*specs.Property{},
		},
	}

	if params.Options != nil {
		result.Options = ParseIntermediateSpecOptions(params.Options.Body)
	}

	for _, attr := range properties {
		results, err := ParseIntermediateProperty(ctx, attr.Name, functions, attr)
		if err != nil {
			return nil, err
		}

		result.Property.Nested[attr.Name] = results
	}

	for _, nested := range params.Nested {
		results, err := ParseIntermediateNestedParameterMap(ctx, nested, functions, nested.Name)
		if err != nil {
			return nil, err
		}

		result.Property.Nested[nested.Name] = results
	}

	for _, repeated := range params.Repeated {
		results, err := ParseIntermediateRepeatedParameterMap(ctx, repeated, functions, repeated.Name)
		if err != nil {
			return nil, err
		}

		result.Property.Nested[repeated.Name] = results
	}

	return &result, nil
}

// ParseIntermediateNestedParameterMap parses the given intermediate parameter map to a spec parameter map
func ParseIntermediateNestedParameterMap(ctx instance.Context, params NestedParameterMap, functions specs.CustomDefinedFunctions, path string) (*specs.Property, error) {
	properties, _ := params.Properties.JustAttributes()
	result := specs.Property{
		Name:   params.Name,
		Path:   path,
		Type:   types.Message,
		Label:  labels.Optional,
		Nested: map[string]*specs.Property{},
	}

	for _, nested := range params.Nested {
		returns, err := ParseIntermediateNestedParameterMap(ctx, nested, functions, specs.JoinPath(path, nested.Name))
		if err != nil {
			return nil, err
		}

		result.Nested[nested.Name] = returns
	}

	for _, repeated := range params.Repeated {
		returns, err := ParseIntermediateRepeatedParameterMap(ctx, repeated, functions, specs.JoinPath(path, repeated.Name))
		if err != nil {
			return nil, err
		}

		result.Nested[repeated.Name] = returns
	}

	for _, attr := range properties {
		returns, err := ParseIntermediateProperty(ctx, specs.JoinPath(path, attr.Name), functions, attr)
		if err != nil {
			return nil, err
		}

		result.Nested[attr.Name] = returns
	}

	return &result, nil
}

// ParseIntermediateRepeatedParameterMap parses the given intermediate repeated parameter map to a spec repeated parameter map
func ParseIntermediateRepeatedParameterMap(ctx instance.Context, params RepeatedParameterMap, functions specs.CustomDefinedFunctions, path string) (*specs.Property, error) {
	properties, _ := params.Properties.JustAttributes()
	result := specs.Property{
		Name:      params.Name,
		Path:      path,
		Reference: specs.ParsePropertyReference(params.Template),
		Type:      types.Message,
		Label:     labels.Optional,
		Nested:    map[string]*specs.Property{},
	}

	for _, nested := range params.Nested {
		returns, err := ParseIntermediateNestedParameterMap(ctx, nested, functions, specs.JoinPath(path, nested.Name))
		if err != nil {
			return nil, err
		}

		result.Nested[nested.Name] = returns
	}

	for _, repeated := range params.Repeated {
		returns, err := ParseIntermediateRepeatedParameterMap(ctx, repeated, functions, specs.JoinPath(path, repeated.Name))
		if err != nil {
			return nil, err
		}

		result.Nested[repeated.Name] = returns
	}

	for _, attr := range properties {
		returns, err := ParseIntermediateProperty(ctx, specs.JoinPath(path, attr.Name), functions, attr)
		if err != nil {
			return nil, err
		}

		result.Nested[attr.Name] = returns
	}

	return &result, nil
}

// ParseIntermediateHeader parses the given intermediate header to a spec header
func ParseIntermediateHeader(ctx instance.Context, header *Header, functions specs.CustomDefinedFunctions) (specs.Header, error) {
	if header == nil {
		return nil, nil
	}

	attributes, _ := header.Body.JustAttributes()
	result := make(specs.Header, len(attributes))

	for _, attr := range attributes {
		results, err := ParseIntermediateProperty(ctx, attr.Name, functions, attr)
		if err != nil {
			return nil, err
		}

		result[attr.Name] = results
	}

	return result, nil
}

// ParseIntermediateSpecOptions parses the given intermediate options to a spec options
func ParseIntermediateSpecOptions(options hcl.Body) specs.Options {
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
func ParseIntermediateNode(ctx instance.Context, node Node, functions specs.CustomDefinedFunctions) (*specs.Node, error) {
	call, err := ParseIntermediateCall(ctx, node.Request, functions)
	if err != nil {
		return nil, err
	}

	rollback, err := ParseIntermediateCall(ctx, node.Rollback, functions)
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
func ParseIntermediateCall(ctx instance.Context, call *Call, functions specs.CustomDefinedFunctions) (*specs.Call, error) {
	if call == nil {
		return nil, nil
	}

	results, err := ParseIntermediateCallParameterMap(ctx, call, functions)
	if err != nil {
		return nil, err
	}

	result := specs.Call{
		Service: call.Service,
		Method:  call.Method,
		Request: results,
	}

	return &result, nil
}

// ParseIntermediateCallParameterMap parses the given intermediate parameter map to a spec parameter map
func ParseIntermediateCallParameterMap(ctx instance.Context, params *Call, functions specs.CustomDefinedFunctions) (*specs.ParameterMap, error) {
	if params == nil {
		return nil, nil
	}

	properties, _ := params.Properties.JustAttributes()

	header, err := ParseIntermediateHeader(ctx, params.Header, functions)
	if err != nil {
		return nil, err
	}

	result := specs.ParameterMap{
		Options: make(specs.Options),
		Header:  header,
		Property: &specs.Property{
			Type:   types.Message,
			Label:  labels.Optional,
			Nested: map[string]*specs.Property{},
		},
	}

	if params.Options != nil {
		result.Options = ParseIntermediateSpecOptions(params.Options.Body)
	}

	for _, attr := range properties {
		results, err := ParseIntermediateProperty(ctx, attr.Name, functions, attr)
		if err != nil {
			return nil, err
		}

		result.Property.Nested[attr.Name] = results
	}

	for _, nested := range params.Nested {
		results, err := ParseIntermediateNestedParameterMap(ctx, nested, functions, nested.Name)
		if err != nil {
			return nil, err
		}

		result.Property.Nested[nested.Name] = results
	}

	for _, repeated := range params.Repeated {
		results, err := ParseIntermediateRepeatedParameterMap(ctx, repeated, functions, repeated.Name)
		if err != nil {
			return nil, err
		}

		result.Property.Nested[repeated.Name] = results
	}

	return &result, nil
}

// ParseIntermediateProperty parses the given intermediate property to a spec property
func ParseIntermediateProperty(ctx instance.Context, path string, functions specs.CustomDefinedFunctions, property *hcl.Attribute) (*specs.Property, error) {
	if property == nil {
		return nil, nil
	}

	ctx.Logger(logger.Core).WithField("path", path).Debug("Parsing intermediate property to specs")

	value, _ := property.Expr.Value(nil)
	result := &specs.Property{
		Name: property.Name,
		Path: path,
		Expr: property.Expr,
	}

	if value.Type() != cty.String || !specs.IsTemplate(value.AsString()) {
		specs.SetDefaultValue(ctx, result, value)
		return result, nil
	}

	result, err := specs.ParseTemplate(ctx, path, functions, value.AsString())
	if err != nil {
		return nil, err
	}

	result.Name = property.Name
	return result, nil
}
