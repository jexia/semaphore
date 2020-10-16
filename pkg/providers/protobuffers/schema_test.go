package protobuffers

import (
	"testing"

	"github.com/jexia/semaphore/pkg/broker"
	"github.com/jexia/semaphore/pkg/broker/logger"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/types"
)

func equal(t *testing.T, name string, actual, expected interface{}) {
	if actual != expected {
		t.Errorf("property %q %T(%v) was expected to be %T(%v)", name, actual, actual, expected, expected)
	}
}

func notNil(t *testing.T, name string, value interface{}) {
	if value == nil {
		t.Errorf("property %q was not expected to be nil", name)
	}
}

func hasField(t *testing.T, value map[string]*specs.Property, key string) {
	if _, ok := value[key]; !ok {
		t.Errorf("message was expected to contain field %q", key)
	}
}

func hasKeys(t *testing.T, enum *specs.Enum, keys []string) {
	if actual, expected := len(enum.Keys), len(keys); actual != expected {
		t.Errorf("enum %q was expected to have %d keys, got %d", enum.Name, expected, actual)
	}

	for _, key := range keys {
		if _, ok := enum.Keys[key]; !ok {
			t.Errorf("enum %q was expected to contain key %q", enum.Name, key)
		}
	}
}

func schemaFromFile(t *testing.T, path string) (request, response *specs.Property) {
	ctx := logger.WithLogger(broker.NewBackground())

	descriptors, err := Collect(ctx, []string{}, path)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	var (
		ok     bool
		schema = NewSchema(descriptors)
	)

	request, ok = schema["semaphore.Request"]
	if !ok {
		t.Fatal(`"semaphore.RecursiveRequest" was expected to be set`)
	}

	if request.Message == nil {
		t.Fatal("no message in the request")
	}

	response, ok = schema["semaphore.Response"]
	if !ok {
		t.Fatal(`"semaphore.RecursiveResponse" was expected to be set`)
	}

	if response.Message == nil {
		t.Fatal("no message in the response")
	}

	return request, response
}

func TestNewSchema(t *testing.T) {
	type test struct {
		path     string
		request  map[string]func(t *testing.T, property *specs.Property)
		response map[string]func(t *testing.T, property *specs.Property)
	}

	tests := map[string]test{
		"enum": {
			path: "./tests/enum.proto",
			request: map[string]func(t *testing.T, property *specs.Property){
				"action": func(t *testing.T, property *specs.Property) {
					equal(t, "Name", property.Name, "action")
					equal(t, "Path", property.Path, "action")
					equal(t, "Position", property.Position, int32(1))
					notNil(t, "Enum", property.Enum)
					hasKeys(t, property.Enum, []string{"UNKNOWN", "SELECT", "INSERT", "UPDATE", "DELETE"})
				},
				"message": func(t *testing.T, property *specs.Property) {
					equal(t, "Name", property.Name, "message")
					equal(t, "Path", property.Path, "message")
					equal(t, "Position", property.Position, int32(2))
					notNil(t, "Scalar", property.Scalar)
					equal(t, "Scalar.Type", property.Scalar.Type, types.String)
				},
			},
			response: map[string]func(t *testing.T, property *specs.Property){},
		},
		"recursive": {
			path:    "./tests/recursive.proto",
			request: map[string]func(t *testing.T, property *specs.Property){},
			response: map[string]func(t *testing.T, property *specs.Property){
				"string": func(t *testing.T, property *specs.Property) {
					equal(t, "Name", property.Name, "string")
					equal(t, "Path", property.Path, "string")
					equal(t, "Position", property.Position, int32(1))
					notNil(t, "Scalar", property.Scalar)
					equal(t, "Scalar.Type", property.Scalar.Type, types.String)
				},
				"integer": func(t *testing.T, property *specs.Property) {
					equal(t, "Name", property.Name, "integer")
					equal(t, "Path", property.Path, "integer")
					equal(t, "Position", property.Position, int32(2))
					notNil(t, "Scalar", property.Scalar)
					equal(t, "Scalar.Type", property.Scalar.Type, types.Int32)
				},
				"recursive": func(t *testing.T, property *specs.Property) {
					equal(t, "Name", property.Name, "recursive")
					equal(t, "Path", property.Path, "recursive")
					equal(t, "Position", property.Position, int32(3))
					notNil(t, "Message", property.Message)

					hasField(t, property.Message, "boolean")
					hasField(t, property.Message, "recursive")
					// check if it is the same pointer
					equal(t, "Pointer", property, property.Message["recursive"])
				},
			},
		},
		"repeated": {
			path: "./tests/repeated.proto",
			request: map[string]func(t *testing.T, property *specs.Property){
				"limit": func(t *testing.T, property *specs.Property) {
					equal(t, "Name", property.Name, "limit")
					equal(t, "Path", property.Path, "limit")
					equal(t, "Position", property.Position, int32(1))
					notNil(t, "Scalar", property.Scalar)
					equal(t, "Scalar.Type", property.Scalar.Type, types.Int32)
				},
				"offset": func(t *testing.T, property *specs.Property) {
					equal(t, "Name", property.Name, "offset")
					equal(t, "Path", property.Path, "offset")
					equal(t, "Position", property.Position, int32(2))
					notNil(t, "Scalar", property.Scalar)
					equal(t, "Scalar.Type", property.Scalar.Type, types.Int32)
				},
			},
			response: map[string]func(t *testing.T, property *specs.Property){
				"users": func(t *testing.T, property *specs.Property) {
					equal(t, "Name", property.Name, "users")
					equal(t, "Path", property.Path, "users")
					equal(t, "Position", property.Position, int32(1))
					notNil(t, "Message", property.Repeated)

					template, err := property.Repeated.Template()
					if err != nil {
						t.Errorf("unexpected error: %s", err)
					}

					notNil(t, "Template.Message", template.Message)
					hasField(t, template.Message, "id")
					hasField(t, template.Message, "name")
					hasField(t, template.Message, "department")
					hasField(t, template.Message, "subordinates")

					subordinates := template.Message["subordinates"]
					recursive := subordinates.Repeated[0].Message["subordinates"]
					// check if it is the same pointer
					equal(t, "Pointer", subordinates, recursive)
				},
			},
		},
	}

	for title, test := range tests {
		t.Run(title, func(t *testing.T) {
			request, response := schemaFromFile(t, test.path)

			if actual, expected := len(request.Message), len(test.request); actual != expected {
				t.Fatalf("request was expected to contain %d fields, got %d", expected, actual)
			}

			if actual, expected := len(response.Message), len(test.response); actual != expected {
				t.Fatalf("response was expected to contain %d fields, got %d", expected, actual)
			}

			for name, assert := range test.request {
				property, ok := request.Message[name]
				if !ok {
					t.Errorf("request field %q was expected to be set", name)
				}

				assert(t, property)
			}

			for name, assert := range test.response {
				property, ok := response.Message[name]
				if !ok {
					t.Errorf("response field %q was expected to be set", name)
				}

				assert(t, property)
			}
		})
	}
}
