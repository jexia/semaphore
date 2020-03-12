package flow

import (
	"sync"
)

// NewTracker constructs a new tracker
func NewTracker(nodes int) *Tracker {
	return &Tracker{
		Nodes: make(map[string]struct{}, nodes),
		Locks: make(map[*Node]*sync.Mutex, nodes),
	}
}

// Tracker represents a structure responsible of tracking nodes
type Tracker struct {
	mutex sync.Mutex
	Nodes map[string]struct{}
	Locks map[*Node]*sync.Mutex
}

// Mark marks the given node as called
func (tracker *Tracker) Mark(node *Node) {
	tracker.mutex.Lock()
	tracker.Nodes[node.Name] = struct{}{}
	tracker.mutex.Unlock()
}

// Met checks whether the given nodes have been called
func (tracker *Tracker) Met(nodes ...*Node) bool {
	tracker.mutex.Lock()
	defer tracker.mutex.Unlock()
	for _, node := range nodes {
		_, has := tracker.Nodes[node.Name]
		if !has {
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
	if mutex == nil {
		mutex = &sync.Mutex{}
		tracker.Locks[node] = mutex
	}
	tracker.mutex.Unlock()
	mutex.Unlock()
}
