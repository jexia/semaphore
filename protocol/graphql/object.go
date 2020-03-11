package graphql

import (
	"errors"

	"github.com/graphql-go/graphql"
	"github.com/jexia/maestro/specs"
	"github.com/jexia/maestro/specs/types"
)

// ErrInvalidObject is thrown when the given property is not of type message
var ErrInvalidObject = errors.New("graphql only supports object types as root elements")

// NewObject constructs a new graphql object of the given specs
func NewObject(name string, prop *specs.Property) (*graphql.Object, error) {
	if prop.Type != types.TypeMessage {
		return nil, ErrInvalidObject
	}

	fields := graphql.Fields{}
	for _, nested := range prop.Nested {
		if nested.Type == types.TypeMessage {
			field := &graphql.Field{
				Description: nested.Desciptor.GetComment(),
			}

			object, err := NewObject(name+"_"+nested.Name, nested)
			if err != nil {
				return nil, err
			}

			if nested.Label == types.LabelRepeated {
				field.Type = graphql.NewList(object)
			} else {
				field.Type = object
			}

			fields[nested.Name] = field
			continue
		}

		field := &graphql.Field{
			Description: nested.Desciptor.GetComment(),
		}

		typ := gtypes[nested.Type]
		if nested.Label == types.LabelRepeated {
			field.Type = graphql.NewList(typ)
		} else {
			field.Type = typ
		}

		fields[nested.Name] = field
	}

	config := graphql.ObjectConfig{
		Name:        name,
		Fields:      fields,
		Description: prop.Desciptor.GetComment(),
	}

	return graphql.NewObject(config), nil
}
