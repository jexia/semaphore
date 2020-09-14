package graphql

import (
	"errors"

	"github.com/graphql-go/graphql"
	"github.com/jexia/semaphore/pkg/broker/trace"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/labels"
	"github.com/jexia/semaphore/pkg/specs/types"
	"github.com/jexia/semaphore/pkg/transport"
)

// ErrInvalidObject is thrown when the given property is not of type message
var ErrInvalidObject = errors.New("graphql only supports object types as root elements")

// NewObject constructs a new graphql object of the given specs
func NewObject(name string, prop *specs.Property) (*graphql.Object, error) {
	if prop.Type != types.Message {
		return nil, ErrInvalidObject
	}

	fields := graphql.Fields{}
	for _, nested := range prop.Repeated {
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

			for key, field := range nested.Enum.Keys {
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
func NewSchemaObject(objects *Objects, name string, object *transport.Object) (*graphql.Object, error) {
	if object == nil {
		return NewEmptyObject(objects, name), nil
	}

	params := object.Definition
	if params == nil {
		return NewEmptyObject(objects, name), nil
	}

	property := params.Property
	_, has := objects.collection[name]
	if has {
		if objects.properties[name] != property {
			return nil, trace.New(trace.WithMessage("duplicate object '%s'", name))
		}

		return objects.collection[name], nil
	}

	obj, err := NewObject(name, property)
	if err != nil {
		return nil, err
	}

	if len(obj.Fields()) == 0 {
		return NewEmptyObject(objects, name), nil
	}

	objects.collection[name] = obj
	objects.properties[name] = property

	return obj, nil
}

// NewEmptyObject constructs a new empty object with the given name and stores it inside the objects collection.
// The object is claimed inside the object properties as nil. If the name is already claimed is the configured graphql object returned.
func NewEmptyObject(objects *Objects, name string) *graphql.Object {
	_, has := objects.properties[name]
	if has {
		return objects.collection[name]
	}

	result := graphql.NewObject(graphql.ObjectConfig{
		Name:        name,
		Description: "empty object",
		Fields: graphql.Fields{
			"_empty": &graphql.Field{
				Type: graphql.Boolean,
			},
		},
	})

	objects.collection[name] = result
	objects.properties[name] = nil

	return result
}
