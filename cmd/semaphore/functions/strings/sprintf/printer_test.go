package sprintf

import (
	"log"
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
			format: "note that %.1f policeman + %.1f policeman != %d policemen", //"%s%s%s",
			store: map[string]interface{}{
				"message": "hello world",
			},
			args: []*specs.Property{
				{
					Name:    "message",
					Path:    "input.message",
					Default: 1.5,
				},
				{
					Name:      "message",
					Path:      "message",
					Reference: nil, // TODO
				},
				{
					Name:    "message",
					Path:    "message",
					Default: 42,
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
				printer = NewPrinter(tokens)
				refs    = references.NewReferenceStore(len(test.store))
			)

			refs.StoreValues(template.InputResource, "", test.store)

			log.Printf("store: %s", refs)

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
