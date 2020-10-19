package xml

import (
	"encoding/xml"
	"errors"
	"io"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/jexia/semaphore/pkg/codec/tests"
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

	manager, err := xml.New("mock", tests.SchemaObject)
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
		schema   *specs.ParameterMap
		expected string
	}

	var tests = map[string]test{
		"array empty": {
			schema: tests.SchemaArrayDefaultEmpty,
		},
		"array default reference": {
			input: map[string]interface{}{
				"string": "foo",
			},
			schema:   tests.SchemaArrayWithValues,
			expected: "<array>foo</array><array>bar</array>",
		},
		"simple": {
			input: map[string]interface{}{
				"message": "hello world",
			},
			schema:   tests.SchemaObjectComplex,
			expected: "<root><message>hello world</message><nested></nested></root>",
		},
		"enum": {
			input: map[string]interface{}{
				"nested": map[string]interface{}{},
				"status": references.Enum("PENDING", 1),
			},
			schema:   tests.SchemaObjectComplex,
			expected: "<root><status>PENDING</status><nested></nested></root>",
		},
		"nested": {
			input: map[string]interface{}{
				"nested": map[string]interface{}{
					"first":  "foo",
					"second": "bar",
				},
			},
			schema:   tests.SchemaObjectComplex,
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
			schema:   tests.SchemaObjectComplex,
			expected: "<root><nested></nested><repeating_string>repeating one</repeating_string><repeating_string>repeating two</repeating_string></root>",
		},
		"repeating enum": {
			input: map[string]interface{}{
				"repeating_enum": []interface{}{
					references.Enum("UNKNOWN", 0),
					references.Enum("PENDING", 1),
				},
			},
			schema:   tests.SchemaObjectComplex,
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
			schema:   tests.SchemaObjectComplex,
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
			schema:   tests.SchemaObjectComplex,
			expected: "<root><numeric>42</numeric><message>hello world</message><nested><first>foo</first><second>bar</second></nested><repeating><value>repeating one</value></repeating><repeating><value>repeating two</value></repeating></root>",
		},
	}

	for title, test := range tests {
		t.Run(title, func(t *testing.T) {
			manager, err := constructor.New("mock", test.schema)
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
		expected tests.Expect
		error    error
	}

	testCases := map[string]test{
		"empty scalar with unexpected element": {
			input: strings.NewReader(
				"<integer><unexpected></integer>",
			),
			schema: &specs.ParameterMap{
				Property: tests.PropInteger(),
			},
			error: errFailedToDecode{
				errStack{
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
		},
		"scalar with unexpected element": {
			input: strings.NewReader(
				"<integer>42<unexpected></integer>",
			),
			schema: &specs.ParameterMap{
				Property: tests.PropInteger(),
			},
			error: errFailedToDecode{
				errStack{
					property: "integer",
					inner: errUnexpectedToken{
						actual: xml.StartElement{},
						expected: []xml.Token{
							xml.EndElement{},
						},
					},
				},
			},
		},
		"scalar with type error": {
			input: strings.NewReader(
				"<integer>foo</integer>",
			),
			schema: &specs.ParameterMap{
				Property: tests.PropInteger(),
			},
			error: errFailedToDecode{
				errStack{
					property: "integer",
					inner:    errors.New(`strconv.ParseInt: parsing "foo": invalid syntax`),
				},
			},
		},
		"scalar with empty value": {
			input: strings.NewReader(
				"<integer></integer>",
			),
			schema: &specs.ParameterMap{
				Property: tests.PropInteger(),
			},
			expected: tests.Expect{
				Nested: map[string]tests.Expect{
					"integer": {
						Value: nil,
					},
				},
			},
		},
		"scalar": {
			input: strings.NewReader(
				"<integer>42</integer>",
			),
			schema: &specs.ParameterMap{
				Property: tests.PropInteger(),
			},
			expected: tests.Expect{
				Nested: map[string]tests.Expect{
					"integer": {
						Value: int32(42),
					},
				},
			},
		},
		"empty enum with unexpected element": {
			input: strings.NewReader(
				"<status><unexpected></status>",
			),
			schema: &specs.ParameterMap{
				Property: tests.PropEnum(),
			},
			error: errFailedToDecode{
				errStack{
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
		},
		"enum with unexpected element": {
			input: strings.NewReader(
				"<status>UNKNOWN<unexpected></status>",
			),
			schema: &specs.ParameterMap{
				Property: tests.PropEnum(),
			},
			error: errFailedToDecode{
				errStack{
					property: "status",
					inner: errUnexpectedToken{
						actual: xml.StartElement{},
						expected: []xml.Token{
							xml.EndElement{},
						},
					},
				},
			},
		},
		"enum with unrecognized value": {
			input: strings.NewReader(
				"<status>foo</status>",
			),
			schema: &specs.ParameterMap{
				Property: tests.PropEnum(),
			},
			error: errFailedToDecode{
				errStack{
					property: "status",
					inner:    errUnknownEnum("foo"),
				},
			},
		},
		"enum with empty value": {
			input: strings.NewReader(
				"<status></status>",
			),
			schema: &specs.ParameterMap{
				Property: tests.PropEnum(),
			},
		},
		"enum": {
			input: strings.NewReader(
				"<status>PENDING</status>",
			),
			schema: &specs.ParameterMap{
				Property: tests.PropEnum(),
			},
			expected: tests.Expect{
				Nested: map[string]tests.Expect{
					"status": {
						Enum: func() *int32 { i := int32(1); return &i }(),
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
			schema: tests.SchemaObject,
			expected: tests.Expect{
				Nested: map[string]tests.Expect{
					"root.status": {
						Enum: func() *int32 { i := int32(1); return &i }(),
					},
					"root.integer": {
						Value: int32(42),
					},
				},
			},
		},
		"error on nested property": {
			input: strings.NewReader(
				`<root>
					<nested>
						<status>PENDING</status>
						<integer>oops</integer>
					</nested>
					<string>foobar</string>
				</root>`,
			),
			schema: tests.SchemaObjectNested,
			error: errFailedToDecode{
				errStack: errStack{
					property: "root",
					inner: errFailedToDecode{
						errStack: errStack{
							property: "nested",
							inner: errFailedToDecode{
								errStack: errStack{
									property: "integer",
									inner:    errors.New("strconv.ParseInt: parsing \"oops\": invalid syntax"),
								},
							},
						},
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
			schema: tests.SchemaObjectNested,
			expected: tests.Expect{
				Nested: map[string]tests.Expect{
					"root.nested.status": {
						Enum: func() *int32 { i := int32(1); return &i }(),
					},
					"root.nested.integer": {
						Value: int32(42),
					},
					"root.string": {
						Value: "foobar",
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
			schema: &specs.ParameterMap{
				Property: tests.PropArray(),
			},
			expected: tests.Expect{
				Nested: map[string]tests.Expect{
					"array": {
						Repeated: []tests.Expect{
							{
								Value: "foo",
							},
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
		"array nested": {
			input: strings.NewReader(
				`<root>
					<array>foo</array>
					<array></array>
					<integer>42</integer>
					<array>bar</array>
				</root>`,
			),
			schema: tests.SchemaNestedArray,
			expected: tests.Expect{
				Nested: map[string]tests.Expect{
					"root.integer": {
						Value: int32(42),
					},
					"root.array": {
						Repeated: []tests.Expect{
							{
								Value: "foo",
							},
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
			schema: tests.SchemaObjectComplex,
			expected: tests.Expect{
				Nested: map[string]tests.Expect{
					"root.repeating_string": {
						Repeated: []tests.Expect{
							{
								Value: "foo",
							},
							{
								Value: "bar",
							},
						},
					},
					"root.message": {
						Value: "hello world",
					},
					"root.nested.first": {
						Value: "foo",
					},
					"root.nested.second": {
						Value: "bar",
					},
					"root.repeating": {
						Repeated: []tests.Expect{
							{
								Nested: map[string]tests.Expect{
									"value": {
										Value: "repeating one",
									},
								},
							},
							{
								Nested: map[string]tests.Expect{
									"value": {
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

	for title, test := range testCases {
		t.Run(title, func(t *testing.T) {
			xml := NewConstructor()
			if xml == nil {
				t.Fatal("unexpected nil")
			}

			manager, err := xml.New("mock", test.schema)
			if err != nil {
				t.Fatal(err)
			}

			var store = references.NewReferenceStore(0)
			err = manager.Unmarshal(test.input, store)

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

			tests.Assert(t, "mock", "", store, test.expected)
		})
	}
}
