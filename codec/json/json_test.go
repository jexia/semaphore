package json

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/jexia/maestro"
	"github.com/jexia/maestro/refs"
	"github.com/jexia/maestro/schema/mock"
	"github.com/jexia/maestro/specs"
)

func FindFlow(manifest *specs.Manifest, name string) *specs.Flow {
	for _, flow := range manifest.Flows {
		if flow.GetName() == name {
			return flow
		}
	}

	return nil
}

func FindCall(flow *specs.Flow, name string) *specs.Call {
	for _, call := range flow.GetCalls() {
		if call.GetName() == name {
			return call
		}
	}

	return nil
}

func NewMock() (*specs.Manifest, error) {
	path, err := filepath.Abs("./tests/schema.yaml")
	if err != nil {
		return nil, err
	}

	reader, err := os.Open(path)
	collection, err := mock.UnmarshalFile(reader)
	if err != nil {
		return nil, err
	}

	manifest, err := maestro.New(maestro.WithPath("./tests", false), maestro.WithSchemaCollection(collection))
	if err != nil {
		return nil, err
	}

	return manifest, nil
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

func BenchmarkSimpleMarshal(b *testing.B) {
	input := map[string]interface{}{
		"message": "message",
	}

	refs := refs.NewStore(len(input))
	refs.StoreValues("input", "", input)

	manifest, err := NewMock()
	if err != nil {
		b.Fatal(err)
	}

	flow := FindFlow(manifest, "simple")
	specs := FindCall(flow, "first").GetRequest()

	manager, err := New("input", nil, specs)
	if err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		reader, err := manager.Marshal(refs)
		if err != nil {
			b.Fatal(err)
		}

		ioutil.ReadAll(reader)
	}
}

func BenchmarkNestedMarshal(b *testing.B) {
	input := map[string]interface{}{
		"nested": map[string]interface{}{
			"value": "message",
		},
	}

	refs := refs.NewStore(len(input))
	refs.StoreValues("input", "", input)

	manifest, err := NewMock()
	if err != nil {
		b.Fatal(err)
	}

	flow := FindFlow(manifest, "nested")
	specs := FindCall(flow, "first").GetRequest()

	manager, err := New("input", nil, specs)
	if err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		reader, err := manager.Marshal(refs)
		if err != nil {
			b.Fatal(err)
		}

		ioutil.ReadAll(reader)
	}
}

func BenchmarkRepeatedMarshal(b *testing.B) {
	input := map[string]interface{}{
		"repeating": []map[string]interface{}{
			{
				"value": "message",
			},
		},
	}

	refs := refs.NewStore(len(input))
	refs.StoreValues("input", "", input)

	manifest, err := NewMock()
	if err != nil {
		b.Fatal(err)
	}

	flow := FindFlow(manifest, "repeated")
	specs := FindCall(flow, "first").GetRequest()

	manager, err := New("input", nil, specs)
	if err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		reader, err := manager.Marshal(refs)
		if err != nil {
			b.Fatal(err)
		}

		ioutil.ReadAll(reader)
	}
}

func BenchmarkSimpleUnmarshal(b *testing.B) {
	input := map[string]interface{}{
		"message": "message",
	}

	bb, err := json.Marshal(input)
	if err != nil {
		b.Fatal(err)
	}

	refs := refs.NewStore(len(input))
	manifest, err := NewMock()
	if err != nil {
		b.Fatal(err)
	}

	flow := FindFlow(manifest, "simple")
	specs := FindCall(flow, "first").GetRequest()

	manager, err := New("input", nil, specs)
	if err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := manager.Unmarshal(bytes.NewBuffer(bb), refs)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkNestedUnmarshal(b *testing.B) {
	input := map[string]interface{}{
		"nested": map[string]interface{}{
			"value": "message",
		},
	}

	bb, err := json.Marshal(input)
	if err != nil {
		b.Fatal(err)
	}

	refs := refs.NewStore(len(input))
	manifest, err := NewMock()
	if err != nil {
		b.Fatal(err)
	}

	flow := FindFlow(manifest, "nested")
	specs := FindCall(flow, "first").GetRequest()

	manager, err := New("input", nil, specs)
	if err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := manager.Unmarshal(bytes.NewBuffer(bb), refs)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkRepeatedUnmarshal(b *testing.B) {
	input := map[string]interface{}{
		"repeating": []map[string]interface{}{
			{
				"value": "message",
			},
		},
	}

	bb, err := json.Marshal(input)
	if err != nil {
		b.Fatal(err)
	}

	refs := refs.NewStore(len(input))
	manifest, err := NewMock()
	if err != nil {
		b.Fatal(err)
	}

	flow := FindFlow(manifest, "repeated")
	specs := FindCall(flow, "first").GetRequest()

	manager, err := New("input", nil, specs)
	if err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := manager.Unmarshal(bytes.NewBuffer(bb), refs)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func TestMarshal(t *testing.T) {
	manifest, err := NewMock()
	if err != nil {
		t.Fatal(err)
	}

	flow := FindFlow(manifest, "complete")
	specs := FindCall(flow, "first").GetRequest()

	manager, err := New("input", nil, specs)
	if err != nil {
		t.Fatal(err)
	}

	tests := map[string]map[string]interface{}{
		"simple": map[string]interface{}{
			"message": "some message",
			"nested":  map[string]interface{}{},
		},
		"nested": map[string]interface{}{
			"nested": map[string]interface{}{
				"value": "some message",
			},
		},
		"repeating": map[string]interface{}{
			"nested": map[string]interface{}{},
			"repeating": []map[string]interface{}{
				{
					"value": "repeating value",
				},
				{
					"value": "repeating value",
				},
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

			refs := refs.NewStore(len(input))
			refs.StoreValues("input", "", input)

			reader, err := manager.Marshal(refs)
			if err != nil {
				t.Fatal(err)
			}

			responseAsJSON, err := ioutil.ReadAll(reader)
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
	manifest, err := NewMock()
	if err != nil {
		t.Fatal(err)
	}

	flow := FindFlow(manifest, "complete")
	specs := FindCall(flow, "first").GetRequest()

	manager, err := New("input", nil, specs)
	if err != nil {
		t.Fatal(err)
	}

	tests := map[string]map[string]interface{}{
		"simple": map[string]interface{}{
			"message": "some message",
			"nested":  map[string]interface{}{},
		},
		"nested": map[string]interface{}{
			"nested": map[string]interface{}{
				"value": "some message",
			},
		},
		"repeating": map[string]interface{}{
			"nested": map[string]interface{}{},
			"repeating": []map[string]interface{}{
				{
					"value": "repeating value",
				},
				{
					"value": "repeating value",
				},
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

			store := refs.NewStore(len(input))
			err = manager.Unmarshal(bytes.NewBuffer(inputAsJSON), store)
			if err != nil {
				t.Fatal(err)
			}

			t.Log(store)

			ValidateStore(t, "input", "", input, store)
		})
	}
}
