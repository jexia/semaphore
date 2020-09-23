package compare

import (
	"fmt"
	"testing"

	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/labels"
	"github.com/jexia/semaphore/pkg/specs/types"
)

func TestCheckPropertyTypes(t *testing.T) {
	createScalar := func() *specs.Property {
		return &specs.Property{
			Name:     "age",
			Path:     "dog",
			Position: 0,
			Label:    labels.Required,
			Template: specs.Template{
				Scalar: &specs.Scalar{
					Type: types.Int32,
				},
			},
		}
	}
	createEnum := func() *specs.Property {
		retriever := &specs.EnumValue{
			Key:      "retriever",
			Position: 0,
		}
		shepherd := &specs.EnumValue{
			Key:      "shepherd",
			Position: 1,
		}

		return &specs.Property{
			Name:     "breed",
			Path:     "dog",
			Position: 0,
			Template: specs.Template{
				Enum: &specs.Enum{
					Name: "breed",
					Keys: map[string]*specs.EnumValue{
						retriever.Key: retriever, shepherd.Key: shepherd,
					},
					Positions: map[int32]*specs.EnumValue{
						retriever.Position: retriever, shepherd.Position: shepherd,
					},
				},
			},
		}
	}
	createRepeated := func() *specs.Property {
		return &specs.Property{
			Name: "hunters",
			Path: "dogs",
			Template: specs.Template{
				Repeated: &specs.Repeated{
					Property: createEnum(),
				},
			},
		}
	}
	createMessage := func() *specs.Property {
		return &specs.Property{
			Name: "dog",
			Path: "request",
			Template: specs.Template{
				Message: specs.Message{
					"age":   createScalar(),
					"breed": createEnum(),
				},
			},
		}
	}

	// createAnother<T> behaves like create<T> with a tiny difference.
	// For example, a scalar type or a name might be different.
	// We use it to test comparison <T> against a bit different <T>

	createAnotherScalar := func() *specs.Property {
		prop := createScalar()
		prop.Scalar.Type = types.String
		return prop
	}
	createAnotherEnum := func() *specs.Property {
		prop := createEnum()
		prop.Enum.Keys["foobar"] = &specs.EnumValue{Key: "foobar", Position: 100}
		prop.Enum.Positions[100] = &specs.EnumValue{Key: "foobar", Position: 100}
		return prop
	}
	createAnotherRepeated := func() *specs.Property {
		prop := createRepeated()
		prop.Repeated.Property.Enum = &specs.Enum{}
		return prop
	}
	createAnotherMessage := func() *specs.Property {
		prop := createMessage()
		prop.Message["age"] = createAnotherScalar()
		return prop
	}

	shouldFail := func(t *testing.T, property, schema *specs.Property) {
		if CheckPropertyTypes(property, schema) == nil {
			t.Fatalf("CheckPropertyTypes() = nil, an error expected")
		}
	}

	shouldMatch := func(t *testing.T, property, schema *specs.Property) {
		got := CheckPropertyTypes(property, schema)
		if got != nil {
			t.Fatalf("CheckPropertyTypes() returns unexpected error: %s", got)
		}
	}

	t.Run("should fail as schema is nil", func(t *testing.T) {
		shouldFail(t, createScalar(), nil)
	})

	t.Run("should fail due to different type", func(t *testing.T) {
		schema := createScalar()
		prop := createScalar()

		prop.Scalar.Type = types.Float

		shouldFail(t, prop, schema)
	})

	t.Run("should fail due to different label", func(t *testing.T) {
		schema := createScalar()
		prop := createScalar()

		prop.Label = labels.Optional

		shouldFail(t, prop, schema)
	})

	t.Run("should fail due to empty schema, but filled property", func(t *testing.T) {
		prop := createScalar()

		prop.Label = labels.Optional

		shouldFail(t, prop, &specs.Property{})
	})

	t.Run("should fail due to empty property, but filled schema", func(t *testing.T) {
		schema := createScalar()

		shouldFail(t, &specs.Property{}, schema)
	})

	t.Run("", func(t *testing.T) {
		props := map[string]*specs.Property{
			"scalar":   createScalar(),
			"enum":     createEnum(),
			"message":  createMessage(),
			"repeated": createRepeated(),

			"another_scalar":   createAnotherScalar(),
			"another_enum":     createAnotherEnum(),
			"another_message":  createAnotherMessage(),
			"another_repeated": createAnotherRepeated(),
		}

		for propKind, prop := range props {
			for schemaKind, schema := range props {
				if propKind == schemaKind {
					t.Run(fmt.Sprintf("should match property '%s' against schema '%s'", propKind, schemaKind), func(t *testing.T) {
						shouldMatch(t, prop, schema)
					})
				} else {
					t.Run(fmt.Sprintf("should not match property '%s' against schema '%s'", propKind, schemaKind), func(t *testing.T) {
						shouldFail(t, prop, schema)
					})
				}
			}
		}
	})
}
