package graphql

import (
	"strings"

	"github.com/graphql-go/graphql"
	"github.com/jexia/semaphore/pkg/broker/trace"
)

// SetField sets the given field inside the given fields on the given path
func SetField(path string, fields graphql.Fields, field *graphql.Field) error {
	if IsNestedPath(path) {
		parts := ParsePath(path)
		key := parts[0]

		target, has := fields[key]
		if !has {
			target = &graphql.Field{
				Name: key,
				Type: graphql.NewObject(graphql.ObjectConfig{
					Name:   key,
					Fields: make(graphql.Fields),
				}),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return p.Source, nil
				},
			}

			fields[key] = target
		}

		nested, is := target.Type.(*graphql.Object)
		if !is {
			return trace.New(trace.WithMessage("unable set field '%s' in '%s'", key, path))
		}

		err := SetFieldPath(nested, NewPath(parts[1:]), field)
		if err != nil {
			return err
		}

		return nil
	}

	fields[path] = field
	return nil
}

// ParsePath returns the path as steps
func ParsePath(path string) []string {
	return strings.Split(path, ".")
}

// NewPath constructs a new path of the given parts
func NewPath(parts []string) string {
	return strings.Join(parts, ".")
}

// IsNestedPath checks whether the given value is a path
func IsNestedPath(value string) bool {
	parts := strings.Split(value, ".")
	if len(parts) == 1 {
		return false
	}

	return true
}

// SetFieldPath sets the given field on the given path
func SetFieldPath(object *graphql.Object, path string, field *graphql.Field) error {
	parts := ParsePath(path)
	key := parts[0]

	fields := object.Fields()
	target, has := fields[key]

	if len(parts) > 1 {
		config := graphql.ObjectConfig{
			Name:   key,
			Fields: make(graphql.Fields),
		}

		nested := graphql.NewObject(config)

		if has {
			result, isObject := target.Type.(*graphql.Object)
			if !isObject {
				return trace.New(trace.WithMessage("unable set field '%s' in '%s'", key, path))
			}

			nested = result
		}

		part := &graphql.Field{
			Name: key,
			Type: nested,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return p.Source, nil
			},
		}

		object.AddFieldConfig(key, part)
		object.Fields()

		return SetFieldPath(nested, NewPath(parts[1:]), field)
	}

	// Check if field is set and path has parts left
	if has {
		return trace.New(trace.WithMessage("field already set '%s' in '%s'", key, path))
	}

	object.AddFieldConfig(key, field)
	object.Fields()

	return nil
}
