package json

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/jexia/semaphore/v2"
	"github.com/jexia/semaphore/v2/cmd/semaphore/daemon/providers"
	"github.com/jexia/semaphore/v2/pkg/broker"
	"github.com/jexia/semaphore/v2/pkg/broker/logger"
	"github.com/jexia/semaphore/v2/pkg/codec/tests"
	"github.com/jexia/semaphore/v2/pkg/functions"
	"github.com/jexia/semaphore/v2/pkg/providers/hcl"
	"github.com/jexia/semaphore/v2/pkg/providers/mock"
	"github.com/jexia/semaphore/v2/pkg/references"
	"github.com/jexia/semaphore/v2/pkg/specs"
)

func NewMock() (specs.FlowListInterface, error) {
	path, err := filepath.Abs("./tests/schema.yaml")
	if err != nil {
		return nil, err
	}

	ctx := logger.WithLogger(broker.NewBackground())
	core, err := semaphore.NewOptions(ctx,
		semaphore.WithFlows(hcl.FlowsResolver("./tests/*.hcl")),
	)
	if err != nil {
		return nil, err
	}

	options, err := providers.NewOptions(ctx, core,
		providers.WithSchema(mock.SchemaResolver(path)),
		providers.WithServices(mock.ServicesResolver(path)),
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
	tracker := references.NewTracker()
	references.StoreValues(refs, tracker, "input:", input)

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

		io.Copy(ioutil.Discard, reader)
	}
}

func BenchmarkNestedMarshal(b *testing.B) {
	input := map[string]interface{}{
		"nested": map[string]interface{}{
			"value": "message",
		},
	}

	refs := references.NewStore(len(input))
	tracker := references.NewTracker()
	references.StoreValues(refs, tracker, "input:", input)

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

		io.Copy(ioutil.Discard, reader)
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

	refs := references.NewStore(len(input))
	tracker := references.NewTracker()
	references.StoreValues(refs, tracker, "input:", input)

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

		if _, err := ioutil.ReadAll(reader); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkRepeatedValuesMarshal(b *testing.B) {
	input := map[string]interface{}{
		"repeating_values": []interface{}{
			"message",
		},
	}

	refs := references.NewStore(len(input))
	tracker := references.NewTracker()
	references.StoreValues(refs, tracker, "input:", input)

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

		if _, err := ioutil.ReadAll(reader); err != nil {
			b.Fatal(err)
		}
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

	refs := references.NewStore(len(input))
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

	refs := references.NewStore(len(input))
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

	refs := references.NewStore(len(input))
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

	refs := references.NewStore(len(input))
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

func TestMarshal(t *testing.T) {
	flows, err := NewMock()
	if err != nil {
		t.Fatal(err)
	}

	flow := flows.Get("complete")
	schema := flow.GetNodes().Get("first").Call.Request

	type test struct {
		input    map[string]interface{}
		schema   *specs.ParameterMap
		expected string
	}

	tests := map[string]test{
		"nil schema": {
			schema: new(specs.ParameterMap),
		},
		"scalar from reference": {
			input: map[string]interface{}{
				"integer": int32(42),
			},
			schema: &specs.ParameterMap{
				Property: func() *specs.Property {
					property := tests.PropInteger()
					property.Reference = &specs.PropertyReference{
						Resource: specs.InputResource,
						Path:     "integer",
					}

					return property
				}(),
			},
			expected: `42`,
		},
		"scalar default value": {
			input: map[string]interface{}{
				"integer": int32(42),
			},
			schema: &specs.ParameterMap{
				Property: func() *specs.Property {
					property := tests.PropInteger()
					property.Scalar.Default = int32(42)

					return property
				}(),
			},
			expected: `42`,
		},
		"array empty": {
			schema:   tests.SchemaArrayDefaultEmpty,
			expected: `[null]`,
		},
		"array default reference": {
			input: map[string]interface{}{
				"string": "foo",
			},
			schema:   tests.SchemaArrayWithValues,
			expected: `["foo","bar"]`,
		},
		"array default no reference value": {
			input: map[string]interface{}{
				"string": nil,
			},
			schema:   tests.SchemaArrayWithValues,
			expected: `[null,"bar"]`,
		},
		"array of arrays": {
			input:    map[string]interface{}{},
			schema:   tests.SchemaArrayOfArrays,
			expected: `[[null,"bar"]]`,
		},
		"simple": {
			input: map[string]interface{}{
				"message": "some message",
				"nested":  map[string]interface{}{},
			},
			schema:   schema,
			expected: `{"message":"some message","nested":{}}`,
		},
		"nested": {
			input: map[string]interface{}{
				"nested": map[string]interface{}{
					"value": "some message",
				},
			},
			schema:   schema,
			expected: `{"nested":{"value":"some message"}}`,
		},
		"enum": {
			input: map[string]interface{}{
				"nested": map[string]interface{}{},
				"enum":   references.Enum("PENDING", 2),
			},
			schema:   schema,
			expected: `{"nested":{},"enum":"PENDING"}`,
		},
		"repeating_enum": {
			input: map[string]interface{}{
				"repeating_enum": []interface{}{
					references.Enum("UNKNOWN", 1),
					references.Enum("PENDING", 2),
				},
			},
			schema:   schema,
			expected: `{"nested":{},"repeating_enum":["UNKNOWN","PENDING"]}`,
		},
		"repeating objects": {
			input: map[string]interface{}{
				"repeating": []map[string]interface{}{
					{
						"value": "repeating value",
					},
					{
						"value": "repeating value",
					},
				},
			},
			schema:   schema,
			expected: `{"nested":{},"repeating":[{"value":"repeating value"},{"value":"repeating value"}]}`,
		},
		"repeating values from reference": {
			input: map[string]interface{}{
				"repeating_values": []interface{}{
					"repeating one",
					"repeating two",
				},
			},
			schema:   schema,
			expected: `{"nested":{},"repeating_values":["repeating one","repeating two"]}`,
		},
		"complex": {
			input: map[string]interface{}{
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
			schema:   schema,
			expected: `{"message":"hello world","nested":{"value":"nested value"},"repeating":[{"value":"repeating value"},{"value":"repeating value"}]}`,
		},
	}

	for key, test := range tests {
		t.Run(key, func(t *testing.T) {
			constructor := &Constructor{}
			manager, err := constructor.New(specs.InputResource, test.schema)
			if err != nil {
				t.Fatal(err)
			}

			store := references.NewStore(len(test.input))
			tracker := references.NewTracker()
			references.StoreValues(store, tracker, specs.ResourcePath(specs.InputResource), test.input)

			reader, err := manager.Marshal(store)
			if err != nil {
				t.Fatal(err)
			}

			if test.expected != "" {
				data, err := ioutil.ReadAll(reader)
				if err != nil {
					t.Fatal(err)
				}

				if actual := string(data); actual != test.expected {
					t.Errorf("unexpected output %s, expected %s", data, test.expected)
				}
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

	var (
		flow   = flows.Get("complete")
		schema = flow.GetNodes().Get("first").Call.Request
	)

	type test struct {
		input    string
		schema   *specs.ParameterMap
		expected map[string]tests.Expect
	}

	cases := map[string]test{
		"nil schema": {
			schema: new(specs.ParameterMap),
		},
		"empty": {
			input:  ``,
			schema: schema,
		},
		"array": {
			input:  `[null,"bar"]`,
			schema: tests.SchemaArrayDefaultEmpty,
			expected: map[string]tests.Expect{
				"[0]": {
					Scalar: nil,
				},
				"[1]": {
					Scalar: "bar",
				},
			},
		},
		"array of arrays": {
			input:  `[[null,"bar"]]`,
			schema: tests.SchemaArrayOfArrays,
			expected: map[string]tests.Expect{
				"[0][0]": {
					Scalar: nil,
				},
				"[0][1]": {
					Scalar: "bar",
				},
			},
		},
		"simple": {
			input:  `{"message":"some message"}`,
			schema: schema,
			expected: map[string]tests.Expect{
				"message": {
					Scalar: "some message",
				},
			},
		},
		"nested": {
			input:  `{"nested":{"value":"some message"}}`,
			schema: schema,
			expected: map[string]tests.Expect{
				"nested.value": {
					Scalar: "some message",
				},
			},
		},
		"enum": {
			input:  `{"enum":"PENDING"}`,
			schema: schema,
			expected: map[string]tests.Expect{
				"enum": {
					Enum: func() *int32 { i := int32(2); return &i }(),
				},
			},
		},
		"repeating_enum": {
			input:  `{"repeating_enum":["UNKNOWN","PENDING"]}`,
			schema: schema,
			expected: map[string]tests.Expect{
				"repeating_enum[0]": {
					Enum: func() *int32 { i := int32(1); return &i }(),
				},
				"repeating_enum[1]": {
					Enum: func() *int32 { i := int32(2); return &i }(),
				},
			},
		},
		"repeating_values": {
			input:  `{"repeating_values":["repeating one","repeating two"]}`,
			schema: schema,
			expected: map[string]tests.Expect{
				"repeating_values[0]": {
					Scalar: "repeating one",
				},
				"repeating_values[1]": {
					Scalar: "repeating two",
				},
			},
		},
		"repeating objects": {
			input:  `{"repeating":[{"value":"repeating one"},{"value":"repeating two"}]}`,
			schema: schema,
			expected: map[string]tests.Expect{
				"repeating[0].value": {
					Scalar: "repeating one",
				},
				"repeating[1].value": {
					Scalar: "repeating two",
				},
			},
		},
		"complex": {
			input:  `{"message":"hello world","nested":{"value":"hello nested world"},"repeating":[{"value":"repeating one"},{"value":"repeating two"}]}`,
			schema: schema,
			expected: map[string]tests.Expect{
				"message": {
					Scalar: "hello world",
				},
				"nested.value": {
					Scalar: "hello nested world",
				},

				"repeating[0].value": {
					Scalar: "repeating one",
				},

				"repeating[1].value": {
					Scalar: "repeating two",
				},
			},
		},
	}

	for key, test := range cases {
		t.Run(key, func(t *testing.T) {
			constructor := &Constructor{}
			manager, err := constructor.New(specs.InputResource, test.schema)
			if err != nil {
				t.Fatal(err)
			}

			store := references.NewStore(0)
			err = manager.Unmarshal(bytes.NewBuffer([]byte(test.input)), store)

			if err != nil {
				t.Fatal(err)
			}

			for path, expect := range test.expected {
				tests.Assert(t, specs.InputResource, path, store, expect)
			}
		})
	}
}
