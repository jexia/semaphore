package specs

import (
	"context"
	"testing"

	"github.com/jexia/maestro/logger"
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
		ctx := context.Background()
		ctx = logger.WithValue(ctx)

		err := CheckManifestDuplicates(ctx, input)
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
		ctx := context.Background()
		ctx = logger.WithValue(ctx)

		err := CheckFlowDuplicates(ctx, input)
		if err == nil {
			t.Fatal("unexpected pass", input)
		}
	}
}
