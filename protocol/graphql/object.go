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
			object, err := NewObject(name+"_"+nested.Name, nested)
			if err != nil {
				return nil, err
			}

			fields[nested.Name] = &graphql.Field{
				Type:        object,
				Description: nested.Desciptor.GetComment(),
			}
			continue
		}

		fields[nested.Name] = &graphql.Field{
			Type:        gtypes[nested.Type],
			Description: nested.Desciptor.GetComment(),
		}
	}

	config := graphql.ObjectConfig{
		Name:        name,
		Fields:      fields,
		Description: prop.Desciptor.GetComment(),
	}

	return graphql.NewObject(config), nil
}
