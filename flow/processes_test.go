package flow

import (
	"errors"
	"testing"
)

func TestProcesses(t *testing.T) {
	expected := errors.New("expected")
	processes := NewProcesses(0)

	processes.Fatal(expected)
	result := processes.Err()
	if result != expected {
		t.Errorf("unexpected result %s, expected %s", result, expected)
	}

	processes.Add(1)
	processes.Done()
}
