package flow

import (
	"sync"
)

// NewTracker constructs a new tracker
func NewTracker(nodes int) *Tracker {
	return &Tracker{
		Nodes: make(map[string]int, nodes),
		Locks: make(map[*Node]*sync.Mutex, nodes),
	}
}

// Tracker represents a structure responsible of tracking nodes
type Tracker struct {
	mutex sync.Mutex
	Nodes map[string]int
	Locks map[*Node]*sync.Mutex
}

// Mark marks the given node as called
func (tracker *Tracker) Mark(node *Node) {
	tracker.mutex.Lock()
	tracker.Nodes[node.Name]++
	tracker.mutex.Unlock()
}

// Reached checks whether the required dependencies counter have been reached
func (tracker *Tracker) Reached(node *Node, nodes int) bool {
	tracker.mutex.Lock()
	defer tracker.mutex.Unlock()
	if tracker.Nodes[node.Name] != nodes {
		return false
	}

	return true
}

// Met checks whether the given nodes have been called
func (tracker *Tracker) Met(nodes ...*Node) bool {
	tracker.mutex.Lock()
	defer tracker.mutex.Unlock()
	for _, node := range nodes {
		value := tracker.Nodes[node.Name]
		if value == 0 {
			return false
		}
	}
	return true
}

// Lock locks the given node
func (tracker *Tracker) Lock(node *Node) {
	tracker.mutex.Lock()
	mutex := tracker.Locks[node]
	if mutex == nil {
		mutex = &sync.Mutex{}
		tracker.Locks[node] = mutex
	}
	tracker.mutex.Unlock()
	mutex.Lock()
}

// Unlock unlocks the given node
func (tracker *Tracker) Unlock(node *Node) {
	tracker.mutex.Lock()
	mutex := tracker.Locks[node]
	tracker.mutex.Unlock()
	mutex.Unlock()
}
