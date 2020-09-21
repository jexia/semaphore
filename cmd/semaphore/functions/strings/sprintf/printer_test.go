package sprintf

import (
	"reflect"
	"testing"

	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/labels"
	"github.com/jexia/semaphore/pkg/specs/template"
	"github.com/jexia/semaphore/pkg/specs/types"
)

func TestTokensPrint(t *testing.T) {
	type test struct {
		format   string
		args     []*specs.Property
		store    map[string]interface{}
		expected string
	}

	var tests = map[string]test{
		"test with string arguments": {
			format:   "note that a %s of the policeman + another %s of the policeman != the %s policeman",
			expected: "note that a half of the policeman + another half of the policeman != the whole policeman",
			store: map[string]interface{}{
				"half": "half",
			},
			args: []*specs.Property{
				{
					Name: "half",
					Path: "half",
					Reference: &specs.PropertyReference{
						Resource: template.InputResource,
						Path:     "half",
					},
				},
				{
					Name: "half",
					Path: "half",
					Reference: &specs.PropertyReference{
						Resource: template.InputResource,
						Path:     "half",
					},
				},
				{
					Name:    "whole",
					Path:    "whole",
					Default: "whole",
				},
			},
		},
		"test with numeric arguments": {
			format:   "note that %.1f policeman + %.1f policeman != %d policemen",
			expected: "note that 1.5 policeman + 0.5 policeman != 2 policemen",
			store: map[string]interface{}{
				"second": 0.5,
			},
			args: []*specs.Property{
				{
					Name:    "first",
					Path:    "first",
					Default: 1.5,
				},
				{
					Name: "second",
					Path: "second",
					Reference: &specs.PropertyReference{
						Resource: template.InputResource,
						Path:     "second",
					},
				},
				{
					Name:    "message",
					Path:    "message",
					Default: 2,
				},
			},
		},
		"test with json formatter": {
			format:   `{"array":%json,"object":%json}`,
			expected: `{"array":[3.14159265,6.62607015],"object":{"e":2.71828182}}`,
			store: map[string]interface{}{
				"object": map[string]interface{}{
					"e": float64(2.71828182),
				},
				"array": []interface{}{
					float64(3.14159265),
					float64(6.62607015),
				},
			},
			args: []*specs.Property{
				{
					Name:  "array",
					Path:  "array",
					Type:  types.Float,
					Label: labels.Repeated,
					Reference: &specs.PropertyReference{
						Resource: template.InputResource,
						Path:     "array",
					},
				},
				{
					Name: "object",
					Path: "object",
					Type: types.Message,
					Nested: map[string]*specs.Property{
						"e": {
							Name:  "e",
							Path:  "object.e",
							Type:  types.Float,
							Label: labels.Optional,
							Reference: &specs.PropertyReference{
								Resource: template.InputResource,
								Path:     "object.e",
							},
						},
					},
					Reference: &specs.PropertyReference{
						Resource: template.InputResource,
						Path:     "object",
					},
				},
			},
		},
	}

	for title, test := range tests {
		t.Run(title, func(t *testing.T) {
			tokens, err := scanner.Scan(test.format)
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			var (
				printer = Tokens(tokens)
				refs    = references.NewReferenceStore(len(test.store))
			)

			refs.StoreValues(template.InputResource, "", test.store)

			actual, err := printer.Print(refs, test.args...)
			if err != nil {
				t.Fatalf("unexpected error %q", err)
			}

			if actual != test.expected {
				t.Errorf("output %q does not match expected %q", actual, test.expected)
			}
		})
	}
}

func TestTokensVerbs(t *testing.T) {
	var (
		tokens   = Tokens{Constant("a"), Verb{Verb: "b"}, Precision{}, Constant("c"), Verb{Verb: "d"}}
		expected = []Verb{{Verb: "b"}, {Verb: "d"}}
		actual   = tokens.Verbs()
	)

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("the output '%v' does not match expected '%v'", actual, expected)
	}
}
