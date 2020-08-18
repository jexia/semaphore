package flow

import "testing"

func TestTrackerMark(t *testing.T) {
	tracker := NewTracker("", 1)
	node := NewMockNode("first", nil, nil)

	if tracker.Met(node) {
		t.Errorf("unexpected result, tracker met node before marked")
	}

	tracker.Mark(node)

	if !tracker.Met(node) {
		t.Errorf("unexpected result, tracker did not met node after marked")
	}
}

func TestTrackerName(t *testing.T) {
	expected := "MockFlow"
	tracker := NewTracker(expected, 1)

	flow := tracker.Flow()
	if flow != expected {
		t.Fatalf("unexpected flow name %s, expected %s", flow, expected)
	}
}
