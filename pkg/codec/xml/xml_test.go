package xml

import (
	"encoding/xml"
	"errors"
	"io"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/template"
)

func TestName(t *testing.T) {
	var (
		xml      = NewConstructor()
		expected = "xml"
	)

	if xml == nil {
		t.Fatal("unexpected nil")
	}

	t.Run("check constuctor name", func(t *testing.T) {
		if actual := xml.Name(); actual != expected {
			t.Errorf("constructor name %q was expected to be %s", actual, expected)
		}
	})

	manager, err := xml.New("mock", SchemaObject)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("check manager name", func(t *testing.T) {
		if actual := manager.Name(); actual != expected {
			t.Errorf("manager name %q was expected to be %s", actual, expected)
		}
	})
}

func TestMarshal(t *testing.T) {
	var constructor = NewConstructor()
	if constructor == nil {
		t.Fatal("unexpected nil")
	}

	type test struct {
		input    map[string]interface{}
		expected string
	}

	var tests = map[string]test{
		"simple": {
			input: map[string]interface{}{
				"message": "hello world",
			},
			expected: "<root><message>hello world</message><nested></nested></root>",
		},
		"enum": {
			input: map[string]interface{}{
				"nested": map[string]interface{}{},
				"status": references.Enum("PENDING", 1),
			},
			expected: "<root><status>PENDING</status><nested></nested></root>",
		},
		"nested": {
			input: map[string]interface{}{
				"nested": map[string]interface{}{
					"first":  "foo",
					"second": "bar",
				},
			},
			expected: "<root><nested><first>foo</first><second>bar</second></nested></root>",
		},
		"repeating string": {
			input: map[string]interface{}{
				"repeating_string": []interface{}{
					"repeating one",
					"repeating two",
					nil, // TODO: nil (null) values should not be ignored
				},
			},
			expected: "<root><nested></nested><repeating_string>repeating one</repeating_string><repeating_string>repeating two</repeating_string></root>",
		},
		"repeating enum": {
			input: map[string]interface{}{
				"repeating_enum": []interface{}{
					references.Enum("UNKNOWN", 0),
					references.Enum("PENDING", 1),
				},
			},
			expected: "<root><nested></nested><repeating_enum>UNKNOWN</repeating_enum><repeating_enum>PENDING</repeating_enum></root>",
		},
		"repeating nested": {
			input: map[string]interface{}{
				"repeating": []map[string]interface{}{
					{
						"value": "repeating one",
					},
					{
						"value": "repeating two",
					},
				},
			},
			expected: "<root><nested></nested><repeating><value>repeating one</value></repeating><repeating><value>repeating two</value></repeating></root>",
		},
		"complex": {
			input: map[string]interface{}{
				"message": "hello world",
				"nested": map[string]interface{}{
					"first":  "foo",
					"second": "bar",
				},
				"numeric": 42,
				"repeating": []map[string]interface{}{
					{
						"value": "repeating one",
					},
					{
						"value": "repeating two",
					},
				},
			},
			expected: "<root><numeric>42</numeric><message>hello world</message><nested><first>foo</first><second>bar</second></nested><repeating><value>repeating one</value></repeating><repeating><value>repeating two</value></repeating></root>",
		},
	}

	for title, test := range tests {
		t.Run(title, func(t *testing.T) {
			manager, err := constructor.New("mock", SchemaObjectComplex)
			if err != nil {
				t.Fatal(err)
			}

			refs := references.NewReferenceStore(len(test.input))
			refs.StoreValues(template.InputResource, "", test.input)

			reader, err := manager.Marshal(refs)
			if err != nil {
				t.Error(err)
			}

			bb, err := ioutil.ReadAll(reader)
			if err != nil {
				t.Fatal(err)
			}

			if actual := string(bb); actual != test.expected {
				t.Errorf("unexpected difference %s, %s", test.expected, actual)
			}
		})
	}
}

type readerFunc func([]byte) (int, error)

func (fn readerFunc) Read(p []byte) (int, error) { return fn(p) }

func errorString(err error) string {
	if err != nil {
		return err.Error()
	}

	return "<nil>"
}

func TestUnmarshal(t *testing.T) {
	type test struct {
		input    io.Reader
		schema   *specs.ParameterMap
		expected expect
		error    error
	}

	tests := map[string]test{
		"empty scalar with unexpected element": {
			input: strings.NewReader(
				"<integer><unexpected></integer>",
			),
			schema: SchemaScalar,
			error: errFailedToDecodeProperty{
				property: "integer",
				inner: errUnexpectedToken{
					actual: xml.StartElement{},
					expected: []xml.Token{
						xml.CharData{},
						xml.EndElement{},
					},
				},
			},
		},
		"scalar with unexpected element": {
			input: strings.NewReader(
				"<integer>42<unexpected></integer>",
			),
			schema: SchemaScalar,
			error: errFailedToDecodeProperty{
				property: "integer",
				inner: errUnexpectedToken{
					actual: xml.StartElement{},
					expected: []xml.Token{
						xml.EndElement{},
					},
				},
			},
		},
		"scalar with type error": {
			input: strings.NewReader(
				"<integer>foo</integer>",
			),
			schema: SchemaScalar,
			error: errFailedToDecodeProperty{
				property: "integer",
				inner:    errors.New(`strconv.ParseInt: parsing "foo": invalid syntax`),
			},
		},
		"scalar with empty value": {
			input: strings.NewReader(
				"<integer></integer>",
			),
			schema: SchemaScalar,
			expected: expect{
				nested: map[string]expect{
					"integer": {
						value: nil,
					},
				},
			},
		},
		"scalar": {
			input: strings.NewReader(
				"<integer>42</integer>",
			),
			schema: SchemaScalar,
			expected: expect{
				nested: map[string]expect{
					"integer": {
						value: int32(42),
					},
				},
			},
		},
		"empty enum with unexpected element": {
			input: strings.NewReader(
				"<status><unexpected></status>",
			),
			schema: SchemaEnum,
			error: errFailedToDecodeProperty{
				property: "status",
				inner: errUnexpectedToken{
					actual: xml.StartElement{},
					expected: []xml.Token{
						xml.CharData{},
						xml.EndElement{},
					},
				},
			},
		},
		"enum with unexpected element": {
			input: strings.NewReader(
				"<status>UNKNOWN<unexpected></status>",
			),
			schema: SchemaEnum,
			error: errFailedToDecodeProperty{
				property: "status",
				inner: errUnexpectedToken{
					actual: xml.StartElement{},
					expected: []xml.Token{
						xml.EndElement{},
					},
				},
			},
		},
		"enum with unrecognized value": {
			input: strings.NewReader(
				"<status>foo</status>",
			),
			schema: SchemaEnum,
			error: errFailedToDecodeProperty{
				property: "status",
				inner:    errUnknownEnum("foo"),
			},
		},
		"enum with empty value": {
			input: strings.NewReader(
				"<status></status>",
			),
			schema: SchemaEnum,
			expected: expect{
				nested: map[string]expect{
					"status": {
						value: nil,
					},
				},
			},
		},
		"enum": {
			input: strings.NewReader(
				"<status>PENDING</status>",
			),
			schema: SchemaEnum,
			expected: expect{
				nested: map[string]expect{
					"status": {
						enum: func() *int32 { i := int32(1); return &i }(),
					},
				},
			},
		},
		"object": {
			input: strings.NewReader(
				`<root>
					<status>PENDING</status>
					<integer>42</integer>
				</root>`,
			),
			schema: SchemaObject,
			expected: expect{
				nested: map[string]expect{
					"root.status": {
						enum: func() *int32 { i := int32(1); return &i }(),
					},
					"root.integer": {
						value: int32(42),
					},
				},
			},
		},
		"object nested": {
			input: strings.NewReader(
				`<root>
					<nested>
						<status>PENDING</status>
						<integer>42</integer>
					</nested>
					<string>foobar</string>
				</root>`,
			),
			schema: SchemaObjectNested,
			expected: expect{
				nested: map[string]expect{
					"root.nested.status": {
						enum: func() *int32 { i := int32(1); return &i }(),
					},
					"root.nested.integer": {
						value: int32(42),
					},
					"root.string": {
						value: "foobar",
					},
				},
			},
		},
		"array of strings": {
			input: strings.NewReader(
				`<array>foo</array>
				<array></array>
				<array>bar</array>`,
			),
			schema: SchemaArray,
			expected: expect{
				nested: map[string]expect{
					"array": {
						repeated: []expect{
							{
								value: "foo",
							},
							{
								value: nil,
							},
							{
								value: "bar",
							},
						},
					},
				},
			},
		},
		"array nested": {
			input: strings.NewReader(
				`<root>
					<array>foo</array>
					<array></array>
					<integer>42</integer>
					<array>bar</array>
				</root>`,
			),
			schema: SchemaNestedArray,
			expected: expect{
				nested: map[string]expect{
					"root.integer": {
						value: int32(42),
					},
					"root.array": {
						repeated: []expect{
							{
								value: "foo",
							},
							{
								value: nil,
							},
							{
								value: "bar",
							},
						},
					},
				},
			},
		},
		"complex": {
			input: strings.NewReader(
				`<root>
				<numeric>42</numeric>
				<message>hello world</message>
				<nested>
					<first>foo</first>
					<second>bar</second>
				</nested>
				<repeating_string>foo</repeating_string>
				<repeating_string>bar</repeating_string>
				<repeating>
					<value>repeating one</value>
				</repeating>
				<repeating>
					<value>repeating two</value>
				</repeating>
			</root>`,
			),
			schema: SchemaObjectComplex,
			expected: expect{
				nested: map[string]expect{
					"root.repeating_string": {
						repeated: []expect{
							{
								value: "foo",
							},
							{
								value: "bar",
							},
						},
					},
					"root.message": {
						value: "hello world",
					},
					"root.nested.first": {
						value: "foo",
					},
					"root.nested.second": {
						value: "bar",
					},
					"root.repeating": {
						repeated: []expect{
							{
								nested: map[string]expect{
									"value": {
										value: "repeating one",
									},
								},
							},
							{
								nested: map[string]expect{
									"value": {
										value: "repeating two",
									},
								},
							},
						},
					},
				},
			},
		},
	}

	for title, test := range tests {
		t.Run(title, func(t *testing.T) {
			xml := NewConstructor()
			if xml == nil {
				t.Fatal("unexpected nil")
			}

			manager, err := xml.New("mock", test.schema)
			if err != nil {
				t.Fatal(err)
			}

			var refs = references.NewReferenceStore(0)
			err = manager.Unmarshal(test.input, refs)

			if test.error != nil {
				if err == nil {
					t.Fatalf("error '%s' was expected", test.error)
				}

				// TODO: find a better way to compare errors
				if err.Error() != test.error.Error() {
					t.Fatalf("error '%s' was expected to be '%s'", errorString(err), test.error)
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error '%s'", err)
				}
			}

			assert(t, "mock", "", refs, test.expected)
		})
	}
}

type expect struct {
	value    interface{}
	enum     *int32
	repeated []expect
	nested   map[string]expect
}

func assert(t *testing.T, resource, path string, store references.Store, input expect) {
	if input.nested != nil {
		for key, value := range input.nested {
			assert(t, resource, key, store, value)
		}

		return
	}

	ref := store.Load(resource, path)

	switch {
	case input.enum != nil:
		if ref == nil || ref.Enum == nil {
			t.Fatalf("reference %q was expected to have a enum value", path)
		}

		if *input.enum != *ref.Enum {
			t.Errorf("reference %q was expected to have enum value [%d], not [%d]", path, *input.enum, *ref.Enum)
		}
	case input.value != nil:
		if ref == nil || ref.Value != input.value {
			t.Errorf("reference %q was expected to be %T(%v), got %T(%v)", path, input.value, input.value, ref.Value, ref.Value)
		}
	case input.repeated != nil:
		if ref == nil || ref.Repeated == nil {
			t.Fatalf("reference %q was expected to have a repeated value", path)
		}

		if expected, actual := len(input.repeated), len(ref.Repeated); actual != expected {
			t.Fatalf("invalid number of repeated values, expected %d, got %d", expected, actual)
		}

		for index, expected := range input.repeated {
			assert(t, "", "", ref.Repeated[index], expected)
		}
	}
}
