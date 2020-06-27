package proto

import (
	"testing"

	"github.com/jexia/maestro/pkg/specs"
	"github.com/jexia/maestro/pkg/specs/labels"
	"github.com/jexia/maestro/pkg/specs/types"
	"github.com/jhump/protoreflect/desc/builder"
)

type mockMethod struct {
	name     string
	request  map[string]*specs.Property
	response map[string]*specs.Property
}

func (method *mockMethod) GetName() string {
	return method.name
}

func (method *mockMethod) GetRequest() map[string]*specs.Property {
	return method.request
}

func (method *mockMethod) GetResponse() map[string]*specs.Property {
	return method.response
}

func TestMessageDescriptorNonMessage(t *testing.T) {
	path := "mock"
	specs := &specs.ParameterMap{
		Property: &specs.Property{
			Type:    types.String,
			Label:   labels.Optional,
			Default: "hello world",
		},
	}

	_, err := NewMessageDescriptor(path, specs)
	if err == nil {
		t.Fatal("unexpected pass, expected a error to be thrown")
	}
}

func TestMessageDescriptor(t *testing.T) {
	path := "mock"
	tests := map[string]*specs.ParameterMap{
		"simple": {
			Property: &specs.Property{
				Type:     types.Message,
				Label:    labels.Optional,
				Comment:  "",
				Position: 1,
				Nested: map[string]*specs.Property{
					"msg": {
						Type:     types.String,
						Label:    labels.Optional,
						Comment:  "",
						Position: 1,
					},
				},
			},
		},
		"nested": {
			Property: &specs.Property{
				Type:     types.Message,
				Label:    labels.Optional,
				Comment:  "",
				Position: 1,
				Nested: map[string]*specs.Property{
					"msg": {
						Type:     types.Message,
						Label:    labels.Optional,
						Comment:  "",
						Position: 1,
						Nested: map[string]*specs.Property{
							"msg": {
								Type:     types.String,
								Label:    labels.Optional,
								Comment:  "",
								Position: 1,
							},
						},
					},
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			descriptor, err := NewMessageDescriptor(path, test)
			if err != nil {
				t.Fatalf("unexpected error, %s", err)
			}

			if len(descriptor) == 0 {
				t.Fatal("descriptor was not set")
			}
		})
	}
}

func TestServiceDescriptor(t *testing.T) {
	tests := map[string]Methods{
		"simple": {
			"append": &mockMethod{
				name: "append",
				request: map[string]*specs.Property{
					"key": {
						Type:     types.String,
						Label:    labels.Optional,
						Comment:  "",
						Position: 1,
					},
				},
				response: map[string]*specs.Property{
					"key": {
						Type:     types.String,
						Label:    labels.Optional,
						Comment:  "",
						Position: 1,
					},
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			file := builder.NewFile("mock")
			err := NewServiceDescriptor(file, "mock", test)
			if err != nil {
				t.Fatalf("unexpected error, %s", err)
			}
		})
	}
}
