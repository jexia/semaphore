package flow

import (
	"context"
	"io"
	"sync"

	"github.com/jexia/maestro/codec"
	"github.com/jexia/maestro/refs"
	"github.com/jexia/maestro/services"
	"github.com/jexia/maestro/specs"
)

// NewManager constructs a new manager for the given flow.
// Branches are constructed for the constructed nodes to optimalise performance.
// Various variables such as the ammount of nodes, references and loose ends are collected to optimalise allocations during runtime.
func NewManager(flow specs.FlowManager, codec codec.Manager, services services.Collection) *Manager {
	nodes := make([]*Node, len(flow.GetCalls()))

	for index, call := range flow.GetCalls() {
		nodes[index] = NewNode(call, services)
	}

	ConstructBranches(nodes)

	manager := &Manager{
		Codec:    codec,
		Starting: FetchStarting(nodes),
		Nodes:    len(nodes),
	}

	ends := make(map[string]*Node, len(nodes))
	for _, node := range manager.Starting {
		node.Walk(ends, func(node *Node) {
			manager.References += len(node.References)
		})
	}

	manager.Ends = len(ends)

	return manager
}

// Manager is responsible for the handling of a flow and its steps
type Manager struct {
	Codec      codec.Manager
	Starting   []*Node
	References int
	Nodes      int
	Ends       int
	wg         sync.WaitGroup
}

// Call calls all the nodes inside the manager if a error is returned is a rollback of all the already executed steps triggered.
// Nodes are executed concurrently to one another.
func (manager *Manager) Call(ctx context.Context, reader io.Reader) (io.Reader, error) {
	manager.wg.Add(1)
	defer manager.wg.Done()

	refs := refs.NewStore(manager.References)
	err := manager.Codec.Unmarshal(reader, refs)
	if err != nil {
		return nil, err
	}

	processes := NewProcesses(len(manager.Starting))
	tracker := NewTracker(manager.Nodes)

	for _, node := range manager.Starting {
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

// Revert reverts the executed nodes found inside the given tracker.
// All nodes that have not been executed will be ignored.
func (manager *Manager) Revert(executed *Tracker, refs *refs.Store) {
	defer manager.wg.Done()

	ctx := context.Background()
	tracker := NewTracker(manager.Nodes)
	ends := make(map[string]*Node, manager.Ends)

	// Include all nodes to the revert tracker that have not been called
	for _, node := range manager.Starting {
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
