package specs

type ResolvedProperty struct {
	ResolvedTemplate
	byPath map[string]struct{}
}

func NewResolvedProperty() *ResolvedProperty {
	return &ResolvedProperty{
		ResolvedTemplate: make(ResolvedTemplate),
		byPath:           make(map[string]struct{}),
	}
}

func (resolved ResolvedProperty) Resolved(property *Property) bool {
	if property.Identifier != "" {
		return resolved.ResolvedTemplate.Resolved(property.Template)
	}

	if property.Path != "" {
		_, ok := resolved.byPath[property.Path]

		return ok
	}

	return false
}

func (resolved ResolvedProperty) Resolve(property *Property) {
	if property.Identifier != "" {
		resolved.ResolvedTemplate.Resolve(property.Template)
	}

	if property.Path != "" {
		resolved.byPath[property.Path] = struct{}{}
	}
}

type ResolvedTemplate map[string]struct{}

func (resolved ResolvedTemplate) Resolved(template Template) bool {
	_, ok := resolved[template.Identifier]

	return ok
}

func (resolved ResolvedTemplate) Resolve(template Template) {
	resolved[template.Identifier] = struct{}{}
}
