package flow

import (
	"testing"

	"github.com/jexia/maestro/specs"
)

func NewMockNodes() []*Node {
	nodes := []*Node{
		NewMockNode("first", nil, nil),
		NewMockNode("second", nil, nil),
		NewMockNode("third", nil, nil),
	}

	nodes[1].References["first"] = &specs.PropertyReference{
		Resource: "first",
	}

	nodes[2].References["second"] = &specs.PropertyReference{
		Resource: "second",
	}

	return nodes
}

func TestConstructBranches(t *testing.T) {
	nodes := NewMockNodes()
	ConstructBranches(nodes)

	if len(nodes[0].Next) != 1 {
		t.Errorf("unexpected next node %+v", nodes[0])
	}

	if nodes[0].Next[0] != nodes[1] {
		t.Errorf("unexpected previous node %+v, expected %+v", nodes[0].Next[0], nodes[1])
	}

	if len(nodes[1].Previous) != 1 {
		t.Errorf("unexpected previous node %+v", nodes[1])
	}

	if nodes[1].Previous[0] != nodes[0] {
		t.Errorf("unexpected previous node %+v, expected %+v", nodes[1].Previous[0], nodes[0])
	}

	if len(nodes[1].Next) != 1 {
		t.Errorf("unexpected next node %+v", nodes[1])
	}

	if nodes[1].Next[0] != nodes[2] {
		t.Errorf("unexpected previous node %+v, expected %+v", nodes[1].Next[0], nodes[2])
	}

	if nodes[2].Previous[0] != nodes[1] {
		t.Errorf("unexpected previous node %+v, expected %+v", nodes[2].Previous[0], nodes[1])
	}

	if len(nodes[2].Previous) != 1 {
		t.Errorf("unexpected previous node %+v", nodes[1])
	}
}

func TestStartNodes(t *testing.T) {
	last := NewMockNode("last", nil, nil)
	nodes := append(NewMockNodes(), last)
	ConstructBranches(nodes)

	start := FetchStarting(nodes)
	expected := []*Node{nodes[0], last}

	if len(start) != 2 {
		t.Errorf("unexpected ammount of start nodes returned %+v", start)
	}

lookup:
	for _, start := range start {
		for _, expected := range expected {
			if start == expected {
				break lookup
			}
		}

		t.Errorf("unexpected start node %+v", start)
	}
}
