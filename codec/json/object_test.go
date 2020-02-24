package json

import (
	"os"
	"testing"

	"github.com/francoispqt/gojay"
	"github.com/jexia/maestro/refs"
	"github.com/jexia/maestro/specs"
	"github.com/jexia/maestro/specs/types"
)

func TestSomething(t *testing.T) {
	params := &specs.ParameterMap{
		Properties: map[string]*specs.Property{
			"message": &specs.Property{
				Default: nil,
				Type:    types.TypeString,
			},
		},
		Nested: map[string]*specs.NestedParameterMap{
			"nested": &specs.NestedParameterMap{
				Properties: map[string]*specs.Property{
					"message": &specs.Property{
						Default: "hello",
						Type:    types.TypeString,
					},
				},
			},
		},
		Repeated: map[string]*specs.RepeatedParameterMap{
			"repeated": &specs.RepeatedParameterMap{
				Template: specs.ParsePropertyReference("input:repeated"),
				Properties: map[string]*specs.Property{
					"message": &specs.Property{
						Default:   nil,
						Type:      types.TypeString,
						Reference: specs.ParsePropertyReference("input:repeated.message"),
					},
				},
				Repeated: map[string]*specs.RepeatedParameterMap{
					"repeated": &specs.RepeatedParameterMap{
						Template: specs.ParsePropertyReference("input:repeated.repeated"),
						Properties: map[string]*specs.Property{
							"message": &specs.Property{
								Default:   nil,
								Type:      types.TypeString,
								Reference: specs.ParsePropertyReference("input:repeated.repeated.message"),
							},
						},
					},
				},
			},
		},
	}

	refs := refs.NewStore(1)
	refs.StoreValues("input", "", map[string]interface{}{
		"repeated": []map[string]interface{}{
			map[string]interface{}{
				"message": "some message",
				"repeated": []map[string]interface{}{
					map[string]interface{}{
						"message": "some message",
					},
				},
			},
		},
	})

	object := NewObject(params, refs)
	encoder := gojay.BorrowEncoder(os.Stdout)
	encoder.Encode(object)
}
