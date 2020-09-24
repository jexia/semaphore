package functions

import (
	"testing"

	"github.com/jexia/semaphore/pkg/broker"
	"github.com/jexia/semaphore/pkg/broker/logger"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/labels"
	"github.com/jexia/semaphore/pkg/specs/types"
)

func TestDefineFunction(t *testing.T) {
	ctx := logger.WithLogger(broker.NewBackground())

	type test struct {
		stack Stack
		node  *specs.Node
		flow  specs.FlowInterface
	}

	tests := map[string]*test{
		"arguments": {
			stack: Stack{
				"sample": &Function{
					Arguments: []*specs.Property{
						{
							Name: "name",
							Template: specs.Template{
								Reference: &specs.PropertyReference{
									Resource: "first",
									Path:     "name",
								},
							},
						},
					},
				},
			},
			node: &specs.Node{
				ID: "second",
			},
			flow: &specs.Flow{
				Nodes: specs.NodeList{
					{
						ID: "first",
						Call: &specs.Call{
							Response: &specs.ParameterMap{
								Property: &specs.Property{
									Label: labels.Optional,
									Template: specs.Template{
										Message: specs.Message{
											"name": &specs.Property{
												Name:  "name",
												Path:  "name",
												Label: labels.Optional,
												Template: specs.Template{
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
						Template: specs.Template{
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
				Nodes: specs.NodeList{
					{
						ID: "first",
						Call: &specs.Call{
							Response: &specs.ParameterMap{
								Property: &specs.Property{
									Label: labels.Optional,
									Template: specs.Template{
										Message: specs.Message{
											"name": &specs.Property{
												Name:  "name",
												Path:  "name",
												Label: labels.Optional,
												Template: specs.Template{
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
				Nodes: specs.NodeList{},
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
