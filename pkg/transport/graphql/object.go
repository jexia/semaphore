package graphql

import (
	"errors"

	"github.com/graphql-go/graphql"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/types"
	"github.com/jexia/semaphore/pkg/transport"
)

// ErrInvalidObject is thrown when the given property is not of type message
var ErrInvalidObject = errors.New("graphql only supports object types as root elements")

// NewObject constructs a new graphql object of the given specs
func NewObject(name string, description string, template specs.Template) (*graphql.Object, error) {
	if template.Type() != types.Message {
		return nil, ErrInvalidObject
	}

	typed, err := NewType(name, description, template)
	if err != nil {
		return nil, err
	}

	object, is := typed.(*graphql.Object)
	if !is {
		return nil, ErrInvalidObject
	}

	return object, nil
}

// NewType constructs a new output type from the given template.
func NewType(name string, description string, property specs.Template) (graphql.Output, error) {
	switch {
	case property.Message != nil:
		fields := graphql.Fields{}

		for _, nested := range property.Message {
			typed, err := NewType(name+"_"+nested.Name, nested.Description, nested.Template)
			if err != nil {
				return nil, err
			}

			field := &graphql.Field{
				Name:        nested.Name,
				Description: nested.Description,
				Type:        typed,
			}

			fields[nested.Name] = field
		}

		object := graphql.NewObject(graphql.ObjectConfig{
			Name:   name,
			Fields: fields,
		})

		return object, nil
	case property.Repeated != nil:
		template, err := property.Repeated.Template()
		if err != nil {
			return nil, err
		}

		typed, err := NewType(name, description, template)
		if err != nil {
			return nil, err
		}

		return graphql.NewList(typed), nil
	case property.Enum != nil:
		values := graphql.EnumValueConfigMap{}

		for key, field := range property.Enum.Keys {
			values[key] = &graphql.EnumValueConfig{
				Value:       key,
				Description: field.Description,
			}
		}

		config := graphql.EnumConfig{
			Name:        name + "_" + property.Enum.Name,
			Description: property.Enum.Description,
			Values:      values,
		}

		return graphql.NewEnum(config), nil
	case property.Scalar != nil:
		return gtypes[property.Scalar.Type], nil
	}

	return nil, nil
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
			return nil, ErrDuplicateObject{
				Name: name,
			}
		}

		return objects.collection[name], nil
	}

	obj, err := NewObject(name, property.Description, property.Template)
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
