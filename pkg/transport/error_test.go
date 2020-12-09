package transport

import (
	"errors"
	"testing"

	"github.com/jexia/semaphore/pkg/specs"
)

func TestUnwrap(t *testing.T) {
	typed := errors.New("mock error")
	handle := &specs.OnError{}
	err := WrapError(typed, handle)

	result := Unwrap(err)
	if result.Ptr() != handle {
		t.Errorf("unexpected result %+v, expected %+v", result, handle)
	}
}

func TestUnwrapNil(t *testing.T) {
	result := Unwrap(nil)
	if result != nil {
		t.Errorf("unexpected result %+v", result)
	}
}

func TestWrapper(t *testing.T) {
	typed := errors.New("mock error")
	handle := &specs.OnError{}

	wrapper := wrapper{
		ErrorHandle: handle,
		err:         typed,
	}

	if wrapper.String() != typed.Error() {
		t.Errorf("unexpected err string %s, expected %s", wrapper.String(), typed.Error())
	}

	if wrapper.Error() != typed.Error() {
		t.Errorf("unexpected err string %s, expected %s", wrapper.Error(), typed.Error())
	}

	if wrapper.Unwrap() != typed {
		t.Errorf("unexpected unwrap %+v, expected %+v", wrapper.Unwrap(), typed)
	}

	if wrapper.Ptr() != handle {
		t.Errorf("unexpected handle %+v, expected %+v", wrapper.Ptr(), handle)
	}
}

func TestWrapperNil(t *testing.T) {
	wrapper := wrapper{}

	if wrapper.String() != "" {
		t.Errorf("unexpected err string %s, expected empty string", wrapper.String())
	}

	if wrapper.Error() != "" {
		t.Errorf("unexpected err string %s, expected empty string", wrapper.Error())
	}

	if wrapper.Unwrap() != nil {
		t.Errorf("unexpected unwrap %+v, expected %+v", wrapper.Unwrap(), nil)
	}

	if wrapper.Ptr() != nil {
		t.Errorf("unexpected handle %+v, expected %+v", wrapper.Ptr(), nil)
	}
}

func TestErrMalformedTemplate(t *testing.T) {
	var (
		err = ErrMalformedTemplate{
			Template: "something went wrong here",
			Position: 21,
			Cause:    "failed to compile the template",
		}

		expected = "failed to compile the template:\n                     ↓\nsomething went wrong here\n"
	)

	if actual := err.Error(); actual != expected {
		t.Errorf("the error:\n%s\nwas expected to be:\n%s", actual, expected)
	}
}
