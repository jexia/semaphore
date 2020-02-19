package flow

import (
	"context"
	"errors"
	"io"
	"testing"

	"github.com/jexia/maestro/refs"
)

type MockCodec struct{}

func (codec *MockCodec) Marshal(*refs.Store) (io.Reader, error) {
	return nil, nil
}

func (codec *MockCodec) Unmarshal(io.Reader, *refs.Store) error {
	return nil
}

type caller struct {
	Counter int
}

func (caller *caller) Call(context.Context, io.Reader) (io.Reader, error) {
	caller.Counter++
	return nil, nil
}

func NewFlowManager(caller Call, revert Call) ([]*Node, *Manager) {
	nodes := []*Node{
		NewNode("first", caller, revert),
		NewNode("second", caller, revert),
		NewNode("third", caller, revert),
		NewNode("fourth", caller, revert),
	}

	nodes[0].Next = []*Node{nodes[1], nodes[2]}

	nodes[1].Previous = []*Node{nodes[0]}
	nodes[1].Next = []*Node{nodes[3]}
	nodes[2].Previous = []*Node{nodes[0]}
	nodes[2].Next = []*Node{nodes[3]}

	nodes[3].Previous = []*Node{nodes[1], nodes[2]}

	return nodes, &Manager{
		Codec:      &MockCodec{},
		Node:       nodes[0],
		References: 0,
		Nodes:      len(nodes),
		Ends:       1,
	}
}

func TestCallFlowManager(t *testing.T) {
	caller := &caller{}
	nodes, manager := NewFlowManager(caller.Call, nil)
	_, err := manager.Call(context.Background(), nil)
	if err != nil {
		t.Error(err)
	}

	if caller.Counter != len(nodes) {
		t.Errorf("unexpected counter total %d, expected %d", caller.Counter, len(nodes))
	}
}

func TestFailFlowManager(t *testing.T) {
	expected := errors.New("something went wrong")
	reverts := 2
	calls := 2

	rollback := &caller{}
	caller := &caller{}

	nodes, manager := NewFlowManager(caller.Call, rollback.Call)

	nodes[2].Call = func(context.Context, io.Reader) (io.Reader, error) {
		return nil, expected
	}

	_, err := manager.Call(context.Background(), nil)
	if err != expected {
		t.Errorf("unexpected result %s, expected %s", err, expected)
	}

	manager.Wait()

	if caller.Counter != calls {
		t.Errorf("unexpected counter total %d, expected %d", caller.Counter, calls)
	}

	if rollback.Counter != reverts {
		t.Errorf("unexpected rollback counter total %d, expected %d", rollback.Counter, reverts)
	}
}
