package sprintf

import (
	"errors"
	"testing"
)

func TestStatefulScanner(t *testing.T) {
	type test struct {
		input        string
		constructors []Constructor
		tokens       []string
		error        error
	}

	tests := map[string]test{
		"missing verb": {
			input:        "%.2",
			constructors: []Constructor{},
			error:        errIncomplete,
		},
		"unknown verb": {
			input:        "%.2unknown",
			constructors: []Constructor{},
			error:        errUnknownFormatter,
		},
		"float": {
			input:        "%f polizei + %4.2f polizei",
			constructors: []Constructor{Float{}},
			tokens: []string{
				"", // empty since there are no static input before the first verb
				"%0.0",
				"f",
				" polizei + ",
				"%4.2",
				"f",
				" polizei",
			},
		},
		"string": {
			input:        "%s, what is that %s?",
			constructors: []Constructor{String{}},
			tokens: []string{
				"", // empty since there are no static input before the first verb
				"%0.0",
				"s",
				", what is that ",
				"%0.0",
				"s",
				"?",
			},
		},
		"json": {
			input:        `{"array":%json,"object":%json}`,
			constructors: []Constructor{JSON{}},
			tokens:       []string{`{"array":`, "%0.0", "json", `,"object":`, "%0.0", "json", "}"},
		},
	}

	for title, test := range tests {
		t.Run(title, func(t *testing.T) {
			var detector = NewRadix()

			for _, constructor := range test.constructors {
				if err := detector.Register(constructor); err != nil {
					t.Fatalf("unexpected error: %s", err)
				}
			}

			var scanner = NewScanner(detector)

			tokens, err := scanner.Scan(test.input)
			if test.error != nil {
				if !errors.Is(err, test.error) {
					t.Errorf("unexpected error: %s", err)
				}
			}

			if test.error != nil {
				if err == nil {
					t.Errorf("error %q was expected", test.error)
				}

				if err != nil && !errors.Is(err, test.error) {
					t.Errorf("error %q, was expected to be %q", err, test.error)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error %q", err)
				}
			}

			if actual, expected := len(tokens), len(test.tokens); actual != expected {
				t.Errorf("got %d tokens but expected %d", actual, expected)

				return
			}

			for index, token := range tokens {
				if actual := token.String(); actual != test.tokens[index] {
					t.Errorf("token %q with position %d does not match expected token %q", actual, index, test.tokens[index])
				}
			}
		})
	}
}
