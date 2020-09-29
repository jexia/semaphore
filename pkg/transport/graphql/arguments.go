package graphql

import (
	"github.com/graphql-go/graphql"
	"github.com/jexia/semaphore/pkg/broker/trace"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/types"
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

	if prop.Type() != types.Message {
		return nil, trace.New(trace.WithMessage("arguments must be a object, received '%s'", prop.Type()))
	}

	if len(prop.Message) == 0 {
		return args, nil
	}

	for _, nested := range prop.Message {
		typ := gtypes[nested.Type()]

		switch {
		case nested.Message != nil:
			result, err := NewInputArgObject(nested)
			if err != nil {
				return nil, err
			}

			typ = result
			break
		case nested.Repeated != nil:
			typ = graphql.NewList(typ)
		case nested.Enum != nil:
			values := graphql.EnumValueConfigMap{}

			for key, field := range nested.Enum.Keys {
				values[key] = &graphql.EnumValueConfig{
					Value:       key,
					Description: field.Description,
				}
			}

			config := graphql.EnumConfig{
				Name:        nested.Name + "_" + nested.Enum.Name,
				Description: nested.Enum.Description,
				Values:      values,
			}

			typ = graphql.NewEnum(config)
			break
		}

		args[nested.Name] = &graphql.ArgumentConfig{
			Type:        typ,
			Description: nested.Description,
		}
	}

	return args, nil
}

// NewInputArgObject constructs a new input argument object
func NewInputArgObject(prop *specs.Property) (*graphql.InputObject, error) {
	if prop.Type() != types.Message {
		return nil, trace.New(trace.WithMessage("expected a message type received '%s'", prop.Type()))
	}

	fields := graphql.InputObjectConfigFieldMap{}

	for _, nested := range prop.Message {
		typ := gtypes[nested.Type()]

		switch {
		case nested.Message != nil:
			result, err := NewInputArgObject(nested)
			if err != nil {
				return nil, err
			}

			typ = result
		case nested.Repeated != nil:
			typ = graphql.NewList(typ)
		case nested.Enum != nil:
			values := graphql.EnumValueConfigMap{}

			for key, field := range nested.Enum.Keys {
				values[key] = &graphql.EnumValueConfig{
					Value:       key,
					Description: field.Description,
				}
			}

			config := graphql.EnumConfig{
				Name:        nested.Name + "_" + nested.Enum.Name,
				Description: nested.Enum.Description,
				Values:      values,
			}

			typ = graphql.NewEnum(config)
		}

		fields[nested.Name] = &graphql.InputObjectFieldConfig{
			Type:        typ,
			Description: nested.Description,
		}
	}

	result := graphql.NewInputObject(graphql.InputObjectConfig{
		Name:        prop.Name,
		Fields:      fields,
		Description: prop.Description,
	})

	return result, nil
}
