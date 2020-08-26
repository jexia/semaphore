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
		Type:  types.Message,
		Label: labels.Optional,
		Nested: map[string]*specs.Property{
			"message": {
				Name:  "message",
				Path:  "message",
				Type:  types.String,
				Label: labels.Optional,
				Reference: &specs.PropertyReference{
					Resource: template.InputResource,
					Path:     "message",
				},
			},
			"repeating_values": {
				Name:  "repeating_values",
				Path:  "repeating_values",
				Type:  types.String,
				Label: labels.Repeated,
				Reference: &specs.PropertyReference{
					Resource: template.InputResource,
					Path:     "repeating_values",
				},
			},
			"status": {
				Name:  "status",
				Path:  "status",
				Type:  types.Enum,
				Label: labels.Optional,
				Enum:  enum,
				Reference: &specs.PropertyReference{
					Resource: template.InputResource,
					Path:     "status",
				},
			},
			"repeating_status": {
				Name:  "repeating_status",
				Path:  "repeating_status",
				Type:  types.Enum,
				Label: labels.Repeated,
				Enum:  enum,
				Reference: &specs.PropertyReference{
					Resource: template.InputResource,
					Path:     "repeating_status",
				},
			},
			"nested": {
				Name:  "nested",
				Path:  "nested",
				Type:  types.Message,
				Label: labels.Optional,
				Nested: map[string]*specs.Property{
					"value": {
						Name:  "value",
						Path:  "nested.value",
						Type:  types.String,
						Label: labels.Optional,
						Reference: &specs.PropertyReference{
							Resource: template.InputResource,
							Path:     "nested.value",
						},
					},
				},
			},
			"repeating": {
				Name:  "repeating",
				Path:  "repeating",
				Type:  types.Message,
				Label: labels.Repeated,
				Nested: map[string]*specs.Property{
					"value": {
						Name:  "value",
						Path:  "repeating.value",
						Type:  types.String,
						Label: labels.Optional,
						Reference: &specs.PropertyReference{
							Resource: template.InputResource,
							Path:     "repeating.value",
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
			expected: "nested.value=nested+value",
			input: map[string]interface{}{
				"nested": map[string]interface{}{
					"value": "nested value",
				},
			},
		},
		"repeating": {
			expected: "repeating%5B0%5D.value=repeating+value",
			input: map[string]interface{}{
				"nested": map[string]interface{}{},
				"repeating": []map[string]interface{}{
					{
						"value": "repeating value",
					},
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
		"repeating_enum": {
			expected: "repeating_status%5B0%5D=PENDING&repeating_status%5B1%5D=UNKOWN",
			input: map[string]interface{}{
				"nested": map[string]interface{}{},
				"repeating_status": []interface{}{
					references.Enum("PENDING", 1),
					references.Enum("UNKNOWN", 0),
				},
			},
		},
		"repeating_values": {
			expected: "repeating_values%5B0%5D=repeating+value&repeating_values%5B1%5D=repeating+value",
			input: map[string]interface{}{
				"nested": map[string]interface{}{},
				"repeating_values": []interface{}{
					"repeating value",
					"repeating value",
				},
			},
		},
		"complex": {
			expected: "message=hello+world&nested.value=nested+value&repeating%5B0%5D.value=repeating+value",
			input: map[string]interface{}{
				"message": "hello world",
				"nested": map[string]interface{}{
					"value": "nested value",
				},
				"repeating": []map[string]interface{}{
					{
						"value": "repeating value",
					},
				},
			},
		},
		"empty_repeating": {
			expected: "message=hello+world&nested.value=nested+value&repeating%5B0%5D.value=repeating+value",
			input: map[string]interface{}{
				"message": "hello world",
				"nested": map[string]interface{}{
					"value": "nested value",
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

			refs := references.NewReferenceStore(0)
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
