package flow

import (
	"context"
	"testing"

	"github.com/jexia/maestro/internal/logger"
	"github.com/jexia/maestro/pkg/instance"
	"github.com/jexia/maestro/pkg/refs"
	"github.com/jexia/maestro/pkg/specs"
)

func NewMockNode(name string, caller Call, rollback Call) *Node {
	ctx := instance.NewContext()
	logger := ctx.Logger(logger.Flow)

	return &Node{
		ctx:        ctx,
		logger:     logger,
		Name:       name,
		Call:       caller,
		Rollback:   rollback,
		DependsOn:  map[string]*specs.Node{},
		References: map[string]*specs.PropertyReference{},
	}
}

func BenchmarkSingleNodeCalling(b *testing.B) {
	caller := &caller{}
	node := NewMockNode("first", caller, nil)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ctx := context.Background()
		tracker := NewTracker(1)
		processes := NewProcesses(1)
		refs := refs.NewReferenceStore(0)

		node.Do(ctx, tracker, processes, refs)
	}
}

func BenchmarkBranchedNodeCalling(b *testing.B) {
	caller := &caller{}
	nodes := []*Node{
		NewMockNode("first", caller, nil),
		NewMockNode("second", caller, nil),
		NewMockNode("third", caller, nil),
	}

	nodes[0].Next = []*Node{nodes[1]}
	nodes[1].Previous = []*Node{nodes[0]}
	nodes[1].Next = []*Node{nodes[2]}
	nodes[2].Previous = []*Node{nodes[1]}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ctx := context.Background()
		tracker := NewTracker(len(nodes))
		processes := NewProcesses(1)
		refs := refs.NewReferenceStore(0)

		nodes[0].Do(ctx, tracker, processes, refs)
	}
}

func TestNodeCalling(t *testing.T) {
	caller := &caller{}

	nodes := []*Node{
		NewMockNode("first", caller, nil),
		NewMockNode("second", caller, nil),
		NewMockNode("third", caller, nil),
	}

	nodes[0].Next = []*Node{nodes[1]}
	nodes[1].Previous = []*Node{nodes[0]}
	nodes[1].Next = []*Node{nodes[2]}
	nodes[2].Previous = []*Node{nodes[1]}

	tracker := NewTracker(len(nodes))
	processes := NewProcesses(1)
	refs := refs.NewReferenceStore(0)

	nodes[0].Do(context.Background(), tracker, processes, refs)
	processes.Wait()

	if processes.Err() != nil {
		t.Error(processes.Err())
	}

	if caller.Counter != len(nodes) {
		t.Errorf("unexpected counter total %d, expected %d", caller.Counter, len(nodes))
	}
}

func TestNodeRevert(t *testing.T) {
	rollback := &caller{}

	nodes := []*Node{
		NewMockNode("first", nil, rollback),
		NewMockNode("second", nil, rollback),
		NewMockNode("third", nil, rollback),
	}

	nodes[0].Next = []*Node{nodes[1]}
	nodes[1].Previous = []*Node{nodes[0]}
	nodes[1].Next = []*Node{nodes[2]}
	nodes[2].Previous = []*Node{nodes[1]}

	tracker := NewTracker(len(nodes))
	processes := NewProcesses(1)
	refs := refs.NewReferenceStore(0)

	nodes[len(nodes)-1].Revert(context.Background(), tracker, processes, refs)
	processes.Wait()

	if processes.Err() != nil {
		t.Error(processes.Err())
	}

	if rollback.Counter != len(nodes) {
		t.Errorf("unexpected counter total %d, expected %d", rollback.Counter, len(nodes))
	}
}

func TestNodeBranchesCalling(t *testing.T) {
	caller := &caller{}

	nodes := []*Node{
		NewMockNode("first", caller, nil),
		NewMockNode("second", caller, nil),
		NewMockNode("third", caller, nil),
		NewMockNode("fourth", caller, nil),
	}

	nodes[0].Next = []*Node{nodes[1], nodes[2]}

	nodes[1].Previous = []*Node{nodes[0]}
	nodes[1].Next = []*Node{nodes[3]}
	nodes[2].Previous = []*Node{nodes[0]}
	nodes[2].Next = []*Node{nodes[3]}

	nodes[3].Previous = []*Node{nodes[1], nodes[2]}

	tracker := NewTracker(len(nodes))
	processes := NewProcesses(1)
	refs := refs.NewReferenceStore(0)

	nodes[0].Do(context.Background(), tracker, processes, refs)
	processes.Wait()

	if processes.Err() != nil {
		t.Error(processes.Err())
	}

	if caller.Counter != len(nodes) {
		t.Errorf("unexpected counter total %d, expected %d", caller.Counter, len(nodes))
	}
}
