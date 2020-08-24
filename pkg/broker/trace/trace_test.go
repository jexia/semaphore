package trace

import "testing"

type expressionFunc func() string

func (fn expressionFunc) Position() string { return fn() }

func TestNew(t *testing.T) {
	expr := expressionFunc(func() string { return "file:10" })

	tests := map[string][]Option{
		"unexpected error":            {WithMessage("unexpected error")},
		"unexpected error: component": {WithMessage("unexpected error: %s", "component")},
		"file:10 unexpected error":    {WithMessage("unexpected error"), WithExpression(expr)},
		"file:10 ":                    {WithExpression(expr)},
	}

	for expected, options := range tests {
		err := New(options...)
		if err.Error() != expected {
			t.Errorf("unexpected result %s, expected %s", err, expected)
		}
	}
}
