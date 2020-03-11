package graphql

import (
	"github.com/graphql-go/graphql"
	"github.com/jexia/maestro/specs"
	"github.com/jexia/maestro/specs/types"
)

// NewArgs construct new field config arguments for the graphql schema
func NewArgs(prop *specs.Property) graphql.FieldConfigArgument {
	args := graphql.FieldConfigArgument{}
	if prop.Type == types.TypeMessage {
		if len(prop.Nested) == 0 {
			return nil
		}

		for _, nested := range prop.Nested {
			if nested.Type == types.TypeMessage {
				args[nested.Name] = &graphql.ArgumentConfig{
					Type: NewInputArgObject(nested),
				}
				continue
			}

			args[nested.Name] = &graphql.ArgumentConfig{
				Type: gtypes[nested.Type],
			}
		}

		return args
	}

	args[prop.Name] = &graphql.ArgumentConfig{
		Type: gtypes[prop.Type],
	}

	return args
}

// NewInputArgObject constructs a new input argument object
func NewInputArgObject(prop *specs.Property) *graphql.InputObject {
	if prop.Type != types.TypeMessage {
		return nil
	}

	fields := map[string]*graphql.InputObjectFieldConfig{}

	for _, nested := range prop.Nested {
		if nested.Type == types.TypeMessage {
			fields[nested.Name] = &graphql.InputObjectFieldConfig{
				Type: NewInputArgObject(nested),
			}

			continue
		}

		fields[nested.Name] = &graphql.InputObjectFieldConfig{
			Type: gtypes[prop.Type],
		}
	}

	return graphql.NewInputObject(graphql.InputObjectConfig{
		Fields: fields,
	})
}
