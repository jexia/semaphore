package specs

//
// import (
// 	"encoding/json"
// 	"reflect"
// 	"testing"
//
// 	"github.com/jexia/semaphore/pkg/specs/labels"
// 	"github.com/jexia/semaphore/pkg/specs/metadata"
// 	"github.com/jexia/semaphore/pkg/specs/types"
// )
//
// func TestPropertyUnmarshalDefaultPropertyRepeated(t *testing.T) {
// 	type test struct {
// 		input    *Property
// 		expected reflect.Kind
// 	}
//
// 	tests := map[string]test{
// 		"int64": {
// 			input: &Property{
// 				Type:    types.Int64,
// 				Label:   labels.Optional,
// 				Default: 100,
// 			},
// 			expected: reflect.Int64,
// 		},
// 		"sint64": {
// 			input: &Property{
// 				Type:    types.Sint64,
// 				Label:   labels.Optional,
// 				Default: 100,
// 			},
// 			expected: reflect.Int64,
// 		},
// 		"sfixed64": {
// 			input: &Property{
// 				Type:    types.Sfixed64,
// 				Label:   labels.Optional,
// 				Default: 100,
// 			},
// 			expected: reflect.Int64,
// 		},
// 		"uint64": {
// 			input: &Property{
// 				Type:    types.Uint64,
// 				Label:   labels.Optional,
// 				Default: 100,
// 			},
// 			expected: reflect.Uint64,
// 		},
// 		"fixed64": {
// 			input: &Property{
// 				Type:    types.Fixed64,
// 				Label:   labels.Optional,
// 				Default: 100,
// 			},
// 			expected: reflect.Uint64,
// 		},
// 		"int32": {
// 			input: &Property{
// 				Type:    types.Int32,
// 				Label:   labels.Optional,
// 				Default: 100,
// 			},
// 			expected: reflect.Int32,
// 		},
// 		"sint32": {
// 			input: &Property{
// 				Type:    types.Sint32,
// 				Label:   labels.Optional,
// 				Default: 100,
// 			},
// 			expected: reflect.Int32,
// 		},
// 		"sfixed32": {
// 			input: &Property{
// 				Type:    types.Sfixed32,
// 				Label:   labels.Optional,
// 				Default: 100,
// 			},
// 			expected: reflect.Int32,
// 		},
// 		"uint32": {
// 			input: &Property{
// 				Type:    types.Uint32,
// 				Label:   labels.Optional,
// 				Default: 100,
// 			},
// 			expected: reflect.Uint32,
// 		},
// 		"fixed32": {
// 			input: &Property{
// 				Type:    types.Fixed32,
// 				Label:   labels.Optional,
// 				Default: 100,
// 			},
// 			expected: reflect.Uint32,
// 		},
// 	}
//
// 	for name, test := range tests {
// 		t.Run(name, func(t *testing.T) {
// 			input := &Property{
// 				Type:   test.input.Type,
// 				Label:  labels.Repeated,
// 				Nested: PropertyList{test.input},
// 			}
//
// 			encoded, err := json.Marshal(input)
// 			if err != nil {
// 				t.Fatal(err)
// 			}
//
// 			prop := Property{}
// 			err = prop.UnmarshalJSON(encoded)
// 			if err != nil {
// 				t.Fatal(err)
// 			}
//
// 			if prop.Nested == nil || len(prop.Nested) != 1 {
// 				t.Fatalf("unexpected prop.Nested %+v", prop.Nested)
// 			}
//
// 			kind := reflect.TypeOf(prop.Nested[0].Default).Kind()
// 			if kind != test.expected {
// 				t.Errorf("unexpected type %+v, expected %+v", kind, test.expected)
// 			}
// 		})
// 	}
// }
//
// func TestPropertyUnmarshalFail(t *testing.T) {
// 	payload := "non json string"
// 	prop := Property{}
// 	err := prop.UnmarshalJSON([]byte(payload))
// 	if err == nil {
// 		t.Error("expected error got nil")
// 	}
// }
//
// func TestPropertyUnmarshalDefaultProperty(t *testing.T) {
// 	type test struct {
// 		input    *Property
// 		expected reflect.Kind
// 	}
//
// 	tests := map[string]test{
// 		"int64": {
// 			input: &Property{
// 				Type:    types.Int64,
// 				Label:   labels.Optional,
// 				Default: 100,
// 			},
// 			expected: reflect.Int64,
// 		},
// 		"sint64": {
// 			input: &Property{
// 				Type:    types.Sint64,
// 				Label:   labels.Optional,
// 				Default: 100,
// 			},
// 			expected: reflect.Int64,
// 		},
// 		"sfixed64": {
// 			input: &Property{
// 				Type:    types.Sfixed64,
// 				Label:   labels.Optional,
// 				Default: 100,
// 			},
// 			expected: reflect.Int64,
// 		},
// 		"uint64": {
// 			input: &Property{
// 				Type:    types.Uint64,
// 				Label:   labels.Optional,
// 				Default: 100,
// 			},
// 			expected: reflect.Uint64,
// 		},
// 		"fixed64": {
// 			input: &Property{
// 				Type:    types.Fixed64,
// 				Label:   labels.Optional,
// 				Default: 100,
// 			},
// 			expected: reflect.Uint64,
// 		},
// 		"int32": {
// 			input: &Property{
// 				Type:    types.Int32,
// 				Label:   labels.Optional,
// 				Default: 100,
// 			},
// 			expected: reflect.Int32,
// 		},
// 		"sint32": {
// 			input: &Property{
// 				Type:    types.Sint32,
// 				Label:   labels.Optional,
// 				Default: 100,
// 			},
// 			expected: reflect.Int32,
// 		},
// 		"sfixed32": {
// 			input: &Property{
// 				Type:    types.Sfixed32,
// 				Label:   labels.Optional,
// 				Default: 100,
// 			},
// 			expected: reflect.Int32,
// 		},
// 		"uint32": {
// 			input: &Property{
// 				Type:    types.Uint32,
// 				Label:   labels.Optional,
// 				Default: 100,
// 			},
// 			expected: reflect.Uint32,
// 		},
// 		"fixed32": {
// 			input: &Property{
// 				Type:    types.Fixed32,
// 				Label:   labels.Optional,
// 				Default: 100,
// 			},
// 			expected: reflect.Uint32,
// 		},
// 	}
//
// 	for name, test := range tests {
// 		t.Run(name, func(t *testing.T) {
// 			input, err := json.Marshal(test.input)
// 			if err != nil {
// 				t.Fatal(err)
// 			}
//
// 			prop := Property{}
// 			err = prop.UnmarshalJSON(input)
// 			if err != nil {
// 				t.Fatal(err)
// 			}
//
// 			kind := reflect.TypeOf(prop.Default).Kind()
// 			if kind != test.expected {
// 				t.Errorf("unexpected type %+v, expected %+v", kind, test.expected)
// 			}
// 		})
// 	}
// }
//
// func TestPropertyReferenceClone(t *testing.T) {
// 	reference := &PropertyReference{
// 		Meta:     metadata.WithValue(nil, nil, nil),
// 		Resource: "resource",
// 		Path:     "path",
// 		Property: &Property{},
// 	}
//
// 	result := reference.Clone()
// 	if result == nil {
// 		t.Error("unexpected result, expected property reference clone to be returned")
// 	}
//
// 	if result.Meta != reference.Meta {
// 		t.Errorf("unexpected meta %+v, expected %+v", result.Meta, reference.Meta)
// 	}
//
// 	if result.Resource != reference.Resource {
// 		t.Errorf("unexpected resource %+v", result.Resource)
// 	}
//
// 	if result.Path != reference.Path {
// 		t.Errorf("unexpected path %+v", result.Path)
// 	}
//
// 	if result.Property != nil {
// 		t.Errorf("unexpected property %+v", result.Property)
// 	}
// }
//
// func TestPropertyReferenceCloneNilValue(t *testing.T) {
// 	var reference *PropertyReference
// 	result := reference.Clone()
// 	if result != nil {
// 		t.Errorf("unexpected result %+v", result)
// 	}
// }
//
// func TestPropertyReferenceString(t *testing.T) {
// 	t.Parallel()
//
// 	tests := map[string]*PropertyReference{
// 		"resource:path": {
// 			Resource: "resource",
// 			Path:     "path",
// 		},
// 		"resource:nested.path": {
// 			Resource: "resource",
// 			Path:     "nested.path",
// 		},
// 		"resource.prop:path": {
// 			Resource: "resource.prop",
// 			Path:     "path",
// 		},
// 		":": {},
// 		"":  nil,
// 	}
//
// 	for expected, reference := range tests {
// 		t.Run(expected, func(t *testing.T) {
// 			result := reference.String()
// 			if result != expected {
// 				t.Fatalf("unexpected result %s, expected %s", result, expected)
// 			}
// 		})
// 	}
// }
//
// func TestObjectsAppend(t *testing.T) {
// 	objects := Schemas{}
//
// 	objects.Append(Schemas{
// 		"first":  &Property{},
// 		"second": &Property{},
// 	})
//
// 	if len(objects) != 2 {
// 		t.Fatalf("unexpected length %+v, expected 2", len(objects))
// 	}
// }
//
// func TestObjectsAppendNilValue(t *testing.T) {
// 	var objects Schemas
// 	objects.Append(nil)
// }
//
// func TestObjectsGet(t *testing.T) {
// 	objects := Schemas{
// 		"first":  &Property{},
// 		"second": &Property{},
// 	}
//
// 	result := objects.Get("second")
// 	if result == nil {
// 		t.Fatal("unexpected empty result")
// 	}
// }
//
// func TestObjectsGetNilValue(t *testing.T) {
// 	var objects Schemas
// 	result := objects.Get("unknown")
// 	if result != nil {
// 		t.Fatalf("unexpected result %+v", result)
// 	}
// }
//
// func TestPropertyClone(t *testing.T) {
// 	property := &Property{
// 		Meta:      metadata.WithValue(nil, nil, nil),
// 		Position:  1,
// 		Comment:   "sample",
// 		Name:      "first",
// 		Path:      "path",
// 		Default:   false,
// 		Type:      types.String,
// 		Label:     labels.Optional,
// 		Reference: &PropertyReference{},
// 		Nested: []*Property{
// 			{Path: "first"},
// 		},
// 		Raw: "first",
// 		Options: Options{
// 			"sample": "option",
// 		},
// 		Enum: &Enum{},
// 	}
//
// 	result := property.Clone()
// 	if result == nil {
// 		t.Error("unexpected result, expected property reference clone to be returned")
// 	}
//
// 	if result.Meta != property.Meta {
// 		t.Errorf("unexpected meta %+v, expected %+v", result.Meta, property.Meta)
// 	}
//
// 	if result.Position != property.Position {
// 		t.Errorf("unexpected position %+v", result.Position)
// 	}
//
// 	if result.Comment != property.Comment {
// 		t.Errorf("unexpected comment %+v", result.Comment)
// 	}
//
// 	if result.Name != property.Name {
// 		t.Errorf("unexpected name %+v", result.Name)
// 	}
//
// 	if result.Path != property.Path {
// 		t.Errorf("unexpected path %+v", result.Path)
// 	}
//
// 	if result.Default != property.Default {
// 		t.Errorf("unexpected default %+v", result.Default)
// 	}
//
// 	if result.Type != property.Type {
// 		t.Errorf("unexpected type %+v", result.Type)
// 	}
//
// 	if result.Label != property.Label {
// 		t.Errorf("unexpected label %+v", result.Label)
// 	}
//
// 	if result.Reference == nil || result.Reference == property.Reference {
// 		t.Errorf("unexpected reference %+v", result.Reference)
// 	}
//
// 	if result.Nested == nil || len(result.Nested) != len(property.Nested) {
// 		t.Errorf("unexpected repeated %+v", result.Nested)
// 	}
//
// 	if result.Raw != property.Raw {
// 		t.Errorf("unexpected raw %+v", result.Raw)
// 	}
//
// 	if len(result.Options) != len(property.Options) {
// 		t.Errorf("unexpected options %+v", result.Options)
// 	}
//
// 	if result.Enum != property.Enum {
// 		t.Errorf("unexpected enum %+v", result.Enum)
// 	}
// }
//
// func TestPropertyListGet(t *testing.T) {
// 	list := PropertyList{
// 		&Property{Name: "first"},
// 		&Property{Name: "second"},
// 	}
//
// 	result := list.Get("second")
// 	if result == nil {
// 		t.Fatal("unexpected empty result when looking up second")
// 	}
//
// 	unexpected := list.Get("unexpected")
// 	if unexpected != nil {
// 		t.Fatal("unexpected lookup returned a unexpected property")
// 	}
// }
//
// func TestPropertyListGetNil(t *testing.T) {
// 	list := PropertyList{
// 		nil,
// 		&Property{Name: "first"},
// 		nil,
// 		&Property{Name: "second"},
// 		nil,
// 	}
//
// 	result := list.Get("second")
// 	if result == nil {
// 		t.Fatal("unexpected empty result when looking up second")
// 	}
//
// 	unexpected := list.Get("unexpected")
// 	if unexpected != nil {
// 		t.Fatal("unexpected lookup returned a unexpected property")
// 	}
// }
//
// func TestParameterMapClone(t *testing.T) {
// 	property := &ParameterMap{
// 		Meta:   metadata.WithValue(nil, nil, nil),
// 		Schema: "com.schema",
// 		Params: map[string]*Property{
// 			"sample": {},
// 		},
// 		Options: Options{
// 			"sample": "option",
// 		},
// 		Header: Header{
// 			"sample": {},
// 		},
// 		Property: &Property{},
// 		Stack: map[string]*Property{
// 			"hash": {},
// 		},
// 	}
//
// 	result := property.Clone()
// 	if result == nil {
// 		t.Error("unexpected result, expected property reference clone to be returned")
// 	}
//
// 	if result.Meta != property.Meta {
// 		t.Errorf("unexpected meta %+v, expected %+v", result.Meta, property.Meta)
// 	}
//
// 	if result.Schema != property.Schema {
// 		t.Errorf("unexpected schema %+v", result.Schema)
// 	}
//
// 	if result.Property == nil || result.Property == property.Property {
// 		t.Errorf("unexpected property %+v", result.Property)
// 	}
//
// 	if len(result.Options) != len(property.Options) {
// 		t.Errorf("unexpected options %+v", result.Options)
// 	}
//
// 	if len(result.Header) != len(property.Header) {
// 		t.Errorf("unexpected header %+v", result.Header)
// 	}
//
// 	if len(result.Stack) != len(property.Stack) {
// 		t.Errorf("unexpected stack %+v", result.Stack)
// 	}
// }
//
// func TestParameterMapCloneNilValue(t *testing.T) {
// 	var params *ParameterMap
//
// 	result := params.Clone()
// 	if result != nil {
// 		t.Errorf("unexpected result %+v", result)
// 	}
// }
