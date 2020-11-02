package specs

// ResolvedProperty tracks resoved/seen/visited properties in order to escape
// from the infinite loop.
type ResolvedProperty struct {
	ResolvedTemplate
	byIdentifier map[string]struct{}
}

// NewResolvedProperty creates a storage for visited properties.
func NewResolvedProperty() *ResolvedProperty {
	return &ResolvedProperty{
		ResolvedTemplate: make(ResolvedTemplate),
		byIdentifier:     make(map[string]struct{}),
	}
}

// Resolved returns true if property was already visited by encoder/resolver etc.
func (resolved ResolvedProperty) Resolved(property *Property) bool {
	if property.Template != nil && property.Template.Identifier != "" {
		return resolved.ResolvedTemplate.Resolved(property.Template)
	}

	if property.Identifier != "" {
		_, ok := resolved.byIdentifier[property.Identifier]

		return ok
	}

	return false
}

// Resolve adds the property to the visited list.
func (resolved ResolvedProperty) Resolve(property *Property) {
	if property.Template != nil && property.Template.Identifier != "" {
		resolved.ResolvedTemplate.Resolve(property.Template)
	}

	if property.Identifier != "" {
		resolved.byIdentifier[property.Identifier] = struct{}{}
	}
}

// ResolvedTemplate contains the list of templates that were already processed.
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
