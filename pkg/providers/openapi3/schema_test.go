package openapi3

import (
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/types"
	"github.com/stretchr/testify/assert"
)

type testFn func(*testing.T, *specs.Template) bool

func shouldBeScalar(t *testing.T, tpl *specs.Template, expected specs.Scalar, opts ...interface{}) bool {
	prefix := ""
	if len(opts) > 0 && opts[0] != nil {
		prefix = opts[0].(string)
	}

	if tpl.Type() != expected.Type {
		t.Errorf("%s should be Scalar(%s), but it is %s", prefix, expected.Type, tpl.Type())
		return false
	}

	if tpl.Scalar.Type != expected.Type || tpl.Scalar.Default != expected.Default {
		t.Errorf(
			"%s scalar should be of type %s with default value %v, but it is %s (%v)", prefix,
			expected.Type, expected.Default,
			tpl.Scalar.Type, tpl.Scalar.Default)
		return false
	}

	return true
}

func shouldBeMessage(t *testing.T, tpl *specs.Template, fields map[string]testFn, opts ...interface{}) bool {
	prefix := ""
	if len(opts) > 0 && opts[0] != nil {
		prefix = opts[0].(string)
	}

	if tpl.Type() != types.Message {
		t.Errorf("%s should be Message, but it is %s", prefix, tpl.Type())
		return false
	}

	if fields == nil {
		return true
	}

	if len(tpl.Message) != len(fields) {
		t.Errorf("%s should have length %d, but it is %d", prefix, len(fields), len(tpl.Message))
		return false
	}

	valid := true
	for field, check := range fields {
		prop, ok := tpl.Message[field]
		if !ok {
			t.Errorf("%s should include field '%s'", prefix, field)
			return false
		}

		if check == nil {
			continue
		}

		if !check(t, prop.Template) {
			valid = false
		}
	}

	return valid
}

func shouldBeRepeated(t *testing.T, tpl *specs.Template, check testFn, opts ...interface{}) bool {
	prefix := ""
	if len(opts) > 0 && opts[0] != nil {
		prefix = opts[0].(string)
	}

	if tpl.Type() != types.Array {
		t.Errorf("%s should be Repeated, but it is %s", prefix, tpl.Type())
		return false
	}

	item, err := tpl.Repeated.Template()
	if err != nil {
		t.Errorf("%s should have valid template, but: %s", prefix, err)
		return false
	}

	if check == nil {
		return true
	}

	return check(t, item)
}

func shouldBeOneOf(t *testing.T, tpl *specs.Template, checks []testFn, opts ...interface{}) bool {
	prefix := ""
	if len(opts) > 0 && opts[0] != nil {
		prefix = opts[0].(string)
	}

	if tpl.Type() != types.OneOf {
		t.Errorf("%s should be oneOf, but it is %s", prefix, tpl.Type())
		return false
	}

	if checks == nil {
		return true
	}

	if len(tpl.OneOf) != len(checks) {
		t.Errorf("%s should have length %d, but it is %d", prefix, len(checks), len(tpl.OneOf))
	}

	valid := true
	for i, item := range tpl.OneOf {
		if !checks[i](t, item) {
			valid = false
		}
	}

	return valid
}

func Test_newTemplate(t *testing.T) {
	loader := openapi3.NewSwaggerLoader()
	doc, err := loader.LoadSwaggerFromFile("./fixtures/petstore.yml")
	if err != nil {
		t.Fatalf("failed to load openapi file: %s", err)
	}

	// a helper function to build a property for a particular component
	build := func(t *testing.T, name string) (*specs.Template, error) {
		schemaRef, ok := doc.Components.Schemas[name]

		if !ok {
			t.Fatalf("failed to get schemaRef '%s'", name)
		}

		if schemaRef.Value == nil {
			t.Fatalf("unexpected nil schemaRef value")
		}

		return newTemplate(schemaRef.Value)
	}

	t.Run("test scalar", func(t *testing.T) {
		got, err := build(t, "RandomName")

		if err != nil {
			t.Errorf("newTemplate() unexpected error: %s", err)
			return
		}

		shouldBeScalar(t, got, specs.Scalar{Type: types.String, Default: "Fido"})
	})

	t.Run("test repeated of scalars", func(t *testing.T) {
		got, err := build(t, "DogNames")

		if err != nil {
			t.Errorf("newTemplate() unexpected error: %s", err)
			return
		}

		shouldBeRepeated(t, got, func(t *testing.T, template *specs.Template) bool {
			return shouldBeScalar(t, template, specs.Scalar{Type: types.String})
		})
	})

	t.Run("test message", func(t *testing.T) {
		got, err := build(t, "Dog")

		if err != nil {
			t.Errorf("newTemplate() unexpected error: %s", err)
			return
		}

		shouldBeMessage(t, got, map[string]testFn{
			"name": func(t *testing.T, template *specs.Template) bool {
				return shouldBeScalar(t, template, specs.Scalar{Type: types.String}, "field 'name'")
			},
			"is_good_boy": func(t *testing.T, template *specs.Template) bool {
				return shouldBeScalar(t, template, specs.Scalar{Type: types.Bool}, "field 'is_good_boy'")
			},
		})
	})

	t.Run("test oneOf", func(t *testing.T) {
		got, err := build(t, "Pet")

		if err != nil {
			t.Errorf("newTemplate() unexpected error: %s", err)
			return
		}

		shouldBeOneOf(t, got, []testFn{
			func(t *testing.T, template *specs.Template) bool {
				return shouldBeMessage(t, template, map[string]testFn{
					"name": func(t *testing.T, template *specs.Template) bool {
						return shouldBeScalar(t, template, specs.Scalar{Type: types.String}, "name")
					},
					"meow": func(t *testing.T, template *specs.Template) bool {
						return shouldBeScalar(t, template, specs.Scalar{Type: types.Bool}, "meow")
					},
				}, "Cat")
			},

			func(t *testing.T, template *specs.Template) bool {
				return shouldBeMessage(t, template, map[string]testFn{
					"name": func(t *testing.T, template *specs.Template) bool {
						return shouldBeScalar(t, template, specs.Scalar{Type: types.String}, "name")
					},
					"is_good_boy": func(t *testing.T, template *specs.Template) bool {
						return shouldBeScalar(t, template, specs.Scalar{Type: types.Bool}, "is_good_boy")
					},
				}, "Dog")
			},
		}, "favorite_pet")
	})
}

func Test_scalar(t *testing.T) {
	t.Run("should support all the openapi types", func(t *testing.T) {
		// map of openapi scalar types to testEndpoint types
		types := map[string]types.Type{
			"string":  types.String,
			"number":  types.Float,
			"integer": types.Int32,
			"boolean": types.Bool,
		}

		for oapiType, want := range types {
			t.Run(oapiType, func(t *testing.T) {
				schema := &openapi3.Schema{Type: oapiType}
				got, err := scalar(schema)

				if err != nil {
					t.Errorf("scalar() error = %v", err)
					return
				}

				if got.Scalar.Type != want {
					t.Errorf("scalar() = %v, want %v", got.Scalar.Type, want)
				}
			})
		}
	})

	t.Run("should support default value", func(t *testing.T) {
		schema := &openapi3.Schema{Type: "string", Default: "foo"}
		got, err := scalar(schema)

		if err != nil {
			t.Errorf("scalar() error = %v", err)
			return
		}

		if got.Scalar.Default != "foo" {
			t.Errorf("scalar().Default = %v, want %v", got.Scalar.Default, schema.Default)
		}
	})
}

func Test_newSchemas(t *testing.T) {
	loader := openapi3.NewSwaggerLoader()
	doc, err := loader.LoadSwaggerFromFile("./fixtures/petstore.yml")
	if err != nil {
		t.Fatalf("failed to load openapi file: %s", err)
	}

	docs := swaggers{
		"fixtures/petstore.yml": doc,
	}

	schemas, err := newSchemas(docs)
	if err != nil {
		t.Fatalf("failed to build schemas: %s", err)
	}

	assert.Len(t, schemas, 6)

	expected := []string{
		"com.semaphore.GET:/pets:Response[application/json][default]",
		"com.semaphore.GET:/pets:Response[application/json][200]",
		"com.semaphore.POST:/pets:Request[application/json]",
		"com.semaphore.POST:/pets:Response[application/json][default]",
		"com.semaphore.GET:/pets/{petId}:Response[application/json][default]",
		"com.semaphore.GET:/pets/{petId}:Response[application/json][200]",
	}

	for _, want := range expected {
		_, ok := schemas[want]
		if !ok {
			assert.Truef(t, ok, "should include property %s", want)
		}
	}
}
