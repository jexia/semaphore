package flow

import (
	"context"

	"github.com/jexia/maestro/codec"
	"github.com/jexia/maestro/refs"
	"github.com/jexia/maestro/services"
	"github.com/jexia/maestro/specs"
	"github.com/jexia/maestro/specs/strict"
)

// NewNode constructs a new node for the given call
func NewNode(call *specs.Call, services services.Collection) *Node {
	service := services.Get(strict.GetService(call.GetEndpoint()))
	return &Node{
		Name:       call.GetName(),
		Previous:   []*Node{},
		Codec:      service.Codec,
		Call:       service.Call,
		Rollback:   service.Rollback,
		DependsOn:  call.DependsOn,
		References: refs.References(call.GetRequest()),
		Next:       []*Node{},
	}
}

// Nodes represents a node colleciton
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
	Name       string
	Previous   Nodes
	Call       services.Call
	Rollback   services.Call
	Codec      codec.Manager
	DependsOn  map[string]*specs.Call
	References map[string]*specs.PropertyReference
	Next       Nodes
}

// Do executes the given node an calls the next nodes, if the call or process failed the next nodes are note called.
func (node *Node) Do(ctx context.Context, tracker *Tracker, processes *Processes, refs *refs.Store) {
	defer processes.Done()

	if !tracker.Met(node.Previous...) {
		return
	}

	if tracker.Met(node) {
		return
	}

	if node.Call != nil {
		err := node.Execute(ctx, node.Call, refs)
		if err != nil {
			processes.Fatal(err)
			return
		}
	}

	tracker.Mark(node)

	if processes.Err() != nil {
		return
	}

	processes.Add(len(node.Next))
	for _, next := range node.Next {
		go next.Do(ctx, tracker, processes, refs)
	}
}

// Revert reverts the given call if a rollback call is defined, once the rollback is completed the previous rollbacks are called.
func (node *Node) Revert(ctx context.Context, tracker *Tracker, processes *Processes, refs *refs.Store) {
	defer processes.Done()

	if !tracker.Met(node.Next...) {
		return
	}

	defer func() {
		processes.Add(len(node.Previous))
		for _, node := range node.Previous {
			go node.Revert(ctx, tracker, processes, refs)
		}
	}()

	if tracker.Met(node) {
		return
	}

	if node.Rollback != nil {
		err := node.Execute(ctx, node.Rollback, refs)
		if err != nil {
			processes.Fatal(err)
			return
		}
	}

	tracker.Mark(node)

	if processes.Err() != nil {
		return
	}
}

// Execute marshals the given reference store into the needed codec and calls the given call.
func (node *Node) Execute(ctx context.Context, caller services.Call, refs *refs.Store) error {
	reader, err := node.Codec.Marshal(refs)
	if err != nil {
		return err
	}

	reader, err = caller(ctx, reader)
	if err != nil {
		return err
	}

	err = node.Codec.Unmarshal(reader, refs)
	if err != nil {
		return err
	}

	return nil
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
