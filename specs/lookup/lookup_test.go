package lookup

import (
	"testing"

	"github.com/jexia/maestro/specs"
	"github.com/jexia/maestro/specs/types"
)

func TestGetFlow(t *testing.T) {
	manifest := specs.Manifest{
		Flows: []*specs.Flow{
			{
				Name: "first",
			},
			{
				Name: "second",
			},
		},
	}

	tests := map[string]*specs.Flow{
		"second":  manifest.Flows[1],
		"first":   manifest.Flows[0],
		"unknown": nil,
	}

	for input, expected := range tests {
		result := GetFlow(manifest, input)
		if result != expected {
			t.Errorf("unexpected result %+v, expected %+v", result, expected)
		}
	}
}

func TestGetDefaultProp(t *testing.T) {
	tests := map[string]string{
		specs.InputResource: specs.ResourceRequest,
		"unknown":           specs.ResourceResponse,
	}

	for input, expected := range tests {
		result := GetDefaultProp(input)
		if result != expected {
			t.Errorf("unexpected result %s, expected %s", result, expected)
		}
	}
}

func NewInputMockProperty() *specs.Property {
	return &specs.Property{
		Path:  "",
		Type:  types.TypeMessage,
		Label: types.LabelOptional,
		Nested: map[string]*specs.Property{
			"message": {
				Path:    "message",
				Default: "hello world",
				Type:    types.TypeString,
				Label:   types.LabelOptional,
			},
			"active": {
				Path:    "active",
				Default: false,
				Type:    types.TypeBool,
				Label:   types.LabelOptional,
			},
			"nested": {
				Path:  "nested",
				Type:  types.TypeMessage,
				Label: types.LabelOptional,
				Nested: map[string]*specs.Property{
					"message": {
						Path:    "nested.message",
						Default: "hello world",
						Type:    types.TypeString,
						Label:   types.LabelOptional,
					},
					"active": {
						Path:    "nested.active",
						Default: false,
						Type:    types.TypeBool,
						Label:   types.LabelOptional,
					},
					"nested": {
						Path:  "nested.nested",
						Type:  types.TypeMessage,
						Label: types.LabelOptional,
						Nested: map[string]*specs.Property{
							"message": {
								Path:    "nested.nested.message",
								Default: "hello world",
								Type:    types.TypeString,
								Label:   types.LabelOptional,
							},
							"active": {
								Path:    "nested.nested.active",
								Default: false,
								Type:    types.TypeBool,
								Label:   types.LabelOptional,
							},
						},
					},
					"repeated": {
						Path:  "nested.repeated",
						Type:  types.TypeMessage,
						Label: types.LabelRepeated,
						Nested: map[string]*specs.Property{
							"message": {
								Path:    "nested.repeated.message",
								Default: "hello world",
								Type:    types.TypeString,
								Label:   types.LabelOptional,
							},
							"active": {
								Path:    "nested.repeated.active",
								Default: false,
								Type:    types.TypeBool,
								Label:   types.LabelOptional,
							},
						},
					},
				},
			},
			"repeated": {
				Path:  "repeated",
				Type:  types.TypeMessage,
				Label: types.LabelRepeated,
				Nested: map[string]*specs.Property{
					"message": {
						Path:    "message",
						Default: "hello world",
						Type:    types.TypeString,
						Label:   types.LabelOptional,
					},
					"active": {
						Path:    "active",
						Default: false,
						Type:    types.TypeBool,
						Label:   types.LabelOptional,
					},
					"repeated": {
						Path:  "repeated.repeated",
						Type:  types.TypeMessage,
						Label: types.LabelRepeated,
						Nested: map[string]*specs.Property{
							"repeated.message": {
								Path:    "message",
								Default: "hello world",
								Type:    types.TypeString,
								Label:   types.LabelOptional,
							},
							"repeated.active": {
								Path:    "active",
								Default: false,
								Type:    types.TypeBool,
								Label:   types.LabelOptional,
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
		Type:  types.TypeMessage,
		Label: types.LabelOptional,
		Nested: map[string]*specs.Property{
			"result": {
				Path:    "result",
				Default: "hello world",
				Type:    types.TypeString,
				Label:   types.LabelOptional,
			},
			"active": {
				Path:    "active",
				Default: false,
				Type:    types.TypeBool,
				Label:   types.LabelOptional,
			},
			"nested": {
				Path:  "nested",
				Type:  types.TypeMessage,
				Label: types.LabelOptional,
				Nested: map[string]*specs.Property{
					"result": {
						Path:    "nested.result",
						Default: "hello world",
						Type:    types.TypeString,
						Label:   types.LabelOptional,
					},
					"active": {
						Path:    "nested.active",
						Default: false,
						Type:    types.TypeBool,
						Label:   types.LabelOptional,
					},
					"nested": {
						Path:  "nested.nested",
						Type:  types.TypeMessage,
						Label: types.LabelOptional,
						Nested: map[string]*specs.Property{
							"result": {
								Path:    "nested.nested.result",
								Default: "hello world",
								Type:    types.TypeString,
								Label:   types.LabelOptional,
							},
							"active": {
								Path:    "nested.nested.active",
								Default: false,
								Type:    types.TypeBool,
								Label:   types.LabelOptional,
							},
						},
					},
					"repeated": {
						Path:  "nested.repeated",
						Type:  types.TypeMessage,
						Label: types.LabelRepeated,
						Nested: map[string]*specs.Property{
							"result": {
								Path:    "nested.repeated.result",
								Default: "hello world",
								Type:    types.TypeString,
								Label:   types.LabelOptional,
							},
							"active": {
								Path:    "nested.repeated.active",
								Default: false,
								Type:    types.TypeBool,
								Label:   types.LabelOptional,
							},
						},
					},
				},
			},
			"repeated": {
				Path:  "repeated",
				Type:  types.TypeMessage,
				Label: types.LabelRepeated,
				Nested: map[string]*specs.Property{
					"result": {
						Path:    "message",
						Default: "hello world",
						Type:    types.TypeString,
						Label:   types.LabelOptional,
					},
					"active": {
						Path:    "active",
						Default: false,
						Type:    types.TypeBool,
						Label:   types.LabelOptional,
					},
					"repeated": {
						Path:  "repeated.repeated",
						Type:  types.TypeMessage,
						Label: types.LabelRepeated,
						Nested: map[string]*specs.Property{
							"result": {
								Path:    "repeated.repeated.result",
								Default: "hello world",
								Type:    types.TypeString,
								Label:   types.LabelOptional,
							},
							"active": {
								Path:    "repeated.repeated.active",
								Default: false,
								Type:    types.TypeBool,
								Label:   types.LabelOptional,
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
		Name: name,
		Call: &specs.Call{
			Request: &specs.ParameterMap{
				Header: specs.Header{
					"cookie": &specs.Property{
						Path:    "cookie",
						Default: "mnomnom",
						Type:    types.TypeString,
						Label:   types.LabelOptional,
					},
				},
				Property: NewInputMockProperty(),
			},
			Response: &specs.ParameterMap{
				Header: specs.Header{
					"cookie": &specs.Property{
						Path:    "cookie",
						Default: "mnomnom",
						Type:    types.TypeString,
						Label:   types.LabelOptional,
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
		Nodes: []*specs.Node{
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
			expected := []string{"input", "first", "second"}

			result := GetAvailableResources(flow, "second")
			return expected, result
		},
		"input": func() ([]string, map[string]ReferenceMap) {
			flow := NewMockFlow("first")
			expected := []string{"input", "first"}

			result := GetAvailableResources(flow, "first")
			return expected, result
		},
		"output": func() ([]string, map[string]ReferenceMap) {
			flow := NewMockFlow("first")
			expected := []string{"input", "first", "second", "third"}

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
			name:  "first",
			empty: []string{specs.ResourceRequest, specs.ResourceHeader},
		},
		{
			name:  "second",
			empty: []string{specs.ResourceRequest, specs.ResourceHeader},
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

func TestGetResourceReference(t *testing.T) {
	flow := NewMockFlow("first")
	references := GetAvailableResources(flow, "output")

	tests := map[*specs.PropertyReference]*specs.Property{
		NewPropertyReference("input", "message"):                           flow.Input.Property.Nested["message"],
		NewPropertyReference("first", "result"):                            flow.Nodes[0].Call.Response.Property.Nested["result"],
		NewPropertyReference("first.response", "result"):                   flow.Nodes[0].Call.Response.Property.Nested["result"],
		NewPropertyReference("first.request", "message"):                   flow.Nodes[0].Call.Request.Property.Nested["message"],
		NewPropertyReference("first.request", "message"):                   flow.Nodes[0].Call.Request.Property.Nested["message"],
		NewPropertyReference("first.header", "cookie"):                     flow.Nodes[0].Call.Response.Header["cookie"],
		NewPropertyReference("first.request", "nested"):                    flow.Nodes[0].Call.Request.Property.Nested["nested"],
		NewPropertyReference("first.request", "nested.message"):            flow.Nodes[0].Call.Request.Property.Nested["nested"].Nested["message"],
		NewPropertyReference("first", "nested.result"):                     flow.Nodes[0].Call.Response.Property.Nested["nested"].Nested["result"],
		NewPropertyReference("first", "nested.nested.result"):              flow.Nodes[0].Call.Response.Property.Nested["nested"].Nested["nested"].Nested["result"],
		NewPropertyReference("first.response", "nested.nested.result"):     flow.Nodes[0].Call.Response.Property.Nested["nested"].Nested["nested"].Nested["result"],
		NewPropertyReference("first.request", "nested.repeated.message"):   flow.Nodes[0].Call.Request.Property.Nested["nested"].Nested["repeated"].Nested["message"],
		NewPropertyReference("first.response", "nested.repeated.result"):   flow.Nodes[0].Call.Response.Property.Nested["nested"].Nested["repeated"].Nested["result"],
		NewPropertyReference("first.response", "repeated.repeated.result"): flow.Nodes[0].Call.Response.Property.Nested["repeated"].Nested["repeated"].Nested["result"],
		NewPropertyReference("first.response", "nested.repeated.result"):   flow.Nodes[0].Call.Response.Property.Nested["nested"].Nested["repeated"].Nested["result"],
		NewPropertyReference("first.response", "nested.nested.result"):     flow.Nodes[0].Call.Response.Property.Nested["nested"].Nested["nested"].Nested["result"],
		NewPropertyReference("first.response", "nested.repeated.result"):   flow.Nodes[0].Call.Response.Property.Nested["nested"].Nested["repeated"].Nested["result"],
		NewPropertyReference("first.request", "nested.repeated.message"):   flow.Nodes[0].Call.Request.Property.Nested["nested"].Nested["repeated"].Nested["message"],
		NewPropertyReference("first.request", "nested.nested.message"):     flow.Nodes[0].Call.Request.Property.Nested["nested"].Nested["nested"].Nested["message"],
		NewPropertyReference("first.request", "nested.repeated.message"):   flow.Nodes[0].Call.Request.Property.Nested["nested"].Nested["repeated"].Nested["message"],
	}

	for input, expected := range tests {
		t.Run(input.String(), func(t *testing.T) {
			result := GetResourceReference(input, references, "output")
			if result == nil {
				t.Fatalf("unexpected result on lookup %s, expected %+v", input, expected)
			}

			if result.Path != expected.Path {
				t.Fatalf("unexpected result %+v, expected %+v", result, expected)
			}
		})
	}
}
