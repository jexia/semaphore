package hcl

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/jexia/semaphore/pkg/broker"
	"github.com/jexia/semaphore/pkg/broker/logger"
	"github.com/jexia/semaphore/pkg/conditions"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/labels"
	"github.com/jexia/semaphore/pkg/specs/template"
	"github.com/jexia/semaphore/pkg/specs/types"
	"github.com/zclconf/go-cty/cty"
	"go.uber.org/zap"
)

/**
 * TODO: this file has to be refactored to avoid code duplication and
 * type casting. A lot of code duplication is aroundnested and repeated parameter
 * maps. Interfaces have to be created for these types to allow to reuse implementations.
 */

// ParseFlows parses the given intermediate manifest to a flows manifest
func ParseFlows(ctx *broker.Context, manifest Manifest) (errObject *specs.ParameterMap, _ specs.FlowListInterface, _ error) {
	logger.Info(ctx, "parsing intermediate manifest to flows manifest")

	result := make(specs.FlowListInterface, 0, len(manifest.Flows)+len(manifest.Proxy))

	if manifest.Error != nil {
		spec, err := ParseIntermediateParameterMap(ctx, manifest.Error)
		if err != nil {
			return errObject, nil, err
		}

		errObject = spec
	}

	for _, flow := range manifest.Flows {
		flow, err := ParseIntermediateFlow(logger.WithFields(ctx, zap.String("flow", flow.Name)), flow)
		if err != nil {
			return errObject, nil, err
		}

		result.Append(specs.FlowListInterface{flow})
	}

	for _, proxy := range manifest.Proxy {
		proxy, err := ParseIntermediateProxy(ctx, proxy)
		if err != nil {
			return errObject, nil, err
		}

		result.Append(specs.FlowListInterface{proxy})
	}

	return errObject, result, nil
}

// ParseEndpoints parses the given intermediate manifest to a endpoints manifest
func ParseEndpoints(ctx *broker.Context, manifest Manifest) (specs.EndpointList, error) {
	logger.Info(ctx, "parsing intermediate manifest to endpoints manifest")

	result := make(specs.EndpointList, len(manifest.Endpoints))
	for index, endpoint := range manifest.Endpoints {
		result[index] = ParseIntermediateEndpoint(logger.WithFields(ctx, zap.String("flow", endpoint.Flow)), endpoint)
	}

	return result, nil
}

// ParseIntermediateEndpoint parses the given intermediate endpoint to a specs endpoint
func ParseIntermediateEndpoint(ctx *broker.Context, endpoint Endpoint) *specs.Endpoint {
	logger.Debug(ctx, "parsing intermediate endpoint to specs")

	result := specs.Endpoint{
		Options:  ParseIntermediateSpecOptions(endpoint.Options),
		Flow:     endpoint.Flow,
		Listener: endpoint.Listener,
	}

	return &result
}

// ParseIntermediateFlow parses the given intermediate flow to a specs flow
func ParseIntermediateFlow(ctx *broker.Context, flow Flow) (*specs.Flow, error) {
	logger.Info(ctx, "parsing intermediate flow to specs")

	input, err := ParseIntermediateInputParameterMap(ctx, flow.Input)
	if err != nil {
		return nil, err
	}

	output, err := ParseIntermediateParameterMap(ctx, flow.Output)
	if err != nil {
		return nil, err
	}

	result := specs.Flow{
		Name:    flow.Name,
		Input:   input,
		Nodes:   specs.NodeList{},
		Output:  output,
		OnError: &specs.OnError{},
	}

	if flow.OnError != nil {
		spec, err := ParseIntermediateOnError(ctx, flow.OnError)
		if err != nil {
			return nil, err
		}

		result.OnError = spec
	}

	if flow.Error != nil {
		spec, err := ParseIntermediateParameterMap(ctx, flow.Error)
		if err != nil {
			return nil, err
		}

		result.OnError.Response = spec
	}

	var before specs.Dependencies
	if flow.Before != nil {
		dependencies, references, resources := ParseIntermediateBefore(ctx, flow.Before)
		before = dependencies

		flow.References = append(flow.References, references...)
		flow.Resources = append(flow.Resources, resources...)
	}

	for _, references := range flow.References {
		nodes, err := ParseIntermediateResources(ctx, before, references)
		if err != nil {
			return nil, err
		}

		result.Nodes = append(result.Nodes, nodes...)
	}

	for _, intermediate := range flow.Resources {
		node, err := ParseIntermediateNode(ctx, before, intermediate)
		if err != nil {
			return nil, err
		}

		result.Nodes = append(result.Nodes, node)
	}

	for _, condition := range flow.Conditions {
		nodes, err := ParseIntermediateCondition(ctx, before, condition)
		if err != nil {
			return nil, err
		}

		result.Nodes = append(result.Nodes, nodes...)
	}

	return &result, nil
}

// DependenciesExcept copies the given dependencies except the given resource
func DependenciesExcept(dependencies specs.Dependencies, resource string) specs.Dependencies {
	result := specs.Dependencies{}
	for key, val := range dependencies {
		if key == resource {
			continue
		}

		result[key] = val
	}

	return result
}

// DependsOn sets the given dependency for the given nodes
func DependsOn(dependency *specs.Node, nodes ...*specs.Node) {
	for _, node := range nodes {
		node.DependsOn[dependency.ID] = dependency
	}
}

// ParseIntermediateInputParameterMap parses the given input parameter map
func ParseIntermediateInputParameterMap(ctx *broker.Context, params *InputParameterMap) (*specs.ParameterMap, error) {
	if params == nil {
		return nil, nil
	}

	result := &specs.ParameterMap{
		Schema:  params.Schema,
		Options: make(specs.Options),
		Header:  make(specs.Header, len(params.Header)),
	}

	for _, key := range params.Header {
		result.Header[key] = &specs.Property{
			Path:  key,
			Name:  key,
			Label: labels.Optional,
			Template: specs.Template{
				Scalar: &specs.Scalar{
					Type: types.String,
				},
			},
		}
	}

	if params.Options != nil {
		result.Options = ParseIntermediateSpecOptions(params.Options.Body)
	}

	return result, nil
}

// ParseIntermediateProxy parses the given intermediate proxy to a specs proxy
func ParseIntermediateProxy(ctx *broker.Context, proxy Proxy) (*specs.Proxy, error) {
	forward, err := ParseIntermediateProxyForward(ctx, proxy.Forward)
	if err != nil {
		return nil, err
	}

	result := specs.Proxy{
		Name:    proxy.Name,
		Nodes:   specs.NodeList{},
		Forward: forward,
		OnError: &specs.OnError{},
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
				Path: key,
				Name: key,
				Template: specs.Template{
					Scalar: &specs.Scalar{
						Type: types.String,
					},
				},
			}
		}

		result.Input = input
	}

	if proxy.OnError != nil {
		spec, err := ParseIntermediateOnError(ctx, proxy.OnError)
		if err != nil {
			return nil, err
		}

		result.OnError = spec
	}

	if proxy.Error != nil {
		spec, err := ParseIntermediateParameterMap(ctx, proxy.Error)
		if err != nil {
			return nil, err
		}

		result.OnError.Response = spec
	}

	var before specs.Dependencies
	if proxy.Before != nil {
		dependencies, references, resources := ParseIntermediateBefore(ctx, proxy.Before)
		before = dependencies

		proxy.References = append(proxy.References, references...)
		proxy.Resources = append(proxy.Resources, resources...)
	}

	for _, references := range proxy.References {
		nodes, err := ParseIntermediateResources(ctx, before, references)
		if err != nil {
			return nil, err
		}

		result.Nodes = append(result.Nodes, nodes...)
	}

	for _, node := range proxy.Resources {
		node, err := ParseIntermediateNode(ctx, before, node)
		if err != nil {
			return nil, err
		}

		result.Nodes = append(result.Nodes, node)
	}

	for _, condition := range proxy.Conditions {
		nodes, err := ParseIntermediateCondition(ctx, before, condition)
		if err != nil {
			return nil, err
		}

		result.Nodes = append(result.Nodes, nodes...)
	}

	return &result, nil
}

// ParseIntermediateProxyForward parses the given intermediate proxy forward to a specs proxy forward
func ParseIntermediateProxyForward(ctx *broker.Context, proxy ProxyForward) (*specs.Call, error) {
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

// ParseIntermediateParameterMap parses the given intermediate parameter map to a spec parameter map
func ParseIntermediateParameterMap(ctx *broker.Context, params *ParameterMap) (*specs.ParameterMap, error) {
	if params == nil {
		return nil, nil
	}

	properties, _ := params.Properties.JustAttributes()

	result := specs.ParameterMap{
		Schema:  params.Schema,
		Options: make(specs.Options),
		Property: &specs.Property{
			Label: labels.Optional,
			Template: specs.Template{
				Message: make(specs.Message),
			},
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
		val, _ := attr.Expr.Value(nil)
		results, err := ParseIntermediateProperty(ctx, "", attr, val)
		if err != nil {
			return nil, err
		}

		result.Property.Message[attr.Name] = results
	}

	for _, nested := range params.Nested {
		results, err := ParseIntermediateNestedParameterMap(ctx, nested, nested.Name)
		if err != nil {
			return nil, err
		}

		result.Property.Message[nested.Name] = results
	}

	for _, repeated := range params.Repeated {
		results, err := ParseIntermediateRepeatedParameterMap(ctx, repeated, repeated.Name)
		if err != nil {
			return nil, err
		}

		result.Property.Message[repeated.Name] = results
	}

	return &result, nil
}

// ParseIntermediateNestedParameterMap parses the given intermediate parameter map to a spec parameter map
func ParseIntermediateNestedParameterMap(ctx *broker.Context, params NestedParameterMap, path string) (*specs.Property, error) {
	properties, _ := params.Properties.JustAttributes()
	result := specs.Property{
		Name:  params.Name,
		Path:  path,
		Label: labels.Optional,
		Template: specs.Template{
			Message: make(specs.Message),
		},
	}

	for _, nested := range params.Nested {
		returns, err := ParseIntermediateNestedParameterMap(ctx, nested, template.JoinPath(path, nested.Name))
		if err != nil {
			return nil, err
		}

		result.Message[nested.Name] = returns
	}

	for _, repeated := range params.Repeated {
		returns, err := ParseIntermediateRepeatedParameterMap(ctx, repeated, template.JoinPath(path, repeated.Name))
		if err != nil {
			return nil, err
		}

		result.Message[repeated.Name] = returns
	}

	for _, attr := range properties {
		val, _ := attr.Expr.Value(nil)
		returns, err := ParseIntermediateProperty(ctx, path, attr, val)
		if err != nil {
			return nil, err
		}

		result.Message[attr.Name] = returns
	}

	return &result, nil
}

// ParseIntermediateRepeatedParameterMap parses the given intermediate repeated parameter map to a spec repeated parameter map
func ParseIntermediateRepeatedParameterMap(ctx *broker.Context, params RepeatedParameterMap, path string) (*specs.Property, error) {
	properties, _ := params.Properties.JustAttributes()
	msg := specs.Template{
		Message: make(specs.Message),
	}

	for _, nested := range params.Nested {
		returns, err := ParseIntermediateNestedParameterMap(ctx, nested, template.JoinPath(path, nested.Name))
		if err != nil {
			return nil, err
		}

		msg.Message[nested.Name] = returns
	}

	for _, repeated := range params.Repeated {
		returns, err := ParseIntermediateRepeatedParameterMap(ctx, repeated, template.JoinPath(path, repeated.Name))
		if err != nil {
			return nil, err
		}

		msg.Message[repeated.Name] = returns
	}

	for _, attr := range properties {
		val, _ := attr.Expr.Value(nil)
		returns, err := ParseIntermediateProperty(ctx, path, attr, val)
		if err != nil {
			return nil, err
		}

		msg.Message[attr.Name] = returns
	}

	result := specs.Property{
		Name:  params.Name,
		Path:  path,
		Label: labels.Optional,
		Template: specs.Template{
			Reference: template.ParsePropertyReference(params.Template),
			Repeated:  specs.Repeated{msg},
		},
	}

	return &result, nil
}

// ParseIntermediateHeader parses the given intermediate header to a spec header
func ParseIntermediateHeader(ctx *broker.Context, header *Header) (specs.Header, error) {
	attributes, _ := header.Body.JustAttributes()
	result := make(specs.Header, len(attributes))

	for _, attr := range attributes {
		val, _ := attr.Expr.Value(nil)
		results, err := ParseIntermediateProperty(ctx, "", attr, val)
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
func ParseIntermediateParameters(ctx *broker.Context, options hcl.Body) (map[string]*specs.Property, error) {
	if options == nil {
		return map[string]*specs.Property{}, nil
	}

	result := map[string]*specs.Property{}
	attrs, _ := options.JustAttributes()

	for key, attr := range attrs {
		val, _ := attr.Expr.Value(nil)
		property, err := ParseIntermediateProperty(ctx, key, attr, val)
		if err != nil {
			return nil, err
		}

		result[key] = property
	}

	return result, nil
}

// ParseIntermediateNode parses the given intermediate call to a spec call
func ParseIntermediateNode(ctx *broker.Context, dependencies specs.Dependencies, node Resource) (*specs.Node, error) {
	call, err := ParseIntermediateCall(ctx, node.Request)
	if err != nil {
		return nil, err
	}

	rollback, err := ParseIntermediateCall(ctx, node.Rollback)
	if err != nil {
		return nil, err
	}

	result := specs.Node{
		Type:         specs.NodeIntermediate,
		DependsOn:    make(specs.Dependencies, len(node.DependsOn)),
		ID:           node.Name,
		Name:         node.Name,
		Call:         call,
		ExpectStatus: node.ExpectStatus,
		Rollback:     rollback,
		OnError:      &specs.OnError{},
	}

	if node.Request != nil {
		result.Type = specs.NodeCall
	}

	if node.OnError != nil {
		spec, err := ParseIntermediateOnError(ctx, node.OnError)
		if err != nil {
			return nil, err
		}

		result.OnError = spec
	}

	if node.Error != nil {
		spec, err := ParseIntermediateParameterMap(ctx, node.Error)
		if err != nil {
			return nil, err
		}

		result.OnError.Response = spec
	}

	for _, dependency := range node.DependsOn {
		result.DependsOn[dependency] = nil
	}

	for key := range DependenciesExcept(dependencies, node.Name) {
		result.DependsOn[key] = nil
	}

	return &result, nil
}

// ParseIntermediateCall parses the given intermediate call to a spec call
func ParseIntermediateCall(ctx *broker.Context, call *Call) (*specs.Call, error) {
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
func ParseIntermediateCallParameterMap(ctx *broker.Context, params *Call) (*specs.ParameterMap, error) {
	if params == nil {
		return nil, nil
	}

	result := specs.ParameterMap{
		Options: make(specs.Options),
		Property: &specs.Property{
			Label: labels.Optional,
			Template: specs.Template{
				Message: make(specs.Message),
			},
		},
	}

	if params.Parameters != nil {
		params, err := ParseIntermediateParameters(ctx, params.Parameters.Body)
		if err != nil {
			return nil, err
		}

		result.Params = params
	}

	if params.Options != nil {
		result.Options = ParseIntermediateSpecOptions(params.Options.Body)
	}

	properties, _ := params.Properties.JustAttributes()

	if params.Header != nil {
		header, err := ParseIntermediateHeader(ctx, params.Header)
		if err != nil {
			return nil, err
		}

		result.Header = header
	}

	for _, attr := range properties {
		val, _ := attr.Expr.Value(nil)
		results, err := ParseIntermediateProperty(ctx, "", attr, val)
		if err != nil {
			return nil, err
		}

		result.Property.Message[results.Name] = results
	}

	for _, nested := range params.Nested {
		results, err := ParseIntermediateNestedParameterMap(ctx, nested, nested.Name)
		if err != nil {
			return nil, err
		}

		result.Property.Message[results.Name] = results
	}

	for _, repeated := range params.Repeated {
		results, err := ParseIntermediateRepeatedParameterMap(ctx, repeated, repeated.Name)
		if err != nil {
			return nil, err
		}

		result.Property.Message[repeated.Name] = results
	}

	return &result, nil
}

// ParseIntermediateProperty parses the given intermediate property to a spec property
func ParseIntermediateProperty(ctx *broker.Context, path string, property *hcl.Attribute, value cty.Value) (*specs.Property, error) {
	if property == nil {
		return nil, nil
	}

	logger.Debug(ctx, "parsing intermediate property to specs", zap.String("path", path))

	fqpath := template.JoinPath(path, property.Name)
	typed := value.Type()
	result := &specs.Property{
		Name:  property.Name,
		Path:  fqpath,
		Expr:  &Expression{property.Expr},
		Label: labels.Optional,
	}

	switch {
	case typed.IsTupleType():
		result.Repeated = make(specs.Repeated, 0, typed.Length())

		value.ForEachElement(func(key cty.Value, value cty.Value) (stop bool) {
			attr := &hcl.Attribute{
				Expr: property.Expr,
			}

			// TODO: refactor to return template
			item, err := ParseIntermediateProperty(ctx, fqpath, attr, value)
			if err != nil {
				return true
			}

			result.Repeated = append(result.Repeated, item.Template)
			return false
		})

		break
	case typed.IsObjectType():
		result.Message = specs.Message{}

		value.ForEachElement(func(key cty.Value, value cty.Value) (stop bool) {
			attr := &hcl.Attribute{
				Name: key.AsString(),
				Expr: property.Expr,
			}

			item, err := ParseIntermediateProperty(ctx, fqpath, attr, value)
			if err != nil {
				return true
			}

			result.Message[key.AsString()] = item
			return false
		})

		break
	case typed == cty.String && template.Is(value.AsString()):
		// TODO: refactor and return template
		returns, err := template.Parse(ctx, path, result.Name, value.AsString())
		if err != nil {
			return nil, err
		}

		result = returns
		break
	default:
		err := SetScalar(ctx, &result.Template, value)
		if err != nil {
			return nil, err
		}

		break
	}

	return result, nil
}

// ParseIntermediateBefore parses the given before into a collection of dependencies
func ParseIntermediateBefore(ctx *broker.Context, before *Before) (dependencies specs.Dependencies, references []Resources, resources []Resource) {
	result := make(specs.Dependencies)

	for _, resources := range before.References {
		attrs, _ := resources.Properties.JustAttributes()
		for _, attr := range attrs {
			result[attr.Name] = nil
		}

		references = append([]Resources{resources}, references...)
	}

	for _, node := range before.Resources {
		result[node.Name] = nil
		resources = append([]Resource{node}, resources...)
	}

	return result, references, resources
}

// ParseIntermediateResources parses the given resources to nodes
func ParseIntermediateResources(ctx *broker.Context, dependencies specs.Dependencies, resources Resources) (specs.NodeList, error) {
	attrs, _ := resources.Properties.JustAttributes()
	nodes := make(specs.NodeList, 0, len(attrs))

	for _, attr := range attrs {
		val, _ := attr.Expr.Value(nil)
		prop, err := ParseIntermediateProperty(ctx, "", attr, val)
		if err != nil {
			return nil, err
		}

		node := &specs.Node{
			Type:      specs.NodeIntermediate,
			DependsOn: DependenciesExcept(dependencies, prop.Name),
			ID:        prop.Name,
			Intermediate: &specs.ParameterMap{
				Property: prop,
			},
		}

		nodes = append(nodes, node)
	}

	return nodes, nil
}

// ParseIntermediateCondition parses the given intermediate condition and returns the compiled nodes
func ParseIntermediateCondition(ctx *broker.Context, dependencies specs.Dependencies, condition Condition) (specs.NodeList, error) {
	expr, err := conditions.NewEvaluableExpression(ctx, condition.Expression)
	if err != nil {
		return nil, err
	}

	expression := &specs.Node{
		Type:      specs.NodeCondition,
		ID:        condition.Expression,
		Name:      "condition",
		DependsOn: specs.Dependencies{},
		Condition: expr,
	}

	result := specs.NodeList{expression}

	for _, references := range condition.References {
		nodes, err := ParseIntermediateResources(ctx, dependencies, references)
		if err != nil {
			return nil, err
		}

		DependsOn(expression, nodes...)
		result = append(result, result...)
	}

	for _, intermediate := range condition.Resources {
		node, err := ParseIntermediateNode(ctx, dependencies, intermediate)
		if err != nil {
			return nil, err
		}

		DependsOn(expression, node)
		result = append(result, node)
	}

	for _, condition := range condition.Conditions {
		nodes, err := ParseIntermediateCondition(ctx, dependencies, condition)
		if err != nil {
			return nil, err
		}

		DependsOn(expression, nodes[0])
		result = append(result, nodes...)
	}

	return result, nil
}

// ParseIntermediateOnError returns a specs on error
func ParseIntermediateOnError(ctx *broker.Context, onError *OnError) (*specs.OnError, error) {
	properties, err := ParseIntermediateParameters(ctx, onError.Body)
	if err != nil {
		return nil, err
	}

	result := &specs.OnError{
		Status:  properties["status"],
		Message: properties["message"],
	}

	if onError.Schema != "" {
		result.Response = &specs.ParameterMap{
			Schema: onError.Schema,
		}
	}

	if onError.Params != nil {
		params, err := ParseIntermediateParameters(ctx, onError.Params.Body)
		if err != nil {
			return nil, err
		}

		result.Params = params
	}

	return result, nil
}
