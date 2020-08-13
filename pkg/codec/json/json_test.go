package json

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/jexia/semaphore"
	"github.com/jexia/semaphore/pkg/broker"
	"github.com/jexia/semaphore/pkg/broker/logger"
	"github.com/jexia/semaphore/pkg/providers/hcl"
	"github.com/jexia/semaphore/pkg/providers/mock"
	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/template"
)

func NewMock() (specs.FlowListInterface, error) {
	path, err := filepath.Abs("./tests/schema.yaml")
	if err != nil {
		return nil, err
	}

	ctx := logger.WithLogger(broker.NewBackground())

	client, err := semaphore.New(
		ctx,
		semaphore.WithFlows(hcl.FlowsResolver("./tests/*.hcl")),
		semaphore.WithSchema(mock.SchemaResolver(path)),
		semaphore.WithServices(mock.ServicesResolver(path)),
	)

	if err != nil {
		return nil, err
	}

	return client.GetFlows(), nil
}

func ValidateStore(t *testing.T, prop *specs.Property, resource string, origin string, input map[string]interface{}, store references.Store) {
	for key, value := range input {
		nprop := prop.Nested[key]
		if nprop == nil {
			nprop = prop
		}

		path := template.JoinPath(origin, key)
		nested, is := value.(map[string]interface{})
		if is {
			ValidateStore(t, nprop, resource, path, nested, store)
			continue
		}

		repeated, is := value.([]map[string]interface{})
		if is {
			repeating := store.Load(resource, path)
			for index, store := range repeating.Repeated {
				ValidateStore(t, nprop, resource, path, repeated[index], store)
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

				ValidateStore(t, nprop, "", "", wrapper, store)
			}
			continue
		}

		ref := store.Load(resource, path)
		if ref == nil {
			t.Fatalf("resource not found %s", path)
		}

		if ref.Enum != nil && nprop.Enum != nil {
			if nprop.Enum.Positions[*ref.Enum] == nil {
				t.Fatalf("unexpected enum value at %s '%+v', expected '%+v'", path, ref.Enum, value)
			}
			continue
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

	refs := references.NewReferenceStore(len(input))
	refs.StoreValues("input", "", input)

	flows, err := NewMock()
	if err != nil {
		b.Fatal(err)
	}

	flow := flows.Get("simple")
	specs := flow.GetNodes().Get("first").Call.Request

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

	refs := references.NewReferenceStore(len(input))
	refs.StoreValues("input", "", input)

	flows, err := NewMock()
	if err != nil {
		b.Fatal(err)
	}

	flow := flows.Get("nested")
	specs := flow.GetNodes().Get("first").Call.Request

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

	refs := references.NewReferenceStore(len(input))
	refs.StoreValues("input", "", input)

	flows, err := NewMock()
	if err != nil {
		b.Fatal(err)
	}

	flow := flows.Get("repeated")
	specs := flow.GetNodes().Get("first").Call.Request

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

	refs := references.NewReferenceStore(len(input))
	refs.StoreValues("input", "", input)

	flows, err := NewMock()
	if err != nil {
		b.Fatal(err)
	}

	flow := flows.Get("repeated_values")
	specs := flow.GetNodes().Get("first").Call.Request

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

	refs := references.NewReferenceStore(len(input))
	flows, err := NewMock()
	if err != nil {
		b.Fatal(err)
	}

	flow := flows.Get("simple")
	specs := flow.GetNodes().Get("first").Call.Request

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

	refs := references.NewReferenceStore(len(input))
	flows, err := NewMock()
	if err != nil {
		b.Fatal(err)
	}

	flow := flows.Get("nested")
	specs := flow.GetNodes().Get("first").Call.Request

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

	refs := references.NewReferenceStore(len(input))
	flows, err := NewMock()
	if err != nil {
		b.Fatal(err)
	}

	flow := flows.Get("repeated")
	specs := flow.GetNodes().Get("first").Call.Request

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

	refs := references.NewReferenceStore(len(input))
	flows, err := NewMock()
	if err != nil {
		b.Fatal(err)
	}

	flow := flows.Get("repeated_values")
	specs := flow.GetNodes().Get("first").Call.Request

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
	flows, err := NewMock()
	if err != nil {
		t.Fatal(err)
	}

	flow := flows.Get("complete")
	req := flow.GetNodes().Get("first").Call.Request

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
		"enum": {
			"nested": map[string]interface{}{},
			"enum":   "PENDING",
		},
		"repeating_enum": {
			"nested": map[string]interface{}{},
			"repeating_enum": []interface{}{
				"UNKNOWN",
				"PENDING",
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

			refs := references.NewReferenceStore(len(input))
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
	flows, err := NewMock()
	if err != nil {
		t.Fatal(err)
	}

	flow := flows.Get("complete")
	req := flow.GetNodes().Get("first").Call.Request

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
		"enum": {
			"enum": "PENDING",
		},
		"repeating_enum": {
			"repeating_enum": []interface{}{
				"UNKNOWN",
				"PENDING",
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

			store := references.NewReferenceStore(len(input))
			err = manager.Unmarshal(bytes.NewBuffer(inputAsJSON), store)
			if err != nil {
				t.Fatal(err)
			}

			t.Log(store)

			ValidateStore(t, req.Property, "input", "", input, store)
		})
	}
}
