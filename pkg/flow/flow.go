package flow

import (
	"context"
	"sync"

	"github.com/jexia/maestro/internal/logger"
	"github.com/jexia/maestro/pkg/instance"
	"github.com/jexia/maestro/pkg/refs"
	"github.com/jexia/maestro/pkg/specs"
	"github.com/sirupsen/logrus"
)

// Call represents a transport caller implementation
type Call interface {
	References() []*specs.Property
	Do(context.Context, refs.Store) error
}

// NewManager constructs a new manager for the given flow.
// Branches are constructed for the constructed nodes to optimalise performance.
// Various variables such as the amount of nodes, references and loose ends are collected to optimalise allocations during runtime.
func NewManager(ctx instance.Context, name string, nodes []*Node) *Manager {
	ConstructBranches(nodes)

	manager := &Manager{
		ctx:      ctx,
		Name:     name,
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
	ctx        instance.Context
	Name       string
	Starting   []*Node
	References int
	Nodes      int
	Ends       int
	wg         sync.WaitGroup
}

// GetName returns the name of the given flow manager
func (manager *Manager) GetName() string {
	return manager.Name
}

// Call calls all the nodes inside the manager if a error is returned is a rollback of all the already executed steps triggered.
// Nodes are executed concurrently to one another.
func (manager *Manager) Call(ctx context.Context, refs refs.Store) error {
	manager.wg.Add(1)
	defer manager.wg.Done()

	manager.ctx.Logger(logger.Flow).WithField("flow", manager.Name).Debug("Executing flow")

	processes := NewProcesses(len(manager.Starting))
	tracker := NewTracker(manager.Nodes)

	for _, node := range manager.Starting {
		go node.Do(ctx, tracker, processes, refs)
	}

	processes.Wait()

	manager.ctx.Logger(logger.Flow).WithField("flow", manager.Name).Debug("Processes completed")

	if processes.Err() != nil {
		manager.ctx.Logger(logger.Flow).WithFields(logrus.Fields{
			"flow": manager.Name,
			"err":  processes.Err(),
		}).Error("An error occurred, executing rollback")

		manager.wg.Add(1)
		go manager.Revert(tracker, refs)
		return processes.Err()
	}

	manager.ctx.Logger(logger.Flow).WithField("flow", manager.Name).Debug("Flow completed")
	return nil
}

// NewStore constructs a new reference store for the given manager
func (manager *Manager) NewStore() refs.Store {
	return refs.NewReferenceStore(manager.References)
}

// Revert reverts the executed nodes found inside the given tracker.
// All nodes that have not been executed will be ignored.
func (manager *Manager) Revert(executed *Tracker, refs refs.Store) {
	defer manager.wg.Done()

	ctx := context.Background()
	tracker := NewTracker(manager.Nodes)
	ends := make(map[string]*Node, manager.Ends)

	// Include all nodes to the revert tracker that have not been called
	for _, node := range manager.Starting {
		node.Walk(ends, func(node *Node) {
			if executed.Reached(node, len(node.Previous)) {
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
	manager.ctx.Logger(logger.Flow).WithField("flow", manager.Name).Info("Awaiting till all processes are completed")
	manager.wg.Wait()
}
