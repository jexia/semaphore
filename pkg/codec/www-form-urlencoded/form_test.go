package formencoded

import (
	"io/ioutil"
	"testing"

	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/labels"
	"github.com/jexia/semaphore/pkg/specs/template"
	"github.com/jexia/semaphore/pkg/specs/types"
)

var enum = &specs.Enum{
	Keys: map[string]*specs.EnumValue{
		"UNKNOWN": {
			Key:      "UNKNOWN",
			Position: 0,
		},
		"PENDING": {
			Key:      "PENDING",
			Position: 1,
		},
	},
	Positions: map[int32]*specs.EnumValue{
		0: {
			Key:      "UNKNOWN",
			Position: 0,
		},
		1: {
			Key:      "PENDING",
			Position: 1,
		},
	},
}

var schemaWithDefaultRepeating = &specs.ParameterMap{
	Property: &specs.Property{
		Name: "repeating_with_defaults",
		Path: "repeating_with_defaults",
		Template: specs.Template{
			Repeated: specs.Repeated{
				{
					Scalar: &specs.Scalar{
						Default: "yes",
						Type:    types.String,
					},
				},
				{
					Scalar: &specs.Scalar{
						Default: "no",
						Type:    types.String,
					},
				},
			},
		},
	},
}

var schema = &specs.ParameterMap{
	Property: &specs.Property{
		Label: labels.Optional,
		Template: specs.Template{
			Message: specs.Message{
				"bad_label": {
					Name:  "bad_label",
					Path:  "bad_label",
					Label: "",
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
						Message: specs.Message{},
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
					Template: specs.Template{
						Reference: &specs.PropertyReference{
							Resource: template.InputResource,
							Path:     "message",
						},
						Scalar: &specs.Scalar{
							Type: types.String,
						},
					},
				},
				"status": {
					Name:  "status",
					Path:  "status",
					Label: labels.Optional,
					Template: specs.Template{
						Reference: &specs.PropertyReference{
							Resource: template.InputResource,
							Path:     "status",
						},
						Enum: enum,
					},
				},
				"nested": {
					Name:  "nested",
					Path:  "nested",
					Label: labels.Optional,
					Template: specs.Template{
						Message: specs.Message{
							"first": {
								Name:  "first",
								Path:  "nested.first",
								Label: labels.Optional,
								Template: specs.Template{
									Reference: &specs.PropertyReference{
										Resource: template.InputResource,
										Path:     "nested.first",
									},
									Scalar: &specs.Scalar{
										Type: types.String,
									},
								},
							},
							"second": {
								Name:  "second",
								Path:  "nested.second",
								Label: labels.Optional,
								Template: specs.Template{
									Reference: &specs.PropertyReference{
										Resource: template.InputResource,
										Path:     "nested.second",
									},
									Scalar: &specs.Scalar{
										Type: types.String,
									},
								},
							},
						},
					},
				},
				"repeating_string": {
					Name: "repeating_string",
					Path: "repeating_string",
					Template: specs.Template{
						Reference: &specs.PropertyReference{
							Resource: template.InputResource,
							Path:     "repeating_string",
						},
						Repeated: specs.Repeated{
							{
								Scalar: &specs.Scalar{Type: types.String},
							},
						},
					},
				},
				"repeating_enum": {
					Name:  "repeating_enum",
					Path:  "repeating_enum",
					Label: labels.Optional,
					Template: specs.Template{
						Reference: &specs.PropertyReference{
							Resource: template.InputResource,
							Path:     "repeating_enum",
						},
						Repeated: specs.Repeated{
							{Enum: enum},
						},
					},
				},
				"repeating": {
					Name:  "repeating",
					Path:  "repeating",
					Label: labels.Optional,
					Template: specs.Template{
						Reference: &specs.PropertyReference{
							Resource: template.InputResource,
							Path:     "repeating",
						},
						Repeated: specs.Repeated{
							{
								Message: specs.Message{
									"value": {
										Name:  "value",
										Path:  "repeating.value",
										Label: labels.Optional,
										Template: specs.Template{
											Reference: &specs.PropertyReference{
												Resource: template.InputResource,
												Path:     "repeating.value",
											},
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
		t.Errorf("unexpected name %s, want %s", name, expected)
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
		t.Errorf("unexpected name %s, want %s", name, expected)
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
		t.Errorf("unexpected property %+v, want %+v", result, expected.Property)
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
		schema   *specs.ParameterMap // if nil, use the default global `schema`
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
		"repeating_with_defaults": {
			schema:   schemaWithDefaultRepeating,
			expected: "repeating_with_defaults%5B0%5D=yes&repeating_with_defaults%5B1%5D=no",
			input: map[string]interface{}{
				"repeating_with_defaults": []interface{}{},
			},
		},
		"repeating_enum": {
			expected: "repeating_enum%5B0%5D=PENDING&repeating_enum%5B1%5D=UNKNOWN",
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
			sch := schema
			if test.schema != nil {
				sch = test.schema
			}
			manager, err := form.New("mock", sch)
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
				t.Errorf("unexpected result %s, want %s", string(bb), test.expected)
			}
		})
	}
}
