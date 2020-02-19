package flow

import (
	"context"
	"io"
	"sync"

	"github.com/jexia/maestro/refs"
	"github.com/jexia/maestro/specs"
)

// NewManager constructs a new manager for the given flow
func NewManager(flow specs.FlowManager, codec Codec, services Services) *Manager {
	nodes := make([]*Node, len(flow.GetCalls()))

	for index, call := range flow.GetCalls() {
		nodes[index] = NewNode(call, services)
	}

	ConstructBranches(nodes)

	manager := &Manager{
		Codec: codec,
		Seed:  nodes,
		Nodes: len(nodes),
	}

	ends := make(map[string]*Node, len(nodes))
	for _, node := range manager.Seed {
		node.Walk(ends, func(node *Node) {
			manager.References += len(node.References)
		})
	}

	manager.Ends = len(ends)

	return manager
}

// Call represents a caller which could be called
type Call func(context.Context, io.Reader) (io.Reader, error)

// Manager is responsible for the handling of a flow and it's setps
type Manager struct {
	Codec      Codec
	Seed       []*Node
	References int
	Nodes      int
	Ends       int
	wg         sync.WaitGroup
}

// Call calls all the steps inside the manager if a error is returned is a rollback of all the already executed steps triggered
func (manager *Manager) Call(ctx context.Context, reader io.Reader) (io.Reader, error) {
	manager.wg.Add(1)
	defer manager.wg.Done()

	refs := refs.NewStore(manager.References)
	err := manager.Codec.Unmarshal(reader, refs)
	if err != nil {
		return nil, err
	}

	processes := NewProcesses(len(manager.Seed))
	tracker := NewTracker(manager.Nodes)

	for _, node := range manager.Seed {
		go node.Do(ctx, tracker, processes, refs)
	}

	processes.Wait()

	if processes.Err() != nil {
		manager.wg.Add(1)
		go manager.Revert(tracker, refs)
		return nil, processes.Err()
	}

	reader, err = manager.Codec.Marshal(refs)
	if err != nil {
		return nil, err
	}

	return reader, nil
}

// Revert reverts the nodes available inside the given tracker
func (manager *Manager) Revert(executed *Tracker, refs *refs.Store) {
	defer manager.wg.Done()

	ctx := context.Background()
	tracker := NewTracker(manager.Nodes)
	ends := make(map[string]*Node, manager.Ends)

	// Include all nodes to the revert tracker that have not been called
	for _, node := range manager.Seed {
		node.Walk(ends, func(node *Node) {
			if !executed.Met(node) {
				tracker.Mark(node)
			}
		})
	}

	processes := NewProcesses(len(ends))

	for _, end := range ends {
		go end.Revert(ctx, tracker, processes, refs)
	}

	processes.Wait()
}

// Wait awaits till all calls and rollbacks are completed
func (manager *Manager) Wait() {
	manager.wg.Wait()
}
