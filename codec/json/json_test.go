package json

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/jexia/maestro"
	"github.com/jexia/maestro/definitions/hcl"
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

func FindNode(flow *specs.Flow, name string) *specs.Node {
	for _, node := range flow.GetNodes() {
		if node.GetName() == name {
			return node
		}
	}

	return nil
}

func NewMock() (*specs.Manifest, error) {
	path, err := filepath.Abs("./tests/schema.yaml")
	if err != nil {
		return nil, err
	}

	client, err := maestro.New(
		maestro.WithDefinitions(hcl.DefinitionResolver("./tests/*.hcl")),
		maestro.WithSchema(mock.SchemaResolver(path)),
	)

	if err != nil {
		return nil, err
	}

	return client.Manifest, nil
}

func ValidateStore(t *testing.T, resource string, origin string, input map[string]interface{}, store *specs.Store) {
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

		values, is := value.([]interface{})
		if is {
			repeating := store.Load(resource, path)
			for index, store := range repeating.Repeated {
				// small wrapper that allows to reuse functionalities
				wrapper := map[string]interface{}{
					"": values[index],
				}

				ValidateStore(t, "", "", wrapper, store)
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

	refs := specs.NewReferenceStore(len(input))
	refs.StoreValues("input", "", input)

	manifest, err := NewMock()
	if err != nil {
		b.Fatal(err)
	}

	flow := FindFlow(manifest, "simple")
	specs := FindNode(flow, "first").Call.GetRequest()

	constructor := &Constructor{}
	manager, err := constructor.New("input", specs)
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

	refs := specs.NewReferenceStore(len(input))
	refs.StoreValues("input", "", input)

	manifest, err := NewMock()
	if err != nil {
		b.Fatal(err)
	}

	flow := FindFlow(manifest, "nested")
	specs := FindNode(flow, "first").Call.GetRequest()

	constructor := &Constructor{}
	manager, err := constructor.New("input", specs)
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

func BenchmarkRepeatedMessagesMarshal(b *testing.B) {
	input := map[string]interface{}{
		"repeating": []map[string]interface{}{
			{
				"value": "message",
			},
		},
	}

	refs := specs.NewReferenceStore(len(input))
	refs.StoreValues("input", "", input)

	manifest, err := NewMock()
	if err != nil {
		b.Fatal(err)
	}

	flow := FindFlow(manifest, "repeated")
	specs := FindNode(flow, "first").Call.GetRequest()

	constructor := &Constructor{}
	manager, err := constructor.New("input", specs)
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

func BenchmarkRepeatedValuesMarshal(b *testing.B) {
	input := map[string]interface{}{
		"repeating_values": []interface{}{
			"message",
		},
	}

	refs := specs.NewReferenceStore(len(input))
	refs.StoreValues("input", "", input)

	manifest, err := NewMock()
	if err != nil {
		b.Fatal(err)
	}

	flow := FindFlow(manifest, "repeated_values")
	specs := FindNode(flow, "first").Call.GetRequest()

	constructor := &Constructor{}
	manager, err := constructor.New("input", specs)
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

	refs := specs.NewReferenceStore(len(input))
	manifest, err := NewMock()
	if err != nil {
		b.Fatal(err)
	}

	flow := FindFlow(manifest, "simple")
	specs := FindNode(flow, "first").Call.GetRequest()

	constructor := &Constructor{}
	manager, err := constructor.New("input", specs)
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

	refs := specs.NewReferenceStore(len(input))
	manifest, err := NewMock()
	if err != nil {
		b.Fatal(err)
	}

	flow := FindFlow(manifest, "nested")
	specs := FindNode(flow, "first").Call.GetRequest()

	constructor := &Constructor{}
	manager, err := constructor.New("input", specs)
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

func BenchmarkRepeatedMessagesUnmarshal(b *testing.B) {
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

	refs := specs.NewReferenceStore(len(input))
	manifest, err := NewMock()
	if err != nil {
		b.Fatal(err)
	}

	flow := FindFlow(manifest, "repeated")
	specs := FindNode(flow, "first").Call.GetRequest()

	constructor := &Constructor{}
	manager, err := constructor.New("input", specs)
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

func BenchmarkRepeatedValuesUnmarshal(b *testing.B) {
	input := map[string]interface{}{
		"repeating_values": []interface{}{
			"message",
		},
	}

	bb, err := json.Marshal(input)
	if err != nil {
		b.Fatal(err)
	}

	refs := specs.NewReferenceStore(len(input))
	manifest, err := NewMock()
	if err != nil {
		b.Fatal(err)
	}

	flow := FindFlow(manifest, "repeated_values")
	specs := FindNode(flow, "first").Call.GetRequest()

	constructor := &Constructor{}
	manager, err := constructor.New("input", specs)
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
	req := FindNode(flow, "first").Call.GetRequest()

	constructor := &Constructor{}
	manager, err := constructor.New("input", req)
	if err != nil {
		t.Fatal(err)
	}

	tests := map[string]map[string]interface{}{
		"simple": {
			"message": "some message",
			"nested":  map[string]interface{}{},
		},
		"nested": {
			"nested": map[string]interface{}{
				"value": "some message",
			},
		},
		"repeating": {
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
		"repeating_values": {
			"nested": map[string]interface{}{},
			"repeating_values": []interface{}{
				"repeating value",
				"repeating value",
			},
		},
		"complex": {
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

			refs := specs.NewReferenceStore(len(input))
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

func TestSimple(t *testing.T) {
	_, err := NewMock()
	if err != nil {
		t.Fatal(err)
	}
}

func TestUnmarshal(t *testing.T) {
	manifest, err := NewMock()
	if err != nil {
		t.Fatal(err)
	}

	flow := FindFlow(manifest, "complete")
	req := FindNode(flow, "first").Call.GetRequest()

	constructor := &Constructor{}
	manager, err := constructor.New("input", req)
	if err != nil {
		t.Fatal(err)
	}

	tests := map[string]map[string]interface{}{
		"simple": {
			"message": "some message",
			"nested":  map[string]interface{}{},
		},
		"nested": {
			"nested": map[string]interface{}{
				"value": "some message",
			},
		},
		"repeating": {
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
		"repeating_values": {
			"nested": map[string]interface{}{},
			"repeating_values": []interface{}{
				"repeating value",
				"repeating value",
			},
		},
		"complex": {
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

			store := specs.NewReferenceStore(len(input))
			err = manager.Unmarshal(bytes.NewBuffer(inputAsJSON), store)
			if err != nil {
				t.Fatal(err)
			}

			t.Log(store)

			ValidateStore(t, "input", "", input, store)
		})
	}
}
