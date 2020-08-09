package forwarding

import (
	"testing"

	"github.com/jexia/semaphore/pkg/broker"
	"github.com/jexia/semaphore/pkg/broker/logger"
	"github.com/jexia/semaphore/pkg/specs"
)

func TestResolveReferences(t *testing.T) {
	tests := map[string]func() (specs.FlowListInterface, *specs.Property, *specs.Property){
		"flow": func() (specs.FlowListInterface, *specs.Property, *specs.Property) {
			expected := &specs.Property{
				Reference: &specs.PropertyReference{
					Resource: "expected",
					Path:     "expected",
				},
			}

			target := &specs.Property{
				Reference: &specs.PropertyReference{
					Resource: "mock",
					Path:     "",
					Property: expected,
				},
			}

			flows := specs.FlowListInterface{
				&specs.Flow{
					Nodes: []*specs.Node{
						{
							Call: &specs.Call{
								Request: &specs.ParameterMap{
									Property: target,
								},
							},
						},
					},
				},
			}

			return flows, target, expected
		},
		"proxy": func() (specs.FlowListInterface, *specs.Property, *specs.Property) {
			expected := &specs.Property{
				Reference: &specs.PropertyReference{
					Resource: "expected",
					Path:     "expected",
				},
			}

			target := &specs.Property{
				Reference: &specs.PropertyReference{
					Resource: "mock",
					Path:     "",
					Property: expected,
				},
			}

			flows := specs.FlowListInterface{
				&specs.Proxy{
					Nodes: []*specs.Node{
						{
							Call: &specs.Call{
								Request: &specs.ParameterMap{
									Property: target,
								},
							},
						},
					},
				},
			}

			return flows, target, expected
		},
		"params": func() (specs.FlowListInterface, *specs.Property, *specs.Property) {
			expected := &specs.Property{
				Reference: &specs.PropertyReference{
					Resource: "expected",
					Path:     "expected",
				},
			}

			target := &specs.Property{
				Reference: &specs.PropertyReference{
					Resource: "mock",
					Path:     "",
					Property: expected,
				},
			}

			flows := specs.FlowListInterface{
				&specs.Flow{
					Nodes: []*specs.Node{
						{
							Call: &specs.Call{
								Request: &specs.ParameterMap{
									Params: map[string]*specs.Property{
										"mock": target,
									},
								},
							},
						},
					},
				},
			}

			return flows, target, expected
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := logger.WithLogger(broker.NewContext())
			flows, target, expected := test()
			ResolveReferences(ctx, flows)

			if target.Reference.String() != expected.Reference.String() {
				t.Errorf("unexpected reference '%s', expected '%s'", target.Reference, expected.Reference)
			}
		})
	}
}

func TestResolveNodeReferences(t *testing.T) {
	tests := map[string]func() (*specs.Node, *specs.Property, *specs.Property){
		"call request": func() (*specs.Node, *specs.Property, *specs.Property) {
			expected := &specs.Property{
				Reference: &specs.PropertyReference{
					Resource: "expected",
					Path:     "expected",
				},
			}

			target := &specs.Property{
				Reference: &specs.PropertyReference{
					Resource: "mock",
					Path:     "",
					Property: expected,
				},
			}

			node := &specs.Node{
				Call: &specs.Call{
					Request: &specs.ParameterMap{
						Property: target,
					},
				},
			}

			return node, target, expected
		},
		"call response": func() (*specs.Node, *specs.Property, *specs.Property) {
			expected := &specs.Property{
				Reference: &specs.PropertyReference{
					Resource: "expected",
					Path:     "expected",
				},
			}

			target := &specs.Property{
				Reference: &specs.PropertyReference{
					Resource: "mock",
					Path:     "",
					Property: expected,
				},
			}

			node := &specs.Node{
				Call: &specs.Call{
					Response: &specs.ParameterMap{
						Property: target,
					},
				},
			}

			return node, target, expected
		},
		"rollback request": func() (*specs.Node, *specs.Property, *specs.Property) {
			expected := &specs.Property{
				Reference: &specs.PropertyReference{
					Resource: "expected",
					Path:     "expected",
				},
			}

			target := &specs.Property{
				Reference: &specs.PropertyReference{
					Resource: "mock",
					Path:     "",
					Property: expected,
				},
			}

			node := &specs.Node{
				Rollback: &specs.Call{
					Request: &specs.ParameterMap{
						Property: target,
					},
				},
			}

			return node, target, expected
		},
		"rollback response": func() (*specs.Node, *specs.Property, *specs.Property) {
			expected := &specs.Property{
				Reference: &specs.PropertyReference{
					Resource: "expected",
					Path:     "expected",
				},
			}

			target := &specs.Property{
				Reference: &specs.PropertyReference{
					Resource: "mock",
					Path:     "",
					Property: expected,
				},
			}

			node := &specs.Node{
				Rollback: &specs.Call{
					Response: &specs.ParameterMap{
						Property: target,
					},
				},
			}

			return node, target, expected
		},
		"call nested": func() (*specs.Node, *specs.Property, *specs.Property) {
			expected := &specs.Property{
				Reference: &specs.PropertyReference{
					Resource: "expected",
					Path:     "expected",
				},
				Nested: map[string]*specs.Property{
					"example": {
						Name: "example",
						Path: "example",
					},
				},
			}

			target := &specs.Property{
				Reference: &specs.PropertyReference{
					Resource: "mock",
					Path:     "",
					Property: expected,
				},
			}

			node := &specs.Node{
				Rollback: &specs.Call{
					Response: &specs.ParameterMap{
						Property: &specs.Property{
							Nested: map[string]*specs.Property{
								"example": target,
							},
						},
					},
				},
			}

			return node, target, expected
		},
		"header": func() (*specs.Node, *specs.Property, *specs.Property) {
			expected := &specs.Property{
				Reference: &specs.PropertyReference{
					Resource: "expected",
					Path:     "expected",
				},
			}

			target := &specs.Property{
				Reference: &specs.PropertyReference{
					Resource: "mock",
					Path:     "",
					Property: expected,
				},
			}

			node := &specs.Node{
				Rollback: &specs.Call{
					Response: &specs.ParameterMap{
						Header: specs.Header{
							"example": target,
						},
					},
				},
			}

			return node, target, expected
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			node, target, expected := test()
			ResolveNodeReferences(node)

			if target.Reference == nil {
				t.Fatal("target reference not set")
			}

			if expected.Reference == nil {
				t.Fatal("expected reference not set")
			}

			if target.Reference.String() != expected.Reference.String() {
				t.Errorf("unexpected reference '%s', expected '%s'", target.Reference, expected.Reference)
			}
		})
	}
}

func TestResolveOutputReferences(t *testing.T) {
	tests := map[string]func() (*specs.Flow, *specs.Property, *specs.Property){
		"simple": func() (*specs.Flow, *specs.Property, *specs.Property) {
			expected := &specs.Property{
				Reference: &specs.PropertyReference{
					Resource: "expected",
					Path:     "expected",
				},
			}

			target := &specs.Property{
				Reference: &specs.PropertyReference{
					Resource: "mock",
					Path:     "",
					Property: expected,
				},
			}

			flow := &specs.Flow{
				Output: &specs.ParameterMap{
					Property: target,
				},
			}

			return flow, target, expected
		},
		"nested": func() (*specs.Flow, *specs.Property, *specs.Property) {
			expected := &specs.Property{
				Reference: &specs.PropertyReference{
					Resource: "expected",
					Path:     "expected",
				},
				Nested: map[string]*specs.Property{
					"example": {
						Name: "example",
						Path: "example",
					},
				},
			}

			target := &specs.Property{
				Reference: &specs.PropertyReference{
					Resource: "mock",
					Path:     "",
					Property: expected,
				},
			}

			flow := &specs.Flow{
				Output: &specs.ParameterMap{
					Property: target,
				},
			}

			return flow, target, expected
		},
		"header": func() (*specs.Flow, *specs.Property, *specs.Property) {
			expected := &specs.Property{
				Reference: &specs.PropertyReference{
					Resource: "expected",
					Path:     "expected",
				},
			}

			target := &specs.Property{
				Reference: &specs.PropertyReference{
					Resource: "mock",
					Path:     "",
					Property: expected,
				},
			}

			flow := &specs.Flow{
				Output: &specs.ParameterMap{
					Header: specs.Header{
						"example": target,
					},
				},
			}

			return flow, target, expected
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			flow, target, expected := test()
			ctx := logger.WithLogger(broker.NewContext())
			ResolveReferences(ctx, specs.FlowListInterface{flow})

			if target.Reference == nil {
				t.Fatal("target reference not set")
			}

			if expected.Reference == nil {
				t.Fatal("expected reference not set")
			}

			if target.Reference.String() != expected.Reference.String() {
				t.Errorf("unexpected reference '%s', expected '%s'", target.Reference, expected.Reference)
			}
		})
	}
}
