package forwarding

import (
	"github.com/jexia/semaphore/pkg/broker"
	"github.com/jexia/semaphore/pkg/broker/logger"
	"github.com/jexia/semaphore/pkg/functions"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/template"
)

// ResolveReferences resolves all references inside the given manifest by forwarding references.
// If a reference is referencing another reference the node is marked as a dependency of the
// references resource and the referenced reference is copied over the the current resource.
func ResolveReferences(ctx *broker.Context, flows specs.FlowListInterface) {
	logger.Info(ctx, "resolving flow references")

	for _, flow := range flows {
		for _, node := range flow.GetNodes() {
			ResolveNodeReferences(node)
		}

		// Error and output dependencies could safely be ignored
		empty := specs.Dependencies{}

		if flow.GetOnError() != nil {
			ResolveOnError(flow.GetOnError(), empty)
		}

		if flow.GetOutput() != nil {
			if flow.GetOutput().DependsOn == nil {
				flow.GetOutput().DependsOn = specs.Dependencies{}
			}

			ResolveHeaderReferences(flow.GetOutput().Header, flow.GetOutput().DependsOn)
			ResolvePropertyReferences(flow.GetOutput().Property, flow.GetOutput().DependsOn)
		}
	}
}

// ResolveNodeReferences resolves the node references found inside the request and response property
func ResolveNodeReferences(node *specs.Node) {
	if node.DependsOn == nil {
		node.DependsOn = specs.Dependencies{}
	}

	if node.OnError != nil {
		ResolveParameterMap(node.OnError.Response, node.DependsOn)
	}

	if node.Condition != nil {
		ResolveParamReferences(node.Condition.Params.Params, node.DependsOn)
	}

	if node.Call != nil {
		ResolveParameterMap(node.Call.Request, node.DependsOn)
		ResolveParameterMap(node.Call.Response, node.DependsOn)
	}

	if node.Rollback != nil {
		ResolveParameterMap(node.Rollback.Request, node.DependsOn)
		ResolveParameterMap(node.Rollback.Response, node.DependsOn)
	}
}

// ResolveParameterMap resolves the params inside the given parameter map
func ResolveParameterMap(parameters *specs.ParameterMap, dependencies specs.Dependencies) {
	if parameters == nil {
		return
	}

	if parameters.DependsOn == nil {
		parameters.DependsOn = specs.Dependencies{}
	}

	ResolveParamReferences(parameters.Params, parameters.DependsOn)
	ResolveHeaderReferences(parameters.Header, parameters.DependsOn)
	ResolvePropertyReferences(parameters.Property, parameters.DependsOn)

	for key, val := range parameters.DependsOn {
		dependencies[key] = val
	}
}

// ResolveOnError resolves the params inside the given parameter map
func ResolveOnError(parameters *specs.OnError, dependencies specs.Dependencies) {
	if parameters == nil {
		return
	}

	if parameters.Response != nil {
		ResolveParamReferences(parameters.Response.Params, dependencies)
		ResolveHeaderReferences(parameters.Response.Header, dependencies)
		ResolvePropertyReferences(parameters.Response.Property, dependencies)
	}
}

// ResolveParamReferences resolves all nested references made inside the given params
func ResolveParamReferences(params map[string]*specs.Property, dependencies specs.Dependencies) {
	if params == nil {
		return
	}

	for _, property := range params {
		ResolvePropertyReferences(property, dependencies)
	}
}

// ResolveFunctionsReferences resolves all references made inside the given function arguments and return value
func ResolveFunctionsReferences(functions functions.Stack, dependencies specs.Dependencies) {
	if functions == nil {
		return
	}

	for _, function := range functions {
		if function.Arguments != nil {
			for _, arg := range function.Arguments {
				ResolvePropertyReferences(arg, dependencies)
			}
		}

		ResolvePropertyReferences(function.Returns, dependencies)
	}
}

// ResolveHeaderReferences resolves all references made inside the header
func ResolveHeaderReferences(header specs.Header, dependencies specs.Dependencies) {
	for _, prop := range header {
		ResolvePropertyReferences(prop, dependencies)
	}
}

// ResolvePropertyReferences moves any property reference into the correct data structure
func ResolvePropertyReferences(property *specs.Property, dependencies specs.Dependencies) {
	if property == nil {
		return
	}

	for _, nested := range property.Nested {
		ResolvePropertyReferences(nested, dependencies)
	}

	if property.Reference == nil || property.Reference.Property == nil {
		return
	}

	dependency := template.SplitPath(property.Reference.Resource)[0]
	dependencies[dependency] = nil

	ScopePropertyReference(property)
}

// ScopePropertyReference ensures that the root property is used inside the
// property reference.
func ScopePropertyReference(property *specs.Property) {
	if property.Reference == nil || property.Reference.Property == nil {
		return
	}

	if property.Reference.Property.Reference == nil {
		return
	}

	property.Reference = property.Reference.Property.Reference
	ScopePropertyReference(property)
}
