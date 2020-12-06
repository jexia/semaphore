package references

import (
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/template"
)

// ReferencedCollection contains a collection of paths
type ReferencedCollection map[string]struct{}

// Set deconstructs the given paths to ensure that all parts to the
// given paths are available. This is used to allow for easy comparisons against
// property paths.
//
// ex: "meta.info.name" results in: "meta", "meta.info", "meta.info.name"
func (collection ReferencedCollection) Set(path string) {
	absolute := ""
	parts := template.SplitPath(path)

	for _, part := range parts {
		current := template.JoinPath(absolute, part)
		collection[current] = struct{}{}
		absolute = current
	}

	collection[path] = struct{}{}
}

// Has verifies whether the given path is available inside the given collection.
func (collection ReferencedCollection) Has(path string) bool {
	_, has := collection[path]
	return has
}

// ReferencedParameterMapPaths constructs a new parameter map containing
// only the paths provided inisde the references collection. All properties
// not found iniside the collection will be ignored.
func ReferencedParameterMapPaths(referenced ReferencedCollection, property *specs.ParameterMap) *specs.ParameterMap {
	result := property.Clone()
	result.Property = ReferencedPathsProperty(referenced, result.Property)

	return result
}

// ReferencedPathsProperty constructs a new property containing
// only the paths provided inisde the references collection. All properties
// not found iniside the collection will be ignored.
func ReferencedPathsProperty(referenced ReferencedCollection, property *specs.Property) *specs.Property {
	if property == nil {
		return nil
	}

	result := property.Clone()
	removeNoneReferencedPathsTemplate(referenced, "", result.Template)

	return result
}

// removeNoneReferencedPathsTemplate removes all properties not referenced inside
// the given paths collection. This ensures that only referenced properties are kept
// inside the given template.
func removeNoneReferencedPathsTemplate(referenced ReferencedCollection, path string, template *specs.Template) {
	if template == nil {
		return
	}

	switch {
	case template.Message != nil:
		for key, nested := range template.Message {
			if nested == nil {
				delete(template.Message, key)
				continue
			}

			if _, has := referenced[nested.Path]; !has {
				delete(template.Message, key)
				continue
			}

			removeNoneReferencedPathsTemplate(referenced, nested.Path, nested.Template)
		}
	case template.Repeated != nil:
		if _, has := referenced[path]; !has {
			template.Repeated = nil
			break
		}

		for _, item := range template.Repeated {
			removeNoneReferencedPathsTemplate(referenced, path, item)
		}
	}
}

// ReferencedResourcePaths defines all paths references to the given resource
// inside the given flow. These references could be used to only include the properties
// used and referenced inside a given flow.
func ReferencedResourcePaths(flow specs.FlowInterface, resource string) ReferencedCollection {
	result := ReferencedCollection{}

	for _, node := range flow.GetNodes() {
		if node.Call != nil {
			parameterMapReferencedResourcePaths(result, node.Call.Request, resource)
			parameterMapReferencedResourcePaths(result, node.Call.Response, resource)
		}

		parameterMapReferencedResourcePaths(result, node.Intermediate, resource)

		if node.Condition != nil {
			parameterMapReferencedResourcePaths(result, node.Condition.Params, resource)
		}
	}

	parameterMapReferencedResourcePaths(result, flow.GetOutput(), resource)

	return result
}

func parameterMapReferencedResourcePaths(target ReferencedCollection, parameters *specs.ParameterMap, resource string) {
	if parameters == nil {
		return
	}

	for _, header := range parameters.Header {
		if header == nil {
			continue
		}

		propertyReferencedResourcePaths(target, header.Template, resource)
	}

	for _, params := range parameters.Params {
		if params == nil {
			continue
		}

		propertyReferencedResourcePaths(target, params.Template, resource)
	}

	if parameters.Property != nil {
		propertyReferencedResourcePaths(target, parameters.Property.Template, resource)
	}

	for _, stack := range parameters.Stack {
		if stack == nil {
			continue
		}

		propertyReferencedResourcePaths(target, stack.Template, resource)
	}
}

func propertyReferencedResourcePaths(collection ReferencedCollection, template *specs.Template, resource string) {
	if template.Reference != nil && template.Reference.Resource == resource {
		collection.Set(template.Reference.Path)
	}

	switch {
	case template.Message != nil:
		for _, nested := range template.Message {
			propertyReferencedResourcePaths(collection, nested.Template, resource)
		}
	case template.Repeated != nil:
		for _, item := range template.Repeated {
			propertyReferencedResourcePaths(collection, item, resource)
		}
	}
}
