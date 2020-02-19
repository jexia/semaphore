package flow

import (
	"testing"

	"github.com/jexia/maestro/specs"
)

func TestConstructBranches(t *testing.T) {
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

	ConstructBranches(nodes)
	t.Log(ConstructSeeds(nodes))
}
