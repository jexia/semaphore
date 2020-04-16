package flow

import (
	"context"
	"errors"
	"io"
	"sync"
	"testing"

	"github.com/jexia/maestro/pkg/instance"
	"github.com/jexia/maestro/pkg/refs"
	"github.com/jexia/maestro/pkg/specs"
)

type MockCodec struct{}

func (codec *MockCodec) Marshal(refs.Store) (io.Reader, error) {
	return nil, nil
}

func (codec *MockCodec) Unmarshal(io.Reader, refs.Store) error {
	return nil
}

type caller struct {
	Counter int
	mutex   sync.Mutex
	Err     error
}

func (caller *caller) References() []*specs.Property {
	return nil
}

func (caller *caller) Do(context.Context, refs.Store) error {
	caller.mutex.Lock()
	caller.Counter++
	caller.mutex.Unlock()
	return caller.Err
}

func NewMockFlowManager(caller Call, revert Call) ([]*Node, *Manager) {
	ctx := instance.NewContext()

	nodes := []*Node{
		NewMockNode("first", caller, revert),
		NewMockNode("second", caller, revert),
		NewMockNode("third", caller, revert),
		NewMockNode("fourth", caller, revert),
	}

	nodes[0].Next = []*Node{nodes[1], nodes[2]}

	nodes[1].Previous = []*Node{nodes[0]}
	nodes[1].Next = []*Node{nodes[3]}
	nodes[2].Previous = []*Node{nodes[0]}
	nodes[2].Next = []*Node{nodes[3]}

	nodes[3].Previous = []*Node{nodes[1], nodes[2]}

	return nodes, &Manager{
		ctx:        ctx,
		Starting:   []*Node{nodes[0]},
		References: 0,
		Nodes:      len(nodes),
		Ends:       1,
	}
}

func TestCallFlowManager(t *testing.T) {
	caller := &caller{}
	nodes, manager := NewMockFlowManager(caller, nil)
	err := manager.Call(context.Background(), nil)
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
	call := &caller{}

	nodes, manager := NewMockFlowManager(call, rollback)

	nodes[2].Call = &caller{Err: expected}

	err := manager.Call(context.Background(), nil)
	if err != expected {
		t.Errorf("unexpected result %s, expected %s", err, expected)
	}

	manager.Wait()

	if call.Counter != calls {
		t.Errorf("unexpected counter total %d, expected %d", call.Counter, calls)
	}

	if rollback.Counter != reverts {
		t.Errorf("unexpected rollback counter total %d, expected %d", rollback.Counter, reverts)
	}
}
