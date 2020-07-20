package dependencies

import (
	"testing"

	"github.com/jexia/semaphore/pkg/core/instance"
	"github.com/jexia/semaphore/pkg/specs"
)

func TestResolveReferences(t *testing.T) {
	tests := map[string]func() (*specs.FlowsManifest, *specs.Property, *specs.Property){
		"flow": func() (*specs.FlowsManifest, *specs.Property, *specs.Property) {
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

			manifest := &specs.FlowsManifest{
				Flows: specs.Flows{
					{
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
				},
			}

			return manifest, target, expected
		},
		"proxy": func() (*specs.FlowsManifest, *specs.Property, *specs.Property) {
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

			manifest := &specs.FlowsManifest{
				Proxy: specs.Proxies{
					{
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
				},
			}

			return manifest, target, expected
		},
		"params": func() (*specs.FlowsManifest, *specs.Property, *specs.Property) {
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

			manifest := &specs.FlowsManifest{
				Flows: specs.Flows{
					{
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
				},
			}

			return manifest, target, expected
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := instance.NewContext()
			manifest, target, expected := test()
			ResolveReferences(ctx, manifest)

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
			ctx := instance.NewContext()
			ResolveReferences(ctx, &specs.FlowsManifest{
				Flows: []*specs.Flow{
					flow,
				},
			})

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
