package specs

import (
	"testing"
)

func TestDuplicateManifests(t *testing.T) {
	tests := []*Manifest{
		{
			Flows: []*Flow{
				{
					Name: "dup",
				},
				{
					Name: "dup",
				},
			},
		},
		{
			Flows: []*Flow{
				{
					Name: "first",
					Nodes: []*Node{
						{
							Name: "dup",
						},
						{
							Name: "dup",
						},
					},
				},
			},
		},
	}

	for _, input := range tests {
		err := CheckManifestDuplicates(input)
		if err == nil {
			t.Fatal("unexpected pass", input)
		}
	}
}

func TestDuplicateFlow(t *testing.T) {
	tests := []*Flow{
		{
			Name: "first",
			Nodes: []*Node{
				{
					Name: "dup",
				},
				{
					Name: "dup",
				},
			},
		},
	}

	for _, input := range tests {
		err := CheckFlowDuplicates(input)
		if err == nil {
			t.Fatal("unexpected pass", input)
		}
	}
}
