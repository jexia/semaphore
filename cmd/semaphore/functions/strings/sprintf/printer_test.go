package sprintf

import (
	"testing"

	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/template"
)

func TestPrinter(t *testing.T) {
	type test struct {
		format   string
		args     []*specs.Property
		store    map[string]interface{}
		expected string
	}

	var tests = map[string]test{
		"": {
			format: "note that %f policeman + %f policeman != %d policemen",
			store: map[string]interface{}{
				"message": "hello world",
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
				printer = NewPrinter(tokens)
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
