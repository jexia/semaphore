package hcl

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/jexia/maestro/pkg/instance"
	"github.com/jexia/maestro/pkg/logger"
	"github.com/jexia/maestro/pkg/specs"
	"github.com/jexia/maestro/pkg/specs/labels"
	"github.com/jexia/maestro/pkg/specs/template"
	"github.com/jexia/maestro/pkg/specs/types"
	"github.com/zclconf/go-cty/cty"
)

// ParseFlows parses the given intermediate manifest to a flows manifest
func ParseFlows(ctx instance.Context, manifest Manifest) (*specs.FlowsManifest, error) {
	ctx.Logger(logger.Core).Info("Parsing intermediate manifest to flows manifest")

	result := &specs.FlowsManifest{
		Flows: make([]*specs.Flow, len(manifest.Flows)),
		Proxy: make([]*specs.Proxy, len(manifest.Proxy)),
	}

	for index, flow := range manifest.Flows {
		flow, err := ParseIntermediateFlow(ctx, flow)
		if err != nil {
			return nil, err
		}

		result.Flows[index] = flow
	}

	for index, proxy := range manifest.Proxy {
		proxy, err := ParseIntermediateProxy(ctx, proxy)
		if err != nil {
			return nil, err
		}

		result.Proxy[index] = proxy
	}

	return result, nil
}

// ParseEndpoints parses the given intermediate manifest to a endpoints manifest
func ParseEndpoints(ctx instance.Context, manifest Manifest) (*specs.EndpointsManifest, error) {
	ctx.Logger(logger.Core).Info("Parsing intermediate manifest to endpoints manifest")

	result := &specs.EndpointsManifest{
		Endpoints: make([]*specs.Endpoint, len(manifest.Endpoints)),
	}

	for index, endpoint := range manifest.Endpoints {
		result.Endpoints[index] = ParseIntermediateEndpoint(ctx, endpoint)
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
func ParseIntermediateFlow(ctx instance.Context, flow Flow) (*specs.Flow, error) {
	ctx.Logger(logger.Core).WithField("flow", flow.Name).Debug("Parsing intermediate flow to specs")

	input, err := ParseIntermediateInputParameterMap(ctx, flow.Input)
	if err != nil {
		return nil, err
	}

	output, err := ParseIntermediateParameterMap(ctx, flow.Output)
	if err != nil {
		return nil, err
	}

	length := len(flow.Resources)
	for _, collection := range flow.References {
		attrs, _ := collection.Properties.JustAttributes()
		length += len(attrs)
	}

	result := specs.Flow{
		Name:   flow.Name,
		Input:  input,
		Nodes:  make([]*specs.Node, 0, length),
		Output: output,
	}

	before := map[string]*specs.Node{}
	if flow.Before != nil {
		for _, resources := range flow.Before.References {
			attrs, _ := resources.Properties.JustAttributes()
			for _, attr := range attrs {
				before[attr.Name] = nil
			}

			flow.References = append([]Resources{resources}, flow.References...)
		}

		for _, node := range flow.Before.Resources {
			before[node.Name] = nil
			flow.Resources = append([]Resource{node}, flow.Resources...)
		}
	}

	for _, resources := range flow.References {
		attrs, _ := resources.Properties.JustAttributes()

		// FIXME: attrs are not always loaded in the same order as they are defined
		for _, attr := range attrs {
			prop, err := ParseIntermediateProperty(ctx, "", attr)
			if err != nil {
				return nil, err
			}

			node := &specs.Node{
				DependsOn: DependenciesExcept(before, prop.Name),
				Name:      prop.Name,
				Call: &specs.Call{
					Response: &specs.ParameterMap{
						Property: prop,
					},
				},
			}

			result.Nodes = append(result.Nodes, node)
		}
	}

	for _, intermediate := range flow.Resources {
		node, err := ParseIntermediateNode(ctx, intermediate)
		if err != nil {
			return nil, err
		}

		for key := range DependenciesExcept(before, node.Name) {
			node.DependsOn[key] = nil
		}

		result.Nodes = append(result.Nodes, node)
	}

	return &result, nil
}

// DependenciesExcept copies the given dependencies except the given resource
func DependenciesExcept(dependencies map[string]*specs.Node, resource string) map[string]*specs.Node {
	result := map[string]*specs.Node{}
	for key, val := range dependencies {
		if key == resource {
			continue
		}

		result[key] = val
	}

	return result
}

// ParseIntermediateInputParameterMap parses the given input parameter map
func ParseIntermediateInputParameterMap(ctx instance.Context, params *InputParameterMap) (*specs.ParameterMap, error) {
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
		results, err := ParseIntermediateProperty(ctx, attr.Name, attr)
		if err != nil {
			return nil, err
		}

		result.Property.Nested[attr.Name] = results
	}

	for _, nested := range params.Nested {
		results, err := ParseIntermediateNestedParameterMap(ctx, nested, nested.Name)
		if err != nil {
			return nil, err
		}

		result.Property.Nested[nested.Name] = results
	}

	for _, intermediate := range params.Repeated {
		repeated := ParseIntermediateInputRepeatedParameterMap(intermediate)
		results, err := ParseIntermediateRepeatedParameterMap(ctx, repeated, repeated.Name)
		if err != nil {
			return nil, err
		}

		result.Property.Nested[repeated.Name] = results
	}

	return result, nil
}

// ParseIntermediateProxy parses the given intermediate proxy to a specs proxy
func ParseIntermediateProxy(ctx instance.Context, proxy Proxy) (*specs.Proxy, error) {
	forward, err := ParseIntermediateProxyForward(ctx, proxy.Forward)
	if err != nil {
		return nil, err
	}

	result := specs.Proxy{
		Name:    proxy.Name,
		Nodes:   make([]*specs.Node, len(proxy.References)+len(proxy.Resources)),
		Forward: forward,
	}

	if proxy.Input != nil {
		input := &specs.ParameterMap{
			Schema: proxy.Input.Params,
			Header: make(specs.Header, len(proxy.Input.Header)),
		}

		if proxy.Input.Options != nil {
			input.Options = ParseIntermediateSpecOptions(proxy.Input.Options.Body)
		}

		for _, key := range proxy.Input.Header {
			input.Header[key] = &specs.Property{
				Path:  key,
				Name:  key,
				Type:  types.String,
				Label: labels.Optional,
			}
		}

		result.Input = input
	}

	for index, node := range proxy.Resources {
		node, err := ParseIntermediateNode(ctx, node)
		if err != nil {
			return nil, err
		}

		result.Nodes[index] = node
	}

	return &result, nil
}

// ParseIntermediateProxyForward parses the given intermediate proxy forward to a specs proxy forward
func ParseIntermediateProxyForward(ctx instance.Context, proxy ProxyForward) (*specs.Call, error) {
	result := specs.Call{
		Service: proxy.Service,
		Request: &specs.ParameterMap{},
	}

	if proxy.Header != nil {
		header, err := ParseIntermediateHeader(ctx, proxy.Header)
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
func ParseIntermediateParameterMap(ctx instance.Context, params *ParameterMap) (*specs.ParameterMap, error) {
	if params == nil {
		return nil, nil
	}

	properties, _ := params.Properties.JustAttributes()

	result := specs.ParameterMap{
		Schema:  params.Schema,
		Options: make(specs.Options),
		Property: &specs.Property{
			Type:   types.Message,
			Label:  labels.Optional,
			Nested: map[string]*specs.Property{},
		},
	}

	if params.Header != nil {
		header, err := ParseIntermediateHeader(ctx, params.Header)
		if err != nil {
			return nil, err
		}

		result.Header = header
	}

	if params.Options != nil {
		result.Options = ParseIntermediateSpecOptions(params.Options.Body)
	}

	for _, attr := range properties {
		results, err := ParseIntermediateProperty(ctx, "", attr)
		if err != nil {
			return nil, err
		}

		result.Property.Nested[attr.Name] = results
	}

	for _, nested := range params.Nested {
		results, err := ParseIntermediateNestedParameterMap(ctx, nested, nested.Name)
		if err != nil {
			return nil, err
		}

		result.Property.Nested[nested.Name] = results
	}

	for _, repeated := range params.Repeated {
		results, err := ParseIntermediateRepeatedParameterMap(ctx, repeated, repeated.Name)
		if err != nil {
			return nil, err
		}

		result.Property.Nested[repeated.Name] = results
	}

	return &result, nil
}

// ParseIntermediateNestedParameterMap parses the given intermediate parameter map to a spec parameter map
func ParseIntermediateNestedParameterMap(ctx instance.Context, params NestedParameterMap, path string) (*specs.Property, error) {
	properties, _ := params.Properties.JustAttributes()
	result := specs.Property{
		Name:   params.Name,
		Path:   path,
		Type:   types.Message,
		Label:  labels.Optional,
		Nested: map[string]*specs.Property{},
	}

	for _, nested := range params.Nested {
		returns, err := ParseIntermediateNestedParameterMap(ctx, nested, template.JoinPath(path, nested.Name))
		if err != nil {
			return nil, err
		}

		result.Nested[nested.Name] = returns
	}

	for _, repeated := range params.Repeated {
		returns, err := ParseIntermediateRepeatedParameterMap(ctx, repeated, template.JoinPath(path, repeated.Name))
		if err != nil {
			return nil, err
		}

		result.Nested[repeated.Name] = returns
	}

	for _, attr := range properties {
		returns, err := ParseIntermediateProperty(ctx, path, attr)
		if err != nil {
			return nil, err
		}

		result.Nested[attr.Name] = returns
	}

	return &result, nil
}

// ParseIntermediateRepeatedParameterMap parses the given intermediate repeated parameter map to a spec repeated parameter map
func ParseIntermediateRepeatedParameterMap(ctx instance.Context, params RepeatedParameterMap, path string) (*specs.Property, error) {
	properties, _ := params.Properties.JustAttributes()
	result := specs.Property{
		Name:      params.Name,
		Path:      path,
		Reference: template.ParsePropertyReference(params.Template),
		Type:      types.Message,
		Label:     labels.Optional,
		Nested:    map[string]*specs.Property{},
	}

	for _, nested := range params.Nested {
		returns, err := ParseIntermediateNestedParameterMap(ctx, nested, template.JoinPath(path, nested.Name))
		if err != nil {
			return nil, err
		}

		result.Nested[nested.Name] = returns
	}

	for _, repeated := range params.Repeated {
		returns, err := ParseIntermediateRepeatedParameterMap(ctx, repeated, template.JoinPath(path, repeated.Name))
		if err != nil {
			return nil, err
		}

		result.Nested[repeated.Name] = returns
	}

	for _, attr := range properties {
		returns, err := ParseIntermediateProperty(ctx, path, attr)
		if err != nil {
			return nil, err
		}

		result.Nested[attr.Name] = returns
	}

	return &result, nil
}

// ParseIntermediateHeader parses the given intermediate header to a spec header
func ParseIntermediateHeader(ctx instance.Context, header *Header) (specs.Header, error) {
	attributes, _ := header.Body.JustAttributes()
	result := make(specs.Header, len(attributes))

	for _, attr := range attributes {
		results, err := ParseIntermediateProperty(ctx, "", attr)
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

// ParseIntermediateParameters parses the given intermediate parameters
func ParseIntermediateParameters(options hcl.Body) map[string]*specs.PropertyReference {
	if options == nil {
		return map[string]*specs.PropertyReference{}
	}

	result := map[string]*specs.PropertyReference{}
	attrs, _ := options.JustAttributes()

	for key, val := range attrs {
		val, _ := val.Expr.Value(nil)
		if val.Type() != cty.String {
			continue
		}

		result[key] = template.ParsePropertyReference(template.GetTemplateContent(val.AsString()))
	}

	return result
}

// ParseIntermediateNode parses the given intermediate call to a spec call
func ParseIntermediateNode(ctx instance.Context, node Resource) (*specs.Node, error) {
	call, err := ParseIntermediateCall(ctx, node.Request)
	if err != nil {
		return nil, err
	}

	rollback, err := ParseIntermediateCall(ctx, node.Rollback)
	if err != nil {
		return nil, err
	}

	result := specs.Node{
		DependsOn: make(map[string]*specs.Node, len(node.DependsOn)),
		Name:      node.Name,
		Call:      call,
		Rollback:  rollback,
	}

	for _, dependency := range node.DependsOn {
		result.DependsOn[dependency] = nil
	}

	return &result, nil
}

// ParseIntermediateCall parses the given intermediate call to a spec call
func ParseIntermediateCall(ctx instance.Context, call *Call) (*specs.Call, error) {
	if call == nil {
		return nil, nil
	}

	results, err := ParseIntermediateCallParameterMap(ctx, call)
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
func ParseIntermediateCallParameterMap(ctx instance.Context, params *Call) (*specs.ParameterMap, error) {
	if params == nil {
		return nil, nil
	}

	properties, _ := params.Properties.JustAttributes()

	result := specs.ParameterMap{
		Options: make(specs.Options),
		Property: &specs.Property{
			Type:   types.Message,
			Label:  labels.Optional,
			Nested: map[string]*specs.Property{},
		},
	}

	if params.Parameters != nil {
		result.Params = ParseIntermediateParameters(params.Parameters.Body)
	}

	if params.Header != nil {
		header, err := ParseIntermediateHeader(ctx, params.Header)
		if err != nil {
			return nil, err
		}

		result.Header = header
	}

	if params.Options != nil {
		result.Options = ParseIntermediateSpecOptions(params.Options.Body)
	}

	for _, attr := range properties {
		results, err := ParseIntermediateProperty(ctx, "", attr)
		if err != nil {
			return nil, err
		}

		result.Property.Nested[attr.Name] = results
	}

	for _, nested := range params.Nested {
		results, err := ParseIntermediateNestedParameterMap(ctx, nested, nested.Name)
		if err != nil {
			return nil, err
		}

		result.Property.Nested[nested.Name] = results
	}

	for _, repeated := range params.Repeated {
		results, err := ParseIntermediateRepeatedParameterMap(ctx, repeated, repeated.Name)
		if err != nil {
			return nil, err
		}

		result.Property.Nested[repeated.Name] = results
	}

	return &result, nil
}

// ParseIntermediateProperty parses the given intermediate property to a spec property
func ParseIntermediateProperty(ctx instance.Context, path string, property *hcl.Attribute) (*specs.Property, error) {
	if property == nil {
		return nil, nil
	}

	ctx.Logger(logger.Core).WithField("path", path).Debug("Parsing intermediate property to specs")

	value, _ := property.Expr.Value(nil)
	result := &specs.Property{
		Name:  property.Name,
		Path:  template.JoinPath(path, property.Name),
		Expr:  property.Expr,
		Label: labels.Optional,
	}

	if value.Type() != cty.String || !template.Is(value.AsString()) {
		SetDefaultValue(ctx, result, value)
		return result, nil
	}

	result, err := template.Parse(ctx, path, property.Name, value.AsString())
	if err != nil {
		return nil, err
	}

	result.Name = property.Name
	return result, nil
}
