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

func TestUndefinedSpecs(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			"return the formatted error",
			"no object specs defined",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := ErrUndefinedSpecs{}
			if got := e.Prettify(); got.Message != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}
