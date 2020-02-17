package specs

import (
	"testing"

	"github.com/jexia/maestro/schema"
	"github.com/jexia/maestro/schema/mock"
	"github.com/jexia/maestro/specs/types"
)

func NewMockObject() schema.Object {
	return &mock.Object{
		Fields: map[string]*mock.Field{
			"message": &mock.Field{
				Name:  "message",
				Type:  types.TypeString,
				Label: types.LabelOptional,
			},
			"nested": &mock.Field{
				Name:  "nested",
				Type:  types.TypeMessage,
				Label: types.LabelOptional,
				Object: &mock.Object{
					Fields: map[string]*mock.Field{
						"key": &mock.Field{
							Name:  "key",
							Type:  types.TypeString,
							Label: types.LabelOptional,
						},
					},
				},
			},
			"repeated": &mock.Field{
				Name:  "repeated",
				Type:  types.TypeMessage,
				Label: types.LabelRepeated,
				Object: &mock.Object{
					Fields: map[string]*mock.Field{
						"key": &mock.Field{
							Name:  "key",
							Type:  types.TypeString,
							Label: types.LabelOptional,
						},
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

	if len(result.Nested) != 1 {
		t.Fatal("no nested parameter maps defined")
	}

	if result.Nested["nested"] == nil {
		t.Fatal("nested value not defined")
	}

	if result.Nested["nested"].Properties["key"] == nil {
		t.Fatal("nested value property not defined")
	}

	if len(result.Repeated) != 1 {
		t.Fatal("no repeated parameter maps defined")
	}

	if result.Repeated["repeated"] == nil {
		t.Fatal("repeated value not defined")
	}

	if result.Repeated["repeated"].Properties["key"] == nil {
		t.Fatal("repeated value property not defined")
	}
}
