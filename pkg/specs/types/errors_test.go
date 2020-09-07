package types

import "testing"

func TestErrUnknownType(t *testing.T) {
	var (
		err      = ErrUnknownType("foo")
		expected = `unknown data type "foo"`
	)

	if actual := err.Error(); actual != expected {
		t.Errorf("error %q was expected to be %q", actual, expected)
	}
}
