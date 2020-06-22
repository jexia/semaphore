package lookup

import (
	"strings"

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
		return template.RequestResource
	}

	return template.ResponseResource
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
			template.RequestResource: PropertyLookup(flow.GetInput().Property),
			template.HeaderResource:  HeaderLookup(flow.GetInput().Header),
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

				references[node.Name][template.ParamsResource] = ParamsLookup(node.Call.Request.Params, flow, breakpoint)
				references[node.Name][template.RequestResource] = PropertyLookup(node.Call.Request.Property)
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

				references[node.Name][template.ResponseResource] = PropertyLookup(node.Call.Response.Property)
				references[node.Name][template.HeaderResource] = VariableHeaderLookup(node.Call.Response.Header)
			}
		}
	}

	return references
}

// GetResourceReference attempts to return the resource reference property
func GetResourceReference(reference *specs.PropertyReference, references map[string]ReferenceMap, breakpoint string) *specs.Property {
	if reference == nil {
		return nil
	}

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
	for key := range header {
		header[key].Path = strings.ToLower(header[key].Path)
	}

	return func(path string) *specs.Property {
		for key, header := range header {
			if strings.ToLower(key) == strings.ToLower(path) {
				return header
			}
		}

		return nil
	}
}

// PropertyLookup attempts to lookup the given path inside the params collection
func PropertyLookup(param *specs.Property) PathLookup {
	return func(path string) *specs.Property {
		if param == nil {
			return nil
		}

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

// ParamsLookup constructs a lookup method able to lookup property references inside the given params map
func ParamsLookup(params map[string]*specs.Property, flow specs.FlowResourceManager, breakpoint string) PathLookup {
	return func(path string) *specs.Property {
		if params == nil {
			return nil
		}

		for key, param := range params {
			if key == path {
				if param.Reference == nil {
					return param
				}

				references := GetAvailableResources(flow, breakpoint)
				reference := GetResourceReference(param.Reference, references, breakpoint)
				if reference == nil {
					return nil
				}

				result := reference.Clone()
				result.Reference = param.Reference
				return result
			}
		}

		return nil
	}
}

// ResolveSelfReference appends the given resource if the path is a self reference
func ResolveSelfReference(path string, resource string) string {
	if string(path[0]) != SelfRef {
		return path
	}

	return template.JoinPath(resource, path[1:])
}
