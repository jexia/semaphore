package lookup

import (
	"testing"

	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/labels"
	"github.com/jexia/semaphore/pkg/specs/template"
	"github.com/jexia/semaphore/pkg/specs/types"
)

func TestGetFlow(t *testing.T) {
	flows := specs.FlowListInterface{
		&specs.Flow{
			Name: "first",
		},
		&specs.Flow{
			Name: "second",
		},
	}

	tests := map[string]specs.FlowInterface{
		"second":  flows[1],
		"first":   flows[0],
		"unknown": nil,
	}

	for input, expected := range tests {
		result := flows.Get(input)
		if result != expected {
			t.Errorf("unexpected result %+v, expected %+v", result, expected)
		}
	}
}

func TestGetDefaultProp(t *testing.T) {
	tests := map[string]string{
		template.InputResource: template.RequestResource,
		"unknown":              template.ResponseResource,
	}

	for input, expected := range tests {
		result := GetDefaultProp(input)
		if result != expected {
			t.Errorf("unexpected result %s, expected %s", result, expected)
		}
	}
}

func TestGetNextResource(t *testing.T) {
	type test struct {
		breakpoint string
		expected   string
		manager    *specs.Flow
	}

	tests := map[string]*test{
		"first": {
			breakpoint: "first",
			expected:   "second",
			manager: &specs.Flow{
				Nodes: specs.NodeList{
					{
						ID: "first",
					},
					{
						ID: "second",
					},
				},
			},
		},
		"second": {
			breakpoint: "second",
			expected:   "third",
			manager: &specs.Flow{
				Nodes: specs.NodeList{
					{
						ID: "first",
					},
					{
						ID: "second",
					},
					{
						ID: "third",
					},
				},
			},
		},
		"output": {
			breakpoint: "last",
			expected:   template.OutputResource,
			manager: &specs.Flow{
				Nodes: specs.NodeList{
					{
						ID: "first",
					},
					{
						ID: "last",
					},
				},
			},
		},
		"unknown": {
			breakpoint: "unknown",
			expected:   "unknown",
			manager: &specs.Flow{
				Nodes: specs.NodeList{
					{
						ID: "first",
					},
					{
						ID: "second",
					},
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := GetNextResource(test.manager, test.breakpoint)
			if result != test.expected {
				t.Fatalf("unexpected result '%s', expected '%s'", result, test.expected)
			}
		})
	}
}

func NewInputMockProperty() *specs.Property {
	return &specs.Property{
		Path:  "",
		Label: labels.Optional,
		Template: specs.Template{
			Message: specs.Message{
				"message": {
					Position: 0,
					Name:     "message",
					Path:     "message",
					Label:    labels.Optional,
					Template: specs.Template{
						Scalar: &specs.Scalar{
							Default: "hello world",
							Type:    types.String,
						},
					},
				},
				"active": {
					Position: 1,
					Name:     "active",
					Path:     "active",
					Label:    labels.Optional,
					Template: specs.Template{
						Scalar: &specs.Scalar{
							Default: false,
							Type:    types.Bool,
						},
					},
				},
				"nested": {
					Position: 2,
					Name:     "nested",
					Path:     "nested",
					Label:    labels.Optional,
					Template: specs.Template{
						Message: specs.Message{
							"message": {
								Position: 0,
								Name:     "message",
								Path:     "nested.message",
								Label:    labels.Optional,
								Template: specs.Template{
									Scalar: &specs.Scalar{
										Default: "hello nested world",
										Type:    types.String,
									},
								},
							},
							"active": {
								Position: 1,
								Name:     "active",
								Path:     "nested.active",
								Label:    labels.Optional,
								Template: specs.Template{
									Scalar: &specs.Scalar{
										Default: false,
										Type:    types.Bool,
									},
								},
							},
							"nested": {
								Position: 2,
								Name:     "nested",
								Path:     "nested.nested",
								Label:    labels.Optional,
								Template: specs.Template{
									Message: specs.Message{
										"message": {
											Name:  "message",
											Path:  "nested.nested.message",
											Label: labels.Optional,
											Template: specs.Template{
												Scalar: &specs.Scalar{
													Default: "hello underground world",
													Type:    types.String,
												},
											},
										},
										"active": {
											Name:  "active",
											Path:  "nested.nested.active",
											Label: labels.Optional,
											Template: specs.Template{
												Scalar: &specs.Scalar{
													Default: false,
													Type:    types.Bool,
												},
											},
										},
									},
								},
							},
							"repeated": {
								Name:  "repeated",
								Path:  "nested.repeated",
								Label: labels.Optional,
								Template: specs.Template{
									Repeated: specs.Repeated{
										{
											Message: specs.Message{
												"message": {
													Name:  "message",
													Path:  "nested.repeated.message",
													Label: labels.Optional,
													Template: specs.Template{
														Scalar: &specs.Scalar{
															Default: "hello repeated underground world",
															Type:    types.String,
														},
													},
												},
												"active": {
													Name:  "active",
													Path:  "nested.repeated.active",
													Label: labels.Optional,
													Template: specs.Template{
														Scalar: &specs.Scalar{
															Default: false,
															Type:    types.Bool,
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
				"repeated": {
					Name:  "repeated",
					Path:  "repeated",
					Label: labels.Optional,
					Template: specs.Template{
						Repeated: specs.Repeated{
							{
								Message: specs.Message{
									"message": {
										Name:  "message",
										Path:  "message",
										Label: labels.Optional,
										Template: specs.Template{
											Scalar: &specs.Scalar{
												Default: "hello repeated world",
												Type:    types.String,
											},
										},
									},
									"active": {
										Name:  "active",
										Path:  "active",
										Label: labels.Optional,
										Template: specs.Template{
											Scalar: &specs.Scalar{
												Default: false,
												Type:    types.Bool,
											},
										},
									},
									"repeated": {
										Name:  "repeated",
										Path:  "repeated.repeated",
										Label: labels.Optional,
										Template: specs.Template{
											Repeated: specs.Repeated{
												{
													Message: specs.Message{
														"message": {
															Name:  "message",
															Path:  "repeated.message",
															Label: labels.Optional,
															Template: specs.Template{
																Scalar: &specs.Scalar{
																	Default: "hello repeated nested world",
																	Type:    types.String,
																},
															},
														},
														"active": {
															Name:  "active",
															Path:  "repeated.active",
															Label: labels.Optional,
															Template: specs.Template{
																Scalar: &specs.Scalar{
																	Default: false,
																	Type:    types.Bool,
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
					},
				},
			},
		},
	}
}

func NewResultMockProperty() *specs.Property {
	return &specs.Property{
		Path:  "",
		Label: labels.Optional,
		Template: specs.Template{
			Message: specs.Message{
				"result": {
					Name:  "result",
					Path:  "result",
					Label: labels.Optional,
					Template: specs.Template{
						Scalar: &specs.Scalar{
							Default: "hello world",
							Type:    types.String,
						},
					},
				},
				"active": {
					Name:  "active",
					Path:  "active",
					Label: labels.Optional,
					Template: specs.Template{
						Scalar: &specs.Scalar{
							Default: false,
							Type:    types.Bool,
						},
					},
				},
				"nested": {
					Name:  "nested",
					Path:  "nested",
					Label: labels.Optional,
					Template: specs.Template{
						Message: specs.Message{
							"result": {
								Name:  "result",
								Path:  "nested.result",
								Label: labels.Optional,
								Template: specs.Template{
									Scalar: &specs.Scalar{
										Default: "hello world",
										Type:    types.String,
									},
								},
							},
							"active": {
								Name:  "active",
								Path:  "nested.active",
								Label: labels.Optional,
								Template: specs.Template{
									Scalar: &specs.Scalar{
										Default: false,
										Type:    types.Bool,
									},
								},
							},
							"nested": {
								Name:  "nested",
								Path:  "nested.nested",
								Label: labels.Optional,
								Template: specs.Template{
									Message: specs.Message{
										"result": {
											Name:  "result",
											Path:  "nested.nested.result",
											Label: labels.Optional,
											Template: specs.Template{
												Scalar: &specs.Scalar{
													Default: "hello world",
													Type:    types.String,
												},
											},
										},
										"active": {
											Name:  "active",
											Path:  "nested.nested.active",
											Label: labels.Optional,
											Template: specs.Template{
												Scalar: &specs.Scalar{
													Default: false,
													Type:    types.Bool,
												},
											},
										},
									},
								},
							},
							"repeated": {
								Name: "repeated",
								Path: "nested.repeated",
								Template: specs.Template{
									Repeated: specs.Repeated{
										{
											Message: specs.Message{
												"result": {
													Name:  "result",
													Path:  "nested.repeated.result",
													Label: labels.Optional,
													Template: specs.Template{
														Scalar: &specs.Scalar{
															Default: "hello world",
															Type:    types.String,
														},
													},
												},
												"active": {
													Name:  "active",
													Path:  "nested.repeated.active",
													Label: labels.Optional,
													Template: specs.Template{
														Scalar: &specs.Scalar{
															Default: false,
															Type:    types.Bool,
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
				"repeated": {
					Name: "repeated",
					Path: "repeated",
					Template: specs.Template{
						Repeated: specs.Repeated{
							{
								Message: specs.Message{
									"result": {
										Name:  "result",
										Path:  "repeated.result",
										Label: labels.Optional,
										Template: specs.Template{
											Scalar: &specs.Scalar{
												Default: "hello repeated world",
												Type:    types.String,
											},
										},
									},
									"active": {
										Name:  "active",
										Path:  "repeated.active",
										Label: labels.Optional,
										Template: specs.Template{
											Scalar: &specs.Scalar{
												Default: false,
												Type:    types.Bool,
											},
										},
									},
									"repeated": {
										Name:  "repeated",
										Path:  "repeated.repeated",
										Label: labels.Optional,
										Template: specs.Template{
											Repeated: specs.Repeated{
												{
													Message: specs.Message{
														"result": {
															Name:  "result",
															Path:  "repeated.repeated.result",
															Label: labels.Optional,
															Template: specs.Template{
																Scalar: &specs.Scalar{
																	Default: "hello repeated nested world",
																	Type:    types.String,
																},
															},
														},
														"active": {
															Name:  "active",
															Path:  "repeated.repeated.active",
															Label: labels.Optional,
															Template: specs.Template{
																Scalar: &specs.Scalar{
																	Default: false,
																	Type:    types.Bool,
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
					},
				},
			},
		},
	}
}

func NewMockCall(name string) *specs.Node {
	return &specs.Node{
		ID: name,
		OnError: &specs.OnError{
			Response: &specs.ParameterMap{
				Header: specs.Header{
					"cookie": &specs.Property{
						Path:  "cookie",
						Label: labels.Optional,
						Template: specs.Template{
							Scalar: &specs.Scalar{
								Default: "mnomnom",
								Type:    types.String,
							},
						},
					},
				},
				Property: NewResultMockProperty(),
			},
			Status: &specs.Property{
				Template: specs.Template{
					Scalar: &specs.Scalar{
						Type: types.Int64,
					},
				},
			},
			Message: &specs.Property{
				Template: specs.Template{
					Message: specs.Message{},
				},
			},
			Params: map[string]*specs.Property{
				"sample": {
					Path:  "sample",
					Label: labels.Optional,
					Template: specs.Template{
						Scalar: &specs.Scalar{
							Default: "mock",
							Type:    types.String,
						},
					},
				},
			},
		},
		Call: &specs.Call{
			Request: &specs.ParameterMap{
				Header: specs.Header{
					"cookie": &specs.Property{
						Path:  "cookie",
						Label: labels.Optional,
						Template: specs.Template{
							Scalar: &specs.Scalar{
								Default: "mnomnom",
								Type:    types.String,
							},
						},
					},
				},
				Stack: map[string]*specs.Property{
					name + "_request": {
						Label: labels.Optional,
						Template: specs.Template{
							Message: specs.Message{
								"nested": {
									Name:  "nested",
									Path:  "nested",
									Label: labels.Optional,
									Template: specs.Template{
										Scalar: &specs.Scalar{
											Type:    types.String,
											Default: "nested string",
										},
									},
								},
							},
						},
					},
				},
				Params: map[string]*specs.Property{
					"message": {
						Path:  "message",
						Label: labels.Optional,
						Template: specs.Template{
							Scalar: &specs.Scalar{
								Default: "default message",
								Type:    types.String,
							},
						},
					},
					"name": {
						Path:  "message",
						Label: labels.Optional,
						Template: specs.Template{
							Scalar: &specs.Scalar{
								Default: "default name",
								Type:    types.String,
							},
						},
					},
					"reference": {
						Path: "reference",
						Template: specs.Template{
							Reference: &specs.PropertyReference{
								Resource: name,
								Path:     "result",
							},
						},
					},
				},
				Property: NewInputMockProperty(),
			},
			Response: &specs.ParameterMap{
				Stack: map[string]*specs.Property{
					name + "_response": {
						Label: labels.Optional,
						Template: specs.Template{
							Message: specs.Message{
								"nested": {
									Name:  "nested",
									Path:  "nested",
									Label: labels.Optional,
									Template: specs.Template{
										Scalar: &specs.Scalar{
											Default: "nested string",
											Type:    types.String,
										},
									},
								},
							},
						},
					},
				},
				Header: specs.Header{
					"cookie": &specs.Property{
						Path:  "cookie",
						Label: labels.Optional,
						Template: specs.Template{
							Scalar: &specs.Scalar{
								Default: "mnomnom",
								Type:    types.String,
							},
						},
					},
				},
				Property: NewResultMockProperty(),
			},
		},
	}
}

func NewMockFlow(name string) *specs.Flow {
	return &specs.Flow{
		Name: name,
		Input: &specs.ParameterMap{
			Property: NewInputMockProperty(),
		},
		OnError: &specs.OnError{
			Response: &specs.ParameterMap{
				Property: NewResultMockProperty(),
				Params: map[string]*specs.Property{
					"message": {
						Path:  "message",
						Label: labels.Optional,
						Template: specs.Template{
							Scalar: &specs.Scalar{
								Default: "default message",
								Type:    types.String,
							},
						},
					},
					"name": {
						Path:  "message",
						Label: labels.Optional,
						Template: specs.Template{
							Scalar: &specs.Scalar{
								Default: "default name",
								Type:    types.String,
							},
						},
					},
					"reference": {
						Path: "reference",
						Template: specs.Template{
							Reference: &specs.PropertyReference{
								Resource: name,
								Path:     "message",
							},
						},
					},
				},
			},
		},
		Nodes: specs.NodeList{
			NewMockCall("first"),
			NewMockCall("second"),
			NewMockCall("third"),
		},
		Output: &specs.ParameterMap{
			Property: NewResultMockProperty(),
		},
	}
}

func TestGetAvailableResources(t *testing.T) {
	tests := map[string]func() ([]string, map[string]ReferenceMap){
		"input and first": func() ([]string, map[string]ReferenceMap) {
			flow := NewMockFlow("first")
			expected := []string{template.StackResource, template.ErrorResource, "input", "first", "second", "third"}

			result := GetAvailableResources(flow, "second")
			return expected, result
		},
		"input": func() ([]string, map[string]ReferenceMap) {
			flow := NewMockFlow("first")
			expected := []string{template.StackResource, template.ErrorResource, "input", "first", "second", "third"}

			result := GetAvailableResources(flow, "first")
			return expected, result
		},
		"output": func() ([]string, map[string]ReferenceMap) {
			flow := NewMockFlow("first")
			expected := []string{template.StackResource, template.ErrorResource, "input", "first", "second", "third"}

			result := GetAvailableResources(flow, "output")
			return expected, result
		},
		"output only": func() ([]string, map[string]ReferenceMap) {
			flow := NewMockFlow("first")

			flow.OnError = nil
			flow.Input = nil
			flow.Nodes = nil
			flow.Output = &specs.ParameterMap{
				Stack: map[string]*specs.Property{
					"hash": NewResultMockProperty(),
				},
			}

			expected := []string{template.StackResource}

			result := GetAvailableResources(flow, "output")
			return expected, result
		},
		"stack lookup request": func() ([]string, map[string]ReferenceMap) {
			flow := NewMockFlow("first")
			expected := []string{template.StackResource, template.ErrorResource, "input", "first", "second", "third"}

			flow.Nodes[0].Call.Request.Stack = map[string]*specs.Property{
				"ref": {
					Path: "ref",
				},
			}

			result := GetAvailableResources(flow, "output")
			return expected, result
		},
		"stack lookup response": func() ([]string, map[string]ReferenceMap) {
			flow := NewMockFlow("first")
			expected := []string{template.StackResource, template.ErrorResource, "input", "first", "second", "third"}

			flow.Nodes[0].Call.Response.Stack = map[string]*specs.Property{
				"ref": {
					Path: "ref",
				},
			}

			result := GetAvailableResources(flow, "output")
			return expected, result
		},
		"stack lookup intermediate": func() ([]string, map[string]ReferenceMap) {
			flow := NewMockFlow("first")
			expected := []string{template.StackResource, template.ErrorResource, "input", "first", "second", "third"}

			flow.Nodes[0].Intermediate = flow.Nodes[0].Call.Request
			flow.Nodes[0].Call = nil
			flow.Nodes[0].Rollback = nil

			flow.Nodes[0].Intermediate.Stack = map[string]*specs.Property{
				"ref": {
					Path: "ref",
				},
			}

			result := GetAvailableResources(flow, "output")
			return expected, result
		},
	}

	for key, test := range tests {
		t.Run(key, func(t *testing.T) {
			expected, result := test()
			if len(expected) != len(result) {
				t.Errorf("unexpected result %+v, expected %+v", result, expected)
			}

			for _, resource := range expected {
				_, has := result[resource]
				if !has {
					t.Errorf("expected resource not found %s, %+v", resource, result)
				}
			}
		})
	}
}

func TestSkipMissingParameters(t *testing.T) {
	flow := NewMockFlow("first")

	flow.Input = nil
	flow.Nodes[0].Call.Request = nil
	flow.Nodes[1].Call.Response = nil

	type test struct {
		name  string
		empty []string
	}

	tests := []test{
		{
			name: template.ErrorResource,
		},
		{
			name: template.StackResource,
		},
		{
			name:  "first",
			empty: []string{template.RequestResource, template.HeaderResource},
		},
		{
			name:  "second",
			empty: []string{template.RequestResource, template.HeaderResource},
		},
		{
			name: "third",
		},
	}

	result := GetAvailableResources(flow, "output")

	if len(result) != len(tests) {
		t.Fatalf("unexpected result (%d, %d) %+v, expected %+v", len(result), len(tests), result, tests)
	}

	for _, test := range tests {
		resource, has := result[test.name]
		if !has {
			t.Fatalf("expected resource not found %s, %+v", test.name, result)
		}

	check:
		for key, val := range resource {
			if val == nil {
				for _, known := range test.empty {
					if key == known {
						continue check
					}
				}

				t.Fatalf("unknown empty resource %s", key)
			}
		}
	}
}

func NewPropertyReference(resource string, path string) *specs.PropertyReference {
	return &specs.PropertyReference{
		Resource: resource,
		Path:     path,
	}
}

func TestGetResourceOutputReference(t *testing.T) {
	var (
		flow       = NewMockFlow("first")
		resources  = GetNextResource(flow, template.OutputResource)
		breakpoint = "first"
		tests      = map[[2]string]*specs.Property{
			{"", "result"}:                                 flow.Nodes[0].Call.Response.Property.Message["result"],
			{"input", "message"}:                           flow.Input.Property.Template.Message["message"],
			{"first", "result"}:                            flow.Nodes[0].Call.Response.Property.Message["result"],
			{"first", "nested.result"}:                     flow.Nodes[0].Call.Response.Property.Message["nested"].Message["result"],
			{"first", "nested.nested.result"}:              flow.Nodes[0].Call.Response.Property.Message["nested"].Message["nested"].Message["result"],
			{"first.header", "cookie"}:                     flow.Nodes[0].Call.Response.Header["cookie"],
			{"first.request", "message"}:                   flow.Nodes[0].Call.Request.Property.Message["message"],
			{"first.request", "nested"}:                    flow.Nodes[0].Call.Request.Property.Message["nested"],
			{"first.request", "nested.message"}:            flow.Nodes[0].Call.Request.Property.Message["nested"].Message["message"],
			{"first.request", "nested.nested.message"}:     flow.Nodes[0].Call.Request.Property.Message["nested"].Message["nested"].Message["message"],
			{"first.request", "nested.repeated.message"}:   flow.Nodes[0].Call.Request.Property.Message["nested"].Message["repeated"].Repeated[0].Message["message"],
			{"first.response", "result"}:                   flow.Nodes[0].Call.Response.Property.Message["result"],
			{"first.response", "nested.nested.result"}:     flow.Nodes[0].Call.Response.Property.Message["nested"].Message["nested"].Message["result"],
			{"first.response", "nested.repeated.result"}:   flow.Nodes[0].Call.Response.Property.Message["nested"].Message["repeated"].Repeated[0].Message["result"],
			{"first.response", "repeated.repeated.result"}: flow.Nodes[0].Call.Response.Property.Message["repeated"].Repeated[0].Message["repeated"].Repeated[0].Message["result"],
			{"first.response", "nested.repeated.result"}:   flow.Nodes[0].Call.Response.Property.Message["nested"].Message["repeated"].Repeated[0].Message["result"],
			{"first.params", "message"}:                    flow.Nodes[0].Call.Request.Params["message"],
			{"first.params", "name"}:                       flow.Nodes[0].Call.Request.Params["name"],
			{"first.params", "reference"}:                  flow.Nodes[0].Call.Request.Params["reference"],
			{"stack.first_request", "."}:                   flow.Nodes[0].Call.Request.Stack["first_request"],
			{"stack.first_request", "nested"}:              flow.Nodes[0].Call.Request.Stack["first_request"].Message["nested"],
			{"stack.first_response", "."}:                  flow.Nodes[0].Call.Response.Stack["first_response"],
			{"stack.first_response", "nested"}:             flow.Nodes[0].Call.Response.Stack["first_response"].Message["nested"],
		}
	)

	for input, expected := range tests {
		var (
			reference = NewPropertyReference(input[0], input[1])
			title     = reference.String()
		)

		t.Run(title, func(t *testing.T) {
			var (
				references = GetAvailableResources(flow, resources)
				result     = GetResourceReference(reference, references, breakpoint)
			)

			if result == nil {
				t.Fatalf("unexpected empty result on lookup '%s', expected '%+v'", input, expected)
			}

			if result.Path != expected.Path {
				t.Fatalf("unexpected result '%+v', expected '%+v'", result, expected)
			}
		})
	}
}

func TestGetResourceReference(t *testing.T) {
	flow := NewMockFlow("first")

	tests := map[[2]string]*specs.Property{
		{"error", "status"}:           flow.Nodes[0].OnError.Status,
		{"error", "message"}:          flow.Nodes[0].OnError.Message,
		{"error.response", "status"}:  flow.Nodes[0].OnError.Status,
		{"error.response", "message"}: flow.Nodes[0].OnError.Message,
		{"error.params", "sample"}:    flow.Nodes[0].OnError.Params["sample"],
		{"first.error", "result"}:     flow.Nodes[0].OnError.Response.Property.Template.Message["result"],
	}

	for input, expected := range tests {
		var (
			reference = NewPropertyReference(input[0], input[1])
			title     = reference.String()
		)

		t.Run(title, func(t *testing.T) {
			var (
				resource, _ = ParseResource(reference.Resource)
				references  = GetAvailableResources(flow, "first")
				result      = GetResourceReference(reference, references, resource)
			)

			if result == nil {
				t.Fatalf("unexpected empty result on lookup '%s', expected '%+v'", input, expected)
			}

			if result.Path != expected.Path {
				t.Fatalf("unexpected result '%+v', expected '%+v'", result, expected)
			}
		})
	}
}

func TestGetIntermediateResourceReference(t *testing.T) {
	flow := NewMockFlow("first")

	// set intermediate property
	flow.Nodes[0].Intermediate = flow.Nodes[0].Call.Response
	flow.Nodes[0].Call = nil
	flow.Nodes[0].Rollback = nil

	tests := map[[2]string]*specs.Property{
		{"first", "result"}:                            flow.Nodes[0].Intermediate.Property.Message["result"],
		{"first", "nested.result"}:                     flow.Nodes[0].Intermediate.Property.Message["nested"].Message["result"],
		{"first", "nested.nested.result"}:              flow.Nodes[0].Intermediate.Property.Message["nested"].Message["nested"].Message["result"],
		{"first.header", "cookie"}:                     flow.Nodes[0].Intermediate.Header["cookie"],
		{"first.response", "result"}:                   flow.Nodes[0].Intermediate.Property.Message["result"],
		{"first.response", "nested.nested.result"}:     flow.Nodes[0].Intermediate.Property.Message["nested"].Message["nested"].Message["result"],
		{"first.response", "nested.repeated.result"}:   flow.Nodes[0].Intermediate.Property.Message["nested"].Message["repeated"].Repeated[0].Message["result"],
		{"first.response", "repeated.repeated.result"}: flow.Nodes[0].Intermediate.Property.Message["repeated"].Repeated[0].Message["repeated"].Repeated[0].Message["result"],
	}

	for input, expected := range tests {
		var (
			reference = NewPropertyReference(input[0], input[1])
			title     = reference.String()
		)

		t.Run(title, func(t *testing.T) {
			var (
				resource, _ = ParseResource(reference.Resource)
				references  = GetAvailableResources(flow, "first")
				result      = GetResourceReference(reference, references, resource)
			)

			if result == nil {
				t.Fatalf("unexpected empty result on lookup '%s', expected '%+v'", reference, expected)
			}

			if result.Path != expected.Path {
				t.Fatalf("unexpected result '%+v', expected '%+v'", result, expected)
			}
		})
	}
}

func TestGetUnknownResourceReference(t *testing.T) {
	var (
		flow       = NewMockFlow("first")
		references = GetAvailableResources(flow, "output")
		breakpoint = "first"
		tests      = map[string]*specs.PropertyReference{
			"unknown": NewPropertyReference("unknown", "unknown"),
		}
	)

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := GetResourceReference(test, references, breakpoint)
			if result != nil {
				t.Fatalf("unexpected result")
			}
		})
	}
}

func TestHeaderLookup(t *testing.T) {
	type test struct {
		path   string
		header specs.Header
	}

	tests := map[string]*test{
		"simple": {
			path: "key",
			header: specs.Header{
				"key": &specs.Property{
					Path: "key",
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			resolver := HeaderLookup(test.header)
			prop := resolver(test.path)
			if prop == nil {
				t.Fatalf("unexpected result expected a prop to be returned for '%s'", test.path)
			}
		})
	}
}

func TestUnknownHeaderLookup(t *testing.T) {
	type test struct {
		path   string
		header specs.Header
	}

	tests := map[string]*test{
		"simple": {
			path:   "key",
			header: specs.Header{},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			resolver := HeaderLookup(test.header)
			prop := resolver(test.path)
			if prop != nil {
				t.Fatalf("unexpected result expected nil to be returned for '%s'", test.path)
			}
		})
	}
}

func TestPropertyLookup(t *testing.T) {
	type test struct {
		path  string
		param *specs.Property
	}

	tests := map[string]*test{
		"self reference": {
			path: ".",
			param: &specs.Property{
				Path: "key",
			},
		},
		"simple": {
			path: "key",
			param: &specs.Property{
				Path: "key",
			},
		},
		"nested": {
			path: "key.nested",
			param: &specs.Property{
				Path: "key",
				Template: specs.Template{
					Message: specs.Message{
						"nested": {
							Name: "nested",
							Path: "key.nested",
						},
					},
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			lookup := PropertyLookup(test.param)
			result := lookup(test.path)
			if result == nil {
				t.Fatal("unexpected empty result")
			}
		})
	}
}

func TestResolveSelfReference(t *testing.T) {
	type test struct {
		path     string
		expected string
		resource string
	}

	tests := map[string]*test{
		"self reference": {
			path:     ".request",
			resource: "input",
			expected: "input.request",
		},
		"resource reference": {
			path:     "input.request",
			resource: "first",
			expected: "input.request",
		},
		"broken path": {
			path:     "input.",
			resource: "first",
			expected: "input.",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := ResolveSelfReference(test.path, test.resource)
			if result != test.expected {
				t.Fatalf("unexpected result '%s', expected '%s'", result, test.expected)
			}
		})
	}
}

func TestGetReference(t *testing.T) {
	path := "key"
	prop := "input.request"

	references := ReferenceMap{
		prop: PropertyLookup(&specs.Property{Path: path}),
	}

	result := GetReference(path, prop, references)
	if result == nil {
		t.Fatal("unexpected empty result")
	}
}

func TestUnknownReference(t *testing.T) {
	path := "key"
	prop := "input.request"

	references := ReferenceMap{
		prop: PropertyLookup(&specs.Property{Path: path}),
	}

	result := GetReference(path, "unknown", references)
	if result != nil {
		t.Fatal("unexpected result")
	}
}
