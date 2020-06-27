package dependencies

import (
	"github.com/jexia/maestro/pkg/functions"
	"github.com/jexia/maestro/pkg/instance"
	"github.com/jexia/maestro/pkg/logger"
	"github.com/jexia/maestro/pkg/specs"
	"github.com/jexia/maestro/pkg/specs/lookup"
	"github.com/jexia/maestro/pkg/specs/template"
)

// ResolveReferences resolves all references inside the given manifest by forwarding references.
// If a reference is referencing another reference the node is marked as a dependency of the
// references resource and the referenced reference is copied over the the current resource.
func ResolveReferences(ctx instance.Context, manifest *specs.FlowsManifest) {
	ctx.Logger(logger.Core).Info("Resolving manifest references")

	for _, flow := range manifest.Flows {
		for _, node := range flow.Nodes {
			ResolveNodeReferences(node)
		}

		// Error and output dependencies could safely be ignored
		empty := map[string]*specs.Node{}

		if flow.OnError != nil {
			ResolveOnError(flow.OnError, empty)
		}

		if flow.Output != nil {
			ResolveHeaderReferences(flow.Output.Header, empty)
			ResolvePropertyReferences(flow.Output.Property, empty)
		}
	}

	for _, proxy := range manifest.Proxy {
		for _, node := range proxy.Nodes {
			ResolveNodeReferences(node)
		}
	}
}

// ResolveNodeReferences resolves the node references found inside the request and response property
func ResolveNodeReferences(node *specs.Node) {
	if node.DependsOn == nil {
		node.DependsOn = map[string]*specs.Node{}
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
func ResolveParameterMap(parameters *specs.ParameterMap, dependencies map[string]*specs.Node) {
	if parameters == nil {
		return
	}

	ResolveParamReferences(parameters.Params, dependencies)
	ResolveHeaderReferences(parameters.Header, dependencies)
	ResolvePropertyReferences(parameters.Property, dependencies)
}

// ResolveOnError resolves the params inside the given parameter map
func ResolveOnError(parameters *specs.OnError, dependencies map[string]*specs.Node) {
	if parameters == nil {
		return
	}

	if parameters.Response != nil {
		ResolveParamReferences(parameters.Response.Params, dependencies)
		ResolveHeaderReferences(parameters.Response.Header, dependencies)
		ResolvePropertyReferences(parameters.Response.Property, dependencies)
	}
}

// ResolveFunctionsReferences resolves all references made inside the given function arguments and return value
func ResolveFunctionsReferences(functions functions.Stack, dependencies map[string]*specs.Node) {
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
func ResolveHeaderReferences(header specs.Header, dependencies map[string]*specs.Node) {
	for _, prop := range header {
		ResolvePropertyReferences(prop, dependencies)
	}
}

// ResolvePropertyReferences moves any property reference into the correct data structure
func ResolvePropertyReferences(property *specs.Property, dependencies map[string]*specs.Node) {
	if property == nil {
		return
	}

	if len(property.Nested) > 0 {
		for _, nested := range property.Nested {
			ResolvePropertyReferences(nested, dependencies)
		}

		return
	}

	if property.Reference == nil {
		return
	}

	if property.Reference.Property == nil {
		return
	}

	resource, _ := lookup.ParseResource(property.Reference.Resource)
	if resource != template.StackResource && (resource != template.InputResource && resource != template.ErrorResource) {
		dependencies[template.SplitPath(property.Reference.Resource)[0]] = nil
	}

	clone := CloneProperty(property.Reference.Property, property.Reference, property.Name, property.Path)
	property.Reference = clone.Reference
	property.Nested = clone.Nested
}

// ResolveParamReferences resolves all nested references made inside the given params
func ResolveParamReferences(params map[string]*specs.Property, dependencies map[string]*specs.Node) {
	if params == nil {
		return
	}

	for key, property := range params {
		if property.Reference == nil {
			continue
		}

		resource, _ := lookup.ParseResource(property.Reference.Resource)
		if resource != template.StackResource && resource != template.InputResource {
			dependency := template.SplitPath(property.Reference.Resource)[0]
			dependencies[dependency] = nil
		}

		clone := CloneProperty(property.Reference.Property, property.Reference, property.Name, property.Path)
		params[key].Reference = clone.Reference
		params[key].Nested = clone.Nested
	}
}

// CloneProperty clones the given property with the given reference, name and path
func CloneProperty(source *specs.Property, reference *specs.PropertyReference, name string, path string) *specs.Property {
	result := source.Clone()
	result.Name = name
	result.Path = path
	result.Reference = reference

	if source != nil {
		if source.Reference != nil {
			result.Reference = &specs.PropertyReference{
				Resource: source.Reference.Resource,
				Path:     source.Reference.Path,
			}
		}

		if len(source.Nested) != 0 {
			result.Nested = make(map[string]*specs.Property, len(source.Nested))

			for key, nested := range source.Nested {
				ref := &specs.PropertyReference{
					Resource: result.Reference.Resource,
					Path:     template.JoinPath(result.Path, key),
				}

				result.Nested[key] = CloneProperty(nested, ref, key, template.JoinPath(path, key))
			}
		}
	}

	return result
}
