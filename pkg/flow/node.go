package flow

import (
	"context"

	"github.com/jexia/maestro/pkg/instance"
	"github.com/jexia/maestro/pkg/logger"
	"github.com/jexia/maestro/pkg/refs"
	"github.com/jexia/maestro/pkg/specs"
	"github.com/sirupsen/logrus"
)

// NewNode constructs a new node for the given call.
// The service called inside the call endpoint is retrieved from the services collection.
// The call, codec and rollback are defined inside the node and used while processing requests.
func NewNode(ctx instance.Context, node *specs.Node, call, rollback Call, middleware *NodeMiddleware) *Node {
	references := refs.References{}

	if middleware == nil {
		middleware = &NodeMiddleware{}
	}

	if node.Call != nil {
		references.MergeLeft(refs.ParameterReferences(node.Call.Request))
	}

	if call != nil {
		for _, prop := range call.References() {
			references.MergeLeft(refs.PropertyReferences(prop))
		}
	}

	if rollback != nil {
		for _, prop := range rollback.References() {
			references.MergeLeft(refs.PropertyReferences(prop))
		}
	}

	logger := ctx.Logger(logger.Flow)

	return &Node{
		BeforeDo:     middleware.BeforeDo,
		BeforeRevert: middleware.BeforeRollback,
		ctx:          ctx,
		logger:       logger,
		Name:         node.Name,
		Previous:     []*Node{},
		Call:         call,
		Revert:       rollback,
		DependsOn:    node.DependsOn,
		References:   references,
		Next:         []*Node{},
		AfterDo:      middleware.AfterDo,
		AfterRevert:  middleware.AfterRollback,
	}
}

// Nodes represents a node collection
type Nodes []*Node

// Has checks whether the given node collection has a node with the given name inside
func (nodes Nodes) Has(name string) bool {
	for _, node := range nodes {
		if node.Name == name {
			return true
		}
	}

	return false
}

// NodeMiddleware holds all the available
type NodeMiddleware struct {
	BeforeDo       BeforeNode
	AfterDo        AfterNode
	BeforeRollback BeforeNode
	AfterRollback  AfterNode
}

// BeforeNode is called before a node is executed
type BeforeNode func(ctx context.Context, node *Node, tracker *Tracker, processes *Processes, store refs.Store) (context.Context, error)

// BeforeNodeHandler wraps the before node function to allow middleware to be chained
type BeforeNodeHandler func(BeforeNode) BeforeNode

// AfterNode is called after a node is executed
type AfterNode func(ctx context.Context, node *Node, tracker *Tracker, processes *Processes, store refs.Store) (context.Context, error)

// AfterNodeHandler wraps the after node function to allow middleware to be chained
type AfterNodeHandler func(AfterNode) AfterNode

// Node represents a collection of callers and rollbacks which could be executed parallel.
type Node struct {
	BeforeDo     BeforeNode
	BeforeRevert BeforeNode
	ctx          instance.Context
	logger       *logrus.Logger
	Name         string
	Previous     Nodes
	Call         Call
	Revert       Call
	DependsOn    map[string]*specs.Node
	References   map[string]*specs.PropertyReference
	Next         Nodes
	AfterDo      AfterNode
	AfterRevert  AfterNode
}

// Do executes the given node an calls the next nodes.
// If one of the nodes fails is the error marked and are the processes aborted.
func (node *Node) Do(ctx context.Context, tracker *Tracker, processes *Processes, refs refs.Store) {
	defer processes.Done()
	node.logger.Debug("Executing node call: ", node.Name)

	tracker.Lock(node)
	defer tracker.Unlock(node)

	if !tracker.Reached(node, len(node.Previous)) {
		node.logger.Debug("Has not met dependencies yet: ", node.Name)
		return
	}

	var err error

	if node.BeforeDo != nil {
		ctx, err = node.BeforeDo(ctx, node, tracker, processes, refs)
		if err != nil {
			node.logger.Error("Node before middleware failed: ", err)
			processes.Fatal(err)
			return
		}
	}

	if node.Call != nil {
		err = node.Call.Do(ctx, refs)
		if err != nil {
			node.logger.Error("Call failed: ", node.Name)
			processes.Fatal(err)
			return
		}
	}

	node.logger.Debug("Marking node as completed: ", node.Name)
	tracker.Mark(node)

	if processes.Err() != nil {
		node.logger.Error("Stopping execution a error has been thrown: ", node.Name)
		return
	}

	processes.Add(len(node.Next))
	for _, next := range node.Next {
		tracker.Mark(next)
		go next.Do(ctx, tracker, processes, refs)
	}

	if node.AfterDo != nil {
		ctx, err = node.AfterDo(ctx, node, tracker, processes, refs)
		if err != nil {
			node.logger.Error("Node after middleware failed: ", err)
			processes.Fatal(err)
			return
		}
	}
}

// Rollback executes the given node rollback an calls the previous nodes.
// If one of the nodes fails is the error marked but execution is not aborted.
func (node *Node) Rollback(ctx context.Context, tracker *Tracker, processes *Processes, refs refs.Store) {
	defer processes.Done()
	node.logger.Debug("Executing node revert ", node.Name)

	tracker.Lock(node)
	defer tracker.Unlock(node)

	if !tracker.Reached(node, len(node.Next)) {
		node.logger.Debug("Has not met dependencies yet: ", node.Name)
		return
	}

	var err error

	if node.BeforeRevert != nil {
		ctx, err = node.BeforeRevert(ctx, node, tracker, processes, refs)
		if err != nil {
			node.logger.Error("Node before middleware failed: ", err)
			processes.Fatal(err)
			return
		}
	}

	defer func() {
		processes.Add(len(node.Previous))
		for _, node := range node.Previous {
			tracker.Mark(node)
			go node.Rollback(ctx, tracker, processes, refs)
		}
	}()

	if node.Revert != nil {
		err = node.Revert.Do(ctx, refs)
		if err != nil {
			processes.Fatal(err)
			return
		}
	}

	node.logger.Debug("Marking node as completed: ", node.Name)
	tracker.Mark(node)

	if node.AfterRevert != nil {
		ctx, err = node.AfterRevert(ctx, node, tracker, processes, refs)
		if err != nil {
			node.logger.Error("Node after middleware failed: ", err)
			processes.Fatal(err)
			return
		}
	}
}

// Walk iterates over all nodes and returns the lose ends nodes
func (node *Node) Walk(result map[string]*Node, fn func(node *Node)) {
	fn(node)

	if len(node.Next) == 0 {
		result[node.Name] = node
	}

	for _, next := range node.Next {
		next.Walk(result, fn)
	}
}
