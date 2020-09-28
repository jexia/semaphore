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
func ResolveReferences(ctx *broker.Context, flows specs.FlowListInterface, mem functions.Collection) {
	logger.Info(ctx, "resolving flow references")

	for _, flow := range flows {
		for _, node := range flow.GetNodes() {
			ResolveNodeReferences(node, mem)
		}

		empty := specs.Dependencies{}

		if flow.GetOnError() != nil {
			ResolveOnErrorReferences(flow.GetOnError(), empty)
		}

		if flow.GetOutput() != nil {
			if flow.GetOutput().DependsOn == nil {
				flow.GetOutput().DependsOn = specs.Dependencies{}
			}

			ResolveParameterMapReferences(flow.GetOutput(), empty, mem)
		}
	}
}

// ResolveNodeReferences resolves the node references found inside the request and response property
func ResolveNodeReferences(node *specs.Node, mem functions.Collection) {
	if node.DependsOn == nil {
		node.DependsOn = specs.Dependencies{}
	}

	if node.OnError != nil {
		ResolveParameterMapReferences(node.OnError.Response, node.DependsOn, mem)
	}

	if node.Condition != nil {
		ResolveParamReferences(node.Condition.Params.Params, node.DependsOn)
	}

	if node.Intermediate != nil {
		ResolveParameterMapReferences(node.Intermediate, node.DependsOn, mem)
	}

	if node.Call != nil {
		ResolveParameterMapReferences(node.Call.Request, node.DependsOn, mem)
		ResolveParameterMapReferences(node.Call.Response, node.DependsOn, mem)
	}

	if node.Rollback != nil {
		ResolveParameterMapReferences(node.Rollback.Request, node.DependsOn, mem)
		ResolveParameterMapReferences(node.Rollback.Response, node.DependsOn, mem)
	}
}

// ResolveParameterMapReferences resolves the params inside the given parameter map
func ResolveParameterMapReferences(parameters *specs.ParameterMap, dependencies specs.Dependencies, mem functions.Collection) {
	if parameters == nil {
		return
	}

	if parameters.DependsOn == nil {
		parameters.DependsOn = specs.Dependencies{}
	}

	stack := mem.Load(parameters)
	ResolveFunctionsReferences(stack, parameters.DependsOn)

	ResolveParamReferences(parameters.Params, parameters.DependsOn)
	ResolveHeaderReferences(parameters.Header, parameters.DependsOn)

	if parameters.Property != nil {
		ResolvePropertyReferences(&parameters.Property.Template, parameters.DependsOn)
	}

	for key, val := range parameters.DependsOn {
		dependencies[key] = val
	}
}

// ResolveOnErrorReferences resolves the params inside the given parameter map
func ResolveOnErrorReferences(parameters *specs.OnError, dependencies specs.Dependencies) {
	if parameters == nil {
		return
	}

	if parameters.Response != nil {
		ResolveParamReferences(parameters.Response.Params, dependencies)
		ResolveHeaderReferences(parameters.Response.Header, dependencies)

		if parameters.Response.Property != nil {
			ResolvePropertyReferences(&parameters.Response.Property.Template, dependencies)
		}
	}
}

// ResolveParamReferences resolves all nested references made inside the given params
func ResolveParamReferences(params map[string]*specs.Property, dependencies specs.Dependencies) {
	if params == nil {
		return
	}

	for _, property := range params {
		ResolvePropertyReferences(&property.Template, dependencies)
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
				ResolvePropertyReferences(&arg.Template, dependencies)
			}
		}

		if function.Returns != nil {
			ResolvePropertyReferences(&function.Returns.Template, dependencies)
		}
	}
}

// ResolveHeaderReferences resolves all references made inside the header
func ResolveHeaderReferences(header specs.Header, dependencies specs.Dependencies) {
	for _, property := range header {
		ResolvePropertyReferences(&property.Template, dependencies)
	}
}

// ResolvePropertyReferences moves any property reference into the correct data structure
func ResolvePropertyReferences(tmpl *specs.Template, dependencies specs.Dependencies) {
	switch {
	case tmpl.Repeated != nil:
		for _, repeated := range tmpl.Repeated {
			ResolvePropertyReferences(&repeated, dependencies)
		}
	case tmpl.Message != nil:
		for _, nested := range tmpl.Message {
			ResolvePropertyReferences(&nested.Template, dependencies)
		}
	}

	if tmpl.Reference == nil {
		return
	}

	dependency := template.SplitPath(tmpl.Reference.Resource)[0]
	dependencies[dependency] = nil

	ScopePropertyReference(tmpl)
}

// ScopePropertyReference ensures that the root property is used inside the
// property reference.
func ScopePropertyReference(tmpl *specs.Template) {
	for {
		if tmpl.Reference == nil || tmpl.Reference.Property == nil {
			return
		}

		if tmpl.Reference == tmpl.Reference.Property.Reference {
			// NOTE: circular reference!
			return
		}

		target := tmpl.Reference.Property.Reference
		if target == nil {
			return
		}

		tmpl.Reference = target
	}
}
