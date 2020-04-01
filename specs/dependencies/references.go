package dependencies

import (
	"github.com/jexia/maestro/instance"
	"github.com/jexia/maestro/logger"
	"github.com/jexia/maestro/specs"
	"github.com/jexia/maestro/specs/lookup"
)

// ResolveReferences resolves all references inside the given manifest by forwarding references.
// If a reference is referencing another reference the node is marked as a dependency of the
// references resource and the referenced reference is copied over the the current resource.
func ResolveReferences(ctx instance.Context, manifest *specs.Manifest) {
	ctx.Logger(logger.Core).Info("Resolving manifest references")

	for _, flow := range manifest.Flows {
		for _, node := range flow.Nodes {
			ResolveNodeReferences(node)
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

	ResolvePropertyReferences(node.Call.Request.Property, node.DependsOn)
	ResolvePropertyReferences(node.Call.Response.Property, node.DependsOn)
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
	if resource != specs.StackResource && resource != specs.InputResource {
		dependencies[property.Reference.Resource] = nil
	}

	clone := CloneProperty(property.Reference.Property, property.Reference, property.Name, property.Path)
	property.Reference = clone.Reference
	property.Nested = clone.Nested
}

// CloneProperty clones the given property with the given reference, name and path
func CloneProperty(source *specs.Property, reference *specs.PropertyReference, name string, path string) *specs.Property {
	result := &specs.Property{
		Name:      name,
		Path:      path,
		Reference: reference,
		Default:   source.Default,
		Type:      source.Type,
		Label:     source.Label,
		Expr:      source.Expr,
		Desciptor: source.Desciptor,
	}

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
				Path:     specs.JoinPath(result.Path, key),
			}

			result.Nested[key] = CloneProperty(nested, ref, key, specs.JoinPath(path, key))
		}
	}

	return result
}
