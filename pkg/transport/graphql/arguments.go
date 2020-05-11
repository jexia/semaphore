package graphql

import (
	"github.com/graphql-go/graphql"
	"github.com/jexia/maestro/pkg/specs"
	"github.com/jexia/maestro/pkg/specs/labels"
	"github.com/jexia/maestro/pkg/specs/trace"
	"github.com/jexia/maestro/pkg/specs/types"
)

// NewArgs construct new field config arguments for the graphql schema
func NewArgs(props *specs.ParameterMap) (graphql.FieldConfigArgument, error) {
	if props == nil {
		return graphql.FieldConfigArgument{}, nil
	}

	if props.Property == nil {
		return graphql.FieldConfigArgument{}, nil
	}

	prop := props.Property
	args := graphql.FieldConfigArgument{}

	if prop.Type != types.Message {
		return nil, trace.New(trace.WithMessage("arguments must be a object, received '%s'", prop.Type))
	}

	if len(prop.Nested) == 0 {
		return args, nil
	}

	for _, nested := range prop.Nested {
		typ := gtypes[nested.Type]
		if nested.Type == types.Message {
			result, err := NewInputArgObject(nested)
			if err != nil {
				return nil, err
			}

			typ = result
		}

		if prop.Label == labels.Repeated {
			typ = graphql.NewList(typ)
		}

		args[nested.Name] = &graphql.ArgumentConfig{
			Type:        typ,
			Description: nested.Comment,
		}
	}

	return args, nil
}

// NewInputArgObject constructs a new input argument object
func NewInputArgObject(prop *specs.Property) (*graphql.InputObject, error) {
	if prop.Type != types.Message {
		return nil, trace.New(trace.WithMessage("expected a message type received '%s'", prop.Type))
	}

	fields := graphql.InputObjectConfigFieldMap{}

	for _, nested := range prop.Nested {
		typ := gtypes[nested.Type]
		if nested.Type == types.Message {
			result, err := NewInputArgObject(nested)
			if err != nil {
				return nil, err
			}

			typ = result
		}

		if prop.Label == labels.Repeated {
			typ = graphql.NewList(typ)
		}

		fields[nested.Name] = &graphql.InputObjectFieldConfig{
			Type:        typ,
			Description: nested.Comment,
		}
	}

	result := graphql.NewInputObject(graphql.InputObjectConfig{
		Name:        prop.Name,
		Fields:      fields,
		Description: prop.Comment,
	})

	return result, nil
}
