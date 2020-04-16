package strict

import (
	"testing"

	"github.com/jexia/maestro/pkg/specs"
	"github.com/jexia/maestro/pkg/specs/labels"
	"github.com/jexia/maestro/pkg/specs/types"
)

func NewMockObject() *specs.Property {
	return &specs.Property{
		Type: types.Message,
		Nested: map[string]*specs.Property{
			"message": {
				Name:  "message",
				Type:  types.String,
				Label: labels.Optional,
			},
			"nested": {
				Name:  "nested",
				Type:  types.Message,
				Label: labels.Optional,
				Nested: map[string]*specs.Property{
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
				Nested: map[string]*specs.Property{
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
	origin := &specs.ParameterMap{
		Options: specs.Options{
			"version": "1.0",
		},
		Header: specs.Header{
			"cookie": &specs.Property{
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
