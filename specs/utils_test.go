package specs

import (
	"testing"

	"github.com/jexia/maestro/schema"
	"github.com/jexia/maestro/schema/mock"
	"github.com/jexia/maestro/specs/types"
)

func NewMockObject() schema.Property {
	return &mock.Property{
		Type: types.TypeMessage,
		Nested: map[string]*mock.Property{
			"message": &mock.Property{
				Name:  "message",
				Type:  types.TypeString,
				Label: types.LabelOptional,
			},
			"nested": &mock.Property{
				Name:  "nested",
				Type:  types.TypeMessage,
				Label: types.LabelOptional,
				Nested: map[string]*mock.Property{
					"key": &mock.Property{
						Name:  "key",
						Type:  types.TypeString,
						Label: types.LabelOptional,
					},
				},
			},
			"repeated": &mock.Property{
				Name:  "repeated",
				Type:  types.TypeMessage,
				Label: types.LabelRepeated,
				Nested: map[string]*mock.Property{
					"key": &mock.Property{
						Name:  "key",
						Type:  types.TypeString,
						Label: types.LabelOptional,
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
				Type:    types.TypeString,
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
