package formencoded

import (
	"errors"
	"io"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/labels"
	"github.com/jexia/semaphore/pkg/specs/template"
	"github.com/jexia/semaphore/pkg/specs/types"
)

var enum = &specs.Enum{
	Keys: map[string]*specs.EnumValue{
		"UNKOWN": {
			Key:      "UNKOWN",
			Position: 0,
		},
		"PENDING": {
			Key:      "PENDING",
			Position: 1,
		},
	},
	Positions: map[int32]*specs.EnumValue{
		0: {
			Key:      "UNKOWN",
			Position: 0,
		},
		1: {
			Key:      "PENDING",
			Position: 1,
		},
	},
}

var schema = &specs.ParameterMap{
	Property: &specs.Property{
		Label: labels.Optional,
		Template: specs.Template{
			Message: &specs.Message{
				Properties: map[string]*specs.Property{
					"bad_label": {
						Name:  "bad_label",
						Path:  "bad_label",
						Label: 0,
						Template: specs.Template{
							Scalar: &specs.Scalar{
								Type: types.String,
							},
						},
					},
					"no_nested_schema": {
						Name:  "no_nested_schema",
						Path:  "no_nested_schema",
						Label: labels.Optional,
						Template: specs.Template{
							Message: &specs.Message{},
						},
					},
					"numeric": {
						Name:  "numeric",
						Path:  "numeric",
						Label: labels.Optional,
						Template: specs.Template{
							Scalar: &specs.Scalar{
								Type: types.Int32,
							},
						},
					},
					"message": {
						Name:  "message",
						Path:  "message",
						Label: labels.Optional,
						Reference: &specs.PropertyReference{
							Resource: template.InputResource,
							Path:     "message",
						},
						Template: specs.Template{
							Scalar: &specs.Scalar{
								Type: types.String,
							},
						},
					},
					"status": {
						Name:  "status",
						Path:  "status",
						Label: labels.Optional,
						Reference: &specs.PropertyReference{
							Resource: template.InputResource,
							Path:     "status",
						},
						Template: specs.Template{
							Enum: enum,
						},
					},
					"nested": {
						Name:  "nested",
						Path:  "nested",
						Label: labels.Optional,
						Template: specs.Template{
							Message: &specs.Message{
								Properties: map[string]*specs.Property{
									"first": {
										Name:  "first",
										Path:  "nested.first",
										Label: labels.Optional,
										Reference: &specs.PropertyReference{
											Resource: template.InputResource,
											Path:     "nested.first",
										},
										Template: specs.Template{
											Scalar: &specs.Scalar{
												Type: types.String,
											},
										},
									},
									"second": {
										Name:  "second",
										Path:  "nested.second",
										Label: labels.Optional,
										Reference: &specs.PropertyReference{
											Resource: template.InputResource,
											Path:     "nested.second",
										},
										Template: specs.Template{
											Scalar: &specs.Scalar{
												Type: types.String,
											},
										},
									},
								},
							},
						},
					},
					"repeating_string": {
						Name: "repeating_string",
						Path: "repeating_string",
						Reference: &specs.PropertyReference{
							Resource: template.InputResource,
							Path:     "repeating_string",
						},
						Template: specs.Template{
							Repeated: &specs.Repeated{
								Property: &specs.Property{
									Template: specs.Template{
										Scalar: &specs.Scalar{
											Type: types.String,
										},
									},
								},
							},
						},
					},
					"repeating_enum": {
						Name:  "repeating_enum",
						Path:  "repeating_enum",
						Label: labels.Optional,
						Reference: &specs.PropertyReference{
							Resource: template.InputResource,
							Path:     "repeating_enum",
						},
						Template: specs.Template{
							Repeated: &specs.Repeated{
								Property: &specs.Property{
									Template: specs.Template{
										Enum: enum,
									},
								},
							},
						},
					},
					"repeating": {
						Name:  "repeating",
						Path:  "repeating",
						Label: labels.Optional,
						Template: specs.Template{
							Repeated: &specs.Repeated{
								Property: &specs.Property{
									Template: specs.Template{
										Message: &specs.Message{
											Properties: map[string]*specs.Property{
												"value": {
													Name:  "value",
													Path:  "repeating.value",
													Label: labels.Optional,
													Reference: &specs.PropertyReference{
														Resource: template.InputResource,
														Path:     "repeating.value",
													},
													Template: specs.Template{
														Scalar: &specs.Scalar{
															Type: types.String,
														},
													},
												},
											},
										},
									},
								},
							},
						},
						Reference: &specs.PropertyReference{
							Resource: template.InputResource,
							Path:     "repeating",
						},
					},
				},
			},
		},
	},
}

func TestNewConstructor(t *testing.T) {
	result := NewConstructor()
	if result == nil {
		t.Fatal("unexpected nil")
	}
}

func TestConstructorName(t *testing.T) {
	form := NewConstructor()
	if form == nil {
		t.Fatal("unexpected nil")
	}

	expected := "form-urlencoded"
	name := form.Name()
	if name != expected {
		t.Errorf("unexpected name %s, expected %s", name, expected)
	}
}

func TestConstructorNew(t *testing.T) {
	form := NewConstructor()
	if form == nil {
		t.Fatal("unexpected nil")
	}

	manager, err := form.New("mock", &specs.ParameterMap{})
	if err != nil {
		t.Fatal(err)
	}

	if manager == nil {
		t.Error("unexpected empty manager")
	}
}

func TestConstructorNewNil(t *testing.T) {
	form := NewConstructor()
	if form == nil {
		t.Fatal("unexpected nil")
	}

	_, err := form.New("mock", nil)
	if err == nil {
		t.Fatal("unexpected pass")
	}
}

func TestManagerName(t *testing.T) {
	form := NewConstructor()
	if form == nil {
		t.Fatal("unexpected nil")
	}

	manager, err := form.New("mock", &specs.ParameterMap{})
	if err != nil {
		t.Fatal(err)
	}

	expected := "form-urlencoded"
	name := manager.Name()
	if name != expected {
		t.Errorf("unexpected name %s, expected %s", name, expected)
	}
}

func TestManagerProperty(t *testing.T) {
	form := NewConstructor()
	if form == nil {
		t.Fatal("unexpected nil")
	}

	expected := &specs.ParameterMap{
		Property: &specs.Property{},
	}

	manager, err := form.New("mock", expected)
	if err != nil {
		t.Fatal(err)
	}

	result := manager.Property()
	if result != expected.Property {
		t.Errorf("unexpected property %+v, expected %+v", result, expected.Property)
	}
}

func TestMarshalNil(t *testing.T) {
	form := NewConstructor()
	if form == nil {
		t.Fatal("unexpected nil")
	}

	expected := &specs.ParameterMap{}
	manager, err := form.New("mock", expected)
	if err != nil {
		t.Fatal(err)
	}

	r, err := manager.Marshal(nil)
	if err != nil {
		t.Error(err)
	}

	if r == nil {
		t.Errorf("unexpected nil reader")
	}
}

func TestMarshal(t *testing.T) {
	form := NewConstructor()
	if form == nil {
		t.Fatal("unexpected nil")
	}

	type test struct {
		expected string
		input    map[string]interface{}
	}

	tests := map[string]test{
		"simple": {
			expected: "message=hello+world",
			input: map[string]interface{}{
				"message": "hello world",
				"nested":  map[string]interface{}{},
			},
		},
		"nesting": {
			expected: "nested.first=nested+value",
			input: map[string]interface{}{
				"nested": map[string]interface{}{
					"first": "nested value",
				},
			},
		},
		"enum": {
			expected: "status=PENDING",
			input: map[string]interface{}{
				"nested": map[string]interface{}{},
				"status": references.Enum("PENDING", 1),
			},
		},
		"repeating": {
			expected: "repeating%5B0%5D.value=repeating+one&repeating%5B1%5D.value=repeating+two",
			input: map[string]interface{}{
				"repeating": []map[string]interface{}{
					{
						"value": "repeating one",
					},
					{
						"value": "repeating two",
					},
				},
			},
		},
		"repeating_enum": {
			expected: "repeating_enum%5B0%5D=PENDING&repeating_enum%5B1%5D=UNKOWN",
			input: map[string]interface{}{
				"nested": map[string]interface{}{},
				"repeating_enum": []interface{}{
					references.Enum("PENDING", 1),
					references.Enum("UNKNOWN", 0),
				},
			},
		},
		"repeating_string": {
			expected: "repeating_string%5B0%5D=repeating+one&repeating_string%5B1%5D=repeating+two",
			input: map[string]interface{}{
				"repeating_string": []interface{}{
					"repeating one",
					"repeating two",
				},
			},
		},
		"complex": {
			expected: "message=hello+world&nested.first=nested+value&repeating%5B0%5D.value=repeating+value",
			input: map[string]interface{}{
				"message": "hello world",
				"nested": map[string]interface{}{
					"first": "nested value",
				},
				"repeating": []map[string]interface{}{
					{
						"value": "repeating value",
					},
				},
			},
		},
		"empty_repeating": {
			expected: "message=hello+world&nested.first=nested+value&repeating%5B0%5D.value=repeating+value",
			input: map[string]interface{}{
				"message": "hello world",
				"nested": map[string]interface{}{
					"first": "nested value",
				},
				"repeating": []map[string]interface{}{
					{
						"value": "repeating value",
					},
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			manager, err := form.New("mock", schema)
			if err != nil {
				t.Fatal(err)
			}

			refs := references.NewReferenceStore(len(test.input))
			refs.StoreValues(template.InputResource, "", test.input)

			r, err := manager.Marshal(refs)
			if err != nil {
				t.Error(err)
			}

			if r == nil {
				t.Errorf("unexpected nil reader")
			}

			bb, err := ioutil.ReadAll(r)
			if err != nil {
				t.Fatal(err)
			}

			if string(bb) != test.expected {
				t.Errorf("unexpected result %s, expected %s", string(bb), test.expected)
			}
		})
	}
}

type readerFunc func([]byte) (int, error)

func (fn readerFunc) Read(p []byte) (int, error) { return fn(p) }

func TestUnmarshal(t *testing.T) {
	type test struct {
		input    io.Reader
		expected map[string]expect
		error    error
	}

	tests := map[string]test{
		"reader error": {
			input: readerFunc(
				func([]byte) (int, error) {
					return 0, errors.New("failed")
				},
			),
			error: errors.New("failed"),
		},
		"type mismatch": {
			input: strings.NewReader("numeric=foo"),
			error: errors.New(""), // error returned by ParseInt()
		},
		"empty nested schema": {
			input: strings.NewReader("no_nested_schema.value=foo"),
			error: errNilSchema,
		},
		"empty reader": {
			input: strings.NewReader(""),
		},
		"error with undefined property": {
			input: strings.NewReader(
				"undefined=hello+world",
			),
			error: errUndefinedProperty(""),
		},
		"error with unknown label": {
			input: strings.NewReader(
				"bad_label=hello+world",
			),
			error: errUnknownLabel(""),
		},
		"simple": {
			input: strings.NewReader(
				"message=hello+world",
			),
			expected: map[string]expect{
				"message": {
					value: "hello world",
				},
			},
		},
		"nested": {
			input: strings.NewReader(
				"nested.first=nested+first&nested.second=nested+second",
			),
			expected: map[string]expect{
				"nested.first": {
					value: "nested first",
				},
				"nested.second": {
					value: "nested second",
				},
			},
		},
		"enum": {
			input: strings.NewReader(
				"status=PENDING",
			),
			expected: map[string]expect{
				"status": {
					enum: func() *int32 { i := int32(1); return &i }(),
				},
			},
		},
		"repeated string": {
			input: strings.NewReader(
				"repeating_string=repeating+one&repeating_string=repeating+two",
			),
			expected: map[string]expect{
				"repeating_string": {
					repeated: []expect{
						{
							value: "repeating one",
						},
						{
							value: "repeating two",
						},
					},
				},
			},
		},
		"repeated nested": {
			input: strings.NewReader(
				"repeating.value=repeating+one&repeating.value=repeating+two",
			),
			expected: map[string]expect{
				"repeating": {
					repeated: []expect{
						{
							nested: map[string]expect{
								"repeating.value": {
									value: "repeating one",
								},
							},
						},
						{
							nested: map[string]expect{
								"repeating.value": {
									value: "repeating two",
								},
							},
						},
					},
				},
			},
		},
		"repeated enum": {
			input: strings.NewReader(
				"repeating_enum=PENDING&repeating_enum=UNKOWN",
			),
			expected: map[string]expect{
				"repeating_enum": {
					repeated: []expect{
						{
							enum: func() *int32 { i := int32(1); return &i }(),
						},
						{
							enum: func() *int32 { i := int32(0); return &i }(),
						},
					},
				},
			},
		},
		"complex": {
			input: strings.NewReader(
				"message=hello+world&nested.first=nested+value&repeating.value=repeating+value",
			),
			expected: map[string]expect{
				"message": {
					value: "hello world",
				},
				"nested.first": {
					value: "nested value",
				},
				"repeating": {
					repeated: []expect{
						{
							nested: map[string]expect{
								"repeating.value": {
									value: "repeating value",
								},
							},
						},
					},
				},
			},
		},
	}

	for title, test := range tests {
		t.Run(title, func(t *testing.T) {
			form := NewConstructor()
			if form == nil {
				t.Fatal("unexpected nil")
			}

			manager, err := form.New("mock", schema)
			if err != nil {
				t.Fatal(err)
			}

			var refs = references.NewReferenceStore(0)
			err = manager.Unmarshal(test.input, refs)

			if test.error != nil {
				if !errors.As(err, &test.error) {
					t.Errorf("error [%s] was expected to be [%s]", err, test.error)
				}
			} else if err != nil {
				t.Errorf("error was not expected: %s", err)
			}

			for path, output := range test.expected {
				assert(t, "mock", path, refs, output)
			}
		})
	}
}

type expect struct {
	value    interface{}
	enum     *int32
	repeated []expect
	nested   map[string]expect
}

func assert(t *testing.T, resource string, path string, store references.Store, output expect) {
	ref := store.Load(resource, path)

	if ref == nil {
		t.Errorf("reference %q was expected to be set", path)
	}

	if output.value != nil {
		if ref.Value != output.value {
			t.Errorf("reference %q was expected to have value [%v], not [%v]", path, output.value, ref.Value)
		}

		return
	}

	if output.enum != nil {
		if ref.Enum == nil {
			t.Errorf("reference %q was expected to have a enum value", path)
		}

		if *output.enum != *ref.Enum {
			t.Errorf("reference %q was expected to have enum value [%d], not [%d]", path, *output.enum, *ref.Enum)
		}

		return
	}

	if output.repeated != nil {
		if ref.Repeated == nil {
			t.Errorf("reference %q was expected to have a repeated value", path)
		}

		if expected, actual := len(ref.Repeated), len(ref.Repeated); actual != expected {
			t.Errorf("invalid number of repeated values, expected %d, got %d", expected, actual)
		}

		for index, expected := range output.repeated {
			if expected.value != nil || expected.enum != nil {
				assert(t, "", "", ref.Repeated[index], expected)

				continue
			}

			if expected.nested != nil {
				for key, expected := range expected.nested {
					assert(t, resource, key, ref.Repeated[index], expected)
				}

				continue
			}
		}
	}
}
