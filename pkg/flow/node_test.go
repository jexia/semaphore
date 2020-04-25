package flow

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jexia/maestro/pkg/instance"
	"github.com/jexia/maestro/pkg/logger"
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

func TestConstructingNode(t *testing.T) {
	type test struct {
		Node     *specs.Node
		Call     Call
		Rollback Call
		Expected int
	}

	tests := map[string]*test{
		"node call": {
			Expected: 2,
			Node: &specs.Node{
				Call: &specs.Call{
					Request: &specs.ParameterMap{
						Property: &specs.Property{
							Nested: map[string]*specs.Property{
								"first": {
									Reference: &specs.PropertyReference{
										Resource: "input",
										Path:     "first",
									},
								},
								"second": {
									Reference: &specs.PropertyReference{
										Resource: "input",
										Path:     "second",
									},
								},
							},
						},
					},
				},
			},
		},
		"combination": {
			Expected: 1,
			Node: &specs.Node{
				Call: &specs.Call{
					Request: &specs.ParameterMap{
						Property: &specs.Property{
							Nested: map[string]*specs.Property{
								"first": {
									Reference: &specs.PropertyReference{
										Resource: "input",
										Path:     "first",
									},
								},
							},
						},
					},
				},
			},
			Call: &caller{
				references: []*specs.Property{
					{
						Reference: &specs.PropertyReference{
							Resource: "input",
							Path:     "first",
						},
					},
				},
			},
			Rollback: &caller{
				references: []*specs.Property{
					{
						Reference: &specs.PropertyReference{
							Resource: "input",
							Path:     "first",
						},
					},
				},
			},
		},
		"call references": {
			Expected: 2,
			Node:     &specs.Node{},
			Call: &caller{
				references: []*specs.Property{
					{
						Reference: &specs.PropertyReference{
							Resource: "input",
							Path:     "first",
						},
					},
					{
						Reference: &specs.PropertyReference{
							Resource: "input",
							Path:     "second",
						},
					},
				},
			},
		},
		"rollback references": {
			Expected: 2,
			Node:     &specs.Node{},
			Rollback: &caller{
				references: []*specs.Property{
					{
						Reference: &specs.PropertyReference{
							Resource: "input",
							Path:     "first",
						},
					},
					{
						Reference: &specs.PropertyReference{
							Resource: "input",
							Path:     "second",
						},
					},
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := instance.NewContext()
			result := NewNode(ctx, test.Node, test.Call, test.Rollback)

			if len(result.References) != test.Expected {
				t.Fatalf("unexpected amount of references %d, expected %d", len(result.References), test.Expected)
			}
		})
	}
}

func TestConstructingNodeReferences(t *testing.T) {
	ctx := instance.NewContext()
	call := &caller{}
	rollback := &caller{}

	node := &specs.Node{
		Name: "mock",
	}

	result := NewNode(ctx, node, call, rollback)
	if result == nil {
		t.Fatal("nil node returned")
	}
}

func TestNodeHas(t *testing.T) {
	nodes := make(Nodes, 2)

	nodes[0] = &Node{Name: "first"}
	nodes[1] = &Node{Name: "second"}

	if !nodes.Has("first") {
		t.Fatal("unexpected result, expected 'first' to be available")
	}

	if nodes.Has("unexpected") {
		t.Fatal("unexpected result, expected 'unexpected' to be unavailable")
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

func TestSlowNodeAbortingOnErr(t *testing.T) {
	slow := &caller{name: "slow"}
	failed := &caller{name: "failed", Err: errors.New("unexpected")}
	caller := &caller{}

	nodes := []*Node{
		NewMockNode("first", caller, nil),
		NewMockNode("second", slow, nil),
		NewMockNode("third", failed, nil),
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

	slow.mutex.Lock()
	failed.mutex.Lock()

	go func() {
		failed.mutex.Unlock()
		time.Sleep(100 * time.Millisecond)
		slow.mutex.Unlock()
	}()

	nodes[0].Do(context.Background(), tracker, processes, refs)

	processes.Wait()

	counter := (caller.Counter + slow.Counter + failed.Counter)
	if counter != 3 {
		t.Fatalf("unexpected counter total %d, expected %d", counter, 3)
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
