package flow

import (
	"errors"
	"testing"

	"github.com/jexia/maestro/pkg/transport"
)

func TestProcesses(t *testing.T) {
	expected := errors.New("expected")
	processes := NewProcesses(0)

	processes.Fatal(transport.WrapError(expected, nil))
	result := processes.Err()
	if !errors.Is(result, expected) {
		t.Errorf("unexpected result %s, expected %s", result, expected)
	}

	processes.Add(1)
	processes.Done()
}

func TestProcessesEmptyFatal(t *testing.T) {
	processes := NewProcesses(0)

	processes.Fatal(nil)
	result := processes.Err()
	if result != nil {
		t.Errorf("unexpected result %s, expected %+v", result, nil)
	}
}

func TestProcessesAlreadyThrownErr(t *testing.T) {
	processes := NewProcesses(0)
	expected := errors.New("expected")
	unexpected := errors.New("unexpected")

	processes.Fatal(transport.WrapError(expected, nil))
	processes.Fatal(transport.WrapError(unexpected, nil))

	result := processes.Err()
	if !errors.Is(result, expected) {
		t.Errorf("unexpected result %s, expected %s", result, expected)
	}
}
