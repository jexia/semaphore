package flow

import (
	"context"
	"io"
	"sync"

	"github.com/jexia/maestro/refs"
)

// Call represents a caller which could be called
type Call func(context.Context, io.Reader) (io.Reader, error)

// Manager is responsible for the handling of a flow and it's setps
type Manager struct {
	Codec      Codec
	Node       *Node
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

	processes := NewProcesses(1)
	tracker := NewTracker(manager.Nodes)

	manager.Node.Do(ctx, tracker, processes, refs)
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
	processes := NewProcesses(1)
	tracker := NewTracker(manager.Nodes)
	ends := make(map[string]*Node, manager.Ends)

	// Include all nodes to the revert tracker that have not been called
	manager.Node.Walk(ends, func(node *Node) {
		if !executed.Met(node) {
			tracker.Mark(node)
		}
	})

	for _, end := range ends {
		go end.Revert(ctx, tracker, processes, refs)
	}

	processes.Wait()
}

// Wait awaits till all calls and rollbacks are completed
func (manager *Manager) Wait() {
	manager.wg.Wait()
}
