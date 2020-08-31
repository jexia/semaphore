package formencoded

import "testing"

func TestErrUndefinedProperty(t *testing.T) {
	var (
		err      = errUndefinedProperty("foo")
		expected = `undefined property "foo"`
	)

	if actual := err.Error(); actual != expected {
		t.Errorf("error %q was expected to be %q", actual, expected)
	}
}

func TestErrUnknownLabel(t *testing.T) {
	var (
		err      = errUnknownLabel("foo")
		expected = `unknown label "foo"`
	)

	if actual := err.Error(); actual != expected {
		t.Errorf("error %q was expected to be %q", actual, expected)
	}
}

func TestErrUnknownType(t *testing.T) {
	var (
		err      = errUnknownType("foo")
		expected = `unknown data type "foo"`
	)

	if actual := err.Error(); actual != expected {
		t.Errorf("error %q was expected to be %q", actual, expected)
	}
}
