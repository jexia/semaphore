package graphql

import (
	"github.com/graphql-go/graphql"
	"github.com/jexia/maestro/pkg/specs"
	"github.com/jexia/maestro/pkg/specs/types"
)

// NewArgs construct new field config arguments for the graphql schema
func NewArgs(props *specs.ParameterMap) graphql.FieldConfigArgument {
	if props == nil {
		return graphql.FieldConfigArgument{}
	}

	prop := props.Property
	args := graphql.FieldConfigArgument{}
	// TODO: support repeated types (#71)

	if prop.Type == types.Message {
		if len(prop.Nested) == 0 {
			return nil
		}

		for _, nested := range prop.Nested {
			if nested.Type == types.Message {
				args[nested.Name] = &graphql.ArgumentConfig{
					Type:        NewInputArgObject(nested),
					Description: nested.Comment,
				}
				continue
			}

			args[nested.Name] = &graphql.ArgumentConfig{
				Type:        gtypes[nested.Type],
				Description: nested.Comment,
			}
		}

		return args
	}

	args[prop.Name] = &graphql.ArgumentConfig{
		Type:        gtypes[prop.Type],
		Description: prop.Comment,
	}

	return args
}

// NewInputArgObject constructs a new input argument object
func NewInputArgObject(prop *specs.Property) *graphql.InputObject {
	if prop.Type != types.Message {
		return nil
	}

	fields := map[string]*graphql.InputObjectFieldConfig{}

	for _, nested := range prop.Nested {
		if nested.Type == types.Message {
			fields[nested.Name] = &graphql.InputObjectFieldConfig{
				Type:        NewInputArgObject(nested),
				Description: nested.Comment,
			}

			continue
		}

		fields[nested.Name] = &graphql.InputObjectFieldConfig{
			Type:        gtypes[prop.Type],
			Description: nested.Comment,
		}
	}

	return graphql.NewInputObject(graphql.InputObjectConfig{
		Fields:      fields,
		Description: prop.Comment,
	})
}
