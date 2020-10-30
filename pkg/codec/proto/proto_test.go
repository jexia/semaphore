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
	"github.com/jexia/semaphore/pkg/codec/tests"
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

func BenchmarkSimpleMarshal(b *testing.B) {
	input := map[string]interface{}{
		"message": "message",
	}

	refs := references.NewStore(len(input))
	refs.Store("input:", &references.Reference{Value: input})

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

	refs := references.NewStore(len(input))
	refs.Store("input:", &references.Reference{Value: input})

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

	refs := references.NewStore(len(input))
	refs.Store("input:", &references.Reference{Value: input})

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

	refs := references.NewStore(len(input))
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

	refs := references.NewStore(len(input))
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

	refs := references.NewStore(len(input))
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
		},
		"nesting": {
			"nested": map[string]interface{}{
				"value": "nested value",
			},
		},
		"repeating": {
			"repeating": []map[string]interface{}{
				{
					"value": "repeating value",
				},
			},
		},
		"enum": {
			"status": references.Enum("PENDING", 1),
		},
		"repeating_enum": {
			"repeating_status": []interface{}{
				references.Enum("PENDING", 1),
				references.Enum("UNKNOWN", 0),
			},
		},
		"repeating_values": {
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

			store := references.NewStore(len(input))
			references.StoreValues(store, references.NewTracker(), "input:", input)

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

	type test struct {
		schema   map[string]interface{}
		expected map[string]tests.Expect
	}

	cases := map[string]test{
		"simple": {
			schema: map[string]interface{}{
				"message": "hello world",
				"nested":  map[string]interface{}{},
			},
			expected: map[string]tests.Expect{
				"message": {Scalar: "hello world"},
			},
		},
		"nested": {
			schema: map[string]interface{}{
				"nested": map[string]interface{}{
					"value": "nested value",
				},
			},
			expected: map[string]tests.Expect{
				"nested.value": {Scalar: "nested value"},
			},
		},
		"repeating": {
			schema: map[string]interface{}{
				"nested": map[string]interface{}{},
				"repeating": []map[string]interface{}{
					{
						"value": "repeating value",
					},
				},
			},
			expected: map[string]tests.Expect{
				"repeating[0].value": {Scalar: "repeating value"},
			},
		},
		"repeating_values": {
			schema: map[string]interface{}{
				"nested": map[string]interface{}{},
				"repeating_values": []interface{}{
					"repeating value",
					"repeating value",
				},
			},
			expected: map[string]tests.Expect{
				"repeating_values[0]": {Scalar: "repeating value"},
				"repeating_values[1]": {Scalar: "repeating value"},
			},
		},
		"enum": {
			schema: map[string]interface{}{
				"nested": map[string]interface{}{},
				"status": references.Enum("PENDING", 1),
			},
			expected: map[string]tests.Expect{
				"status": {Enum: func() *int32 { i := int32(1); return &i }()},
			},
		},
		"repeating_enum": {
			schema: map[string]interface{}{
				"nested": map[string]interface{}{},
				"repeating_status": []interface{}{
					references.Enum("PENDING", 1),
					references.Enum("UNKNOWN", 0),
				},
			},
			expected: map[string]tests.Expect{
				"repeating_status[0]": {Enum: func() *int32 { i := int32(1); return &i }()},
				"repeating_status[1]": {Enum: func() *int32 { i := int32(0); return &i }()},
			},
		},
		"complex": {
			schema: map[string]interface{}{
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
			expected: map[string]tests.Expect{
				"message":            {Scalar: "hello world"},
				"nested.value":       {Scalar: "nested value"},
				"repeating[0].value": {Scalar: "repeating value"},
			},
		},
	}

	for key, test := range cases {
		t.Run(key, func(t *testing.T) {
			inputAsJSON, err := json.Marshal(test.schema)
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

			store := references.NewStore(len(test.schema))

			constructor := NewConstructor()
			manager, err := constructor.New(template.InputResource, req)
			if err != nil {
				t.Fatal(err)
			}

			err = manager.Unmarshal(bytes.NewReader(bb), store)
			if err != nil {
				t.Fatal(err)
			}

			for path, expect := range test.expected {
				tests.Assert(t, template.InputResource, path, store, expect)
			}
		})
	}
}
