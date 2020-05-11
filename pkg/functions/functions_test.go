package functions

import (
	"errors"
	"testing"

	"github.com/jexia/maestro/pkg/instance"
	"github.com/jexia/maestro/pkg/specs"
	"github.com/jexia/maestro/pkg/specs/labels"
	"github.com/jexia/maestro/pkg/specs/types"
)

type counter struct {
	total int
	err   error
}

func (counter *counter) fn(args ...*specs.Property) (*specs.Property, Exec, error) {
	counter.total++

	result := &specs.Property{
		Type:    types.String,
		Label:   labels.Optional,
		Default: "",
	}

	return result, nil, counter.err
}

func CompareProperties(t *testing.T, left specs.Property, right specs.Property) {
	if left.Default != right.Default {
		t.Errorf("unexpected default '%s', expected '%s'", left.Default, right.Default)
	}

	if left.Type != right.Type {
		t.Errorf("unexpected type '%s', expected '%s'", left.Type, right.Type)
	}

	if left.Label != right.Label {
		t.Errorf("unexpected label '%s', expected '%s'", left.Label, right.Label)
	}

	if right.Reference != nil && left.Reference == nil {
		t.Error("reference not set but expected")
	}

	if right.Reference != nil {
		if left.Reference.Resource != right.Reference.Resource {
			t.Errorf("unexpected reference resource '%s', expected '%s'", left.Reference.Resource, right.Reference.Resource)
		}

		if left.Reference.Path != right.Reference.Path {
			t.Errorf("unexpected reference path '%s', expected '%s'", left.Reference.Path, right.Reference.Path)
		}
	}
}

func TestParseFunction(t *testing.T) {
	static := specs.Property{
		Path:    "message",
		Default: "message",
		Type:    types.String,
		Label:   labels.Optional,
	}

	custom := Custom{
		"static": func(args ...*specs.Property) (*specs.Property, Exec, error) {
			return &static, nil, nil
		},
	}

	// NOTE: testing of sub functions is a function specific implementation and is not part of the template library
	tests := map[string]specs.Property{
		"static()": static,
	}

	for input, expected := range tests {
		t.Run(input, func(t *testing.T) {
			ctx := instance.NewContext()
			prop := specs.Property{
				Name: "",
				Path: "message",
				Raw:  input,
			}

			err := PrepareFunction(ctx, nil, nil, &prop, make(Stack), custom)
			if err != nil {
				t.Error(err)
			}

			if prop.Reference.Property == nil {
				t.Fatalf("unexpected reference property, reference property not set '%+v'", prop)
			}

			CompareProperties(t, prop, expected)
		})
	}
}

func TestParseUnavailableFunction(t *testing.T) {
	custom := Custom{}
	tests := []string{
		"add()",
	}

	for _, input := range tests {
		prop := &specs.Property{
			Name: "",
			Path: "message",
			Raw:  input,
		}

		ctx := instance.NewContext()
		err := PrepareFunction(ctx, nil, nil, prop, make(Stack), custom)
		if err == nil {
			t.Error("unexpected pass")
		}
	}
}

func TestPrepareManifestFunctions(t *testing.T) {
	type test struct {
		expected    int
		collections int
		manifest    *specs.FlowsManifest
	}

	tests := map[string]test{
		"flow": {
			expected:    3,
			collections: 3,
			manifest: &specs.FlowsManifest{
				Flows: specs.Flows{
					&specs.Flow{
						Nodes: []*specs.Node{
							{
								Call: &specs.Call{
									Request: &specs.ParameterMap{
										Property: &specs.Property{
											Type:  types.Message,
											Label: labels.Optional,
											Nested: map[string]*specs.Property{
												"fn": {
													Name: "fn",
													Path: "fn",
													Raw:  "mock()",
												},
											},
										},
									},
									Response: &specs.ParameterMap{
										Property: &specs.Property{
											Type:  types.Message,
											Label: labels.Optional,
											Nested: map[string]*specs.Property{
												"fn": {
													Name: "fn",
													Path: "fn",
													Raw:  "mock()",
												},
											},
										},
									},
								},
							},
						},
						Output: &specs.ParameterMap{
							Property: &specs.Property{
								Type:  types.Message,
								Label: labels.Optional,
								Nested: map[string]*specs.Property{
									"fn": {
										Name: "fn",
										Path: "fn",
										Raw:  "mock()",
									},
								},
							},
						},
					},
				},
			},
		},
		"flow rollback": {
			expected:    6,
			collections: 2,
			manifest: &specs.FlowsManifest{
				Flows: specs.Flows{
					&specs.Flow{
						Nodes: []*specs.Node{
							{
								Rollback: &specs.Call{
									Request: &specs.ParameterMap{
										Header: specs.Header{
											"fn": &specs.Property{
												Name: "fn",
												Path: "fn",
												Raw:  "mock(mock())",
											},
										},
										Property: &specs.Property{
											Type:  types.Message,
											Label: labels.Optional,
											Nested: map[string]*specs.Property{
												"fn": {
													Name: "fn",
													Path: "fn",
													Raw:  "mock()",
												},
											},
										},
									},
									Response: &specs.ParameterMap{
										Header: specs.Header{
											"fn": &specs.Property{
												Name: "fn",
												Path: "fn",
												Raw:  "mock(mock())",
											},
										},
										Property: &specs.Property{
											Type:  types.Message,
											Label: labels.Optional,
											Nested: map[string]*specs.Property{
												"fn": {
													Name: "fn",
													Path: "fn",
													Raw:  "mock()",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		"flow header": {
			expected:    4,
			collections: 2,
			manifest: &specs.FlowsManifest{
				Flows: specs.Flows{
					&specs.Flow{
						Nodes: []*specs.Node{
							{
								Call: &specs.Call{
									Request: &specs.ParameterMap{
										Header: specs.Header{
											"fn": &specs.Property{
												Name: "fn",
												Path: "fn",
												Raw:  "mock(mock())",
											},
										},
									},
								},
							},
						},
						Output: &specs.ParameterMap{
							Header: specs.Header{
								"fn": &specs.Property{
									Name: "fn",
									Path: "fn",
									Raw:  "mock(mock())",
								},
							},
						},
					},
				},
			},
		},
		"proxy": {
			expected:    2,
			collections: 2,
			manifest: &specs.FlowsManifest{
				Proxy: specs.Proxies{
					&specs.Proxy{
						Nodes: []*specs.Node{
							{
								Call: &specs.Call{
									Request: &specs.ParameterMap{
										Property: &specs.Property{
											Type:  types.Message,
											Label: labels.Optional,
											Nested: map[string]*specs.Property{
												"fn": {
													Name: "fn",
													Path: "fn",
													Raw:  "mock()",
												},
											},
										},
									},
									Response: &specs.ParameterMap{
										Property: &specs.Property{
											Type:  types.Message,
											Label: labels.Optional,
											Nested: map[string]*specs.Property{
												"fn": {
													Name: "fn",
													Path: "fn",
													Raw:  "mock()",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		"proxy rollback": {
			expected:    6,
			collections: 2,
			manifest: &specs.FlowsManifest{
				Proxy: specs.Proxies{
					&specs.Proxy{
						Nodes: []*specs.Node{
							{
								Rollback: &specs.Call{
									Request: &specs.ParameterMap{
										Header: specs.Header{
											"fn": &specs.Property{
												Name: "fn",
												Path: "fn",
												Raw:  "mock(mock())",
											},
										},
										Property: &specs.Property{
											Type:  types.Message,
											Label: labels.Optional,
											Nested: map[string]*specs.Property{
												"fn": {
													Name: "fn",
													Path: "fn",
													Raw:  "mock()",
												},
											},
										},
									},
									Response: &specs.ParameterMap{
										Header: specs.Header{
											"fn": &specs.Property{
												Name: "fn",
												Path: "fn",
												Raw:  "mock(mock())",
											},
										},
										Property: &specs.Property{
											Type:  types.Message,
											Label: labels.Optional,
											Nested: map[string]*specs.Property{
												"fn": {
													Name: "fn",
													Path: "fn",
													Raw:  "mock()",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		"proxy header": {
			expected:    4,
			collections: 2,
			manifest: &specs.FlowsManifest{
				Proxy: specs.Proxies{
					&specs.Proxy{
						Nodes: []*specs.Node{
							{
								Call: &specs.Call{
									Request: &specs.ParameterMap{
										Header: specs.Header{
											"fn": &specs.Property{
												Name: "fn",
												Path: "fn",
												Raw:  "mock(mock())",
											},
										},
									},
								},
							},
						},
						Forward: &specs.Call{
							Request: &specs.ParameterMap{
								Header: specs.Header{
									"fn": &specs.Property{
										Name: "fn",
										Path: "fn",
										Raw:  "mock(mock())",
									},
								},
							},
						},
					},
				},
			},
		},
		"nil request header property": {
			expected:    0,
			collections: 1,
			manifest: &specs.FlowsManifest{
				Proxy: specs.Proxies{
					&specs.Proxy{
						Nodes: []*specs.Node{
							{
								Call: &specs.Call{
									Request: &specs.ParameterMap{
										Header: specs.Header{
											"fn": nil,
										},
									},
								},
							},
						},
					},
				},
			},
		},
		"nil response header property": {
			expected:    0,
			collections: 1,
			manifest: &specs.FlowsManifest{
				Proxy: specs.Proxies{
					&specs.Proxy{
						Nodes: []*specs.Node{
							{
								Call: &specs.Call{
									Response: &specs.ParameterMap{
										Header: specs.Header{
											"fn": nil,
										},
									},
								},
							},
						},
					},
				},
			},
		},
		"nil request property": {
			expected:    0,
			collections: 1,
			manifest: &specs.FlowsManifest{
				Proxy: specs.Proxies{
					&specs.Proxy{
						Nodes: []*specs.Node{
							{
								Call: &specs.Call{
									Request: &specs.ParameterMap{
										Property: nil,
									},
								},
							},
						},
					},
				},
			},
		},
		"nil response property": {
			expected:    0,
			collections: 1,
			manifest: &specs.FlowsManifest{
				Proxy: specs.Proxies{
					&specs.Proxy{
						Nodes: []*specs.Node{
							{
								Call: &specs.Call{
									Response: &specs.ParameterMap{
										Property: nil,
									},
								},
							},
						},
					},
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			counter := counter{}
			functions := Custom{
				"mock": counter.fn,
			}

			ctx := instance.NewContext()
			mem := Collection{}

			err := PrepareManifestFunctions(ctx, mem, functions, test.manifest)
			if err != nil {
				t.Fatal(err)
			}

			if counter.total != test.expected {
				t.Fatalf("unexpected counter %d, expected %d functions to be called", counter.total, test.expected)
			}

			if len(mem) != test.collections {
				t.Fatalf("unexpected collections length %d, expected %d collections to be defined", len(mem), test.collections)
			}
		})
	}
}

func TestPrepareManifestFunctionsErr(t *testing.T) {
	type test struct {
		expected    int
		collections int
		manifest    *specs.FlowsManifest
	}

	tests := map[string]test{
		"flow": {
			expected:    3,
			collections: 3,
			manifest: &specs.FlowsManifest{
				Flows: specs.Flows{
					&specs.Flow{
						Nodes: []*specs.Node{
							{
								Call: &specs.Call{
									Request: &specs.ParameterMap{
										Property: &specs.Property{
											Type:  types.Message,
											Label: labels.Optional,
											Nested: map[string]*specs.Property{
												"fn": {
													Name: "fn",
													Path: "fn",
													Raw:  "mock()",
												},
											},
										},
									},
									Response: &specs.ParameterMap{
										Property: &specs.Property{
											Type:  types.Message,
											Label: labels.Optional,
											Nested: map[string]*specs.Property{
												"fn": {
													Name: "fn",
													Path: "fn",
													Raw:  "mock()",
												},
											},
										},
									},
								},
							},
						},
						Output: &specs.ParameterMap{
							Property: &specs.Property{
								Type:  types.Message,
								Label: labels.Optional,
								Nested: map[string]*specs.Property{
									"fn": {
										Name: "fn",
										Path: "fn",
										Raw:  "mock()",
									},
								},
							},
						},
					},
				},
			},
		},
		"flow rollback": {
			expected:    6,
			collections: 2,
			manifest: &specs.FlowsManifest{
				Flows: specs.Flows{
					&specs.Flow{
						Nodes: []*specs.Node{
							{
								Rollback: &specs.Call{
									Request: &specs.ParameterMap{
										Header: specs.Header{
											"fn": &specs.Property{
												Name: "fn",
												Path: "fn",
												Raw:  "mock(mock())",
											},
										},
										Property: &specs.Property{
											Type:  types.Message,
											Label: labels.Optional,
											Nested: map[string]*specs.Property{
												"fn": {
													Name: "fn",
													Path: "fn",
													Raw:  "mock()",
												},
											},
										},
									},
									Response: &specs.ParameterMap{
										Header: specs.Header{
											"fn": &specs.Property{
												Name: "fn",
												Path: "fn",
												Raw:  "mock(mock())",
											},
										},
										Property: &specs.Property{
											Type:  types.Message,
											Label: labels.Optional,
											Nested: map[string]*specs.Property{
												"fn": {
													Name: "fn",
													Path: "fn",
													Raw:  "mock()",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		"flow header": {
			expected:    4,
			collections: 2,
			manifest: &specs.FlowsManifest{
				Flows: specs.Flows{
					&specs.Flow{
						Nodes: []*specs.Node{
							{
								Call: &specs.Call{
									Request: &specs.ParameterMap{
										Header: specs.Header{
											"fn": &specs.Property{
												Name: "fn",
												Path: "fn",
												Raw:  "mock(mock())",
											},
										},
									},
								},
							},
						},
						Output: &specs.ParameterMap{
							Header: specs.Header{
								"fn": &specs.Property{
									Name: "fn",
									Path: "fn",
									Raw:  "mock(mock())",
								},
							},
						},
					},
				},
			},
		},
		"proxy": {
			expected:    2,
			collections: 2,
			manifest: &specs.FlowsManifest{
				Proxy: specs.Proxies{
					&specs.Proxy{
						Nodes: []*specs.Node{
							{
								Call: &specs.Call{
									Request: &specs.ParameterMap{
										Property: &specs.Property{
											Type:  types.Message,
											Label: labels.Optional,
											Nested: map[string]*specs.Property{
												"fn": {
													Name: "fn",
													Path: "fn",
													Raw:  "mock()",
												},
											},
										},
									},
									Response: &specs.ParameterMap{
										Property: &specs.Property{
											Type:  types.Message,
											Label: labels.Optional,
											Nested: map[string]*specs.Property{
												"fn": {
													Name: "fn",
													Path: "fn",
													Raw:  "mock()",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		"proxy rollback": {
			expected:    6,
			collections: 2,
			manifest: &specs.FlowsManifest{
				Proxy: specs.Proxies{
					&specs.Proxy{
						Nodes: []*specs.Node{
							{
								Rollback: &specs.Call{
									Request: &specs.ParameterMap{
										Header: specs.Header{
											"fn": &specs.Property{
												Name: "fn",
												Path: "fn",
												Raw:  "mock(mock())",
											},
										},
										Property: &specs.Property{
											Type:  types.Message,
											Label: labels.Optional,
											Nested: map[string]*specs.Property{
												"fn": {
													Name: "fn",
													Path: "fn",
													Raw:  "mock()",
												},
											},
										},
									},
									Response: &specs.ParameterMap{
										Header: specs.Header{
											"fn": &specs.Property{
												Name: "fn",
												Path: "fn",
												Raw:  "mock(mock())",
											},
										},
										Property: &specs.Property{
											Type:  types.Message,
											Label: labels.Optional,
											Nested: map[string]*specs.Property{
												"fn": {
													Name: "fn",
													Path: "fn",
													Raw:  "mock()",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		"proxy header": {
			expected:    4,
			collections: 2,
			manifest: &specs.FlowsManifest{
				Proxy: specs.Proxies{
					&specs.Proxy{
						Nodes: []*specs.Node{
							{
								Call: &specs.Call{
									Request: &specs.ParameterMap{
										Header: specs.Header{
											"fn": &specs.Property{
												Name: "fn",
												Path: "fn",
												Raw:  "mock(mock())",
											},
										},
									},
								},
							},
						},
						Forward: &specs.Call{
							Request: &specs.ParameterMap{
								Header: specs.Header{
									"fn": &specs.Property{
										Name: "fn",
										Path: "fn",
										Raw:  "mock(mock())",
									},
								},
							},
						},
					},
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			counter := counter{
				err: errors.New("unexpected err"),
			}

			functions := Custom{
				"mock": counter.fn,
			}

			ctx := instance.NewContext()
			mem := Collection{}

			err := PrepareManifestFunctions(ctx, mem, functions, test.manifest)
			if err == nil {
				t.Fatal("unexpected pass expected prepare to fail")
			}
		})
	}
}

func TestPrepareParameterMapFunctions(t *testing.T) {
	type test struct {
		expected int
		params   *specs.ParameterMap
	}

	tests := map[string]test{
		"simple": {
			expected: 1,
			params: &specs.ParameterMap{
				Property: &specs.Property{
					Type:  types.Message,
					Label: labels.Optional,
					Nested: map[string]*specs.Property{
						"fn": {
							Name: "fn",
							Path: "fn",
							Raw:  "mock()",
						},
					},
				},
			},
		},
		"header": {
			expected: 1,
			params: &specs.ParameterMap{
				Header: specs.Header{
					"sample": &specs.Property{
						Name: "fn",
						Path: "fn",
						Raw:  "mock()",
					},
				},
			},
		},
		"nested": {
			expected: 1,
			params: &specs.ParameterMap{
				Property: &specs.Property{
					Type:  types.Message,
					Label: labels.Optional,
					Nested: map[string]*specs.Property{
						"nested": {
							Type:  types.Message,
							Label: labels.Optional,
							Nested: map[string]*specs.Property{
								"fn": {
									Name: "fn",
									Path: "fn",
									Raw:  "mock()",
								},
							},
						},
					},
				},
			},
		},
		"multiple": {
			expected: 3,
			params: &specs.ParameterMap{
				Property: &specs.Property{
					Type:  types.Message,
					Label: labels.Optional,
					Nested: map[string]*specs.Property{
						"fn": {
							Name: "fn",
							Path: "fn",
							Raw:  "mock(mock(mock()))",
						},
					},
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			counter := counter{}
			functions := Custom{
				"mock": counter.fn,
			}

			ctx := instance.NewContext()
			stack := Stack{}

			err := PrepareParameterMapFunctions(ctx, nil, nil, stack, test.params, functions)
			if err != nil {
				t.Fatal(err)
			}

			if counter.total != test.expected {
				t.Fatalf("unexpected counter result %d, expected %d functions to be called", counter.total, test.expected)
			}
		})
	}
}
