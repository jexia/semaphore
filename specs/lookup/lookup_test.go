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
		"second": manifest.Flows[1],
		"first":  manifest.Flows[0],
		"unkown": nil,
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
		"unkown":            specs.ResourceResponse,
	}

	for input, expected := range tests {
		result := GetDefaultProp(input)
		if result != expected {
			t.Errorf("unexpected result %s, expected %s", result, expected)
		}
	}
}

func NewInputMockProperties(path string) map[string]*specs.Property {
	return map[string]*specs.Property{
		"message": &specs.Property{
			Path:    specs.JoinPath(path, "message"),
			Default: "hello world",
			Type:    types.TypeString,
		},
		"active": &specs.Property{
			Path:    specs.JoinPath(path, "active"),
			Default: false,
			Type:    types.TypeBool,
		},
	}
}

func NewOutputMockProperties(path string) map[string]*specs.Property {
	return map[string]*specs.Property{
		"result": &specs.Property{
			Path:    specs.JoinPath(path, "result"),
			Default: "hello world",
			Type:    types.TypeString,
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
					},
				},
				Properties: NewInputMockProperties(""),
				Nested: map[string]*specs.NestedParameterMap{
					"nested": &specs.NestedParameterMap{
						Path:       "nested",
						Properties: NewInputMockProperties("nested"),
						Nested: map[string]*specs.NestedParameterMap{
							"nested": &specs.NestedParameterMap{
								Path:       "nested.nested",
								Properties: NewInputMockProperties("nested.nested"),
							},
						},
						Repeated: map[string]*specs.RepeatedParameterMap{
							"repeated": &specs.RepeatedParameterMap{
								Path:       "nested.repeated",
								Properties: NewInputMockProperties("nested.repeated"),
							},
						},
					},
				},
				Repeated: map[string]*specs.RepeatedParameterMap{
					"repeated": &specs.RepeatedParameterMap{
						Path:       "repeated",
						Properties: NewInputMockProperties("repeated"),
						Nested: map[string]*specs.NestedParameterMap{
							"nested": &specs.NestedParameterMap{
								Path:       "repeated.nested",
								Properties: NewInputMockProperties("repeated.nested"),
							},
						},
						Repeated: map[string]*specs.RepeatedParameterMap{
							"repeated": &specs.RepeatedParameterMap{
								Path:       "repeated.repeated",
								Properties: NewInputMockProperties("repeated.nested"),
							},
						},
					},
				},
			},
			Response: &specs.ParameterMap{
				Header: specs.Header{
					"cookie": &specs.Property{
						Path:    "cookie",
						Default: "mnomnom",
						Type:    types.TypeString,
					},
				},
				Properties: NewOutputMockProperties(""),
				Nested: map[string]*specs.NestedParameterMap{
					"nested": &specs.NestedParameterMap{
						Path:       "nested",
						Properties: NewOutputMockProperties("nested"),
						Nested: map[string]*specs.NestedParameterMap{
							"nested": &specs.NestedParameterMap{
								Path:       "nested.nested",
								Properties: NewOutputMockProperties("nested.nested"),
							},
						},
						Repeated: map[string]*specs.RepeatedParameterMap{
							"repeated": &specs.RepeatedParameterMap{
								Path:       "nested.repeated",
								Properties: NewOutputMockProperties("nested.repeated"),
							},
						},
					},
				},
				Repeated: map[string]*specs.RepeatedParameterMap{
					"repeated": &specs.RepeatedParameterMap{
						Path:       "repeated",
						Properties: NewOutputMockProperties("repeated"),
						Nested: map[string]*specs.NestedParameterMap{
							"nested": &specs.NestedParameterMap{
								Path:       "repeated.nested",
								Properties: NewOutputMockProperties("repeated.nested"),
							},
						},
						Repeated: map[string]*specs.RepeatedParameterMap{
							"repeated": &specs.RepeatedParameterMap{
								Path:       "repeated.repeated",
								Properties: NewOutputMockProperties("repeated.repeated"),
							},
						},
					},
				},
			},
		},
	}
}

func NewMockFlow(name string) *specs.Flow {
	return &specs.Flow{
		Name: name,
		Input: &specs.ParameterMap{
			Properties: NewInputMockProperties(""),
		},
		Nodes: []*specs.Node{
			NewMockCall("first"),
			NewMockCall("second"),
			NewMockCall("third"),
		},
		Output: &specs.ParameterMap{
			Properties: NewOutputMockProperties(""),
		},
	}
}

func TestGetAvailableResources(t *testing.T) {
	tests := map[string]func() ([]string, map[string]ReferenceMap){
		"input and first": func() ([]string, map[string]ReferenceMap) {
			flow := NewMockFlow("first")
			expected := []string{"input", "first"}

			result := GetAvailableResources(flow, "second")
			return expected, result
		},
		"input": func() ([]string, map[string]ReferenceMap) {
			flow := NewMockFlow("first")
			expected := []string{"input", "first"}

			result := GetAvailableResources(flow, "second")
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
		test{
			name:  "first",
			empty: []string{specs.ResourceRequest, specs.ResourceRequestHeader},
		},
		test{
			name:  "second",
			empty: []string{specs.ResourceRequest, specs.ResourceRequestHeader},
		},
		test{
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

				t.Fatalf("unkown empty resource %s", key)
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

	tests := map[*specs.PropertyReference]Reference{
		NewPropertyReference("input", "message"):                           flow.Input.Properties["message"],
		NewPropertyReference("first", "result"):                            flow.Nodes[0].Call.Response.Properties["result"],
		NewPropertyReference("first.response", "result"):                   flow.Nodes[0].Call.Response.Properties["result"],
		NewPropertyReference("first.request", "message"):                   flow.Nodes[0].Call.Request.Properties["message"],
		NewPropertyReference("first.request", "message"):                   flow.Nodes[0].Call.Request.Properties["message"],
		NewPropertyReference("first.request.header", "cookie"):             flow.Nodes[0].Call.Request.Header["cookie"],
		NewPropertyReference("first.request", "nested"):                    flow.Nodes[0].Call.Request.Nested["nested"],
		NewPropertyReference("first.request", "nested.message"):            flow.Nodes[0].Call.Request.Nested["nested"].Properties["message"],
		NewPropertyReference("first", "nested.result"):                     flow.Nodes[0].Call.Response.Nested["nested"].Properties["result"],
		NewPropertyReference("first", "nested.nested.result"):              flow.Nodes[0].Call.Response.Nested["nested"].Nested["nested"].Properties["result"],
		NewPropertyReference("first.response", "nested.nested.result"):     flow.Nodes[0].Call.Response.Nested["nested"].Nested["nested"].Properties["result"],
		NewPropertyReference("first.request", "nested.repeated.message"):   flow.Nodes[0].Call.Request.Nested["nested"].Repeated["repeated"].Properties["message"],
		NewPropertyReference("first.response", "nested.repeated.result"):   flow.Nodes[0].Call.Response.Nested["nested"].Repeated["repeated"].Properties["result"],
		NewPropertyReference("first.response", "repeated.repeated.result"): flow.Nodes[0].Call.Response.Repeated["repeated"].Repeated["repeated"].Properties["result"],
		NewPropertyReference("first.response", "nested.repeated.result"):   flow.Nodes[0].Call.Response.Nested["nested"].Repeated["repeated"].Properties["result"],
		NewPropertyReference("first.response", "nested.nested.result"):     flow.Nodes[0].Call.Response.Nested["nested"].Nested["nested"].Properties["result"],
		NewPropertyReference("first.response", "nested.repeated.result"):   flow.Nodes[0].Call.Response.Nested["nested"].Repeated["repeated"].Properties["result"],
		NewPropertyReference("first.request", "nested.repeated.message"):   flow.Nodes[0].Call.Request.Nested["nested"].Repeated["repeated"].Properties["message"],
		NewPropertyReference("first.request", "nested.nested.message"):     flow.Nodes[0].Call.Request.Nested["nested"].Nested["nested"].Properties["message"],
		NewPropertyReference("first.request", "nested.repeated.message"):   flow.Nodes[0].Call.Request.Nested["nested"].Repeated["repeated"].Properties["message"],
	}

	for input, expected := range tests {
		t.Run(input.String(), func(t *testing.T) {
			result := GetResourceReference(input, references)
			if result == nil {
				t.Fatalf("unexpected result expected %+v", expected)
			}

			if result.GetPath() != expected.GetPath() {
				t.Fatalf("unexpected result %+v, expected %+v", result, expected)
			}
		})
	}
}
