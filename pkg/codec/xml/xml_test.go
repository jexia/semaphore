package xml

import (
	"errors"
	"io"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/jexia/semaphore/pkg/references"
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

	manager, err := xml.New("mock", schema)
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
	var xml = NewConstructor()
	if xml == nil {
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
			expected: "<mock><message>hello world</message><nested></nested></mock>",
		},
		"enum": {
			input: map[string]interface{}{
				"nested": map[string]interface{}{},
				"status": references.Enum("PENDING", 1),
			},
			expected: "<mock><nested></nested><status>PENDING</status></mock>",
		},
		"nested": {
			input: map[string]interface{}{
				"nested": map[string]interface{}{
					"first":  "foo",
					"second": "bar",
				},
			},
			expected: "<mock><nested><first>foo</first><second>bar</second></nested></mock>",
		},
		"repeating string": {
			input: map[string]interface{}{
				"repeating_string": []interface{}{
					"repeating one",
					"repeating two",
					nil, // TODO: nil (null) values should not be ignored
				},
			},
			expected: "<mock><nested></nested><repeating_string>repeating one</repeating_string><repeating_string>repeating two</repeating_string></mock>",
		},
		"repeating enum": {
			input: map[string]interface{}{
				"repeating_enum": []interface{}{
					"UNKNOWN",
					"PENDING",
				},
			},
			expected: "<mock><nested></nested><repeating_enum>UNKNOWN</repeating_enum><repeating_enum>PENDING</repeating_enum></mock>",
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
			expected: "<mock><nested></nested><repeating><value>repeating one</value></repeating><repeating><value>repeating two</value></repeating></mock>",
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
			expected: "<mock><message>hello world</message><nested><first>foo</first><second>bar</second></nested><repeating><value>repeating one</value></repeating><repeating><value>repeating two</value></repeating></mock>",
		},
	}

	for title, test := range tests {
		t.Run(title, func(t *testing.T) {
			manager, err := xml.New("mock", schema)
			if err != nil {
				t.Fatal(err)
			}

			refs := references.NewReferenceStore(len(test.input))
			refs.StoreValues(template.InputResource, "", test.input)

			r, err := manager.Marshal(refs)
			if err != nil {
				t.Error(err)
			}

			bb, err := ioutil.ReadAll(r)
			if err != nil {
				t.Fatal(err)
			}

			if actual := string(bb); actual != test.expected {
				t.Errorf("output %q was expectetd to be %q", actual, test.expected)
			}
		})
	}
}

type readerFunc func([]byte) (int, error)

func (fn readerFunc) Read(p []byte) (int, error) { return fn(p) }

func TestUnmarshal(t *testing.T) {
	type test struct {
		input    io.Reader
		expected map[string]expect
		error    error
	}

	tests := map[string]test{
		"reader error": {
			input: readerFunc(
				func([]byte) (int, error) {
					return 0, errors.New("failed")
				},
			),
			error: errors.New("failed"),
		},
		"unknown enum value": {
			input: strings.NewReader(
				"<mock><status>PENDING</status><another_status>DONE</another_status></mock>",
			),
			error: errUnknownEnum("DONE"),
		},
		"unexpected token enum value": {
			input: strings.NewReader(
				"<mock><repeating_numeric><oops></oops></repeating_numeric></mock>",
			),
			error: errNotAnObject,
		},
		"unknown enum value (repeated)": {
			input: strings.NewReader(
				"<mock><repeating_enum>DONE</repeating_enum></mock>",
			),
			error: errUnknownEnum("DONE"),
		},
		"type mismatch": {
			input: strings.NewReader(
				"<mock><numeric>not a number</numeric></mock>",
			),
			error: errors.New(""), // error returned by ParseInt()
		},
		"type mismatch (repeated)": {
			input: strings.NewReader(
				"<mock><repeating_numeric>not a number</repeating_numeric></mock>",
			),
			error: errors.New(""), // error returned by ParseInt()
		},
		"empty reader": {
			input: strings.NewReader(""),
		},
		"simple": {
			input: strings.NewReader(
				"<mock><nested></nested><message>hello world</message><another_message>dlrow olleh</another_message></mock>",
			),
			expected: map[string]expect{
				"message": {
					value: "hello world",
				},
				"another_message": {
					value: "dlrow olleh",
				},
			},
		},
		"simple with unknown property": {
			input: strings.NewReader(
				"<mock><unknown></unknown><message>hello world</message><another_message>dlrow olleh</another_message></mock>",
			),
			expected: map[string]expect{
				"message": {
					value: "hello world",
				},
				"another_message": {
					value: "dlrow olleh",
				},
			},
		},
		"enum": {
			input: strings.NewReader(
				"<mock><status>PENDING</status><another_status>UNKNOWN</another_status></mock>",
			),
			expected: map[string]expect{
				"status": {
					enum: func() *int32 { i := int32(1); return &i }(),
				},
				"another_status": {
					enum: func() *int32 { i := int32(0); return &i }(),
				},
			},
		},
		"nested": {
			input: strings.NewReader(
				"<mock><nested><first>foo</first><second>bar</second></nested></mock>",
			),
			expected: map[string]expect{
				"nested.first": {
					value: "foo",
				},
				"nested.second": {
					value: "bar",
				},
			},
		},
		"repeated string": {
			//  TODO: do not ignore empty blocks
			input: strings.NewReader(
				"<mock><repeating_string>repeating one</repeating_string><repeating_string></repeating_string><repeating_string>repeating two</repeating_string></mock>",
			),
			expected: map[string]expect{
				"repeating_string": {
					repeated: []expect{
						{
							value: "repeating one",
						},
						{
							value: "repeating two",
						},
					},
				},
			},
		},
		"repeated enum": {
			input: strings.NewReader(
				"<mock><repeating_enum>UNKNOWN</repeating_enum><repeating_enum>PENDING</repeating_enum></mock>",
			),
			expected: map[string]expect{
				"repeating_enum": {
					repeated: []expect{
						{
							enum: func() *int32 { i := int32(0); return &i }(),
						},
						{
							enum: func() *int32 { i := int32(1); return &i }(),
						},
					},
				},
			},
		},
		"repeated nested": {
			input: strings.NewReader(
				"<mock><repeating><value>repeating one</value></repeating><repeating><value>repeating two</value></repeating></mock>",
			),
			expected: map[string]expect{
				"repeating": {
					repeated: []expect{
						{
							nested: map[string]expect{
								"repeating.value": {
									value: "repeating one",
								},
							},
						},
						{
							nested: map[string]expect{
								"repeating.value": {
									value: "repeating two",
								},
							},
						},
					},
				},
			},
		},
		"complex": {
			input: strings.NewReader(
				`<mock>
				<repeating_string>repeating one</repeating_string>
				<repeating_string>repeating two</repeating_string>
				<message>hello world</message>
					<nested>
						<first>foo</first>
						<second>bar</second>
					</nested><repeating>
				<value>repeating one</value>
				</repeating>
				<repeating>  
					<value>repeating two</value>
				</repeating>
			</mock>`,
			),
			expected: map[string]expect{
				"repeating_string": {
					repeated: []expect{
						{
							value: "repeating one",
						},
						{
							value: "repeating two",
						},
					},
				},
				"message": {
					value: "hello world",
				},
				"nested.first": {
					value: "foo",
				},
				"nested.second": {
					value: "bar",
				},
				"repeating": {
					repeated: []expect{
						{
							nested: map[string]expect{
								"repeating.value": {
									value: "repeating one",
								},
							},
						},
						{
							nested: map[string]expect{
								"repeating.value": {
									value: "repeating two",
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

			manager, err := xml.New("mock", schema)
			if err != nil {
				t.Fatal(err)
			}

			var refs = references.NewReferenceStore(0)
			err = manager.Unmarshal(test.input, refs)

			if test.error != nil {
				if !errors.As(err, &test.error) {
					t.Errorf("error [%s] was expected to be [%s]", err, test.error)
				}
			} else if err != nil {
				t.Errorf("error was not expected: %s", err)
			}

			for path, output := range test.expected {
				assert(t, "mock", path, refs, output)
			}
		})
	}
}

type expect struct {
	value    interface{}
	enum     *int32
	repeated []expect
	nested   map[string]expect
}

func assert(t *testing.T, resource string, path string, store references.Store, output expect) {
	var ref = store.Load(resource, path)

	if ref == nil {
		t.Errorf("reference %q was expected to be set", path)

		return
	}

	if output.value != nil {
		if ref.Value != output.value {
			t.Errorf("reference %q [%v] was expected to have value [%v]", path, output.value, ref.Value)
		}

		return
	}

	if output.enum != nil {
		if ref.Enum == nil {
			t.Errorf("reference %q was expected to have a enum value", path)
		}

		if *output.enum != *ref.Enum {
			t.Errorf("reference %q was expected to have enum value [%d], not [%d]", path, *output.enum, *ref.Enum)
		}

		return
	}

	if output.repeated != nil {
		if ref.Repeated == nil {
			t.Errorf("reference %q was expected to have a repeated value", path)

			return
		}

		if expected, actual := len(ref.Repeated), len(ref.Repeated); actual != expected {
			t.Errorf("invalid number of repeated values, expected %d, got %d", expected, actual)

			return
		}

		for index, expected := range output.repeated {
			if expected.value != nil || expected.enum != nil {
				assert(t, "", "", ref.Repeated[index], expected)

				continue
			}

			if expected.nested != nil {
				for key, expected := range expected.nested {
					assert(t, resource, key, ref.Repeated[index], expected)
				}

				continue
			}
		}
	}
}
