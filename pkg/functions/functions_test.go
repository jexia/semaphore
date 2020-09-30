package functions

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/jexia/semaphore/pkg/broker"
	"github.com/jexia/semaphore/pkg/broker/logger"
	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/labels"
	"github.com/jexia/semaphore/pkg/specs/template"
	"github.com/jexia/semaphore/pkg/specs/types"
)

type counter struct {
	total int
	err   error
}

func (counter *counter) fn(args ...*specs.Property) (*specs.Property, Exec, error) {
	counter.total++

	result := &specs.Property{
		Template: specs.Template{
			Scalar: &specs.Scalar{
				Default: "",
				Type:    types.String,
			},
		},
		Label: labels.Optional,
	}

	return result, nil, counter.err
}

func CompareScalarProperties(t *testing.T, left specs.Property, right specs.Property) {
	if left.Scalar.Default != right.Scalar.Default {
		t.Errorf("unexpected default '%s', expected '%s'", left.Scalar.Default, right.Scalar.Default)
	}

	if left.Type() != right.Type() {
		t.Errorf("unexpected type '%s', expected '%s'", left.Type(), right.Type())
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

func TestUndefinedFunction(t *testing.T) {
	type fields struct {
		Function string
		Property string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"return the formatted error",
			fields{Function: "getsources", Property: "add"},
			"undefined custom function 'getsources' in 'add'",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := ErrUndefinedFunction{
				Function: "getsources",
				Property: "add",
			}
			if got := e.Prettify(); got.Message != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseFunction(t *testing.T) {
	static := specs.Property{
		Path: "message",
		Template: specs.Template{
			Scalar: &specs.Scalar{
				Default: "message",
				Type:    types.String,
			},
		},
		Label: labels.Optional,
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

			CompareScalarProperties(t, prop, expected)
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
										Label: labels.Optional,
										Template: specs.Template{
											Message: specs.Message{
												"fn": &specs.Property{
													Name: "fn",
													Path: "fn",
													Raw:  "mock()",
												},
											},
										},
									},
								},
								Response: &specs.ParameterMap{
									Property: &specs.Property{
										Label: labels.Optional,
										Template: specs.Template{
											Message: specs.Message{
												"fn": &specs.Property{
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
					Output: &specs.ParameterMap{
						Property: &specs.Property{
							Label: labels.Optional,
							Template: specs.Template{
								Message: specs.Message{
									"fn": &specs.Property{
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
										Label: labels.Optional,
										Template: specs.Template{
											Message: specs.Message{
												"fn": &specs.Property{
													Name: "fn",
													Path: "fn",
													Raw:  "mock()",
												},
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
										Label: labels.Optional,
										Template: specs.Template{
											Message: specs.Message{
												"fn": &specs.Property{
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
										Label: labels.Optional,
										Template: specs.Template{
											Message: specs.Message{
												"fn": &specs.Property{
													Name: "fn",
													Path: "fn",
													Raw:  "mock()",
												},
											},
										},
									},
								},
								Response: &specs.ParameterMap{
									Property: &specs.Property{
										Label: labels.Optional,
										Template: specs.Template{
											Message: specs.Message{
												"fn": &specs.Property{
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
										Label: labels.Optional,
										Template: specs.Template{
											Message: specs.Message{
												"fn": &specs.Property{
													Name: "fn",
													Path: "fn",
													Raw:  "mock()",
												},
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
										Label: labels.Optional,
										Template: specs.Template{
											Message: specs.Message{
												"fn": &specs.Property{
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
										Label: labels.Optional,
										Template: specs.Template{
											Message: specs.Message{
												"fn": &specs.Property{
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
		"intermediate": {
			expected:    1,
			collections: 1,
			flows: specs.FlowListInterface{
				&specs.Proxy{
					Nodes: specs.NodeList{
						{
							Intermediate: &specs.ParameterMap{
								Property: &specs.Property{
									Label: labels.Optional,
									Template: specs.Template{
										Message: specs.Message{
											"fn": &specs.Property{
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

			err := PrepareFunctions(ctx, mem, functions, test.flows)
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
		flows specs.FlowListInterface
	}

	tests := map[string]test{
		"flow": {
			flows: specs.FlowListInterface{
				&specs.Flow{
					Nodes: specs.NodeList{
						{
							Call: &specs.Call{
								Request: &specs.ParameterMap{
									Property: &specs.Property{
										Label: labels.Optional,
										Template: specs.Template{
											Message: specs.Message{
												"fn": &specs.Property{
													Name: "fn",
													Path: "fn",
													Raw:  "mock()",
												},
											},
										},
									},
								},
								Response: &specs.ParameterMap{
									Property: &specs.Property{
										Label: labels.Optional,
										Template: specs.Template{
											Message: specs.Message{
												"fn": &specs.Property{
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
					Output: &specs.ParameterMap{
						Property: &specs.Property{
							Label: labels.Optional,
							Template: specs.Template{
								Message: specs.Message{
									"fn": &specs.Property{
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
										Label: labels.Optional,
										Template: specs.Template{
											Message: specs.Message{
												"fn": &specs.Property{
													Name: "fn",
													Path: "fn",
													Raw:  "mock()",
												},
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
										Label: labels.Optional,
										Template: specs.Template{
											Message: specs.Message{
												"fn": &specs.Property{
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
			flows: specs.FlowListInterface{
				&specs.Proxy{
					Nodes: specs.NodeList{
						{
							Call: &specs.Call{
								Request: &specs.ParameterMap{
									Property: &specs.Property{
										Label: labels.Optional,
										Template: specs.Template{
											Message: specs.Message{
												"fn": &specs.Property{
													Name: "fn",
													Path: "fn",
													Raw:  "mock()",
												},
											},
										},
									},
								},
								Response: &specs.ParameterMap{
									Property: &specs.Property{
										Label: labels.Optional,
										Template: specs.Template{
											Message: specs.Message{
												"fn": &specs.Property{
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
										Label: labels.Optional,
										Template: specs.Template{
											Message: specs.Message{
												"fn": &specs.Property{
													Name: "fn",
													Path: "fn",
													Raw:  "mock()",
												},
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
										Label: labels.Optional,
										Template: specs.Template{
											Message: specs.Message{
												"fn": &specs.Property{
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
			flows: specs.FlowListInterface{
				&specs.Flow{
					Nodes: specs.NodeList{
						{
							Condition: &specs.Condition{
								Params: &specs.ParameterMap{
									Property: &specs.Property{
										Label: labels.Optional,
										Template: specs.Template{
											Message: specs.Message{
												"fn": &specs.Property{
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
		"intermediate": {
			flows: specs.FlowListInterface{
				&specs.Flow{
					Nodes: specs.NodeList{
						{
							Intermediate: &specs.ParameterMap{
								Property: &specs.Property{
									Label: labels.Optional,
									Template: specs.Template{
										Message: specs.Message{
											"fn": &specs.Property{
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
		"call response": {
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

			err := PrepareFunctions(ctx, mem, functions, test.flows)
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
					Label: labels.Optional,
					Template: specs.Template{
						Message: specs.Message{
							"fn": &specs.Property{
								Name: "fn",
								Path: "fn",
								Raw:  "mock()",
							},
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
					Label: labels.Optional,
					Template: specs.Template{
						Message: specs.Message{
							"fn": &specs.Property{
								Name: "nested",
								Path: "nested",

								Template: specs.Template{
									Message: specs.Message{
										"fn": &specs.Property{
											Name: "fn",
											Path: "nested.fn",
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
		"multiple": {
			expected: 3,
			params: &specs.ParameterMap{
				Property: &specs.Property{
					Label: labels.Optional,
					Template: specs.Template{
						Message: specs.Message{
							"fn": &specs.Property{
								Name: "fn",
								Path: "nested.fn",
								Raw:  "mock(mock(mock()))",
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

func TestFunctionsNestedReferences(t *testing.T) {
	result := &specs.Property{
		Label: labels.Optional,
		Template: specs.Template{
			Message: specs.Message{
				"id": &specs.Property{
					Name:  "id",
					Path:  "id",
					Label: labels.Optional,

					Template: specs.Template{
						Scalar: &specs.Scalar{
							Default: "abc",
							Type:    types.String,
						},
					},
				},
			},
		},
	}

	functions := Custom{
		"mock": func(args ...*specs.Property) (*specs.Property, Exec, error) {
			return result, func(references.Store) error { return nil }, nil
		},
	}

	property := &specs.Property{
		Raw: "mock()",
	}

	ctx := logger.WithLogger(broker.NewBackground())
	stack := Stack{}

	err := PrepareFunction(ctx, nil, nil, property, stack, functions)
	if err != nil {
		t.Fatal(err)
	}

	if property.Reference == nil || property.Reference.Property == nil {
		t.Fatal("property reference not set")
	}

	if len(property.Message) != len(result.Message) {
		t.Fatal("property reference nested is not equal to result")
	}

	for _, nested := range property.Message {
		if nested.Reference == nil {
			t.Fatal("nested reference not set")
		}

		if !strings.HasPrefix(nested.Reference.Resource, template.StackResource) {
			t.Errorf("nested does not reference stack, %+v", nested.Reference.Resource)
		}

		if nested.Reference.Property == nil {
			t.Error("nested property is not set")
		}
	}
}
