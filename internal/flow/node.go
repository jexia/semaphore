package flow

import (
	"context"

	"github.com/jexia/maestro/internal/instance"
	"github.com/jexia/maestro/internal/logger"
	"github.com/jexia/maestro/specs"
	"github.com/sirupsen/logrus"
)

// NewNode constructs a new node for the given call.
// The service called inside the call endpoint is retrieved from the services collection.
// The call, codec and rollback are defined inside the node and used while processing requests.
func NewNode(ctx instance.Context, node *specs.Node, call, rollback Call) *Node {
	references := specs.References{}

	if node.Call != nil {
		references.MergeLeft(specs.ParameterReferences(node.Call.Request))
	}

	if call != nil {
		for _, prop := range call.References() {
			references.MergeLeft(specs.PropertyReferences(prop))
		}
	}

	if rollback != nil {
		for _, prop := range rollback.References() {
			references.MergeLeft(specs.PropertyReferences(prop))
		}
	}

	logger := ctx.Logger(logger.Flow)

	return &Node{
		ctx:        ctx,
		logger:     logger,
		Name:       node.Name,
		Previous:   []*Node{},
		Call:       call,
		Rollback:   rollback,
		DependsOn:  node.DependsOn,
		References: references,
		Next:       []*Node{},
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

// Node represents a collection of callers and rollbacks which could be executed parallel.
type Node struct {
	ctx        instance.Context
	logger     *logrus.Logger
	Name       string
	Previous   Nodes
	Call       Call
	Rollback   Call
	DependsOn  map[string]*specs.Node
	References map[string]*specs.PropertyReference
	Next       Nodes
}

// Do executes the given node an calls the next nodes.
// If one of the nodes fails is the error marked and are the processes aborted.
func (node *Node) Do(ctx context.Context, tracker *Tracker, processes *Processes, refs specs.Store) {
	defer processes.Done()
	node.logger.WithField("node", node.Name).Debug("Executing node call")

	if !tracker.Met(node.Previous...) {
		node.logger.WithField("node", node.Name).Debug("Has not met dependencies yet")
		return
	}

	tracker.Lock(node)
	if tracker.Met(node) {
		node.logger.WithField("node", node.Name).Debug("Node already executed")
		return
	}

	if node.Call != nil {
		err := node.Call.Do(ctx, refs)
		if err != nil {
			node.logger.WithField("node", node.Name).Error("Call failed")
			processes.Fatal(err)
			return
		}
	}

	node.logger.WithField("node", node.Name).Debug("Marking node as completed")

	tracker.Mark(node)
	tracker.Unlock(node)

	if processes.Err() != nil {
		node.logger.WithField("node", node.Name).Error("Stopping flow execution a error has been thrown")
		return
	}

	processes.Add(len(node.Next))
	for _, next := range node.Next {
		go next.Do(ctx, tracker, processes, refs)
	}
}

// Revert executes the given node rollback an calls the previous nodes.
// If one of the nodes fails is the error marked but execution is not aborted.
func (node *Node) Revert(ctx context.Context, tracker *Tracker, processes *Processes, refs specs.Store) {
	defer processes.Done()
	node.logger.WithField("node", node.Name).Debug("Executing node revert")

	if !tracker.Met(node.Next...) {
		node.ctx.Logger(logger.Flow).WithField("node", node.Name).Debug("Has not met dependencies yet")
		return
	}

	defer func() {
		processes.Add(len(node.Previous))
		for _, node := range node.Previous {
			go node.Revert(ctx, tracker, processes, refs)
		}
	}()

	tracker.Lock(node)
	if tracker.Met(node) {
		node.logger.WithField("node", node.Name).Debug("Node already executed")
		return
	}

	if node.Rollback != nil {
		err := node.Rollback.Do(ctx, refs)
		if err != nil {
			processes.Fatal(err)
			return
		}
	}

	node.logger.WithField("node", node.Name).Debug("Marking node as completed")

	tracker.Mark(node)
	tracker.Unlock(node)
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
