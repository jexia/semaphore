package functions

import (
	"testing"

	"github.com/jexia/maestro/pkg/core/instance"
	"github.com/jexia/maestro/pkg/specs"
	"github.com/jexia/maestro/pkg/specs/labels"
	"github.com/jexia/maestro/pkg/specs/types"
)

func TestDefineFunction(t *testing.T) {
	ctx := instance.NewContext()

	type test struct {
		stack Stack
		node  *specs.Node
		flow  specs.FlowResourceManager
	}

	tests := map[string]*test{
		"arguments": {
			stack: Stack{
				"sample": &Function{
					Arguments: []*specs.Property{
						{
							Name: "name",
							Reference: &specs.PropertyReference{
								Resource: "first",
								Path:     "name",
							},
						},
					},
				},
			},
			node: &specs.Node{
				Name: "second",
			},
			flow: &specs.Flow{
				Nodes: []*specs.Node{
					{
						Name: "first",
						Call: &specs.Call{
							Response: &specs.ParameterMap{
								Property: &specs.Property{
									Type:  types.Message,
									Label: labels.Optional,
									Nested: map[string]*specs.Property{
										"name": {
											Name:  "name",
											Path:  "name",
											Type:  types.String,
											Label: labels.Optional,
										},
									},
								},
							},
						},
					},
					{
						Name: "second",
					},
				},
			},
		},
		"returns": {
			stack: Stack{
				"sample": &Function{
					Returns: &specs.Property{
						Reference: &specs.PropertyReference{
							Resource: "first",
							Path:     "name",
						},
					},
				},
			},
			node: &specs.Node{
				Name: "second",
			},
			flow: &specs.Flow{
				Nodes: []*specs.Node{
					{
						Name: "first",
						Call: &specs.Call{
							Response: &specs.ParameterMap{
								Property: &specs.Property{
									Type:  types.Message,
									Label: labels.Optional,
									Nested: map[string]*specs.Property{
										"name": {
											Name:  "name",
											Path:  "name",
											Type:  types.String,
											Label: labels.Optional,
										},
									},
								},
							},
						},
					},
					{
						Name: "second",
					},
				},
			},
		},
		"empty": {
			stack: nil,
			node:  nil,
			flow: &specs.Flow{
				Nodes: []*specs.Node{},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := DefineFunctions(ctx, test.stack, test.node, test.flow)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}
