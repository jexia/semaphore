package specs

import (
	"testing"

	"github.com/jexia/maestro/schema"
	"github.com/jexia/maestro/schema/mock"
	"github.com/jexia/maestro/specs/labels"
	"github.com/jexia/maestro/specs/types"
)

func NewMockObject() schema.Property {
	return &mock.Property{
		Type: types.Message,
		Nested: map[string]*mock.Property{
			"message": {
				Name:  "message",
				Type:  types.String,
				Label: labels.Optional,
			},
			"nested": {
				Name:  "nested",
				Type:  types.Message,
				Label: labels.Optional,
				Nested: map[string]*mock.Property{
					"key": {
						Name:  "key",
						Type:  types.String,
						Label: labels.Optional,
					},
				},
			},
			"repeated": {
				Name:  "repeated",
				Type:  types.Message,
				Label: labels.Repeated,
				Nested: map[string]*mock.Property{
					"key": {
						Name:  "key",
						Type:  types.String,
						Label: labels.Optional,
					},
				},
			},
		},
	}
}

func TestToParameterMap(t *testing.T) {
	origin := &ParameterMap{
		Options: Options{
			"version": "1.0",
		},
		Header: Header{
			"cookie": &Property{
				Path:    "cookie",
				Default: "mnomnom",
				Type:    types.String,
			},
		},
	}

	path := ""
	object := NewMockObject()

	result := ToParameterMap(origin, path, object)
	if result == nil {
		t.Fatal("result empty")
	}

	if result.Property.Nested["nested"] == nil {
		t.Fatal("nested value not defined")
	}

	if len(result.Property.Nested) != 3 {
		t.Fatalf("expected 3 parameters to be defined received %d", len(result.Property.Nested))
	}

	if result.Property.Nested["nested"].Nested["key"] == nil {
		t.Fatal("nested value property not defined")
	}

	if result.Property.Nested["repeated"] == nil {
		t.Fatal("repeated value not defined")
	}

	if result.Property.Nested["repeated"].Nested["key"] == nil {
		t.Fatal("repeated value property not defined")
	}
}
