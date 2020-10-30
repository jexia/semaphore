package openapi3

// Options represent the available OpenAPI3 provider options
type Options int8

// Has checks whether the given key is available inside the options collection
func (options Options) Has(key Options) bool {
	return options&key != 0
}

const (
	// IncludeNotReferenced includes not referenced properties into the OpenAPI3 specification
	IncludeNotReferenced Options = 1 << iota
)

// DefaultOption represents the set of default options
const DefaultOption Options = 0
