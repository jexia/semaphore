package providers

import (
	"testing"

	"github.com/jexia/semaphore/pkg/broker"
	"github.com/jexia/semaphore/pkg/broker/logger"
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
					Name:   "unknown",
					Input:  "com.mock.unknown",
					Output: "com.mock.unknown",
				},
			},
		},
	}
}

func NewMockSchemas() specs.Schemas {
	return specs.Schemas{
		"com.mock.message": &specs.Property{
			Label: labels.Optional,
			Template: specs.Template{
				Message: specs.Message{
					"value": {
						Name:  "value",
						Path:  "value",
						Label: labels.Optional,
						Template: specs.Template{
							Scalar: &specs.Scalar{
								Type: types.String,
							},
						},
					},
					"meta": {
						Name:  "meta",
						Path:  "meta",
						Label: labels.Optional,
						Template: specs.Template{
							Message: specs.Message{
								"id": {
									Name:  "id",
									Path:  "meta.id",
									Label: labels.Optional,
									Template: specs.Template{
										Scalar: &specs.Scalar{
											Type: types.String,
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
			ctx := logger.WithLogger(broker.NewBackground())
			services := NewMockServices()
			schemas := NewMockSchemas()

			err := ResolveSchemas(ctx, services, schemas, flows)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestDefineSchemasUnknown(t *testing.T) {
	t.Parallel()

	tests := map[string]specs.FlowListInterface{
		"input": {
			&specs.Flow{
				Input: &specs.ParameterMap{
					Schema: "com.mock.unknown",
				},
			},
		},
		"output": {
			&specs.Flow{
				Output: &specs.ParameterMap{
					Schema: "com.mock.unknown",
				},
			},
		},
		"node call request": {
			&specs.Flow{
				Nodes: specs.NodeList{
					&specs.Node{
						Call: &specs.Call{
							Request: &specs.ParameterMap{
								Schema: "com.mock.unknown",
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
								Schema: "com.mock.unknown",
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
							Method:   "unknown",
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
								Schema: "com.mock.unknown",
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
								Schema: "com.mock.unknown",
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
								Schema: "com.mock.unknown",
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
								Schema: "com.mock.unknown",
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
						Schema: "com.mock.unknown",
					},
				},
			},
		},
	}

	for name, flows := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := logger.WithLogger(broker.NewBackground())
			services := NewMockServices()
			schemas := NewMockSchemas()

			err := ResolveSchemas(ctx, services, schemas, flows)
			if err == nil {
				t.Fatal("unexpected pass")
			}
		})
	}
}

func TestErrUndefinedObject(t *testing.T) {
	type fields struct {
		Schema string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"return the formatted error",
			fields{Schema: "get"},
			"object 'get', is unavailable inside the schema collection",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := ErrUndefinedObject{
				Schema: "get",
			}
			if got := e.Prettify(); got.Message != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrUndefinedService(t *testing.T) {
	type fields struct {
		Service string
		Flow    string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"return the formatted error",
			fields{Service: "getsources", Flow: "add"},
			"undefined service 'getsources' in flow 'add'",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := ErrUndefinedService{
				Service: "getsources",
				Flow:    "add",
			}
			if got := e.Prettify(); got.Message != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrUndefinedMethod(t *testing.T) {
	type fields struct {
		Method string
		Flow   string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"return the formatted error",
			fields{Method: "get", Flow: "add"},
			"undefined method 'get' in flow 'add'",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := ErrUndefinedMethod{
				Method: "get",
				Flow:   "add",
			}
			if got := e.Prettify(); got.Message != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrUndefinedOutput(t *testing.T) {
	type fields struct {
		Output string
		Flow   string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"return the formatted error",
			fields{Output: "json", Flow: "add"},
			"undefined method output property 'json' in flow 'add'",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := ErrUndefinedOutput{
				Output: "json",
				Flow:   "add",
			}
			if got := e.Prettify(); got.Message != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrUndefinedProperty(t *testing.T) {
	type fields struct {
		Property string
		Flow     string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"return the formatted error",
			fields{Property: "getdata", Flow: "add"},
			"undefined schema nested message property 'getdata' in flow 'add'",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := ErrUndefinedProperty{
				Property: "getdata",
				Flow:     "add",
			}
			if got := e.Prettify(); got.Message != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
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
				Template: specs.Template{
					Message: specs.Message{
						"meta": nil,
					},
				},
			},
		},
		"nested": {
			Schema: "com.mock.message",
			Property: &specs.Property{
				Template: specs.Template{
					Message: specs.Message{
						"meta": {
							Name: "meta",
							Path: "meta",
							Template: specs.Template{
								Message: specs.Message{
									"id": nil,
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
			ctx := logger.WithLogger(broker.NewBackground())
			schemas := NewMockSchemas()

			flow := &specs.Flow{
				Name: "mock",
			}

			err := ResolveParameterMap(ctx, schemas, test, flow)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestSettingUndefinedNested(t *testing.T) {
	test := &specs.ParameterMap{
		Schema:   "com.mock.message",
		Property: &specs.Property{},
	}

	ctx := logger.WithLogger(broker.NewBackground())
	schemas := NewMockSchemas()

	flow := &specs.Flow{
		Name: "mock",
	}

	err := ResolveParameterMap(ctx, schemas, test, flow)
	if err != nil {
		t.Fatal(err)
	}
}
