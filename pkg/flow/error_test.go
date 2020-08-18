package flow

import (
	"testing"

	"github.com/jexia/semaphore/pkg/specs"
)

func TestNewErrorNil(t *testing.T) {
	handle := NewOnError(nil, nil, nil, nil)
	if handle == nil {
		t.Fatal("empty handle")
	}
}

func TestNewErrorErr(t *testing.T) {
	err := &specs.OnError{
		Status:  &specs.Property{},
		Message: &specs.Property{},
	}

	handle := NewOnError(nil, nil, nil, err)
	if handle == nil {
		t.Fatal("empty handle")
	}

	if handle.status != err.Status {
		t.Errorf("unexpected handle status, expected to be the same as on error handle")
	}

	if handle.message != err.Message {
		t.Errorf("unexpected handle message, expected to be the same as on error handle")
	}
}
