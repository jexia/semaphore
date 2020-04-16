package lookup

import (
	"github.com/jexia/maestro/pkg/specs"
	"github.com/jexia/maestro/pkg/specs/template"
	"github.com/jexia/maestro/pkg/specs/types"
)

// SelfRef represents the syntax used to reference the entire object
const SelfRef = "."

// ReferenceMap holds the resource references and their representing parameter map
type ReferenceMap map[string]PathLookup

// PathLookup represents a lookup method that returns the property available on the given path
type PathLookup func(path string) *specs.Property

// GetFlow attempts to find the given flow inside the given manifest
func GetFlow(manifest specs.FlowsManifest, name string) *specs.Flow {
	for _, flow := range manifest.Flows {
		if flow.Name == name {
			return flow
		}
	}

	return nil
}

// ParseResource parses the given resource into the resource and props
func ParseResource(resource string) (string, string) {
	resources := template.SplitPath(resource)

	target := resources[0]
	prop := GetDefaultProp(target)

	if len(resources) > 1 {
		prop = template.JoinPath(resources[1:]...)
	}

	return target, prop
}

// GetDefaultProp returns the default resource for the given resource
func GetDefaultProp(resource string) string {
	if resource == template.InputResource {
		return template.ResourceRequest
	}

	return template.ResourceResponse
}

// GetNextResource returns the resource after the given breakpoint
func GetNextResource(flow specs.FlowResourceManager, breakpoint string) string {
	for index, node := range flow.GetNodes() {
		if node.Name == breakpoint {
			nodes := flow.GetNodes()
			next := index + 1
			if next >= len(nodes) {
				return template.OutputResource
			}

			return nodes[next].Name
		}
	}

	return breakpoint
}

// GetAvailableResources fetches the available resources able to be referenced
// until the given breakpoint (call.Name) has been reached.
func GetAvailableResources(flow specs.FlowResourceManager, breakpoint string) map[string]ReferenceMap {
	length := len(flow.GetNodes()) + 2
	references := make(map[string]ReferenceMap, length)
	references[template.StackResource] = ReferenceMap{}

	if flow.GetInput() != nil {
		references[template.InputResource] = ReferenceMap{
			template.ResourceRequest: PropertyLookup(flow.GetInput().Property),
			template.ResourceHeader:  HeaderLookup(flow.GetInput().Header),
		}
	}

	for _, node := range flow.GetNodes() {
		references[node.Name] = ReferenceMap{}
		if node.Call != nil {
			if node.Call.Request != nil {
				if node.Call.Request.Stack != nil {
					for key, returns := range node.Call.Request.Stack {
						references[template.StackResource][key] = PropertyLookup(returns)
					}
				}

				references[node.Name][template.ResourceRequest] = PropertyLookup(node.Call.Request.Property)
			}
		}

		if node.Name == breakpoint {
			break
		}

		if node.Call != nil {
			if node.Call.Response != nil {
				if node.Call.Response.Stack != nil {
					for key, returns := range node.Call.Response.Stack {
						references[template.StackResource][key] = PropertyLookup(returns)
					}
				}

				references[node.Name][template.ResourceResponse] = PropertyLookup(node.Call.Response.Property)
				references[node.Name][template.ResourceHeader] = VariableHeaderLookup(node.Call.Response.Header)
			}
		}
	}

	return references
}

// GetResourceReference attempts to return the resource reference property
func GetResourceReference(reference *specs.PropertyReference, references map[string]ReferenceMap, breakpoint string) *specs.Property {
	target, prop := ParseResource(reference.Resource)
	if target == "" {
		target = breakpoint
	}

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

// VariableHeaderLookup returns a string property for the given path and sets the header property
func VariableHeaderLookup(header specs.Header) PathLookup {
	return func(path string) *specs.Property {
		header[path] = &specs.Property{
			Path:    path,
			Type:    types.String,
			Default: "",
		}

		return header[path]
	}
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

// PropertyLookup attempts to lookup the given path inside the params collection
func PropertyLookup(param *specs.Property) PathLookup {
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
			lookup := PropertyLookup(param)(path)
			if lookup != nil {
				return lookup
			}
		}

		return nil
	}
}
