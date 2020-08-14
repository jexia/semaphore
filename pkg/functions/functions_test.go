package functions

import (
	"errors"
	"reflect"
	"testing"

	"github.com/jexia/semaphore/pkg/broker"
	"github.com/jexia/semaphore/pkg/broker/logger"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/labels"
	"github.com/jexia/semaphore/pkg/specs/types"
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

func TestCollectionReserve(t *testing.T) {
	key := &specs.ParameterMap{}
	expected := Stack{}

	collection := Collection{
		key: expected,
	}

	result := collection.Reserve(key)
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("unexpected result %+v, expected %+v", result, expected)
	}

	unknown := collection.Reserve(&specs.ParameterMap{})
	if len(unknown) != 0 {
		t.Fatalf("unknown stack %+v, expected a empty stack", unknown)
	}
}

func TestCollectionLoad(t *testing.T) {
	key := &specs.ParameterMap{}
	expected := Stack{}

	collection := Collection{
		key: expected,
	}

	result := collection.Load(key)
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("unexpected result %+v, expected %+v", result, expected)
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
			ctx := logger.WithLogger(broker.NewBackground())
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

		ctx := logger.WithLogger(broker.NewBackground())
		err := PrepareFunction(ctx, nil, nil, prop, make(Stack), custom)
		if err == nil {
			t.Error("unexpected pass")
		}
	}
}

func TestPrepareFunctions(t *testing.T) {
	type test struct {
		expected    int
		collections int
		flows       specs.FlowListInterface
	}

	tests := map[string]test{
		"flow": {
			expected:    3,
			collections: 3,
			flows: specs.FlowListInterface{
				&specs.Flow{
					Nodes: specs.NodeList{
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
		"flow rollback": {
			expected:    6,
			collections: 2,
			flows: specs.FlowListInterface{
				&specs.Flow{
					Nodes: specs.NodeList{
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
		"flow header": {
			expected:    4,
			collections: 2,
			flows: specs.FlowListInterface{
				&specs.Flow{
					Nodes: specs.NodeList{
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
		"proxy": {
			expected:    2,
			collections: 2,
			flows: specs.FlowListInterface{
				&specs.Proxy{
					Nodes: specs.NodeList{
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
		"proxy rollback": {
			expected:    6,
			collections: 2,
			flows: specs.FlowListInterface{
				&specs.Proxy{
					Nodes: specs.NodeList{
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
		"proxy header": {
			expected:    4,
			collections: 2,
			flows: specs.FlowListInterface{
				&specs.Proxy{
					Nodes: specs.NodeList{
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
		"nil request header property": {
			expected:    0,
			collections: 1,
			flows: specs.FlowListInterface{
				&specs.Proxy{
					Nodes: specs.NodeList{
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
		"nil response header property": {
			expected:    0,
			collections: 1,
			flows: specs.FlowListInterface{
				&specs.Proxy{
					Nodes: specs.NodeList{
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
		"nil request property": {
			expected:    0,
			collections: 1,
			flows: specs.FlowListInterface{
				&specs.Proxy{
					Nodes: specs.NodeList{
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
		"nil response property": {
			expected:    0,
			collections: 1,
			flows: specs.FlowListInterface{
				&specs.Proxy{
					Nodes: specs.NodeList{
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
		"condition": {
			expected:    1,
			collections: 1,
			flows: specs.FlowListInterface{
				&specs.Proxy{
					Nodes: specs.NodeList{
						{
							Condition: &specs.Condition{
								Params: &specs.ParameterMap{
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
		"intermediate": {
			expected:    1,
			collections: 1,
			flows: specs.FlowListInterface{
				&specs.Proxy{
					Nodes: specs.NodeList{
						{
							Intermediate: &specs.ParameterMap{
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
		"params": {
			expected:    2,
			collections: 1,
			flows: specs.FlowListInterface{
				&specs.Flow{
					Nodes: specs.NodeList{
						{
							Call: &specs.Call{
								Response: &specs.ParameterMap{
									Params: map[string]*specs.Property{
										"fn": {
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
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			counter := counter{}
			functions := Custom{
				"mock": counter.fn,
			}

			ctx := logger.WithLogger(broker.NewBackground())
			mem := Collection{}

			err := PrepareManifestFunctions(ctx, mem, functions, test.flows)
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

func TestPrepareFunctionsErr(t *testing.T) {
	type test struct {
		expected    int
		collections int
		flows       specs.FlowListInterface
	}

	tests := map[string]test{
		"flow": {
			expected:    3,
			collections: 3,
			flows: specs.FlowListInterface{
				&specs.Flow{
					Nodes: specs.NodeList{
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
		"flow rollback": {
			expected:    6,
			collections: 2,
			flows: specs.FlowListInterface{
				&specs.Flow{
					Nodes: specs.NodeList{
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
		"flow header": {
			expected:    4,
			collections: 2,
			flows: specs.FlowListInterface{
				&specs.Flow{
					Nodes: specs.NodeList{
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
		"proxy": {
			expected:    2,
			collections: 2,
			flows: specs.FlowListInterface{
				&specs.Proxy{
					Nodes: specs.NodeList{
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
		"proxy rollback": {
			expected:    6,
			collections: 2,
			flows: specs.FlowListInterface{
				&specs.Proxy{
					Nodes: specs.NodeList{
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
		"proxy header": {
			expected:    4,
			collections: 2,
			flows: specs.FlowListInterface{
				&specs.Proxy{
					Nodes: specs.NodeList{
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
		"condition": {
			expected:    1,
			collections: 1,
			flows: specs.FlowListInterface{
				&specs.Flow{
					Nodes: specs.NodeList{
						{
							Condition: &specs.Condition{
								Params: &specs.ParameterMap{
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
		"intermediate": {
			expected:    1,
			collections: 1,
			flows: specs.FlowListInterface{
				&specs.Flow{
					Nodes: specs.NodeList{
						{
							Intermediate: &specs.ParameterMap{
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
		"call response": {
			expected:    1,
			collections: 1,
			flows: specs.FlowListInterface{
				&specs.Flow{
					Nodes: specs.NodeList{
						{
							Call: &specs.Call{
								Response: &specs.ParameterMap{
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
		},
		"params": {
			expected:    2,
			collections: 1,
			flows: specs.FlowListInterface{
				&specs.Flow{
					Nodes: specs.NodeList{
						{
							Call: &specs.Call{
								Response: &specs.ParameterMap{
									Params: map[string]*specs.Property{
										"fn": {
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
		},
		"output": {
			expected:    1,
			collections: 1,
			flows: specs.FlowListInterface{
				&specs.Flow{
					Output: &specs.ParameterMap{
						Params: map[string]*specs.Property{
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
		"forward": {
			expected:    1,
			collections: 1,
			flows: specs.FlowListInterface{
				&specs.Proxy{
					Forward: &specs.Call{
						Request: &specs.ParameterMap{
							Header: map[string]*specs.Property{
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
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			counter := counter{
				err: errors.New("unexpected err"),
			}

			functions := Custom{
				"mock": counter.fn,
			}

			ctx := logger.WithLogger(broker.NewBackground())
			mem := Collection{}

			err := PrepareManifestFunctions(ctx, mem, functions, test.flows)
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

			ctx := logger.WithLogger(broker.NewBackground())
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

func TestPrepareParamsFunctionsNil(t *testing.T) {
	err := PrepareParamsFunctions(nil, nil, nil, nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
}

func TestPreparePropertyFunctionsNil(t *testing.T) {
	err := PreparePropertyFunctions(nil, nil, nil, nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
}

func TestPrepareFunctionNil(t *testing.T) {
	err := PrepareFunction(nil, nil, nil, nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
}
