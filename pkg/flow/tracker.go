package flow

import (
	"sync"
)

// Tracker represents a structure responsible of tracking nodes
type Tracker interface {
	// Flow returns the flow name of the assigned tracker
	Flow() string
	// Mark marks the given node as called
	Mark(node *Node)
	// Skip marks the given node as marked and flag the given node as skipped
	Skip(node *Node)
	// Skipped returns a boolean representing whether the given node has been skipped
	Skipped(node *Node) bool
	// Reached checks whether the required dependencies counter have been reached
	Reached(node *Node, nodes int) bool
	// Met checks whether the given nodes have been called
	Met(nodes ...*Node) bool
	// Lock locks the given node
	Lock(node *Node)
	// Unlock unlocks the given node
	Unlock(node *Node)
}

// NewTracker constructs a new tracker
func NewTracker(flow string, nodes int) Tracker {
	return &tracker{
		flow:  flow,
		nodes: make(map[string]int, nodes),
		locks: make(map[*Node]*sync.Mutex, nodes),
	}
}

// Tracker represents a structure responsible of tracking nodes
type tracker struct {
	flow  string
	mutex sync.Mutex
	nodes map[string]int
	locks map[*Node]*sync.Mutex
}

// Flow returns the flow name of the assigned tracker
func (tracker *tracker) Flow() string {
	return tracker.flow
}

// Mark marks the given node as called
func (tracker *tracker) Mark(node *Node) {
	tracker.mutex.Lock()
	tracker.nodes[node.Name]++
	tracker.mutex.Unlock()
}

// Skip marks the given node as marked and flag the given node as skipped
func (tracker *tracker) Skip(node *Node) {
	tracker.mutex.Lock()
	tracker.nodes[node.Name] = -1
	tracker.mutex.Unlock()
}

// Skip marks the given node as marked and flag the given node as skipped
func (tracker *tracker) Skipped(node *Node) bool {
	tracker.mutex.Lock()
	value := tracker.nodes[node.Name]
	tracker.mutex.Unlock()
	return value < 0
}

// Reached checks whether the required dependencies counter have been reached
func (tracker *tracker) Reached(node *Node, nodes int) bool {
	tracker.mutex.Lock()
	defer tracker.mutex.Unlock()
	if tracker.nodes[node.Name] != nodes {
		return false
	}

	return true
}

// Met checks whether the given nodes have been called
func (tracker *tracker) Met(nodes ...*Node) bool {
	tracker.mutex.Lock()
	defer tracker.mutex.Unlock()
	for _, node := range nodes {
		value := tracker.nodes[node.Name]
		if value <= 0 {
			return false
		}
	}
	return true
}

// Lock locks the given node
func (tracker *tracker) Lock(node *Node) {
	tracker.mutex.Lock()
	mutex := tracker.locks[node]
	if mutex == nil {
		mutex = &sync.Mutex{}
		tracker.locks[node] = mutex
	}
	tracker.mutex.Unlock()
	mutex.Lock()
}

// Unlock unlocks the given node
func (tracker *tracker) Unlock(node *Node) {
	tracker.mutex.Lock()
	mutex := tracker.locks[node]
	tracker.mutex.Unlock()
	mutex.Unlock()
}
