package proto

import (
	"testing"

	"github.com/jexia/maestro/schema"
	"github.com/jexia/maestro/specs"
	"github.com/jexia/maestro/specs/labels"
	"github.com/jexia/maestro/specs/types"
	"github.com/jhump/protoreflect/desc/builder"
)

type mockSchemaProperty struct {
	position int32
	comment  string
}

func (prop *mockSchemaProperty) GetComment() string {
	return prop.comment
}

func (prop *mockSchemaProperty) GetPosition() int32 {
	return prop.position
}

func (prop *mockSchemaProperty) GetName() string                       { return "" }
func (prop *mockSchemaProperty) GetType() types.Type                   { return "" }
func (prop *mockSchemaProperty) GetLabel() labels.Label                { return "" }
func (prop *mockSchemaProperty) GetNested() map[string]schema.Property { return nil }
func (prop *mockSchemaProperty) GetOptions() schema.Options            { return nil }

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
		"simple": &specs.ParameterMap{
			Property: &specs.Property{
				Type:  types.Message,
				Label: labels.Optional,
				Nested: map[string]*specs.Property{
					"msg": &specs.Property{
						Type:  types.String,
						Label: labels.Optional,
						Desciptor: &mockSchemaProperty{
							comment:  "",
							position: 1,
						},
					},
				},
				Desciptor: &mockSchemaProperty{
					comment:  "",
					position: 1,
				},
			},
		},
		"nested": &specs.ParameterMap{
			Property: &specs.Property{
				Type:  types.Message,
				Label: labels.Optional,
				Nested: map[string]*specs.Property{
					"msg": &specs.Property{
						Type:  types.Message,
						Label: labels.Optional,
						Desciptor: &mockSchemaProperty{
							comment:  "",
							position: 1,
						},
						Nested: map[string]*specs.Property{
							"msg": &specs.Property{
								Type:  types.String,
								Label: labels.Optional,
								Desciptor: &mockSchemaProperty{
									comment:  "",
									position: 1,
								},
							},
						},
					},
				},
				Desciptor: &mockSchemaProperty{
					comment:  "",
					position: 1,
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
		"simple": Methods{
			"append": &mockMethod{
				name: "append",
				request: map[string]*specs.Property{
					"key": &specs.Property{
						Type:  types.String,
						Label: labels.Optional,
						Desciptor: &mockSchemaProperty{
							comment:  "",
							position: 1,
						},
					},
				},
				response: map[string]*specs.Property{
					"key": &specs.Property{
						Type:  types.String,
						Label: labels.Optional,
						Desciptor: &mockSchemaProperty{
							comment:  "",
							position: 1,
						},
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
