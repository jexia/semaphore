package proto

import (
	"testing"

	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/labels"
	"github.com/jexia/semaphore/pkg/specs/types"
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
						Label:    labels.Required,
						Comment:  "",
						Position: 1,
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
