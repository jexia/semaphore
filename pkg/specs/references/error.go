package references

import "github.com/jexia/maestro/pkg/specs"

// ToError constructs a new error specification from the given parameters
func ToError(prop *specs.Property, object *specs.Error) *specs.Error {
	return &specs.Error{
		Schema:   object.Schema,
		Header:   object.Header,
		Property: prop,
	}
}
