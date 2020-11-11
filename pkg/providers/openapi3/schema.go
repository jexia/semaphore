package openapi3

import (
	"errors"
	"fmt"

	openapi "github.com/getkin/kin-openapi/openapi3"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/labels"
	"github.com/jexia/semaphore/pkg/specs/types"
)

func newSchemas(docs swaggers) (specs.Schemas, error) {
	var result = make(specs.Schemas)

	for file, doc := range docs {
		endpoints, err := scanPaths(doc.Paths)
		if err != nil {
			return nil, fmt.Errorf("failed to scan document paths '%s': %w", file, err)
		}

		for _, endpoints := range endpoints {
			for name, object := range endpoints.objects() {
				name := getCanonicalName(doc, name)
				// empty schemaRef, skipping
				if object == nil || object.ref == nil {
					continue
				}

				tpl, err := newTemplate(object.ref.Value)

				if err != nil {
					return nil, fmt.Errorf("failed to build property '%s' from file '%s': %w", name, file, err)
				}

				prop := &specs.Property{
					Template: tpl,
					Name:     name,
					Label:    labels.Optional,
				}

				result[prop.Name] = prop
			}
		}
	}

	return result, nil
}

// newTemplate builds a Property from the given swagger schemaRef component.
// the function returns a property without name, the name should be set in the top caller.
func newTemplate(schema *openapi.Schema) (tpl specs.Template, err error) {
	switch schema.Type {
	// handle as Message
	case "object":
		tpl, err = message(schema)
		if err != nil {
			return tpl, fmt.Errorf("failed to create object: %w", err)
		}

	// handle as Repeated
	case "array":
		tpl, err = repeated(schema)
		if err != nil {
			return tpl, fmt.Errorf("failed to create array: %w", err)
		}

	// check for OneOf, AnyOf, AllOf, Not
	// and handle according to its type
	case "":
		switch {
		case schema.OneOf != nil:
			tpl, err = oneOf(schema)
			if err != nil {
				return tpl, fmt.Errorf("failed to create oneOf type: %w", err)
			}
		default:
			return tpl, fmt.Errorf("unsupported mixed type")
		}

	// handle as scalar: string, integer, etc
	default:
		tpl, err = scalar(schema)
		if err != nil {
			return tpl, fmt.Errorf("failed to create scalar: %w", err)
		}
	}

	return tpl, nil
}

// builds a message template
func message(s *openapi.Schema) (specs.Template, error) {
	var (
		message = make(specs.Message, len(s.Properties))
	)

	for fieldName, ref := range s.Properties {
		if ref.Value == nil {
			return specs.Template{}, fmt.Errorf("field '%s' does not have schemaRef", fieldName)
		}

		tpl, err := newTemplate(ref.Value)
		if err != nil {
			return specs.Template{}, fmt.Errorf("failed to build field '%s': %w", fieldName, err)
		}

		message[fieldName] = &specs.Property{
			Name:     fieldName,
			Template: tpl,
		}
	}

	return specs.Template{
		Message: message,
	}, nil
}

// builds a scalar template
func scalar(s *openapi.Schema) (specs.Template, error) {
	var t types.Type

	switch s.Type {
	case "string":
		t = types.String
	case "boolean":
		t = types.Bool
	case "number":
		t = types.Float
	case "integer":
		t = types.Int32
	default:
		return specs.Template{}, fmt.Errorf("unknown type %s", s.Type)
	}

	return specs.Template{
		Scalar: &specs.Scalar{
			Type:    t,
			Default: s.Default,
		},
	}, nil
}

// builds a repeated template
func repeated(s *openapi.Schema) (specs.Template, error) {
	if s.Items == nil {
		return specs.Template{}, errors.New("empty item schemaRef")
	}

	item, err := newTemplate(s.Items.Value)
	if err != nil {
		return specs.Template{}, fmt.Errorf("failed to build item schemaRef: %w", err)
	}

	return specs.Template{
		Repeated: []specs.Template{item},
	}, nil
}

// builds an oneOf template
func oneOf(s *openapi.Schema) (specs.Template, error) {
	oneOf := make(specs.OneOf, len(s.OneOf))

	for id, ref := range s.OneOf {
		if ref.Value == nil {
			return specs.Template{}, fmt.Errorf("type at index %d does not have schemaRef", id)
		}

		tpl, err := newTemplate(ref.Value)
		if err != nil {
			return specs.Template{}, fmt.Errorf("failed to build type at index %d: %w", id, err)
		}

		oneOf[ref.Value.Type] = &specs.Property{
			Template: tpl,
		}
	}

	return specs.Template{
		OneOf: oneOf,
	}, nil
}
