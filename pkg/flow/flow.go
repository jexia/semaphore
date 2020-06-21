package flow

import (
	"context"
	"sync"

	"github.com/jexia/maestro/pkg/functions"
	"github.com/jexia/maestro/pkg/instance"
	"github.com/jexia/maestro/pkg/logger"
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
func NewManager(ctx instance.Context, name string, nodes []*Node, after functions.Stack, middleware *ManagerMiddleware) *Manager {
	ConstructBranches(nodes)

	if middleware == nil {
		middleware = &ManagerMiddleware{}
	}

	manager := &Manager{
		BeforeDo:       middleware.BeforeDo,
		BeforeRollback: middleware.BeforeRollback,
		ctx:            ctx,
		Name:           name,
		Starting:       FetchStarting(nodes),
		Nodes:          len(nodes),
		AfterFunctions: after,
		AfterDo:        middleware.AfterDo,
		AfterRollback:  middleware.AfterRollback,
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

// ManagerMiddleware holds the available middleware options for a flow manager
type ManagerMiddleware struct {
	BeforeDo       BeforeManager
	AfterDo        AfterManager
	BeforeRollback BeforeManager
	AfterRollback  AfterManager
}

// BeforeManager is called before a manager get's calles
type BeforeManager func(ctx context.Context, manager *Manager, store refs.Store) (context.Context, error)

// BeforeManagerHandler wraps the before call function to allow middleware to be chained
type BeforeManagerHandler func(BeforeManager) BeforeManager

// AfterManager is called after a manager is called
type AfterManager func(ctx context.Context, manager *Manager, store refs.Store) (context.Context, error)

// AfterManagerHandler wraps the after call function to allow middleware to be chained
type AfterManagerHandler func(AfterManager) AfterManager

// Manager is responsible for the handling of a flow and its steps
type Manager struct {
	BeforeDo       BeforeManager
	BeforeRollback BeforeManager
	ctx            instance.Context
	Name           string
	Starting       []*Node
	References     int
	Nodes          int
	Ends           int
	wg             sync.WaitGroup
	AfterFunctions functions.Stack
	AfterDo        AfterManager
	AfterRollback  AfterManager
}

// GetName returns the name of the given flow manager
func (manager *Manager) GetName() string {
	return manager.Name
}

// Do calls all the nodes inside the manager if a error is returned is a rollback of all the already executed steps triggered.
// Nodes are executed concurrently to one another.
func (manager *Manager) Do(ctx context.Context, refs refs.Store) (err error) {
	if manager.BeforeDo != nil {
		ctx, err = manager.BeforeDo(ctx, manager, refs)
		if err != nil {
			return err
		}
	}

	manager.wg.Add(1)
	defer manager.wg.Done()

	manager.ctx.Logger(logger.Flow).WithField("flow", manager.Name).Debug("Executing flow")

	processes := NewProcesses(len(manager.Starting))
	tracker := NewTracker(manager.Name, manager.Nodes)

	for _, node := range manager.Starting {
		go node.Do(ctx, tracker, processes, refs)
	}

	processes.Wait()

	manager.ctx.Logger(logger.Flow).WithField("flow", manager.Name).Debug("Processes completed")

	if manager.AfterFunctions != nil && processes.Err() == nil {
		err = ExecuteFunctions(manager.AfterFunctions, refs)
		if err != nil {
			processes.Fatal(err)
		}
	}

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

	if manager.AfterDo != nil {
		ctx, err = manager.AfterDo(ctx, manager, refs)
		if err != nil {
			return err
		}
	}

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

	var err error
	ctx := context.Background()

	if manager.BeforeRollback != nil {
		ctx, err = manager.BeforeRollback(ctx, manager, refs)
		if err != nil {
			manager.ctx.Logger(logger.Flow).Error("Revert failed before rollback returned a error: ", err)
			return
		}
	}

	tracker := NewTracker(manager.Name, manager.Nodes)
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
		go end.Rollback(ctx, tracker, processes, refs)
	}

	processes.Wait()

	if manager.AfterRollback != nil {
		ctx, err = manager.AfterRollback(ctx, manager, refs)
		if err != nil {
			manager.ctx.Logger(logger.Flow).Error("Revert failed after rollback returned a error: ", err)
			return
		}
	}
}

// Wait awaits till all calls and rollbacks are completed
func (manager *Manager) Wait() {
	manager.ctx.Logger(logger.Flow).WithField("flow", manager.Name).Info("Awaiting till all processes are completed")
	manager.wg.Wait()
}
