package listeners

import (
	"testing"

	"github.com/jexia/semaphore/pkg/prettyerr"
)

func TestErrNoListener_Error(t *testing.T) {
	type fields struct {
		Listener string
	}

	type test struct {
		fields   fields
		expected string
	}

	tests := map[string]test{
		"return formatted error": {
			fields:   fields{Listener: "dogs"},
			expected: "unknown listener 'dogs'",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			err := ErrNoListener{
				Listener: tt.fields.Listener,
			}

			result := err.Error()
			if result != tt.expected {
				t.Errorf("%v, want %v", result, tt.expected)
			}
		})
	}
}

func TestErrNoListener_Prettify(t *testing.T) {
	type fields struct {
		Listener string
	}

	type test struct {
		fields   fields
		expected prettyerr.Error
	}

	tests := map[string]test{
		"return pretty error": {
			fields: fields{Listener: "dogs"},
			expected: prettyerr.Error{
				Message: "unknown listener 'dogs'",
				Details: map[string]interface{}{"listener": "dogs"},
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			err := ErrNoListener{
				Listener: tt.fields.Listener,
			}

			result := err.Prettify()
			if result.Error() != tt.expected.Error() {
				t.Errorf("%v, want %v", result, tt.expected)
			}
		})
	}
}
