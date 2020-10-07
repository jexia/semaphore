package lookup

import (
	"strings"

	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/labels"
	"github.com/jexia/semaphore/pkg/specs/template"
	"github.com/jexia/semaphore/pkg/specs/types"
)

// SelfRef represents the syntax used to reference the entire object
const SelfRef = "."

// ReferenceMap holds the resource references and their representing parameter map
type ReferenceMap map[string]PathLookup

// PathLookup represents a lookup method that returns the property available on the given path
type PathLookup func(path string) *specs.Property

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

	if resource == template.ErrorResource {
		return template.ResponseResource
	}

	return template.ResponseResource
}

// GetNextResource returns the resource after the given breakpoint
func GetNextResource(flow specs.FlowInterface, breakpoint string) string {
	for index, node := range flow.GetNodes() {
		if node.ID == breakpoint {
			nodes := flow.GetNodes()
			next := index + 1
			if next >= len(nodes) {
				return template.OutputResource
			}

			return nodes[next].ID
		}
	}

	return breakpoint
}

// GetAvailableResources fetches the available resources able to be referenced
// until the given breakpoint (call.Name) has been reached.
func GetAvailableResources(flow specs.FlowInterface, breakpoint string) map[string]ReferenceMap {
	length := len(flow.GetNodes()) + 2
	references := make(map[string]ReferenceMap, length)
	references[template.StackResource] = ReferenceMap{}

	if flow.GetInput() != nil {
		references[template.InputResource] = ReferenceMap{
			template.RequestResource: PropertyLookup(flow.GetInput().Property),
			template.HeaderResource:  HeaderLookup(flow.GetInput().Header),
		}
	}

	if breakpoint == template.OutputResource {
		if flow.GetOnError() != nil {
			references[template.ErrorResource] = ReferenceMap{
				template.ResponseResource: OnErrLookup(template.OutputResource, flow.GetOnError()),
				template.ParamsResource:   ParamsLookup(flow.GetOnError().Params, flow, ""),
			}
		}
	}

	for _, node := range flow.GetNodes() {
		references[node.ID] = ReferenceMap{}

		if node.Intermediate != nil {
			if node.Intermediate.Stack != nil {
				for key, returns := range node.Intermediate.Stack {
					references[template.StackResource][key] = PropertyLookup(returns)
				}
			}

			references[node.ID][template.ResponseResource] = PropertyLookup(node.Intermediate.Property)
			references[node.ID][template.HeaderResource] = VariableHeaderLookup(node.Intermediate.Header)
		}

		if node.Call != nil {
			if node.Call.Request != nil {
				if node.Call.Request.Stack != nil {
					for key, returns := range node.Call.Request.Stack {
						references[template.StackResource][key] = PropertyLookup(returns)
					}
				}

				references[node.ID][template.ParamsResource] = ParamsLookup(node.Call.Request.Params, flow, breakpoint)
				references[node.ID][template.RequestResource] = PropertyLookup(node.Call.Request.Property)
			}

			if node.Call.Response != nil {
				if node.Call.Response.Stack != nil {
					for key, returns := range node.Call.Response.Stack {
						references[template.StackResource][key] = PropertyLookup(returns)
					}
				}

				references[node.ID][template.ResponseResource] = PropertyLookup(node.Call.Response.Property)
				references[node.ID][template.HeaderResource] = VariableHeaderLookup(node.Call.Response.Header)
			}
		}

		if node.ID == breakpoint {
			if node.GetOnError() != nil {
				references[template.ErrorResource] = ReferenceMap{
					template.ResponseResource: OnErrLookup(node.ID, node.GetOnError()),
					template.ParamsResource:   ParamsLookup(node.GetOnError().Params, flow, breakpoint),
				}

				if node.GetOnError().Response != nil {
					references[node.ID][template.ErrorResource] = PropertyLookup(node.GetOnError().Response.Property)
				}
			}
		}
	}

	if flow.GetOutput() != nil {
		if flow.GetOutput().Stack != nil {
			for key, returns := range flow.GetOutput().Stack {
				references[template.StackResource][key] = PropertyLookup(returns)
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
			Path: path,
			Template: specs.Template{
				Scalar: &specs.Scalar{
					Type:    types.String,
					Default: "",
				},
			},
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
			if strings.EqualFold(key, path) {
				return header
			}
		}

		return nil
	}
}

// PropertyLookup attempts to lookup the given path inside the params collection
func PropertyLookup(param *specs.Property) PathLookup {
	return func(path string) *specs.Property {
		switch {
		case param == nil:
			return nil
		case path == SelfRef:
			return param
		case param.Path == path:
			return param
		case param.Repeated != nil:
			// TODO: allow to reference indexes
			var template, _ = param.Repeated.Template()
			return PropertyLookup(
				&specs.Property{
					Template: template,
				},
			)(path)
		case param.Message != nil:
			for _, nested := range param.Template.Message {
				lookup := PropertyLookup(nested)(path)
				if lookup != nil {
					return lookup
				}
			}

			return nil
		default:
			return nil
		}
	}
}

// ParamsLookup constructs a lookup method able to lookup property references inside the given params map
func ParamsLookup(params map[string]*specs.Property, flow specs.FlowInterface, breakpoint string) PathLookup {
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

				if param.Scalar != nil && result.Scalar != nil {
					param.Scalar.Type = result.Scalar.Type
				}

				param.Label = result.Label
				param.Reference.Property = result

				return param
			}
		}

		return nil
	}
}

// OnErrLookup constructs a lookup method able to lookup error references
func OnErrLookup(node string, spec *specs.OnError) PathLookup {
	if spec == nil {
		spec = &specs.OnError{}
	}

	if spec.Message == nil {
		spec.Message = &specs.Property{
			Label: labels.Optional,
			Template: specs.Template{
				Scalar: &specs.Scalar{
					Type: types.String,
				},
			},
		}
	}

	if spec.Status == nil {
		spec.Status = &specs.Property{
			Label: labels.Optional,
			Template: specs.Template{
				Scalar: &specs.Scalar{
					Type: types.Int64,
				},
			},
		}
	}

	if spec.Message.Reference != nil && spec.Message.Reference.Resource == template.ErrorResource {
		spec.Message.Reference.Property = spec.Message
	}

	if spec.Status.Reference != nil && spec.Status.Reference.Resource == template.ErrorResource {
		spec.Status.Reference.Property = spec.Status
	}

	return func(path string) *specs.Property {
		// The references {{ error:message }} and {{ error:status }} are always available and set by the flow manager on error
		if path == "message" {
			return spec.Message
		}

		if path == "status" {
			return spec.Status
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
