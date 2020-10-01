package forwarding

import (
	"testing"

	"github.com/jexia/semaphore/pkg/broker"
	"github.com/jexia/semaphore/pkg/broker/logger"
	"github.com/jexia/semaphore/pkg/functions"
	"github.com/jexia/semaphore/pkg/specs"
)

func TestResolveReferences(t *testing.T) {
	tests := map[string]func() (specs.FlowListInterface, *specs.Property, *specs.Property, functions.Collection){
		"flow": func() (specs.FlowListInterface, *specs.Property, *specs.Property, functions.Collection) {
			expected := &specs.Property{
				Template: specs.Template{
					Reference: &specs.PropertyReference{
						Resource: "expected",
						Path:     "expected",
					},
				},
			}

			target := &specs.Property{
				Template: specs.Template{
					Reference: &specs.PropertyReference{
						Resource: "mock",
						Path:     "",
						Property: expected,
					},
				},
			}

			flows := specs.FlowListInterface{
				&specs.Flow{
					Nodes: specs.NodeList{
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

			return flows, target, expected, functions.Collection{}
		},
		"intermediate": func() (specs.FlowListInterface, *specs.Property, *specs.Property, functions.Collection) {
			expected := &specs.Property{
				Template: specs.Template{
					Reference: &specs.PropertyReference{
						Resource: "expected",
						Path:     "expected",
					},
				},
			}

			target := &specs.Property{
				Template: specs.Template{
					Reference: &specs.PropertyReference{
						Resource: "mock",
						Path:     "",
						Property: expected,
					},
				},
			}

			flows := specs.FlowListInterface{
				&specs.Flow{
					Nodes: specs.NodeList{
						{
							Intermediate: &specs.ParameterMap{
								Property: target,
							},
						},
					},
				},
			}

			return flows, target, expected, functions.Collection{}
		},
		"intermediate function": func() (specs.FlowListInterface, *specs.Property, *specs.Property, functions.Collection) {
			expected := &specs.Property{
				Template: specs.Template{
					Reference: &specs.PropertyReference{
						Resource: "expected",
						Path:     "expected",
					},
				},
			}

			target := &specs.Property{
				Template: specs.Template{
					Reference: &specs.PropertyReference{
						Resource: "mock",
						Path:     "",
						Property: expected,
					},
				},
			}

			param := &specs.ParameterMap{}
			mem := functions.Collection{}
			stack := mem.Reserve(param)

			stack["hash"] = &functions.Function{
				Arguments: []*specs.Property{
					target,
				},
			}

			flows := specs.FlowListInterface{
				&specs.Flow{
					Nodes: specs.NodeList{
						{
							Intermediate: param,
						},
					},
				},
			}

			return flows, target, expected, mem
		},
		"call function": func() (specs.FlowListInterface, *specs.Property, *specs.Property, functions.Collection) {
			expected := &specs.Property{
				Template: specs.Template{
					Reference: &specs.PropertyReference{
						Resource: "expected",
						Path:     "expected",
					},
				},
			}

			target := &specs.Property{
				Template: specs.Template{
					Reference: &specs.PropertyReference{
						Resource: "mock",
						Path:     "",
						Property: expected,
					},
				},
			}

			param := &specs.ParameterMap{}
			mem := functions.Collection{}
			stack := mem.Reserve(param)

			stack["hash"] = &functions.Function{
				Arguments: []*specs.Property{
					target,
				},
			}

			flows := specs.FlowListInterface{
				&specs.Flow{
					Nodes: specs.NodeList{
						{
							Call: &specs.Call{
								Request: param,
							},
						},
					},
				},
			}

			return flows, target, expected, mem
		},
		"condition": func() (specs.FlowListInterface, *specs.Property, *specs.Property, functions.Collection) {
			expected := &specs.Property{
				Template: specs.Template{
					Reference: &specs.PropertyReference{
						Resource: "expected",
						Path:     "expected",
					},
				},
			}

			target := &specs.Property{
				Template: specs.Template{
					Reference: &specs.PropertyReference{
						Resource: "mock",
						Path:     "",
						Property: expected,
					},
				},
			}

			flows := specs.FlowListInterface{
				&specs.Flow{
					Nodes: specs.NodeList{
						{
							Condition: &specs.Condition{
								Params: &specs.ParameterMap{
									Params: map[string]*specs.Property{
										"mock": target,
									},
								},
							},
						},
					},
				},
			}

			return flows, target, expected, functions.Collection{}
		},
		"flow on_error property": func() (specs.FlowListInterface, *specs.Property, *specs.Property, functions.Collection) {
			expected := &specs.Property{
				Template: specs.Template{
					Reference: &specs.PropertyReference{
						Resource: "expected",
						Path:     "expected",
					},
				},
			}

			target := &specs.Property{
				Template: specs.Template{
					Reference: &specs.PropertyReference{
						Resource: "mock",
						Path:     "",
						Property: expected,
					},
				},
			}

			flows := specs.FlowListInterface{
				&specs.Flow{
					OnError: &specs.OnError{
						Response: &specs.ParameterMap{
							Property: target,
						},
					},
				},
			}

			return flows, target, expected, functions.Collection{}
		},
		"flow on_error header": func() (specs.FlowListInterface, *specs.Property, *specs.Property, functions.Collection) {
			expected := &specs.Property{
				Template: specs.Template{
					Reference: &specs.PropertyReference{
						Resource: "expected",
						Path:     "expected",
					},
				},
			}

			target := &specs.Property{
				Template: specs.Template{
					Reference: &specs.PropertyReference{
						Resource: "mock",
						Path:     "",
						Property: expected,
					},
				},
			}

			flows := specs.FlowListInterface{
				&specs.Flow{
					OnError: &specs.OnError{
						Response: &specs.ParameterMap{
							Header: specs.Header{
								"mock": target,
							},
						},
					},
				},
			}

			return flows, target, expected, functions.Collection{}
		},
		"flow on_error param": func() (specs.FlowListInterface, *specs.Property, *specs.Property, functions.Collection) {
			expected := &specs.Property{
				Template: specs.Template{
					Reference: &specs.PropertyReference{
						Resource: "expected",
						Path:     "expected",
					},
				},
			}

			target := &specs.Property{
				Template: specs.Template{
					Reference: &specs.PropertyReference{
						Resource: "mock",
						Path:     "",
						Property: expected,
					},
				},
			}

			flows := specs.FlowListInterface{
				&specs.Flow{
					OnError: &specs.OnError{
						Response: &specs.ParameterMap{
							Params: map[string]*specs.Property{
								"mock": target,
							},
						},
					},
				},
			}

			return flows, target, expected, functions.Collection{}
		},
		"node on_error property": func() (specs.FlowListInterface, *specs.Property, *specs.Property, functions.Collection) {
			expected := &specs.Property{
				Template: specs.Template{
					Reference: &specs.PropertyReference{
						Resource: "expected",
						Path:     "expected",
					},
				},
			}

			target := &specs.Property{
				Template: specs.Template{
					Reference: &specs.PropertyReference{
						Resource: "mock",
						Path:     "",
						Property: expected,
					},
				},
			}

			flows := specs.FlowListInterface{
				&specs.Flow{
					Nodes: specs.NodeList{
						{
							OnError: &specs.OnError{
								Response: &specs.ParameterMap{
									Property: target,
								},
							},
						},
					},
				},
			}

			return flows, target, expected, functions.Collection{}
		},
		"node on_error header": func() (specs.FlowListInterface, *specs.Property, *specs.Property, functions.Collection) {
			expected := &specs.Property{
				Template: specs.Template{
					Reference: &specs.PropertyReference{
						Resource: "expected",
						Path:     "expected",
					},
				},
			}

			target := &specs.Property{
				Template: specs.Template{
					Reference: &specs.PropertyReference{
						Resource: "mock",
						Path:     "",
						Property: expected,
					},
				},
			}

			flows := specs.FlowListInterface{
				&specs.Flow{
					Nodes: specs.NodeList{
						{
							OnError: &specs.OnError{
								Response: &specs.ParameterMap{
									Header: specs.Header{
										"mock": target,
									},
								},
							},
						},
					},
				},
			}

			return flows, target, expected, functions.Collection{}
		},
		"node on_error param": func() (specs.FlowListInterface, *specs.Property, *specs.Property, functions.Collection) {
			expected := &specs.Property{
				Template: specs.Template{
					Reference: &specs.PropertyReference{
						Resource: "expected",
						Path:     "expected",
					},
				},
			}

			target := &specs.Property{
				Template: specs.Template{
					Reference: &specs.PropertyReference{
						Resource: "mock",
						Path:     "",
						Property: expected,
					},
				},
			}

			flows := specs.FlowListInterface{
				&specs.Flow{
					Nodes: specs.NodeList{
						{
							OnError: &specs.OnError{
								Response: &specs.ParameterMap{
									Params: map[string]*specs.Property{
										"mock": target,
									},
								},
							},
						},
					},
				},
			}

			return flows, target, expected, functions.Collection{}
		},
		"proxy": func() (specs.FlowListInterface, *specs.Property, *specs.Property, functions.Collection) {
			expected := &specs.Property{
				Template: specs.Template{
					Reference: &specs.PropertyReference{
						Resource: "expected",
						Path:     "expected",
					},
				},
			}

			target := &specs.Property{
				Template: specs.Template{
					Reference: &specs.PropertyReference{
						Resource: "mock",
						Path:     "",
						Property: expected,
					},
				},
			}

			flows := specs.FlowListInterface{
				&specs.Proxy{
					Nodes: specs.NodeList{
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

			return flows, target, expected, functions.Collection{}
		},
		"params": func() (specs.FlowListInterface, *specs.Property, *specs.Property, functions.Collection) {
			expected := &specs.Property{
				Template: specs.Template{
					Reference: &specs.PropertyReference{
						Resource: "expected",
						Path:     "expected",
					},
				},
			}

			target := &specs.Property{
				Template: specs.Template{
					Reference: &specs.PropertyReference{
						Resource: "mock",
						Path:     "",
						Property: expected,
					},
				},
			}

			flows := specs.FlowListInterface{
				&specs.Flow{
					Nodes: specs.NodeList{
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

			return flows, target, expected, functions.Collection{}
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := logger.WithLogger(broker.NewBackground())
			flows, target, expected, mem := test()
			ResolveReferences(ctx, flows, mem)

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
				Template: specs.Template{
					Reference: &specs.PropertyReference{
						Resource: "expected",
						Path:     "expected",
					},
				},
			}

			target := &specs.Property{
				Template: specs.Template{
					Reference: &specs.PropertyReference{
						Resource: "mock",
						Path:     "",
						Property: expected,
					},
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
				Template: specs.Template{
					Reference: &specs.PropertyReference{
						Resource: "expected",
						Path:     "expected",
					},
				},
			}

			target := &specs.Property{
				Template: specs.Template{
					Reference: &specs.PropertyReference{
						Resource: "mock",
						Path:     "",
						Property: expected,
					},
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
				Template: specs.Template{
					Reference: &specs.PropertyReference{
						Resource: "expected",
						Path:     "expected",
					},
				},
			}

			target := &specs.Property{
				Template: specs.Template{
					Reference: &specs.PropertyReference{
						Resource: "mock",
						Path:     "",
						Property: expected,
					},
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
				Template: specs.Template{
					Reference: &specs.PropertyReference{
						Resource: "expected",
						Path:     "expected",
					},
				},
			}

			target := &specs.Property{
				Template: specs.Template{
					Reference: &specs.PropertyReference{
						Resource: "mock",
						Path:     "",
						Property: expected,
					},
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
				Template: specs.Template{
					Reference: &specs.PropertyReference{
						Resource: "expected",
						Path:     "expected",
					},
					Message: specs.Message{
						"example": {
							Name: "example",
							Path: "example",
						},
					},
				},
			}

			target := &specs.Property{
				Name: "example",
				Path: "example",
				Template: specs.Template{
					Reference: &specs.PropertyReference{
						Resource: "mock",
						Path:     "",
						Property: expected,
					},
				},
			}

			node := &specs.Node{
				Rollback: &specs.Call{
					Response: &specs.ParameterMap{
						Property: &specs.Property{
							Template: specs.Template{
								Message: specs.Message{
									target.Name: target,
								},
							},
						},
					},
				},
			}

			return node, target, expected
		},
		"header": func() (*specs.Node, *specs.Property, *specs.Property) {
			expected := &specs.Property{
				Template: specs.Template{
					Reference: &specs.PropertyReference{
						Resource: "expected",
						Path:     "expected",
					},
				},
			}

			target := &specs.Property{
				Template: specs.Template{
					Reference: &specs.PropertyReference{
						Resource: "mock",
						Path:     "",
						Property: expected,
					},
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
			mem := functions.Collection{}
			ResolveNodeReferences(node, mem)

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
				Template: specs.Template{
					Reference: &specs.PropertyReference{
						Resource: "expected",
						Path:     "expected",
					},
				},
			}

			target := &specs.Property{
				Template: specs.Template{
					Reference: &specs.PropertyReference{
						Resource: "mock",
						Path:     "",
						Property: expected,
					},
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
				Template: specs.Template{
					Reference: &specs.PropertyReference{
						Resource: "expected",
						Path:     "expected",
					},
					Message: specs.Message{
						"example": {
							Name: "example",
							Path: "example",
						},
					},
				},
			}

			target := &specs.Property{
				Template: specs.Template{
					Reference: &specs.PropertyReference{
						Resource: "mock",
						Path:     "",
						Property: expected,
					},
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
				Template: specs.Template{
					Reference: &specs.PropertyReference{
						Resource: "expected",
						Path:     "expected",
					},
				},
			}

			target := &specs.Property{
				Template: specs.Template{
					Reference: &specs.PropertyReference{
						Resource: "mock",
						Path:     "",
						Property: expected,
					},
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
			ctx := logger.WithLogger(broker.NewBackground())
			mem := functions.Collection{}
			ResolveReferences(ctx, specs.FlowListInterface{flow}, mem)

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
