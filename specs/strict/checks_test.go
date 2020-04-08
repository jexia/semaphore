package strict

import (
	"testing"

	"github.com/jexia/maestro/internal/instance"
	"github.com/jexia/maestro/specs"
)

func TestDuplicateManifests(t *testing.T) {
	tests := []*specs.Manifest{
		{
			Flows: []*specs.Flow{
				{
					Name: "dup",
				},
				{
					Name: "dup",
				},
			},
		},
		{
			Flows: []*specs.Flow{
				{
					Name: "first",
					Nodes: []*specs.Node{
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
		ctx := instance.NewContext()
		err := CheckManifestDuplicates(ctx, input)
		if err == nil {
			t.Fatal("unexpected pass", input)
		}
	}
}

func TestDuplicateFlow(t *testing.T) {
	tests := []*specs.Flow{
		{
			Name: "first",
			Nodes: []*specs.Node{
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
		ctx := instance.NewContext()
		err := CheckFlowDuplicates(ctx, input)
		if err == nil {
			t.Fatal("unexpected pass", input)
		}
	}
}
