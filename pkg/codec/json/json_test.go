package json

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/jexia/semaphore"
	"github.com/jexia/semaphore/cmd/semaphore/daemon/providers"
	"github.com/jexia/semaphore/pkg/broker"
	"github.com/jexia/semaphore/pkg/broker/logger"
	"github.com/jexia/semaphore/pkg/codec/tests"
	"github.com/jexia/semaphore/pkg/functions"
	"github.com/jexia/semaphore/pkg/providers/hcl"
	"github.com/jexia/semaphore/pkg/providers/mock"
	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
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

		if _, err := ioutil.ReadAll(reader); err != nil {
			b.Fatal(err)
		}
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

		if _, err := ioutil.ReadAll(reader); err != nil {
			b.Fatal(err)
		}
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
					var property = tests.PropInteger()
					property.Reference = &specs.PropertyReference{
						Resource: "input",
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
					var property = tests.PropInteger()
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
			input:    map[string]interface{}{},
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
			expected: `{"message":"some message","nested":{},"repeating":[],"repeating_values":[],"repeating_enum":[]}`,
		},
		"nested": {
			input: map[string]interface{}{
				"nested": map[string]interface{}{
					"value": "some message",
				},
			},
			schema:   schema,
			expected: `{"nested":{"value":"some message"},"repeating":[],"repeating_values":[],"repeating_enum":[]}`,
		},
		"enum": {
			input: map[string]interface{}{
				"nested": map[string]interface{}{},
				"enum":   references.Enum("PENDING", 2),
			},
			schema:   schema,
			expected: `{"nested":{},"repeating":[],"repeating_values":[],"enum":"PENDING","repeating_enum":[]}`,
		},
		"repeating_enum": {
			input: map[string]interface{}{
				"repeating_enum": []interface{}{
					references.Enum("UNKNOWN", 1),
					references.Enum("PENDING", 2),
				},
			},
			schema:   schema,
			expected: `{"nested":{},"repeating":[],"repeating_values":[],"repeating_enum":["UNKNOWN","PENDING"]}`,
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
			expected: `{"nested":{},"repeating":[{"value":"repeating value"},{"value":"repeating value"}],"repeating_values":[],"repeating_enum":[]}`,
		},
		"repeating values from reference": {
			input: map[string]interface{}{
				"repeating_values": []interface{}{
					"repeating one",
					"repeating two",
				},
			},
			schema:   schema,
			expected: `{"nested":{},"repeating":[],"repeating_values":["repeating one","repeating two"],"repeating_enum":[]}`,
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
			expected: `{"message":"hello world","nested":{"value":"nested value"},"repeating":[{"value":"repeating value"},{"value":"repeating value"}],"repeating_values":[],"repeating_enum":[]}`,
		},
	}

	for key, test := range tests {
		t.Run(key, func(t *testing.T) {
			constructor := &Constructor{}
			manager, err := constructor.New("input", test.schema)
			if err != nil {
				t.Fatal(err)
			}

			refs := references.NewReferenceStore(len(test.input))
			refs.StoreValues("input", "", test.input)

			reader, err := manager.Marshal(refs)
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
		expected tests.Expect
	}

	testsCases := map[string]test{
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
			expected: tests.Expect{
				Nested: map[string]tests.Expect{
					"array": {
						Repeated: []tests.Expect{
							{
								Value: nil,
							},
							{
								Value: "bar",
							},
						},
					},
				},
			},
		},
		"array of arrays": {
			input:  `[[null,"bar"]]`,
			schema: tests.SchemaArrayOfArrays,
			expected: tests.Expect{
				Nested: map[string]tests.Expect{
					"array": {
						Repeated: []tests.Expect{
							{
								Repeated: []tests.Expect{
									{
										Value: nil,
									},
									{
										Value: "bar",
									},
								},
							},
						},
					},
				},
			},
		},
		"simple": {
			input:  `{"message":"some message"}`,
			schema: schema,
			expected: tests.Expect{
				Nested: map[string]tests.Expect{
					"message": {
						Value: "some message",
					},
				},
			},
		},
		"nested": {
			input:  `{"nested":{"value":"some message"}}`,
			schema: schema,
			expected: tests.Expect{
				Nested: map[string]tests.Expect{
					"nested.value": {
						Value: "some message",
					},
				},
			},
		},
		"enum": {
			input:  `{"enum":"PENDING"}`,
			schema: schema,
			expected: tests.Expect{
				Nested: map[string]tests.Expect{
					"enum": {
						Enum: func() *int32 { i := int32(2); return &i }(),
					},
				},
			},
		},
		"repeating_enum": {
			input:  `{"repeating_enum":["UNKNOWN","PENDING"]}`,
			schema: schema,
			expected: tests.Expect{
				Nested: map[string]tests.Expect{
					"repeating_enum": {
						Repeated: []tests.Expect{
							{
								Enum: func() *int32 { i := int32(1); return &i }(),
							},
							{
								Enum: func() *int32 { i := int32(2); return &i }(),
							},
						},
					},
				},
			},
		},
		"repeating_values": {
			input:  `{"repeating_values":["repeating one","repeating two"]}`,
			schema: schema,
			expected: tests.Expect{
				Nested: map[string]tests.Expect{
					"repeating_values": {
						Repeated: []tests.Expect{
							{
								Value: "repeating one",
							},
							{
								Value: "repeating two",
							},
						},
					},
				},
			},
		},
		"repeating objects": {
			input:  `{"repeating":[{"value":"repeating one"},{"value":"repeating two"}]}`,
			schema: schema,
			expected: tests.Expect{
				Nested: map[string]tests.Expect{
					"repeating": {
						Repeated: []tests.Expect{
							{
								Nested: map[string]tests.Expect{
									"repeating.value": {
										Value: "repeating one",
									},
								},
							},
							{
								Nested: map[string]tests.Expect{
									"repeating.value": {
										Value: "repeating two",
									},
								},
							},
						},
					},
				},
			},
		},
		"complex": {
			input:  `{"message":"hello world","nested":{"value":"hello nested world"},"repeating":[{"value":"repeating one"},{"value":"repeating two"}]}`,
			schema: schema,
			expected: tests.Expect{
				Nested: map[string]tests.Expect{
					"message": {
						Value: "hello world",
					},
					"nested": {
						Nested: map[string]tests.Expect{
							"value": {
								Value: "hello nested world",
							},
						},
					},
					"repeating": {
						Repeated: []tests.Expect{
							{
								Nested: map[string]tests.Expect{
									"repeating.value": {
										Value: "repeating one",
									},
								},
							},
							{
								Nested: map[string]tests.Expect{
									"repeating.value": {
										Value: "repeating two",
									},
								},
							},
						},
					},
				},
			},
		},
	}

	for key, test := range testsCases {
		t.Run(key, func(t *testing.T) {
			constructor := &Constructor{}
			manager, err := constructor.New("input", test.schema)
			if err != nil {
				t.Fatal(err)
			}

			store := references.NewReferenceStore(0)
			err = manager.Unmarshal(bytes.NewBuffer([]byte(test.input)), store)

			if err != nil {
				t.Fatal(err)
			}

			tests.Assert(t, "input", "", store, test.expected)
		})
	}
}
