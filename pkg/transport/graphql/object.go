package graphql

import (
	"errors"

	"github.com/graphql-go/graphql"
	"github.com/jexia/maestro/pkg/specs"
	"github.com/jexia/maestro/pkg/specs/labels"
	"github.com/jexia/maestro/pkg/specs/trace"
	"github.com/jexia/maestro/pkg/specs/types"
)

// ErrInvalidObject is thrown when the given property is not of type message
var ErrInvalidObject = errors.New("graphql only supports object types as root elements")

// NewObject constructs a new graphql object of the given specs
func NewObject(name string, prop *specs.Property) (*graphql.Object, error) {
	if prop.Type != types.Message {
		return nil, ErrInvalidObject
	}

	fields := graphql.Fields{}
	for _, nested := range prop.Nested {
		if nested.Type == types.Message {
			field := &graphql.Field{
				Description: nested.Comment,
			}

			object, err := NewObject(name+"_"+nested.Name, nested)
			if err != nil {
				return nil, err
			}

			if nested.Label == labels.Repeated {
				field.Type = graphql.NewList(object)
			} else {
				field.Type = object
			}

			fields[nested.Name] = field
			continue
		}

		field := &graphql.Field{
			Description: nested.Comment,
		}

		field.Type = gtypes[nested.Type]

		if nested.Label == labels.Repeated {
			field.Type = graphql.NewList(gtypes[nested.Type])
		}

		if nested.Type == types.Enum {
			values := graphql.EnumValueConfigMap{}

			for key, field := range nested.Enum.Values {
				values[key] = &graphql.EnumValueConfig{
					Value:       key,
					Description: field.Description,
				}
			}

			config := graphql.EnumConfig{
				Name:        name + "_" + nested.Enum.Name,
				Description: nested.Enum.Description,
				Values:      values,
			}

			field.Type = graphql.NewEnum(config)
		}

		fields[nested.Name] = field
	}

	config := graphql.ObjectConfig{
		Name:        name,
		Fields:      fields,
		Description: prop.Comment,
	}

	return graphql.NewObject(config), nil
}

// NewObjects constructs a new objects collection
func NewObjects() *Objects {
	return &Objects{
		properties: map[string]*specs.Property{},
		collection: map[string]*graphql.Object{},
	}
}

// Objects represents a schema objects collection
type Objects struct {
	properties map[string]*specs.Property
	collection map[string]*graphql.Object
}

// NewSchemaObject constructs a new object for the given property with the given name.
// If a object with the same name already exists is it used instead.
func NewSchemaObject(objects *Objects, name string, params *specs.ParameterMap) (*graphql.Object, error) {
	if params == nil {
		return &graphql.Object{}, nil
	}

	property := params.Property
	_, has := objects.collection[name]
	if has {
		if objects.properties[name] != property {
			return nil, trace.New(trace.WithMessage("duplicate object '%s'", name))
		}

		return objects.collection[name], nil
	}

	object, err := NewObject(name, property)
	if err != nil {
		return nil, err
	}

	objects.collection[name] = object
	objects.properties[name] = property

	return object, nil
}
