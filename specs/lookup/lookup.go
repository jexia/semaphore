package lookup

import (
	"github.com/jexia/maestro/specs"
)

// SelfRef represents the syntax used to reference the entire object
const SelfRef = "."

// ReferenceMap holds the resource references and their representing parameter map
type ReferenceMap map[string]PathLookup

// PathLookup represents a lookup method that returns the property available on the given path
type PathLookup func(path string) *specs.Property

// GetFlow attempts to find the given flow inside the given manifest
func GetFlow(manifest specs.Manifest, name string) *specs.Flow {
	for _, flow := range manifest.Flows {
		if flow.Name == name {
			return flow
		}
	}

	return nil
}

// ParseResource parses the given resource into the resource and props
func ParseResource(resource string) (string, string) {
	resources := specs.SplitPath(resource)

	target := resources[0]
	prop := GetDefaultProp(target)

	if len(resources) > 1 {
		prop = specs.JoinPath(resources[1:]...)
	}

	return target, prop
}

// GetDefaultProp returns the default resource for the given resource
func GetDefaultProp(resource string) string {
	if resource == specs.InputResource {
		return specs.ResourceRequest
	}

	return specs.ResourceResponse
}

// GetAvailableResources fetches the available resources able to be referenced
// until the given breakpoint (call.Name) has been reached.
func GetAvailableResources(flow specs.FlowManager, breakpoint string) map[string]ReferenceMap {
	references := make(map[string]ReferenceMap, len(flow.GetNodes())+1)

	if flow.GetInput() != nil {
		references[specs.InputResource] = ReferenceMap{
			specs.ResourceRequest: ParameterMapLookup(flow.GetInput().Property),
			specs.ResourceHeader:  HeaderLookup(flow.GetInput().Header),
		}
	}

	for _, node := range flow.GetNodes() {
		if node.Name == breakpoint {
			break
		}

		resources := ReferenceMap{}

		if node.Call != nil {
			if node.Call.Request != nil {
				resources[specs.ResourceRequest] = ParameterMapLookup(node.Call.Request.Property)
			}

			if node.Call.Response != nil {
				resources[specs.ResourceResponse] = ParameterMapLookup(node.Call.Response.Property)
				resources[specs.ResourceHeader] = HeaderLookup(node.Call.Response.Header)
			}
		}

		references[node.Name] = resources
	}

	return references
}

// GetResourceReference attempts to return the resource reference property
func GetResourceReference(reference *specs.PropertyReference, references map[string]ReferenceMap) *specs.Property {
	target, prop := ParseResource(reference.Resource)

	for resource, refs := range references {
		if resource != target {
			continue
		}

		return GetReference(reference.Path, prop, refs)
	}

	return nil
}

// GetReference attempts to lookup and return the available property on the given path
func GetReference(path string, prop string, references ReferenceMap) *specs.Property {
	lookup, has := references[prop]
	if !has {
		return nil
	}

	return lookup(path)
}

// HeaderLookup attempts to lookup the given path inside the header
func HeaderLookup(header specs.Header) PathLookup {
	return func(path string) *specs.Property {
		for key, header := range header {
			if key == path {
				return header
			}
		}

		return nil
	}
}

// ParameterMapLookup attempts to lookup the given path inside the params collection
func ParameterMapLookup(param *specs.Property) PathLookup {
	return func(path string) *specs.Property {
		if path == SelfRef {
			return param
		}

		if param.Path == path {
			return param
		}

		if param.Nested == nil {
			return nil
		}

		for _, param := range param.Nested {
			lookup := ParameterMapLookup(param)(path)
			if lookup != nil {
				return lookup
			}
		}

		return nil
	}
}
