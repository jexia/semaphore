package specs

// ResolvedProperty tracks resoved/seen/visited properties in order to escape
// from the infinite loop.
type ResolvedProperty struct {
	ResolvedTemplate
	byPath map[string]struct{}
}

// NewResolvedProperty creates a storage for visited properties.
func NewResolvedProperty() *ResolvedProperty {
	return &ResolvedProperty{
		ResolvedTemplate: make(ResolvedTemplate),
		byPath:           make(map[string]struct{}),
	}
}

// Resolved returns true if property was already visited by encoder/resolver etc.
func (resolved ResolvedProperty) Resolved(property *Property) bool {
	if property.Template != nil && property.Identifier != "" {
		return resolved.ResolvedTemplate.Resolved(property.Template)
	}

	if property.Path != "" {
		_, ok := resolved.byPath[property.Path]

		return ok
	}

	return false
}

// Resolve adds the property to the visited list.
func (resolved ResolvedProperty) Resolve(property *Property) {
	if property.Template != nil && property.Identifier != "" {
		resolved.ResolvedTemplate.Resolve(property.Template)
	}

	if property.Path != "" {
		resolved.byPath[property.Path] = struct{}{}
	}
}

type ResolvedTemplate map[string]struct{}

// Resolved returns true if the template was already resolved (by Identifier).
func (resolved ResolvedTemplate) Resolved(template *Template) bool {
	if template.Identifier == "" {
		return false
	}

	_, ok := resolved[template.Identifier]

	return ok
}

// Resolve provided template by Identifier.
func (resolved ResolvedTemplate) Resolve(template *Template) {
	if template.Identifier != "" {
		resolved[template.Identifier] = struct{}{}
	}
}
