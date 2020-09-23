package specs

import (
	"encoding/json"
	"reflect"
	"sort"
	"testing"

	"github.com/jexia/semaphore/pkg/specs/labels"
	"github.com/jexia/semaphore/pkg/specs/metadata"
	"github.com/jexia/semaphore/pkg/specs/types"
)

func TestScalarUnmarshalInvalidJSON(t *testing.T) {
	payload := "non json string"
	prop := Property{
		Template: Template{
			Scalar: &Scalar{},
		},
	}

	err := prop.Scalar.UnmarshalJSON([]byte(payload))
	if err == nil {
		t.Error("expected error got nil")
	}
}

func TestScalarUnmarshalNil(t *testing.T) {
	var scalar *Scalar
	err := scalar.UnmarshalJSON(nil)
	if err != nil {
		t.Error(err)
	}
}

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
						Repeated: &Repeated{
							Default: map[uint]*Property{
								0: test.input,
							},
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

				property := result.Repeated.Default[0]
				if property == nil {
					t.Fatalf("result repeated default property %s is not set", name)
				}

				kind := reflect.TypeOf(property.Scalar.Default).Kind()
				if kind != test.expected {
					t.Errorf("unexpected type %+v, expected %+v", kind, test.expected)
				}
			})
		})
	}
}

func TestPropertyReferenceClone(t *testing.T) {
	reference := &PropertyReference{
		Meta:     metadata.WithValue(nil, nil, nil),
		Resource: "resource",
		Path:     "path",
		Property: &Property{},
	}

	result := reference.Clone()
	if result == nil {
		t.Error("unexpected result, expected property reference clone to be returned")
	}

	if result.Meta != reference.Meta {
		t.Errorf("unexpected meta %+v, expected %+v", result.Meta, reference.Meta)
	}

	if result.Resource != reference.Resource {
		t.Errorf("unexpected resource %+v", result.Resource)
	}

	if result.Path != reference.Path {
		t.Errorf("unexpected path %+v", result.Path)
	}

	if result.Property != nil {
		t.Errorf("unexpected property %+v", result.Property)
	}
}

func TestPropertyReferenceCloneNilValue(t *testing.T) {
	var reference *PropertyReference
	result := reference.Clone()
	if result != nil {
		t.Errorf("unexpected result %+v", result)
	}
}

func TestPropertyReferenceString(t *testing.T) {
	t.Parallel()

	tests := map[string]*PropertyReference{
		"resource:path": {
			Resource: "resource",
			Path:     "path",
		},
		"resource:nested.path": {
			Resource: "resource",
			Path:     "nested.path",
		},
		"resource.prop:path": {
			Resource: "resource.prop",
			Path:     "path",
		},
		":": {},
		"":  nil,
	}

	for expected, reference := range tests {
		t.Run(expected, func(t *testing.T) {
			result := reference.String()
			if result != expected {
				t.Fatalf("unexpected result %s, expected %s", result, expected)
			}
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
			Scalar: &Scalar{
				Default: false,
				Type:    types.String,
			},
			Repeated: &Repeated{
				Default: map[uint]*Property{},
			},
			Message: Message{
				"first": {Name: "first", Path: "first"},
			},
			Enum: &Enum{
				Name: "unknown",
			},
		},
		Reference: &PropertyReference{},
		Raw:       "first",
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

	if result.Repeated.Default == nil {
		t.Fatalf("repeated default not set")
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

func TestPropertyListSort(t *testing.T) {
	list := PropertyList{
		&Property{Name: "third", Position: 2},
		&Property{Name: "first", Position: 0},
		&Property{Name: "second", Position: 1},
	}

	sort.Sort(list)

	for index, item := range list {
		if int(item.Position) != index {
			t.Fatalf("unexpected property list order %d, expected %d", item.Position, index)
		}
	}
}

func TestPropertyListGet(t *testing.T) {
	list := PropertyList{
		&Property{Name: "first"},
		&Property{Name: "second"},
	}

	result := list.Get("second")
	if result == nil {
		t.Fatal("unexpected empty result when looking up second")
	}

	unexpected := list.Get("unexpected")
	if unexpected != nil {
		t.Fatal("unexpected lookup returned a unexpected property")
	}
}

func TestPropertyListGetNil(t *testing.T) {
	list := PropertyList{
		nil,
		&Property{Name: "first"},
		nil,
		&Property{Name: "second"},
		nil,
	}

	result := list.Get("second")
	if result == nil {
		t.Fatal("unexpected empty result when looking up second")
	}

	unexpected := list.Get("unexpected")
	if unexpected != nil {
		t.Fatal("unexpected lookup returned a unexpected property")
	}
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
