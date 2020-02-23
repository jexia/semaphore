package flow

import (
	"sync"
)

// NewTracker constructs a new tracker
func NewTracker(nodes int) *Tracker {
	return &Tracker{
		Nodes: make(map[string]struct{}, nodes),
	}
}

// Tracker represents a structure responsible of tracking nodes
type Tracker struct {
	mutex sync.Mutex
	Nodes map[string]struct{}
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
