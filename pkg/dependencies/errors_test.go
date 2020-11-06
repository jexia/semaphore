package dependencies

import "testing"

func TestErrCircularDependency(t *testing.T) {
	var (
		expected = "circular resource dependency detected: flow.a <-> flow.b"

		err = ErrCircularDependency{
			Flow: "flow",
			From: "a",
			To:   "b",
		}
	)

	if actual := err.Error(); actual != expected {
		t.Errorf("error message %q was expected to be %q", actual, expected)
	}
}
