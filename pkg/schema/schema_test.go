package schema

import (
	"testing"

	"github.com/jexia/semaphore/pkg/core/instance"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/labels"
	"github.com/jexia/semaphore/pkg/specs/types"
)

func NewMockServices() specs.ServiceList {
	return specs.ServiceList{
		&specs.Service{
			FullyQualifiedName: "com.mock.greeter",
			Package:            "com.mock",
			Name:               "greeter",
			Methods: []*specs.Method{
				{
					Name:   "echo",
					Input:  "com.mock.message",
					Output: "com.mock.message",
				},
				{
					Name:   "unkown",
					Input:  "com.mock.unkown",
					Output: "com.mock.unkown",
				},
			},
		},
	}
}

func NewMockSchemas() specs.Objects {
	return specs.Objects{
		"com.mock.message": &specs.Property{
			Type:  types.Message,
			Label: labels.Optional,
			Nested: map[string]*specs.Property{
				"value": {
					Type:  types.String,
					Label: labels.Optional,
				},
				"meta": {
					Type:  types.Message,
					Label: labels.Optional,
					Nested: map[string]*specs.Property{
						"id": {
							Type:  types.String,
							Label: labels.Optional,
						},
					},
				},
			},
		},
	}
}

func TestDefineSchemas(t *testing.T) {
	t.Parallel()

	tests := map[string]specs.FlowListInterface{
		"input": {
			&specs.Flow{
				Input: &specs.ParameterMap{
					Schema: "com.mock.message",
				},
			},
		},
		"output": {
			&specs.Flow{
				Output: &specs.ParameterMap{
					Schema: "com.mock.message",
				},
			},
		},
		"node call request": {
			&specs.Flow{
				Nodes: specs.NodeList{
					&specs.Node{
						Call: &specs.Call{
							Request: &specs.ParameterMap{
								Schema: "com.mock.message",
							},
						},
					},
				},
			},
		},
		"node call response": {
			&specs.Flow{
				Nodes: specs.NodeList{
					&specs.Node{
						Call: &specs.Call{
							Response: &specs.ParameterMap{
								Schema: "com.mock.message",
							},
						},
					},
				},
			},
		},
		"node service method": {
			&specs.Flow{
				Nodes: specs.NodeList{
					&specs.Node{
						Call: &specs.Call{
							Service: "com.mock.greeter",
							Method:  "echo",
							Request: &specs.ParameterMap{},
							Response: &specs.ParameterMap{
								Schema: "com.mock.message",
							},
						},
					},
				},
			},
		},
		"node on error": {
			&specs.Flow{
				Nodes: specs.NodeList{
					&specs.Node{
						OnError: &specs.OnError{
							Response: &specs.ParameterMap{
								Schema: "com.mock.message",
							},
						},
					},
				},
			},
		},
		"node condition": {
			&specs.Flow{
				Nodes: specs.NodeList{
					&specs.Node{
						Condition: &specs.Condition{
							Params: &specs.ParameterMap{
								Schema: "com.mock.message",
							},
						},
					},
				},
			},
		},
		"node rollback request": {
			&specs.Flow{
				Nodes: specs.NodeList{
					&specs.Node{
						Rollback: &specs.Call{
							Request: &specs.ParameterMap{
								Schema: "com.mock.message",
							},
						},
					},
				},
			},
		},
		"node rollback response": {
			&specs.Flow{
				Nodes: specs.NodeList{
					&specs.Node{
						Rollback: &specs.Call{
							Response: &specs.ParameterMap{
								Schema: "com.mock.message",
							},
						},
					},
				},
			},
		},
		"on error": {
			&specs.Flow{
				OnError: &specs.OnError{
					Response: &specs.ParameterMap{
						Schema: "com.mock.message",
					},
				},
			},
		},
	}

	for name, flows := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := instance.NewContext()
			services := NewMockServices()
			schemas := NewMockSchemas()

			err := Define(ctx, services, schemas, flows)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestDefineSchemasUnkown(t *testing.T) {
	t.Parallel()

	tests := map[string]specs.FlowListInterface{
		"input": {
			&specs.Flow{
				Input: &specs.ParameterMap{
					Schema: "com.mock.unkown",
				},
			},
		},
		"output": {
			&specs.Flow{
				Output: &specs.ParameterMap{
					Schema: "com.mock.unkown",
				},
			},
		},
		"node call request": {
			&specs.Flow{
				Nodes: specs.NodeList{
					&specs.Node{
						Call: &specs.Call{
							Request: &specs.ParameterMap{
								Schema: "com.mock.unkown",
							},
						},
					},
				},
			},
		},
		"node call response": {
			&specs.Flow{
				Nodes: specs.NodeList{
					&specs.Node{
						Call: &specs.Call{
							Response: &specs.ParameterMap{
								Schema: "com.mock.unkown",
							},
						},
					},
				},
			},
		},
		"node service method": {
			&specs.Flow{
				Nodes: specs.NodeList{
					&specs.Node{
						Call: &specs.Call{
							Service:  "com.mock.greeter",
							Method:   "undefined",
							Request:  &specs.ParameterMap{},
							Response: &specs.ParameterMap{},
						},
					},
				},
			},
		},
		"node service": {
			&specs.Flow{
				Nodes: specs.NodeList{
					&specs.Node{
						Call: &specs.Call{
							Service:  "com.mock.undefined",
							Method:   "echo",
							Request:  &specs.ParameterMap{},
							Response: &specs.ParameterMap{},
						},
					},
				},
			},
		},
		"node service output": {
			&specs.Flow{
				Nodes: specs.NodeList{
					&specs.Node{
						Call: &specs.Call{
							Service:  "com.mock.greeter",
							Method:   "unkown",
							Request:  &specs.ParameterMap{},
							Response: &specs.ParameterMap{},
						},
					},
				},
			},
		},
		"node on error": {
			&specs.Flow{
				Nodes: specs.NodeList{
					&specs.Node{
						OnError: &specs.OnError{
							Response: &specs.ParameterMap{
								Schema: "com.mock.unkown",
							},
						},
					},
				},
			},
		},
		"node condition": {
			&specs.Flow{
				Nodes: specs.NodeList{
					&specs.Node{
						Condition: &specs.Condition{
							Params: &specs.ParameterMap{
								Schema: "com.mock.unkown",
							},
						},
					},
				},
			},
		},
		"node rollback request": {
			&specs.Flow{
				Nodes: specs.NodeList{
					&specs.Node{
						Rollback: &specs.Call{
							Request: &specs.ParameterMap{
								Schema: "com.mock.unkown",
							},
						},
					},
				},
			},
		},
		"node rollback response": {
			&specs.Flow{
				Nodes: specs.NodeList{
					&specs.Node{
						Rollback: &specs.Call{
							Response: &specs.ParameterMap{
								Schema: "com.mock.unkown",
							},
						},
					},
				},
			},
		},
		"on error": {
			&specs.Flow{
				OnError: &specs.OnError{
					Response: &specs.ParameterMap{
						Schema: "com.mock.unkown",
					},
				},
			},
		},
	}

	for name, flows := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := instance.NewContext()
			services := NewMockServices()
			schemas := NewMockSchemas()

			err := Define(ctx, services, schemas, flows)
			if err == nil {
				t.Fatal("unexpected pass")
			}
		})
	}
}

func TestUndefinedNestedSchemaProperty(t *testing.T) {
	t.Parallel()

	tests := map[string]*specs.ParameterMap{
		"single": {
			Schema: "com.mock.message",
			Property: &specs.Property{
				Nested: map[string]*specs.Property{
					"nil": nil,
				},
			},
		},
		"nested": {
			Schema: "com.mock.message",
			Property: &specs.Property{
				Nested: map[string]*specs.Property{
					"meta": {
						Nested: map[string]*specs.Property{
							"nil": nil,
						},
					},
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := instance.NewContext()
			schemas := NewMockSchemas()

			flow := &specs.Flow{
				Name: "mock",
			}

			err := DefineParameterMap(ctx, schemas, test, flow)
			if err == nil {
				t.Fatal("unexpected pass")
			}
		})
	}
}

func TestSettingUndefinedNested(t *testing.T) {
	test := &specs.ParameterMap{
		Schema:   "com.mock.message",
		Property: &specs.Property{},
	}

	ctx := instance.NewContext()
	schemas := NewMockSchemas()

	flow := &specs.Flow{
		Name: "mock",
	}

	err := DefineParameterMap(ctx, schemas, test, flow)
	if err != nil {
		t.Fatal(err)
	}
}
