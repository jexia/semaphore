package endpoints

import (
	"testing"

	"github.com/jexia/semaphore/pkg/prettyerr"
)

func TestErrNoServiceForMethod(t *testing.T) {
	type fields struct {
		Method string
	}

	type test struct {
		fields   fields
		expected string
	}

	tests := map[string]test{
		"return the formatted error": {
			fields:   fields{Method: "get"},
			expected: "unknown service 'get'",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			e := ErrUnknownService{
				Service: tt.fields.Method,
			}
			if got := e.Error(); got != tt.expected {
				t.Errorf("%v, want %v", got, tt.expected)
			}
		})
	}
}

func TestErrNoServiceForMethodPretty(t *testing.T) {
	type fields struct {
		Method string
	}

	type test struct {
		fields   fields
		expected prettyerr.Error
	}

	tests := map[string]test{
		"return pretty error": {
			fields: fields{
				Method: "get",
			},
			expected: prettyerr.Error{
				Message: "unknown service 'get'",
				Details: map[string]interface{}{"method": "get"},
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			err := ErrUnknownService{
				Service: tt.fields.Method,
			}

			result := err.Prettify()
			if result.Error() != tt.expected.Error() {
				t.Errorf("%v, want %v", result, tt.expected)
			}
		})
	}
}
