package specs

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/jexia/semaphore/pkg/specs/labels"
	"github.com/jexia/semaphore/pkg/specs/metadata"
	"github.com/jexia/semaphore/pkg/specs/types"
)

func TestPropertyUnmarshalDefault(t *testing.T) {
	t.Parallel()

	type test struct {
		input    *Property
		expected reflect.Kind
	}

	tests := map[string]test{
		"int64": {
			input: &Property{
				Label: labels.Optional,
				Template: Template{
					Scalar: &Scalar{
						Type:    types.Int64,
						Default: 100,
					},
				},
			},
			expected: reflect.Int64,
		},
		"sint64": {
			input: &Property{
				Label: labels.Optional,
				Template: Template{
					Scalar: &Scalar{
						Type:    types.Sint64,
						Default: 100,
					},
				},
			},
			expected: reflect.Int64,
		},
		"sfixed64": {
			input: &Property{
				Label: labels.Optional,
				Template: Template{
					Scalar: &Scalar{
						Type:    types.Sfixed64,
						Default: 100,
					},
				},
			},
			expected: reflect.Int64,
		},
		"uint64": {
			input: &Property{
				Label: labels.Optional,
				Template: Template{
					Scalar: &Scalar{
						Type:    types.Uint64,
						Default: 100,
					},
				},
			},
			expected: reflect.Uint64,
		},
		"fixed64": {
			input: &Property{
				Label: labels.Optional,
				Template: Template{
					Scalar: &Scalar{
						Type:    types.Fixed64,
						Default: 100,
					},
				},
			},
			expected: reflect.Uint64,
		},
		"int32": {
			input: &Property{
				Label: labels.Optional,
				Template: Template{
					Scalar: &Scalar{
						Type:    types.Int32,
						Default: 100,
					},
				},
			},
			expected: reflect.Int32,
		},
		"sint32": {
			input: &Property{
				Label: labels.Optional,
				Template: Template{
					Scalar: &Scalar{
						Type:    types.Sint32,
						Default: 100,
					},
				},
			},
			expected: reflect.Int32,
		},
		"sfixed32": {
			input: &Property{
				Label: labels.Optional,
				Template: Template{
					Scalar: &Scalar{
						Type:    types.Sfixed32,
						Default: 100,
					},
				},
			},
			expected: reflect.Int32,
		},
		"uint32": {
			input: &Property{
				Label: labels.Optional,
				Template: Template{
					Scalar: &Scalar{
						Type:    types.Uint32,
						Default: 100,
					},
				},
			},
			expected: reflect.Uint32,
		},
		"fixed32": {
			input: &Property{
				Label: labels.Optional,
				Template: Template{
					Scalar: &Scalar{
						Type:    types.Fixed32,
						Default: 100,
					},
				},
			},
			expected: reflect.Uint32,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			t.Run("scalar", func(t *testing.T) {
				input, err := json.Marshal(test.input)
				if err != nil {
					t.Fatal(err)
				}

				result := Property{}
				err = json.Unmarshal(input, &result)
				if err != nil {
					t.Fatal(err)
				}

				if result.Scalar == nil {
					t.Fatalf("result scalar has not been set")
				}

				kind := reflect.TypeOf(result.Scalar.Default).Kind()
				if kind != test.expected {
					t.Errorf("unexpected type %+v, expected %+v", kind, test.expected)
				}
			})

			t.Run("message", func(t *testing.T) {
				message := &Property{
					Template: Template{
						Message: Message{
							name: test.input,
						},
					},
				}

				input, err := json.Marshal(message)
				if err != nil {
					t.Fatal(err)
				}

				result := Property{}
				err = json.Unmarshal(input, &result)
				if err != nil {
					t.Fatal(err)
				}

				if result.Message == nil {
					t.Fatalf("result message is not set")
				}

				property := result.Message[name]
				if property == nil {
					t.Fatalf("result message property %s is not set", name)
				}

				kind := reflect.TypeOf(property.Scalar.Default).Kind()
				if kind != test.expected {
					t.Errorf("unexpected type %+v, expected %+v", kind, test.expected)
				}
			})

			t.Run("repeated", func(t *testing.T) {
				message := &Property{
					Template: Template{
						Repeated: Repeated{
							test.input.Template,
						},
					},
				}

				input, err := json.Marshal(message)
				if err != nil {
					t.Fatal(err)
				}

				result := Property{}
				err = json.Unmarshal(input, &result)
				if err != nil {
					t.Fatal(err)
				}

				if result.Repeated == nil {
					t.Fatalf("result repeated is not set")
				}

				template := result.Repeated[0]
				if template.Scalar == nil {
					t.Fatalf("result repeated default property %s is not set", name)
				}

				kind := reflect.TypeOf(template.Scalar.Default).Kind()
				if kind != test.expected {
					t.Errorf("unexpected type %+v, expected %+v", kind, test.expected)
				}
			})
		})
	}
}

func TestObjectsAppend(t *testing.T) {
	objects := Schemas{}

	objects.Append(Schemas{
		"first":  &Property{},
		"second": &Property{},
	})

	if len(objects) != 2 {
		t.Fatalf("unexpected length %+v, expected 2", len(objects))
	}
}

func TestObjectsAppendNilValue(t *testing.T) {
	var objects Schemas
	objects.Append(nil)
}

func TestObjectsGet(t *testing.T) {
	objects := Schemas{
		"first":  &Property{},
		"second": &Property{},
	}

	result := objects.Get("second")
	if result == nil {
		t.Fatal("unexpected empty result")
	}
}

func TestObjectsGetNilValue(t *testing.T) {
	var objects Schemas
	result := objects.Get("unknown")
	if result != nil {
		t.Fatalf("unexpected result %+v", result)
	}
}

func TestPropertyClone(t *testing.T) {
	property := &Property{
		Meta:        metadata.WithValue(nil, nil, nil),
		Position:    1,
		Description: "sample",
		Name:        "first",
		Path:        "path",
		Label:       labels.Optional,
		Template: Template{
			Reference: &PropertyReference{},
			Scalar: &Scalar{
				Default: false,
				Type:    types.String,
			},
			Repeated: Repeated{},
			Message: Message{
				"first": {Name: "first", Path: "first"},
			},
			Enum: &Enum{
				Name: "unknown",
			},
		},
		Raw: "first",
		Options: Options{
			"sample": "option",
		},
	}

	result := property.Clone()
	if result == nil {
		t.Error("unexpected result, expected property reference clone to be returned")
	}

	if result.Meta != property.Meta {
		t.Errorf("unexpected meta %+v, expected %+v", result.Meta, property.Meta)
	}

	if result.Position != property.Position {
		t.Errorf("unexpected position %+v", result.Position)
	}

	if result.Description != property.Description {
		t.Errorf("unexpected description %+v", result.Description)
	}

	if result.Name != property.Name {
		t.Errorf("unexpected name %+v", result.Name)
	}

	if result.Path != property.Path {
		t.Errorf("unexpected path %+v", result.Path)
	}

	if result.Scalar == nil {
		t.Fatalf("scalar not set")
	}

	if result.Scalar.Default != property.Scalar.Default {
		t.Errorf("unexpected default %+v", result.Scalar.Default)
	}

	if result.Scalar.Type != property.Scalar.Type {
		t.Errorf("unexpected type %+v", result.Scalar.Type)
	}

	if result.Label != property.Label {
		t.Errorf("unexpected label %+v", result.Label)
	}

	if result.Reference == nil || result.Reference == property.Reference {
		t.Errorf("unexpected reference %+v", result.Reference)
	}

	if result.Message == nil {
		t.Fatalf("message not set")
	}

	if len(result.Message) != len(property.Message) {
		t.Errorf("unexpected message properties %+v", result.Message)
	}

	if result.Repeated == nil {
		t.Fatalf("repeated not set")
	}

	if result.Repeated == nil {
		t.Fatalf("repeated is not set")
	}

	if result.Raw != property.Raw {
		t.Errorf("unexpected raw %+v", result.Raw)
	}

	if len(result.Options) != len(property.Options) {
		t.Errorf("unexpected options %+v", result.Options)
	}

	if result.Enum == nil {
		t.Fatal("enum not set")
	}

	if result.Enum.Name != property.Enum.Name {
		t.Errorf("unexpected enum %+v", result.Enum)
	}
}

func TestPropertyCompare(t *testing.T) {
	var (
		createScalar = func() *Property {
			return &Property{
				Name:     "age",
				Path:     "dog",
				Position: 0,
				Label:    labels.Required,
				Template: Template{
					Scalar: &Scalar{
						Type: types.Int32,
					},
				},
			}
		}

		createEnum = func() *Property {
			retriever := &EnumValue{
				Key:      "retriever",
				Position: 0,
			}
			shepherd := &EnumValue{
				Key:      "shepherd",
				Position: 1,
			}

			return &Property{
				Name:     "breed",
				Path:     "dog",
				Position: 0,
				Template: Template{
					Enum: &Enum{
						Name: "breed",
						Keys: map[string]*EnumValue{
							retriever.Key: retriever, shepherd.Key: shepherd,
						},
						Positions: map[int32]*EnumValue{
							retriever.Position: retriever, shepherd.Position: shepherd,
						},
					},
				},
			}
		}

		createRepeated = func() *Property {
			return &Property{
				Name: "hunters",
				Path: "dogs",
				Template: Template{
					Repeated: Repeated{
						createEnum().Template,
					},
				},
			}
		}

		createMessage = func() *Property {
			return &Property{
				Name: "dog",
				Path: "request",
				Template: Template{
					Message: Message{
						"age":   createScalar(),
						"breed": createEnum(),
					},
				},
			}
		}

		// createAnother<T> behaves like create<T> with a tiny difference.
		// For example, a scalar type or a name might be different.
		// We use it to test comparison <T> against a bit different <T>

		createAnotherScalar = func() *Property {
			prop := createScalar()
			prop.Scalar.Type = types.String
			return prop
		}

		createAnotherEnum = func() *Property {
			prop := createEnum()
			prop.Enum.Keys["foobar"] = &EnumValue{Key: "foobar", Position: 100}
			prop.Enum.Positions[100] = &EnumValue{Key: "foobar", Position: 100}
			return prop
		}

		createAnotherRepeated = func() *Property {
			prop := createRepeated()
			prop.Repeated = Repeated{
				{
					Enum: &Enum{},
				},
			}
			return prop
		}

		createAnotherMessage = func() *Property {
			prop := createMessage()
			prop.Message["age"] = createAnotherScalar()
			return prop
		}

		shouldFail = func(t *testing.T, property, schema *Property) {
			if property.Compare(schema) == nil {
				t.Fatalf("nil, an error expected")
			}
		}

		shouldMatch = func(t *testing.T, property, schema *Property) {
			if err := property.Compare(schema); err != nil {
				t.Fatalf("returns unexpected error: %s", err)
			}
		}
	)

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

		shouldFail(t, prop, &Property{})
	})

	t.Run("should fail due to empty property, but filled schema", func(t *testing.T) {
		schema := createScalar()

		shouldFail(t, &Property{}, schema)
	})

	t.Run("", func(t *testing.T) {
		props := map[string]*Property{
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

func TestParameterMapClone(t *testing.T) {
	property := &ParameterMap{
		Meta:   metadata.WithValue(nil, nil, nil),
		Schema: "com.schema",
		Params: map[string]*Property{
			"sample": {},
		},
		Options: Options{
			"sample": "option",
		},
		Header: Header{
			"sample": {},
		},
		Property: &Property{},
		Stack: map[string]*Property{
			"hash": {},
		},
	}

	result := property.Clone()
	if result == nil {
		t.Error("unexpected result, expected property reference clone to be returned")
	}

	if result.Meta != property.Meta {
		t.Errorf("unexpected meta %+v, expected %+v", result.Meta, property.Meta)
	}

	if result.Schema != property.Schema {
		t.Errorf("unexpected schema %+v", result.Schema)
	}

	if result.Property == nil || result.Property == property.Property {
		t.Errorf("unexpected property %+v", result.Property)
	}

	if len(result.Options) != len(property.Options) {
		t.Errorf("unexpected options %+v", result.Options)
	}

	if len(result.Header) != len(property.Header) {
		t.Errorf("unexpected header %+v", result.Header)
	}

	if len(result.Stack) != len(property.Stack) {
		t.Errorf("unexpected stack %+v", result.Stack)
	}
}

func TestParameterMapCloneNilValue(t *testing.T) {
	var params *ParameterMap

	result := params.Clone()
	if result != nil {
		t.Errorf("unexpected result %+v", result)
	}
}

func TestPropertyDefaultValue(t *testing.T) {
	type test struct {
		property *Property
		expected interface{}
	}

	var tests = []test{
		{
			property: &Property{
				Template: Template{
					Scalar: &Scalar{
						Type:    types.String,
						Default: "develop smart not hard",
					},
				},
			},
			expected: "develop smart not hard",
		},
		{
			property: &Property{
				Template: Template{
					Enum: &Enum{},
				},
			},
			expected: nil,
		},
		{
			property: &Property{
				Template: Template{
					Message: Message{},
				},
			},
			expected: nil,
		},
		{
			property: &Property{
				Template: Template{
					Repeated: Repeated{},
				},
			},
			expected: nil,
		},
		{
			property: &Property{},
			expected: nil,
		},
	}

	for _, test := range tests {
		t.Run(string(test.property.Type()), func(t *testing.T) {
			if actual := test.property.DefaultValue(); actual != test.expected {
				t.Errorf("default value '%+v' was expected to be '%+v'", actual, test.expected)
			}
		})
	}
}
