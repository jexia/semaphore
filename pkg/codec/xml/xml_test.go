package xml

import (
	"io/ioutil"
	"testing"

	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs/template"
)

func TestMarshal(t *testing.T) {
	xml := NewConstructor()
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
