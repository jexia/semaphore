package proto

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/jexia/semaphore"
	"github.com/jexia/semaphore/cmd/semaphore/daemon/providers"
	"github.com/jexia/semaphore/pkg/broker"
	"github.com/jexia/semaphore/pkg/broker/logger"
	"github.com/jexia/semaphore/pkg/functions"
	"github.com/jexia/semaphore/pkg/providers/hcl"
	"github.com/jexia/semaphore/pkg/providers/protobuffers"
	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/template"
	"github.com/jhump/protoreflect/dynamic"
)

func NewMock() (specs.FlowListInterface, error) {
	ctx := logger.WithLogger(broker.NewBackground())
	core, err := semaphore.NewOptions(ctx,
		semaphore.WithFlows(hcl.FlowsResolver("./tests/*.hcl")),
	)

	if err != nil {
		return nil, err
	}

	options, err := providers.NewOptions(ctx, core,
		providers.WithServices(protobuffers.ServiceResolver([]string{"./tests"}, "./tests/*.proto")),
		providers.WithSchema(protobuffers.SchemaResolver([]string{"./tests"}, "./tests/*.proto")),
	)

	if err != nil {
		return nil, err
	}

	stack := functions.Collection{}
	collection, err := providers.Resolve(ctx, stack, options)
	if err != nil {
		return nil, err
	}

	return collection.FlowListInterface, nil
}

func ValidateStore(t *testing.T, resource, path string, tmpl specs.Template, input interface{}, store references.Store) {
	switch typed := input.(type) {
	case map[string]interface{}:
		for key, value := range typed {
			property := tmpl.Message[key]
			if property == nil {
				t.Fatalf("property (%s) does not exist in map %s:%s", key, resource, path)
			}

			path := template.JoinPath(path, key)

			ValidateStore(t, resource, path, property.Template, value, store)
		}
	case []map[string]interface{}:
		repeating := store.Load(resource, path)
		if repeating == nil {
			t.Fatalf("repeating message does not exist in store '%s:%s'", resource, path)
		}

		tmpl, err := tmpl.Repeated.Template()
		if err != nil {
			t.Fatal(err)
		}

		for index, store := range repeating.Repeated {
			ValidateStore(t, resource, path, tmpl, typed[index], store)
		}
	case []interface{}:
		repeating := store.Load(resource, path)
		if repeating == nil {
			t.Fatalf("resource not found %s:%s", resource, path)
		}

		template, err := tmpl.Repeated.Template()
		if err != nil {
			t.Fatal(err)
		}

		for index, store := range repeating.Repeated {
			ValidateStore(t, "", "", template, typed[index], store)
		}
	case *references.EnumVal:
		ref := store.Load(resource, path)
		if ref == nil {
			t.Fatalf("resource not found %s:%s", resource, path)
		}

		if ref.Enum == nil {
			t.Fatalf("reference enum not set %s:%s", resource, path)
		}

		if *ref.Enum != typed.Pos() {
			t.Fatalf("unexpected enum value at %s:%s '%+v', expected '%+v'", resource, path, ref.Enum, typed.Pos())
		}
	default:
		ref := store.Load(resource, path)
		if ref == nil {
			t.Fatalf("resource not found %s:%s", resource, path)
		}

		if ref.Value != typed {
			t.Fatalf("unexpected value at %s '%+v', expected '%+v'", path, ref.Value, typed)
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

	constructor := NewConstructor()
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

	flow := flows.Get("simple")
	specs := flow.GetNodes().Get("first").Call.Request

	constructor := NewConstructor()
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

func BenchmarkRepeatedMarshal(b *testing.B) {
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

	flow := flows.Get("simple")
	specs := flow.GetNodes().Get("first").Call.Request

	constructor := NewConstructor()
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

	jsonBB, err := json.Marshal(input)
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

	desc, err := NewMessage("MockRequest", specs.Property.Message)
	if err != nil {
		b.Fatal(err)
	}

	msg := dynamic.NewMessage(desc)
	err = msg.UnmarshalJSON(jsonBB)
	if err != nil {
		b.Fatal(err)
	}

	bb, err := msg.Marshal()
	if err != nil {
		b.Fatal(err)
	}

	constructor := NewConstructor()
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

	jsonBB, err := json.Marshal(input)
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

	desc, err := NewMessage("MockRequest", specs.Property.Message)
	if err != nil {
		b.Fatal(err)
	}

	msg := dynamic.NewMessage(desc)
	err = msg.UnmarshalJSON(jsonBB)
	if err != nil {
		b.Fatal(err)
	}

	bb, err := msg.Marshal()
	if err != nil {
		b.Fatal(err)
	}

	constructor := NewConstructor()
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

func BenchmarkRepeatedUnmarshal(b *testing.B) {
	input := map[string]interface{}{
		"repeating": []map[string]interface{}{
			{
				"value": "message",
			},
		},
	}

	jsonBB, err := json.Marshal(input)
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

	desc, err := NewMessage("MockRequest", specs.Property.Message)
	if err != nil {
		b.Fatal(err)
	}

	msg := dynamic.NewMessage(desc)
	err = msg.UnmarshalJSON(jsonBB)
	if err != nil {
		b.Fatal(err)
	}

	bb, err := msg.Marshal()
	if err != nil {
		b.Fatal(err)
	}

	constructor := NewConstructor()
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

func TestUndefinedSpecs(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "return the formatted error",
			want: "no object specs defined",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := ErrUndefinedSpecs{}
			if got := e.Prettify(); got.Message != tt.want {
				t.Errorf("%v, want %v", got, tt.want)
			}
		})
	}
}

func TestNonRootMessage(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "return the formatted error",
			want: "protobuffer messages root property should be a message",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := ErrNonRootMessage{}
			if got := e.Prettify(); got.Message != tt.want {
				t.Errorf("%v, want %v", got, tt.want)
			}
		})
	}
}

func TestMarshal(t *testing.T) {
	flows, err := NewMock()
	if err != nil {
		t.Fatal(err)
	}

	flow := flows.Get("complete")
	req := flow.GetNodes().Get("first").Call.Request
	desc, err := NewMessage("mock", req.Property.Message)
	if err != nil {
		t.Fatal(err)
	}

	response := dynamic.NewMessage(desc)

	constructor := NewConstructor()
	manager, err := constructor.New("input", req)
	if err != nil {
		t.Fatal(err)
	}

	tests := map[string]map[string]interface{}{
		"simple": {
			"message": "hello world",
			"nested":  map[string]interface{}{},
		},
		"nesting": {
			"nested": map[string]interface{}{
				"value": "nested value",
			},
		},
		"repeating": {
			"nested": map[string]interface{}{},
			"repeating": []map[string]interface{}{
				{
					"value": "repeating value",
				},
			},
		},
		"enum": {
			"nested": map[string]interface{}{},
			"status": references.Enum("PENDING", 1),
		},
		"repeating_enum": {
			"nested": map[string]interface{}{},
			"repeating_status": []interface{}{
				references.Enum("PENDING", 1),
				references.Enum("UNKNOWN", 0),
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
					"value": "first repeating value",
				},
				{
					"value": "second repeating value",
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

			store := references.NewReferenceStore(0)
			store.StoreValues("input", "", input)

			reader, err := manager.Marshal(store)
			if err != nil {
				t.Fatal(err)
			}

			bb, err := ioutil.ReadAll(reader)
			if err != nil {
				t.Fatal(err)
			}

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
	flows, err := NewMock()
	if err != nil {
		t.Fatal(err)
	}

	flow := flows.Get("complete")
	req := flow.GetNodes().Get("first").Call.Request

	tests := map[string]map[string]interface{}{
		"simple": {
			"message": "hello world",
			"nested":  map[string]interface{}{},
		},
		"nested": {
			"nested": map[string]interface{}{
				"value": "nested value",
			},
		},
		"repeating": {
			"nested": map[string]interface{}{},
			"repeating": []map[string]interface{}{
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
		"enum": {
			"nested": map[string]interface{}{},
			"status": references.Enum("PENDING", 1),
		},
		"repeating_enum": {
			"nested": map[string]interface{}{},
			"repeating_status": []interface{}{
				references.Enum("PENDING", 1),
				references.Enum("UNKNOWN", 0),
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
			},
		},
	}

	for key, input := range tests {
		t.Run(key, func(t *testing.T) {
			inputAsJSON, err := json.Marshal(input)
			if err != nil {
				t.Fatal(err)
			}

			desc, err := NewMessage("input", req.Property.Message)
			if err != nil {
				t.Fatal(err)
			}

			inputAsProto := dynamic.NewMessage(desc)
			err = inputAsProto.UnmarshalJSON(inputAsJSON)
			if err != nil {
				t.Fatal(err)
			}

			bb, err := inputAsProto.Marshal()
			if err != nil {
				t.Fatal(err)
			}

			store := references.NewReferenceStore(len(input))

			constructor := NewConstructor()
			manager, err := constructor.New("input", req)
			if err != nil {
				t.Fatal(err)
			}

			err = manager.Unmarshal(bytes.NewReader(bb), store)
			if err != nil {
				t.Fatal(err)
			}

			ValidateStore(t, template.InputResource, "", req.Property.Template, input, store)
		})
	}
}
