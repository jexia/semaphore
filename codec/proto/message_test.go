package proto

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/jexia/maestro"
	"github.com/jexia/maestro/refs"
	"github.com/jexia/maestro/schema/protoc"
	"github.com/jexia/maestro/specs"
	"github.com/jhump/protoreflect/dynamic"
)

func NewMock(t *testing.T) (protoc.Object, specs.Object) {
	collection, err := protoc.Collect(nil, "./tests")
	if err != nil {
		t.Fatal(err)
	}

	manifest, err := maestro.New(maestro.WithPath("./tests", false), maestro.WithSchemaCollection(collection))
	if err != nil {
		t.Fatal(err)
	}

	method := collection.GetService("proto.Logger").GetMethod("Append")
	schema := method.GetInput().(protoc.Object)
	specs := manifest.Flows[0].GetCalls()[0].Request

	return schema, specs
}

func ValidateStore(t *testing.T, resource string, origin string, input map[string]interface{}, store *refs.Store) {
	for key, value := range input {
		path := specs.JoinPath(origin, key)
		nested, is := value.(map[string]interface{})
		if is {
			ValidateStore(t, resource, path, nested, store)
			continue
		}

		repeated, is := value.([]map[string]interface{})
		if is {
			repeating := store.Load(resource, path)
			for index, store := range repeating.Repeated {
				ValidateStore(t, resource, path, repeated[index], store)
			}
			continue
		}

		ref := store.Load(resource, path)
		if ref == nil {
			t.Fatalf("resource not found %s", path)
		}

		if ref.Value != value {
			t.Fatalf("unexpected value at %s '%+v', expected '%+v'", path, ref.Value, value)
		}
	}
}

func TestMarshal(t *testing.T) {
	schema, specs := NewMock(t)
	tests := map[string]map[string]interface{}{
		"simple": map[string]interface{}{
			"message": "hello world",
			"nested":  map[string]interface{}{},
		},
		"nested": map[string]interface{}{
			"nested": map[string]interface{}{
				"value": "nested value",
			},
		},
		"complex": map[string]interface{}{
			"message": "hello world",
			"nested": map[string]interface{}{
				"value": "nested value",
			},
			"repeating": []map[string]interface{}{
				{
					"value": "repeating value",
				},
			},
		},
	}

	for key, input := range tests {
		t.Run(key, func(t *testing.T) {
			inputAsJSON, err := json.Marshal(input)
			if err != nil {
				t.Fatal(err)
			}

			store := refs.NewStore(3)
			store.StoreValues("input", "", input)

			manager, err := New("input", schema, specs)
			if err != nil {
				t.Fatal(err)
			}

			reader, err := manager.Marshal(store)
			if err != nil {
				t.Fatal(err)
			}

			bb, err := ioutil.ReadAll(reader)
			if err != nil {
				t.Fatal(err)
			}

			response := dynamic.NewMessage(schema.GetDescriptor())
			err = response.Unmarshal(bb)
			if err != nil {
				t.Fatal(err)
			}

			responseAsJSON, err := response.MarshalJSON()
			if err != nil {
				t.Fatal(err)
			}

			result := map[string]interface{}{}
			err = json.Unmarshal(responseAsJSON, &result)
			if err != nil {
				t.Fatal(err)
			}

			expected := map[string]interface{}{}
			err = json.Unmarshal(inputAsJSON, &expected)
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(expected, result) {
				t.Errorf("unexpected response %s, expected %s", string(responseAsJSON), string(inputAsJSON))
			}
		})
	}
}

func TestUnmarshal(t *testing.T) {
	schema, specs := NewMock(t)
	tests := map[string]map[string]interface{}{
		"simple": map[string]interface{}{
			"message": "hello world",
			"nested":  map[string]interface{}{},
		},
		"nested": map[string]interface{}{
			"nested": map[string]interface{}{
				"value": "nested value",
			},
		},
		"complex": map[string]interface{}{
			"message": "hello world",
			"nested": map[string]interface{}{
				"value": "nested value",
			},
			"repeating": []map[string]interface{}{
				{
					"value": "repeating value",
				},
			},
		},
	}

	for key, input := range tests {
		t.Run(key, func(t *testing.T) {
			inputAsJSON, err := json.Marshal(input)
			if err != nil {
				t.Fatal(err)
			}

			inputAsProto := dynamic.NewMessage(schema.GetDescriptor())
			err = inputAsProto.UnmarshalJSON(inputAsJSON)
			if err != nil {
				t.Fatal(err)
			}

			bb, _ := inputAsProto.Marshal()
			store := refs.NewStore(3)

			manager, err := New("input", schema, specs)
			if err != nil {
				t.Fatal(err)
			}

			err = manager.Unmarshal(bytes.NewBuffer(bb), store)
			if err != nil {
				t.Fatal(err)
			}

			t.Log(store)

			ValidateStore(t, "input", "", input, store)
		})
	}
}
