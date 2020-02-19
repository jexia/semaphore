package flow

import "testing"

func TestTrackerMark(t *testing.T) {
	tracker := NewTracker(1)
	node := NewMockNode("first", nil, nil)

	if tracker.Met(node) {
		t.Errorf("unexpected result, tracker met node before marked")
	}

	tracker.Mark(node)

	if !tracker.Met(node) {
		t.Errorf("unexpected result, tracker dit not met node after marked")
	}
}
