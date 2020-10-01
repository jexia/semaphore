package proto

import (
	"testing"

	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/labels"
	"github.com/jexia/semaphore/pkg/specs/types"
)

type mockMethod struct {
	name     string
	request  specs.Message
	response specs.Message
}

func (method *mockMethod) GetName() string {
	return method.name
}

func (method *mockMethod) GetRequest() specs.Message {
	return method.request
}

func (method *mockMethod) GetResponse() specs.Message {
	return method.response
}

func TestServiceDescriptor(t *testing.T) {
	tests := map[string]Methods{
		"simple": {
			"append": &mockMethod{
				name: "append",
				request: specs.Message{
					"key": {
						Name:        "key",
						Path:        "key",
						Label:       labels.Optional,
						Description: "",
						Position:    1,
						Template: specs.Template{
							Scalar: &specs.Scalar{
								Type: types.String,
							},
						},
					},
				},
				response: specs.Message{
					"key": {
						Name:        "key",
						Path:        "key",
						Label:       labels.Optional,
						Description: "",
						Position:    1,
						Template: specs.Template{
							Scalar: &specs.Scalar{
								Type: types.String,
							},
						},
					},
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			var service = &Service{
				Package: "pkg",
				Name:    "test",
				Methods: test,
			}

			_, err := service.FileDescriptor()
			if err != nil {
				t.Fatalf("unexpected error, %s", err)
			}
		})
	}
}
