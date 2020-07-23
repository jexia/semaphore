package functions

import (
	"testing"

	"github.com/jexia/semaphore/pkg/core/instance"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/labels"
	"github.com/jexia/semaphore/pkg/specs/types"
)

func TestDefineFunction(t *testing.T) {
	ctx := instance.NewContext()

	type test struct {
		stack Stack
		node  *specs.Node
		flow  specs.FlowsInterface
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
				ID: "second",
			},
			flow: &specs.Flow{
				Nodes: []*specs.Node{
					{
						ID: "first",
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
						ID: "second",
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
				ID: "second",
			},
			flow: &specs.Flow{
				Nodes: []*specs.Node{
					{
						ID: "first",
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
						ID: "second",
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
